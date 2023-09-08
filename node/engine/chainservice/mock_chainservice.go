package chainservice

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

// MockChainService adheres to the ChainService interface. The constructor accepts a MockChain, which allows multiple clients to share the same, in-memory chain.
type MockChainService struct {
	chain     *MockChain
	eventFeed <-chan Event
}

// NewMockChainService returns a new MockChainService.
func NewMockChainService(chain *MockChain, address common.Address) *MockChainService {
	mc := MockChainService{chain: chain}
	mc.eventFeed = chain.SubscribeToEvents(address)
	return &mc
}

// SendTransaction responds to the given tx.
func (mc *MockChainService) SendTransaction(tx protocols.ChainTransaction) error {
	return mc.chain.SubmitTransaction(tx)
}

// GetConsensusAppAddress returns the zero address, since the mock chain will not run any application logic.
func (mc *MockChainService) GetConsensusAppAddress() types.Address {
	return types.Address{}
}

// GetVirtualPaymentAppAddress returns the zero address, since the mock chain will not run any application logic.
func (mc *MockChainService) GetVirtualPaymentAppAddress() types.Address {
	return types.Address{}
}

func (mc *MockChainService) EventFeed() <-chan Event {
	return mc.eventFeed
}

func (mc *MockChainService) GetChainId() (*big.Int, error) {
	return big.NewInt(TEST_CHAIN_ID), nil
}

func (mc *MockChainService) GetLastConfirmedBlockNum() uint64 {
	mc.chain.BlockNumMu.Lock()
	blockNum := mc.chain.BlockNum
	mc.chain.BlockNumMu.Unlock()

	return blockNum
}

func (mc *MockChainService) Close() error {
	return nil
}
