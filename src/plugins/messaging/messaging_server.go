//go:build !darwin

package messaging

import (
	"net"

	"github.com/nats-io/nats-server/v2/server"
)

var nats_server *server.Server
var nats_port uint16

func ShutdownServer() {
	if nats_server != nil && nats_server.Running() {
		nats_server.Shutdown()
	}
}

func getFreeTCPPort() (port int, err error) {
	var a *net.TCPAddr
	if a, err = net.ResolveTCPAddr("tcp", "127.0.0.1:0"); err == nil {
		var l *net.TCPListener
		if l, err = net.ListenTCP("tcp", a); err == nil {
			defer l.Close()
			return l.Addr().(*net.TCPAddr).Port, nil
		}
	}
	return getFreeTCPPort()
}

func InitServer() (port uint16) {

	free_port, _ := getFreeTCPPort()

	opts := &server.Options{
		Port: free_port,
		Websocket: server.WebsocketOpts{
			Port:  free_port + 1,
			NoTLS: true,
		},
	}

	ns, err := server.NewServer(opts)

	if err != nil {
		panic(err)
	}

	ns.Start()

	nats_server = ns

	nats_port = uint16(free_port)

	return uint16(free_port)

}
