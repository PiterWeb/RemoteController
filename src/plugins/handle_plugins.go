package plugins

import (
	"fmt"
	"strings"

	ipc "github.com/james-barrow/golang-ipc"
	"github.com/pion/webrtc/v3"
)

func HandleServerPlugins(d *webrtc.DataChannel) {

	if !strings.Contains(d.Label(), "plugin:") {
		return
	}

	plugins := LoadPlugins()

	for _, p := range plugins {

		if d.Label() != "plugin:"+p.name {
			continue
		}

		var ipcClient *ipc.Client

		d.OnOpen(func() {

			p.Init_server()

			var err error

			ipcClient, err = ipc.StartClient("from_rmc_plugin:"+p.name, &ipc.ClientConfig{})

			if err != nil {
				fmt.Println(err)
			}

			ipcServer, err := ipc.StartServer("to_rmc_plugin:"+p.name, &ipc.ServerConfig{})

			if err != nil {
				fmt.Println(err)
			}

			for {

				msg, _ := ipcServer.Read()

				d.Send(msg.Data)
			}

		})

		d.OnMessage(func(msg webrtc.DataChannelMessage) {

			_ = ipcClient.Write(1, msg.Data)

		})

		break

	}

}
