package net

import (
	"fmt"

	"github.com/pion/mediadevices"
	"github.com/pion/mediadevices/pkg/codec/openh264" // This is required to use h264 video encoder
	"github.com/pion/mediadevices/pkg/codec/opus"     // This is required to use opus audio encoder
	"github.com/pion/mediadevices/pkg/prop"
	"github.com/pion/webrtc/v4"

	_ "github.com/pion/mediadevices/pkg/driver/screen"
)

func getMediaEngine() (*mediadevices.CodecSelector, *webrtc.MediaEngine, error) {

	openh264Params, err := openh264.NewParams()
	if err != nil {
		return nil, nil, err
	}
	openh264Params.BitRate = 1_000_000 // 1mbps

	opusParams, err := opus.NewParams()
	if err != nil {
		return nil, nil, err
	}
	codecSelector := mediadevices.NewCodecSelector(
		mediadevices.WithVideoEncoders(&openh264Params),
		mediadevices.WithAudioEncoders(&opusParams),
	)

	mediaEngine := webrtc.MediaEngine{}
	codecSelector.Populate(&mediaEngine)

	return codecSelector, &mediaEngine, nil

}

func getDisplayMedia(peerConnection *webrtc.PeerConnection, codecSelector *mediadevices.CodecSelector) (*mediadevices.MediaStream, error) {

	mediaStream, err := mediadevices.GetDisplayMedia(mediadevices.MediaStreamConstraints{
		Video: func(c *mediadevices.MediaTrackConstraints) {
			c.FrameRate = prop.Float(30)
			fmt.Println("MediaTrackConstraints: ", c)
		},
		Audio: func(c *mediadevices.MediaTrackConstraints) {},
		Codec: codecSelector,
	})

	if err != nil {
		return nil, err
	}

	for _, track := range mediaStream.GetTracks() {
		track.OnEnded(func(err error) {
			fmt.Printf("Track (ID: %s) ended with error: %v\n",
				track.ID(), err)
		})

		_, err = peerConnection.AddTransceiverFromTrack(track,
			webrtc.RTPTransceiverInit{
				Direction: webrtc.RTPTransceiverDirectionSendonly,
			},
		)
		if err != nil {
			return nil, err
		}
	}

	return &mediaStream, nil

}
