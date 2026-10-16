package streaming_signal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"sync/atomic"
	"time"

	"LibreRemotePlay/src/cli"
	"LibreRemotePlay/src/devices"
	"github.com/pion/webrtc/v4"
)

type whipConfig struct {
	Port       uint16
	OfferChan  chan string
	AnswerChan chan string
	Enabled    *devices.DeviceEnabled
	ICEServers atomic.Pointer[[]webrtc.ICEServer]
}

var whipConfigEnabled = &devices.DeviceEnabled{}

var WhipConfig = &whipConfig{
	Port:       cli.GetConfig().GetWhipServerPort(),
	OfferChan:  make(chan string),
	AnswerChan: make(chan string),
	Enabled:    whipConfigEnabled.Disable(),
	ICEServers: atomic.Pointer[[]webrtc.ICEServer]{},
}

func InitWhipServer(ctx context.Context, config *whipConfig) error {

	var answerChan <-chan string = config.AnswerChan
	var offerChan chan<- string = config.OfferChan

	defer func() {
		close(config.OfferChan)
		close(config.AnswerChan)
	}()

	httpServerMux := http.NewServeMux()

	httpServerMux.HandleFunc("OPTIONS /whip", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Methods", "POST")
		w.Header().Add("Access-Control-Allow-Headers", "*")
		w.Header().Add("Access-Control-Allow-Headers", "Authorization")
	})
	
	httpServerMux.HandleFunc("POST /whip", func(w http.ResponseWriter, r *http.Request) {
		
		if !WhipConfig.Enabled.IsEnabled() {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("Whip not enabled"))
			return
		}

		defer func() {
			if err := recover(); err != nil {
				log.Printf("Fatal error %s\n", err)
				w.WriteHeader(http.StatusInternalServerError)
				_, _ =w.Write([]byte("Fatal error"))
				return
			}
		}()
		
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Methods", "POST")
		w.Header().Add("Access-Control-Allow-Headers", "*")
		w.Header().Add("Access-Control-Allow-Headers", "Authorization")

		log.Println(*WhipConfig.ICEServers.Load())
		
		for _, s := range *WhipConfig.ICEServers.Load() {
			for _, url := range s.URLs {
				w.Header().Add("Link", fmt.Sprintf("<%s>; rel=\"ice-server\"", url))
			}
		}

		offer, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("Error reading request body"))
			return
		}

		if len(offer) == 0 {
			log.Println("WHIP - Offer received is too short")
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("Error request body too short"))
			return
		}
		
		log.Println("WHIP - Offer received")
		offerChan <- string(offer)

		log.Println("WHIP - Waiting for answer")
		
		rawAnswer, ok := <-answerChan

		if !ok {
			log.Println("WHIP - Answer channel closed")
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("Error getting an answer"))
			return
		}

		type answerT struct {
			Answer struct {
				SDP string
			}
		}

		answerStruct := answerT{}

		err = json.Unmarshal([]byte(rawAnswer), &answerStruct)

		if err != nil {
			log.Printf("WHIP - Error parsing answer: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("Error getting an answer"))
			return
		}

		log.Println("WHIP - Sending answer")

		w.Header().Set("Content-Type", "application/sdp")
		w.Header().Add("Location", fmt.Sprintf("http://localhost:%d/whip", config.Port))

		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(answerStruct.Answer.SDP))

	})

	httpServer := &http.Server{
		Handler:      httpServerMux,
		Addr:         fmt.Sprintf("127.0.0.1:%d", config.Port),
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 30,
		BaseContext: func(l net.Listener) context.Context {
			return ctx
		},
	}

	err := httpServer.ListenAndServe()
	
	if err != nil {
		errors.Join(err, errors.New("whip server error"))
	}
	
	return err

}

func handleWhipOffer(ctx context.Context, streamingSignalChannel *webrtc.DataChannel) {

	for offer := range WhipConfig.OfferChan {
		if WhipConfig.Enabled.IsEnabled() {

			offerMap := map[string]any{}

			ctx, cancelCtx :=  context.WithTimeout(ctx, time.Second)
			
			defer cancelCtx()
			
			finalOffer := getFinalOffer(ctx, offer)		
			
			offerSdp := webrtc.SessionDescription{
				Type: webrtc.SDPTypeOffer, SDP: finalOffer,
			}

			offerMap["offer"] = offerSdp
			offerMap["role"] = "host"
			offerMap["type"] = "offer"

			offerJSONBytes, err := json.Marshal(offerMap)

			if err != nil {
				log.Printf("WHIP - Error encoding offer: %s\n", err)
				continue
			}

			if err := streamingSignalChannel.SendText(string(offerJSONBytes)); err != nil {
				break
			}

		}
	}

}

func getFinalOffer(ctx context.Context, offer string) string {
	tempOffer := offer
	for {
		select {
			case tempOffer = <- WhipConfig.OfferChan:
			continue
			case <-ctx.Done():
			return tempOffer
		}
	}
}

func handleWhipAnswer(msg []byte) {

	type dataT struct {
		Type string
	}

	dataStruct := dataT{}

	err := json.Unmarshal(msg, &dataStruct)

	if err != nil {
		return
	}

	if dataStruct.Type == "answer" {
		WhipConfig.AnswerChan <- string(msg)
	}

}
