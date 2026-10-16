package oninit

import (
	"context"
	"embed"
	"log"
	"net/http"

	LRPSignals "github.com/PiterWeb/LibreRemotePlaySignals/v1"
	"LibreRemotePlay/src/cli"
	"LibreRemotePlay/src/devices/gamepad"
	net_http "LibreRemotePlay/src/net/http"
)

func Execute(ctx context.Context, assets embed.FS) error {
	err := gamepad.InitViGEm()

	if err != nil {
		return err
	}

	httpServerMux := http.NewServeMux()

	ips_channel := make(chan []string, 1)
	errChan := make(chan error, 2)

	go func() {
		err := net_http.InitHTTPAssets(ctx, httpServerMux, assets)

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
