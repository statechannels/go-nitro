// Package virtualfund implements an off-chain protocol to virtually fund a channel.
package virtualfund // import "github.com/statechannels/go-nitro/virtualfund"

import (
	"bytes"
	"encoding/json"
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
	WaitingForCompletePrefund  protocols.WaitingFor = "WaitingForCompletePrefund"  // Round 1
	WaitingForCompleteFunding  protocols.WaitingFor = "WaitingForCompleteFunding"  // Round 2
	WaitingForCompletePostFund protocols.WaitingFor = "WaitingForCompletePostFund" // Round 3
	WaitingForNothing          protocols.WaitingFor = "WaitingForNothing"          // Finished
)

const (
	SignedStatePayload protocols.PayloadType = "SignedStatePayload"
)

const ObjectivePrefix = "VirtualFund-"

// GuaranteeInfo contains the information used to generate the expected guarantees.
type GuaranteeInfo struct {
	Left                 types.Destination
	Right                types.Destination
	LeftAmount           types.Funds
	RightAmount          types.Funds
	GuaranteeDestination types.Destination
}
type Connection struct {
	Channel       *consensus_channel.ConsensusChannel
	GuaranteeInfo GuaranteeInfo
}

// insertGuaranteeInfo mutates the receiver Connection struct.
func (c *Connection) insertGuaranteeInfo(a0 types.Funds, b0 types.Funds, vId types.Destination, left types.Destination, right types.Destination) error {
	guaranteeInfo := GuaranteeInfo{
		Left:                 left,
		Right:                right,
		LeftAmount:           a0,
		RightAmount:          b0,
		GuaranteeDestination: vId,
	}

	// Check that the guarantee metadata can be encoded. This allows us to avoid clunky error-return-chains for getExpectedGuarantees
	metadata := outcome.GuaranteeMetadata{
		Left:  guaranteeInfo.Left,
		Right: guaranteeInfo.Right,
	}

	if _, err := metadata.Encode(); err != nil {
		return err
	}

	// the metadata can be encoded, so update the connection's guarantee
	c.GuaranteeInfo = guaranteeInfo
	return nil
}

// handleProposal receives a signed proposal and acts according to the leader / follower
// status of the Connection's ConsensusChannel
func (c *Connection) handleProposal(sp consensus_channel.SignedProposal) error {
	if c == nil {
		return fmt.Errorf("nil connection should not handle proposals")
	}

	if sp.Proposal.LedgerID != c.Channel.Id {
		return consensus_channel.ErrIncorrectChannelID
	}

	if c.Channel != nil {
		err := c.Channel.Receive(sp)
		// Ignore stale or future proposals
		if errors.Is(err, consensus_channel.ErrInvalidTurnNum) {
			return nil
		}
	}

	return nil
}

// IsFundingTheTarget computes whether the ledger channel on the receiver funds the guarantee expected by this connection
func (c *Connection) IsFundingTheTarget() bool {
	g := c.getExpectedGuarantee()
	return c.Channel.Includes(g)
}

// getExpectedGuarantee returns a map of asset addresses to guarantees for a Connection.
func (c *Connection) getExpectedGuarantee() consensus_channel.Guarantee {
	amountFunds := c.GuaranteeInfo.LeftAmount.Add(c.GuaranteeInfo.RightAmount)

	// HACK: GuaranteeInfo stores amounts as types.Funds.
	// We only expect a single asset type, and we want to know how much is to be
	// diverted for that asset type.
	// So, we loop through amountFunds and break after the first asset type ...
	var amount *big.Int
	for _, val := range amountFunds {
		amount = val
		break
	}

	target := c.GuaranteeInfo.GuaranteeDestination
	left := c.GuaranteeInfo.Left
	right := c.GuaranteeInfo.Right

	return consensus_channel.NewGuarantee(amount, target, left, right)
}

