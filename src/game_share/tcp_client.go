package game_share

import "github.com/pion/webrtc/v3"

func clientTCP(port int, datachannel *webrtc.DataChannel) error {

	datachannel.OnMessage(func(msg webrtc.DataChannelMessage) {

	})

	return nil
}
