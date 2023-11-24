package gamepad

type GamepadAPIState struct {
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
