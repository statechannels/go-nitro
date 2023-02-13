package wss

import (
	"encoding/json"
	"errors"
	"log"
	"net"
	"net/http"

	"github.com/fasthttp/websocket"
	"github.com/statechannels/go-nitro/rpc/serde"
)

type webSocketConnection struct {
	*websocket.Conn
}

func NewWebSocketConnectionAsClient(url string) (*webSocketConnection, error) {
	c, _, err := websocket.DefaultDialer.Dial(url+"/wss", nil)
	if err != nil {
		return nil, err
	}
	return &webSocketConnection{c}, nil
}

func NewWebSocketConnectionAsServer(port string) *webSocketConnection {

	wsc := &webSocketConnection{}

	handshaker := func(w http.ResponseWriter, r *http.Request) {
		var upgrader = websocket.Upgrader{} // use default options
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}
		wsc.Conn = c
	}

	http.HandleFunc("/wss", handshaker)
	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(err)
	}
	go func() {
		err := http.Serve(l, nil)
		if err != nil {
			panic(err)
		}
	}()
	return wsc

}

func (c *webSocketConnection) Request(method serde.RequestMethod, data []byte) ([]byte, error) {
	jsonRequest := serde.JsonRpcRequestAny[serde.RequestMethod]{}
	err := json.Unmarshal(data, &jsonRequest)
	if err != nil {
		return nil, err
	}
	responseChan := make(chan []byte, 1)
	go func() {
		responseChan <- c.listenForResponse(jsonRequest.Id)
	}()

	err = send(c, method, data)
	if err != nil {
		return nil, err
	}

	return <-responseChan, nil
}

// Respond subscribes to a topic and calls the handler function when a message is received
// It returns an error if the subscription fails
// The handler processes the incoming data and returns the response data
func (c *webSocketConnection) Respond(topic serde.RequestMethod, handler func([]byte) []byte) error {
	if c == nil {
		return errors.New("no websocket connection yet (client not yet connected)")
	}

	listen := func() {
		for {
			var message []byte
			var mt int
			mt, message, err := c.ReadMessage()
			if err != nil {
				panic(err)
			}

			if mt == websocket.TextMessage {
				jsonRpc := serde.JsonRpcRequestAny[serde.RequestMethod]{}
				err := json.Unmarshal(message, &jsonRpc)
				if err != nil {
					panic(err)
				}

				if jsonRpc.Method == topic {
					response := handler(message)
					err = send(c, topic, response)
					if err != nil {
						panic(err)
					}
					return
				}
			}
		}
	}
	go listen()
	return nil
}

func (c *webSocketConnection) Notify(topic serde.NotificationMethod, data []byte) error {
	return send(c, topic, data)

}

func (c *webSocketConnection) Subscribe(topic serde.NotificationMethod, handler func([]byte)) error {
	if c == nil {
		return errors.New("no websocket connection yet (client not yet connected)")
	}

	listen := func() {
		for {
			var message []byte
			var mt int
			mt, message, err := c.ReadMessage()
			if err != nil {
				panic(err)
			}

			if mt == websocket.TextMessage {
				jsonRpc := serde.JsonRpcRequestAny[serde.NotificationMethod]{}
				err := json.Unmarshal(message, &jsonRpc)
				if err != nil {
					panic(err)
				}

				if jsonRpc.Method == topic {
					handler(message)
				}
			}
		}
	}
	go listen()
	return nil
}

func send[T serde.NotificationOrRequest](c *webSocketConnection, method T, data []byte) error {
	// TODO longer term, the interface should not contain "method", since "data" contains "method" anyway.
	if c == nil {
		return errors.New("no websocket connection yet (not yet connected to server)")
	}

	return c.WriteMessage(websocket.TextMessage, data)
}

func (c *webSocketConnection) listenForResponse(id uint64) []byte {
	for {
		var message []byte
		var mt int
		mt, message, err := c.ReadMessage()
		if err != nil {
			panic(err)
		}

		if mt == websocket.TextMessage {
			jsonRpc := serde.JsonRpcResponseAny{}
			err := json.Unmarshal(message, &jsonRpc)
			if err != nil {
				panic(err)
			}

			if jsonRpc.Id == id {
				return message
			}
		}
	}
}

func (c *webSocketConnection) Close() {
	c.Conn.Close() // TODO there is probably a more graceful protocol https://github.com/fasthttp/websocket/blob/master/_examples/echo/client.go
}

func (c *webSocketConnection) Url() string {
	return c.Conn.RemoteAddr().String()
}
