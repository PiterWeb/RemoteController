package multimedia

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/pion/webrtc/v3"
	"github.com/pion/webrtc/v3/pkg/media"
	"github.com/pion/webrtc/v3/pkg/media/ivfreader"
)

type MultimediaVideo struct {
}

type MultimediaAudio struct {
}

func (MultimediaVideo) Read(b []byte) (int, error) {
	return 0, nil
}

func (MultimediaAudio) Read(b []byte) (int, error) {
	return 0, nil
}

func SendVideoAndAudio(peerConnection *webrtc.PeerConnection, video MultimediaVideo, audio MultimediaAudio) {

	_, header, openErr := ivfreader.NewWith(video)

	if openErr != nil {
		panic(openErr)
	}

	var trackCodec string
	switch header.FourCC {
	case "AV01":
		trackCodec = webrtc.MimeTypeAV1
	case "VP90":
		trackCodec = webrtc.MimeTypeVP9
	case "VP80":
		trackCodec = webrtc.MimeTypeVP8
	default:
		panic(fmt.Sprintf("Unable to handle FourCC %s", header.FourCC))
	}

	videoTrack, videoTrackErr := webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: trackCodec}, "video", "pion")
	if videoTrackErr != nil {
		panic(videoTrackErr)
	}

	rtpSender, videoTrackErr := peerConnection.AddTrack(videoTrack)
	if videoTrackErr != nil {
		panic(videoTrackErr)
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

	go func() {

		ivf, header, ivfErr := ivfreader.NewWith(video)
		if ivfErr != nil {
			panic(ivfErr)
		}

		// Send our video file frame at a time. Pace our sending so we send it at the same speed it should be played back as.
		// This isn't required since the video is timestamped, but we will such much higher loss if we send all at once.
		//
		// It is important to use a time.Ticker instead of time.Sleep because
		// * avoids accumulating skew, just calling time.Sleep didn't compensate for the time spent parsing the data
		// * works around latency issues with Sleep (see https://github.com/golang/go/issues/44343)
		ticker := time.NewTicker(time.Millisecond * time.Duration((float32(header.TimebaseNumerator)/float32(header.TimebaseDenominator))*1000))
		for ; true; <-ticker.C {
			frame, _, ivfErr := ivf.ParseNextFrame()
			if errors.Is(ivfErr, io.EOF) {
				fmt.Printf("All video frames parsed and sent")
				os.Exit(0)
			}

			if ivfErr != nil {
				panic(ivfErr)
			}

			if ivfErr = videoTrack.WriteSample(media.Sample{Data: frame, Duration: time.Second}); ivfErr != nil {
				panic(ivfErr)
			}
		}
	}()

	audioTrack, audioTrackErr := webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeOpus}, "audio", "pion")
	if audioTrackErr != nil {
		panic(audioTrackErr)
	}

	rtpSender, audioTrackErr = peerConnection.AddTrack(audioTrack)
	if audioTrackErr != nil {
		panic(audioTrackErr)
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

	go func() {

		ivf, header, ivfErr := ivfreader.NewWith(audio)
		if ivfErr != nil {
			panic(ivfErr)
		}

		// Send our video file frame at a time. Pace our sending so we send it at the same speed it should be played back as.
		// This isn't required since the video is timestamped, but we will such much higher loss if we send all at once.
		//
		// It is important to use a time.Ticker instead of time.Sleep because
		// * avoids accumulating skew, just calling time.Sleep didn't compensate for the time spent parsing the data
		// * works around latency issues with Sleep (see https://github.com/golang/go/issues/44343)
		ticker := time.NewTicker(time.Millisecond * time.Duration((float32(header.TimebaseNumerator)/float32(header.TimebaseDenominator))*1000))
		for ; true; <-ticker.C {
			frame, _, ivfErr := ivf.ParseNextFrame()
			if errors.Is(ivfErr, io.EOF) {
				fmt.Printf("All video frames parsed and sent")
				os.Exit(0)
			}

			if ivfErr != nil {
				panic(ivfErr)
			}

			if ivfErr = videoTrack.WriteSample(media.Sample{Data: frame, Duration: time.Second}); ivfErr != nil {
				panic(ivfErr)
			}
		}
	}()

}
