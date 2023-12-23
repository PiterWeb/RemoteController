package streaming_signal

import (
	"context"
	"fmt"

	"github.com/pion/webrtc/v3"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func HandleStreamingSignal(ctx context.Context, streamingSignalChannel *webrtc.DataChannel) {

	if streamingSignalChannel.Label() != "streaming-signal" {
		return
	}

	streamingSignalChannel.OnMessage(func(msg webrtc.DataChannelMessage) {

		runtime.EventsEmit(ctx, "streaming-signal-client", string(msg.Data))

	})

	runtime.EventsOn(ctx, "streaming-signal-server", func(data ...interface{}) {

		signalingData, ok := data[0].(string)

		if !ok {
			fmt.Println(data[0], ok)
			return
		}

		streamingSignalChannel.SendText(signalingData)

	})

}
