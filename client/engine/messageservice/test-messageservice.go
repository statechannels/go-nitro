/*
TestMessageService is an implementaion of the MessageService interface
for use in
*/

package messageservice

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

type TestNetwork struct {
	addressToPeer map[[]byte]types.Address
}

func (tn TestNetwork) Connect (alice, bob TestMessageService) {
	ch := make(chan protocols.Message)

}

type TestMessageService struct {
	// a map of gochans to pass messages to each speific peer
	toPeers map[types.Address]chan protocols.Message
	out     chan protocols.Message

	// a collection of gochans to recieve messages from specific peers
	fromPeers []chan protocols.Message
	in      chan protocols.Message
}

func (t TestMessageService) GetReceiveChan() <-chan protocols.Message {
	return t.in
}

func (t TestMessageService) GetSendChan() chan<- protocols.Message {
	return t.out
}

func (t TestMessageService) Send(message protocols.Message) {
	t.out <- message;

	
}

func (t TestMessageService) forward(message protocols.Message) {
	peerChan, ok := t.toPeers[message.To]
	if !ok {
		t.toPeers[message.To] = make(chan<- protocols.Message)
		t.Send(message)
	}
	peerChan <- message

	addr := common.HexToAddress(string(message.To))
	t.toPeers[addr] <- message
}

func (t TestMessageService) Run() {
	for {
		select {
		case outgoing := <-t.out:
			t.forward(outgoing)
		case incoming :=
		}

	}
}
