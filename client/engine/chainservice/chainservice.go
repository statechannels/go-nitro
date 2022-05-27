// Package chainservice is a chain service responsible for submitting blockchain transactions and relaying blockchain events.
package chainservice // import "github.com/statechannels/go-nitro/client/chainservice"

import (
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
