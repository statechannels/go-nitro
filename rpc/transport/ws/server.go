package ws

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/statechannels/go-nitro/internal/safesync"
	"github.com/statechannels/go-nitro/rand"
	"nhooyr.io/websocket"
)

const webscocketServerAddress = "127.0.0.1:"

type serverWebSocketTransport struct {
	httpServer            *http.Server
	requestHandler        func([]byte) []byte
	port                  string
	notificationListeners safesync.Map[chan []byte]
}

// NewWebSocketTransportAsServer starts an http server that accepts websocket connections
func NewWebSocketTransportAsServer(port string) (*serverWebSocketTransport, error) {
	wsc := &serverWebSocketTransport{port: port, notificationListeners: safesync.Map[chan []byte]{}}

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
	wsc.notificationListeners.Range(func(key string, value chan []byte) bool {
		value <- data
		return true
	})
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

func (wsc *serverWebSocketTransport) subscribe(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		panic(err)
	}
	defer c.Close(websocket.StatusInternalError, "server initiated websocket close")
	notificationChan := make(chan []byte)
	key := strconv.Itoa(int(rand.Uint64()))
	wsc.notificationListeners.Store(key, notificationChan)
	defer wsc.notificationListeners.Delete(key)

	// A client closes a connection by sending a message over the websocket
	closeChan := make(chan error)
	go func() {
		_, _, err := c.Read(r.Context())
		closeChan <- err
	}()

	done := false
	for !done {
		select {
		case err = <-closeChan:
			done = true
		case <-r.Context().Done():
			err = r.Context().Err()
			done = true
		case notificationData := <-notificationChan:
			err := c.Write(r.Context(), websocket.MessageText, notificationData)
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
