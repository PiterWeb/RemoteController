package net

import (
	"fmt"
	"strings"

	"github.com/PiterWeb/RemoteController/src/gamepad"
	"github.com/pion/webrtc/v3"
)

func InitHost(offerEncodedWithCandidates string, answerResponse chan<- string, triggerEnd <-chan struct{}) {

	candidates := []webrtc.ICECandidateInit{}

	// Prepare the configuration
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19305", "stun:stun.l.google.com:19302", "stun:stun.ipfire.org:3478"},
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

	// Register data channel creation handling
	peerConnection.OnDataChannel(func(d *webrtc.DataChannel) {

		gamepad.HandleGamepad(d)

	})

	peerConnection.OnICECandidate(func(c *webrtc.ICECandidate) {

		if c == nil {
			answerResponse <- signalEncode(*peerConnection.LocalDescription()) + ";" + signalEncode(candidates)
			return
		}

		candidates = append(candidates, (*c).ToJSON())

	})

	// Set the handler for Peer connection state
	// This will notify you when the peer has connected/disconnected
	peerConnection.OnConnectionStateChange(func(s webrtc.PeerConnectionState) {
		fmt.Printf("Peer Connection State has changed: %s\n", s.String())

		if s == webrtc.PeerConnectionStateFailed {
			if err := peerConnection.Close(); err != nil {
				panic(err)
			}
		}
	})

	offerEncodedWithCandidatesSplited := strings.Split(offerEncodedWithCandidates, ";")

	offer := webrtc.SessionDescription{}
	signalDecode(offerEncodedWithCandidatesSplited[0], &offer)

	var receivedCandidates []webrtc.ICECandidateInit

	signalDecode(offerEncodedWithCandidatesSplited[1], &receivedCandidates)

	if err := peerConnection.SetRemoteDescription(offer); err != nil {
		panic(err)
	}

	for _, candidate := range receivedCandidates {
		if err := peerConnection.AddICECandidate(candidate); err != nil {
			panic(err)
		}
	}

	// Create an answer to send to the other process
	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		panic(err)
	}

	// Sets the LocalDescription, and starts our UDP listeners
	err = peerConnection.SetLocalDescription(answer)
	if err != nil {
		panic(err)
	}

	// Block until cancel by user
	<-triggerEnd

}
