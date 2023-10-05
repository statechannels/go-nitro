package messageservice

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	p2pms "github.com/statechannels/go-nitro/node/engine/messageservice/p2p-message-service"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/rand"
	"github.com/statechannels/go-nitro/types"
)

// TestMessageService is an implementation of the MessageService interface
// for use in multi-engine test environments.
//
// It allows for individual nitro nodes / engines to:
//  1. be instantiated together via test setup data
//  2. "connect" with one another via gochans
//  3. run independently in information-silo goroutines, while
//     communicating on the simulated network
type TestMessageService struct {
	address types.Address

	// connection to Engine:
	out          chan protocols.Message      // for sending message to engine
	signRequests chan p2pms.SignatureRequest // for sending signature requests to engine
	maxDelay     time.Duration               // the max delay for messages

	broker Broker
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
		address:      address,
		out:          make(chan protocols.Message, 5),
		signRequests: make(chan p2pms.SignatureRequest, 5),
		maxDelay:     maxDelay,
		broker:       broker,
	}

	tms.connect(broker)
	return tms
}

func (t TestMessageService) P2PMessages() <-chan protocols.Message {
	return t.out
}

func (t TestMessageService) SignRequests() <-chan p2pms.SignatureRequest {
	return t.signRequests
}

// dispatchMessage is responsible for dispatching a message to the appropriate peer message service.
// If there is a mean delay it will wait a random amount of time(based on meanDelay) before sending the message.
func (t TestMessageService) dispatchMessage(message protocols.Message) {
	if t.maxDelay > 0 {
		randomDelay := time.Duration(rand.Int63n(t.maxDelay.Nanoseconds()))
		time.Sleep(randomDelay)
	}

	peer, ok := t.broker.services[message.To]
	if ok {
		// To mimic a proper message service, we serialize and then
		// deserialize the message

		serializedMsg, err := message.Serialize()
		if err != nil {
			panic(`could not serialize message`)
		}
		peer.HandleMessage([]byte(serializedMsg))
	} else {
		panic(fmt.Sprintf("node %v has no connection to node %v",
			t.address, message.To))
	}
}

// connect registers the message service with the broker
func (tms TestMessageService) connect(b Broker) {
	b.services[tms.address] = tms
}

// Send dispatches messages
func (tms TestMessageService) Send(msg protocols.Message) error {
	tms.dispatchMessage(msg)
	return nil
}

// HandleMessage deserialize the message and feed it to the engine
func (tms TestMessageService) HandleMessage(message []byte) {
	msg, err := protocols.DeserializeMessage(string(message))
	if err != nil {
		panic(fmt.Errorf("could not deserialize message :%w", err))
	}
	tms.out <- msg
}

// Close stops the TestMessagerService from sending or receiving messages.
func (tms TestMessageService) Close() error {
	return nil
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
//                                v HandleMessage
// ┌──────────┐toMsg       in┌───────────┐
// │          │  ───────────►|           │
// │  Engine  │              │  Message  │
// │          │fromMsg    out│  Service  │
// │    B     │  ◄───────────┤    B      │
// └──────────┘              └───────────┘
