package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"os"

	"LibreRemotePlay/src/cli"
	_ "LibreRemotePlay/src/cli"
	appLogger "LibreRemotePlay/src/logger"
	"LibreRemotePlay/src/net/webrtc/streaming_signal"
	"LibreRemotePlay/src/oninit"
	"LibreRemotePlay/src/state"

	"fyne.io/systray"
	"fyne.io/systray/example/icon"
	"github.com/pkg/browser"
)

//go:embed frontend/build/*
var assets embed.FS

func main() {
	
	logFile := appLogger.InitLogger()

	defer logFile.Close()

	log.Println("LibreRemotePlay Starting app")

	ctx, cancelCtx := context.WithCancel(context.Background())

	app := state.App{
		Ctx: ctx,
		CancelCtx: cancelCtx,
		Assets: assets,
		Dev: os.Getenv("LRP_DEV") == "true",
	}
	
	go systray.Run(onSystrayReady(app), onSysTrayExit)
	defer systray.Quit()

	go func() {
		if err := oninit.Execute(app); err != nil {
			log.Println(err)
		}
	}()
	
	go func () {
		if err := streaming_signal.InitWhipServer(app, streaming_signal.WhipConfig); err != nil {
			log.Println(err)
		}
	}()
	
	<-ctx.Done()
}

func addQuitItem(app state.App) {
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")
	go func() {
		for  {
			select {
				case <-mQuit.ClickedCh: 
					log.Println("Requesting exit")
					app.CancelCtx()
				case <-app.Ctx.Done():
					return
			}
		}
	}()
}

func addOpenBrowserItem(app state.App) {
	mOpenBrowser := systray.AddMenuItem("Open", "Open LibreRemotePlay browser tab")
	go func() {
		for  {
			select {
				case <-mOpenBrowser.ClickedCh: 
					log.Println("Opening web browser")
					var url string

					if app.Dev {
						url = "http://localhost:5173"
					} else {
						url = fmt.Sprintf("http://localhost:%d", cli.GetConfig().GetHTTPPort())
					}
					
					err := browser.OpenURL(url)
					if err != nil {
						log.Println(err)
					}
				case <-app.Ctx.Done():
					return
			}
		}
	}()
}

func onSystrayReady(app state.App) func() {
	return func() {
		systray.SetIcon(icon.Data)
		systray.SetTitle("LibreRemotePlay")
		systray.SetTooltip("Desktop App")
		
		addQuitItem(app)
		addOpenBrowserItem(app)
	}
}

func onSysTrayExit() {}
