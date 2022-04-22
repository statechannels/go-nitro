// Package store contains the interface for a go-nitro store.
package store // import "github.com/statechannels/go-nitro/client/engine/store"

import (
	"errors"

	"github.com/statechannels/go-nitro/channel"
	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

var ErrNoSuchObjective error = errors.New("store: no such objective")
var ErrNoSuchChannel error = errors.New("store: failed to find required channel data")

// Store is responsible for persisting objectives, objective metadata, states, signatures, private keys and blockchain data
type Store interface {
	GetChannelSecretKey() *[]byte // Get a pointer to a secret key for signing channel updates
	GetAddress() *types.Address   // Get the (Ethereum) address associated with the ChannelSecretKey

	ObjectiveGetter
	ObjectiveSetter

	ChannelGetter
	ChannelSetter
	ReleaseChannelFromOwnership(types.Destination) // Release channel from being owned by any objective

	ConsensusChannelGetter
	ConsensusChannelSetter
	GetTwoPartyLedger(firstParty types.Address, secondParty types.Address) (channel *channel.TwoPartyLedger, ok bool) // deprecated

}

type ObjectiveGetter interface {
	GetObjectiveById(protocols.ObjectiveId) (protocols.Objective, error)          // Read an existing objective
	GetObjectiveByChannelId(types.Destination) (obj protocols.Objective, ok bool) // Get the objective that currently owns the channel with the supplied ChannelId

}

type ObjectiveSetter interface {
	SetObjective(protocols.Objective) error // Write an objective
}

type ChannelGetter interface {
	GetChannelById(id types.Destination) (c *channel.Channel, ok bool)
}

type ChannelSetter interface {
	SetChannel(*channel.Channel) error
}

type ConsensusChannelGetter interface {
	GetConsensusChannel(counterparty types.Address) (channel *consensus_channel.ConsensusChannel, ok bool)
}
type ConsensusChannelSetter interface {
	SetConsensusChannel(*consensus_channel.ConsensusChannel) error
}
