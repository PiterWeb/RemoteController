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

	d.OnMessage(func(msg webrtc.DataChannelMessage) {

		fmt.Println("keyboard message: ", msg.Data)

		if !msg.IsString || msg.Data == nil {
			return
		}

		fmt.Println("keyboard message: ", string(msg.Data))

		keyJS := strings.ToUpper(string(msg.Data))

		switch keyJS {
		case "SHIFTLEFT_1":
			kb.HasSHIFT(true)
			return
		case "SHIFTLEFT_0":
			kb.HasSHIFT(false)
		case "CTRLLEFT_1":
			kb.HasCTRL(true)
			return
		case "CTRLLEFT_0":
			kb.HasCTRL(false)
		case "ALTLEFT_1":
			kb.HasALT(true)
			return
		case "ALTLEFT_0":
			kb.HasALT(false)
			return
		default:
			key := keyBoardJSToGolang[keyJS]
			kb.SetKeys(key)
			_ = kb.Launching()
		}

	})

}
