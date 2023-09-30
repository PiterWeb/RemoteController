package net

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/pquerna/ffjson/ffjson"

	"github.com/PiterWeb/RemoteController/src/gamepad"
	"github.com/pion/webrtc/v3"
)

func InitOffer(offerChan chan<- string, answerResponseEncoded <-chan string) {

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

	// Create a datachannel with label 'data'
	dataChannel, err := peerConnection.CreateDataChannel("data", nil)
	if err != nil {
		panic(err)
	}

	// Set a handler for when a new remote track starts
	peerConnection.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		fmt.Printf("Track has started streamId(%s) id(%s) rid(%s) \n", track.StreamID(), track.ID(), track.RID())

		for {
			// Read the RTCP packets as they become available for our new remote track
			rtcpPackets, _, rtcpErr := receiver.ReadRTCP()
			if rtcpErr != nil {
				panic(rtcpErr)
			}

			for _, r := range rtcpPackets {
				// Print a string description of the packets
				if stringer, canString := r.(fmt.Stringer); canString {
					fmt.Printf("Received RTCP Packet: %v", stringer.String())
				}
			}
		}
	})

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
		fmt.Printf("Peer Connection State has changed: %s\n", s.String())

		if s == webrtc.PeerConnectionStateFailed {
			// Wait until PeerConnection has had no network activity for 30 seconds or another failure. It may be reconnected using an ICE Restart.
			// Use webrtc.PeerConnectionStateDisconnected if you are interested in detecting faster timeout.
			// Note that the PeerConnection may come back from PeerConnectionStateDisconnected.
			fmt.Println("Peer Connection has gone to failed exiting")
			os.Exit(0)
		}
	})

	// Gamepad update loop
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

				padRaw, _ := ffjson.Marshal(*pad)

				err := dataChannel.Send(padRaw)

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

	// Block forever
	select {}
}
