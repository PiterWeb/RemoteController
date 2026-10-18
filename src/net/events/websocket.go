// Websocket server
package events

import (
	"context"
	"io"
	"log"
	"net/http"

	"github.com/coder/websocket"
)

type topicT string
type addrStrT string

var conns = map[topicT]map[addrStrT]*websocket.Conn{}

func SetupWebsocketHandler(ctx context.Context, mux *http.ServeMux) {

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {

		c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
			OriginPatterns: []string{"localhost:5173"},
		})
		if err != nil {
			log.Println(err)
			return
		}

		defer func() {
			err := c.CloseNow()

			if err != nil {
				log.Println(err)
			}
		}()

		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		wsBroadcast(ctx, r, c)

		c.Close(websocket.StatusNormalClosure, "Client connection closed")

	})

}

func wsBroadcast(ctx context.Context, r *http.Request, ws *websocket.Conn) {

	addrStr := addrStrT(r.RemoteAddr)
	topic := topicT(r.URL.Path)

	if _, ok := conns[topic]; !ok {
		conns[topic] = map[addrStrT]*websocket.Conn{}
	}
	
	conns[topic][addrStr] = ws

	defer func() {

		// If panic it will recover and liberate resources
		_ = recover()

		for addr := range conns[topic] {
			if addrStr == addr {
				delete(conns[topic], addr)
				break
			}
		}

	}()

	for {
		typ, reader, err := ws.Reader(ctx)
		if err != nil {
			log.Println(err)
			break
		}

		for addr, con := range conns[topic] {

			if addrStr == addr {
				continue
			}

			writer, err := con.Writer(ctx, typ)

			if err != nil {
				log.Println(err)
				continue
			}

			_, err = io.Copy(writer, reader)

			if err != nil {
				log.Println(err)
				continue
			}

			// log.Println("Message sended to ", addr)

			err = writer.Close()

			if err != nil {
				log.Println(err)
				continue
			}

		}

	}
}
