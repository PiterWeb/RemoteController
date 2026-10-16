package oninit

import (
	"log"
	"net/http"

	"LibreRemotePlay/src/cli"
	"LibreRemotePlay/src/devices/gamepad"
	net_http "LibreRemotePlay/src/net/http"
	"LibreRemotePlay/src/state"

	LRPSignals "github.com/PiterWeb/LibreRemotePlaySignals/v1"
)

func Execute(app state.App) error {
	err := gamepad.InitViGEm()

	if err != nil {
		return err
	}

	httpServerMux := http.NewServeMux()

	ips_channel := make(chan []string, 1)
	errChan := make(chan error, 2)

	go func() {
		err := net_http.InitHTTPAssets(app, httpServerMux)

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
			errChan <- err
		}
	}()

	return <-errChan

}
