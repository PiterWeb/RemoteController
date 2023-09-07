package net

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/PiterWeb/RemoteController/src/gamepad"
	"github.com/pion/webrtc/v3"
)

func InitAnswer(offerEncoded string, answerResponse chan<- string) {

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

	peerConnection.OnICECandidate(func(c *webrtc.ICECandidate) {

		if c == nil {
			return
		}

		candidatesMux.Lock()
		defer candidatesMux.Unlock()

		desc := peerConnection.RemoteDescription()

		if desc != nil {
			log.Println((*c).ToJSON().Candidate)
			candidates = append(candidates, (*c).ToJSON().Candidate)
		}

	})

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
			// fmt.Printf("Message from DataChannel '%s': '%s'\n", d.Label(), string(msg.Data))

			var pad gamepad.State

			err = json.Unmarshal(msg.Data, &pad)

			go gamepad.UpdateVirtualDevice(virtualDevice, pad, virtualState)

		})

	})

	offer := webrtc.SessionDescription{}
	signalDecode(offerEncoded, &offer)

	log.Println("Answer remote description")
	if err := peerConnection.SetRemoteDescription(offer); err != nil {
		panic(err)
	}

	log.Println("Answer created")
	// Create an answer to send to the other process
	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		panic(err)
	}
	
	gatherComplete := webrtc.GatheringCompletePromise(peerConnection)

	log.Println("Answer local description")
	// Sets the LocalDescription, and starts our UDP listeners
	err = peerConnection.SetLocalDescription(answer)
	if err != nil {
		panic(err)
	}

	<-gatherComplete

	answerResponse <- signalEncode(*peerConnection.LocalDescription()) + ";" + signalEncode(candidates)
	
	// Block forever
	select {}

}
