//go:build windows
package gamepad

import (
	"math/rand"
	"os"
	"testing"
	"time"
)

func TestUpdateVirtualDevice(t *testing.T) {

	vigemSetup, err := os.Create("./.vigemsetup")

	if err != nil {
		t.Fatalf("Failed to create ViGEm setup file %s", err.Error())
	}

	defer func() {
		vigemSetup.Close()
		os.Remove(".vigemsetup")
	}()

	if err := InitViGEm(); err != nil {
		t.Fatalf("Failed to init ViGEm %s", err.Error())
	}

	defer func() {
		CloseViGEmDLL()
		os.Remove("ViGEmClient.dll")
	}()

	timer := time.NewTimer(5 * time.Second).C
	ticker := time.NewTicker(166 * time.Millisecond).C

	virtualState := new(ViGEmState)

	virtualDevice, err := GenerateVirtualDevice()

	if err != nil {
		t.Fatalf("Failed to generate ViGEm virtual device %s", err.Error())
	}

	defer FreeTargetAndDisconnect(virtualDevice)

	gamepadState := &GamepadAPIXState{
		Index: 0,
	}

	for {
		select {
		case <-ticker:

			gamepadState.Axes = [4]float64{rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64()}

			for i := range gamepadState.Buttons {
				gamepadState.Buttons[i].Pressed = true
				gamepadState.Buttons[i].Value = rand.Float64()
			}

			if err != nil {
				t.Fatalf("Failed to create gamepad state: %s", err.Error())
			}

			go UpdateVirtualDevice(virtualDevice, *gamepadState, virtualState)
		case <-timer:
			return
		}
	}

}
