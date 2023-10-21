package desktop

import (
	"context"

	"github.com/PiterWeb/RemoteController/src/net"
)

var answerResponse = make(chan string)

func createHost(ctx context.Context, offerEncoded string, triggerEnd <-chan struct{}) string {

	answerResponse := make(chan string)

	go net.InitHost(ctx, offerEncoded, answerResponse, triggerEnd)

	return <-answerResponse

}

func createClient(ctx context.Context, triggerEnd <-chan struct{}) string {

	offer := make(chan string)

	go net.InitClient(ctx, offer, answerResponse, triggerEnd)

	return <-offer

}

func connectToHost(response string) {

	answerResponse <- response

}
