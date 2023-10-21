package net

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/pquerna/ffjson/ffjson"

	"github.com/PiterWeb/RemoteController/src/gamepad"
	"github.com/pion/webrtc/v3"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func InitClient(ctx context.Context, offerChan chan<- string, answerResponseEncoded <-chan string, triggerEnd <-chan struct{}) {

	var candidatesMux sync.Mutex
	candidates := make([]string, 0)

	// Prepare the configuration
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	// Create a new RTCPeerConnection
	peerConnection, err := webrtc.NewAPI().NewPeerConnection(config)
	if err != nil {
		panic(err)
	}
	defer func() {
		if cErr := peerConnection.Close(); cErr != nil {
			fmt.Printf("cannot close peerConnection: %v\n", cErr)
		}
	}()

	// Create a datachannel with label 'controller'
	controllerChannel, err := peerConnection.CreateDataChannel("controller", nil)
	if err != nil {
		panic(err)
	}

	// When an ICE candidate is available send to the other Pion instance
	// the other Pion instance will add this candidate by calling AddICECandidate
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

		defer func() {

		}()

		fmt.Printf("Peer Connection State has changed: %s\n", s.String())

		if s == webrtc.PeerConnectionStateFailed {
			// Wait until PeerConnection has had no network activity for 30 seconds or another failure. It may be reconnected using an ICE Restart.
			// Use webrtc.PeerConnectionStateDisconnected if you are interested in detecting faster timeout.
			// Note that the PeerConnection may come back from PeerConnectionStateDisconnected.
			fmt.Println("Peer Connection has gone to failed exiting")
			panic("Peer Connection has gone to failed exiting")
		}
	})

	peerConnection.OnDataChannel(func(d *webrtc.DataChannel) {

		fmt.Println(d.Label())

		if d.Label() == "streaming" {

			d.OnOpen(func() {
				fmt.Println("Streaming Started: client")
			})

			d.OnMessage(func(msg webrtc.DataChannelMessage) {
				fmt.Println("Packet Emited to Frontend")
				runtime.EventsEmit(ctx, "receive-streaming", string(msg.Data))
			})

		}

	})

	// Gamepad update loop
	controllerChannel.OnOpen(func() {

		gamepads := gamepad.All{}

		defer func() {

			if err := recover(); err != nil {
				fmt.Println(err)
			}

		}()

		for range time.NewTicker(1 * time.Millisecond).C {
			gamepads.Update()
			for i := range gamepads {
				pad := &gamepads[i]

				if !pad.Connected {
					continue
				}

				padRaw, _ := ffjson.Marshal(*pad)

				err := controllerChannel.Send(padRaw)

				if err != nil {
					panic(err)
				}

			}
		}

	})

	// Create an offer to send to the other process
	offer, err := peerConnection.CreateOffer(nil)
	if err != nil {
		panic(err)
	}

	gatherComplete := webrtc.GatheringCompletePromise(peerConnection)

	if err = peerConnection.SetLocalDescription(offer); err != nil {
		panic(err)
	}

	<-gatherComplete

	offerChan <- signalEncode(offer)

	answerResponse := strings.Split(<-answerResponseEncoded, ";")

	if len(answerResponse) != 2 {
		panic("No candidate or answer comming")
	}

	answer := webrtc.SessionDescription{}

	signalDecode(answerResponse[0], &answer)

	remoteCandidates := []string{}

	signalDecode(answerResponse[1], &remoteCandidates)

	if err = peerConnection.SetRemoteDescription(answer); err != nil {
		panic(err)
	}

	for _, candidate := range remoteCandidates {
		err := peerConnection.AddICECandidate(webrtc.ICECandidateInit{
			Candidate: candidate,
		})

		if err != nil {
			panic(err)
		}
	}

	// Block until cancel by user
	<-triggerEnd
}
