package ws

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net"
	"net/http"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/statechannels/go-nitro/internal/safesync"
	"github.com/statechannels/go-nitro/rand"
)

const (
	httpServerAddress = "127.0.0.1:"
	maxRequestSize    = 8192
	apiVersionPath    = "/api/v1"
)

type serverHttpTransport struct {
	httpServer            *http.Server
	requestHandlers       map[string]func([]byte) []byte
	port                  string
	notificationListeners safesync.Map[chan []byte]
	logger                *slog.Logger

	wg *sync.WaitGroup
}

// NewHttpSocketTransportAsServer starts an http server
func NewHttpSocketTransportAsServer(port string) (*serverHttpTransport, error) {
	transport := &serverHttpTransport{port: port, notificationListeners: safesync.Map[chan []byte]{}, logger: slog.Default()}

	tcpListener, err := net.Listen("tcp", ":"+transport.port)
	if err != nil {
		return nil, err
	}

	var serveMux http.ServeMux

	serveMux.HandleFunc(apiVersionPath, transport.request)
	serveMux.HandleFunc(path.Join(apiVersionPath, "subscribe"), transport.subscribe)
	transport.httpServer = &http.Server{
		Handler:      &serveMux,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}

	transport.requestHandlers = make(map[string]func([]byte) []byte)
	transport.wg = &sync.WaitGroup{}

	transport.wg.Add(1)
	go transport.serveHttp(tcpListener)

	return transport, nil
}

func (t *serverHttpTransport) serveHttp(tcpListener net.Listener) {
	defer t.wg.Done()

	err := t.httpServer.Serve(tcpListener)
	if err != nil && errors.Is(err, http.ErrServerClosed) {
		return
	}
	if err != nil {
		panic(err)
	}
}

func (t *serverHttpTransport) RegisterRequestHandler(apiVersion string, handler func([]byte) []byte) error {
	t.requestHandlers[apiVersion] = handler
	return nil
}

func (t *serverHttpTransport) Notify(data []byte) error {
	t.notificationListeners.Range(func(key string, value chan []byte) bool {
		value <- data
		return true
	})
	return nil
}

func (t *serverHttpTransport) Close() error {
	// This will cause the serveHttp gand listenForClose goroutines to exit
	err := t.httpServer.Shutdown(context.Background())
	if err != nil {
		return err
	}

	t.wg.Wait()
	return nil
}

func (t *serverHttpTransport) Url() string {
	return httpServerAddress + t.port + apiVersionPath
}

func (t *serverHttpTransport) request(w http.ResponseWriter, r *http.Request) {
	// Pull api version from the url and determine if the version is supported
	pathSegments := strings.Split(r.URL.Path, "/")
	if len(pathSegments) < 3 {
		http.Error(w, "Invalid API version", http.StatusBadRequest)
		return
	}

	apiVersion := pathSegments[2] // first segment is an empty string
	handler, ok := t.requestHandlers[apiVersion]
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

var upgrader = websocket.Upgrader{} // use default options
func (t *serverHttpTransport) subscribe(w http.ResponseWriter, r *http.Request) {
	// TODO: We currently allow requests from any origins. We should probably use a whitelist.
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}

	defer c.Close()
	notificationChan := make(chan []byte)
	key := strconv.Itoa(int(rand.Uint64()))
	t.notificationListeners.Store(key, notificationChan)
	t.logger.Debug("Websocket transport added a notification listener")
	defer t.notificationListeners.Delete(key)

	closeChan := make(chan error)

	closeHandler := c.CloseHandler()
	c.SetCloseHandler(func(code int, text string) error {
		closeChan <- nil
		return closeHandler(code, text)
	})

EventLoop:
	for {
		select {
		case err = <-closeChan:
			break EventLoop
		case notificationData := <-notificationChan:
			err := c.WriteMessage(websocket.TextMessage, notificationData)
			if err != nil {
				break EventLoop
			}
		}
	}

	if err != nil {
		panic(err)
	}
}

// enableCors sets the CORS headers on the response allowing all origins
func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
