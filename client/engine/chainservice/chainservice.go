// Package chainservice is a chain service responsible for submitting blockchain transactions and relaying blockchain events.
package chainservice // import "github.com/statechannels/go-nitro/client/chainservice"

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

// Event dictates which methods all chain events must implement
type Event interface {
	ChannelID() types.Destination
}

// commonEvent declares fields shared by all chain events
type commonEvent struct {
	channelID types.Destination
	BlockNum  uint64
}

func (ce commonEvent) ChannelID() types.Destination {
	return ce.channelID
}

type assetAndAmount struct {
	AssetAddress common.Address
	AssetAmount  *big.Int
}

// DepositedEvent is an internal representation of the deposited blockchain event
type DepositedEvent struct {
	commonEvent
	assetAndAmount
	NowHeld *big.Int
}

// AllocationUpdated is an internal representation of the AllocatonUpdated blockchain event
// The event includes the token address and amount at the block that generated the event
type AllocationUpdatedEvent struct {
	commonEvent
	assetAndAmount
}

// ConcludedEvent is an internal representation of the Concluded blockchain event
type ConcludedEvent struct {
	commonEvent
}

func NewDepositedEvent(channelId types.Destination, blockNum uint64, assetAddress common.Address, assetAmount *big.Int, nowHeld *big.Int) DepositedEvent {
	return DepositedEvent{commonEvent{channelId, blockNum}, assetAndAmount{AssetAddress: assetAddress, AssetAmount: assetAmount}, nowHeld}
}

func NewAllocationUpdatedEvent(channelId types.Destination, blockNum uint64, assetAddress common.Address, assetAmount *big.Int) AllocationUpdatedEvent {
	return AllocationUpdatedEvent{commonEvent{channelId, blockNum}, assetAndAmount{AssetAddress: assetAddress, AssetAmount: assetAmount}}
}

// todo implement other event types
// ChallengeRegistered
// ChallengeCleared

// ChainEventHandler describes an objective that can handle chain events
type ChainEventHandler interface {
	UpdateWithChainEvent(event Event) (protocols.Objective, error)
}

type ChainService interface {
	// EventFeed returns a chan for receiving events from the chain service.
	EventFeed() <-chan Event
	// SendTransaction is for sending transactions with the chain service
	SendTransaction(protocols.ChainTransaction) error
	// GetConsensusAppAddress returns the address of a deployed ConsensusApp (for ledger channels)
	GetConsensusAppAddress() types.Address
}

type chainServiceBase struct {
	out chan Event
}

// newChainServiceBase constructs a ChainServiceBase. Only implementations of ChainService interface should call the constructor.
func newChainServiceBase() chainServiceBase {
	// Use a buffered channel so we don't have to worry about blocking on writing to the channel.
	return chainServiceBase{out: make(chan Event, 10)}
}

// EventFeed returns the out chan, and narrows the type so that external consumers may only receive on it.
func (csb *chainServiceBase) EventFeed() <-chan Event {
	return csb.out
}

func (csb *chainServiceBase) broadcast(event Event) {
	select {
	case csb.out <- event:
	default:
	}
}
