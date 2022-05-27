package chainservice

import (
	"fmt"

	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

// MockChain provides an interface which simulates a blockchain network. It is designed for use as a central service which multiple
// ChainServices connect to with Go chans.
//
// It keeps a record of of holdings and adjudication status for each channel, accepts transactions and emits events.
type MockChain struct {
	out map[types.Address]chan Event // out is a mapping with a chan for each connected ChainService, used to send Events to that service

	transListener chan protocols.ChainTransaction   // this is used to broadcast transactions that have been received
	holdings      map[types.Destination]types.Funds // holdings tracks funds for each channel
	blockNum      *uint64
}

// Out returns the out chan for a particular ChainService, and narrows the type so that external consumers may only receive on it.
func (mc MockChain) EventFeed(a types.Address) (<-chan Event, error) {
	feed, ok := mc.out[a]
	if !ok {
		return nil, fmt.Errorf("no subscription for address %v", a)
	}
	return feed, nil
}

// NewMockChain returns a new MockChain.
func NewMockChain() MockChain {
	mc := MockChain{}
	mc.out = make(map[types.Address]chan Event)
	mc.holdings = make(map[types.Destination]types.Funds)
	mc.blockNum = new(uint64)
	*mc.blockNum = 1

	return mc
}

// Subscribe inserts a go chan (for the supplied address) into the MockChain.
func (mc MockChain) SubscribeToEvents(a types.Address) <-chan Event {
	// Use a buffered channel so we don't have to worry about blocking on writing to the channel.
	mc.out[a] = make(chan Event, 10)
	return mc.out[a]
}

// handleTx responds to the given tx.
func (mc MockChain) SendTransaction(tx protocols.ChainTransaction) {
	*mc.blockNum++

	if tx.Deposit.IsNonZero() {
		mc.holdings[tx.ChannelId] = mc.holdings[tx.ChannelId].Add(tx.Deposit)
	}
	var event Event
	switch tx.Type {
	case protocols.DepositTransactionType:
		event = DepositedEvent{
			CommonEvent: CommonEvent{
				channelID: tx.ChannelId,
				BlockNum:  *mc.blockNum},

			Holdings: mc.holdings[tx.ChannelId],
		}
	case protocols.WithdrawAllTransactionType:
		mc.holdings[tx.ChannelId] = types.Funds{}
		event = AllocationUpdatedEvent{
			CommonEvent: CommonEvent{
				channelID: tx.ChannelId,
				BlockNum:  *mc.blockNum},

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
