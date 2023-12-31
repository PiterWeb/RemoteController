package desktop

import (
	"context"
	"fmt"
	"runtime"
	"strings"
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

	return false
}

// Shutdown is called at application termination
func (a *App) Shutdown(ctx context.Context) {
	// Perform your teardown here
	a.TryClosePeerConnection()

}

// Create a Host Peer, it receives the offer encoded and returns the encoded answer response
func (a *App) TryCreateHost(offerEncoded string) (value string) {

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

	value = createHost(a.ctx, offerEncoded, triggerEnd)

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
