package virtualdefund

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/channel"
	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/payments"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

const (
	WaitingForLatestVoucher           protocols.WaitingFor = "WaitingForLatestVoucher"
	WaitingForCompleteFinal           protocols.WaitingFor = "WaitingForCompleteFinal"           // Round 1
	WaitingForCompleteLedgerDefunding protocols.WaitingFor = "WaitingForCompleteLedgerDefunding" // Round 2
	WaitingForNothing                 protocols.WaitingFor = "WaitingForNothing"                 // Finished
)

// The turn number used for the final state
const FinalTurnNum = 2

// Objective contains relevant information for the defund objective
type Objective struct {
	Status protocols.ObjectiveStatus

	// InitialOutcome is the initial outcome of the virtual channel
	InitialOutcome outcome.SingleAssetExit

	// VFixed is the fixed channel information for the virtual channel
	VFixed state.FixedPart

	// Signatures are the signatures for the final virtual state from each participant
	// Signatures are ordered by participant order: Signatures[0] is Alice's signature, Signatures[1] is Irene's signature, Signatures[2] is Bob's signature
	// Signatures gets updated as participants sign and send states to each other.
	Signatures [3]state.Signature

	Vouchers [3]payments.Voucher

	// TODO: Do we need this? If so it should probably be private
	VoucherSent bool

	ToMyLeft  *consensus_channel.ConsensusChannel
	ToMyRight *consensus_channel.ConsensusChannel

	// MyRole is the index of the participant in the participants list
	// 0 is Alice
	// 1 is Irene
	// 2 is Bob
	MyRole uint
}

const ObjectivePrefix = "VirtualDefund-"

// GetChannelByIdFunction specifies a function that can be used to retrieve channels from a store.
type GetVoucherByChannelIdFunction func(cId types.Destination) (voucher payments.Voucher, err error)

// GetChannelByIdFunction specifies a function that can be used to retrieve channels from a store.
type GetChannelByIdFunction func(id types.Destination) (channel *channel.Channel, ok bool)

// GetTwoPartyConsensusLedgerFuncion describes functions which return a ConsensusChannel ledger channel between
// the calling client and the given counterparty, if such a channel exists.
type GetTwoPartyConsensusLedgerFunction func(counterparty types.Address) (ledger *consensus_channel.ConsensusChannel, ok bool)

// NewObjective constructs a new virtual defund objective
func NewObjective(request ObjectiveRequest,
	preApprove bool,
	myAddress types.Address,
	getChannel GetChannelByIdFunction,
	getConsensusChannel GetTwoPartyConsensusLedgerFunction,
	getVoucher GetVoucherByChannelIdFunction,
) (Objective, error) {
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

	initialOutcome := V.PostFundState().Outcome[0]

	// This logic assumes a single hop virtual channel.
	// Currently this is the only type of virtual channel supported.
	alice := V.Participants[0]
	intermediary := V.Participants[1]
	bob := V.Participants[2]

	var toMyLeft, toMyRight *consensus_channel.ConsensusChannel
	var ok bool

	switch myAddress {
	case alice:
		toMyRight, ok = getConsensusChannel(intermediary)
		if !ok {
			return Objective{}, fmt.Errorf("could not find a ledger channel between %v and %v", alice, intermediary)
		}
	case intermediary:
		toMyLeft, ok = getConsensusChannel(alice)
		if !ok {
			return Objective{}, fmt.Errorf("could not find a ledger channel between %v and %v", alice, intermediary)
		}
		toMyRight, ok = getConsensusChannel(bob)
		if !ok {
			return Objective{}, fmt.Errorf("could not find a ledger channel between %v and %v", intermediary, bob)
		}
	case bob:
		toMyLeft, ok = getConsensusChannel(intermediary)
		if !ok {
			return Objective{}, fmt.Errorf("could not find a ledger channel between %v and %v", intermediary, bob)
		}
	default:
		return Objective{}, fmt.Errorf("client address not found in an expected participant index")

	}

	var latestVoucher payments.Voucher
	if myAddress == alice {
		var err error
		latestVoucher, err = getVoucher(request.ChannelId)
		if err != nil {
			return Objective{}, fmt.Errorf("could not get the latest voucher for channel %s: %w", request.ChannelId, err)
		}
	}
	vouchers := [3]payments.Voucher{}
	vouchers[V.MyIndex] = latestVoucher
	return Objective{
		Status:         status,
		InitialOutcome: initialOutcome, VFixed: V.FixedPart,
		Signatures: [3]state.Signature{},
		MyRole:     V.MyIndex,
		ToMyLeft:   toMyLeft,
		ToMyRight:  toMyRight,
		Vouchers:   vouchers,
	}, nil

}

