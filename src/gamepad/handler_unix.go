//go:build unix

package gamepad

import (
	// "fmt"

	"github.com/pion/webrtc/v3"
	// "github.com/pquerna/ffjson/ffjson"
)

func HandleGamepad(gamepadChannel *webrtc.DataChannel) {

	if gamepadChannel.Label() != "controller" {
		return
	}

	// Create a virtual device
	gamepadChannel.OnOpen(func() {

	})

	// Update the virtual device
	gamepadChannel.OnMessage(func(msg webrtc.DataChannelMessage) {

		// var pad clonedGamepad

		// ffjson.Unmarshal(msg.Data, &pad)

		// fmt.Println(pad)

	})
}
