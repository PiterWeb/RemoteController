package gamepad

import (
	"github.com/pion/webrtc/v3"
	"github.com/pquerna/ffjson/ffjson"
)

func HandleGamepad(gamepadChannel *webrtc.DataChannel) {

	if gamepadChannel.Label() != "controller" {
		return
	}

	var virtualDevice EmulatedDevice
	defer FreeTargetAndDisconnect(virtualDevice)

	virtualState := new(ViGEmState)

	// Create a virtual device
	gamepadChannel.OnOpen(func() {

		var err error = nil
		virtualDevice, err = GenerateVirtualDevice()

		if err != nil {
			panic(err)
		}

	})

	// Update the virtual device
	gamepadChannel.OnMessage(func(msg webrtc.DataChannelMessage) {

		var pad receivedGamepad

		ffjson.Unmarshal(msg.Data, &pad)

		go UpdateVirtualDevice(virtualDevice, pad, virtualState)

	})

}

func receivedGamepadToXInput(rp receivedGamepad) XInputState {

	if len(rp.Axes) != 6 {
		return XInputState{
			ID:        ID(rp.Index), // You may need to adjust this based on your requirements
			Connected: rp.Connected,
			Packet:    uint32(rp.Index), // You may need to set a proper value for Packet
			Raw: RawControls{
				Buttons:      convertGamepadButtons(rp.Buttons),
				LeftTrigger:  convertFloatToUint8(0),
				RightTrigger: convertFloatToUint8(0),
				ThumbLX:      convertFloatToInt16(0),
				ThumbLY:      convertFloatToInt16(0),
				ThumbRX:      convertFloatToInt16(0),
				ThumbRY:      convertFloatToInt16(0),
			},
		}
	}

	xinputState := XInputState{
		ID:        ID(rp.Index), // You may need to adjust this based on your requirements
		Connected: rp.Connected,
		Packet:    0, // You may need to set a proper value for Packet
		Raw: RawControls{
			Buttons:      convertGamepadButtons(rp.Buttons),
			LeftTrigger:  convertFloatToUint8(rp.Axes[4]),
			RightTrigger: convertFloatToUint8(rp.Axes[5]),
			ThumbLX:      convertFloatToInt16(rp.Axes[0]),
			ThumbLY:      convertFloatToInt16(rp.Axes[1]),
			ThumbRX:      convertFloatToInt16(rp.Axes[2]),
			ThumbRY:      convertFloatToInt16(rp.Axes[3]),
		},
	}
	return xinputState

}

func convertGamepadButtons(buttons []gamepadButton) Button {
	var result Button
	for i, button := range buttons {
		if button.Pressed {
			result |= 1 << i
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
