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
	BlockNum   uint64
	blockNumMu sync.Mutex
	// holdings tracks funds for each channel.
	holdings map[types.Destination]types.Funds
	// out maps addresses to an Event channel. Given that MockChainServices only subscribe
	// (and never unsubscribe) to events, this can be converted to a list.
	out safesync.Map[chan Event]
}

// NewMockChain creates a new MockChain
func NewMockChain() *MockChain {
	chain := MockChain{}
	chain.BlockNum = 1
	chain.holdings = map[types.Destination]types.Funds{}
	chain.out = safesync.Map[chan Event]{}
	return &chain
}

// SubmitTransaction updates internal state and broadcasts events
// unlike an ethereum blockchain, MockChain accepts go-nitro protocols.ChainTransaction
func (mc *MockChain) SubmitTransaction(tx protocols.ChainTransaction) error {
	eventsToBroadcast := []Event{}
	mc.blockNumMu.Lock()
	mc.BlockNum++
	h := mc.holdings[tx.ChannelId()] // ignore `ok` because the returned zero-value is what we want
	switch tx := tx.(type) {
	case protocols.DepositTransaction:
		if tx.Deposit.IsNonZero() {
			mc.holdings[tx.ChannelId()] = h.Add(tx.Deposit)
		}

		for address := range tx.Deposit {
			event := NewDepositedEvent(tx.ChannelId(), mc.BlockNum, address, h.Add(tx.Deposit)[address])
			eventsToBroadcast = append(eventsToBroadcast, event)
		}
	case protocols.WithdrawAllTransaction:
		for assetAddress := range h {
			event := NewAllocationUpdatedEvent(tx.ChannelId(), mc.BlockNum, assetAddress, common.Big0)
			eventsToBroadcast = append(eventsToBroadcast, event)
		}
		mc.holdings[tx.ChannelId()] = types.Funds{}
	default:
		return fmt.Errorf("unexpected transaction type %T", tx)
	}
	mc.blockNumMu.Unlock()
	for _, event := range eventsToBroadcast {
		mc.broadcastEvent(event)
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
