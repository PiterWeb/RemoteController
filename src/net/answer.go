package net

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/PiterWeb/RemoteController/src/gamepad"
	"github.com/pion/webrtc/v3"
)

func InitAnswer(offerEncoded string, answerResponse chan<- string) {

	// Prepare the configuration
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	// Create a new RTCPeerConnection
	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := peerConnection.Close(); err != nil {
			fmt.Printf("cannot close peerConnection: %v\n", err)
		}
	}()

	// Set the handler for Peer connection state
	// This will notify you when the peer has connected/disconnected
	peerConnection.OnConnectionStateChange(func(s webrtc.PeerConnectionState) {
		fmt.Printf("Peer Connection State has changed: %s\n", s.String())

		if s == webrtc.PeerConnectionStateFailed {
			// Wait until PeerConnection has had no network activity for 30 seconds or another failure. It may be reconnected using an ICE Restart.
			// Use webrtc.PeerConnectionStateDisconnected if you are interested in detecting faster timeout.
			// Note that the PeerConnection may come back from PeerConnectionStateDisconnected.
			fmt.Println("Peer Connection has gone to failed exiting")
			os.Exit(0)
		}
	})

	// Register data channel creation handling
	peerConnection.OnDataChannel(func(d *webrtc.DataChannel) {
		fmt.Printf("New DataChannel %s %d\n", d.Label(), d.ID())

		var virtualDevice gamepad.EmulatedDevice
		defer gamepad.FreeTargetAndDisconnect(virtualDevice)
		
		virtualState := new(gamepad.ViGEmState)

		// Register channel opening handling
		d.OnOpen(func() {

			virtualDevice, err = gamepad.GenerateVirtualDevice()

			if err != nil {
				panic(err)
			}

		})

		// Register text message handling
		d.OnMessage(func(msg webrtc.DataChannelMessage) {
			fmt.Printf("Message from DataChannel '%s': '%s'\n", d.Label(), string(msg.Data))

			var pad gamepad.State

			err = json.Unmarshal(msg.Data, &pad)

			go gamepad.UpdateVirtualDevice(virtualDevice, pad, virtualState)

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

	// Create channel that is blocked until ICE Gathering is complete
	gatherComplete := webrtc.GatheringCompletePromise(peerConnection)

	// Sets the LocalDescription, and starts our UDP listeners
	err = peerConnection.SetLocalDescription(answer)
	if err != nil {
		panic(err)
	}

	<-gatherComplete

	// Output the answer in base64 so we can paste it in browser

	answerResponse <- signalEncode(*peerConnection.LocalDescription())

	// Block forever
	select {}

}
