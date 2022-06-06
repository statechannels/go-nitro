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
	ChainServiceBase

	holdings map[types.Destination]types.Funds // holdings tracks funds for each channel
	blockNum *uint64                           // MockChain is often passed around by value. The pointer allows for shared state.
}

// NewMockChain returns a new MockChain.
func NewMockChain() *MockChain {
	mc := MockChain{ChainServiceBase: newChainServiceBase()}
	mc.holdings = make(map[types.Destination]types.Funds)
	mc.blockNum = new(uint64)
	*mc.blockNum = 1

	return &mc
}

// SendTransaction responds to the given tx.
func (mc *MockChain) SendTransaction(tx protocols.ChainTransaction) {
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

	mc.broadcast(event)
}
