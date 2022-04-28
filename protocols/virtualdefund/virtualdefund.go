package virtualdefund

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/statechannels/go-nitro/channel"
	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

const (
	WaitingForCompleteFinal           protocols.WaitingFor = "WaitingForCompleteFinal"           // Round 1
	WaitingForCompleteLedgerDefunding protocols.WaitingFor = "WaitingForCompleteLedgerDefunding" // Round 2
	WaitingForNothing                 protocols.WaitingFor = "WaitingForNothing"                 // Finished
)

// The turn number used for the final state
const FinalTurnNum = 2

// errors
var ErrNotApproved = errors.New("objective not approved")

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

// GetChannelByIdFunction specifies a function that can be used to retreive channels from a store.
type GetChannelByIdFunction func(id types.Destination) (channel *channel.Channel, ok bool)

// GetTwoPartyConsensusLedgerFuncion describes functions which return a ConsensusChannel ledger channel between
// the calling client and the given counterparty, if such a channel exists.
type GetTwoPartyConsensusLedgerFunction func(counterparty types.Address) (ledger *consensus_channel.ConsensusChannel, ok bool)

// NewObjective constructs a new virtual defund objective
func NewObjective(preApprove bool,
	request ObjectiveRequest,
	getChannel GetChannelByIdFunction,
	getConsensusChannel GetTwoPartyConsensusLedgerFunction) (Objective, error) {
	var status protocols.ObjectiveStatus

	if preApprove {
		status = protocols.Approved
	} else {
		status = protocols.Unapproved
	}

	V, found := getChannel(request.ChannelId)
	if !found {
		return Objective{}, fmt.Errorf("could not find channel %s", request.ChannelId)
	}
	if V.MyIndex != 0 {
		return Objective{}, errors.New("only participant 0 may instigate closing the channel")
	}
	initialOutcome := V.PostFundState().Outcome[0]

	var toMyLeft *consensus_channel.ConsensusChannel // Only Alice calls NewObjective, so there is no ledger channel to her left.

	intermediary := V.Participants[1] // This mirrors how it is done in virtualfund. We should probably extract to a shared constant somewhere.

	toMyRight, found := getConsensusChannel(intermediary)
	if !found {
		return Objective{}, fmt.Errorf("Error getting consensus channel")
	}

	return Objective{
		Status:         status,
		InitialOutcome: initialOutcome,
		PaidToBob:      request.PaidToBob,
		VFixed:         V.FixedPart,
		Signatures:     [3]state.Signature{},
		MyRole:         V.MyIndex,
		ToMyLeft:       toMyLeft,
		ToMyRight:      toMyRight,
	}, nil

}

// ConstructObjectiveFromState takes in a message and constructs an objective from it.
// It accepts the message, myAddress, and a function to to retrieve ledgers from a store.
func ConstructObjectiveFromState(
	initialState state.State,
	getChannel GetChannelByIdFunction,
	getTwoPartyConsensusLedger GetTwoPartyConsensusLedgerFunction,
) (Objective, error) {
	channelId, err := initialState.ChannelId()
	if err != nil {
		return Objective{}, err
	}
	paidToBob := big.NewInt(0) // TODO how to set this properly?
	return NewObjective(true,
		ObjectiveRequest{channelId, paidToBob},
		getChannel,
		getTwoPartyConsensusLedger)
}

// IsVirtualDefundObjective inspects a objective id and returns true if the objective id is for a virtualdefund objective.
func IsVirtualDefundObjective(id protocols.ObjectiveId) bool {
	return strings.HasPrefix(string(id), ObjectivePrefix)
}

// signedFinalState returns the final state for the virtual channel
func (o Objective) signedFinalState() (state.SignedState, error) {
	signed := state.NewSignedState(o.finalState())
	for _, sig := range o.Signatures {
		if !isZero(sig) {
			err := signed.AddSignature(sig)
			if err != nil {
				return state.SignedState{}, err
			}
		}
	}
	return signed, nil
}

// finalState returns the final state for the virtual channel
func (o Objective) finalState() state.State {
	vp := state.VariablePart{Outcome: outcome.Exit{o.finalOutcome()}, TurnNum: FinalTurnNum, IsFinal: true}
	return state.StateFromFixedAndVariablePart(o.VFixed, vp)
}

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

// OwnsChannel returns the channel that the objective is funding.
func (o Objective) OwnsChannel() types.Destination {
	vId, _ := o.VFixed.ChannelId()
	return vId
}

// GetStatus returns the status of the objective.
func (o Objective) GetStatus() protocols.ObjectiveStatus {
	return o.Status
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
		clone.Signatures[i] = state.CloneSignature(s)
	}
	clone.MyRole = o.MyRole

	// TODO: Properly clone the consensus channels
	if o.ToMyLeft != nil {
		clone.ToMyLeft = o.ToMyLeft
	}
	if o.ToMyRight != nil {
		clone.ToMyRight = o.ToMyRight
	}

	return clone
}

