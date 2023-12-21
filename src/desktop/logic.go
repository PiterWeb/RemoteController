package desktop

import (
	"context"

	"github.com/PiterWeb/RemoteController/src/net"
)

func createHost(ctx context.Context, offerEncoded string, triggerEnd <-chan struct{}) string {

	answerResponse := make(chan string)

	go net.InitHost(ctx, offerEncoded, answerResponse, triggerEnd)

	return <-answerResponse

}
