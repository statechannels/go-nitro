// Package chainservice is a chain service responsible for submitting blockchain transactions and relaying blockchain events.
package chainservice // import "github.com/statechannels/go-nitro/client/chainservice"

import (
	"fmt"

	"github.com/statechannels/go-nitro/client/engine/store/safesync"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

// Event dictates which methods all chain events must implement
type Event interface {
	ChannelID() types.Destination
}

// CommonEvent declares fields shared by all chain events
type CommonEvent struct {
	channelID types.Destination
	BlockNum  uint64
}

func (ce CommonEvent) ChannelID() types.Destination {
	return ce.channelID
}

// DepositedEvent is an internal representation of the deposited blockchain event
type DepositedEvent struct {
	CommonEvent
	Holdings types.Funds // indexed by asset
}

// AllocationUpdated is an internal representation of the AllocatonUpdated blockchain event
type AllocationUpdatedEvent struct {
	CommonEvent
	Holdings types.Funds // indexed by asset
}

type ConcludedEvent struct {
	CommonEvent
}

// todo implement other event types
// Concluded
// ChallengeRegistered
// ChallengeCleared

// ChainEventHandler describes an objective that can handle chain events
type ChainEventHandler interface {
	UpdateWithChainEvent(event Event) (protocols.Objective, error)
}

type ChainService interface {
	// EventFeed returns a chan for receiving events from the chain service. An error is returned if no subscription exists
	EventFeed(types.Address) (<-chan Event, error)
	// SubscribeToEvents creates and returs a subscription channel.
	SubscribeToEvents(types.Address) <-chan Event
	// SendTransaction is for sending transactions with the chain service
	SendTransaction(protocols.ChainTransaction)
}

type ChainServiceBase struct {
	out safesync.Map[chan Event]
}

// newChainServiceBase constructs a ChainServiceBase. Only implementations of ChainService interface should call the constructor.
func newChainServiceBase() ChainServiceBase {
	return ChainServiceBase{out: safesync.Map[chan Event]{}}
}

// Subscribe inserts a go chan (for the supplied address) into the ChainService.
func (csb *ChainServiceBase) SubscribeToEvents(a types.Address) <-chan Event {
	// Use a buffered channel so we don't have to worry about blocking on writing to the channel.
	c := make(chan Event, 10)
	csb.out.Store(a.String(), c)
	return c
}

// EventFeed returns the out chan for a particular ChainService, and narrows the type so that external consumers may only receive on it.
func (csb *ChainServiceBase) EventFeed(a types.Address) (<-chan Event, error) {
	c, ok := csb.out.Load(a.String())
	if !ok {
		return nil, fmt.Errorf("no subscription for address %v", a)
	}
	return c, nil
}

func (csb *ChainServiceBase) broadcast(event Event) {
	csb.out.Range(func(_ string, channel chan Event) bool {
		attemptSend(channel, event)
		return true
	})
}

// attemptSend sends event to the supplied chan, and drops it if the chan is full
func attemptSend(out chan Event, event Event) {
	select {
	case out <- event:
	default:
	}
}