// Crank inspects the extended state and declares a list of Effects to be executed
func (o Objective) Crank(secretKey *[]byte) (protocols.Objective, protocols.SideEffects, protocols.WaitingFor, error) {
	updated := o.clone()
	sideEffects := protocols.SideEffects{}

	// Input validation
	if updated.Status != protocols.Approved {
		return &updated, sideEffects, WaitingForNothing, ErrNotApproved
	}

	// Signing of the final state
	if !updated.signedByMe() {

		sig, err := o.finalState().Sign(*secretKey)
		if err != nil {
			return &updated, sideEffects, WaitingForNothing, fmt.Errorf("could not sign final state: %w", err)
		}
		// Update the signature stored on the objective
		updated.Signatures[updated.MyRole] = sig

		// Send out the signature (using a signed state) to fellow participants
		signedFinal, err := updated.signedFinalState()
		if err != nil {
			return &updated, sideEffects, WaitingForNothing, fmt.Errorf("could not generate signed final state: %w", err)
		}
		messages := protocols.CreateSignedStateMessages(updated.Id(), signedFinal, updated.MyRole)
		sideEffects.MessagesToSend = append(sideEffects.MessagesToSend, messages...)
	}

	// Check if all participants have signed the final state
	if !updated.fullySigned() {
		return &updated, sideEffects, WaitingForCompleteFinal, nil
	}

	if !updated.isAlice() && !updated.isLeftDefunded() {
		ledgerSideEffects, err := updated.updateLedgerToRemoveGuarantee(updated.ToMyLeft, secretKey)
		if err != nil {
			return &o, protocols.SideEffects{}, WaitingForNothing, fmt.Errorf("error updating ledger funding: %w", err)
		}
		sideEffects.Merge(ledgerSideEffects)
	}

	if !updated.isBob() && !updated.isRightDefunded() {
		ledgerSideEffects, err := updated.updateLedgerToRemoveGuarantee(updated.ToMyRight, secretKey)
		if err != nil {
			return &o, protocols.SideEffects{}, WaitingForNothing, fmt.Errorf("error updating ledger funding: %w", err)
		}
		sideEffects.Merge(ledgerSideEffects)
	}

	if fullyDefunded := updated.isLeftDefunded() && updated.isRightDefunded(); !fullyDefunded {
		return &updated, sideEffects, WaitingForCompleteLedgerDefunding, nil
	}
	return &updated, sideEffects, WaitingForNothing, nil

}

// fullySigned returns whether we have a signature from every partciapant
func (o Objective) fullySigned() bool {
	for _, sig := range o.Signatures {
		if isZero(sig) {
			return false
		}
	}
	return true
}

// isAlice returns true if the receiver represents participant 0 in the virtualdefund protocol
func (o Objective) isAlice() bool {
	return o.MyRole == 0
}

// isBob returns true if the receiver represents participant 2 in the virtualdefund protocol
func (o Objective) isBob() bool {
	return o.MyRole == 2
}

// ledgerProposal generates a ledger proposal to remove the guarantee for V for ledger
func (o Objective) ledgerProposal(ledger *consensus_channel.ConsensusChannel) consensus_channel.Proposal {
	left := o.finalOutcome().Allocations[0].Amount
	right := o.finalOutcome().Allocations[1].Amount

	return consensus_channel.NewRemoveProposal(ledger.Id, FinalTurnNum, o.VId(), left, right)
}

// updateLedgerToRemoveGuarantee updates the ledger channel to remove the guarantee that funds V.
func (o *Objective) updateLedgerToRemoveGuarantee(ledger *consensus_channel.ConsensusChannel, sk *[]byte) (protocols.SideEffects, error) {

	var sideEffects protocols.SideEffects

	proposed := ledger.HasRemovalBeenProposedFor(o.VId())

	if ledger.IsLeader() {
		if proposed { // If we've already proposed a remove proposal we can return
			return protocols.SideEffects{}, nil
		}

		sp, err := ledger.Propose(o.ledgerProposal(ledger), *sk)
		if err != nil {
			return protocols.SideEffects{}, fmt.Errorf("error proposing ledger update: %w", err)
		}

		message := protocols.CreateSignedProposalMessage(ledger, sp)
		sideEffects.MessagesToSend = append(sideEffects.MessagesToSend, message)

	} else {

		if proposed {
			sp, err := ledger.SignNextProposal(o.ledgerProposal(ledger), *sk)

			if err != nil {
				return protocols.SideEffects{}, fmt.Errorf("could not sign proposal: %w", err)
			}

			message := protocols.CreateSignedProposalMessage(ledger, sp)
			sideEffects.MessagesToSend = append(sideEffects.MessagesToSend, message)
		}
	}

	return sideEffects, nil
}

// VId returns the channel id of the virtual channel.
func (o Objective) VId() types.Destination {
	vId, _ := o.VFixed.ChannelId() // TODO Deal with error

	return vId
}

