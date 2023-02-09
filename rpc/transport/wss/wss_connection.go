package wss

import (
	"errors"
	"log"
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
	defer c.Close()

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
	go http.ListenAndServe(port, nil)
	return wsc

}

func (c *webSocketConnection) Request(method serde.RequestMethod, data []byte) ([]byte, error) {
	if c == nil {
		return []byte{}, errors.New("No websocket connection yet (not yet connected to server)")
	}
	var prefix byte
	switch method {
	case serde.DirectFundRequestMethod:
		prefix = 0
	case serde.DirectDefundRequestMethod:
		prefix = 1
	case serde.VirtualFundRequestMethod:
		prefix = 2
	case serde.VirtualDefundRequestMethod:
		prefix = 3
	case serde.PayRequestMethod:
		prefix = 4
	default:
		panic("unimplemented!")
	}

	data = append([]byte{prefix}, data...)

	err := c.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		log.Fatal(err)
	}

	return []byte{}, nil // TODO grab the "return value" of this message / request and return it

}

// Subscribe subscribes to a topic and calls the handler function when a message is received
// It returns an error if the subscription fails
// The handler processes the incoming data and returns the response data
func (c *webSocketConnection) Subscribe(topic serde.RequestMethod, handler func([]byte) []byte) error {
	if c == nil {
		return errors.New("No websocket connection yet (client not yet connected)")
	}

	listen := func() {

		for {
			var message []byte
			var err error
			var mt int
			mt, message, err = c.ReadMessage()

			if err != nil {
				log.Println(err)
			}
			switch message[0] {
			case 0:
				if topic != serde.DirectFundRequestMethod {
					continue
				}
			case 1:
				if topic != serde.DirectDefundRequestMethod {
					continue
				}
			case 2:
				if topic != serde.VirtualFundRequestMethod {
					continue
				}
			case 3:
				if topic != serde.VirtualDefundRequestMethod {
					continue
				}
			case 4:
				if topic != serde.PayRequestMethod {
					continue
				}

			}
			if mt == websocket.TextMessage {
				handler(message[1:])
			}
		}
	}
	go listen()
	return nil
}

func (c *webSocketConnection) Close() {
	c.Close() // TODO there is probably a more graceful protocol https://github.com/fasthttp/websocket/blob/master/_examples/echo/client.go
}
