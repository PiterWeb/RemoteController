package messaging

import (
	"fmt"
	"sync"

	"github.com/nats-io/nats.go"
)

var client_lock = &sync.Mutex{}
var clientInstance *nats.Conn

func Get_Client() *nats.Conn {
	if clientInstance != nil {
		client_lock.Lock()
		defer client_lock.Unlock()
		if clientInstance == nil {
			var err error
			clientInstance, err = nats.Connect(fmt.Sprintf("nats://localhost:%d", nats_port))

			if err != nil {
				return nil
			}
		} else if clientInstance.IsClosed() {
			var err error
			clientInstance, err = nats.Connect(fmt.Sprintf("nats://localhost:%d", nats_port))

			if err != nil {
				return nil
			}

		}
	}

	return clientInstance
}
