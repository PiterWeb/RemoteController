package desktop

import (
	"os"
	"runtime"
)

// This struct represents our main window.

func RunDesktop() {

	runtime.LockOSThread()

	mainWindow := initWindow() // instantiate
	initLogic(mainWindow)
	mainWindow.wnd.RunAsMain() // ...and run

	runtime.UnlockOSThread()
	os.Exit(0)

}
