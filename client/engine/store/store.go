// Package Store contains the interface for a go-nitro store.
package store // import "github.com/statechannels/go-nitro/client/engine/store"

import (
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

// Store is responsible for persisting objectives, objective metadata, states, signatures, private keys and blockchain data
type Store interface {
	GetObjectiveById(protocols.ObjectiveId) protocols.Objective // Read an existing objective
	GetObjectiveByChannelId(types.Bytes32) protocols.Objective  // Get the objective that currently owns the channel with the supplied ChannelId
	SetObjective(protocols.Objective) error                     // Write an objective

	UpdateProgressLastMadeAt(protocols.ObjectiveId, protocols.WaitingFor) // updates progressLastMadeAt information for an objective
}