// signedBy returns whether we have a valid signature for the given participant
func (o Objective) signedBy(participant uint) bool {
	return !isZero(o.Signatures[participant])
}

// signedByMe returns whether the current participant has signed the final state
func (o Objective) signedByMe() bool {
	return o.signedBy(o.MyRole)

}

// isRightDefunded returns whether the ledger channel ToMyRight has been defunded
// If ToMyRight==nil then we return true
func (o Objective) isRightDefunded() bool {
	if o.ToMyRight == nil {
		return true
	}

	included := o.ToMyRight.IncludesTarget(o.VId())
	return !included
}

// isLeftDefunded returns whether the ledger channel ToMyLeft has been defunded
// If ToMyLeft==nil then we return true
func (o Objective) isLeftDefunded() bool {
	if o.ToMyLeft == nil {
		return true
	}

	included := o.ToMyLeft.IncludesTarget(o.VId())
	return !included
}

// validateSignature returns whether the given signature is valid for the given participant
// If a signature is invalid an error will be returned containing the reason
func (o Objective) validateSignature(sig state.Signature, participantIndex uint) (bool, error) {
	if participantIndex > 2 {
		return false, fmt.Errorf("participant index %d is out of bounds", participantIndex)
	}

	finalState := o.finalState()
	signer, err := finalState.RecoverSigner(sig)
	if err != nil {
		return false, fmt.Errorf("failed to recover signer from signature: %w", err)
	}
	if signer != o.VFixed.Participants[participantIndex] {
		return false, fmt.Errorf("signature is for %s, expected signature from %s ", signer, o.VFixed.Participants[participantIndex])
	}
	return true, nil
}

// Update receives an protocols.ObjectiveEvent, applies all applicable event data to the VirtualDefundObjective,
// and returns the updated state.
func (o Objective) Update(event protocols.ObjectiveEvent) (protocols.Objective, error) {
	if o.Id() != event.ObjectiveId {
		return &o, fmt.Errorf("event and objective Ids do not match: %s and %s respectively", string(event.ObjectiveId), string(o.Id()))
	}

	updated := o.clone()

	for _, ss := range event.SignedStates {
		incomingChannelId, _ := ss.State().ChannelId() // TODO handle error
		vChannelId, _ := updated.VFixed.ChannelId()    // TODO handle error

		if incomingChannelId != vChannelId {
			return &o, errors.New("event channelId out of scope of objective")
		} else {
			incomingSignatures := ss.Signatures()
			for i := uint(0); i < 3; i++ {
				existingSig := o.Signatures[i]
				incomingSig := incomingSignatures[i]

				// If the incoming signature is zeroed we ignore it
				if isZero(incomingSig) {
					continue
				}
				// If the existing signature is not zeroed we check that it matches the incoming signature
				if !isZero(existingSig) {
					if existingSig.Equal(incomingSig) {
						continue
					} else {
						return &o, fmt.Errorf("incoming signature %+v does not match existing %+v", incomingSig, existingSig)
					}
				}
				// Otherwise we validate the incoming signature and update our signatures
				isValid, err := updated.validateSignature(incomingSig, i)
				if isValid {
					// Update the signature
					updated.Signatures[i] = incomingSig
				} else {
					return &o, fmt.Errorf("failed to validate signature: %w", err)
				}
			}
		}
	}
	var toMyLeftId types.Destination
	var toMyRightId types.Destination

	if o.ToMyLeft != nil {
		toMyLeftId = o.ToMyLeft.Id
	}
	if o.ToMyRight != nil {
		toMyRightId = o.ToMyRight.Id
	}

	for _, sp := range event.SignedProposals {
		var err error
		switch sp.Proposal.ChannelID {
		case types.Destination{}:
			return &o, fmt.Errorf("signed proposal is for a zero-addressed ledger channel") // catch this case to avoid unspecified behaviour -- because of Alice or Bob we allow a null channel.
		case toMyLeftId:
			err = updated.ToMyLeft.Receive(sp)
		case toMyRightId:
			err = updated.ToMyRight.Receive(sp)
		default:
			return &o, fmt.Errorf("signed proposal is not addressed to a known ledger connection")
		}

		if err != nil {
			return &o, fmt.Errorf("error incorporating signed proposal into objective: %w", err)
		}
	}
	return &updated, nil

}

// isZero returns true if every byte field on the signature is zero
func isZero(sig state.Signature) bool {
	zeroSig := state.Signature{}
	return sig.Equal(zeroSig)
}

// ObjectiveRequest represents a request to create a new direct defund objective.
type ObjectiveRequest struct {
	ChannelId types.Destination
	PaidToBob *big.Int
}

// Id returns the objective id for the request.
func (r ObjectiveRequest) Id() protocols.ObjectiveId {
	return protocols.ObjectiveId(ObjectivePrefix + r.ChannelId.String())
}
