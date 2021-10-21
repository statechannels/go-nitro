// Package chain is a chain service responsible for submitting blockchain transactions and relaying blockchain events
package chain // import "github.com/statechannels/go-nitro/client/chain"

import (
	"math/big"

	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

// ChainEvent is an internal representation of a blockchain event
type Event struct {
	ChannelId          types.Bytes32
	Holdings           map[types.Address]big.Int // indexed by asset
	AdjudicationStatus protocols.AdjudicationStatus
}

type Chain interface {
	GetRecieveChan() chan Event
	Submit(tx protocols.Transaction)
}
