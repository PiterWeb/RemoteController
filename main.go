package main

import (
	"time"

	"github.com/PiterWeb/RemoteController/gamepad"
)

func main() {
	gamepads := gamepad.All{}

	virtualDevice, err := gamepad.GenerateVirtualDevice()

	if err != nil {
		panic(err)
	}

	defer gamepad.FreeTargetAndDisconnect(virtualDevice)

	virtualState := new(gamepad.XInputState)

	for range time.Tick(1 * time.Millisecond) {
		gamepads.Update()
		for i := range gamepads {
			pad := &gamepads[i]

			if !pad.Connected {
				continue
			}

			go gamepad.UpdateVirtualDevice(virtualDevice, *pad, virtualState)

		}
	}
}
