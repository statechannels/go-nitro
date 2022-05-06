package chainservice

import (
	"math/rand"
	"time"

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
	blockNum uint64
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
func NewMockChain() MockChain {

	mc := MockChain{}
	mc.out = make(map[types.Address]chan Event)
	mc.in = make(chan protocols.ChainTransaction)
	mc.holdings = make(map[types.Destination]types.Funds)
	mc.blockNum = 1

	go mc.Run()
	return mc
}

func (mc *MockChain) Subscribe(a types.Address) {
	// Use a buffered channel so we don't have to worry about blocking on writing to the channel.
	mc.out[a] = make(chan Event, 10)
}

// Run starts a listener for transactions on the MockChain's in chan.
func (mc MockChain) Run() {
	for tx := range mc.in {
		mc.blockNum++
		mc.handleTx(tx)
	}
}

// handleTx responds to the given tx.
func (mc MockChain) handleTx(tx protocols.ChainTransaction) {
	if tx.Deposit.IsNonZero() {
		mc.holdings[tx.ChannelId] = mc.holdings[tx.ChannelId].Add(tx.Deposit)
	}
	maxDelay := time.Millisecond * 100
	randomDelay := time.Duration(rand.Int63n(maxDelay.Nanoseconds()))
	time.Sleep(randomDelay)

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
	maxDelay := time.Millisecond * 100
	randomDelay := time.Duration(rand.Int63n(maxDelay.Nanoseconds()))
	time.Sleep(randomDelay)

	select {
	case out <- event:
	default:
	}
}
