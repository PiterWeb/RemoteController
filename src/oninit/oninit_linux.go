package oninit

import "github.com/PiterWeb/RemoteController/src/net/websocket"

func Execute() error {

	wsPort := 8080
	clientPort := 80

	err := websocket.InitWebsocketServer(wsPort, clientPort)

	if err != nil {
		return err
	}

	return nil

}