// ConstructObjectiveFromVoucher takes in a message and constructs an objective from it.
// It accepts the message, myAddress, and a function to to retrieve ledgers from a store.
func ConstructObjectiveFromVoucher(
	fixedPart state.FixedPart,
	initialVoucher payments.Voucher,
	preapprove bool,
	myAddress types.Address,
	getChannel GetChannelByIdFunction,
	getTwoPartyConsensusLedger GetTwoPartyConsensusLedgerFunction,
	getVoucher GetVoucherByChannelIdFunction,
) (Objective, error) {

	err := fixedPart.Validate()
	if err != nil {
		return Objective{}, err
	}

	channelId := fixedPart.ChannelId()

	if err != nil {
		return Objective{}, err
	}
	return NewObjective(
		ObjectiveRequest{channelId},
		preapprove,
		myAddress,
		getChannel,
		getTwoPartyConsensusLedger,
		getVoucher)
}

// IsVirtualDefundObjective inspects a objective id and returns true if the objective id is for a virtualdefund objective.
func IsVirtualDefundObjective(id protocols.ObjectiveId) bool {
	return strings.HasPrefix(string(id), ObjectivePrefix)
}

// signedFinalState returns the final state for the virtual channel
func (o *Objective) signedFinalState() (state.SignedState, error) {
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
func (o *Objective) finalState() state.State {
	vp := state.VariablePart{Outcome: outcome.Exit{o.finalOutcome()}, TurnNum: FinalTurnNum, IsFinal: true}
	return state.StateFromFixedAndVariablePart(o.VFixed, vp)
}

// Returns the proposal with the largest amount
func (o *Objective) LargestProposal() *payments.Voucher {
	largest := big.NewInt(0)
	index := 0
	for i, v := range o.Vouchers {
		if v.Amount.Cmp(largest) > 0 {
			index = i
			largest = v.Amount
		}
	}
	return &o.Vouchers[index]
}

// finalOutcome returns the outcome for the final state calculated from the InitialOutcome and PaidToBob
func (o *Objective) finalOutcome() outcome.SingleAssetExit {
	finalOutcome := o.InitialOutcome.Clone()
	finalOutcome.Allocations[0].Amount.Sub(finalOutcome.Allocations[0].Amount, o.LargestProposal().Amount)
	finalOutcome.Allocations[1].Amount.Add(finalOutcome.Allocations[1].Amount, o.LargestProposal().Amount)

	return finalOutcome
}

// Id returns the objective id.
func (o *Objective) Id() protocols.ObjectiveId {
	vId := o.VFixed.ChannelId() //TODO: Handle error
	return protocols.ObjectiveId(ObjectivePrefix + vId.String())

}

// Approve returns an approved copy of the objective.
func (o *Objective) Approve() protocols.Objective {
	updated := o.clone()
	// todo: consider case of s.Status == Rejected
	updated.Status = protocols.Approved
	return &updated
}

// Reject returns a rejected copy of the objective.
func (o *Objective) Reject() (protocols.Objective, protocols.SideEffects) {
	updated := o.clone()
	updated.Status = protocols.Rejected
	peers := []common.Address{}
	for i, peer := range o.VFixed.Participants {
		if i != int(o.MyRole) {
			peers = append(peers, peer)
		}
	}
	messages := protocols.CreateRejectionNoticeMessage(o.Id(), peers...)

	return &updated, protocols.SideEffects{MessagesToSend: messages}
}

// OwnsChannel returns the channel that the objective is funding.
func (o *Objective) OwnsChannel() types.Destination {
	vId := o.VFixed.ChannelId()
	return vId
}

// GetStatus returns the status of the objective.
func (o *Objective) GetStatus() protocols.ObjectiveStatus {
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

	for i, v := range o.Vouchers {
		clone.Vouchers[i] = *v.Clone()
	}

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
	clone.VoucherSent = o.VoucherSent

	return clone
}

// Crank inspects the extended state and declares a list of Effects to be executed
func (o *Objective) Crank(secretKey *[]byte) (protocols.Objective, protocols.SideEffects, protocols.WaitingFor, error) {
	updated := o.clone()
	sideEffects := protocols.SideEffects{}

	// Input validation
	if updated.Status != protocols.Approved {
		return &updated, sideEffects, WaitingForNothing, protocols.ErrNotApproved
	}

	if !updated.VoucherSent {

		vMsgs := protocols.CreateVoucherMessage(updated.Vouchers[o.MyRole], updated.VFixed.Participants...)
		// We send out the state we think is the final state
		sMsgs := protocols.CreateSignedStateMessages(
			updated.Id(),
			state.NewSignedState(updated.finalState()),
			updated.MyRole)
		msgs := append(vMsgs, sMsgs...)
		sideEffects.MessagesToSend = append(sideEffects.MessagesToSend, msgs...)
		updated.VoucherSent = true
	}

	if !updated.haveAllVouchers() {
		return &updated, sideEffects, WaitingForLatestVoucher, nil
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
			return o, protocols.SideEffects{}, WaitingForNothing, fmt.Errorf("error updating ledger funding: %w", err)
		}
		sideEffects.Merge(ledgerSideEffects)
	}

	if !updated.isBob() && !updated.isRightDefunded() {
		ledgerSideEffects, err := updated.updateLedgerToRemoveGuarantee(updated.ToMyRight, secretKey)
		if err != nil {
			return o, protocols.SideEffects{}, WaitingForNothing, fmt.Errorf("error updating ledger funding: %w", err)
		}
		sideEffects.Merge(ledgerSideEffects)
	}

	if fullyDefunded := updated.isLeftDefunded() && updated.isRightDefunded(); !fullyDefunded {
		return &updated, sideEffects, WaitingForCompleteLedgerDefunding, nil
	}

	// Mark the objective as done
	updated.Status = protocols.Completed
	return &updated, sideEffects, WaitingForNothing, nil

}

