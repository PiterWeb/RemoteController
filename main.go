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

	virtualStates := [gamepad.ID(4)]*gamepad.XInputState{}

	for range time.Tick(1 * time.Millisecond) {
		gamepads.Update()
		for i := range gamepads {
			pad := &gamepads[i]
			virtualState := virtualStates[i]

			if !pad.Connected {
				continue
			}

			if virtualStates[i] == nil {
				virtualState = new(gamepad.XInputState)
			}

			gamepad.UpdateVirtualDevice(virtualDevice, *pad, virtualState)

		}
	}
}
