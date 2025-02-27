package streaming_signal

import (
	"context"
	"log"

	"github.com/coder/websocket"
	"github.com/pion/webrtc/v3"
)

func HandleStreamingSignal(ctx context.Context, streamingSignalChannel *webrtc.DataChannel) {

	if streamingSignalChannel.Label() != "streaming-signal" {
		return
	}

	wsClient, _, err := websocket.Dial(context.Background(), "ws://localhost:8080/ws", nil)

	if err != nil {
		log.Println(err)
		return
	}

	defer func() {
		if err := recover(); err != nil {
			wsClient.Close(websocket.StatusInternalError, "Fatal error on client")
		}
	}()

	go func() {

		defer wsClient.Close(websocket.StatusInternalError, "Client terminated")

		for {
			t, data, err := wsClient.Read(context.Background())

			if err != nil {
				log.Println(err)
				continue
			}

			if t != websocket.MessageText {
				continue
			}

			streamingSignalChannel.SendText(string(data))

		}
	}()

	streamingSignalChannel.OnMessage(func(msg webrtc.DataChannelMessage) {

		wsClient.Write(context.Background(), websocket.MessageText, msg.Data)

	})

}
