// Package store contains the interface for a go-nitro store.
package store // import "github.com/statechannels/go-nitro/client/engine/store"

import (
	"errors"

	"github.com/statechannels/go-nitro/channel"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

var ErrNoSuchObjective error = errors.New("store: no such objective")
var ErrNoSuchChannel error = errors.New("store: failed to find required channel data")

// Store is responsible for persisting objectives, objective metadata, states, signatures, private keys and blockchain data
type Store interface {
	GetChannelSecretKey() *[]byte // Get a pointer to a secret key for signing channel updates
	GetAddress() *types.Address   // Get the (Ethereum) address associated with the ChannelSecretKey

	GetObjectiveById(protocols.ObjectiveId) (protocols.Objective, error)          // Read an existing objective
	GetObjectiveByChannelId(types.Destination) (obj protocols.Objective, ok bool) // Get the objective that currently owns the channel with the supplied ChannelId
	SetObjective(protocols.Objective) error                                       // Write an objective
	GetTwoPartyLedger(firstParty types.Address, secondParty types.Address) (channel *channel.TwoPartyLedger, ok bool)
	SetChannel(*channel.Channel) error
}
