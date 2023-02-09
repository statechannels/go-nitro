package wss

import (
	"errors"
	"log"
	"net/http"

	"github.com/fasthttp/websocket"
	"github.com/statechannels/go-nitro/rpc/serde"
)

type webSocketConnection struct {
	createLedgerChannelConnection *websocket.Conn // we will have one connection per method initially, and look to combine them later on.
}

func NewWebSocketConnectionAsClient(url string) *webSocketConnection {
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()
	return &webSocketConnection{c}
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
		wsc.createLedgerChannelConnection = c
	}

	http.HandleFunc("/wss", handshaker)
	go http.ListenAndServe(port, nil)
	return wsc

}

func (c *webSocketConnection) Request(method serde.RequestMethod, data []byte) ([]byte, error) {
	if c.createLedgerChannelConnection == nil {
		return []byte{}, errors.New("No websocket connection yet (not yet connected to server)")
	}
	switch method {
	case serde.DirectFundRequestMethod:
		err := c.createLedgerChannelConnection.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			log.Fatal(err)
		}
	default:
		panic("unimplemented!")
	}

	return []byte{}, nil // TODO grab the "return value" of this message / request and return it

}

// Subscribe subscribes to a topic and calls the handler function when a message is received
// It returns an error if the subscription fails
// The handler processes the incoming data and returns the response data
func (c *webSocketConnection) Subscribe(topic serde.RequestMethod, handler func([]byte) []byte) error {
	if c.createLedgerChannelConnection == nil {
		return errors.New("No websocket connection yet (client not yet connected)")
	}

	listen := func() {

		for {
			var message []byte
			var err error
			var mt int
			switch topic {
			case serde.DirectFundRequestMethod:
				mt, message, err = c.createLedgerChannelConnection.ReadMessage()
			default:
				panic("unimplemented")
			}
			if err != nil {
				log.Println(err)
			}
			if mt == websocket.TextMessage {
				handler(message)
			}
		}
	}
	go listen()
	return nil
}

func (c *webSocketConnection) Close() {
	c.createLedgerChannelConnection.Close() // TODO there is probably a more graceful protocol https://github.com/fasthttp/websocket/blob/master/_examples/echo/client.go
}
