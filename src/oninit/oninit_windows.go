package oninit

import (
	"LibreRemotePlay/src/devices/gamepad"
	"LibreRemotePlay/src/state"
)

func executePlatform(_ state.App) error {
	err := gamepad.InitViGEm()

	return err
}
