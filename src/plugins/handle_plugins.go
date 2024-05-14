package plugins

import (
	"fmt"
	"strings"

	"github.com/PiterWeb/RemoteController/src/plugins/messaging"
	"github.com/nats-io/nats.go"
	"github.com/pion/webrtc/v3"
)

var messaging_port uint16

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

		d.OnOpen(func() {

			// Message comming from plugin to app
			messaging_client.Subscribe(p.name+":app", func(msg *nats.Msg) {
				msg.Ack()
				d.Send(msg.Data)
			})

			// We pass the messaging port to the plugin so it can communicate with the app
			_, err := p.Init_host(uint16(messaging_port))

			if err != nil {
				fmt.Println(err)
			}

		})

		d.OnMessage(func(msg webrtc.DataChannelMessage) {

			// Send message from app to plugin
			messaging_client.Publish("app:"+p.name, msg.Data)

		})

		break

	}

}
