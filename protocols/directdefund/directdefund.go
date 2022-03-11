// Package directfund implements an off-chain protocol to defund a directly-funded channel.
package directfund // import "github.com/statechannels/go-nitro/directfund"

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/statechannels/go-nitro/channel"
	"github.com/statechannels/go-nitro/protocols"
)

const (
	WaitingForFinalization protocols.WaitingFor = "WaitingForFinalization"
	WaitingForWithdraw     protocols.WaitingFor = "WaitingForWithdraw"
	WaitingForNothing      protocols.WaitingFor = "WaitingForNothing" // Finished
)

const ObjectivePrefix = "DirectDefunding-"

// errors
var ErrNotApproved = errors.New("objective not approved")
var ErrChannelUpdateInProgress = errors.New("cannot defund channel with unsupported, non-final latest state")

// Objective is a cache of data computed by reading from the store. It stores (potentially) infinite data
type Objective struct {
	Status protocols.ObjectiveStatus
	C      *channel.Channel
}

func isUpdateInProgress(c *channel.Channel) bool {
	latestSS, err := c.LatestSignedState()
	// There are no signed states
	if err != nil {
		return false
	}
	return !latestSS.HasAllSignatures() && !latestSS.State().IsFinal
}

// NewObjective initiates an Objective with the supplied channel
func NewObjective(
	preApprove bool,
	c *channel.Channel,
) (Objective, error) {
	// We are chosing to disallow creating an objective if the channel has an in-progress update
	if isUpdateInProgress(c) {
		return Objective{}, ErrChannelUpdateInProgress
	}

	var init = Objective{}

	if preApprove {
		init.Status = protocols.Approved
	} else {
		init.Status = protocols.Unapproved
	}
	init.C = c.Clone()

	return init, nil
}

// Public methods on the DirectDefundingObjective
func (o Objective) Id() protocols.ObjectiveId {
	return protocols.ObjectiveId(ObjectivePrefix + o.C.Id.String())
}

func (o Objective) Approve() Objective {
	updated := o.clone()
	// todo: consider case of o.Status == Rejected
	updated.Status = protocols.Approved

	return updated
}

func (o Objective) Reject() Objective {
	updated := o.clone()
	updated.Status = protocols.Rejected
	return updated
}

// Update receives an ObjectiveEvent, applies all applicable event data to the DirectDefundingObjective,
// and returns the updated objective
func (o Objective) Update(event protocols.ObjectiveEvent) (Objective, error) {
	if o.Id() != event.ObjectiveId {
		return o, fmt.Errorf("event and objective Ids do not match: %s and %s respectively", string(event.ObjectiveId), string(o.Id()))
	}

	if len(event.SignedStates) > 0 {
		for _, ss := range event.SignedStates {
			if !ss.State().IsFinal {
				return o, errors.New("direct defund objective can only be updated with final states")
			}
		}
	}

	updated := o.clone()
	updated.C.AddSignedStates(event.SignedStates)

	return updated, nil
}

// Crank inspects the extended state and declares a list of Effects to be executed
func (o Objective) Crank(secretKey *[]byte) (Objective, protocols.SideEffects, protocols.WaitingFor, error) {
	updated := o.clone()

	sideEffects := protocols.SideEffects{}

	// Input validation
	if updated.Status != protocols.Approved {
		return updated, sideEffects, WaitingForNothing, ErrNotApproved
	}

	latestSignedState, err := updated.C.LatestSignedState()
	if err != nil {
		return updated, sideEffects, WaitingForNothing, errors.New("The channel must contain at least one signed state to crank the defund objective")
	}

	// Sign a final state if it has not been signed
	if !latestSignedState.State().IsFinal || !latestSignedState.HasSignatureForParticipant(updated.C.MyIndex) {
		stateToSign := latestSignedState.State().Clone()
		if !stateToSign.IsFinal {
			stateToSign.TurnNum += 1
			stateToSign.IsFinal = true
		}
		ss, err := updated.C.SignAndAddState(stateToSign, secretKey)
		if err != nil {
			return updated, protocols.SideEffects{}, WaitingForFinalization, fmt.Errorf("could not sign final state %w", err)
		}
		messages := protocols.CreateSignedStateMessages(updated.Id(), ss, updated.C.MyIndex)
		sideEffects.MessagesToSend = append(sideEffects.MessagesToSend, messages...)
	}

	latestSupportedState, err := updated.C.LatestSupportedState()
	if err != nil {
		return updated, sideEffects, WaitingForFinalization, fmt.Errorf("error finding a supported state: %w", err)
	}
	if !latestSupportedState.IsFinal {
		return updated, sideEffects, WaitingForFinalization, nil
	}

	// Withdrawal of funds
	if !updated.fullyWithdrawn() {
		// TODO before submiting a withdrawal transaction, we should check if a withdrawal transaction has already been submitted
		// The first participant in the channel submits the withdrawAll transaction
		if o.C.MyIndex == 0 {
			withdrawAll := protocols.ChainTransaction{Type: protocols.WithdrawAllTransactionType, ChannelId: updated.C.Id}
			sideEffects.TransactionsToSubmit = append(sideEffects.TransactionsToSubmit, withdrawAll)
		}
		// Every participant waits for all channel funds to be distributed, even if the participant has no funds in the channel
		return updated, sideEffects, WaitingForWithdraw, nil
	}

	return updated, sideEffects, WaitingForNothing, nil
}

// Equal returns true if the supplied Objective is deeply equal to the receiver.
func (o Objective) Equal(r Objective) bool {
	return o.Status == r.Status &&
		o.C.Equal(*r.C)
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

	return clone
}
