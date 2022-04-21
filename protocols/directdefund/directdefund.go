// Package directdefund implements an off-chain protocol to defund a directly-funded channel.
package directdefund // import "github.com/statechannels/go-nitro/directfund"

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/statechannels/go-nitro/channel"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

const (
	WaitingForFinalization protocols.WaitingFor = "WaitingForFinalization"
	WaitingForWithdraw     protocols.WaitingFor = "WaitingForWithdraw"
	WaitingForNothing      protocols.WaitingFor = "WaitingForNothing" // Finished
)

const ObjectivePrefix = "DirectDefunding-"

// errors
var ErrNotApproved = errors.New("objective not approved")
var ErrChannelUpdateInProgress = errors.New("can only defund a channel when the latest state is supported or when the channel has a final state")

// Objective is a cache of data computed by reading from the store. It stores (potentially) infinite data
type Objective struct {
	Status       protocols.ObjectiveStatus
	C            *channel.Channel
	finalTurnNum uint64
}

// isInConsensusOrFinalState returns true if the channel has a final state or latest state that is supported
func isInConsensusOrFinalState(c *channel.Channel) (bool, error) {
	latestSS, err := c.LatestSignedState()
	// There are no signed states. We consider this as consensus
	if err != nil && err.Error() == "No states are signed" {
		return true, nil
	}
	if latestSS.State().IsFinal {
		return true, nil
	}

	latestSupportedState, err := c.LatestSupportedState()
	if err != nil {
		return false, err
	}

	return cmp.Equal(latestSS.State(), latestSupportedState), nil
}

// GetChannelByIdFunction specifies a function that can be used to retreive channels from a store.
type GetChannelByIdFunction func(id types.Destination) (channel *channel.Channel, ok bool)

// NewObjective initiates an Objective with the supplied channel
func NewObjective(
	preApprove bool,
	channelId types.Destination,
	getChannel GetChannelByIdFunction,
) (Objective, error) {
	c, ok := getChannel(channelId)

	if !ok {
		return Objective{}, fmt.Errorf("could not find channel %s", channelId)
	}
	// We choose to disallow creating an objective if the channel has an in-progress update.
	// We allow the creation of of an objective if the channel has some final states.
	// In the future, we can add a restriction that only defund objectives can add final states to the channel.
	canCreateObjective, err := isInConsensusOrFinalState(c)
	if err != nil {
		return Objective{}, err
	}
	if !canCreateObjective {
		return Objective{}, ErrChannelUpdateInProgress
	}

	var init = Objective{}

	if preApprove {
		init.Status = protocols.Approved
	} else {
		init.Status = protocols.Unapproved
	}
	init.C = c.Clone()

	latestSS, err := c.LatestSupportedState()
	if err != nil {
		return init, err
	}

	if !latestSS.IsFinal {
		init.finalTurnNum = latestSS.TurnNum + 1

	} else {
		init.finalTurnNum = latestSS.TurnNum
	}

	return init, nil
}

var ErrNoFinalState = errors.New("Cannot spawn direct defund objective without a final state")

// ConstructObjectiveFromMessage takes in a message and constructs an objective from it.
func ConstructObjectiveFromMessage(
	m protocols.Message,
	getChannel GetChannelByIdFunction,
) (Objective, error) {
	preApprove := true
	// TODO: do not blindly preapprove
	// See https://github.com/statechannels/go-nitro/issues/213

	// Implicit in the wire protocol is that the message signalling
	// closure of a channel includes an isFinal state (in the 0 slot of the message)
	//
	if !m.SignedStates[0].State().IsFinal {
		return Objective{}, ErrNoFinalState
	}

	cId, err := m.SignedStates[0].State().ChannelId()
	if err != nil {
		return Objective{}, err
	}
	return NewObjective(preApprove, cId, getChannel)
}

// Public methods on the DirectDefundingObjective

// Id returns the unique id of the objective
func (o Objective) Id() protocols.ObjectiveId {
	return protocols.ObjectiveId(ObjectivePrefix + o.C.Id.String())
}

func (o Objective) Approve() protocols.Objective {
	updated := o.clone()
	// todo: consider case of o.Status == Rejected
	updated.Status = protocols.Approved

	return &updated
}

func (o Objective) Reject() protocols.Objective {
	updated := o.clone()
	updated.Status = protocols.Rejected
	return &updated
}

// OwnsChannel returns the channel that the objective is funding.
func (ddo Objective) OwnsChannel() types.Destination {
	return ddo.C.Id
}

// GetStatus returns the status of the objective.
func (ddo Objective) GetStatus() protocols.ObjectiveStatus {
	return ddo.Status
}

func (o Objective) Related() []protocols.Storable {
	return []protocols.Storable{o.C}
}