// Objective is a cache of data computed by reading from the store. It stores (potentially) infinite data.
type Objective struct {
	Status protocols.ObjectiveStatus
	V      *channel.VirtualChannel

	ToMyLeft  *Connection
	ToMyRight *Connection

	n      uint // number of intermediaries
	MyRole uint // index in the virtual funding protocol. 0 for Alice, n+1 for Bob. Otherwise, one of the intermediaries.

	a0 types.Funds // Initial balance for Alice
	b0 types.Funds // Initial balance for Bob
}

// NewObjective creates a new virtual funding objective from a given request.
func NewObjective(request ObjectiveRequest, preApprove bool, myAddress types.Address, chainId *big.Int, getTwoPartyConsensusLedger GetTwoPartyConsensusLedgerFunction) (Objective, error) {
	var rightCC *consensus_channel.ConsensusChannel
	ok := false

	if len(request.Intermediaries) > 0 {
		rightCC, ok = getTwoPartyConsensusLedger(request.Intermediaries[0])
	} else {
		rightCC, ok = getTwoPartyConsensusLedger(request.CounterParty)
	}

	if !ok {
		return Objective{}, fmt.Errorf("could not find ledger for %s and %s", myAddress, request.Intermediaries[0])
	}
	var leftCC *consensus_channel.ConsensusChannel

	participants := []types.Address{myAddress}
	participants = append(participants, request.Intermediaries...)
	participants = append(participants, request.CounterParty)

	objective, err := constructFromState(preApprove,
		state.State{
			Participants:      participants,
			ChannelNonce:      request.Nonce,
			ChallengeDuration: request.ChallengeDuration,
			Outcome:           request.Outcome,
			TurnNum:           0,
			IsFinal:           false,
		},
		myAddress,
		leftCC, rightCC)
	if err != nil {
		return Objective{}, fmt.Errorf("error creating objective: %w", err)
	}
	return objective, nil
}

// constructFromState initiates an Objective from an initial state and set of ledgers.
func constructFromState(
	preApprove bool,
	initialStateOfV state.State,
	myAddress types.Address,
	consensusChannelToMyLeft *consensus_channel.ConsensusChannel,
	consensusChannelToMyRight *consensus_channel.ConsensusChannel,
) (Objective, error) {
	var init Objective

	if preApprove {
		init.Status = protocols.Approved
	} else {
		init.Status = protocols.Unapproved
	}

	// Infer MyRole
	found := false
	for i, addr := range initialStateOfV.Participants {
		if bytes.Equal(addr[:], myAddress[:]) {
			init.MyRole = uint(i)
			found = true
			continue
		}
	}
	if !found {
		return Objective{}, errors.New("not a participant in V")
	}

	// Initialize virtual channel
	v, err := channel.NewVirtualChannel(initialStateOfV, init.MyRole)
	if err != nil {
		return Objective{}, err
	}

	init.V = v

	init.n = uint(len(initialStateOfV.Participants)) - 2 // NewSingleHopVirtualChannel will error unless there are at least 3 participants

	init.a0 = make(map[types.Address]*big.Int)
	init.b0 = make(map[types.Address]*big.Int)

	// Compute a0 and b0 from the initial state of J
	for i := range initialStateOfV.Outcome {
		asset := initialStateOfV.Outcome[i].Asset
		if initialStateOfV.Outcome[i].Allocations[0].Destination != types.AddressToDestination(initialStateOfV.Participants[0]) {
			return Objective{}, errors.New("allocation in slot 0 does not correspond to participant 0")
		}
		amount0 := initialStateOfV.Outcome[i].Allocations[0].Amount
		if initialStateOfV.Outcome[i].Allocations[1].Destination != types.AddressToDestination(initialStateOfV.Participants[init.n+1]) {
			return Objective{}, errors.New("allocation in slot 1 does not correspond to participant " + fmt.Sprint(init.n+1))
		}
		amount1 := initialStateOfV.Outcome[i].Allocations[1].Amount
		if init.a0[asset] == nil {
			init.a0[asset] = big.NewInt(0)
		}
		if init.b0[asset] == nil {
			init.b0[asset] = big.NewInt(0)
		}
		init.a0[asset].Add(init.a0[asset], amount0)
		init.b0[asset].Add(init.b0[asset], amount1)
	}

	// Setup Ledger Channel Connections and expected guarantees
	if !init.isAlice() { // everyone other than Alice has a left-channel
		init.ToMyLeft = &Connection{}

		if consensusChannelToMyLeft == nil {
			return Objective{}, fmt.Errorf("non-alice virtualfund objective requires non-nil left ledger channel")
		}

		init.ToMyLeft.Channel = consensusChannelToMyLeft
		err = init.ToMyLeft.insertGuaranteeInfo(
			init.a0,
			init.b0,
			init.V.Id,
			types.AddressToDestination(init.V.Participants[init.MyRole-1]),
			types.AddressToDestination(init.V.Participants[init.MyRole]),
		)
		if err != nil {
			return Objective{}, err
		}
	}

	if !init.isBob() { // everyone other than Bob has a right-channel
		init.ToMyRight = &Connection{}

		if consensusChannelToMyRight == nil {
			return Objective{}, fmt.Errorf("non-bob virtualfund objective requires non-nil right ledger channel")
		}

		init.ToMyRight.Channel = consensusChannelToMyRight
		err = init.ToMyRight.insertGuaranteeInfo(
			init.a0,
			init.b0,
			init.V.Id,
			types.AddressToDestination(init.V.Participants[init.MyRole]),
			types.AddressToDestination(init.V.Participants[init.MyRole+1]),
		)
		if err != nil {
			return Objective{}, err
		}
	}

	return init, nil
}

