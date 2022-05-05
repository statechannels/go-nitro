package messageservice

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/internal/protocols"
	"github.com/statechannels/go-nitro/internal/types"
)

// TestMessageService is an implementaion of the MessageService interface
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
	in       chan protocols.Message // for recieving messages from engine
	out      chan protocols.Message // for sending message to engine
	maxDelay time.Duration          // the max delay for messages
}

// A Broker manages a mapping from identifying address to a TestMessageService,
// allowing messages sent from one message service to be directed to the intended
// recipient
type Broker struct {
	services map[types.Address]TestMessageService
}

func NewBroker() Broker {
	b := Broker{
		services: make(map[common.Address]TestMessageService),
	}

	return b
}

// NewTestMessageService returns a running TestMessageService
// It accepts an address, a broker, and a max delay for messages.
// Messages will be handled with a random delay between 0 and maxDelay
func NewTestMessageService(address types.Address, broker Broker, maxDelay time.Duration) TestMessageService {
	tms := TestMessageService{
		address:  address,
		in:       make(chan protocols.Message, 5),
		out:      make(chan protocols.Message, 5),
		maxDelay: maxDelay,
	}

	tms.connect(broker)
	return tms
}

func (t TestMessageService) Out() <-chan protocols.Message {
	return t.out
}

func (t TestMessageService) In() chan<- protocols.Message {
	return t.in
}

// dispatchMessage is responsible for dispatching a message to the appropriate peer message service.
// If there is a mean delay it will wait a random amount of time(based on meanDelay) before sending the message.
// It serializes and deserializes a message to mimic a real message service.
func (t TestMessageService) dispatchMessage(message protocols.Message, b Broker) {
	if t.maxDelay > 0 {
		randomDelay := time.Duration(rand.Int63n(t.maxDelay.Nanoseconds()))
		time.Sleep(randomDelay)
	}

	peerChan, ok := b.services[message.To]
	if ok {
		// To mimic a proper message service, we serialize and then
		// deserialize the message

		serializedMsg, err := message.Serialize()
		if err != nil {
			panic(`could not serialize message`)
		}
		deserializedMsg, err := protocols.DeserializeMessage(serializedMsg)
		if err != nil {
			panic(`could not deserialize message`)
		}
		peerChan.out <- deserializedMsg
	} else {
		panic(fmt.Sprintf("client %v has no connection to client %v",
			t.address, message.To))
	}
}

// connect creates a gochan for message service to send messages to the given peer.
func (t TestMessageService) connect(b Broker) {
	go func() {
		for message := range t.in {

			go t.dispatchMessage(message, b)
		}

	}()

	b.services[t.address] = t
}

// ┌──────────┐toMsg       in┌───────────┐
// │          │  ───────────►|           │
// │  Engine  │              │  Message  │
// │          │fromMsg    out│  Service  │
// │    A     │  ◄───────────┤    A      │
// └──────────┘              └────┬──────┘
//                                │toPeers[B]
//                                │
//                                │
//                     ┌──────────┘
//                     │
//                     │
//                     │
// ┌──────────┐toMsg   │   in┌───────────┐
// │          │  ──────┼────►|           │
// │  Engine  │        │     │  Message  │
// │          │fromMsg │  out│  Service  │
// │    B     │  ◄─────┴─────┤    B      │
// └──────────┘              └───────────┘
