package client

import (
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

type Store interface {
	GetObjectiveById(protocols.ObjectiveId) protocols.Objective
	GetObjectiveByChannelId(types.Bytes32) protocols.Objective
	SetObjective(protocols.Objective) error
	ApproveObjective(protocols.ObjectiveId)
	RejectObjective(protocols.ObjectiveId)

	EvaluateProgress(protocols.ObjectiveId, protocols.WaitingFor) // sets waitingFor, checks to see if objective has stalled
	GetWaitingFor(protocols.ObjectiveId) protocols.WaitingFor
	SetWaitingFor(protocols.ObjectiveId, protocols.WaitingFor)
}
