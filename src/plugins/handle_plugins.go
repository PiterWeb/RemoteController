package plugins

import (
	"fmt"
	"strings"

	"github.com/PiterWeb/RemoteController/src/plugins/messaging"
	"github.com/nats-io/nats.go"
	"github.com/pion/webrtc/v3"
)

type NATS_PORT struct {
	value uint16
}

func (m NATS_PORT) Get() uint16 {
	return m.value
}

var MessagingPort NATS_PORT

func init() {
	MessagingPort = NATS_PORT{value: messaging.InitServer()}
}

func HandleServerPlugins(d *webrtc.DataChannel) {

	if !strings.Contains(d.Label(), "plugin:") {
		return
	}

	plugins := GetPlugins()

	messaging_client := messaging.Get_Client()

	for _, p := range plugins {

		if p.IsEnabled() {
			continue
		}

		if d.Label() != "plugin:"+p.Name {
			continue
		}

		d.OnOpen(func() {

			// Message comming from plugin to app
			messaging_client.Subscribe(p.Name+":app", func(msg *nats.Msg) {
				msg.Ack()
				d.Send(msg.Data)
			})

			// We pass the messaging port to the plugin so it can communicate with the app
			_, err := p.Init_host(MessagingPort.Get())

			if err != nil {
				fmt.Println(err)
			}

		})

		d.OnMessage(func(msg webrtc.DataChannelMessage) {

			// Send message from app to plugin
			messaging_client.Publish("app:"+p.Name, msg.Data)

		})

		break

	}

}
