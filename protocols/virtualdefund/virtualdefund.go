package virtualdefund

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/channel"
	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

const (
	WaitingForFinalStateFromAlice protocols.WaitingFor = "WaitingForFinalStateFromAlice"
	WaitingForSupportedFinalState protocols.WaitingFor = "WaitingForSupportedFinalState" // Round 1
	WaitingForDefundingOnMyLeft   protocols.WaitingFor = "WaitingForDefundingOnMyLeft"   // Round 2
	WaitingForDefundingOnMyRight  protocols.WaitingFor = "WaitingForDefundingOnMyRight"  // Round 2
	WaitingForNothing             protocols.WaitingFor = "WaitingForNothing"             // Finished
)

const (
	// SignedStatePayload indicates that the payload is a json serialized signed state
	SignedStatePayload protocols.PayloadType = "SignedStatePayload"
	// RequestFinalStatePayload indicates that the payload is a request for the final state
	// The actual payload is simply the channel id that the final state is for
	RequestFinalStatePayload protocols.PayloadType = "RequestFinalStatePayload"
)

// The turn number used for the final state
const FinalTurnNum = 2

// Objective contains relevant information for the defund objective
type Objective struct {
	Status protocols.ObjectiveStatus

	// MinimumPaymentAmount is the latest payment amount we have received from Alice before starting defunding.
	// This is set by Bob so he can ensure he receives the latest amount from any vouchers he's received.
	// If this is not set then virtual defunding will accept any final outcome from Alice.
	MinimumPaymentAmount *big.Int

	V *channel.Channel // This is a Virtual Channel. TODO should this be a literal channel.VirtualChannel?

	ToMyLeft  *consensus_channel.ConsensusChannel
	ToMyRight *consensus_channel.ConsensusChannel

	// MyRole is the index of the participant in the participants list
	// 0 is Alice
	// 1...n is Irene, Ivan, ... (the n intermediaries)
	// n+1 is Bob
	MyRole uint
}

const ObjectivePrefix = "VirtualDefund-"

// GetChannelByIdFunction specifies a function that can be used to retrieve channels from a store.
type GetChannelByIdFunction func(id types.Destination) (channel *channel.Channel, ok bool)

// GetTwoPartyConsensusLedgerFuncion describes functions which return a ConsensusChannel ledger channel between
// the calling client and the given counterparty, if such a channel exists.
type GetTwoPartyConsensusLedgerFunction func(counterparty types.Address) (ledger *consensus_channel.ConsensusChannel, ok bool)

// NewObjective constructs a new virtual defund objective
func NewObjective(request ObjectiveRequest,
	preApprove bool,
	myAddress types.Address,
	largestPaymentAmount *big.Int,
	getChannel GetChannelByIdFunction,
	getConsensusChannel GetTwoPartyConsensusLedgerFunction,
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

	alice := V.Participants[0]
	bob := V.Participants[len(V.Participants)-1]

	var leftLedger, rightLedger *consensus_channel.ConsensusChannel
	var ok bool

	if myAddress == alice {
		rightOfAlice := V.Participants[1]
		rightLedger, ok = getConsensusChannel(rightOfAlice)
		if !ok {
			return Objective{}, fmt.Errorf("could not find a ledger channel between %v and %v", alice, rightOfAlice)
		}
	} else if myAddress == bob {
		leftOfBob := V.Participants[len(V.Participants)-2]
		leftLedger, ok = getConsensusChannel(leftOfBob)
		if !ok {
			return Objective{}, fmt.Errorf("could not find a ledger channel between %v and %v", leftOfBob, bob)
		}
	} else {
		intermediaries := V.Participants[1 : len(V.Participants)-1]
		foundMyself := false

		for i, intermediary := range intermediaries {
			if myAddress == intermediary {
				foundMyself = true
				// I am intermediary `i` and participant `p`
				p := i + 1 // participants[p] === intermediaries[i]

				leftOfMe := V.Participants[p-1]
				rightOfMe := V.Participants[p+1]

				leftLedger, ok = getConsensusChannel(leftOfMe)
				if !ok {
					return Objective{}, fmt.Errorf("could not find a ledger channel between %v and %v", leftOfMe, myAddress)
				}
				rightLedger, ok = getConsensusChannel(rightOfMe)
				if !ok {
					return Objective{}, fmt.Errorf("could not find a ledger channel between %v and %v", myAddress, rightOfMe)
				}

				break
			}
		}

		if !foundMyself {
			return Objective{}, fmt.Errorf("client address not found in an expected participant index")
		}
	}

	if largestPaymentAmount == nil {
		largestPaymentAmount = big.NewInt(0)
	}

	return Objective{
		Status:               status,
		MinimumPaymentAmount: largestPaymentAmount,
		V:                    V,
		MyRole:               V.MyIndex,
		ToMyLeft:             leftLedger,
		ToMyRight:            rightLedger,
	}, nil
}

