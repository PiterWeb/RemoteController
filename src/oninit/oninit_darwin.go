package oninit

import (
	"context"
	"embed"
	"log"

	LRPSignals "github.com/PiterWeb/LibreRemotePlaySignals/v1"
	"LibreRemotePlay/src/cli"
)

func Execute(ctx context.Context, assets embed.FS) error {

	ips_channel := make(chan []string, 1)
	defer close(ips_channel)

	options := LRPSignals.ServerOptions{
		Port: cli.GetConfig().GetEasyConnectPort(),
	}

	log.Printf("Easy Connect Server started on port %d\n", options.Port)

	err := LRPSignals.InitServer(options, ips_channel)

	return err

}
