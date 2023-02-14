package wss

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/rpc/serde"
	"nhooyr.io/websocket"
)

type serverWebSocketConnection struct {
	serveMux        http.ServeMux
	httpServer      *http.Server
	logger          zerolog.Logger
	requestHandlers map[serde.RequestMethod]func([]byte) []byte
	serverWebsocket *websocket.Conn
	port            string
}

type clientWebSocketConnection struct {
	logger               zerolog.Logger
	notificationHandlers map[serde.NotificationMethod]func([]byte)
	responseHandlers     map[uint64]chan ([]byte)
	clientWebsocket      *websocket.Conn
}

func NewWebSocketConnectionAsServer(port string) (*serverWebSocketConnection, error) {
	wsc := &serverWebSocketConnection{}
	wsc.requestHandlers = make(map[serde.RequestMethod]func([]byte) []byte)
	wsc.port = port
	wsc.serveMux.HandleFunc("/", wsc.subscribeRequestHandler)

	tcpListener, err := net.Listen("tcp", "127.0.0.1:"+wsc.port)
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

func (wsc *serverWebSocketConnection) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	wsc.serveMux.ServeHTTP(w, r)
}

func (wsc *serverWebSocketConnection) subscribeRequestHandler(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		panic(err)
	}
	wsc.serverWebsocket = c
	defer c.Close(websocket.StatusInternalError, "")

	err = wsc.readRequests(r.Context())
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

func (wsc *serverWebSocketConnection) readRequests(ctx context.Context) error {
	for {
		_, data, err := wsc.serverWebsocket.Read(ctx)
		if err != nil {
			return err
		}
		jsonRpc := serde.JsonRpcRequestAny[serde.RequestMethod]{}
		err = json.Unmarshal(data, &jsonRpc)
		if err != nil {
			return err
		}
		wsc.logger.Trace().Msgf("Received message: %v", jsonRpc)
		responseData := wsc.requestHandlers[jsonRpc.Method](data)
		err = wsc.serverWebsocket.Write(ctx, websocket.MessageText, responseData)
		if err != nil {
			return err
		}
	}
}

// Respond subscribes to a topic and calls the handler function when a message is received
// It returns an error if the subscription fails
// The handler processes the incoming data and returns the response data
func (wsc *serverWebSocketConnection) Respond(topic serde.RequestMethod, handler func([]byte) []byte) error {
	wsc.requestHandlers[topic] = handler
	return nil
}

func (wsc *serverWebSocketConnection) Notify(topic serde.NotificationMethod, data []byte) error {
	return wsc.serverWebsocket.Write(context.Background(), websocket.MessageText, data)
}

func (wsc *serverWebSocketConnection) Close() {
	err := wsc.httpServer.Shutdown(context.Background())
	if err != nil {
		panic(err)
	}
}

func (wsc *serverWebSocketConnection) Url() string {
	return "ws://127.0.0.1:" + wsc.port
}

func NewWebSocketConnectionAsClient(url string) (*clientWebSocketConnection, error) {
	wsc := &clientWebSocketConnection{}
	wsc.responseHandlers = make(map[uint64]chan ([]byte))
	wsc.notificationHandlers = make(map[serde.NotificationMethod]func([]byte))

	conn, _, err := websocket.Dial(context.Background(), url, &websocket.DialOptions{})
	if err != nil {
		return nil, err
	}
	wsc.clientWebsocket = conn
	go func() { wsc.readMessages(context.Background()) }()
	return wsc, nil
}

func (wsc *clientWebSocketConnection) readMessages(ctx context.Context) {
	for {
		_, data, err := wsc.clientWebsocket.Read(ctx)
		if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
			return
		}
		wsc.logger.Trace().Msgf("Received message: %s", string(data))

		// Is this a notification?
		unmarshaledData := map[string]any{}
		err = json.Unmarshal(data, &unmarshaledData)
		if err != nil {
			panic(err)
		}
		wsc.logger.Trace().Msgf("Received message: %v", unmarshaledData)

		// Is this a notification?
		if unmarshaledData["method"] != nil {
			unmarshaledNotifacation := serde.JsonRpcRequest[protocols.ObjectiveId]{}
			err = json.Unmarshal(data, &unmarshaledNotifacation)
			if err != nil {
				panic(err)
			}
			wsc.notificationHandlers[serde.NotificationMethod(unmarshaledNotifacation.Method)](data)
		} else {
			// Or is this a reply?
			unmarshaledResponse := serde.JsonRpcResponseAny{}
			err = json.Unmarshal(data, &unmarshaledResponse)
			if err != nil {
				panic(err)
			}
			wsc.responseHandlers[unmarshaledResponse.Id] <- data
		}
	}
}

func (wsc *clientWebSocketConnection) Request(method serde.RequestMethod, data []byte) ([]byte, error) {
	jsonRequest := serde.JsonRpcRequestAny[serde.RequestMethod]{}
	err := json.Unmarshal(data, &jsonRequest)
	if err != nil {
		return nil, err
	}
	responseChan := make(chan []byte, 1)
	wsc.responseHandlers[jsonRequest.Id] = responseChan
	go func() {
		responseChan <- <-responseChan
		delete(wsc.responseHandlers, jsonRequest.Id)
	}()

	err = wsc.clientWebsocket.Write(context.Background(), websocket.MessageText, data)
	if err != nil {
		return nil, err
	}

	return <-responseChan, nil
}

func (wsc *clientWebSocketConnection) Subscribe(topic serde.NotificationMethod, handler func([]byte)) error {
	wsc.notificationHandlers[topic] = handler
	return nil
}

func (wsc *clientWebSocketConnection) Close() {
	// Clients initiate and close websockets
	if wsc.clientWebsocket != nil {
		wsc.clientWebsocket.Close(websocket.StatusNormalClosure, "client initiated close")
	}
}
