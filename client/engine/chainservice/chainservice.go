// Package chainservice is a chain service responsible for submitting blockchain transactions and relaying blockchain events.
package chainservice // import "github.com/statechannels/go-nitro/client/chainservice"

import (
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

type Event interface {
	GetChannelId() types.Destination
	GetBlockNum() uint64
}

type CommonEvent struct {
	ChannelId          types.Destination
	AdjudicationStatus protocols.AdjudicationStatus
	BlockNum           uint64
}

// DepositedEvent is an internal representation of the deposited blockchain event
type DepositedEvent struct {
	CommonEvent
	Holdings types.Funds // indexed by asset
}

func (de DepositedEvent) GetChannelId() types.Destination {
	return de.ChannelId
}

func (de DepositedEvent) GetBlockNum() uint64 {
	return de.BlockNum
}

// AllocationUpdated is an internal representation of the AllocatonUpdated blockchain event
type AllocationUpdatedEvent struct {
	CommonEvent
	Holdings types.Funds // indexed by asset
}

func (de AllocationUpdatedEvent) GetChannelId() types.Destination {
	return de.ChannelId
}

func (de AllocationUpdatedEvent) GetBlockNum() uint64 {
	return de.BlockNum
}

// todo implement other event types
// AllocationUpdated
// Concluded
// ChallengeRegistered
// ChallengeCleared

type ChainEventHandler interface {
	UpdateWithChainEvent(event Event) (protocols.Objective, error)
}

type ChainService interface {
	Out() <-chan Event
	In() chan<- protocols.ChainTransaction
}
