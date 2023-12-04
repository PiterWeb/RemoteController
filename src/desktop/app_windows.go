package desktop

import (
	"github.com/PiterWeb/RemoteController/src/gamepad"
)

func (a *App) OpenViGEmWizard() (err string) {

	return gamepad.OpenViGEmWizard().Error()

}