// ConstructObjectiveFromPayload takes in a message payload and constructs an objective from it.
func ConstructObjectiveFromPayload(
	p protocols.ObjectivePayload,
	preapprove bool,
	myAddress types.Address,
	getChannel GetChannelByIdFunction,
	getTwoPartyConsensusLedger GetTwoPartyConsensusLedgerFunction,
	latestVoucherAmount *big.Int,
) (Objective, error) {
	if latestVoucherAmount == nil {
		latestVoucherAmount = big.NewInt(0)
	}
	var cId types.Destination
	var err error
	switch p.Type {
	case RequestFinalStatePayload:
		cId, err = getRequestFinalStatePayload(p.PayloadData)
	case SignedStatePayload:
		var ss state.SignedState
		ss, err = getSignedStatePayload(p.PayloadData)
		cId = ss.ChannelId()
	default:
		return Objective{}, fmt.Errorf("unknown payload type %s", p.Type)
	}
	if err != nil {
		return Objective{}, err
	}
	return NewObjective(
		NewObjectiveRequest(cId),
		preapprove,
		myAddress,
		latestVoucherAmount,
		getChannel,
		getTwoPartyConsensusLedger)
}

// IsVirtualDefundObjective inspects a objective id and returns true if the objective id is for a virtualdefund objective.
func IsVirtualDefundObjective(id protocols.ObjectiveId) bool {
	return strings.HasPrefix(string(id), ObjectivePrefix)
}

// finalState returns the final state for the virtual channel
func (o *Objective) finalState() state.State {
	return o.V.SignedStateForTurnNum[FinalTurnNum].State()
}

func (o *Objective) initialOutcome() outcome.SingleAssetExit {
	return o.V.PostFundState().Outcome[0]
}

func (o *Objective) generateFinalOutcome() outcome.SingleAssetExit {
	if o.MyRole != 0 {
		panic("Only Alice should call generateFinalOutcome")
	}
	// Since Alice is responsible for issuing vouchers she always has the largest payment amount
	// This means she can just set her FinalOutcomeFromAlice based on the largest voucher amount she has sent
	finalOutcome := o.initialOutcome().Clone()
	finalOutcome.Allocations[0].Amount.Sub(finalOutcome.Allocations[0].Amount, o.MinimumPaymentAmount)
	finalOutcome.Allocations[1].Amount.Add(finalOutcome.Allocations[1].Amount, o.MinimumPaymentAmount)
	return finalOutcome
}

// finalState returns the final state for the virtual channel
func (o *Objective) generateFinalState() state.State {
	vp := state.VariablePart{Outcome: outcome.Exit{o.generateFinalOutcome()}, TurnNum: FinalTurnNum, IsFinal: true}
	return state.StateFromFixedAndVariablePart(o.V.FixedPart, vp)
}

// Id returns the objective id.
func (o *Objective) Id() protocols.ObjectiveId {
	id := o.VId().String()
	return protocols.ObjectiveId(ObjectivePrefix + id)
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
	for i, peer := range o.V.Participants {
		if i != int(o.MyRole) {
			peers = append(peers, peer)
		}
	}
	messages := protocols.CreateRejectionNoticeMessage(o.Id(), peers...)

	return &updated, protocols.SideEffects{MessagesToSend: messages}
}

