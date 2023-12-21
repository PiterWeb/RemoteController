package streaming_signal

import "github.com/pion/webrtc/v3"

func HandleStreamingSignal(streamingSignalChannel *webrtc.DataChannel) {

	if streamingSignalChannel.Label() != "streaming-signal" {
		return
	}

}