// Id returns the objective id.
func (o *Objective) Id() protocols.ObjectiveId {
	return protocols.ObjectiveId(ObjectivePrefix + o.V.Id.String())
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

	messages := protocols.CreateRejectionNoticeMessage(o.Id(), o.otherParticipants()...)
	sideEffects := protocols.SideEffects{MessagesToSend: messages}
	return &updated, sideEffects
}

// OwnsChannel returns the channel that the objective is funding.
func (o *Objective) OwnsChannel() types.Destination {
	return o.V.Id
}

// GetStatus returns the status of the objective.
func (o *Objective) GetStatus() protocols.ObjectiveStatus {
	return o.Status
}

func (o *Objective) otherParticipants() []types.Address {
	otherParticipants := make([]types.Address, 0)
	for i, p := range o.V.Participants {
		if i != int(o.MyRole) {
			otherParticipants = append(otherParticipants, p)
		}
	}
	return otherParticipants
}

func (o *Objective) getPayload(raw protocols.ObjectivePayload) (*state.SignedState, error) {
	payload := &state.SignedState{}

	err := json.Unmarshal(raw.PayloadData, payload)
	if err != nil {
		return nil, err
	}
	return payload, nil
}

func (o *Objective) ReceiveProposal(sp consensus_channel.SignedProposal) (protocols.ProposalReceiver, error) {
	if pId := protocols.GetProposalObjectiveId(sp.Proposal); o.Id() != pId {
		return o, fmt.Errorf("sp and objective Ids do not match: %s and %s respectively", string(pId), string(o.Id()))
	}

	updated := o.clone()

	var toMyLeftId types.Destination
	var toMyRightId types.Destination

	if !o.isAlice() {
		toMyLeftId = o.ToMyLeft.Channel.Id // Avoid this if it is nil
	}
	if !o.isBob() {
		toMyRightId = o.ToMyRight.Channel.Id // Avoid this if it is nil
	}

	if sp.Proposal.Target() == o.V.Id {
		var err error

		switch sp.Proposal.LedgerID {
		case types.Destination{}:
			return o, fmt.Errorf("signed proposal is for a zero-addressed ledger channel") // catch this case to avoid unspecified behaviour -- because if Alice or Bob we allow a null channel.
		case toMyLeftId:
			err = updated.ToMyLeft.handleProposal(sp)
		case toMyRightId:
			err = updated.ToMyRight.handleProposal(sp)
		default:
			return o, fmt.Errorf("signed proposal is not addressed to a known ledger connection")
		}

		if err != nil {
			return o, fmt.Errorf("error incorporating signed proposal %+v into objective: %w", sp, err)
		}
	}
	return &updated, nil
}

