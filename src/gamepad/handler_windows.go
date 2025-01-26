package gamepad

import (
	"time"

	"github.com/pion/webrtc/v3"
	"github.com/pquerna/ffjson/ffjson"
)

var buttonValueToHexMap = map[int]uint16{
	0:  0x1000,
	1:  0x2000,
	2:  0x4000,
	3:  0x8000,
	4:  0x0100,
	5:  0x0200,
	8:  0x0020,
	9:  0x0010,
	10: 0x0040,
	11: 0x0080,
	12: 0x0001,
	13: 0x0002,
	14: 0x0004,
	15: 0x0008,
}

var virtualDevice EmulatedDevice

func HandleGamepad(gamepadChannel *webrtc.DataChannel) {

	if gamepadChannel.Label() != "controller" {
		return
	}

	defer func() {
		if err := recover(); err != nil {
			FreeTargetAndDisconnect(virtualDevice)
		}
	}()

	virtualState := new(ViGEmState)

	// Create a virtual device
	gamepadChannel.OnOpen(func() {

		var err error = nil
		virtualDevice, err = GenerateVirtualDevice()

		if err != nil {
			panic(err)
		}

	})

	gamepadChannel.OnClose(func() {

		FreeTargetAndDisconnect(virtualDevice)

	})

	// Update the virtual device
	gamepadChannel.OnMessage(func(msg webrtc.DataChannelMessage) {

		var pad GamepadAPIXState

		ffjson.Unmarshal(msg.Data, &pad)

		go UpdateVirtualDevice(virtualDevice, pad, virtualState)

	})

}

func gamepadAPIXToXInput(gms GamepadAPIXState) XInputState {

	return XInputState{
		ID:        ID(gms.Index),
		Connected: gms.Connected,
		Packet:    uint32(time.Now().Nanosecond()), // Different values trigger update
		Raw: RawControls{
			Buttons:      convertGamepadButtons(gms.Buttons),
			LeftTrigger:  convertFloatToUint8(gms.Buttons[6].Value),
			RightTrigger: convertFloatToUint8(gms.Buttons[7].Value),
			ThumbLX:      convertFloatToInt16(gms.Axes[0]),
			ThumbLY:      convertFloatToInt16(fixLYAxis(gms.Axes[1])),
			ThumbRX:      convertFloatToInt16(gms.Axes[2]),
			ThumbRY:      convertFloatToInt16(fixRYAxis(gms.Axes[3])),
		},
	}

}

func convertGamepadButtons(buttons [16]gamepadButton) Button {
	var result Button

	for i, button := range buttons {
		if button.Pressed {
			result += Button(buttonValueToHexMap[i])
		}
	}
	return result
}

func convertFloatToUint8(value float64) uint8 {
	return uint8(value * 255)
}

func convertFloatToInt16(value float64) int16 {
	return int16(value * 32767)
}
