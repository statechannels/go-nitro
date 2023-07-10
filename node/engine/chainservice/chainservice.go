// Package chainservice is a chain service responsible for submitting blockchain transactions and relaying blockchain events.
package chainservice // import "github.com/statechannels/go-nitro/node/chainservice"

import (
	"fmt"
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

func (aaa assetAndAmount) String() string {
	return aaa.AssetAmount.String() + " units of " + aaa.AssetAddress.Hex() + " token"
}

// DepositedEvent is an internal representation of the deposited blockchain event
type DepositedEvent struct {
	commonEvent
	Asset   types.Address
	NowHeld *big.Int
}

func (de DepositedEvent) String() string {
	return "Deposited " + de.Asset.String() + " leaving " + de.NowHeld.String() + " now held against channel " + de.channelID.String() + " at Block " + fmt.Sprint(de.BlockNum)
}

// AllocationUpdated is an internal representation of the AllocatonUpdated blockchain event
// The event includes the token address and amount at the block that generated the event
type AllocationUpdatedEvent struct {
	commonEvent
	assetAndAmount
}

func (aue AllocationUpdatedEvent) String() string {
	return "Channel " + aue.channelID.String() + " has had allocation updated to " + aue.assetAndAmount.String() + " at Block " + fmt.Sprint(aue.BlockNum)
}

// ConcludedEvent is an internal representation of the Concluded blockchain event
type ConcludedEvent struct {
	commonEvent
}

func (ce ConcludedEvent) String() string {
	return "Channel " + ce.channelID.String() + " concluded at Block " + fmt.Sprint(ce.BlockNum)
}

func NewDepositedEvent(channelId types.Destination, blockNum uint64, assetAddress common.Address, nowHeld *big.Int) DepositedEvent {
	return DepositedEvent{commonEvent{channelId, blockNum}, assetAddress, nowHeld}
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
	// GetVirtualPaymentAppAddress returns the address of a deployed VirtualPaymentApp
	GetVirtualPaymentAppAddress() types.Address
	// GetChainId returns the id of the chain the service is connected to
	GetChainId() (*big.Int, error)
	// Close closes the ChainService
	Close() error
}
