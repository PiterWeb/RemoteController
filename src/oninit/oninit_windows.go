package oninit

import "github.com/PiterWeb/RemoteController/src/gamepad"

func Execute() error {
	err := gamepad.InitViGEm()

	if err != nil {
		return err
	}

	return nil
}
