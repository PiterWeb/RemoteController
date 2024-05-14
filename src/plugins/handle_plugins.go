package plugins

import (
	"fmt"
	"strings"

	"github.com/PiterWeb/RemoteController/src/plugins/messaging"
	"github.com/nats-io/nats.go"
	"github.com/pion/webrtc/v3"
)

var messaging_port uint16
var nats_client *nats.Conn

func init() {
	messaging_port = messaging.InitServer()
}

func HandleServerPlugins(d *webrtc.DataChannel) {

	if !strings.Contains(d.Label(), "plugin:") {
		return
	}

	plugins := LoadPlugins()

	messaging_client := messaging.Get_Client()

	for _, p := range plugins {

		if d.Label() != "plugin:"+p.name {
			continue
		}

		plugin_port := 0

		d.OnOpen(func() {

			_, err := p.Init_host(uint16(messaging_port))

			// ipcClient, err = ipc.StartClient("from_rmc_plugin:"+p.name, &ipc.ClientConfig{})

			if err != nil {
				fmt.Println(err)
			}

			// ipcServer, err := ipc.StartServer("to_rmc_plugin:"+p.name, &ipc.ServerConfig{})

			if err != nil {
				fmt.Println(err)
			}

		})

		d.OnMessage(func(msg webrtc.DataChannelMessage) {

			if plugin_port == 0 {
				return
			}

			// _ = ipcClient.Write(1, msg.Data)

		})

		break

	}

}
