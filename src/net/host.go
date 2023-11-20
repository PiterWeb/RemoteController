package net

import (
	"context"
	// "errors"
	// "io"

	// "encoding/json"
	"fmt"
	"sync"

	"github.com/PiterWeb/RemoteController/src/gamepad"
	"github.com/pion/webrtc/v3"
)

var videoTrack *webrtc.TrackLocalStaticRTP

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

	videoTrack, err = webrtc.NewTrackLocalStaticRTP(webrtc.RTPCodecCapability{
		MimeType: webrtc.MimeTypeVP9,
	}, "video", "pion")

	if err != nil {
		panic(err)
	}

	rtpSender, err := peerConnection.AddTrack(videoTrack)

	if err != nil {
		panic(err)
	}

	// Read incoming RTCP packets
	// Before these packets are returned they are processed by interceptors. For things
	// like NACK this needs to be called.
	go func() {
		rtcpBuf := make([]byte, 1500)
		for {
			if _, _, rtcpErr := rtpSender.Read(rtcpBuf); rtcpErr != nil {
				return
			}
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

		gamepad.HandleGamepad(d)

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

	stopStreaming := make(chan struct{})

	// startStreamingPeer(stopStreaming)

	// Block until cancel by user
	<-triggerEnd
	stopStreaming <- struct{}{}

}

// func startStreamingPeer(stopStreaming <-chan struct{}) {

// 	var candidatesMux sync.Mutex
// 	candidates := []string{}

// 	// Prepare the configuration
// 	config := webrtc.Configuration{
// 		ICEServers: []webrtc.ICEServer{
// 			{
// 				URLs: []string{"stun:stun.l.google.com:19302"},
// 			},
// 		},
// 	}

// 	peerConnection, err := webrtc.NewAPI().NewPeerConnection(config)
// 	if err != nil {
// 		panic(err)
// 	}

// 	defer func() {
// 		if err := peerConnection.Close(); err != nil {
// 			fmt.Printf("cannot close peerConnection: %v\n", err)
// 		}
// 	}()

// 	peerConnection.OnTrack(func(tr *webrtc.TrackRemote, r *webrtc.RTPReceiver) {

// 		go func() {
// 			rtcpBuf := make([]byte, 1500)
// 			for {
// 				if _, _, rtcpErr := r.Read(rtcpBuf); rtcpErr != nil {
// 					return
// 				}
// 			}
// 		}()

// 		rtcpBuf := make([]byte, 1400)

// 		for {
// 			i, _, readErr := tr.Read(rtcpBuf)

// 			if readErr != nil {
// 				fmt.Println(readErr)
// 			}

// 			_, err := videoTrack.Write(rtcpBuf[:i])

// 			// ErrClosedPipe means we don't have any subscribers, this is ok if no peers have connected yet
// 			if err != nil && errors.Is(readErr, io.ErrClosedPipe) {
// 				fmt.Println(readErr)
// 			}

// 		}

// 	})

// 	peerConnection.OnICECandidate(func(c *webrtc.ICECandidate) {

// 		if c == nil {
// 			return
// 		}

// 		candidatesMux.Lock()
// 		defer candidatesMux.Unlock()

// 		desc := peerConnection.RemoteDescription()

// 		if desc != nil {
// 			candidates = append(candidates, (*c).ToJSON().Candidate)
// 		}

// 	})

// 	// Set the handler for Peer connection state
// 	// This will notify you when the peer has connected/disconnected
// 	peerConnection.OnConnectionStateChange(func(s webrtc.PeerConnectionState) {
// 		fmt.Printf("Peer Connection State has changed: %s\n", s.String())

// 		if s == webrtc.PeerConnectionStateFailed {
// 			if closeErr := peerConnection.Close(); closeErr != nil {
// 				panic(closeErr)
// 			}
// 		}
// 	})

// 	offer := webrtc.SessionDescription{}
// 	signalDecode(offerEncoded, &offer)

// 	if err := peerConnection.SetRemoteDescription(offer); err != nil {
// 		panic(err)
// 	}

// 	// Create an answer to send to the other process
// 	answer, err := peerConnection.CreateAnswer(nil)
// 	if err != nil {
// 		panic(err)
// 	}

// 	gatherComplete := webrtc.GatheringCompletePromise(peerConnection)

// 	// Sets the LocalDescription, and starts our UDP listeners
// 	err = peerConnection.SetLocalDescription(answer)
// 	if err != nil {
// 		panic(err)
// 	}

// 	<-gatherComplete

// 	answerResponse <- signalEncode(*peerConnection.LocalDescription()) + ";" + signalEncode(candidates)

// 	<-stopStreaming

// }
