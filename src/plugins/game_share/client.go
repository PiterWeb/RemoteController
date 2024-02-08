package game_share

import "errors"

func Client(protocol protocolType, port int) error {

	if protocol == tcp {
		return clientTCP(port)
	}

	if protocol == udp {
		return clientUDP(port)
	}

	return errors.New("Invalid protocol")

}
