// Package Store contains the interface for a go-nitro store.
package store // import "github.com/statechannels/go-nitro/client/engine/store"

import (
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

// Store is responsible for persisting objectives, objective metadata, states, signatures, private keys and blockchain data
type Store interface {
	SetChannelSecretKey([]byte)   // Store the channel secret key
	GetChannelSecretKey() *[]byte // Get a pointer to a secret key for signing channel updates

	GetObjectiveById(protocols.ObjectiveId) (obj protocols.Objective, ok bool)    // Read an existing objective
	GetObjectiveByChannelId(types.Destination) (obj protocols.Objective, ok bool) // Get the objective that currently owns the channel with the supplied ChannelId
	SetObjective(protocols.Objective) error                                       // Write an objective

	UpdateProgressLastMadeAt(protocols.ObjectiveId, protocols.WaitingFor) // updates progressLastMadeAt information for an objective
}
