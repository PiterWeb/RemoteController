// Binding for JS to Go
// This package is responsible for the communication between the JS and Go code.
package bindings

import (
	// "context"
	// "embed"

	//	"log"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	//	"strings"
	// "sync"
	// 	"time"
	//
	// 	"runtime"
	//
	// 	"LibreRemotePlay/src/cli"
	// 	"LibreRemotePlay/src/devices/audio"
	// 	"LibreRemotePlay/src/devices/gamepad"
	// 	"LibreRemotePlay/src/devices/keyboard"
	// 	"LibreRemotePlay/src/devices/mouse"
	net_webrtc "LibreRemotePlay/src/net/webrtc"
	// 	"LibreRemotePlay/src/net/webrtc/streaming_signal"
	// 	"LibreRemotePlay/src/onfinish"
	// 	"LibreRemotePlay/src/oninit"
	"LibreRemotePlay/src/net/webrtc/streaming_signal"
	"LibreRemotePlay/src/state"

	"github.com/pion/webrtc/v4"
	// "github.com/pion/webrtc/v4"
	// wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

type UsedPorts struct {
	HTTP uint16
	EasyConnect uint16
	WHIP uint16
}

// App struct
// type App struct {
// 	ctx    context.Context
// 	assets embed.FS 
// 	openPeer bool
// 	mutex *sync.Mutex
// 	pidAudioChan chan uint32
// 	cancelConn context.CancelFunc
// }

func SetupBindings(app state.App, mux *http.ServeMux) {

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Bindings OK"))
	})

	mux.Handle("PATCH /notify/open/peer", notifyOpenPeer(app))
	mux.Handle("PATCH /notify/close/peer", notifyClosePeer(app))

	mux.Handle("POST /conn/host", tryCreateHost(app))
}

func notifyOpenPeer(app state.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		app.Conn.Peer.Open()
		log.Println("NotifyCreateClient")
		w.WriteHeader(200)
	}
}

func notifyClosePeer(app state.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		app.Conn.Peer.Close()
		log.Println("NotifyCloseClient")
		w.WriteHeader(200)
	}
}

// Create a Host Peer, it receives the offer encoded and returns the encoded answer response
func tryCreateHost(app state.App) http.HandlerFunc {

	type Req struct {
		ICEServers []webrtc.ICEServer `json:"ice_servers"`
		OfferEncoded string `json:"offer_encoded"`
	}
	
	type Res struct {
		Response string `json:"data"`
	}
	
	return func(w http.ResponseWriter, r *http.Request) {
		// var ICEServers []webrtc.ICEServer
		// var offerEncoded string
		
		if app.Conn.Peer.IsOpen() {
			app.Conn.CancelCtx()
		}

		app.Conn.Peer.Open()
		
		answerResponse := make(chan string, 1) // Buffer of 1 to not block thread on cancel
	
		connCtx, cancelConn := context.WithCancel(app.Ctx)
		
		app.Conn.CancelCtx = cancelConn
		app.Conn.Ctx = connCtx

		rawReq, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println("Error reading req body: ", err)
			w.WriteHeader(500)
			return
		}

		var req Req
		err = json.Unmarshal(rawReq, &req)
		if err != nil {
			log.Println("Error unmarshalling req body: ", err)
			w.WriteHeader(400)
			return 
		}
		
		log.Println(req.ICEServers)
		streaming_signal.WhipConfig.ICEServers.Store(&req.ICEServers)

		pidAudioChan := make(chan uint32)
		
		go net_webrtc.InitHost(app.Ctx, connCtx, req.ICEServers, req.OfferEncoded, answerResponse, pidAudioChan)
	
		answer := <-answerResponse
	
		if answer == net_webrtc.ERROR_ANSWER {
			app.Conn.Peer.Close()
			log.Println("Error on WebRTC host connection")
			w.WriteHeader(500)
		}

		res := Res{
			Response: answer,
		}

		resRaw, err := json.Marshal(res)
		if err != nil {
			log.Println("Error marshaling response: ", err)
			w.WriteHeader(500)
			return
		}
		
		w.Write(resRaw)
	}
}

// // Closes the peer connection and returns a boolean indication if a connection existed and was closed or not
// func (a *App) TryClosePeerConnection() bool {
// 
// 	a.mutex.Lock()
// 	defer a.mutex.Unlock()
// 
// 	if !a.openPeer {
// 		return false
// 	}
// 
// 	a.cancelConn()
// 
// 	a.openPeer = false
// 
// 	return true
// 
// }

