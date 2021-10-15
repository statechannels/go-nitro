package protocols

import (
	"math/big"

	"github.com/statechannels/go-nitro/channel/state"
)

// TODO these are placeholders for now
type SideEffects []interface{}
type WaitingFor string

type Status struct {
	TurnNumRecord uint // TODO add other fields
}
type ObjectiveEvent struct {
	ChannelId          string
	Sigs               map[*state.State]state.Signature // mapping from state to signature TODO consider using a hash of the state
	Holdings           big.Int                          // TODO allow for multiple assets
	AdjudicationStatus Status
}

// Objective is the interface for off-chain protocols
type Objective interface {
	Id() ObjectiveId
	Initialize(initialState state.State) Objective // returns the initial Protocol object, does not declare effects

	Approve()                              // returns an updated Protocol (a copy, no mutation allowed), does not declare effects
	Reject()                               // returns an updated Protocol (a copy, no mutation allowed), does not declare effects
	Update(event ObjectiveEvent) Objective // returns an updated Protocol (a copy, no mutation allowed), does not declare effects

	Crank() (SideEffects, WaitingFor, error) // does *not* accept an event, but *does* declare side effects, does *not* return an updated Protocol
}

// TODO these are placeholders for now (they are the fundamental events the wallet reacts to)
type ObjectiveId string
