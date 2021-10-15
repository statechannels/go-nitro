package client

import (
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

// Store is responsible for persisting objectives, objective metadata, states, signatures, private keys and blockchain data
type Store interface {
	GetObjectiveById(protocols.ObjectiveId) protocols.Objective // Read an existing objective
	GetObjectiveByChannelId(types.Bytes32) protocols.Objective  // Get the objective that currently owns the channel with the supplied ChannelId
	SetObjective(protocols.Objective) error                     // Write an objective

	EvaluateProgress(protocols.ObjectiveId, protocols.WaitingFor) // checks to see if objective has stalled
}
