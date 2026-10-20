// Websocket server
package events

import (
	"LibreRemotePlay/src/cli"
	"LibreRemotePlay/src/state"
	"context"
	"io"
	"log"
	"net/http"

	"github.com/coder/websocket"
)

type addrStrT string

var conns = map[topicT]map[addrStrT]*websocket.Conn{}

func SetupWebsocketHandler(app state.App, mux *http.ServeMux) {

	acceptOptions := &websocket.AcceptOptions{}
	
	if app.Dev {
		acceptOptions.OriginPatterns = []string{"localhost:5173"}
	}

	if cli.GetConfig().GetNetworkVisible() {
		acceptOptions.OriginPatterns = []string{}
		acceptOptions.InsecureSkipVerify = true	
	}
	
	mux.HandleFunc("GET /{topic}", func(w http.ResponseWriter, r *http.Request) {
		
		c, err := websocket.Accept(w, r, acceptOptions)
		
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

		ctx, cancel := context.WithCancel(app.Ctx)
		defer cancel()

		wsBroadcast(ctx, r, c)

		c.Close(websocket.StatusNormalClosure, "Client connection closed")

	})

}

func wsBroadcast(ctx context.Context, r *http.Request, ws *websocket.Conn) {

	addrStr := addrStrT(r.RemoteAddr)
	topic := topicT(r.PathValue("topic"))

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
