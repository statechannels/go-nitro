package chainservice

import (
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

// simpleChainService forwards inputted transactions to a MockChain, and passes Events straight back.
type simpleChainService struct {
	out chan Event                      // out is the chan used to send Events to the engine
	in  chan protocols.ChainTransaction // in is the chan used to receive Transactions from the engine
}

// NewSimpleChainService returns a SimpleChainService which is listening for transactions and events.
func NewSimpleChainService(mc *MockChain, address types.Address) ChainService {
	mcs := simpleChainService{}
	mcs.out = make(chan Event)
	mcs.in = mc.in
	mc.Subscribe(address)

	go mcs.forwardEvents(mc, address)

	return mcs
}

// Out returns the out chan but narrows the type so that external consumers may only receive on it.
func (mcs simpleChainService) Out() <-chan Event {
	return mcs.out
}

// In returns the in chan but narrows the type so that external consumers may only send on it.
func (mcs simpleChainService) In() chan<- protocols.ChainTransaction {
	return mcs.in
}

// forwardEvents pipes events from the MockChain
func (mcs simpleChainService) forwardEvents(mc *MockChain, address types.Address) {
	for event := range mc.Out(address) {
		mcs.out <- event
	}
}
