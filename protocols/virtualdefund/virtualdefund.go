package virtualdefund

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/protocols"
)

const (
	WaitingForCompleteFinal           protocols.WaitingFor = "WaitingForCompleteFinal"           // Round 1
	WaitingForCompleteLedgerDefunding protocols.WaitingFor = "WaitingForCompleteLedgerDefunding" // Round 2
	WaitingForNothing                 protocols.WaitingFor = "WaitingForNothing"                 // Finished
)

// The turn number used for the final state
const FinalTurnNum = 3

// Objective contains relevent information for the defund objective
type Objective struct {
	Status protocols.ObjectiveStatus

	// InitialOutcome is the initial outcome of the virtual channel
	InitialOutcome outcome.SingleAssetExit

	// PaidToBob is the amount that should be paid from Alice (participant 0) to Bob (participant 2)
	PaidToBob *big.Int

	// VFixed is the fixed channel information for the virtual channel
	VFixed state.FixedPart

	// Signatures are the signatures for the final virtual state from each participant
	// Signatures are ordered by participant order: Signatures[0] is Alice's signature, Signatures[1] is Irene's signature, Signatures[2] is Bob's signature
	// Signatures gets updated as participants sign and send states to each other.
	Signatures [3]state.Signature

	ToMyLeft  *consensus_channel.ConsensusChannel
	ToMyRight *consensus_channel.ConsensusChannel

	// MyRole is the index of the participant in the participants list
	// 0 is Alice
	// 1 is Irene
	// 2 is Bob
	MyRole uint
}

const ObjectivePrefix = "VirtualDefund-"

//nolint:unused // not used yet
// finalState returns the final state for the virtual channel
func (o Objective) finalState() state.State {
	vp := state.VariablePart{Outcome: outcome.Exit{o.finalOutcome()}, TurnNum: FinalTurnNum, IsFinal: true}
	return state.StateFromFixedAndVariablePart(o.VFixed, vp)
}

//nolint:unused // not used yet
// finalOutcome returns the outcome for the final state calculated from the InitialOutcome and PaidToBob
func (o Objective) finalOutcome() outcome.SingleAssetExit {
	finalOutcome := o.InitialOutcome.Clone()

	finalOutcome.Allocations[0].Amount.Sub(finalOutcome.Allocations[0].Amount, o.PaidToBob)
	finalOutcome.Allocations[1].Amount.Add(finalOutcome.Allocations[1].Amount, o.PaidToBob)

	return finalOutcome
}

// Id returns the objective id.
func (o Objective) Id() protocols.ObjectiveId {
	vId, _ := o.VFixed.ChannelId() //TODO: Handle error
	return protocols.ObjectiveId(ObjectivePrefix + vId.String())

}

// Approve returns an approved copy of the objective.
func (o Objective) Approve() protocols.Objective {
	updated := o.clone()
	// todo: consider case of s.Status == Rejected
	updated.Status = protocols.Approved
	return &updated
}

// Approve returns a rejected copy of the objective.
func (o Objective) Reject() protocols.Objective {
	updated := o.clone()
	updated.Status = protocols.Rejected
	return &updated
}

// Relable returns related channels that need to be stored along with the objective.
func (o *Objective) Related() []protocols.Storable {
	relatable := []protocols.Storable{}

	if o.ToMyLeft != nil {
		relatable = append(relatable, o.ToMyLeft)
	}

	if o.ToMyRight != nil {
		relatable = append(relatable, o.ToMyRight)
	}
	return relatable
}

// Clone returns a deep copy of the receiver.
func (o *Objective) clone() Objective {
	clone := Objective{}
	clone.Status = o.Status

	clone.VFixed = o.VFixed.Clone()
	clone.InitialOutcome = o.InitialOutcome.Clone()
	clone.PaidToBob = big.NewInt(0).Set(o.PaidToBob)

	clone.Signatures = [3]state.Signature{}
	for i, s := range o.Signatures {
		clone.Signatures[i] = s
	}
	clone.MyRole = o.MyRole

	return clone
}

// Crank inspects the extended state and declares a list of Effects to be executed
func (o Objective) Crank(secretKey *[]byte) (protocols.Objective, protocols.SideEffects, protocols.WaitingFor, error) {
	return &o, protocols.SideEffects{}, WaitingForCompleteFinal, errors.New("TODO: UNIMPLEMENTED")
}

// Update receives an protocols.ObjectiveEvent, applies all applicable event data to the VirtualDefundObjective,
// and returns the updated state.
func (o Objective) Update(event protocols.ObjectiveEvent) (protocols.Objective, error) {
	if o.Id() != event.ObjectiveId {
		return &o, fmt.Errorf("event and objective Ids do not match: %s and %s respectively", string(event.ObjectiveId), string(o.Id()))
	}
	return &o, errors.New("TODO: UNIMPLEMENTED")

}
