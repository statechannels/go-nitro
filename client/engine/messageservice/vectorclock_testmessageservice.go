package messageservice

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/DistributedClocks/GoVector/govec"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

// VectorClockTestMessageService embeds a TestMessageService and extends it with instrumentation outputing vector clock logs.
// See https://en.wikipedia.org/wiki/Vector_clock and https://github.com/DistributedClocks/GoVector
type VectorClockTestMessageService struct {
	TestMessageService
	goveclogger *govec.GoLog // vector clock logger

}

// NewVectorClockTestMessageService returns a running VectorClockTestMessageService
// It accepts an address, a broker, a max delay for messages and a prettyName for use in the log output.
// Messages will be handled with a random delay between 0 and maxDelay
func NewVectorClockTestMessageService(
	address types.Address,
	broker Broker,
	maxDelay time.Duration,
	logDir string,
	prettyName string,
) VectorClockTestMessageService {

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
// Messages sends are intercepted by the vector clock logger, which runs the vector clock algorithm and wraps the message accordingly.
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

// summarizeMessageSend returns a string which tersely summarizes the supplied message.
// It may be used to make logs more readable.
func summarizeMessageSend(msg protocols.Message) string {
	summary := ""
	for _, entry := range msg.SignedProposals() {
		summary += `propose `
		summary += fmt.Sprint(entry.Payload.Proposal.ChannelID)[1:8]
		summary += ` funds `
		summary += fmt.Sprint(entry.Payload.Proposal.ToAdd.Target())[1:8]
	}
	for _, entry := range msg.SignedStates() {
		summary += `send `
		if len(entry.Payload.State().Participants) == 3 {
			summary += `V`
		} else {
			summary += `L`
		}
		summary += fmt.Sprint(entry.Payload.ChannelId())[1:8]
		summary += fmt.Sprint(entry.Payload.TurnNum())
		summary += ` @turn `
		summary += fmt.Sprint(entry.Payload.TurnNum())

	}
	return summary
}

// routeToPeers listens for messages from the engine, and dispatches them
func (vctms VectorClockTestMessageService) routeToPeers(b Broker) {
	for message := range vctms.in {
		go vctms.dispatchMessage(message, b)
	}
}

// routeFromPeers listens for messages from peers, deserializes them and feeds them to the engine.
// Inbound messages are intercepted by the vector clock logger, which unwraps the message (stripping off the vector clock header) and runs the vector clock algorithm.
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
