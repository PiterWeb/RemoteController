package gamepad

import "math"

const (
	threshold float64 = 1e-9
)

var (
	prevThumbLY float64
	prevThumbRY float64
)

// Struct for GamepadAPI for XINPUT gamepads
type GamepadAPIXState struct {
	Axes      [4]float64
	Buttons   [16]gamepadButton
	Connected bool
	ID        string
	Index     int
}

type gamepadButton struct {
	Pressed bool
	Value   float64
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
