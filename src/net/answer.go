package net

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/pquerna/ffjson/ffjson"

	"github.com/PiterWeb/RemoteController/src/gamepad"
	"github.com/PiterWeb/RemoteController/src/net/multimedia"
	"github.com/pion/rtcp"
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
	webrtcWebm, saver := multimedia.InitWebRTCAnswer()

	peerConnection, err := webrtcWebm.NewPeerConnection(config)
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

			err = ffjson.Unmarshal(msg.Data, &pad)

			go gamepad.UpdateVirtualDevice(virtualDevice, pad, virtualState)

		})

	})

	peerConnection.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {

		go func() {
			ticker := time.NewTicker(time.Second * 3)
			for range ticker.C {
				errSend := peerConnection.WriteRTCP([]rtcp.Packet{&rtcp.PictureLossIndication{MediaSSRC: uint32(track.SSRC())}})
				if errSend != nil {
					fmt.Println(errSend)
				}
			}
		}()

		fmt.Printf("Track has started, of type %d: %s \n", track.PayloadType(), track.Codec().RTPCodecCapability.MimeType)
		for {
			// Read RTP packets being sent to Pion
			rtp, _, readErr := track.ReadRTP()
			if readErr != nil {
				if readErr == io.EOF {
					return
				}
				panic(readErr)
			}
			switch track.Kind() {
			case webrtc.RTPCodecTypeAudio:
				saver.PushOpus(rtp)
			case webrtc.RTPCodecTypeVideo:
				saver.PushVP8(rtp)
			}
		}

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

	// Block forever
	select {}

}
