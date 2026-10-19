//go:build unix

package gamepad

import (
	"context"
	"log"
	"testing"
	"time"
)

func TestUpdateVirtualDevice(t *testing.T) {
	virtualGamepad, err := generateVirtualDevice()
	if err != nil {
		log.Println(err)
		log.Println("VirtualGamepad is not defined, gamepad disabled by default")
	}

	defer virtualGamepad.Unregister()

	var lastPad GamepadAPIXState
	var actualPad GamepadAPIXState

	ctx, cancelCtx := context.WithTimeout(context.Background(), time.Second * 20)
	defer cancelCtx()
	
	ticker := time.NewTicker(1 * time.Second)

	for {
		select {
			case <-ticker.C: 		
				actualPad.Buttons[12].Pressed = !actualPad.Buttons[12].Pressed
				
				updateVirtualDevice(virtualGamepad, actualPad, lastPad)
				lastPad = actualPad
			case <-ctx.Done():
				return
		}
	}
	
}