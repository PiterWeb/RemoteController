// Package webrtc provides all the related functions for WebRTC communication
package webrtc

import (
	"context"
	"log"
	"runtime"
	"strings"

	// "LibreRemotePlay/src/plugins"
	// "LibreRemotePlay/src/devices/audio"
	"LibreRemotePlay/src/devices/gamepad"
	"LibreRemotePlay/src/devices/keyboard"
	"LibreRemotePlay/src/devices/mouse"
	"LibreRemotePlay/src/net/webrtc/streaming_signal"
	"github.com/pion/webrtc/v4"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

var defaultSTUNServers = []string{"stun:stun.l.google.com:19305", "stun:stun.l.google.com:19302", "stun:stun.ipfire.org:3478"}

const ERROR_ANSWER = "ERROR"

func InitHost(wailsCtx context.Context, ctx context.Context, ICEServers []webrtc.ICEServer, offerEncodedWithCandidates string, answerResponse chan<- string, pidChan <-chan uint32) {

	connCtx, cancelConn := context.WithCancel(ctx)
	
	candidates := []webrtc.ICECandidateInit{}

	if len(ICEServers) == 0 {
		ICEServers = []webrtc.ICEServer{
			{
				URLs: defaultSTUNServers,
			},
		}
	}

	// Prepare the configuration
	config := webrtc.Configuration{
		ICEServers: ICEServers,
	}

	defer func() {
		if err := recover(); err != nil {
			answerResponse <- ERROR_ANSWER
			cancelConn()
		}
	}()

	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := peerConnection.Close(); err != nil {
			log.Printf("cannot close peerConnection: %v\n", err)
		}
	}()

	// audioTrack, err := webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypePCMA, Channels: 2}, "audio", "app-audio")

	// if err != nil {
	// 	panic(err)
	// } else if _, err := peerConnection.AddTrack(audioTrack); err != nil {
	// 	panic(err)
	// }

	// Defer close of the wails signaling channel
	defer wailsRuntime.EventsOff(wailsCtx, "streaming-signal-server")

	// Reload plugins in case a new plugin was added or configuration changed
	// plugins.ReloadPlugins()

	// Register data channel creation handling
	peerConnection.OnDataChannel(func(d *webrtc.DataChannel) {

		gamepad.HandleGamepad(d)
		if runtime.GOOS == "linux" {
			streaming_signal.HandleStreamingSignal(connCtx, d)
		} else {
			streaming_signal.HandleStreamingSignal(wailsCtx, d)
		}
		keyboard.HandleKeyboard(d)
		mouse.HandleMouse(wailsCtx, d)
		// plugins.HandleServerPlugins(d)

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
		log.Printf("Peer Connection State has changed: %s\n", s.String())

		wailsRuntime.EventsEmit(wailsCtx, "connection_state", s.String())

		if s == webrtc.PeerConnectionStateFailed {
			cancelConn()
			peerConnection.Close()
		}

		if s == webrtc.PeerConnectionStateClosed {
			cancelConn()
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
			log.Println(err)
			continue
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

	// var audioCtx context.Context
	var cancelAudioCtx context.CancelFunc = func() {}

	for {
		select {
		case pid := <-pidChan:
			log.Printf("Audio pid: %d\n", pid)
			cancelAudioCtx()
			// Pid value to not stream audio
			if pid == 0 {
				continue
			}
			// audioCtx, cancelAudioCtx = context.WithCancel(context.WithValue(context.Background(), "pid", pid))
			// go func() {
				// if err := audio.HandleAudio(audioCtx, audioTrack); err != nil {
				// 	log.Println(err)
				// }
			// }()
		case <-connCtx.Done(): // Block until failed/clossed/canceled by user peerconnection
			answerResponse <- ERROR_ANSWER
			peerConnection.Close()
			cancelAudioCtx()
			return
		}
	}

}
