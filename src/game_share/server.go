package game_share

import (
	"errors"

	"github.com/pion/webrtc/v3"
)

func HandleServer(protocolChan <-chan ProtocolType, portChan <-chan int, d *webrtc.DataChannel) error {

	if d.Label() != "game_share_server" || d.Label() != "game_share_client" {
		return nil
	}

	protocol := <-protocolChan
	port := <-portChan

	if protocol == tcp {
		return serverTCP(port, d)
	}

	if protocol == udp {
		return serverUDP(port, d)
	}

	return errors.New("invalid protocol")

}
