package desktop

import (
	"context"

	"github.com/PiterWeb/RemoteController/src/customctx"
)

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

// DomReady is called after front-end resources have been loaded
func (a App) DomReady(ctx context.Context) {
	customctx.DomReadyCtx = ctx
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
}

func (a *App) ConnectToHost(offerEncoded string) string {

	var value string

	defer func() {

		if err := recover(); err != nil {
			value = "ERROR"
		}

	}()

	value = connectToHost(offerEncoded)

	return value
}

func (a *App) CreateHost() string {

	var value string

	defer func() {

		if err := recover(); err != nil {
			value = "ERROR"
		}

	}()

	value = createHost()

	return value
}

func (a *App) ConnectToClient(response string) string {

	var value string

	defer func() {

		if err := recover(); err != nil {
			value = "ERROR"
		}

	}()

	value = "OK"
	connectToClient(response)

	return value

}
