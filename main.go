package main

import (
	"time"

	"github.com/PiterWeb/RemoteController/src/gamepad"
	"github.com/PiterWeb/RemoteController/src/desktop"
)

func main() {

	go desktop.InitWindow()

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
