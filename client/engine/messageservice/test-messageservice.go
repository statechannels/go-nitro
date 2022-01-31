package messageservice

import (
	"fmt"

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

	// connection to peer message services
	toPeers map[types.Address]chan<- protocols.Message

	// connection to Engine:
	out chan protocols.Message // for recieving messages from engine
	in  chan protocols.Message // for sending message to engine
}

// NewTestMessageService returns a running TestMessageService
func NewTestMessageService(address types.Address) TestMessageService {
	tms := TestMessageService{
		address: address,
		toPeers: make(map[types.Address]chan<- protocols.Message),
		in:      make(chan protocols.Message),
		out:     make(chan protocols.Message),
	}
	tms.run()
	return tms
}

func (t TestMessageService) run() {
	go t.routeOutgoing()
}

func (t TestMessageService) GetReceiveChan() <-chan protocols.Message {
	return t.in
}

func (t TestMessageService) GetSendChan() chan<- protocols.Message {
	return t.out
}

func (t TestMessageService) Send(message protocols.Message) {
	t.out <- message
}

// Connect creates a gochan for message service to send messages to the given peer.
func (t TestMessageService) Connect(peer TestMessageService) {
	toPeer := make(chan protocols.Message)

	t.toPeers[peer.address] = toPeer

	go func() {
		for msg := range toPeer {
			peer.in <- msg // send messages directly to peer's engine, bypassing their message service
		}
	}()
}

// forward finds the appropriate gochan for the message recipient,
// and sends the message. It panics if no such channel exists.
func (t TestMessageService) forward(message protocols.Message) {
	peerChan, ok := t.toPeers[message.To]
	if ok {
		peerChan <- message
	} else {
		panic(fmt.Sprintf("client %v has no connection to client %v",
			t.address, message.To))
	}
}

// routeOutgoing listens to the messageService's outbox and passes
// messages to the forwarding function
func (t TestMessageService) routeOutgoing() {
	for msg := range t.out {
		t.forward(msg)
	}
}

// ┌──────────┐          out┌───────────┐
// │          │  ───────────►           │
// │  Engine  │             │  Message  │
// │          │           in│  Service  │
// │    A     │  ◄──────────┤    A      │
// └──────────┘             └────┬──────┘
// 							     │toPeers[B]
// 							     │
// 							     │
// 				    ┌────────────┘
// 				    │
// 				    │
// 				    │
// ┌──────────┐     │    out┌───────────┐
// │          │  ───┼───────►           │
// │  Engine  │     │       │  Message  │
// │          │     │     in│  Service  │
// │    B     │  ◄──┴───────┤    B      │
// └──────────┘             └───────────┘
