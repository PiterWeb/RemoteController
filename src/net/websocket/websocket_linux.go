package websocket

import (
	"net/http"

	"golang.org/x/net/websocket"
)

func SetupWebsocketHandler() {

	http.Handle("GET /ws", websocket.Handler(echoServer))

}

var conns = []*websocket.Conn{}

func echoServer(ws *websocket.Conn) {
	conns = append(conns, ws)

	defer func() {
		for i, con := range conns {
			if con.RemoteAddr().String() == ws.RemoteAddr().String() {
				conns = append(conns[:i], conns[i+1:]...)
				break
			}
		}
		ws.Close()
	}()

	for {
		msg := []byte{}
		if _, err := ws.Read(msg); err != nil {
			break
		}

		go func() {
			for _, con := range conns {
				if con.RemoteAddr().String() == ws.RemoteAddr().String() {
					continue
				}
				con.Write(msg)
			}
		}()
	}
}
