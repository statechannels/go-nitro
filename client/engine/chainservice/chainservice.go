// Package chainservice is a chain service responsible for submitting blockchain transactions and relaying blockchain events.
package chainservice // import "github.com/statechannels/go-nitro/client/chainservice"

import (
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

// ChainEvent is an internal representation of a blockchain event
type Event struct {
	ChannelId          types.Destination
	Holdings           types.Funds // indexed by asset
	AdjudicationStatus protocols.AdjudicationStatus
}

type ChainService interface {
	// Out is a channel for sending events out of the ChainService
	Out() <-chan Event
	// In is a channel for receiving transactions into the ChainService
	In() chan<- protocols.ChainTransaction
}
