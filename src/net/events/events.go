package events

import (
	"LibreRemotePlay/src/cli"
	"context"
	"fmt"
	"log"

	"github.com/coder/websocket"
)

type topicT string

const (
	EVENT_CONNECTION_STATE topicT = "connection_state"
	EVENT_ERROR topicT = "ERROR" 
	EVENT_STREAMING_SIGNAL_CLIENT topicT = "streaming-signal-client"
	EVENT_STREAMING_SIGNAL_SERVER topicT = "streaming-signal-server"
)

func EventsEmit(ctx context.Context, eventName topicT, data []byte) {

	wsUrl := fmt.Sprintf("ws://localhost:%d/ws/%s", cli.GetConfig().GetHTTPPort(), eventName)

	wsClient, _, err := websocket.Dial(ctx, wsUrl, nil)
	if err != nil {
		log.Println(err)
		return
	}

	wsClient.Write(ctx, websocket.MessageText, data)
	
}