// Update receives an protocols.ObjectiveEvent, applies all applicable event data to the VirtualFundObjective,
// and returns the updated state.
func (o *Objective) Update(raw protocols.ObjectivePayload) (protocols.Objective, error) {
	if o.Id() != raw.ObjectiveId {
		return o, fmt.Errorf("raw and objective Ids do not match: %s and %s respectively", string(raw.ObjectiveId), string(o.Id()))
	}
	payload, err := o.getPayload(raw)
	if err != nil {
		return o, fmt.Errorf("error parsing payload: %w", err)
	}
	updated := o.clone()

	if ss := payload; len(ss.Signatures()) != 0 {
		updated.V.AddSignedState(*ss)
	}

	return &updated, nil
}

// Crank inspects the extended state and declares a list of Effects to be executed
// It's like a state machine transition function where the finite / enumerable state is returned (computed from the extended state)
// rather than being independent of the extended state; and where there is only one type of event ("the crank") with no data on it at all.
func (o *Objective) Crank(secretKey *[]byte) (protocols.Objective, protocols.SideEffects, protocols.WaitingFor, error) {
	updated := o.clone()

	sideEffects := protocols.SideEffects{}
	// Input validation
	if updated.Status != protocols.Approved {
		return &updated, sideEffects, WaitingForNothing, protocols.ErrNotApproved
	}

	// Prefunding

	if !updated.V.PreFundSignedByMe() {
		ss, err := updated.V.SignAndAddPrefund(secretKey)
		if err != nil {
			return o, protocols.SideEffects{}, WaitingForNothing, err
		}

		messages := protocols.CreateObjectivePayloadMessage(o.Id(), ss, SignedStatePayload, o.otherParticipants()...)
		sideEffects.MessagesToSend = append(sideEffects.MessagesToSend, messages...)
	}

	if !updated.V.PreFundComplete() {
		return &updated, sideEffects, WaitingForCompletePrefund, nil
	}

	// Funding

	if !updated.isAlice() && !updated.ToMyLeft.IsFundingTheTarget() {

		ledgerSideEffects, err := updated.updateLedgerWithGuarantee(*updated.ToMyLeft, secretKey)
		if err != nil {
			return o, protocols.SideEffects{}, WaitingForNothing, fmt.Errorf("error updating ledger funding: %w", err)
		}
		sideEffects.Merge(ledgerSideEffects)
	}

	if !updated.isBob() && !updated.ToMyRight.IsFundingTheTarget() {
		ledgerSideEffects, err := updated.updateLedgerWithGuarantee(*updated.ToMyRight, secretKey)
		if err != nil {
			return o, protocols.SideEffects{}, WaitingForNothing, fmt.Errorf("error updating ledger funding: %w", err)
		}
		sideEffects.Merge(ledgerSideEffects)
	}

	if !updated.fundingComplete() {
		return &updated, sideEffects, WaitingForCompleteFunding, nil
	}

	// Postfunding
	if !updated.V.PostFundSignedByMe() {
		ss, err := updated.V.SignAndAddPostfund(secretKey)
		if err != nil {
			return o, protocols.SideEffects{}, WaitingForNothing, err
		}

		messages := protocols.CreateObjectivePayloadMessage(o.Id(), ss, SignedStatePayload, o.otherParticipants()...)
		sideEffects.MessagesToSend = append(sideEffects.MessagesToSend, messages...)
	}

	// Alice and Bob require a complete post fund round to know that vouchers may be enforced on chain.
	// Intermediaries do not require the complete post fund, so we allow them to finish the protocol early.
	// If they need to recover funds, they can force V to close by challenging with the pre fund state.
	// Alice and Bob may counter-challenge with a postfund state plus a redemption state.
	// See ADR-0009.
	if !updated.V.PostFundComplete() && (updated.isAlice() || updated.isBob()) {
		return &updated, sideEffects, WaitingForCompletePostFund, nil
	}

	// Completion
	updated.Status = protocols.Completed
	return &updated, sideEffects, WaitingForNothing, nil
}

