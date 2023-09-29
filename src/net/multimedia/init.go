package multimedia

import (
	"github.com/pion/ice/v2"
	"github.com/pion/webrtc/v3"
)

func InitWebRTCAnswer() (*webrtc.API, *webmSaver) {
	// Create a MediaEngine object to configure the supported codec
	m := &webrtc.MediaEngine{}

	settingEngine := webrtc.SettingEngine{}

	mux, err := ice.NewMultiUDPMuxFromPort(8443)

	if err != nil {
		panic(err)
	}

	settingEngine.SetICEUDPMux(mux)

	// Setup the codecs you want to use.
	// Only support VP8 and OPUS, this makes our WebM muxer code simpler
	if err := m.RegisterCodec(webrtc.RTPCodecParameters{
		RTPCodecCapability: webrtc.RTPCodecCapability{MimeType: "video/VP8", ClockRate: 90000, Channels: 0, SDPFmtpLine: "", RTCPFeedback: nil},
		PayloadType:        96,
	}, webrtc.RTPCodecTypeVideo); err != nil {
		panic(err)
	}
	if err := m.RegisterCodec(webrtc.RTPCodecParameters{
		RTPCodecCapability: webrtc.RTPCodecCapability{MimeType: "audio/opus", ClockRate: 48000, Channels: 0, SDPFmtpLine: "", RTCPFeedback: nil},
		PayloadType:        111,
	}, webrtc.RTPCodecTypeAudio); err != nil {
		panic(err)
	}

	// Create the API object with the MediaEngine
	return webrtc.NewAPI(webrtc.WithMediaEngine(m), webrtc.WithSettingEngine(settingEngine)), newWebmSaver()
}

func InitWebRTCOffer() *webrtc.API {

	mediaEngine := &webrtc.MediaEngine{}

	return webrtc.NewAPI(webrtc.WithMediaEngine(mediaEngine))

}
