package gamepad

type receivedGamepad struct {
	Axes      []float64
	Buttons   []gamepadButton
	Connected bool
	ID        string
	Index     int
}

type gamepadButton struct {
	Pressed bool
	Value   float64
}
