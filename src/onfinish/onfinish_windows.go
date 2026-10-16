package onfinish

import "LibreRemotePlay/src/devices/gamepad"

func Execute() error {

	if err := gamepad.CloseViGEmDLL(); err != nil {
		return err
	}

	return nil
}