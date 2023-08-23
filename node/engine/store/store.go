// Package store contains the interface for a go-nitro store.
package store // import "github.com/statechannels/go-nitro/node/engine/store"

import (
	"io"

	"github.com/statechannels/go-nitro/channel"
	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/payments"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

const (
	ErrNoSuchObjective    = types.ConstError("store: no such objective")
	ErrNoSuchChannel      = types.ConstError("store: failed to find required channel data")
	ErrLoadVouchers       = types.ConstError("store: could not load vouchers")
	lastBlockProcessedKey = "lastBlockProcessed"
)

// Store is responsible for persisting objectives, objective metadata, states, signatures, private keys and blockchain data
type Store interface {
	GetChannelSecretKey() *[]byte                                                 // Get a pointer to a secret key for signing channel updates
	GetAddress() *types.Address                                                   // Get the (Ethereum) address associated with the ChannelSecretKey
	GetObjectiveById(protocols.ObjectiveId) (protocols.Objective, error)          // Read an existing objective
	GetObjectiveByChannelId(types.Destination) (obj protocols.Objective, ok bool) // Get the objective that currently owns the channel with the supplied ChannelId
	SetObjective(protocols.Objective) error                                       // Write an objective
	GetChannelsByIds(ids []types.Destination) ([]*channel.Channel, error)         // Returns a collection of channels with the given ids
	GetChannelById(id types.Destination) (c *channel.Channel, ok bool)
	GetChannelsByParticipant(participant types.Address) ([]*channel.Channel, error) // Returns any channels that includes the given participant
	SetChannel(*channel.Channel) error
	DestroyChannel(id types.Destination) error
	GetChannelsByAppDefinition(appDef types.Address) ([]*channel.Channel, error) // Returns any channels that includes the given app definition
	ReleaseChannelFromOwnership(types.Destination) error                         // Release channel from being owned by any objective
	GetLastBlockProcessed() (uint64, error)
	SetLastBlockProcessed(uint64) error

	ConsensusChannelStore
	payments.VoucherStore
	io.Closer
}

type ConsensusChannelStore interface {
	GetAllConsensusChannels() ([]*consensus_channel.ConsensusChannel, error)
	GetConsensusChannel(counterparty types.Address) (channel *consensus_channel.ConsensusChannel, ok bool)
	GetConsensusChannelById(id types.Destination) (channel *consensus_channel.ConsensusChannel, err error)
	SetConsensusChannel(*consensus_channel.ConsensusChannel) error
	DestroyConsensusChannel(id types.Destination) error
}
