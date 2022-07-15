package chainservice

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

// MockChainService adheres to the ChainService interface. The constructor accepts a MockChain, which allows multiple clients to share the same, in-memory chain.
type MockChainService struct {
	chain      *MockChainImpl
	txListener chan protocols.ChainTransaction // this is used to broadcast transactions that have been received
	eventFeed  <-chan Event
}

// NewMockChainService returns a new MockChainService.
func NewMockChainService(chain *MockChainImpl, address common.Address) *MockChainService {
	mc := MockChainService{chain: chain}
	mc.eventFeed = chain.SubscribeToEvents(address)
	return &mc
}

// NewMockChainWithTransactionListener returns a new MockChainService that will send transactions to the supplied chan.
// This lets us easily rebroadcast transactions to other MockChainServices.
func NewMockChainWithTransactionListener(chain *MockChainImpl, address common.Address, txListener chan protocols.ChainTransaction) *MockChainService {
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

func (mc *MockChainService) EventFeed() <-chan Event {
	return mc.eventFeed
}
