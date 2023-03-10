package ws

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"nhooyr.io/websocket"
)

const webscocketServerAddress = "127.0.0.1:"

type serverWebSocketTransport struct {
	serveMux        http.ServeMux
	httpServer      *http.Server
	logger          zerolog.Logger
	requestHandler  func([]byte) []byte
	serverWebsocket *websocket.Conn
	port            string
}

// NewWebSocketTransportAsServer starts an http server that accepts websocket connections
func NewWebSocketTransportAsServer(port string) (*serverWebSocketTransport, error) {
	wsc := &serverWebSocketTransport{}
	wsc.port = port
	wsc.serveMux.HandleFunc("/", wsc.subscribeRequestHandler)

	tcpListener, err := net.Listen("tcp", webscocketServerAddress+wsc.port)
	if err != nil {
		return nil, err
	}

	wsc.httpServer = &http.Server{
		Handler:      wsc,
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

// ServeHTTP is a required method for the http.Handler interface
func (wsc *serverWebSocketTransport) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	wsc.serveMux.ServeHTTP(w, r)
}

func (wsc *serverWebSocketTransport) Respond(handler func([]byte) []byte) error {
	wsc.requestHandler = handler
	return nil
}

func (wsc *serverWebSocketTransport) Notify(data []byte) error {
	return wsc.serverWebsocket.Write(context.Background(), websocket.MessageText, data)
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

func (wsc *serverWebSocketTransport) subscribeRequestHandler(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		panic(err)
	}
	wsc.serverWebsocket = c
	defer c.Close(websocket.StatusInternalError, "server initiated websocket close")

	err = wsc.readRequests(r.Context())
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

func (wsc *serverWebSocketTransport) readRequests(ctx context.Context) error {
	for {
		_, data, err := wsc.serverWebsocket.Read(ctx)
		if err != nil {
			return err
		}
		wsc.logger.Trace().Msgf("Received message: %v", data)
		responseData := wsc.requestHandler(data)
		err = wsc.serverWebsocket.Write(ctx, websocket.MessageText, responseData)
		if err != nil {
			return err
		}
	}
}
