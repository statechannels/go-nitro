package ws

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/rs/zerolog"
	"nhooyr.io/websocket"
)

type clientWebSocketTransport struct {
	logger           zerolog.Logger
	notificationChan chan []byte
	clientWebsocket  *websocket.Conn
	url              string
}

// NewWebSocketTransportAsClient creates a websocket connection that can be used to send requests and listen for notifications
func NewWebSocketTransportAsClient(url string) (*clientWebSocketTransport, error) {
	wsc := &clientWebSocketTransport{}
	wsc.notificationChan = make(chan []byte)
	wsc.url = url

	conn, _, err := websocket.Dial(context.Background(), url+"/subscribe", &websocket.DialOptions{})
	if err != nil {
		return nil, err
	}
	wsc.clientWebsocket = conn
	go func() { wsc.readMessages(context.Background()) }()
	return wsc, nil
}

func (wsc *clientWebSocketTransport) Request(data []byte) ([]byte, error) {
	resp, err := http.Post("http"+wsc.url[2:], "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
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
		wsc.notificationChan <- data
	}
}
