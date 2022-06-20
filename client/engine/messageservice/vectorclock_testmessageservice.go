package messageservice

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/DistributedClocks/GoVector/govec"
	"github.com/statechannels/go-nitro/channel/consensus_channel"
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
) VectorClockTestMessageService {

	vctms := VectorClockTestMessageService{
		TestMessageService: TestMessageService{
			address:   address,
			out:       make(chan protocols.Message, 5),
			maxDelay:  maxDelay,
			fromPeers: make(chan []byte, 5),
			broker:    broker,
		},
		goveclogger: govec.InitGoVector(address.String(), logDir+"/"+address.String(), govec.GetDefaultConfig()),
	}

	vctms.connect(broker)
	go vctms.routeFromPeers()

	return vctms
}

// dispatchMessage is responsible for dispatching a message to the appropriate peer message service.
// If there is a mean delay it will wait a random amount of time(based on meanDelay) before sending the message.
// Messages sends are intercepted by the vector clock logger, which runs the vector clock algorithm and wraps the message accordingly.
func (t VectorClockTestMessageService) dispatchMessage(message protocols.Message) {
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
	str, _ := msg.Serialize()
	size := len([]byte(str))
	summary := fmt.Sprint(size) + "B:"
	for _, entry := range msg.SignedProposals() {
		summary += `propose `
		summary += fmt.Sprint(entry.Payload.Proposal.LedgerID)
		if (entry.Payload.Proposal.ToAdd != consensus_channel.Add{}) {
			summary += ` funds `
			summary += fmt.Sprint(entry.Payload.Proposal.ToAdd.Target())
		}
		if (entry.Payload.Proposal.ToRemove != consensus_channel.Remove{}) {
			summary += ` defunds `
			summary += fmt.Sprint(entry.Payload.Proposal.ToRemove.Target)
		}
	}
	for _, entry := range msg.SignedStates() {
		summary += `send `
		_, turnNum := entry.Payload.SortInfo()
		summary += fmt.Sprint(entry.Payload.ChannelId())
		summary += ` @turn `
		summary += fmt.Sprint(turnNum)

	}
	return summary
}

// Send dispatches messages
func (vctms VectorClockTestMessageService) Send(msg protocols.Message) {
	vctms.dispatchMessage(msg)
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
