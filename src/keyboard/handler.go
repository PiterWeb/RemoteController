package keyboard

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/micmonay/keybd_event"
	"github.com/pion/webrtc/v3"
)

func HandleKeyboard(d *webrtc.DataChannel) {

	if d.Label() != "keyboard" {
		return
	}

	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		panic(err)
	}

	// For linux, it is very important to wait 2 seconds
	if runtime.GOOS == "linux" {
		time.Sleep(2 * time.Second)
	}

	d.OnOpen(func() {
		fmt.Println("keyboard data channel is open")
	})

	var SHIFT_KEY_PRESSED, CTRL_KEY_PRESSED, ALT_KEY_PRESSED bool

	d.OnMessage(func(msg webrtc.DataChannelMessage) {

		fmt.Println("keyboard message: ", msg.Data)

		if !msg.IsString || msg.Data == nil {
			return
		}

		fmt.Println("keyboard message: ", string(msg.Data))

		keyJS := strings.ToUpper(string(msg.Data))

		switch keyJS {
		case "SHIFTLEFT_1":
			SHIFT_KEY_PRESSED = true
			return
		case "SHIFTLEFT_0":
			SHIFT_KEY_PRESSED = false
		case "CTRLLEFT_1":
			CTRL_KEY_PRESSED = true
			return
		case "CTRLLEFT_0":
			CTRL_KEY_PRESSED = false
		case "ALTLEFT_1":
			ALT_KEY_PRESSED = true
			return
		case "ALTLEFT_0":
			ALT_KEY_PRESSED = false
			return
		default:
			if key, ok := keyBoardJSToGolang[keyJS]; ok {
				kb.SetKeys(key)

				kb.HasALT(ALT_KEY_PRESSED)
				kb.HasCTRL(CTRL_KEY_PRESSED)
				kb.HasSHIFTR(SHIFT_KEY_PRESSED)

				kb.Press()
				// This sleep is arbitrary, it may be necessary to adjust it
				time.Sleep(10 * time.Millisecond)
				kb.Release()
			}
		}

	})

}
