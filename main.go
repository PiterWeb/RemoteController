package main

import (
	"embed"
	"log"

	"github.com/PiterWeb/RemoteController/src/desktop"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed frontend/build/*
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := desktop.NewApp()

	// Create application with options
	// Create application with options
	err := wails.Run(&options.App{
		Title:             "Remote Controller",
		Width:             1024,
		Height:            768,
		DisableResize:     false,
		Fullscreen:        false,
		StartHidden:       false,
		HideWindowOnClose: false,
		BackgroundColour:  &options.RGBA{R: 75, G: 107, B: 251, A: 255},
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		Menu:                     nil,
		Logger:                   nil,
		LogLevel:                 logger.DEBUG,
		OnStartup:                app.Startup,
		OnBeforeClose:            app.BeforeClose,
		OnShutdown:               app.Shutdown,
		WindowStartState:         options.Normal,
		EnableDefaultContextMenu: true,
		Bind: []interface{}{
			app,
		},
		// Windows platform specific options
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			DisableWindowIcon:    true,
			Theme:                windows.Theme(windows.Acrylic),
			// DisableFramelessWindowDecorations: false,
			WebviewUserDataPath: "",
		},
	})

	if err != nil {
		log.Println(err)
	}
}
