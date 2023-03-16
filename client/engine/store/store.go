// Package store contains the interface for a go-nitro store.
package store // import "github.com/statechannels/go-nitro/client/engine/store"

import (
	"errors"
	"io"

	"github.com/statechannels/go-nitro/channel"
	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

var (
	ErrNoSuchObjective error = errors.New("store: no such objective")
	ErrNoSuchChannel   error = errors.New("store: failed to find required channel data")
)

// Store is responsible for persisting objectives, objective metadata, states, signatures, private keys and blockchain data
type Store interface {
	GetChannelSecretKey() *[]byte // Get a pointer to a secret key for signing channel updates
	GetAddress() *types.Address   // Get the (Ethereum) address associated with the ChannelSecretKey

	GetObjectiveById(protocols.ObjectiveId) (protocols.Objective, error)          // Read an existing objective
	GetObjectiveByChannelId(types.Destination) (obj protocols.Objective, ok bool) // Get the objective that currently owns the channel with the supplied ChannelId
	SetObjective(protocols.Objective) error                                       // Write an objective

	GetChannelById(id types.Destination) (c *channel.Channel, ok bool)
	GetChannelsByParticipant(participant types.Address) []*channel.Channel // Returns any channels that includes the given participant
	SetChannel(*channel.Channel) error
	DestroyChannel(id types.Destination)

	ReleaseChannelFromOwnership(types.Destination) // Release channel from being owned by any objective

	ConsensusChannelStore
	io.Closer
}

type ConsensusChannelStore interface {
	GetConsensusChannel(counterparty types.Address) (channel *consensus_channel.ConsensusChannel, ok bool)
	GetConsensusChannelById(id types.Destination) (channel *consensus_channel.ConsensusChannel, err error)
	SetConsensusChannel(*consensus_channel.ConsensusChannel) error
	DestroyConsensusChannel(id types.Destination)
}