// fullySigned returns whether we have a signature from every partciapant
func (o *Objective) fullySigned() bool {
	for _, sig := range o.Signatures {
		if isZero(sig) {
			return false
		}
	}
	return true
}

// isAlice returns true if the receiver represents participant 0 in the virtualdefund protocol
func (o *Objective) isAlice() bool {
	return o.MyRole == 0
}

// isBob returns true if the receiver represents participant 2 in the virtualdefund protocol
func (o *Objective) isBob() bool {
	return o.MyRole == 2
}

// ledgerProposal generates a ledger proposal to remove the guarantee for V for ledger
func (o *Objective) ledgerProposal(ledger *consensus_channel.ConsensusChannel) consensus_channel.Proposal {
	left := o.finalOutcome().Allocations[0].Amount

	return consensus_channel.NewRemoveProposal(ledger.Id, o.VId(), left)
}

// updateLedgerToRemoveGuarantee updates the ledger channel to remove the guarantee that funds V.
func (o *Objective) updateLedgerToRemoveGuarantee(ledger *consensus_channel.ConsensusChannel, sk *[]byte) (protocols.SideEffects, error) {

	var sideEffects protocols.SideEffects

	proposed := ledger.HasRemovalBeenProposed(o.VId())

	if ledger.IsLeader() {
		if proposed { // If we've already proposed a remove proposal we can return
			return protocols.SideEffects{}, nil
		}

		_, err := ledger.Propose(o.ledgerProposal(ledger), *sk)
		if err != nil {
			return protocols.SideEffects{}, fmt.Errorf("error proposing ledger update: %w", err)
		}
		recipient := ledger.Follower()
		message := protocols.CreateSignedProposalMessage(recipient, ledger.ProposalQueue()...)
		sideEffects.MessagesToSend = append(sideEffects.MessagesToSend, message)

	} else {
		// If the proposal is next in the queue we accept it
		proposedNext := ledger.HasRemovalBeenProposedNext(o.VId())
		if proposedNext {
			sp, err := ledger.SignNextProposal(o.ledgerProposal(ledger), *sk)

			if err != nil {
				return protocols.SideEffects{}, fmt.Errorf("could not sign proposal: %w", err)
			}
			// ledger sideEffect
			if proposals := ledger.ProposalQueue(); len(proposals) != 0 {
				sideEffects.ProposalsToProcess = append(sideEffects.ProposalsToProcess, proposals[0].Proposal)
			}

			// messaging sideEffect
			recipient := ledger.Leader()
			message := protocols.CreateSignedProposalMessage(recipient, sp)
			sideEffects.MessagesToSend = append(sideEffects.MessagesToSend, message)
		}
	}

	return sideEffects, nil
}

