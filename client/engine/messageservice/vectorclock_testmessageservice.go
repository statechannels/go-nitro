package messageservice

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/DistributedClocks/GoVector/govec"
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
type VectorClockTestMessageService struct {
	TestMessageService
	goveclogger *govec.GoLog // vector clock logger

}

// NewTestMessageService returns a running TestMessageService
// It accepts an address, a broker, and a max delay for messages.
// Messages will be handled with a random delay between 0 and maxDelay
func NewVectorClockTestMessageService(address types.Address, broker Broker, maxDelay time.Duration, logDir string, prettyName string) VectorClockTestMessageService {

	vctms := VectorClockTestMessageService{
		TestMessageService: TestMessageService{
			address:   address,
			in:        make(chan protocols.Message, 5),
			out:       make(chan protocols.Message, 5),
			maxDelay:  maxDelay,
			fromPeers: make(chan []byte, 5),
		},
		goveclogger: govec.InitGoVector(prettyName, logDir+"/"+address.String(), govec.GetDefaultConfig()),
	}

	vctms.connect(broker)
	go vctms.routeFromPeers()
	go vctms.routeToPeers(broker)

	return vctms
}

// dispatchMessage is responsible for dispatching a message to the appropriate peer message service.
// If there is a mean delay it will wait a random amount of time(based on meanDelay) before sending the message.
func (t VectorClockTestMessageService) dispatchMessage(message protocols.Message, b Broker) {
	if t.maxDelay > 0 {
		randomDelay := time.Duration(rand.Int63n(t.maxDelay.Nanoseconds()))
		time.Sleep(randomDelay)
	}

	peer, ok := b.services[message.To]
	if ok {
		// To mimic a proper message service, we serialize and then
		// deserialize the message

		serializedMsg, err := message.Serialize()
		if err != nil {
			panic(`could not serialize message`)
		}
		vectorClockMessage := t.goveclogger.PrepareSend(summarizeMessageSend(message), serializedMsg, govec.GetDefaultLogOptions())
		peer.fromPeers <- vectorClockMessage
	} else {
		panic(fmt.Sprintf("client %v has no connection to client %v",
			t.address, message.To))
	}
}

func summarizeMessageSend(msg protocols.Message) string {
	return "Send: " + string(msg.Payloads[0].ObjectiveId)
}

// routeToPeers listens for messages from the engine, and dispatches them
func (vctms VectorClockTestMessageService) routeToPeers(b Broker) {
	for message := range vctms.in {
		go vctms.dispatchMessage(message, b)
	}
}

// routeFromPeers listens for messages from peers, deserializes them and feeds them to the engine
func (vctms VectorClockTestMessageService) routeFromPeers() {
	for vectorClockMessage := range vctms.fromPeers {
		message := []byte("")
		vctms.goveclogger.UnpackReceive("Receiving Message", vectorClockMessage, &message, govec.GetDefaultLogOptions())
		msg, err := protocols.DeserializeMessage(string(message))
		if err != nil {
			panic(err)
		}
		vctms.out <- msg
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
