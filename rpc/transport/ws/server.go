package ws

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"time"

	"nhooyr.io/websocket"
)

const webscocketServerAddress = "127.0.0.1:"

type serverWebSocketTransport struct {
	httpServer       *http.Server
	requestHandler   func([]byte) []byte
	serverWebsocket  *websocket.Conn
	port             string
	notificationChan chan []byte
}

// NewWebSocketTransportAsServer starts an http server that accepts websocket connections
func NewWebSocketTransportAsServer(port string) (*serverWebSocketTransport, error) {
	wsc := &serverWebSocketTransport{port: port, notificationChan: make(chan []byte)}

	tcpListener, err := net.Listen("tcp", webscocketServerAddress+wsc.port)
	if err != nil {
		return nil, err
	}

	var serveMux http.ServeMux
	serveMux.HandleFunc("/", wsc.request)
	serveMux.HandleFunc("/subscribe", wsc.subscribe)
	wsc.httpServer = &http.Server{
		Handler:      &serveMux,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}

	go func() {
		err = wsc.httpServer.Serve(tcpListener)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	return wsc, nil
}

func (wsc *serverWebSocketTransport) RegisterRequestHandler(handler func([]byte) []byte) error {
	wsc.requestHandler = handler
	return nil
}

func (wsc *serverWebSocketTransport) Notify(data []byte) error {
	wsc.notificationChan <- data
	return nil
}

func (wsc *serverWebSocketTransport) Close() {
	err := wsc.httpServer.Shutdown(context.Background())
	if err != nil {
		panic(err)
	}
}

func (wsc *serverWebSocketTransport) Url() string {
	return "ws://" + webscocketServerAddress + wsc.port
}

func (wsc *serverWebSocketTransport) request(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	body := http.MaxBytesReader(w, r.Body, 8192)
	msg, err := io.ReadAll(body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusRequestEntityTooLarge), http.StatusRequestEntityTooLarge)
		return
	}

	_, err = w.Write(wsc.requestHandler(msg))
	if err != nil {
		panic(err)
	}
}

// Todo: should allow for multiple notification subscribers
func (wsc *serverWebSocketTransport) subscribe(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		panic(err)
	}
	wsc.serverWebsocket = c
	defer c.Close(websocket.StatusInternalError, "server initiated websocket close")

	done := false
	for !done {
		select {
		case <-r.Context().Done():
			err = r.Context().Err()
			done = true
		case notificationData := <-wsc.notificationChan:
			err := wsc.serverWebsocket.Write(r.Context(), websocket.MessageText, notificationData)
			if err != nil {
				done = true
			}
		}
	}
	if errors.Is(err, context.Canceled) {
		return
	}
	if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
		return
	}
	if err != nil {
		panic(err)
	}
}
