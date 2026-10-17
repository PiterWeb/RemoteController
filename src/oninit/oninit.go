package oninit

import (
	"errors"
	"log"
	"net/http"

	"LibreRemotePlay/src/bindings"
	"LibreRemotePlay/src/cli"
	"LibreRemotePlay/src/net/assets"
	"LibreRemotePlay/src/net/webrtc/streaming_signal"
	"LibreRemotePlay/src/net/events"
	"LibreRemotePlay/src/state"

	LRPSignals "github.com/PiterWeb/LibreRemotePlaySignals/v1"
)

func Execute(app state.App) error {

	err := executePlatform(app)
	if err != nil {
		return err
	}
	
	httpServerMux := http.NewServeMux()

	bindingsServerMux := http.NewServeMux()
	httpServerMux.Handle("/bindings", bindingsServerMux)
	bindings.SetupBindings(app, bindingsServerMux)

	wsServerMux := http.NewServeMux()
	httpServerMux.Handle("/ws", wsServerMux)
	events.SetupWebsocketHandler(app.Ctx, wsServerMux)

	ips_channel := make(chan []string, 1)
	errChan := make(chan error, 2)

	go func() {
		err := assets.ServeHTTPAssets(app, httpServerMux)

		if err != nil {
			errChan <- err
		}
	}()
	
	go func() {

		options := LRPSignals.ServerOptions{
			Port: cli.GetConfig().GetEasyConnectPort(),
		}
		
		log.Printf("Easy Connect Server started on port %d\n", options.Port)
		err := LRPSignals.InitServer(options, ips_channel)
		if err != nil {
			errChan <- errors.Join(err, errors.New("easy connect server error"))
		}
	}()

	go func() {
		err := streaming_signal.InitWhipServer(app, streaming_signal.WhipConfig)
		if err != nil {
			errChan <- errors.Join(err, errors.New("whip server error"))
		}
	}()
	
	return <-errChan

}