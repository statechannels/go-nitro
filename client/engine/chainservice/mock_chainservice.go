package chainservice

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

// MockChainService provides an interface which simulates a blockchain network. It is designed for use as a central service which multiple
// ChainServices connect to with Go chans.
//
// It keeps a record of of holdings and adjudication status for each channel, accepts transactions and emits events.
type MockChainService struct {
	chainServiceBase

	holdings   map[types.Destination]types.Funds // holdings tracks funds for each channel
	blockNum   *uint64                           // MockChain is often passed around by value. The pointer allows for shared state.
	txListener chan protocols.ChainTransaction   // this is used to broadcast transactions that have been received
}

// NewMockChainService returns a new MockChain.
func NewMockChainService() *MockChainService {
	mc := MockChainService{chainServiceBase: newChainServiceBase()}
	mc.holdings = make(map[types.Destination]types.Funds)
	mc.blockNum = new(uint64)
	*mc.blockNum = 1

	return &mc
}

// NewMockChainWithTransactionListener returns a new mock chain that will send transactions to the supplied chan.
// This lets us easily rebroadcast transactions to other mock chains.
func NewMockChainWithTransactionListener(txListener chan protocols.ChainTransaction) *MockChainService {
	mc := NewMockChainService()
	mc.txListener = txListener
	return mc
}

// SendTransaction responds to the given tx.
func (mc *MockChainService) SendTransaction(tx protocols.ChainTransaction) error {
	*mc.blockNum++
	if mc.txListener != nil {
		mc.txListener <- tx
	}
	switch tx := tx.(type) {
	case protocols.DepositTransaction:
		if tx.Deposit.IsNonZero() {
			mc.holdings[tx.ChannelId()] = mc.holdings[tx.ChannelId()].Add(tx.Deposit)
		}
		for address, amount := range tx.Deposit {
			event := NewDepositedEvent(tx.ChannelId(), *mc.blockNum, address, amount, mc.holdings[tx.ChannelId()][address])
			mc.broadcast(event)
		}
	case protocols.WithdrawAllTransaction:
		for assetAddress := range mc.holdings[tx.ChannelId()] {
			event := NewAllocationUpdatedEvent(tx.ChannelId(), *mc.blockNum, assetAddress, common.Big0)
			mc.broadcast(event)
		}
		mc.holdings[tx.ChannelId()] = types.Funds{}
	default:
		return fmt.Errorf("unexpected transaction type %T", tx)
	}
	return nil
}

// GetConsensusAppAddress returns the zero address, since the mock chain will not run any application logic.
func (mc *MockChain) GetConsensusAppAddress() types.Address {
	return types.Address{}
}
