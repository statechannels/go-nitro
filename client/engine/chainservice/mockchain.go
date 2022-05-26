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
	in  chan protocols.ChainTransaction // in is the chan used to receive Transactions from multiple ChainServices

	transListener chan protocols.ChainTransaction   // this is used to broadcast transactions that have been received
	holdings      map[types.Destination]types.Funds // holdings tracks funds for each channel
	blockNum      uint64
}

// Out returns the out chan for a particular ChainService, and narrows the type so that external consumers may only receive on it.
func (mc MockChain) Out(a types.Address) <-chan Event {
	return mc.out[a]
}

// In returns the in chan but narrows the type so that external consumers may only send on it.
func (mc MockChain) In() chan<- protocols.ChainTransaction {
	return mc.in
}

// NewMockChain returns a new MockChain.
func NewMockChain() MockChain {
	return NewMockChainWithTransactionListener(nil)
}

// NewMockChainWithTransactionListener returns a new MockChain with the supplied transaction listener.
// The transaction listener will receive all transactions that are sent to the MockChain.
func NewMockChainWithTransactionListener(transactionListener chan protocols.ChainTransaction) MockChain {

	mc := MockChain{}
	mc.out = make(map[types.Address]chan Event)
	mc.in = make(chan protocols.ChainTransaction, 1000)
	mc.holdings = make(map[types.Destination]types.Funds)
	mc.transListener = transactionListener
	mc.blockNum = 1

	go mc.Run()
	return mc
}

// Subscribe inserts a go chan (for the supplied address) into the MockChain.
func (mc *MockChain) Subscribe(a types.Address) {
	// Use a buffered channel so we don't have to worry about blocking on writing to the channel.
	mc.out[a] = make(chan Event, 1000)
}

// Run starts a listener for transactions on the MockChain's in chan.
func (mc MockChain) Run() {
	for tx := range mc.in {
		mc.sendToTransListener(tx)
		mc.blockNum++
		mc.handleTx(tx)
	}
}

// sendToTransListener sends the transaction to the transListener if not nil and the chan is not full.
func (mc *MockChain) sendToTransListener(tx protocols.ChainTransaction) {
	if mc.transListener != nil {
		// Send to transListener and ignore if the chan is full
		select {
		case mc.transListener <- tx:
		default:
		}
	}
}

// handleTx responds to the given tx.
func (mc MockChain) handleTx(tx protocols.ChainTransaction) {
	if tx.Deposit.IsNonZero() {
		mc.holdings[tx.ChannelId] = mc.holdings[tx.ChannelId].Add(tx.Deposit)
	}
	var event Event
	switch tx.Type {
	case protocols.DepositTransactionType:
		event = DepositedEvent{
			CommonEvent: CommonEvent{
				channelID:          tx.ChannelId,
				AdjudicationStatus: protocols.AdjudicationStatus{TurnNumRecord: 0},
				BlockNum:           mc.blockNum},

			Holdings: mc.holdings[tx.ChannelId],
		}
	case protocols.WithdrawAllTransactionType:
		mc.holdings[tx.ChannelId] = types.Funds{}
		event = AllocationUpdatedEvent{
			CommonEvent: CommonEvent{
				channelID:          tx.ChannelId,
				AdjudicationStatus: protocols.AdjudicationStatus{TurnNumRecord: 0},
				BlockNum:           mc.blockNum},

			Holdings: mc.holdings[tx.ChannelId],
		}
	default:
		panic("unexpected chain transaction")
	}

	for _, out := range mc.out {
		attemptSend(out, event)
	}

}

// attemptSend sends event to the supplied chan, and drops it if the chan is full
func attemptSend(out chan Event, event Event) {
	select {
	case out <- event:
	default:
	}
}