// OwnsChannel returns the channel that the objective is funding.
func (o *Objective) OwnsChannel() types.Destination {
	return o.VId()
}

// GetStatus returns the status of the objective.
func (o *Objective) GetStatus() protocols.ObjectiveStatus {
	return o.Status
}

// Related returns channels that need to be stored along with the objective.
func (o *Objective) Related() []protocols.Storable {
	related := []protocols.Storable{}

	related = append(related, o.V)

	if o.ToMyLeft != nil {
		related = append(related, o.ToMyLeft)
	}

	if o.ToMyRight != nil {
		related = append(related, o.ToMyRight)
	}
	return related
}

// Clone returns a deep copy of the receiver.
func (o *Objective) clone() Objective {
	clone := Objective{}
	clone.Status = o.Status

	clone.V = o.V.Clone()

	if o.MinimumPaymentAmount != nil {
		clone.MinimumPaymentAmount = big.NewInt(0).Set(o.MinimumPaymentAmount)
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

// otherParticipants returns the participants in the channel that are not the current participant.
func (o *Objective) otherParticipants() []types.Address {
	others := make([]types.Address, 0)
	for i, p := range o.V.Participants {
		if i != int(o.MyRole) {
			others = append(others, p)
		}
	}
	return others
}

func (o *Objective) hasFinalStateFromAlice() bool {
	ss, ok := o.V.SignedStateForTurnNum[FinalTurnNum]
	return ok && ss.State().IsFinal && !isZero(ss.Signatures()[0])
}

// Crank inspects the extended state and declares a list of Effects to be executed.
func (o *Objective) Crank(secretKey *[]byte) (protocols.Objective, protocols.SideEffects, protocols.WaitingFor, error) {
	updated := o.clone()
	sideEffects := protocols.SideEffects{}

	// Input validation
	if updated.Status != protocols.Approved {
		return &updated, sideEffects, WaitingForNothing, protocols.ErrNotApproved
	}

	// If we don't know the amount yet we send a message to alice to request it
	if !updated.isAlice() && !updated.hasFinalStateFromAlice() {
		alice := o.V.Participants[0]
		messages := protocols.CreateObjectivePayloadMessage(updated.Id(), o.VId(), RequestFinalStatePayload, alice)
		sideEffects.MessagesToSend = append(sideEffects.MessagesToSend, messages...)

		return &updated, sideEffects, WaitingForFinalStateFromAlice, nil
	}

	// Signing of the final state
	if !updated.V.FinalSignedByMe() {
		var s state.State
		if updated.isAlice() {
			s = updated.generateFinalState()
		} else {
			s = updated.finalState()
		}
		// Sign and store:
		ss, err := updated.V.SignAndAddState(s, secretKey)
		if err != nil {
			return &updated, sideEffects, WaitingForNothing, fmt.Errorf("could not sign final state: %w", err)
		}

		if err != nil {
			return &updated, sideEffects, WaitingForNothing, fmt.Errorf("could not get signed final state: %w", err)
		}
		messages := protocols.CreateObjectivePayloadMessage(updated.Id(), ss, SignedStatePayload, o.otherParticipants()...)
		sideEffects.MessagesToSend = append(sideEffects.MessagesToSend, messages...)

	}

	// Check if all participants have signed the final state
	if !updated.V.FinalCompleted() {
		return &updated, sideEffects, WaitingForSupportedFinalState, nil
	}

	if !updated.isAlice() && !updated.leftHasDefunded() {
		ledgerSideEffects, err := updated.updateLedgerToRemoveGuarantee(updated.ToMyLeft, secretKey)
		if err != nil {
			return o, protocols.SideEffects{}, WaitingForNothing, fmt.Errorf("error updating ledger funding: %w", err)
		}

		sideEffects.Merge(ledgerSideEffects)
	}

	if !updated.isBob() && !updated.rightHasDefunded() {
		ledgerSideEffects, err := updated.updateLedgerToRemoveGuarantee(updated.ToMyRight, secretKey)
		if err != nil {
			return o, protocols.SideEffects{}, WaitingForNothing, fmt.Errorf("error updating ledger funding: %w", err)
		}

		sideEffects.Merge(ledgerSideEffects)

	}

	if !updated.leftHasDefunded() {
		return &updated, sideEffects, WaitingForDefundingOnMyLeft, nil
	}

	if !updated.rightHasDefunded() {
		return &updated, sideEffects, WaitingForDefundingOnMyRight, nil
	}

	// Mark the objective as done
	updated.Status = protocols.Completed
	return &updated, sideEffects, WaitingForNothing, nil
}

// isAlice returns true if the receiver represents participant 0 in the virtualdefund protocol.
func (o *Objective) isAlice() bool {
	return o.MyRole == 0
}

// isBob returns true if the receiver represents the last participant in the virtualdefund protocol.
func (o *Objective) isBob() bool {
	return int(o.MyRole) == len(o.V.Participants)-1
}

// ledgerProposal generates a ledger proposal to remove the guarantee for V for ledger
func (o *Objective) ledgerProposal(ledger *consensus_channel.ConsensusChannel) consensus_channel.Proposal {
	left := o.finalState().Outcome[0].Allocations[0].Amount
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

		// Since the proposal queue is constructed with consecutive turn numbers, we can pass it straight in
		// to create a valid message with ordered proposals:
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
	return o.V.ChannelId()
}

// rightHasDefunded returns whether the ledger channel ToMyRight has removed
// its funding for the target channel.
//
// If ToMyRight==nil then we return true.
func (o *Objective) rightHasDefunded() bool {
	if o.ToMyRight == nil {
		return true
	}

	included := o.ToMyRight.IncludesTarget(o.VId())
	return !included
}

// leftHasDefunded returns whether the ledger channel ToMyLeft has removed
// its funding for the target channel.
//
// If ToMyLeft==nil then we return true.
func (o *Objective) leftHasDefunded() bool {
	if o.ToMyLeft == nil {
		return true
	}

	included := o.ToMyLeft.IncludesTarget(o.VId())
	return !included
}

// getSignedStatePayload takes in a serialized signed state payload and returns the deserialized SignedState.
func getSignedStatePayload(b []byte) (state.SignedState, error) {
	ss := state.SignedState{}
	err := json.Unmarshal(b, &ss)
	if err != nil {
		return ss, fmt.Errorf("could not unmarshal signed state: %w", err)
	}
	return ss, nil
}

// getRequestFinalStatePayload takes in a serialized channel id payload and returns the deserialized channel id.
func getRequestFinalStatePayload(b []byte) (types.Destination, error) {
	cId := types.Destination{}
	err := json.Unmarshal(b, &cId)
	if err != nil {
		return cId, fmt.Errorf("could not unmarshal signatures: %w", err)
	}
	return cId, nil
}

// Update receives an protocols.ObjectiveEvent, applies all applicable event data to the VirtualDefundObjective,
// and returns the updated state.
func (o *Objective) Update(op protocols.ObjectivePayload) (protocols.Objective, error) {
	if o.Id() != op.ObjectiveId {
		return o, fmt.Errorf("event and objective Ids do not match: %s and %s respectively", string(op.ObjectiveId), string(o.Id()))
	}

	switch op.Type {
	case SignedStatePayload:
		ss, err := getSignedStatePayload(op.PayloadData)
		if err != nil {
			return &Objective{}, err
		}
		updated := o.clone()
		err = validateFinalOutcome(updated.V.FixedPart, updated.initialOutcome(), ss.State().Outcome[0], o.V.Participants[o.MyRole], updated.MinimumPaymentAmount)
		if err != nil {
			return o, fmt.Errorf("outcome failed validation %w", err)
		}
		ok := updated.V.AddSignedState(ss)

		if !ok {
			return o, fmt.Errorf("could not add signed state %v", ss)
		}

		return &updated, nil
	case RequestFinalStatePayload:
		// Since the objective is already created we don't need to do anything else with the payload
		return o, nil
	default:
		return o, fmt.Errorf("unknown payload type %s", op.Type)
	}
}

// ReceiveProposal receives a signed proposal and returns an updated VirtualDefund objective.
func (o *Objective) ReceiveProposal(sp consensus_channel.SignedProposal) (protocols.ProposalReceiver, error) {
	var toMyLeftId types.Destination
	var toMyRightId types.Destination

	if o.ToMyLeft != nil {
		toMyLeftId = o.ToMyLeft.Id
	}
	if o.ToMyRight != nil {
		toMyRightId = o.ToMyRight.Id
	}

	updated := o.clone()

	if sp.Proposal.Target() == o.VId() {
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
			return o, fmt.Errorf("error incorporating signed proposal %+v into objective: %w", sp, err)
		}
	}
	return &updated, nil
}

// isZero returns true if every byte field on the signature is zero
func isZero(sig state.Signature) bool {
	zeroSig := state.Signature{}
	return sig.Equal(zeroSig)
}

// ObjectiveRequest represents a request to create a new virtual defund objective.
type ObjectiveRequest struct {
	ChannelId        types.Destination
	objectiveStarted chan struct{}
}

// NewObjectiveRequest creates a new ObjectiveRequest.
func NewObjectiveRequest(channelId types.Destination) ObjectiveRequest {
	return ObjectiveRequest{
		ChannelId:        channelId,
		objectiveStarted: make(chan struct{}),
	}
}

// Id returns the objective id for the request.
func (r ObjectiveRequest) Id(types.Address, *big.Int) protocols.ObjectiveId {
	return protocols.ObjectiveId(ObjectivePrefix + r.ChannelId.String())
}

// GetVirtualChannelFromObjectiveId gets the virtual channel id from the objective id.
func GetVirtualChannelFromObjectiveId(id protocols.ObjectiveId) (types.Destination, error) {
	if !strings.HasPrefix(string(id), ObjectivePrefix) {
		return types.Destination{}, fmt.Errorf("id %s does not have prefix %s", id, ObjectivePrefix)
	}
	raw := string(id)[len(ObjectivePrefix):]
	return types.Destination(common.HexToHash(raw)), nil
}

// validateFinalOutcome is a helper function that validates a final outcome from Alice is valid.
func validateFinalOutcome(vFixed state.FixedPart, initialOutcome outcome.SingleAssetExit, finalOutcome outcome.SingleAssetExit, me types.Address, minAmount *big.Int) error {
	// Check the outcome participants are correct
	alice, bob := vFixed.Participants[0], vFixed.Participants[len(vFixed.Participants)-1]
	if initialOutcome.Allocations[0].Destination != types.AddressToDestination(alice) {
		return fmt.Errorf("0th allocation is not to Alice but to %s", initialOutcome.Allocations[0].Destination)
	}
	if initialOutcome.Allocations[1].Destination != types.AddressToDestination(bob) {
		return fmt.Errorf("1st allocation is not to Bob but to %s", initialOutcome.Allocations[0].Destination)
	}

	// Check the amounts are correct
	initialAliceAmount, initialBobAmount := initialOutcome.Allocations[0].Amount, initialOutcome.Allocations[1].Amount
	finalAliceAmount, finalBobAmount := finalOutcome.Allocations[0].Amount, finalOutcome.Allocations[1].Amount
	paidToBob := big.NewInt(0).Sub(finalBobAmount, initialBobAmount)
	paidFromAlice := big.NewInt(0).Sub(initialAliceAmount, finalAliceAmount)
	if paidToBob.Cmp(paidFromAlice) != 0 {
		return fmt.Errorf("final outcome is not balanced: Alice paid %d, Bob received %d", paidFromAlice, paidToBob)
	}

	// if we're Bob we want to make sure the final state Alice sent is equal to or larger than the payment we already have
	if me == bob {
		if paidToBob.Cmp(minAmount) < 0 {
			return fmt.Errorf("payment amount %d is less than the minimum payment amount %d", paidToBob, minAmount)
		}
	}
	return nil
}

// SignalObjectiveStarted is used by the engine to signal the objective has been started.
func (r ObjectiveRequest) SignalObjectiveStarted() {
	close(r.objectiveStarted)
}

// WaitForObjectiveToStart blocks until the objective starts
func (r ObjectiveRequest) WaitForObjectiveToStart() {
	<-r.objectiveStarted
}
