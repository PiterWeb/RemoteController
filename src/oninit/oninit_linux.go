package oninit

import (
	"context"
	"embed"
	"errors"
	"log"
	"net/http"

	LRPSignals "github.com/PiterWeb/LibreRemotePlaySignals/v1"
	"LibreRemotePlay/src/cli"
	net_http "LibreRemotePlay/src/net/http"
	"LibreRemotePlay/src/net/websocket"
)

func Execute(ctx context.Context, assets embed.FS) error {

	httpServerMux := http.NewServeMux()
	
	websocket.SetupWebsocketHandler(ctx, httpServerMux)

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
			errChan <- errors.Join(err, errors.New("easy connect server error"))
		}
	}()

	return <-errChan

}
