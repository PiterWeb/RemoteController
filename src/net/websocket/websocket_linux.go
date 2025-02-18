package websocket

import (
	"context"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/coder/websocket"
)

var conns = map[string]*websocket.Conn{}

func SetupWebsocketHandler() {

	http.HandleFunc("GET /ws", func(w http.ResponseWriter, r *http.Request) {

		c, err := websocket.Accept(w, r, nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Set the context as needed. Use of r.Context() is not recommended
		// to avoid surprising behavior (see http.Hijacker).
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		wsBroadcast(ctx, r, c)

	})

}

func wsBroadcast(ctx context.Context, r *http.Request, ws *websocket.Conn) {
	conns[r.RemoteAddr] = ws

	defer func() {
		for addr := range conns {
			if r.RemoteAddr == addr {
				delete(conns, addr)
				break
			}
		}

		err := ws.CloseNow()

		if err != nil {
			log.Println(err)
		}
	}()

	for {
		typ, reader, err := ws.Reader(ctx)
		if err != nil {
			log.Println(err)
			return
		}

		for addr, con := range conns {
			go func() {

				if r.RemoteAddr == addr {
					return
				}

				writer, err := con.Writer(ctx, typ)

				if err != nil {
					log.Println(err)
					return
				}

				_, err = io.Copy(writer, reader)

				if err != nil {
					log.Println(err)
					return
				}
			}()
		}
	}
}
