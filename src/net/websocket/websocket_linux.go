package websocket

import (
	"fmt"

	"golang.org/x/net/websocket"
)

func InitWebsocketServer(wsPort int, clientPort int) error {

	// origin := fmt.Sprintf("http://localhost:%d/", clientPort)
	origin := "*"

	url := fmt.Sprintf("ws://localhost:%d/ws", wsPort)

	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		return err
	}

	defer ws.Close()

	msg := make([]byte, 512)

	// We have to handle the messages to do the signalization of the streaming
	for {
		if _, err := ws.Read(msg); err != nil {
			return err
		}

	}

}
