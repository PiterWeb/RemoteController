package streaming_signal

import (
	"context"

	"github.com/pion/webrtc/v3"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func HandleStreamingSignal(ctx context.Context, streamingSignalChannel *webrtc.DataChannel) {

	if streamingSignalChannel.Label() != "streaming-signal" {
		return
	}

	streamingSignalChannel.OnMessage(func(msg webrtc.DataChannelMessage) {

		runtime.EventsEmit(ctx, "streaming-signal", msg.Data)

	})

}
