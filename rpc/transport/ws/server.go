package ws

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/statechannels/go-nitro/internal/safesync"
	"github.com/statechannels/go-nitro/rand"
	"nhooyr.io/websocket"
)

const (
	websocketServerAddress = "127.0.0.1:"
	maxRequestSize         = 8192
	apiVersionPath         = "/api/v1"
)

type serverWebSocketTransport struct {
	httpServer            *http.Server
	requestHandlers       map[string]func([]byte) []byte
	port                  string
	notificationListeners safesync.Map[chan []byte]
}

// NewWebSocketTransportAsServer starts an http server that accepts websocket connections
func NewWebSocketTransportAsServer(port string) (*serverWebSocketTransport, error) {
	wsc := &serverWebSocketTransport{port: port, notificationListeners: safesync.Map[chan []byte]{}}

	tcpListener, err := net.Listen("tcp", websocketServerAddress+wsc.port)
	if err != nil {
		return nil, err
	}

	var serveMux http.ServeMux

	serveMux.HandleFunc(apiVersionPath, wsc.request)
	serveMux.HandleFunc(path.Join(apiVersionPath, "subscribe"), wsc.subscribe)
	wsc.httpServer = &http.Server{
		Handler:      &serveMux,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}

	wsc.requestHandlers = make(map[string]func([]byte) []byte)

	go func() {
		err = wsc.httpServer.Serve(tcpListener)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	return wsc, nil
}

func (wsc *serverWebSocketTransport) RegisterRequestHandler(apiVersion string, handler func([]byte) []byte) error {
	wsc.requestHandlers[apiVersion] = handler
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
	return websocketServerAddress + wsc.port + apiVersionPath
}

func (wsc *serverWebSocketTransport) request(w http.ResponseWriter, r *http.Request) {
	// Pull api version from the url and determine if the version is supported
	pathSegments := strings.Split(r.URL.Path, "/")
	if len(pathSegments) < 3 {
		http.Error(w, "Invalid API version", http.StatusBadRequest)
		return
	}

	apiVersion := pathSegments[2] // first segment is an empty string
	handler, ok := wsc.requestHandlers[apiVersion]
	if !ok {
		http.Error(w, "Invalid API version", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "OPTIONS": // OPTIONS is used for a pre-flight CORS check by the browser before POST
		enableCors(&w)
		// This header value indicates which request headers can be used when making the actual request.
		// See https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Headers
		w.Header().Set("Access-Control-Allow-Headers", "*")
	case "POST":
		enableCors(&w)
		body := http.MaxBytesReader(w, r.Body, maxRequestSize)
		msg, err := io.ReadAll(body)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusRequestEntityTooLarge), http.StatusRequestEntityTooLarge)
			return
		}
		_, err = w.Write(handler(msg))
		if err != nil {
			panic(err)
		}
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (wsc *serverWebSocketTransport) subscribe(w http.ResponseWriter, r *http.Request) {
	// TODO: We currently allow requests from any origins. We should probably use a whitelist.
	opts := &websocket.AcceptOptions{InsecureSkipVerify: true}
	c, err := websocket.Accept(w, r, opts)
	if err != nil {
		panic(err)
	}
	defer c.Close(websocket.StatusInternalError, "server initiated websocket close")
	notificationChan := make(chan []byte)
	key := strconv.Itoa(int(rand.Uint64()))
	wsc.notificationListeners.Store(key, notificationChan)
	defer wsc.notificationListeners.Delete(key)

	// A client closes a connection by sending a message over the websocket
	// This code converts the socket `Read` call to a channel.
	// Ideally, all signals would be managed in the select statement below. But there is no channel API for the websocket read.
	closeChan := make(chan error)
	go func() {
		_, _, err := c.Read(r.Context())
		closeChan <- err
	}()

EventLoop:
	for {
		select {
		case err = <-closeChan:
			break EventLoop
		case <-r.Context().Done():
			err = r.Context().Err()
			break EventLoop
		case notificationData := <-notificationChan:
			err := c.Write(r.Context(), websocket.MessageText, notificationData)
			if err != nil {
				break EventLoop
			}
		}
	}
	if errors.Is(err, context.Canceled) {
		return
	}
	if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
		websocket.CloseStatus(err) == websocket.StatusGoingAway {
		return
	}
	if err != nil {
		panic(err)
	}
}

// enableCors sets the CORS headers on the response allowing all origins
func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