// // NewApp creates a new App application struct
// func NewApp(assets embed.FS) *App {
// 	return &App{
// 		ctx:    context.Background(),
// 		assets: assets,
// 		openPeer: false,
// 		mutex: &sync.Mutex{},
// 		pidAudioChan: make(chan uint32),
// 	}
// }
// 
// // Startup is called at application Startup
// func (a *App) Startup(ctx context.Context) {
// 	// Perform your setup here
// 
// 	a.ctx = ctx
// 	
// 	go func() {
// 		if err := oninit.Execute(a.ctx, a.assets); err != nil {
// 			time.Sleep(time.Second * 5) // Add some time to load the UI
// 			wailsRuntime.EventsEmit(ctx, "ERROR", err.Error())
// 			log.Println(err)
// 		}
// 	}()
// 	
// 	go func () {
// 		if err := streaming_signal.InitWhipServer(a.ctx, streaming_signal.WhipConfig); err != nil {
// 			time.Sleep(time.Second * 5) // Add some time to load the UI
// 			wailsRuntime.EventsEmit(ctx, "ERROR", err.Error())
// 			log.Println(err)
// 		}
// 	}()
// 
// }
// 
// // BeforeClose is called when the application is about to quit,
// // either by clicking the window close button or calling runtime.Quit.
// // Returning true will cause the application to continue, false will continue shutdown as normal.
// func (a *App) BeforeClose(ctx context.Context) (prevent bool) {
// 
// 	a.mutex.Lock()
// 	defer a.mutex.Unlock()
// 
// 	if !a.openPeer {
// 		prevent = false
// 		return prevent
// 	}
// 
// 	// Show a dialog to confirm the user wants to quit
// 	option, err := wailsRuntime.MessageDialog(ctx, wailsRuntime.MessageDialogOptions{
// 		Type:          wailsRuntime.QuestionDialog,
// 		Title:         "Quit",
// 		Message:       "Are you sure you want to quit?",
// 		Buttons:       []string{"Yes", "No"},
// 		DefaultButton: "No",
// 		CancelButton:  "No",
// 	})
// 
// 	if err != nil {
// 		return a.BeforeClose(ctx)
// 	}
// 
// 	if option == "Yes" {
// 		prevent = false
// 		return prevent
// 	}
// 
// 	prevent = true
// 	return prevent
// }
// 
// // Shutdown is called at application termination
// func (a *App) Shutdown(ctx context.Context) {
// 	// Perform your teardown here
// 	a.TryClosePeerConnection()
// 	if err := onfinish.Execute(); err != nil {
// 		log.Printf("Error onfinish: %s", err.Error())
// 	}
// }
// 

// 
// // Create a Host Peer, it receives the offer encoded and returns the encoded answer response
// func (a *App) TryCreateHost(ICEServers []webrtc.ICEServer, offerEncoded string) (value string) {
// 
// 	a.mutex.Lock()
// 
// 	if a.openPeer {
// 		a.cancelConn()
// 	}
// 
// 	a.openPeer = true
// 	
// 	answerResponse := make(chan string, 1) // Buffer of 1 to not block thread on cancel
// 
// 	connCtx, cancelConn := context.WithCancel(a.ctx)
// 	
// 	a.cancelConn = cancelConn
// 	a.mutex.Unlock()
// 	
// 	log.Println(ICEServers)
// 	streaming_signal.WhipConfig.ICEServers.Store(&ICEServers)
// 	
// 	go net.InitHost(a.ctx, connCtx, ICEServers, offerEncoded, answerResponse, a.pidAudioChan)
// 
// 	response := <-answerResponse
// 
// 	if response == net.ERROR_ANSWER {
// 		a.mutex.Lock()
// 		defer a.mutex.Unlock()
// 		a.openPeer = false
// 		log.Println("Error on WebRTC host connection")
// 	}
// 
// 	return response
// 
// }
// 
// // Closes the peer connection and returns a boolean indication if a connection existed and was closed or not
// func (a *App) TryClosePeerConnection() bool {
// 
// 	a.mutex.Lock()
// 	defer a.mutex.Unlock()
// 
// 	if !a.openPeer {
// 		return false
// 	}
// 
// 	a.cancelConn()
// 
// 	a.openPeer = false
// 
// 	return true
// 
// }
// 
// func (a *App) GetUsedPorts() UsedPorts {
// 	
// 	config := cli.GetConfig()
// 	
// 	return UsedPorts{
// 		HTTP: config.GetHTTPPort(),
// 		EasyConnect: config.GetEasyConnectPort(),
// 		WHIP: config.GetWhipServerPort(),
// 	}
// }
// 
// // Gamepad
// 
// func (a *App) ToogleGamepad() {
// 	gamepad.GamepadEnabled.Toogle()
// }
// 
// func (a *App) IsGamepadEnabled() bool {
// 	return gamepad.GamepadEnabled.IsEnabled()
// }
// 
// // Keyboard
// 
// func (a *App) ToogleKeyboard() {
// 	keyboard.KeyboardEnabled.Toogle()
// }
// 
// func (a *App) IsKeyboardEnabled() bool {
// 	return keyboard.KeyboardEnabled.IsEnabled()
// }
// 
// // Mouse
// 
// func (a *App) ToogleMouse() {
// 	mouse.MouseEnabled.Toogle()
// }
// 
// func (a *App) IsMouseEnabled() bool {
// 	return mouse.MouseEnabled.IsEnabled()
// }
// 
// // WHIP
// 
// func (a *App) ToogleWhip() {
// 	streaming_signal.WhipConfig.Enabled.Toogle()
// }
// 
// func (a *App) IsWhipEnabled() bool {
// 	return streaming_signal.WhipConfig.Enabled.IsEnabled()
// }
// 
// // OS Detection
// 
// func (a *App) GetCurrentOS() string {
// 	return strings.ToUpper(runtime.GOOS)
// }
// 
// // Log
// 
// func (a *App) LogPrintln(info string) {
// 	log.Println(info)
// }
// 
// // Audio
// 
// func (a *App) SetAudioPid(pid uint32) {
// 	a.pidAudioChan <- pid
// }
// 
// func (a *App) GetAudioProcess() []audio.AudioProcess {
// 	return audio.GetAudioProcess()
// }
