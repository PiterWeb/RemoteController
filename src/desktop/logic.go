package desktop

import (
	"github.com/PiterWeb/RemoteController/src/net"
)

func createHost(offerEncoded string, triggerEnd <-chan struct{}) string {

	answerResponse := make(chan string)

	go net.InitHost(offerEncoded, answerResponse, triggerEnd)

	return <-answerResponse

}
