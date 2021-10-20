package protocols

import (
	"math/big"

	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/types"
)

// SideEffects is an list of effects to be executed by an imperative shell
type SideEffects []string

// WaitingFor is an enumerable "pause-point" computed from an Objective. It describes how the objective is blocked on actions by third parties (i.e. co-participants or the blockchain).
type WaitingFor string

// AdjudicationStatus mirrors the on chain adjudication status of a particular channel.
// Everything that is stored on chain, other than holdings.
type AdjudicationStatus struct {
	TurnNumRecord uint
	// TODO This struct is a placeholder for the time being, until we add a chain-service
	// TODO eventually this struct will contain the other fields stored in (or committed to by) the adjudicator
}

// ObjectiveEvent holds information used to update an Objective. Some fields may be nil.
type ObjectiveEvent struct {
	ChannelId          string                            // Must be defined
	Sigs               map[types.Bytes32]state.Signature // mapping from state hash to signature
	Holdings           map[types.Address]big.Int         // mapping from asset identifier to amount
	AdjudicationStatus AdjudicationStatus
}

// Objective is the interface for off-chain protocols.
// The lifecycle of an objective is as follows:
// 	* It is initialized by a single client (passing in various parameters). It is implicitly approved by that client. It is communicated to the other clients.
// 	* It is stored and then approved or rejected by the other clients
// 	* It is updated with external information arriving to the client
// 	* After each update, it is cranked. This generates side effects and other metadata
// 	* The metadata will eventually indicate that the Objective has stalled OR the Objective has completed successfully
type Objective interface {
	Id() ObjectiveId
	Initialize(initialState state.State) Objective // returns the initial Objective, does not declare effects

	Approve() Objective                    // returns an updated Objective (a copy, no mutation allowed), does not declare effects
	Reject() Objective                     // returns an updated Objective (a copy, no mutation allowed), does not declare effects
	Update(event ObjectiveEvent) Objective // returns an updated Objective (a copy, no mutation allowed), does not declare effects

	Crank() (SideEffects, WaitingFor, error) // does *not* accept an event, but *does* declare side effects, does *not* return an updated Protocol
}

// ObjectiveId is a unique identifier for an Objective.
type ObjectiveId string
