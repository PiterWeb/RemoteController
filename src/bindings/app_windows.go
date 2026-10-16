package bindings

import (
	"LibreRemotePlay/src/devices/gamepad"
)

func (a *App) OpenViGEmWizard() (err string) {

	return gamepad.OpenViGEmWizard().Error()

}
