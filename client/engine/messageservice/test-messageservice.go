package messageservice

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/statechannels/go-nitro/client/engine/store/safesync"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

// TestMessageService is an implementation of the MessageService interface
// for use in multi-engine test environments.
//
// It allows for individual nitro-clients / engines to:
//  1. be instantiated together via test setup data
//  2. "connect" with one another via gochans
//  3. run independently in information-silo goroutines, while
//     communicating on the simulated network
type TestMessageService struct {
	address types.Address

	// connection to Engine:
	out      chan protocols.Message // for sending message to engine
	maxDelay time.Duration          // the max delay for messages

	// connection with Peers:
	fromPeers chan []byte // for receiving serialized messages from peers

	broker *Broker
}

// A Broker manages a mapping from identifying address to a TestMessageService,
// allowing messages sent from one message service to be directed to the intended
// recipient
type Broker struct {
	services safesync.Map[*TestMessageService]
}

func NewBroker() *Broker {
	b := Broker{
		services: safesync.Map[*TestMessageService]{},
	}

	return &b
}

// NewTestMessageService returns a running TestMessageService
// It accepts an address, a broker, and a max delay for messages.
// Messages will be handled with a random delay between 0 and maxDelay
func NewTestMessageService(address types.Address, broker *Broker, maxDelay time.Duration) TestMessageService {
	tms := TestMessageService{
		address:   address,
		out:       make(chan protocols.Message, 5),
		maxDelay:  maxDelay,
		fromPeers: make(chan []byte, 5),
		broker:    broker,
	}

	tms.connect(broker)
	go tms.routeFromPeers()
	return tms
}

func (t *TestMessageService) Out() <-chan protocols.Message {
	return t.out
}

// dispatchMessage is responsible for dispatching a message to the appropriate peer message service.
// If there is a mean delay it will wait a random amount of time(based on meanDelay) before sending the message.
func (t *TestMessageService) dispatchMessage(message protocols.Message) {
	if t.maxDelay > 0 {
		randomDelay := time.Duration(rand.Int63n(t.maxDelay.Nanoseconds()))
		time.Sleep(randomDelay)
	}

	peer, ok := t.broker.services.Load(message.To.Hex())
	if ok {
		// To mimic a proper message service, we serialize and then
		// deserialize the message

		serializedMsg, err := message.Serialize()
		if err != nil {
			panic(`could not serialize message`)
		}
		peer.fromPeers <- []byte(serializedMsg)
	} else {
		panic(fmt.Sprintf("client %v has no connection to client %v",
			t.address, message.To))
	}
}

// connect registers the message service with the broker
func (tms *TestMessageService) connect(b *Broker) {
	b.services.Store(tms.address.Hex(), tms)
}

// Send dispatches messages
func (tms *TestMessageService) Send(msg protocols.Message) {
	tms.dispatchMessage(msg)
}

// routeFromPeers listens for messages from peers, deserializes them and feeds them to the engine
func (tms *TestMessageService) routeFromPeers() {
	for message := range tms.fromPeers {
		msg, err := protocols.DeserializeMessage(string(message))
		if err != nil {
			panic(fmt.Errorf("could not deserialize message :%w", err))
		}
		tms.out <- msg
	}
}

// ┌──────────┐toMsg       in┌───────────┐
// │          │  ───────────►|           │
// │  Engine  │              │  Message  │
// │          │fromMsg    out│  Service  │
// │    A     │  ◄───────────┤    A      │
// └──────────┘              └───────────┘
//                                │
//                                │
//                                │
//                                │
//                                │
//                                v fromPeers
// ┌──────────┐toMsg       in┌───────────┐
// │          │  ───────────►|           │
// │  Engine  │              │  Message  │
// │          │fromMsg    out│  Service  │
// │    B     │  ◄───────────┤    B      │
// └──────────┘              └───────────┘
