package desktop

import (
	"os"
	"runtime"

	"github.com/PiterWeb/RemoteController/src/net"
	"github.com/rodrigocfd/windigo/ui"
)

// This struct represents our main window.
type MyWindow struct {
	wnd     ui.WindowMain
	lblName ui.Static
	txtName ui.Edit
	btnShow ui.Button
}

func RunDesktop() {

	go net.CreateWebRTCConn()

	runtime.LockOSThread()

	myWindow := newWindow()  // instantiate
	myWindow.wnd.RunAsMain() // ...and run

	runtime.UnlockOSThread()
	os.Exit(0)

}
