package desktop

import (
	"context"

	"github.com/PiterWeb/RemoteController/src/net"
	"github.com/pion/webrtc/v3"
)

func createHost(ctx context.Context, ICEServers []webrtc.ICEServer, offerEncoded string, triggerEnd <-chan struct{}) string {

	answerResponse := make(chan string)

	go net.InitHost(ctx, ICEServers, offerEncoded, answerResponse, triggerEnd)

	return <-answerResponse

}