// Update receives an ObjectiveEvent, applies all applicable event data to the DirectDefundingObjective,
// and returns the updated objective
func (o Objective) Update(event protocols.ObjectiveEvent) (protocols.Objective, error) {
	if o.Id() != event.ObjectiveId {
		return &o, fmt.Errorf("event and objective Ids do not match: %s and %s respectively", string(event.ObjectiveId), string(o.Id()))
	}

	if len(event.SignedStates) > 0 {
		for _, ss := range event.SignedStates {
			if !ss.State().IsFinal {
				return &o, errors.New("direct defund objective can only be updated with final states")
			}
			if o.finalTurnNum != ss.State().TurnNum {
				return &o, fmt.Errorf("expected state with turn number %d, received turn number %d", o.finalTurnNum, ss.State().TurnNum)
			}
		}
	}

	updated := o.clone()
	updated.C.AddSignedStates(event.SignedStates)

	return &updated, nil
}

// UpdateWithChainEvent updates the objective with observed on-chain data.
//
// Only Allocation Updated events are currently handled.
func (o Objective) UpdateWithChainEvent(event chainservice.Event) (protocols.Objective, error) {
	updated := o.clone()
	de, ok := event.(chainservice.AllocationUpdatedEvent)
	if !ok {
		return &updated, fmt.Errorf("objective %+v cannot handle event %+v", updated, event)
	}
	// todo: check block number
	if de.Holdings != nil {
		updated.C.OnChainFunding = de.Holdings.Clone()
	}

	return &updated, nil

}

// Crank inspects the extended state and declares a list of Effects to be executed
func (o Objective) Crank(secretKey *[]byte) (protocols.Objective, protocols.SideEffects, protocols.WaitingFor, error) {
	updated := o.clone()

	sideEffects := protocols.SideEffects{}

	if updated.Status != protocols.Approved {
		return &updated, sideEffects, WaitingForNothing, ErrNotApproved
	}

	latestSignedState, err := updated.C.LatestSignedState()
	if err != nil {
		return &updated, sideEffects, WaitingForNothing, errors.New("The channel must contain at least one signed state to crank the defund objective")
	}

	// Finalize and sign a state if no supported, finalized state exists
	if !latestSignedState.State().IsFinal || !latestSignedState.HasSignatureForParticipant(updated.C.MyIndex) {
		stateToSign := latestSignedState.State().Clone()
		if !stateToSign.IsFinal {
			stateToSign.TurnNum += 1
			stateToSign.IsFinal = true
		}
		ss, err := updated.C.SignAndAddState(stateToSign, secretKey)
		if err != nil {
			return &updated, protocols.SideEffects{}, WaitingForFinalization, fmt.Errorf("could not sign final state %w", err)
		}
		messages := protocols.CreateSignedStateMessages(updated.Id(), ss, updated.C.MyIndex)
		sideEffects.MessagesToSend = append(sideEffects.MessagesToSend, messages...)
	}

	latestSupportedState, err := updated.C.LatestSupportedState()
	if err != nil {
		return &updated, sideEffects, WaitingForFinalization, fmt.Errorf("error finding a supported state: %w", err)
	}
	if !latestSupportedState.IsFinal {
		return &updated, sideEffects, WaitingForFinalization, nil
	}

	// Withdrawal of funds
	if !updated.fullyWithdrawn() {
		// TODO #314: before submiting a withdrawal transaction, we should check if a withdrawal transaction has already been submitted
		// The first participant in the channel submits the withdrawAll transaction
		if updated.C.MyIndex == 0 {
			withdrawAll := protocols.ChainTransaction{Type: protocols.WithdrawAllTransactionType, ChannelId: updated.C.Id}
			sideEffects.TransactionsToSubmit = append(sideEffects.TransactionsToSubmit, withdrawAll)
		}
		// Every participant waits for all channel funds to be distributed, even if the participant has no funds in the channel
		return &updated, sideEffects, WaitingForWithdraw, nil
	}

	return &updated, sideEffects, WaitingForNothing, nil
}

// IsDirectDefundObjective inspects a objective id and returns true if the objective id is for a direct defund objective.
func IsDirectDefundObjective(id protocols.ObjectiveId) bool {
	return strings.HasPrefix(string(id), ObjectivePrefix)
}

//  Private methods on the DirectDefundingObjective

// fullyWithdrawn returns true if the channel contains no assets on chain
func (o Objective) fullyWithdrawn() bool {
	for _, holdings := range o.C.OnChainFunding {
		if holdings.Cmp(big.NewInt(0)) != 0 {
			return false
		}
	}
	return true
}

// clone returns a deep copy of the receiver.
func (o Objective) clone() Objective {
	clone := Objective{}
	clone.Status = o.Status

	cClone := o.C.Clone()
	clone.C = cClone
	clone.finalTurnNum = o.finalTurnNum

	return clone
}

// ObjectiveRequest represents a request to create a new direct defund objective.
type ObjectiveRequest struct {
	ChannelId types.Destination
}

// Id returns the objective id for the request.
func (r ObjectiveRequest) Id() protocols.ObjectiveId {
	return protocols.ObjectiveId(ObjectivePrefix + r.ChannelId.String())
}
