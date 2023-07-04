package ws

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	urlUtil "net/url"

	"github.com/rs/zerolog"
	"nhooyr.io/websocket"
)

type clientWebSocketTransport struct {
	logger           zerolog.Logger
	notificationChan chan []byte
	clientWebsocket  *websocket.Conn
	url              string
	ctx              context.Context
	cancel           context.CancelFunc
	closed           chan struct{}
}

// NewWebSocketTransportAsClient creates a websocket connection that can be used to send requests and listen for notifications
func NewWebSocketTransportAsClient(url string) (*clientWebSocketTransport, error) {
	wsc := &clientWebSocketTransport{}
	wsc.notificationChan = make(chan []byte)
	wsc.url = url

	subscribeUrl, err := urlUtil.JoinPath("ws://", url, "subscribe")
	if err != nil {
		return nil, err
	}
	conn, _, err := websocket.Dial(context.Background(), subscribeUrl, &websocket.DialOptions{})
	if err != nil {
		return nil, err
	}
	wsc.clientWebsocket = conn

	wsc.ctx, wsc.cancel = context.WithCancel(context.Background())
	wsc.closed = make(chan struct{})

	go wsc.readMessages()

	return wsc, nil
}

func (wsc *clientWebSocketTransport) Request(data []byte) ([]byte, error) {
	requestUrl, err := urlUtil.JoinPath("http://", wsc.url)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(requestUrl, "application/json", bytes.NewBuffer(data))
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

func (wsc *clientWebSocketTransport) Close() error {
	wsc.cancel()
	err := wsc.clientWebsocket.Close(websocket.StatusNormalClosure, "client initiated close")
	if err != nil {
		wsc.logger.Error().Err(err).Msgf("Error closing websocket %v", err)
	}
	<-wsc.closed
	return err
}

func (wsc *clientWebSocketTransport) readMessages() {
	defer close(wsc.closed)
	defer close(wsc.notificationChan)

	for {
		//_, data, err := wsc.clientWebsocket.Read(wsc.ctx)
		_, data, err := wsc.clientWebsocket.Read(context.Background())
		fmt.Printf("Received message: %s %v \n", string(data), err)
		if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
			break
		}
		if err != nil {
			wsc.logger.Error().Err(err).Msgf("Error reading from websocket %v", err)
			break
		}
		wsc.logger.Trace().Msgf("Received message: %s", string(data))
		wsc.notificationChan <- data
	}
}
