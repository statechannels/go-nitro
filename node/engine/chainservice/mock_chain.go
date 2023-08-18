package chainservice

import (
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/internal/safesync"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

// MockChain mimics the Ethereum blockchain by keeping track of block numbers and account balances in memory.
// MockChain accepts transactions and broadcasts events.
type MockChain struct {
	blockNum uint64
	// holdings tracks funds for each channel.
	holdings safesync.Map[types.Funds]
	// out maps addresses to an Event channel. Given that MockChainServices only subscribe
	// (and never unsubscribe) to events, this can be converted to a list.
	out safesync.Map[chan Event]

	txLock *sync.Mutex
}

// NewMockChain creates a new MockChain
func NewMockChain() *MockChain {
	chain := MockChain{}
	chain.blockNum = 1
	chain.holdings = safesync.Map[types.Funds]{}
	chain.out = safesync.Map[chan Event]{}
	chain.txLock = &sync.Mutex{}
	return &chain
}

// SubmitTransaction updates internal state and broadcasts events
// unlike an ethereum blockchain, MockChain accepts go-nitro protocols.ChainTransaction
func (mc *MockChain) SubmitTransaction(tx protocols.ChainTransaction) error {
	mc.txLock.Lock()
	defer mc.txLock.Unlock()
	mc.blockNum++
	channelIdString := tx.ChannelId().String()
	h, _ := mc.holdings.Load(channelIdString) // ignore `ok` because the returned zero-value is what we want
	switch tx := tx.(type) {
	case protocols.DepositTransaction:
		if tx.Deposit.IsNonZero() {
			mc.holdings.Store(channelIdString, h.Add(tx.Deposit))
		}

		for address := range tx.Deposit {
			event := NewDepositedEvent(tx.ChannelId(), mc.blockNum, address, h.Add(tx.Deposit)[address])
			mc.broadcastEvent(event)
		}
	case protocols.WithdrawAllTransaction:
		for assetAddress := range h {
			event := NewAllocationUpdatedEvent(tx.ChannelId(), mc.blockNum, assetAddress, common.Big0)
			mc.broadcastEvent(event)
		}
		mc.holdings.Store(channelIdString, types.Funds{})
	default:
		return fmt.Errorf("unexpected transaction type %T", tx)
	}
	return nil
}

func (mc *MockChain) broadcastEvent(event Event) {
	mc.out.Range(func(_ string, channel chan Event) bool {
		channel <- event
		return true
	})
}

// SubscribeToEvents creates, stores, and returns a new Event channel that produces all chain Events
func (mc *MockChain) SubscribeToEvents(a types.Address) <-chan Event {
	// Use a buffered channel so we don't have to worry about blocking on writing to the channel.
	c := make(chan Event, 10)
	mc.out.Store(a.String(), c)
	return c
}

func (mc *MockChain) Close() error {
	f := func(key string, value chan Event) bool {
		close(value)
		return true
	}
	mc.out.Range(f)
	return nil
}
