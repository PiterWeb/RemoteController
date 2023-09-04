package net

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/PiterWeb/RemoteController/src/gamepad"
	"github.com/pion/webrtc/v3"
)

func InitOffer(offerChan chan<- string, answerResponse <-chan string) {

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
		if cErr := peerConnection.Close(); cErr != nil {
			fmt.Printf("cannot close peerConnection: %v\n", cErr)
		}
	}()

	// Create a datachannel with label 'data'
	dataChannel, err := peerConnection.CreateDataChannel("data", nil)
	if err != nil {
		panic(err)
	}

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

	// Register channel opening handling
	dataChannel.OnOpen(func() {

		gamepads := gamepad.All{}

		if err != nil {
			panic(err)
		}

		for range time.NewTicker(1 * time.Millisecond).C {
			gamepads.Update()
			for i := range gamepads {
				pad := &gamepads[i]

				if !pad.Connected {
					continue
				}

				padBytes, _ := json.Marshal(*pad)

				err := dataChannel.Send(padBytes)

				if err != nil {
					panic(err)
				}

			}
		}

	})

	// Register text message handling
	dataChannel.OnMessage(func(msg webrtc.DataChannelMessage) {
		fmt.Printf("Message from DataChannel '%s': '%s'\n", dataChannel.Label(), string(msg.Data))
	})

	// Create an offer to send to the other process
	offer, err := peerConnection.CreateOffer(nil)
	if err != nil {
		panic(err)
	}

	gatherComplete := webrtc.GatheringCompletePromise(peerConnection)

	// Sets the LocalDescription, and starts our UDP listeners
	// Note: this will start the gathering of ICE candidates
	if err = peerConnection.SetLocalDescription(offer); err != nil {
		panic(err)
	}

	// Block until ICE Gathering is complete, disabling trickle ICE
	// we do this because we only can exchange one signaling message
	// in a production application you should exchange ICE Candidates via OnICECandidate
	<-gatherComplete

	offerChan <- signalEncode(offer)

	answer := webrtc.SessionDescription{}

	signalDecode(<-answerResponse, &answer)

	if err = peerConnection.SetRemoteDescription(answer); err != nil {
		panic(err)
	}
	// Block forever
	select {}
}
