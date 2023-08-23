package ws

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	urlUtil "net/url"
	"sync"

	"github.com/gorilla/websocket"
)

type clientWebSocketTransport struct {
	logger           *slog.Logger
	notificationChan chan []byte
	clientWebsocket  *websocket.Conn
	url              string
	wg               *sync.WaitGroup
}

// NewWebSocketTransportAsClient creates a websocket connection that can be used to send requests and listen for notifications
func NewWebSocketTransportAsClient(url string) (*clientWebSocketTransport, error) {
	subscribeUrl, err := urlUtil.JoinPath("ws://", url, "subscribe")
	if err != nil {
		return nil, err
	}
	conn, _, err := websocket.DefaultDialer.Dial(subscribeUrl, nil)
	if err != nil {
		return nil, err
	}

	wsc := &clientWebSocketTransport{notificationChan: make(chan []byte, 10), clientWebsocket: conn, url: url, wg: &sync.WaitGroup{}, logger: slog.Default()}

	wsc.wg.Add(1)
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
	// This will also cause the go-routine to unblock waiting on `ReadMessage` and thus serves as a signal to exit
	err := wsc.clientWebsocket.Close()
	if err != nil {
		return err
	}
	wsc.wg.Wait()

	close(wsc.notificationChan)
	return nil
}

func (wsc *clientWebSocketTransport) readMessages() {
	wsc.logger.Debug("Starting to read websocket messages")
	for {
		_, data, err := wsc.clientWebsocket.ReadMessage()
		if err != nil {
			wsc.logger.Info("Websocket read error", "error", err)
			wsc.wg.Done()
			return
		}
		wsc.logger.Debug("Websocket received message", "data", string(data))

		wsc.notificationChan <- data
	}
}
