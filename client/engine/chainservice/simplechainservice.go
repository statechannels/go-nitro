package chainservice

import (
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

// SimpleChainService forwards inputted transactions to a MockChain, and passes Events straight back.
type SimpleChainService struct {
	out chan Event                      // out is the chan used to send Events to the engine
	in  chan protocols.ChainTransaction // in is the chan used to receive Transactions from the engine

	address types.Address // address is used to subscribe to the MockChain's Out chan
	chain   *MockChain
}

// NewSimpleChainService returns a SimpleChainService which is listening for transactions and events.
func NewSimpleChainService(mc *MockChain, address types.Address) ChainService {
	mcs := SimpleChainService{}
	mcs.out = make(chan Event)
	mcs.in = make(chan protocols.ChainTransaction)
	mcs.chain = mc
	mcs.address = address

	go mcs.forwardEvents()
	go mcs.forwardTransactions()

	return mcs
}

// Out returns the out chan but narrows the type so that external consumers may only receive on it.
func (mcs SimpleChainService) Out() <-chan Event {
	return mcs.out
}

// In returns the in chan but narrows the type so that external consumers may only send on it.
func (mcs SimpleChainService) In() chan<- protocols.ChainTransaction {
	return mcs.in
}

// forwardTransactions pipes transactions to the MockChain
func (mcs SimpleChainService) forwardTransactions() {
	for tx := range mcs.in {
		mcs.chain.In() <- tx
	}
}

// forwardEvents pipes events from the MockChain
func (mcs SimpleChainService) forwardEvents() {
	for event := range mcs.chain.Out(mcs.address) {
		mcs.out <- event
	}
}
