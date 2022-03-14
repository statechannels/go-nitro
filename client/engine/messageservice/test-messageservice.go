package messageservice

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
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
	in  chan protocols.Message // for recieving messages from engine
	out chan protocols.Message // for sending message to engine
}

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
func NewTestMessageService(address types.Address, broker Broker) TestMessageService {
	tms := TestMessageService{
		address: address,
		in:      make(chan protocols.Message, 5),
		out:     make(chan protocols.Message, 5),
	}

	tms.Connect(broker)
	return tms
}

func (t TestMessageService) Out() <-chan protocols.Message {
	return t.out
}

func (t TestMessageService) In() chan<- protocols.Message {
	return t.in
}

// Connect creates a gochan for message service to send messages to the given peer.
func (t TestMessageService) Connect(b Broker) {
	go func() {
		for message := range t.in {
			peerChan, ok := b.services[message.To]
			if ok {
				msg, err := message.Serialize()
				if err != nil {
					panic(`could not serialize message`)
				}
				m, err := protocols.DeserializeMessage(msg)
				if err != nil {
					panic(`could not deserialize message`)
				}
				peerChan.out <- m
			} else {
				panic(fmt.Sprintf("client %v has no connection to client %v",
					t.address, message.To))
			}
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
