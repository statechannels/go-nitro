package chainservice

import (
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

// MockChain provides an interface which simulates a blockchain network. It is designed for use as a central service which multiple
// ChainServices connect to with Go chans.
//
// It keeps a record of of holdings and adjudication status for each channel, accepts transactions and emits events.
type MockChain struct {
	out map[types.Address]chan Event    // out is a mapping with a chan for each connected ChainService, used to send Events to that service
	in  chan protocols.ChainTransaction // in is the chan used to recieve Transactions from multiple ChainServices

	holdings map[types.Destination]types.Funds // holdings tracks funds for each channel
}

// Out returns the out chan for a particular ChainService, and narrows the type so that external consumers may only receive on it.
func (mc MockChain) Out(a types.Address) <-chan Event {
	return mc.out[a]
}

// In returns the in chan but narrows the type so that external consumers may only send on it.
func (mc MockChain) In() chan<- protocols.ChainTransaction {
	return mc.in
}

// NewMockChain returns a new MockChain with an out chan initialized for each of the addresses passed in.
func NewMockChain(addresses []types.Address) MockChain {

	mc := MockChain{}
	mc.out = make(map[types.Address]chan Event)
	mc.in = make(chan protocols.ChainTransaction)
	mc.holdings = make(map[types.Destination]types.Funds)

	for _, a := range addresses {
		// Use a buffered channel so we don't have to worry about blocking on writing to the channel.
		mc.out[a] = make(chan Event, 10)
	}

	go mc.Run()
	return mc
}

// Run starts a listener for transactions on the MockChain's in chan.
func (mc MockChain) Run() {
	for tx := range mc.in {
		mc.handleTx(tx)
	}
}

// handleTx responds to the given tx.
func (mc MockChain) handleTx(tx protocols.ChainTransaction) {
	if tx.Deposit.IsNonZero() {
		mc.holdings[tx.ChannelId] = mc.holdings[tx.ChannelId].Add(tx.Deposit)
	}
	event := Event{
		ChannelId:          tx.ChannelId,
		Holdings:           mc.holdings[tx.ChannelId],
		AdjudicationStatus: protocols.AdjudicationStatus{TurnNumRecord: 0},
	}
	for _, out := range mc.out {
		sendEvent(out, event)
	}

}

// sendEvent sends event to the supplied chan
func sendEvent(out chan Event, event Event) {
	out <- event
}
