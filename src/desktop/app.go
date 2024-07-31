package desktop

import (
	"context"
	"fmt"
	"strings"

	"runtime"

	"github.com/PiterWeb/RemoteController/src/plugins/messaging"
	"github.com/pion/webrtc/v3"
	wRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

var triggerEnd chan struct{} = make(chan struct{})
var openPeer bool = false

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// Startup is called at application Startup
func (a *App) Startup(ctx context.Context) {
	// Perform your setup here
	a.ctx = ctx
}

// BeforeClose is called when the application is about to quit,
// either by clicking the window close button or calling runtime.Quit.
// Returning true will cause the application to continue, false will continue shutdown as normal.
func (a *App) BeforeClose(ctx context.Context) (prevent bool) {

	if !openPeer {
		return false
	}

	// Show a dialog to confirm the user wants to quit
	option, err := wRuntime.MessageDialog(ctx, wRuntime.MessageDialogOptions{
		Type:          wRuntime.QuestionDialog,
		Title:         "Quit",
		Message:       "Are you sure you want to quit?",
		Buttons:       []string{"Yes", "No"},
		DefaultButton: "No",
		CancelButton:  "No",
	})

	if err != nil {
		return a.BeforeClose(ctx)
	}

	if option == "Yes" {
		return false
	}

	return true
}

// Shutdown is called at application termination
func (a *App) Shutdown(ctx context.Context) {
	// Perform your teardown here
	a.TryClosePeerConnection()
	messaging.ShutdownServer()
	messaging.Get_Client().Close()
	close(triggerEnd)
}

func (a *App) NotifyCreateClient() {
	openPeer = true
	println("NotifyCreateClient")
}

func (a *App) NotifyCloseClient() {
	openPeer = false
	println("NotifyCloseClient")
}

// Create a Host Peer, it receives the offer encoded and returns the encoded answer response
func (a *App) TryCreateHost(ICEServers []webrtc.ICEServer, offerEncoded string) (value string) {

	if openPeer {
		triggerEnd <- struct{}{}
	}

	openPeer = true

	defer func() {

		if err := recover(); err != nil {

			fmt.Println(err)

			openPeer = false
			value = "ERROR"
		}

	}()

	value = createHost(a.ctx, ICEServers, offerEncoded, triggerEnd)

	return value
}

// Closes the peer connection and returns a boolean indication if a connection existed and was closed or not
func (a *App) TryClosePeerConnection() bool {

	if !openPeer {
		return false
	}

	triggerEnd <- struct{}{}

	openPeer = false

	return true

}

func (a *App) GetCurrentOS() string {
	return strings.ToUpper(runtime.GOOS)
}

// GetPlugins returns the list of plugins
// func (a *App) GetPlugins() []plugins.Plugin {
// 	return plugins.GetPlugins()
// }

// ReloadPlugins will reload the plugins from the json file
// func (a *App) ReloadPlugins() {
// 	plugins.ReloadPlugins()
// }

// GetPlugin returns a plugin by name
// func (a *App) GetPlugin(pluginName string) *plugins.Plugin {
// 	plugins := plugins.GetPlugins()

// 	for _, plugin := range plugins {
// 		if plugin.Name == pluginName {
// 			return &plugin
// 		}
// 	}

// 	return nil
// }

// TooglePlugin will toogle a plugin by name
// func (a *App) TooglePlugin(pluginName string) {

// 	plugins := plugins.GetPlugins()

// 	for _, plugin := range plugins {
// 		if plugin.Name == pluginName {
// 			plugin.Toogle()
// 			break
// 		}
// 	}

// }

// InitClientPlugin will initialize the client logic of plugin by name
// func (a *App) InitClientPlugin(pluginName string) {

// 	messaging_port := plugins.MessagingPort.Get()

// 	// We use the next port for the websocket connection
// 	messaging_port_ws := messaging_port + 1

// 	plugins := plugins.GetPlugins()

// 	for _, plugin := range plugins {
// 		if plugin.Name == pluginName && plugin.IsEnabled() {
// 			plugin.Init_client(messaging_port_ws)
// 			break
// 		}
// 	}

// }

// func (a *App) GetNATSPort() uint16 {
// 	return plugins.MessagingPort.Get()
// }

// // ModifyPluginArgs will modify the arguments of a plugin
// func (a *App) ModifyPluginArgs(pluginName string, args []plugins.Plugin_arg, mode string) {
// 	plugins.ModifyArgs(pluginName, args, mode)
// }
