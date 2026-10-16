package gamepad

import (
	"math"

	"LibreRemotePlay/src/devices"
)

const (
	threshold float64 = 1e-9
)

var (
	prevThumbLY float64
	prevThumbRY float64
)

var GamepadEnabled = new(devices.DeviceEnabled).Enable()

// Struct for GamepadAPI for XINPUT gamepads
type GamepadAPIXState struct {
	Axes      [4]float64
	Buttons   [16]gamepadButton
	// ID        string
	Index     int
	// Connected bool
}

type gamepadButton struct {
	Value   float64
	Pressed bool
}

func fixLYAxis(value float64) float64 {

	if math.Abs(value-prevThumbLY) <= threshold {
		return prevThumbLY
	}

	prevThumbLY = -value
	return -value

}

func fixRYAxis(value float64) float64 {

	if math.Abs(value-prevThumbRY) <= threshold {
		return prevThumbRY
	}

	prevThumbRY = -value
	return -value

}
