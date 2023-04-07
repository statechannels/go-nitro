package ws

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/statechannels/go-nitro/internal/safesync"
	"github.com/statechannels/go-nitro/rpc/serde"
	"nhooyr.io/websocket"
)

type clientWebSocketTransport struct {
	logger           zerolog.Logger
	notificationChan chan []byte
	responseHandlers safesync.Map[chan []byte]
	clientWebsocket  *websocket.Conn
}

// NewWebSocketTransportAsClient creates a websocket connection that can be used to send requests and listen for notifications
func NewWebSocketTransportAsClient(url string) (*clientWebSocketTransport, error) {
	wsc := &clientWebSocketTransport{}
	wsc.responseHandlers = safesync.Map[chan []byte]{}
	wsc.notificationChan = make(chan []byte)

	conn, _, err := websocket.Dial(context.Background(), url, &websocket.DialOptions{})
	if err != nil {
		return nil, err
	}
	wsc.clientWebsocket = conn
	go func() { wsc.readMessages(context.Background()) }()
	return wsc, nil
}

func (wsc *clientWebSocketTransport) Request(data []byte) ([]byte, error) {
	responseChan := make(chan []byte, 1)
	unmarshaledRequest := serde.JsonRpcMessage{}
	err := json.Unmarshal(data, &unmarshaledRequest)
	if err != nil {
		return nil, err
	}
	wsc.responseHandlers.Store(fmt.Sprintf("%v", unmarshaledRequest.Id), responseChan)

	err = wsc.clientWebsocket.Write(context.Background(), websocket.MessageText, data)
	if err != nil {
		return nil, err
	}

	return <-responseChan, nil
}

func (wsc *clientWebSocketTransport) Subscribe() (<-chan []byte, error) {
	return wsc.notificationChan, nil
}

func (wsc *clientWebSocketTransport) Close() {
	// Clients initiate and close websockets{
	wsc.clientWebsocket.Close(websocket.StatusNormalClosure, "client initiated close")
	close(wsc.notificationChan)
}

func (wsc *clientWebSocketTransport) readMessages(ctx context.Context) {
	for {
		_, data, err := wsc.clientWebsocket.Read(ctx)
		if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
			return
		}
		wsc.logger.Trace().Msgf("Received message: %s", string(data))

		unmarshaledNotification := serde.JsonRpcMessage{}
		err = json.Unmarshal(data, &unmarshaledNotification)
		if err != nil {
			panic(err)
		}
		wsc.logger.Trace().Msgf("Received message: %v", unmarshaledNotification)

		// Is this a notification?
		if unmarshaledNotification.Method != "" {
			wsc.notificationChan <- data
			// Or is this a response?
		} else {
			sId := fmt.Sprintf("%v", unmarshaledNotification.Id)
			responseHandler, ok := wsc.responseHandlers.Load(sId)
			if !ok {
				panic(fmt.Errorf("Expected a response handler for id %v", unmarshaledNotification.Id))
			}
			responseHandler <- data
			wsc.responseHandlers.Delete(sId)
		}
	}
}
