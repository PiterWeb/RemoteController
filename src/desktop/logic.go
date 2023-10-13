package desktop

import (
	"github.com/PiterWeb/RemoteController/src/net"
)

func connectToHost(offerEncoded string) string {

	answerResponse := make(chan string)

	go net.InitAnswer(offerEncoded, answerResponse)

	return <-answerResponse

}

var answerResponse = make(chan string)

func createHost() string {

	offer := make(chan string)

	go net.InitOffer(offer, answerResponse)

	return <-offer

}

func connectToClient(response string) {

	answerResponse <- response

}
