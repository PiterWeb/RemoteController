package mouse

import (
	"bytes"
	"context"
	"encoding/binary"
	"log"

	"LibreRemotePlay/src/devices"
	"github.com/go-vgo/robotgo"
	"github.com/pion/webrtc/v4"
)

type screenT struct {
	IsCurrent bool `json:"isCurrent"`
	IsPrimary bool `json:"isPrimary"`
	// PhysicalSize is the physical size of the screen in pixels
	PhysicalSize screenSizeT `json:"physicalSize"`
}

type screenSizeT struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

var MouseEnabled = new(devices.DeviceEnabled).Disable()

func HandleMouse(ctx context.Context, d *webrtc.DataChannel) {

	if d.Label() != "mouse" {
		return
	}

	d.OnOpen(func() {
		log.Println("mouse data channel is open")
	})

	// TODO: change mock values for screen
	screens := []screenT{{
		IsCurrent: true,
		IsPrimary: true,
		PhysicalSize: screenSizeT{
			Width: 1920,
			Height: 1080,
		},
	}}
	
	moveFun, unregisterMouse := moveMouseHandler(screens)

	d.OnClose(func() {
		if err := unregisterMouse(); err != nil {
			log.Printf("Mouse error: %s", err)
		}
	})
	
	d.OnMessage(func(msg webrtc.DataChannelMessage) {

		if !MouseEnabled.IsEnabled() {
			return
		}

		if msg.Data == nil || msg.IsString {
			return
		}
		
		if len(msg.Data) <= 1 {
			return
		}

		msgBuf := bytes.NewBuffer(msg.Data)
		
		msgType, err := msgBuf.ReadByte()
		
		if err != nil {
			return
		}
		
		// Handle click event
		if msgType == typeMsgClick {
			
			clickBtn, err := msgBuf.ReadByte()
			
			if err != nil {
				return
			}
			
			btnState, err := msgBuf.ReadByte()
			
			if err != nil {
				return
			}
			
			var state string
			
			switch btnState {
				case mouseDown:
				state = "down"
				case mouseUp:
				state = "up"
			}
						
			switch clickBtn {
				case mouseLeft:
				robotgo.Toggle("left", state)
				case mouseCentral:
				robotgo.Toggle("center", state)
				case mouseRight:
				robotgo.Toggle("right", state)
			}
			
		} else if msgType == typeMsgMove { // Handle move event
			
			x := make([]byte, 2)
			_, err := msgBuf.Read(x)
			
			if err != nil {
				return
			}
						
			y := make([]byte, 2)
			_, err = msgBuf.Read(y)
			
			if err != nil {
				return
			}
			
			videoW := make([]byte, 2)
			_, err = msgBuf.Read(videoW)
			
			if err != nil {
				return
			}
			
			videoH := make([]byte, 2)
			_, err = msgBuf.Read(videoH)
			
			if err != nil {
				return
			}
			
			xNum := binary.BigEndian.Uint16(x)
			yNum := binary.BigEndian.Uint16(y)
			videoWNum := binary.BigEndian.Uint16(videoW)
			videoHNum := binary.BigEndian.Uint16(videoH)
			
			log.Printf("Mouse x: %d, y:%d\n", int(xNum), int(yNum))
			
			moveFun(int32(xNum), int32(yNum), int32(videoWNum), int32(videoHNum))
			
		}

	})

}
