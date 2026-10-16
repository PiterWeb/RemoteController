package streaming_signal

import (
	"context"
	"log"

	"github.com/pion/webrtc/v4"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func HandleStreamingSignal(ctx context.Context, streamingSignalChannel *webrtc.DataChannel) {

	if streamingSignalChannel.Label() != "streaming-signal" {
		return
	}

	go handleWhipOffer(ctx, streamingSignalChannel)

	streamingSignalChannel.OnMessage(func(msg webrtc.DataChannelMessage) {

		if WhipConfig.Enabled.IsEnabled() {
			handleWhipAnswer(msg.Data)
		} else {
			runtime.EventsEmit(ctx, "streaming-signal-client", string(msg.Data))
		}

	})

	runtime.EventsOn(ctx, "streaming-signal-server", func(data ...any) {

		if WhipConfig.Enabled.IsEnabled() {
			return
		}

		if len(data) == 0 {
			return
		}

		signalingData, ok := data[0].(string)

		if !ok {
			log.Println(data[0], ok)
			return
		}

		_ = streamingSignalChannel.SendText(signalingData)

	})

}
