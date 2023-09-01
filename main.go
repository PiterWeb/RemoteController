package main

import (
	"time"
	"github.com/PiterWeb/RemoteController/src/desktop"
	"github.com/PiterWeb/RemoteController/src/gamepad"
)

func main() {

	go desktop.RunDesktop()

	gamepads := gamepad.All{}

	virtualDevice, err := gamepad.GenerateVirtualDevice()

	if err != nil {
		panic(err)
	}

	defer gamepad.FreeTargetAndDisconnect(virtualDevice)

	virtualState := new(gamepad.ViGEmState)

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
