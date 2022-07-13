package chainservice

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/client/engine/store/safesync"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

type MockChain interface {
	SubmitTransaction(protocols.ChainTransaction) error
	SubscribeToEvents(a types.Address) <-chan Event
}

type MockChainImpl struct {
	blockNum uint64
	holdings map[types.Destination]types.Funds // holdings tracks funds for each channel
	out      safesync.Map[chan Event]
}

func NewMockChainImpl() *MockChainImpl {
	chain := MockChainImpl{}
	chain.blockNum = 1
	chain.holdings = make(map[types.Destination]types.Funds)
	chain.out = safesync.Map[chan Event]{}
	return &chain
}

func (mc *MockChainImpl) SubmitTransaction(tx protocols.ChainTransaction) error {
	mc.blockNum++
	switch tx := tx.(type) {
	case protocols.DepositTransaction:
		if tx.Deposit.IsNonZero() {
			mc.holdings[tx.ChannelId()] = mc.holdings[tx.ChannelId()].Add(tx.Deposit)
		}
		for address, amount := range tx.Deposit {
			event := NewDepositedEvent(tx.ChannelId(), mc.blockNum, address, amount, mc.holdings[tx.ChannelId()][address])
			mc.broadcastEvent(event)
		}
	case protocols.WithdrawAllTransaction:
		for assetAddress := range mc.holdings[tx.ChannelId()] {
			event := NewAllocationUpdatedEvent(tx.ChannelId(), mc.blockNum, assetAddress, common.Big0)
			mc.broadcastEvent(event)
		}
		mc.holdings[tx.ChannelId()] = types.Funds{}
	default:
		return fmt.Errorf("unexpected transaction type %T", tx)
	}
	return nil
}

func (mc *MockChainImpl) broadcastEvent(event Event) {
	mc.out.Range(func(_ string, channel chan Event) bool {
		channel <- event
		return true
	})
}

func (mc *MockChainImpl) SubscribeToEvents(a types.Address) <-chan Event {
	// Use a buffered channel so we don't have to worry about blocking on writing to the channel.
	c := make(chan Event, 10)
	mc.out.Store(a.String(), c)
	return c
}

// MockChainService provides an interface which simulates a blockchain network. It is designed for use as a central service which multiple
// ChainServices connect to with Go chans.
//
// It keeps a record of of holdings and adjudication status for each channel, accepts transactions and emits events.
type MockChainService struct {
	chainServiceBase

	chain      MockChain
	txListener chan protocols.ChainTransaction // this is used to broadcast transactions that have been received
}

// NewMockChainService returns a new MockChain.
func NewMockChainService(chain MockChain, address common.Address) *MockChainService {
	mc := MockChainService{chainServiceBase: newChainServiceBase()}
	mc.chain = chain
	in := chain.SubscribeToEvents(address)
	go func() {
		for e := range in {
			mc.out <- e
		}
	}()
	return &mc
}

// NewMockChainWithTransactionListener returns a new mock chain that will send transactions to the supplied chan.
// This lets us easily rebroadcast transactions to other mock chains.
func NewMockChainWithTransactionListener(chain MockChain, address common.Address, txListener chan protocols.ChainTransaction) *MockChainService {
	mc := NewMockChainService(chain, address)
	mc.txListener = txListener
	return mc
}

// SendTransaction responds to the given tx.
func (mc *MockChainService) SendTransaction(tx protocols.ChainTransaction) error {
	if mc.txListener != nil {
		mc.txListener <- tx
	}

	return mc.chain.SubmitTransaction(tx)
}

// GetConsensusAppAddress returns the zero address, since the mock chain will not run any application logic.
func (mc *MockChainService) GetConsensusAppAddress() types.Address {
	return types.Address{}
}