func (o *Objective) Related() []protocols.Storable {
	ret := []protocols.Storable{&o.V.Channel}

	if o.ToMyLeft != nil {
		ret = append(ret, o.ToMyLeft.Channel)
	}
	if o.ToMyRight != nil {
		ret = append(ret, o.ToMyRight.Channel)
	}

	return ret
}

//////////////////////////////////////////////////
//  Private methods on the VirtualFundObjective //
//////////////////////////////////////////////////

// fundingComplete returns true if the appropriate ledger channel guarantees sufficient funds for J
func (o *Objective) fundingComplete() bool {
	// Each peer commits to an update in L_{i-1} and L_i including the guarantees G_{i-1} and {G_i} respectively, and deducting b_0 from L_{I-1} and a_0 from L_i.
	// A = P_0 and B=P_n are special cases. A only does the guarantee for L_0 (deducting a0), and B only foes the guarantee for L_n (deducting b0).

	switch {
	case o.isAlice():
		return o.ToMyRight.IsFundingTheTarget()
	case o.isBob():
		return o.ToMyLeft.IsFundingTheTarget()
	default: // Intermediary
		return o.ToMyRight.IsFundingTheTarget() && o.ToMyLeft.IsFundingTheTarget()
	}
}

// Clone returns a deep copy of the receiver.
func (o *Objective) clone() Objective {
	clone := Objective{}
	clone.Status = o.Status
	vClone := o.V.Clone()
	clone.V = vClone

	if o.ToMyLeft != nil {
		lClone := o.ToMyLeft.Channel.Clone()
		clone.ToMyLeft = &Connection{
			Channel:       lClone,
			GuaranteeInfo: o.ToMyLeft.GuaranteeInfo,
		}
	}

	if o.ToMyRight != nil {
		rClone := o.ToMyRight.Channel.Clone()
		clone.ToMyRight = &Connection{
			Channel:       rClone,
			GuaranteeInfo: o.ToMyRight.GuaranteeInfo,
		}
	}

	clone.n = o.n
	clone.MyRole = o.MyRole

	clone.a0 = o.a0
	clone.b0 = o.b0
	return clone
}

// isAlice returns true if the receiver represents participant 0 in the virtualfund protocol.
func (o *Objective) isAlice() bool {
	return o.MyRole == 0
}

// isBob returns true if the receiver represents participant n+1 in the virtualfund protocol.
func (o *Objective) isBob() bool {
	return o.MyRole == o.n+1
}

// GetTwoPartyConsensusLedgerFuncion describes functions which return a ConsensusChannel ledger channel between
// the calling client and the given counterparty, if such a channel exists.
type GetTwoPartyConsensusLedgerFunction func(counterparty types.Address) (ledger *consensus_channel.ConsensusChannel, ok bool)

