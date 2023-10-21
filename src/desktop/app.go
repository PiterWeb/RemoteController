package desktop

import (
	"context"
)

var triggerEnd chan struct{} = make(chan struct{})
var processActive bool = false

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
}

func (a *App) CreateClient() string {

	if processActive {
		triggerEnd <- struct{}{}
	}

	processActive = true

	var value string

	defer func() {

		if err := recover(); err != nil {
			value = "ERROR"
			processActive = false
		}

	}()

	value = createClient(a.ctx, triggerEnd)

	return value
}

func (a *App) CreateHost(offerEncoded string) string {

	if processActive {
		triggerEnd <- struct{}{}
	}

	processActive = true

	var value string

	defer func() {

		if err := recover(); err != nil {
			value = "ERROR"
			processActive = false
		}

	}()

	value = createHost(a.ctx, offerEncoded, triggerEnd)

	return value
}

func (a *App) ConnectToHost(response string) string {

	var value string

	defer func() {

		if err := recover(); err != nil {
			value = "ERROR"
			processActive = false
		}

	}()

	value = "OK"
	connectToHost(response)

	return value

}

func (a *App) CloseConnection() bool {

	if processActive {
		triggerEnd <- struct{}{}
	} else {
		return false
	}

	processActive = false

	return true

}
