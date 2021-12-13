/*
TestMessageService is an implementaion of the MessageService interface
for use in
*/

package messageservice

import (
	"fmt"

	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

type TestMessageService struct {
	address types.Address

	// a map of gochans to pass messages to each speific peer
	toPeers map[types.Address]chan<- protocols.Message
	out     chan protocols.Message

	in chan protocols.Message
}

func (t TestMessageService) Run() {
	go t.routeOutgoing()
}

func (t TestMessageService) GetReceiveChan() chan protocols.Message {
	return t.in
}

func (t TestMessageService) GetSendChan() chan<- protocols.Message {
	return t.out
}

func (t TestMessageService) Send(message protocols.Message) {
	t.out <- message
}

// Connect creates a gochan for message service t to communicate with
// the given peer.
func (t TestMessageService) Connect(peer TestMessageService) {
	toPeer := make(chan protocols.Message)

	t.toPeers[peer.address] = toPeer

	go func() {
		for msg := range toPeer {
			peer.in <- msg
		}
	}()
}

func (t TestMessageService) forward(message protocols.Message) {
	peerChan, ok := t.toPeers[message.To]
	if ok {
		peerChan <- message
	} else {
		panic(fmt.Sprintf("client %v has no connection to client %v",
			t.address, message.To))
	}
}

func (t TestMessageService) routeOutgoing() {
	for msg := range t.out {
		t.forward(msg)
	}
}