// ConstructObjectiveFromPayload takes in a message and constructs an objective from it.
// It accepts the message, myAddress, and a function to to retrieve ledgers from a store.
func ConstructObjectiveFromPayload(
	p protocols.ObjectivePayload,
	preapprove bool,
	myAddress types.Address,
	getTwoPartyConsensusLedger GetTwoPartyConsensusLedgerFunction,
) (Objective, error) {
	initialState, err := getSignedStatePayload(p.PayloadData)
	if err != nil {
		return Objective{}, fmt.Errorf("could not get signed state payload: %w", err)
	}

	participants := initialState.State().Participants

	var leftC *consensus_channel.ConsensusChannel
	var rightC *consensus_channel.ConsensusChannel
	var ok bool

	if myAddress == participants[0] {
		// I am Alice
		return Objective{}, errors.New("participant[0] should not construct objectives from peer messages")
	} else if myAddress == participants[len(participants)-1] {

		// I am Bob
		leftOfBob := participants[len(participants)-2]
		leftC, ok = getTwoPartyConsensusLedger(leftOfBob)
		if !ok {
			return Objective{}, fmt.Errorf("could not find a left ledger channel between %v and %v", leftOfBob, myAddress)
		}

	} else {
		intermediaries := participants[1 : len(participants)-1]
		foundMyself := false

		for i, intermediary := range intermediaries {
			if myAddress == intermediary {
				foundMyself = true
				// I am intermediary `i` and participant `p`
				p := i + 1 // participants[p] === intermediaries[i]

				leftOfMe := participants[p-1]
				rightOfMe := participants[p+1]

				leftC, ok = getTwoPartyConsensusLedger(leftOfMe)
				if !ok {
					return Objective{}, fmt.Errorf("could not find a left ledger channel between %v and %v", leftOfMe, myAddress)
				}

				rightC, ok = getTwoPartyConsensusLedger(rightOfMe)
				if !ok {
					return Objective{}, fmt.Errorf("could not find a right ledger channel between %v and %v", myAddress, rightOfMe)
				}

				break
			}
		}

		if !foundMyself {
			return Objective{}, fmt.Errorf("client address not found in the participant list")
		}
	}

	return constructFromState(
		preapprove,
		initialState.State(),
		myAddress,
		leftC,
		rightC,
	)
}

// IsVirtualFundObjective inspects a objective id and returns true if the objective id is for a virtual fund objective.
func IsVirtualFundObjective(id protocols.ObjectiveId) bool {
	return strings.HasPrefix(string(id), ObjectivePrefix)
}

func (c *Connection) expectedProposal() consensus_channel.Proposal {
	g := c.getExpectedGuarantee()

	var leftAmount *big.Int
	for _, val := range c.GuaranteeInfo.LeftAmount {
		leftAmount = val
		break
	}
	proposal := consensus_channel.NewAddProposal(c.Channel.Id, g, leftAmount)

	return proposal
}

// proposeLedgerUpdate will propose a ledger update to the channel by crafting a new state
func (o *Objective) proposeLedgerUpdate(connection Connection, sk *[]byte) (protocols.SideEffects, error) {
	ledger := connection.Channel

	if !ledger.IsLeader() {
		return protocols.SideEffects{}, errors.New("only the leader can propose a ledger update")
	}

	sideEffects := protocols.SideEffects{}

	_, err := ledger.Propose(connection.expectedProposal(), *sk)
	if err != nil {
		return protocols.SideEffects{}, err
	}

	recipient := ledger.Follower()

	// Since the proposal queue is constructed with consecutive turn numbers, we can pass it straight in
	// to create a valid message with ordered proposals:
	message := protocols.CreateSignedProposalMessage(recipient, connection.Channel.ProposalQueue()...)

	sideEffects.MessagesToSend = append(sideEffects.MessagesToSend, message)

	return sideEffects, nil
}

// acceptLedgerUpdate checks for a ledger state proposal and accepts that proposal if it satisfies the expected guarantee.
func (o *Objective) acceptLedgerUpdate(c Connection, sk *[]byte) (protocols.SideEffects, error) {
	ledger := c.Channel
	sp, err := ledger.SignNextProposal(c.expectedProposal(), *sk)
	if err != nil {
		return protocols.SideEffects{}, fmt.Errorf("no proposed state found for ledger channel %w", err)
	}
	sideEffects := protocols.SideEffects{}

	// ledger sideEffect
	if proposals := ledger.ProposalQueue(); len(proposals) != 0 {
		sideEffects.ProposalsToProcess = append(sideEffects.ProposalsToProcess, proposals[0].Proposal)
	}

	// message sideEffect
	recipient := ledger.Leader()
	message := protocols.CreateSignedProposalMessage(recipient, sp)
	sideEffects.MessagesToSend = append(sideEffects.MessagesToSend, message)
	return sideEffects, nil
}

