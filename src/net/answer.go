package net

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/PiterWeb/RemoteController/src/gamepad"
	"github.com/pion/webrtc/v3"
	"github.com/pquerna/ffjson/ffjson"
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

	peerConnection, err := webrtc.NewAPI().NewPeerConnection(config)
	if err != nil {
		panic(err)
	}

	// Open a UDP Listener for RTP Packets on port 5004
	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 5004})
	if err != nil {
		panic(err)
	}

	// Increase the UDP receive buffer size
	// Default UDP buffer sizes vary on different operating systems
	bufferSize := 300000 // 300KB
	err = listener.SetReadBuffer(bufferSize)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = listener.Close(); err != nil {
			panic(err)
		}
	}()

	// Create a video track
	videoTrack, err := webrtc.NewTrackLocalStaticRTP(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeVP8}, "video", "pion")
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

			fmt.Println("RTCP packet received")
			fmt.Println(rtcpBuf)
		}
	}()

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
			if closeErr := peerConnection.Close(); closeErr != nil {
				panic(closeErr)
			}
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

	inboundRTPPacket := make([]byte, 1600) // UDP MTU
	for {
		n, _, err := listener.ReadFrom(inboundRTPPacket)
		if err != nil {
			panic(fmt.Sprintf("error during read: %s", err))
		}

		if _, err = videoTrack.Write(inboundRTPPacket[:n]); err != nil {
			if errors.Is(err, io.ErrClosedPipe) {
				// The peerConnection has been closed.
				return
			}

			panic(err)
		}
	}

	// // Block forever
	// select {}

}
