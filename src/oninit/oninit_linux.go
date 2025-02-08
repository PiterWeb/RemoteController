package oninit

import (
	"embed"
	"log"

	"github.com/PiterWeb/RemoteController/src/net/http_assets"
	"github.com/PiterWeb/RemoteController/src/net/websocket"
)

func Execute(assets embed.FS) error {

	wsPort := 8080
	clientPort := 8081

	go func() {

		err := websocket.InitWebsocketServer(wsPort, clientPort)

		if err != nil {
			log.Println(err)
		}
	}()

	err := http_assets.InitHTTPAssets(clientPort, assets)

	if err != nil {
		return err
	}

	return nil

}
