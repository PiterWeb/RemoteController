package net

import (
	"context"
	"fmt"
	"sync"

	"github.com/PiterWeb/RemoteController/src/gamepad"
	"github.com/pion/webrtc/v3"
	"github.com/pquerna/ffjson/ffjson"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func InitHost(ctx context.Context, offerEncoded string, answerResponse chan<- string, triggerEnd <-chan struct{}) {

	var candidatesMux sync.Mutex
	candidates := []string{}

	// Prepare the configuration
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	peerConnection, err := webrtc.NewAPI().NewPeerConnection(config)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := peerConnection.Close(); err != nil {
			fmt.Printf("cannot close peerConnection: %v\n", err)
		}
	}()

	// Create a datachannel with label 'controller'
	streamingChannel, err := peerConnection.CreateDataChannel("streaming", nil)
	if err != nil {
		panic(err)
	}

	peerConnection.OnICECandidate(func(c *webrtc.ICECandidate) {

		if c == nil {
			return
		}

		candidatesMux.Lock()
		defer candidatesMux.Unlock()

		desc := peerConnection.RemoteDescription()

		if desc != nil {
			candidates = append(candidates, (*c).ToJSON().Candidate)
		}

	})

	// Set the handler for Peer connection state
	// This will notify you when the peer has connected/disconnected
	peerConnection.OnConnectionStateChange(func(s webrtc.PeerConnectionState) {
		fmt.Printf("Peer Connection State has changed: %s\n", s.String())

		if s == webrtc.PeerConnectionStateFailed {
			if closeErr := peerConnection.Close(); closeErr != nil {
				panic(closeErr)
			}
		}
	})

	// Register data channel creation handling
	peerConnection.OnDataChannel(func(d *webrtc.DataChannel) {

		if d.Label() == "controller" {

			var virtualDevice gamepad.EmulatedDevice
			defer gamepad.FreeTargetAndDisconnect(virtualDevice)

			virtualState := new(gamepad.ViGEmState)

			// Create a virtual device
			d.OnOpen(func() {

				virtualDevice, err = gamepad.GenerateVirtualDevice()

				if err != nil {
					panic(err)
				}

			})

			// Update the virtual device
			d.OnMessage(func(msg webrtc.DataChannelMessage) {

				var pad gamepad.State

				ffjson.Unmarshal(msg.Data, &pad)

				go gamepad.UpdateVirtualDevice(virtualDevice, pad, virtualState)

			})
		}

	})

	streamingChannel.OnOpen(func() {

		fmt.Println("Streaming channel openned")

		runtime.EventsOn(ctx, "send-streaming", func(optionalData ...interface{}) {
			fmt.Println("sending streaming")

			err := streamingChannel.SendText(optionalData[0].(string))

			if err != nil {
				fmt.Println(err)
			}
		})

	})

	offer := webrtc.SessionDescription{}
	signalDecode(offerEncoded, &offer)

	if err := peerConnection.SetRemoteDescription(offer); err != nil {
		panic(err)
	}

	// Create an answer to send to the other process
	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		panic(err)
	}

	gatherComplete := webrtc.GatheringCompletePromise(peerConnection)

	// Sets the LocalDescription, and starts our UDP listeners
	err = peerConnection.SetLocalDescription(answer)
	if err != nil {
		panic(err)
	}

	<-gatherComplete

	answerResponse <- signalEncode(*peerConnection.LocalDescription()) + ";" + signalEncode(candidates)

	// Block until cancel by user
	<-triggerEnd

}
