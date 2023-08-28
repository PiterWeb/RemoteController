package main

import (
	// "fmt"
	"time"
	"github.com/PiterWeb/RemoteController/gamepad"
)



func main() {
	gamepads := gamepad.All{}

	// prev := int16(0)

	go func ()  {
		
		for range time.Tick(1 * time.Millisecond) {
			gamepads.Update()
			for i := range gamepads {
				pad := &gamepads[i]
				if !pad.Connected {
					continue
				}
				
				// if pad.Raw.ThumbLX != prev {
				// 	fmt.Println(time.Now().Nanosecond(), pad.Raw.ThumbLX, pad.Raw.ThumbRX)
				// 	prev = pad.Raw.ThumbLX
				// }
			}
		}
	}()

	virtualDevice, err := gamepad.InitializeEmulatedDevice()

	if err != nil {
		panic(err)
	}

	err = gamepad.EmulateDevice(virtualDevice)

	if err != nil {
		panic(err)
	}
}