// VId returns the channel id of the virtual channel.
func (o *Objective) VId() types.Destination {
	vId := o.VFixed.ChannelId() // TODO Deal with error

	return vId
}

// signedBy returns whether we have a valid signature for the given participant
func (o *Objective) signedBy(participant uint) bool {
	return !isZero(o.Signatures[participant])
}

// signedByMe returns whether the current participant has signed the final state
func (o *Objective) signedByMe() bool {
	return o.signedBy(o.MyRole)

}

func (o *Objective) haveAllVouchers() bool {

	for _, v := range o.Vouchers {
		if v.IsZero() {
			return false
		}
	}
	return true
}

// isRightDefunded returns whether the ledger channel ToMyRight has been defunded
// If ToMyRight==nil then we return true
func (o *Objective) isRightDefunded() bool {
	if o.ToMyRight == nil {
		return true
	}

	included := o.ToMyRight.IncludesTarget(o.VId())
	return !included
}

// isLeftDefunded returns whether the ledger channel ToMyLeft has been defunded
// If ToMyLeft==nil then we return true
func (o *Objective) isLeftDefunded() bool {
	if o.ToMyLeft == nil {
		return true
	}

	included := o.ToMyLeft.IncludesTarget(o.VId())
	return !included
}

// validateSignature returns whether the given signature is valid for the given participant
// If a signature is invalid an error will be returned containing the reason
func (o *Objective) validateSignature(sig state.Signature, participantIndex uint) (bool, error) {
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
func (o *Objective) Update(event protocols.ObjectiveEvent) (protocols.Objective, error) {
	if o.Id() != event.ObjectiveId {
		return o, fmt.Errorf("event and objective Ids do not match: %s and %s respectively", string(event.ObjectiveId), string(o.Id()))
	}

	updated := o.clone()

	if ss := event.SignedState; len(ss.Signatures()) != 0 {
		incomingChannelId := ss.State().ChannelId()
		vChannelId := updated.VFixed.ChannelId()

		if incomingChannelId != vChannelId {
			return o, errors.New("event channelId out of scope of objective")
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
						return o, fmt.Errorf("incoming signature %+v does not match existing %+v", incomingSig, existingSig)
					}
				}
				updated.Signatures[i] = incomingSig

				// Otherwise we validate the incoming signature and update our signatures
				isValid, err := updated.validateSignature(incomingSig, i)
				if isValid {
					// Update the signature
					updated.Signatures[i] = incomingSig
				} else {
					return o, fmt.Errorf("failed to validate signature: %w", err)
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

	if sp := event.SignedProposal; sp.Proposal.Target() == o.VId() {
		var err error
		switch sp.Proposal.LedgerID {
		case types.Destination{}:
			return o, fmt.Errorf("signed proposal is for a zero-addressed ledger channel") // catch this case to avoid unspecified behaviour -- because of Alice or Bob we allow a null channel.
		case toMyLeftId:
			err = updated.ToMyLeft.Receive(sp)
		case toMyRightId:
			err = updated.ToMyRight.Receive(sp)
		default:
			return o, fmt.Errorf("signed proposal is not addressed to a known ledger connection %+v", sp)
		}
		// Ignore stale or future proposals.
		if errors.Is(err, consensus_channel.ErrInvalidTurnNum) {
			return &updated, nil
		}

		if err != nil {
			return o, fmt.Errorf("error incorporating signed proposal %+v into objective: %w", protocols.SummarizeProposal(event.ObjectiveId, sp), err)
		}
	}

	if !event.Voucher.IsZero() {
		signer, err := event.Voucher.RecoverSigner()
		if err != nil {
			return o, fmt.Errorf("error recovering voucher signer: %w", err)
		}
		// This logic assumes a single hop virtual channel.
		// Currently this is the only type of virtual channel supported.
		alice := updated.VFixed.Participants[0]
		intermediary := updated.VFixed.Participants[1]
		bob := updated.VFixed.Participants[2]
		switch {
		case signer == alice:
			updated.Vouchers[0] = event.Voucher
		case signer == intermediary:
			updated.Vouchers[1] = event.Voucher
		case signer == bob:
			updated.Vouchers[2] = event.Voucher
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
}

// Id returns the objective id for the request.
func (r ObjectiveRequest) Id(types.Address) protocols.ObjectiveId {
	return protocols.ObjectiveId(ObjectivePrefix + r.ChannelId.String())
}
