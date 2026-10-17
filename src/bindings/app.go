// Binding for JS to Go
// This package is responsible for the communication between the JS and Go code.
package bindings

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"runtime"
	"strings"

	"LibreRemotePlay/src/cli"
	"LibreRemotePlay/src/devices/gamepad"
	"LibreRemotePlay/src/devices/keyboard"
	"LibreRemotePlay/src/devices/mouse"
	net_webrtc "LibreRemotePlay/src/net/webrtc"

	"LibreRemotePlay/src/net/webrtc/streaming_signal"
	"LibreRemotePlay/src/state"

	"github.com/pion/webrtc/v4"
)

type Res[T any] struct {
	Result T `json:"result"`
}

type Req[T any] struct {
	Args T `json:"args"`
}

func SetupBindings(app state.App, mux *http.ServeMux) {

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Bindings OK"))
	})

	mux.Handle("PUT /notify/peer/open", notifyOpenPeer(app))
	mux.Handle("PUT /notify/peer/close", notifyClosePeer(app))

	mux.Handle("POST /conn/host", tryCreateHost(app))
	mux.Handle("DELETE /conn/host", tryClosePeerConnection(app))

	mux.Handle("GET /ports", getUsedPorts())
	
	mux.Handle("GET /gamepad", isGamepadEnabled())
	mux.Handle("PUT /gamepad", toogleGamepad())

	mux.Handle("GET /keyboard", isKeyboardEnabled())
	mux.Handle("PUT /keyboard", toogleKeyboard())

	mux.Handle("GET /mouse", isMouseEnabled())
	mux.Handle("PUT /mouse", toogleMouse())

	mux.Handle("GET /whip", isWhipEnabled())
	mux.Handle("PUT /whip", toogleWhip())

	mux.Handle("GET /os", getCurrentOS())

	mux.Handle("POST /log", logPrintln())

	mux.Handle("GET /vigem", openViGEmWizard())
	
}

func notifyOpenPeer(app state.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		app.Conn.Peer.Open()
		log.Println("NotifyCreateClient")
		w.WriteHeader(200)
	}
}

func notifyClosePeer(app state.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		app.Conn.Peer.Close()
		log.Println("NotifyCloseClient")
		w.WriteHeader(200)
	}
}

// Create a Host Peer, it receives the offer encoded and returns the encoded answer response
func tryCreateHost(app state.App) http.HandlerFunc {

	type Args struct {
		ICEServers []webrtc.ICEServer `json:"ice_servers"`
		OfferEncoded string `json:"offer_encoded"`
	}
	
	return func(w http.ResponseWriter, r *http.Request) {
		
		if app.Conn.Peer.IsOpen() {
			app.Conn.CancelCtx()
		}

		app.Conn.Peer.Open()
		
		answerResponse := make(chan string, 1) // Buffer of 1 to not block thread on cancel
	
		connCtx, cancelConn := context.WithCancel(app.Ctx)
		
		app.Conn.CancelCtx = cancelConn
		app.Conn.Ctx = connCtx

		rawReq, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println("Error reading req body: ", err)
			w.WriteHeader(500)
			return
		}

		var req Req[Args]
		err = json.Unmarshal(rawReq, &req)
		if err != nil {
			log.Println("Error unmarshalling req body: ", err)
			w.WriteHeader(400)
			return 
		}
		
		log.Println(req.Args.ICEServers)
		streaming_signal.WhipConfig.ICEServers.Store(&req.Args.ICEServers)

		pidAudioChan := make(chan uint32)
		
		go net_webrtc.InitHost(connCtx, req.Args.ICEServers, req.Args.OfferEncoded, answerResponse, pidAudioChan)
	
		answer := <-answerResponse
	
		if answer == net_webrtc.ERROR_ANSWER {
			app.Conn.Peer.Close()
			log.Println("Error on WebRTC host connection")
			w.WriteHeader(500)
		}

		res := Res[string]{
			Result: answer,
		}

		resRaw, err := json.Marshal(res)
		if err != nil {
			log.Println("Error marshaling response: ", err)
			w.WriteHeader(500)
			return
		}
		
		w.Write(resRaw)
	}
}

