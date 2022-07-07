package chainservice

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

// MockChain provides an interface which simulates a blockchain network. It is designed for use as a central service which multiple
// ChainServices connect to with Go chans.
//
// It keeps a record of of holdings and adjudication status for each channel, accepts transactions and emits events.
type MockChain struct {
	ChainServiceBase

	holdings   map[types.Destination]types.Funds // holdings tracks funds for each channel
	blockNum   *uint64                           // MockChain is often passed around by value. The pointer allows for shared state.
	txListener chan protocols.ChainTransaction   // this is used to broadcast transactions that have been received
}

// NewMockChain returns a new MockChain.
func NewMockChain() *MockChain {
	mc := MockChain{ChainServiceBase: newChainServiceBase()}
	mc.holdings = make(map[types.Destination]types.Funds)
	mc.blockNum = new(uint64)
	*mc.blockNum = 1

	return &mc
}

// NewMockChainWithTransactionListener returns a new mock chain that will send transactions to the supplied chan.
// This lets us easily rebroadcast transactions to other mock chains.
func NewMockChainWithTransactionListener(txListener chan protocols.ChainTransaction) *MockChain {
	mc := NewMockChain()
	mc.txListener = txListener
	return mc
}

// SendTransaction responds to the given tx.
func (mc *MockChain) SendTransaction(tx protocols.ChainTransaction) {
	*mc.blockNum++
	if mc.txListener != nil {
		mc.txListener <- tx
	}
	var event Event
	switch tx := tx.(type) {
	case protocols.DepositTransaction:
		if tx.Deposit.IsNonZero() {
			mc.holdings[tx.ChannelId()] = mc.holdings[tx.ChannelId()].Add(tx.Deposit)
		}
		event = DepositedEvent{
			CommonEvent: CommonEvent{
				channelID: tx.ChannelId(),
				BlockNum:  *mc.blockNum},

			Holdings: mc.holdings[tx.ChannelId()],
		}
	case protocols.WithdrawAllTransaction:
		for assetAddress := range mc.holdings[tx.ChannelId()] {
			event = AllocationUpdatedEvent{
				CommonEvent: CommonEvent{
					channelID: tx.ChannelId(),
					BlockNum:  *mc.blockNum},
				AssetAddress: assetAddress,
				AssetAmount:  common.Big0,
			}
		}
		mc.holdings[tx.ChannelId()] = types.Funds{}
	default:
		panic("unexpected chain transaction")
	}

	mc.broadcast(event)
}
