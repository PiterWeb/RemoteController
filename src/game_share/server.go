package game_share

import (
	"errors"
)

func Server(protocol protocolType, port int) error {

	if protocol == tcp {
		return serverTCP(port)
	}

	if protocol == udp {
		return serverUDP(port)
	}

	return errors.New("invalid protocol")

}
