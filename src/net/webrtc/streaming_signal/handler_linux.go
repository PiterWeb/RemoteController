package streaming_signal

import (
	"context"
	"fmt"
	"log"

	"LibreRemotePlay/src/cli"
	"github.com/coder/websocket"
	"github.com/pion/webrtc/v4"
)

func HandleStreamingSignal(ctx context.Context, streamingSignalChannel *webrtc.DataChannel) {

	if streamingSignalChannel.Label() != "streaming-signal" {
		return
	}

	wsClient, _, err := websocket.Dial(ctx, fmt.Sprintf("ws://localhost:%d/ws", cli.GetConfig().GetHTTPPort()), nil)

	if err != nil {
		log.Println(err)
		return
	}

	defer func() {
		if err := recover(); err != nil {
			wsClient.Close(websocket.StatusInternalError, "Fatal error on client")
		}
	}()

	go handleWhipOffer(ctx, streamingSignalChannel)

	go func() {

		defer wsClient.Close(websocket.StatusInternalError, "Client terminated")

		for {
			select {
				case <-ctx.Done():
					return
				default:
					t, data, err := wsClient.Read(ctx)
	
					if err != nil {
						
						log.Println(err)
						continue
					}
	
					if WhipConfig.Enabled.IsEnabled() {
						continue
					}
	
					if t != websocket.MessageText {
						continue
					}
	
					streamingSignalChannel.SendText(string(data))
			}
		}
	}()

	streamingSignalChannel.OnMessage(func(msg webrtc.DataChannelMessage) {

		if WhipConfig.Enabled.IsEnabled() {
			handleWhipAnswer(msg.Data)
		} else {
			wsClient.Write(ctx, websocket.MessageText, msg.Data)
		}

	})

}
