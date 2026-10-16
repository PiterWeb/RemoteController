package oninit

import (
	"log"

	"LibreRemotePlay/src/cli"
	"LibreRemotePlay/src/state"

	LRPSignals "github.com/PiterWeb/LibreRemotePlaySignals/v1"
)

func Execute(app state.App) error {

	ips_channel := make(chan []string, 1)
	defer close(ips_channel)

	options := LRPSignals.ServerOptions{
		Port: cli.GetConfig().GetEasyConnectPort(),
	}

	log.Printf("Easy Connect Server started on port %d\n", options.Port)

	err := LRPSignals.InitServer(options, ips_channel)

	return err

}
