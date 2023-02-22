package ws

import (
	"context"
	"encoding/json"

	"github.com/rs/zerolog"
	"github.com/statechannels/go-nitro/rpc/serde"
	"nhooyr.io/websocket"
)

type clientWebSocketConnection struct {
	logger           zerolog.Logger
	notificationChan chan []byte
	responseHandlers map[uint64]chan ([]byte)
	clientWebsocket  *websocket.Conn
}

// NewWebSocketConnectionAsClient creates a websocket connection that can be used to send requests and listen for notifications
func NewWebSocketConnectionAsClient(url string) (*clientWebSocketConnection, error) {
	wsc := &clientWebSocketConnection{}
	wsc.responseHandlers = make(map[uint64]chan ([]byte))
	wsc.notificationChan = make(chan []byte)

	conn, _, err := websocket.Dial(context.Background(), url, &websocket.DialOptions{})
	if err != nil {
		return nil, err
	}
	wsc.clientWebsocket = conn
	go func() { wsc.readMessages(context.Background()) }()
	return wsc, nil
}

func (wsc *clientWebSocketConnection) Request(data []byte) ([]byte, error) {
	// Any request payload type will do here since we are only interested in the id
	responseChan := make(chan []byte, 1)
	unmarshaledRequest := serde.JsonRpcMessage{}
	err := json.Unmarshal(data, &unmarshaledRequest)
	if err != nil {
		return nil, err
	}
	wsc.responseHandlers[unmarshaledRequest.Id] = responseChan

	err = wsc.clientWebsocket.Write(context.Background(), websocket.MessageText, data)
	if err != nil {
		return nil, err
	}

	return <-responseChan, nil
}

func (wsc *clientWebSocketConnection) Subscribe() (<-chan []byte, error) {
	return wsc.notificationChan, nil
}

func (wsc *clientWebSocketConnection) Close() {
	// Clients initiate and close websockets{
	wsc.clientWebsocket.Close(websocket.StatusNormalClosure, "client initiated close")
}

func (wsc *clientWebSocketConnection) readMessages(ctx context.Context) {
	for {
		_, data, err := wsc.clientWebsocket.Read(ctx)
		if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
			return
		}
		wsc.logger.Trace().Msgf("Received message: %s", string(data))

		// Is this a notification?
		// Any payload type will do here since we are only interested in the method value
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
			wsc.responseHandlers[unmarshaledNotification.Id] <- data
			delete(wsc.responseHandlers, unmarshaledNotification.Id)
		}
	}
}
