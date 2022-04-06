package virtualdefund

import (
	"errors"
	"fmt"

	"github.com/statechannels/go-nitro/channel"
	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

const (
	WaitingForCompleteFinal           protocols.WaitingFor = "WaitingForCompleteFinal"           // Round 1
	WaitingForCompleteLedgerDefunding protocols.WaitingFor = "WaitingForCompleteLedgerDefunding" // Round 2
	WaitingForNothing                 protocols.WaitingFor = "WaitingForNothing"                 // Finished
)

type GuaranteeInfo struct {
	Left                 types.Destination
	Right                types.Destination
	LeftAmount           types.Funds
	RightAmount          types.Funds
	GuaranteeDestination types.Destination
}

type Connection struct {
	ConsensusChannel *consensus_channel.ConsensusChannel
	GuaranteeInfo    GuaranteeInfo
}

// Objective is a cache of data computed by reading from the store. It stores (potentially) infinite data.
type Objective struct {
	Status protocols.ObjectiveStatus
	V      *channel.SingleHopVirtualChannel

	ToMyLeft  *Connection
	ToMyRight *Connection

	MyRole uint // index in the virtual funding protocol. 0 for Alice, 2 for Bob. Otherwise 1 for the intermediary.

}

const ObjectivePrefix = "VirtualDefund-"

// Id returns the objective id.
func (o Objective) Id() protocols.ObjectiveId {
	return protocols.ObjectiveId(ObjectivePrefix + o.V.Id.String())

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

func (o *Objective) Related() []protocols.Storable {
	relatable := []protocols.Storable{o.V}

	if o.ToMyLeft != nil {
		relatable = append(relatable, o.ToMyLeft.ConsensusChannel)
	}

	if o.ToMyRight != nil {
		relatable = append(relatable, o.ToMyRight.ConsensusChannel)
	}
	return relatable
}

// Clone returns a deep copy of the receiver.
func (o *Objective) clone() Objective {
	clone := Objective{}
	clone.Status = o.Status
	vClone := o.V.Clone()
	clone.V = vClone

	if o.ToMyLeft != nil {
		// todo: #420 consider cloning for consensusChannels
		clone.ToMyLeft = &Connection{
			GuaranteeInfo: o.ToMyLeft.GuaranteeInfo,
		}
	}

	if o.ToMyRight != nil {
		// todo: #420 consider cloning for consensusChannels
		clone.ToMyRight = &Connection{
			GuaranteeInfo: o.ToMyRight.GuaranteeInfo,
		}
	}

	clone.MyRole = o.MyRole

	return clone
}

func (o Objective) Crank(secretKey *[]byte) (protocols.Objective, protocols.SideEffects, protocols.WaitingFor, error) {
	return &o, protocols.SideEffects{}, WaitingForCompleteFinal, errors.New("TODO: UNIMPLEMENTED")
}

// Update receives an protocols.ObjectiveEvent, applies all applicable event data to the VirtualFundObjective,
// and returns the updated state.
func (o Objective) Update(event protocols.ObjectiveEvent) (protocols.Objective, error) {
	if o.Id() != event.ObjectiveId {
		return &o, fmt.Errorf("event and objective Ids do not match: %s and %s respectively", string(event.ObjectiveId), string(o.Id()))
	}

	return &o, errors.New("TODO: UNIMPLEMENTED")

}