// updateLedgerWithGuarantee updates the ledger channel funding to include the guarantee.
// If the user is the proposer a new ledger state will be created and signed.
// If the user is the follower then they will sign a ledger state proposal if it satisfies their expected guarantees.
func (o *Objective) updateLedgerWithGuarantee(ledgerConnection Connection, sk *[]byte) (protocols.SideEffects, error) {
	ledger := ledgerConnection.Channel

	var sideEffects protocols.SideEffects
	g := ledgerConnection.getExpectedGuarantee()
	proposed, err := ledger.IsProposed(g)
	if err != nil {
		return protocols.SideEffects{}, err
	}

	if ledger.IsLeader() { // If the user is the proposer craft a new proposal
		if proposed {
			return protocols.SideEffects{}, nil
		}
		se, err := o.proposeLedgerUpdate(ledgerConnection, sk)
		if err != nil {
			return protocols.SideEffects{}, fmt.Errorf("error proposing ledger update: %w", err)
		}
		sideEffects = se
	} else {
		if err != nil {
			return protocols.SideEffects{}, err
		}
		// If the proposal is next in the queue we accept it
		proposedNext, _ := ledger.IsProposedNext(g)
		if proposedNext {

			se, err := o.acceptLedgerUpdate(ledgerConnection, sk)
			if err != nil {
				return protocols.SideEffects{}, fmt.Errorf("error proposing ledger update: %w", err)
			}

			sideEffects = se
		}
	}

	return sideEffects, nil
}

// ObjectiveRequest represents a request to create a new virtual funding objective.
type ObjectiveRequest struct {
	Intermediaries    []types.Address
	CounterParty      types.Address
	ChallengeDuration uint32
	Outcome           outcome.Exit
	Nonce             uint64
	AppDefinition     types.Address
	objectiveStarted  chan struct{}
}

// NewObjectiveRequest creates a new ObjectiveRequest.
func NewObjectiveRequest(intermediaries []types.Address,
	counterparty types.Address,
	challengeDuration uint32,
	outcome outcome.Exit,
	nonce uint64,
	appDefinition types.Address,
) ObjectiveRequest {
	return ObjectiveRequest{
		Intermediaries:    intermediaries,
		CounterParty:      counterparty,
		ChallengeDuration: challengeDuration,
		Outcome:           outcome,
		Nonce:             nonce,
		AppDefinition:     appDefinition,
		objectiveStarted:  make(chan struct{}),
	}
}

// Id returns the objective id for the request.
func (r ObjectiveRequest) Id(myAddress types.Address, chainId *big.Int) protocols.ObjectiveId {
	idStr := r.channelID(myAddress).String()
	return protocols.ObjectiveId(ObjectivePrefix + idStr)
}

// SignalObjectiveStarted is used by the engine to signal the objective has been started.
func (r ObjectiveRequest) SignalObjectiveStarted() {
	close(r.objectiveStarted)
}

// WaitForObjectiveToStart blocks until the objective starts
func (r ObjectiveRequest) WaitForObjectiveToStart() {
	<-r.objectiveStarted
}

// ObjectiveResponse is the type returned across the API in response to the ObjectiveRequest.
type ObjectiveResponse struct {
	Id        protocols.ObjectiveId
	ChannelId types.Destination
}

// Response computes and returns the appropriate response from the request.
func (r ObjectiveRequest) Response(myAddress types.Address) ObjectiveResponse {
	channelId := r.channelID(myAddress)

	return ObjectiveResponse{
		Id:        protocols.ObjectiveId(ObjectivePrefix + channelId.String()),
		ChannelId: channelId,
	}
}

func (r ObjectiveRequest) channelID(myAddress types.Address) types.Destination {
	participants := []types.Address{myAddress}
	participants = append(participants, r.Intermediaries...)
	participants = append(participants, r.CounterParty)

	fixedPart := state.FixedPart{
		Participants:      participants,
		ChannelNonce:      r.Nonce,
		ChallengeDuration: r.ChallengeDuration,
	}

	return fixedPart.ChannelId()
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