// Closes the peer connection and returns a boolean indication if a connection existed and was closed or not
func tryClosePeerConnection(app state.App) http.HandlerFunc {
	
	return func(w http.ResponseWriter, r *http.Request) {
		
		if !app.Conn.Peer.IsOpen() {
			w.WriteHeader(409)
			return
		}
	
		app.Conn.CancelCtx()
	
		app.Conn.Peer.Close()
		
		w.WriteHeader(200)
	}

}

func getUsedPorts() http.HandlerFunc {
	
	config := cli.GetConfig()

	type UsedPorts struct {
		HTTP uint16 `json:"http"`
		EasyConnect uint16 `json:"easy_connect"`
		WHIP uint16 `json:"whip"`
	}

	usedPorts := UsedPorts{
		HTTP: config.GetHTTPPort(),
		EasyConnect: config.GetEasyConnectPort(),
		WHIP: config.GetWhipServerPort(),
	}
	
	return func(w http.ResponseWriter, r *http.Request) {

		res := Res[UsedPorts]{
			Result: usedPorts,
		}

		rawRes, err := json.Marshal(res)
		if err != nil {
			log.Println("Error mashalling res: ", err)
			w.WriteHeader(500)
			return
		}

		w.Write(rawRes)
	} 
}

func toogleHelper(toogle func()) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		toogle()
		w.WriteHeader(200)
	}
}

func isEnabledHelper(isEnabled func() bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := Res[bool]{
			Result: isEnabled(),
		}

		rawRes, err := json.Marshal(res)
		if err != nil {
			log.Println("Error mashalling res: ", err)
			w.WriteHeader(500)
			return
		}

		w.Write(rawRes)
	}
}

// Gamepad

func toogleGamepad() http.HandlerFunc {
	return toogleHelper(gamepad.GamepadEnabled.Toogle)
}

func isGamepadEnabled() http.HandlerFunc {
	return isEnabledHelper(gamepad.GamepadEnabled.IsEnabled)
}

// Keyboard

func toogleKeyboard() http.HandlerFunc {
	return toogleHelper(keyboard.KeyboardEnabled.Toogle)
}

func isKeyboardEnabled() http.HandlerFunc {
	return isEnabledHelper(keyboard.KeyboardEnabled.IsEnabled)
}

// Mouse

func toogleMouse() http.HandlerFunc {
	return toogleHelper(mouse.MouseEnabled.Toogle)
}

func isMouseEnabled() http.HandlerFunc {
	return isEnabledHelper(mouse.MouseEnabled.IsEnabled)
}

// WHIP

func toogleWhip() http.HandlerFunc {
	return toogleHelper(streaming_signal.WhipConfig.Enabled.Toogle)
}

func isWhipEnabled() http.HandlerFunc {
	return isEnabledHelper(streaming_signal.WhipConfig.Enabled.IsEnabled)
}

// OS Detection

func getCurrentOS() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		res := Res[string]{
			Result: strings.ToUpper(runtime.GOOS),
		}

		rawRes, err := json.Marshal(res)
		if err != nil {
			log.Println("Error mashalling res: ", err)
			w.WriteHeader(500)
			return
		}

		w.Write(rawRes)
	}
}

// Log

func logPrintln() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		rawBody, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println("Error reading req body: ", err)
			w.WriteHeader(500)
			return
		}

		var req Req[string]
		err = json.Unmarshal(rawBody, &req)
		if err != nil {
			log.Println("Error unmarshalling req: ", err)
			w.WriteHeader(400)
			return
		}

		log.Println(req.Args)

		w.WriteHeader(200)
	}
}

// // Audio
// 
// func  SetAudioPid(pid uint32) {
// 	a.pidAudioChan <- pid
// }
// 
// func  GetAudioProcess() []audio.AudioProcess {
// 	return audio.GetAudioProcess()
// }
