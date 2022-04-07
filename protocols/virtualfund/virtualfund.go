// Package virtualfund implements an off-chain protocol to virtually fund a channel.
package virtualfund // import "github.com/statechannels/go-nitro/virtualfund"

import (
	"bytes"
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

const ObjectivePrefix = "VirtualFund-"

// errors
var ErrNotApproved = errors.New("objective not approved")

// GuaranteeInfo contains the information used to generate the expected guarantees.
type GuaranteeInfo struct {
	Left                 types.Destination
	Right                types.Destination
	LeftAmount           types.Funds
	RightAmount          types.Funds
	GuaranteeDestination types.Destination
}
type Connection struct {
	Channel          *channel.TwoPartyLedger // todo: #420 deprecate in favor of
	ConsensusChannel *consensus_channel.ConsensusChannel
	GuaranteeInfo    GuaranteeInfo
}

// ledgerChannelAffordsExpectedGuarantees returns true if, for the channel inside the connection, and for each asset keying the input variables, the channel can afford the allocation given the funding.
// The decision is made based on the latest supported state of the channel.
//
// Both arguments are maps keyed by the same asset.
func (c *Connection) ledgerChannelAffordsExpectedGuarantees() bool {
	return c.Channel.Affords(c.getExpectedGuarantees(), c.Channel.OnChainFunding) // todo: #420 deprecate
}

// insertGuaranteeInfo mutates the reciever Connection struct.
func (c *Connection) insertGuaranteeInfo(a0 types.Funds, b0 types.Funds, vId types.Destination, left types.Destination, right types.Destination) error {

	c.GuaranteeInfo = GuaranteeInfo{
		Left:                 left,
		Right:                right,
		LeftAmount:           a0,
		RightAmount:          b0,
		GuaranteeDestination: vId,
	}

	// Check that the guarantee metadata can be encoded. This allows us to avoid clunky error-return-chains for getExpectedGuarantees
	metadata := outcome.GuaranteeMetadata{
		Left:  c.GuaranteeInfo.Left,
		Right: c.GuaranteeInfo.Right,
	}
	_, err := metadata.Encode()
	if err != nil {
		return err
	}

	return nil
}

// handleProposal recieves a signed proposal and acts according to the leader / follower
// status of the Connection's ConsensusChannel
func (c *Connection) handleProposal(sp consensus_channel.SignedProposal) error {
	if c == nil {
		return fmt.Errorf("nil connection should not handle proposals")
	}

	if sp.Proposal.ChannelID != c.ConsensusChannel.Id {
		return consensus_channel.ErrIncorrectChannelID
	}

	if c.ConsensusChannel != nil {
		return c.ConsensusChannel.Receive(sp)
	}

	return nil
}

// getExpectedGuarantees returns a map of asset addresses to guarantees for a Connection.
func (c *Connection) getExpectedGuarantees() map[types.Address]outcome.Allocation {
	expectedGuaranteesForLedgerChannel := make(map[types.Address]outcome.Allocation)
	metadata := outcome.GuaranteeMetadata{
		Left:  c.GuaranteeInfo.Left,
		Right: c.GuaranteeInfo.Right,
	}
	encodedGuarantee, err := metadata.Encode()
	// This error is unexpected. insertGuaranteeInfo checks that the guarantee metadata can be encoded.
	// If this panic is triggered, GuranteeInfo has been modified after creation
	if err != nil {
		panic(err)
	}

	channelFunds := c.GuaranteeInfo.LeftAmount.Add(c.GuaranteeInfo.RightAmount)

	for asset, amount := range channelFunds {
		expectedGuaranteesForLedgerChannel[asset] = outcome.Allocation{
			Destination:    c.GuaranteeInfo.GuaranteeDestination,
			Amount:         amount,
			AllocationType: outcome.GuaranteeAllocationType,
			Metadata:       encodedGuarantee,
		}
	}

	return expectedGuaranteesForLedgerChannel
}

// Objective is a cache of data computed by reading from the store. It stores (potentially) infinite data.
type Objective struct {
	Status protocols.ObjectiveStatus
	V      *channel.SingleHopVirtualChannel

	ToMyLeft  *Connection
	ToMyRight *Connection

	n      uint // number of intermediaries
	MyRole uint // index in the virtual funding protocol. 0 for Alice, n for Bob. Otherwise, one of the intermediaries.

	a0 types.Funds // Initial balance for Alice
	b0 types.Funds // Initial balance for Bob

}

// NewObjective creates a new virtual funding objective from a given request.
func NewObjective(request ObjectiveRequest, getTwoPartyLedger GetTwoPartyLedgerFunction, getTwoPartyConsensusLedger GetTwoPartyConsensusLedgerFunction) (Objective, error) {
	right, ok := getTwoPartyLedger(request.MyAddress, request.Intermediary)
	if !ok {
		return Objective{}, fmt.Errorf("could not find ledger for %s and %s", request.MyAddress, request.Intermediary)
	}

	rightCC, ok := getTwoPartyConsensusLedger(request.Intermediary)

	if !ok {
		return Objective{}, fmt.Errorf("could not find ledger for %s and %s", request.MyAddress, request.Intermediary)
	}
	var left *channel.TwoPartyLedger
	var leftCC *consensus_channel.ConsensusChannel

	objective, err := constructFromState(true,
		state.State{
			ChainId:           big.NewInt(0), // TODO
			Participants:      []types.Address{request.MyAddress, request.Intermediary, request.CounterParty},
			ChannelNonce:      big.NewInt(request.Nonce),
			ChallengeDuration: request.ChallengeDuration,
			AppData:           request.AppData,
			Outcome:           request.Outcome,
			TurnNum:           0,
			IsFinal:           false,
		},
		request.MyAddress,
		left, leftCC, right, rightCC) // todo: replace nil consensusChannels
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
	ledgerChannelToMyLeft *channel.TwoPartyLedger,
	consensusChannelToMyLeft *consensus_channel.ConsensusChannel,
	ledgerChannelToMyRight *channel.TwoPartyLedger,
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
	v, err := channel.NewSingleHopVirtualChannel(initialStateOfV, init.MyRole)
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
			return Objective{}, errors.New("Allocation in slot 0 does not correspond to participant 0")
		}
		amount0 := initialStateOfV.Outcome[i].Allocations[0].Amount
		if initialStateOfV.Outcome[i].Allocations[1].Destination != types.AddressToDestination(initialStateOfV.Participants[init.n+1]) {
			return Objective{}, errors.New("Allocation in slot 1 does not correspond to participant " + fmt.Sprint(init.n+1))
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
		init.ToMyLeft.Channel = ledgerChannelToMyLeft
		init.ToMyLeft.ConsensusChannel = consensusChannelToMyLeft
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
		init.ToMyRight.Channel = ledgerChannelToMyRight
		init.ToMyRight.ConsensusChannel = consensusChannelToMyRight
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

// Update receives an protocols.ObjectiveEvent, applies all applicable event data to the VirtualFundObjective,
// and returns the updated state.
func (o Objective) Update(event protocols.ObjectiveEvent) (protocols.Objective, error) {
	if o.Id() != event.ObjectiveId {
		return &o, fmt.Errorf("event and objective Ids do not match: %s and %s respectively", string(event.ObjectiveId), string(o.Id()))
	}

	updated := o.clone()

	var toMyLeftId types.Destination
	var toMyRightId types.Destination

	if !o.isAlice() {
		toMyLeftId = o.ToMyLeft.Channel.Id // Avoid this if it is nil // todo: #420 deprecate
	}
	if !o.isBob() {
		toMyRightId = o.ToMyRight.Channel.Id // Avoid this if it is nil // todo: #420 deprecate
	}

	for _, sp := range event.SignedProposals {
		var err error
		switch sp.Proposal.ChannelID {
		case types.Destination{}:
			return &o, errors.New("signed proposal is not addressed to a ledger channel") // catch this case to avoid unspecified behaviour -- because if Alice or Bob we allow a null channel.
		case toMyLeftId:
			err = updated.ToMyLeft.handleProposal(sp)
		case toMyRightId:
			err = updated.ToMyRight.handleProposal(sp)
		default:
			return &o, fmt.Errorf("signed proposal is not addressed to a known ledger connection")
		}

		if err != nil {
			return &o, fmt.Errorf("error incorporating signed proposal into objective: %w", err)
		}
	}

	for _, ss := range event.SignedStates {
		channelId, _ := ss.State().ChannelId() // TODO handle error
		switch channelId {
		case types.Destination{}:
			return &o, errors.New("null channel id") // catch this case to avoid a panic below -- because if Alice or Bob we allow a null channel.
		case o.V.Id:
			updated.V.AddSignedState(ss)
			// We expect pre and post fund state signatures.
		case toMyLeftId:
			updated.ToMyLeft.Channel.AddSignedState(ss) // todo: #420 depricate
			// We expect a countersigned state including an outcome with expected guarantee. We don't know the exact statehash, though.
		case toMyRightId:
			updated.ToMyRight.Channel.AddSignedState(ss) // todo: #420 depricate
			// We expect a countersigned state including an outcome with expected guarantee. We don't know the exact statehash, though.
		default:
			return &o, errors.New("event channelId out of scope of objective")
		}
	}

	four20StateAssertions(o.ToMyLeft)
	four20StateAssertions(o.ToMyRight)

	return &updated, nil
}

// four20StateAssertions performs some sanity checks against the state of the Channel
// and ConsensusChannel on a connection
// todo 420: DELETE ME
func four20StateAssertions(c *Connection) {
	if c == nil {
		return
	}

	if c.ConsensusChannel == nil {
		return
	}

	supported, err := c.Channel.LatestSupportedState()
	if err != nil {
		panic(err)
	}
	sanityCheck(c.ConsensusChannel.ConsensusVars(), supported)

	latest, err := c.Channel.LatestSignedState()
	if err != nil {
		panic(err)
	}
	sanityCheck(c.ConsensusChannel.LatestProposedVars(), latest.State())
}

// sanityCheck asserts that v and s are equivalent, by checking their TurnNum and Outcome
// todo 420: DELETE ME
func sanityCheck(v consensus_channel.Vars, s state.State) {
	if !s.Outcome.Equal(v.Outcome.AsOutcome()) {
		panic("supported outcome differs")
	}
	if s.TurnNum != v.TurnNum {
		panic("wrong turn number")
	}
}

// Crank inspects the extended state and declares a list of Effects to be executed
// It's like a state machine transition function where the finite / enumerable state is returned (computed from the extended state)
// rather than being independent of the extended state; and where there is only one type of event ("the crank") with no data on it at all.
func (o Objective) Crank(secretKey *[]byte) (protocols.Objective, protocols.SideEffects, protocols.WaitingFor, error) {
	updated := o.clone()

	sideEffects := protocols.SideEffects{}
	// Input validation
	if updated.Status != protocols.Approved {
		return &updated, sideEffects, WaitingForNothing, ErrNotApproved
	}

	// Prefunding

	if !updated.V.PreFundSignedByMe() {
		ss, err := updated.V.SignAndAddPrefund(secretKey)
		if err != nil {
			return &o, protocols.SideEffects{}, WaitingForNothing, err
		}
		messages := protocols.CreateSignedStateMessages(o.Id(), ss, o.V.MyIndex)
		sideEffects.MessagesToSend = append(sideEffects.MessagesToSend, messages...)
	}

	if !updated.V.PreFundComplete() {
		return &updated, sideEffects, WaitingForCompletePrefund, nil
	}

	// Funding

	if !updated.isAlice() && !updated.ToMyLeft.ledgerChannelAffordsExpectedGuarantees() {

		ledgerSideEffects, err := updated.updateLedgerWithGuarantee(*updated.ToMyLeft, secretKey)
		if err != nil {
			return &o, protocols.SideEffects{}, WaitingForNothing, fmt.Errorf("error updating ledger funding: %w", err)
		}
		sideEffects.Merge(ledgerSideEffects)
	}

	if !updated.isBob() && !updated.ToMyRight.ledgerChannelAffordsExpectedGuarantees() {
		ledgerSideEffects, err := updated.updateLedgerWithGuarantee(*updated.ToMyRight, secretKey)
		if err != nil {
			return &o, protocols.SideEffects{}, WaitingForNothing, fmt.Errorf("error updating ledger funding: %w", err)
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
			return &o, protocols.SideEffects{}, WaitingForNothing, err
		}
		messages := protocols.CreateSignedStateMessages(o.Id(), ss, o.V.MyIndex)
		sideEffects.MessagesToSend = append(sideEffects.MessagesToSend, messages...)
	}

	if !updated.V.PostFundComplete() {
		return &updated, sideEffects, WaitingForCompletePostFund, nil
	}

	// Completion
	return &updated, sideEffects, WaitingForNothing, nil
}

func (o Objective) Related() []protocols.Storable {
	// todo: #420 consider usage & signature of this fcn
	ret := []protocols.Storable{&o.V.Channel}

	if o.ToMyLeft != nil {
		ret = append(ret, &o.ToMyLeft.Channel.Channel) // todo: #420 depricate
	}
	if o.ToMyRight != nil {
		ret = append(ret, &o.ToMyRight.Channel.Channel) // todo: #420 depricate
	}

	return ret
}

//////////////////////////////////////////////////
//  Private methods on the VirtualFundObjective //
//////////////////////////////////////////////////

// fundingComplete returns true if the appropriate ledger channel guarantees sufficient funds for J
func (o Objective) fundingComplete() bool {

	// Each peer commits to an update in L_{i-1} and L_i including the guarantees G_{i-1} and {G_i} respectively, and deducting b_0 from L_{I-1} and a_0 from L_i.
	// A = P_0 and B=P_n are special cases. A only does the guarantee for L_0 (deducting a0), and B only foes the guarantee for L_n (deducting b0).

	switch {
	case o.isAlice(): // Alice
		return o.ToMyRight.ledgerChannelAffordsExpectedGuarantees()
	default: // Intermediary
		return o.ToMyRight.ledgerChannelAffordsExpectedGuarantees() && o.ToMyLeft.ledgerChannelAffordsExpectedGuarantees()
	case o.isBob(): // Bob
		return o.ToMyLeft.ledgerChannelAffordsExpectedGuarantees()
	}

}

// Clone returns a deep copy of the receiver.
func (o *Objective) clone() Objective {
	clone := Objective{}
	clone.Status = o.Status
	vClone := o.V.Clone()
	clone.V = vClone

	if o.ToMyLeft != nil {
		// todo: #420 consider cloning for consensusChannels
		lClone := o.ToMyLeft.Channel.Clone()
		clone.ToMyLeft = &Connection{
			Channel:       lClone,
			GuaranteeInfo: o.ToMyLeft.GuaranteeInfo,
		}
	}

	if o.ToMyRight != nil {
		// todo: #420 consider cloning for consensusChannels
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

// isAlice returns true if the reciever represents participant 0 in the virtualfund protocol.
func (o *Objective) isAlice() bool {
	return o.MyRole == 0
}

// isBob returns true if the reciever represents participant n+1 in the virtualfund protocol.
func (o *Objective) isBob() bool {
	return o.MyRole == o.n+1
}

// GetTwoPartyLedgerFunction specifies a function that can be used to retreive ledgers from a store.
type GetTwoPartyLedgerFunction func(firstParty types.Address, secondParty types.Address) (ledger *channel.TwoPartyLedger, ok bool)

// todo: #420 assume name and godoc from GetTwoPartyLedgerFunction
type GetTwoPartyConsensusLedgerFunction func(counterparty types.Address) (ledger *consensus_channel.ConsensusChannel, ok bool)

// ConstructObjectiveFromMessage takes in a message and constructs an objective from it.
// It accepts the message, myAddress, and a function to to retrieve ledgers from a store.
func ConstructObjectiveFromMessage(
	m protocols.Message,
	myAddress types.Address,
	getTwoPartyLedger GetTwoPartyLedgerFunction,
	getTwoPartyConsensusLedger GetTwoPartyConsensusLedgerFunction,
) (Objective, error) {
	if len(m.SignedStates) != 1 {
		return Objective{}, errors.New("expected exactly one signed state in the message")
	}

	initialState := m.SignedStates[0].State()
	participants := initialState.Participants

	// This logic assumes a single hop virtual channel.
	// Currently this is the only type of virtual channel supported.
	alice := participants[0]
	intermediary := participants[1]
	bob := participants[2]

	var left *channel.TwoPartyLedger
	var right *channel.TwoPartyLedger
	var leftC *consensus_channel.ConsensusChannel
	var rightC *consensus_channel.ConsensusChannel
	var ok bool

	if myAddress == alice {
		return Objective{}, errors.New("participant[0] should not construct objectives from peer messages")
	} else if myAddress == bob {
		left, ok = getTwoPartyLedger(intermediary, bob)
		if !ok {
			return Objective{}, fmt.Errorf("could not find a left ledger channel between %v and %v", intermediary, bob)
		}

		leftC, _ = getTwoPartyConsensusLedger(intermediary)
		// if !ok {
		// 	// todo: #420 handle, after consensus_channel lifecycle is defined & channels present
		// 	return Objective{}, fmt.Errorf("could not find a left ledger channel between %v and %v", intermediary, bob)
		// }

	} else if myAddress == intermediary {
		left, ok = getTwoPartyLedger(alice, intermediary)
		if !ok {
			return Objective{}, fmt.Errorf("could not find a left ledger channel between %v and %v", alice, intermediary)
		}

		leftC, _ = getTwoPartyConsensusLedger(alice)
		// if !ok {
		// 	// todo: #420 handle, after consensus_channel lifecycle is defined & channels present
		// 	return Objective{}, fmt.Errorf("could not find a left ledger channel between %v and %v", alice, intermediary)
		// }

		right, ok = getTwoPartyLedger(intermediary, bob)
		if !ok {
			return Objective{}, fmt.Errorf("could not find a right ledger channel between %v and %v", intermediary, bob)
		}

		rightC, _ = getTwoPartyConsensusLedger(bob)
		// if !ok {
		// 	// todo: #420 handle, after consensus_channel lifecycle is defined & channels present
		// 	return Objective{}, fmt.Errorf("could not find a right ledger channel between %v and %v", intermediary, bob)
		// }

	} else {
		return Objective{}, fmt.Errorf("client address not found in an expected participant index")
	}

	return constructFromState(
		true, // TODO ensure objective in only approved if the application has given permission somehow
		initialState,
		myAddress,
		left, leftC, // todo: replace nil consensusChannels
		right, rightC,
	)
}

// IsVirtualFundObjective inspects a objective id and returns true if the objective id is for a virtual fund objective.
func IsVirtualFundObjective(id protocols.ObjectiveId) bool {
	return strings.HasPrefix(string(id), ObjectivePrefix)
}

// proposeLedgerUpdate will propose a ledger update to the channel by crafting a new state
func (o *Objective) proposeLedgerUpdate(connection Connection, sk *[]byte) (protocols.SideEffects, error) {
	ledger := connection.Channel // todo: #420 deprecate - replace with LeaderChannel.Propose workflow
	left := connection.GuaranteeInfo.Left
	right := connection.GuaranteeInfo.Right
	leftAmount := connection.GuaranteeInfo.LeftAmount
	rightAmount := connection.GuaranteeInfo.RightAmount

	if !ledger.IsProposer() {
		return protocols.SideEffects{}, errors.New("only the proposer can propose a ledger update")
	}

	sideEffects := protocols.SideEffects{}

	supported, err := ledger.LatestSupportedState()
	if err != nil {
		return protocols.SideEffects{}, fmt.Errorf("error finding a supported state: %w", err)
	}

	proposed, proposedFound := ledger.Proposed()

	// Clone the state and update the turn to the next turn.
	// We want to use the most recent proposal if it exists, otherwise we use the support.
	var nextState state.State
	if proposedFound {
		nextState = proposed.Clone()
	} else {
		nextState = supported.Clone()
	}
	if nextState.Outcome.Affords(connection.getExpectedGuarantees(), ledger.OnChainFunding) {
		return protocols.SideEffects{}, nil
	}
	nextState.TurnNum = nextState.TurnNum + 1
	// Update the outcome with the guarantee.
	nextState.Outcome, err = nextState.Outcome.DivertToGuarantee(left, right, leftAmount, rightAmount, o.V.Id)
	if err != nil {
		return protocols.SideEffects{}, fmt.Errorf("error updating ledger channel outcome: %w", err)
	}

	// Sign the state and add it to the ledger.
	ss, err := ledger.SignAndAddState(nextState, sk)
	if err != nil {
		return protocols.SideEffects{}, fmt.Errorf("error adding signed state: %w", err)
	}

	// Add a message with the signed state
	messages := protocols.CreateSignedStateMessages(o.Id(), ss, ledger.MyIndex)
	sideEffects.MessagesToSend = append(sideEffects.MessagesToSend, messages...)

	return sideEffects, nil
}

// acceptLedgerUpdate checks for a ledger state proposal and accepts that proposal if it satisfies the expected guarantee.
func (o *Objective) acceptLedgerUpdate(ledgerConnection Connection, sk *[]byte) (protocols.SideEffects, error) {
	ledger := ledgerConnection.Channel // todo: #420 deprecate - replace with FollowerChannel.Receive workflow
	proposed, ok := ledger.Proposed()

	if !ok {
		return protocols.SideEffects{}, fmt.Errorf("no proposed state found for ledger channel %s", ledger.Id)
	}

	supported, err := ledger.LatestSupportedState()
	if err != nil {
		return protocols.SideEffects{}, fmt.Errorf("no supported state found for ledger channel %s: %w", ledger.Id, err)
	}

	// Determine if we are left or right in the guarantee and determine our amounts.
	ourAddress := types.AddressToDestination(ledger.Participants[ledger.MyIndex])
	var ourDeposit types.Funds
	if ledger.MyIndex == 0 {
		ourDeposit = ledgerConnection.GuaranteeInfo.LeftAmount
	} else if ledger.MyIndex == 1 {
		ourDeposit = ledgerConnection.GuaranteeInfo.RightAmount
	}

	ourPreviousTotal := supported.Outcome.TotalAllocatedFor(ourAddress)
	ourNewTotal := proposed.Outcome.TotalAllocatedFor(ourAddress)

	// Our new total should just be our previous total minus our deposit.
	proposedMaintainsOurFunds := ourNewTotal.Add(ourDeposit).Equal(ourPreviousTotal)

	proposedAffordsGuarantee := proposed.Outcome.Affords(ledgerConnection.getExpectedGuarantees(), ledger.OnChainFunding)

	if proposedMaintainsOurFunds && proposedAffordsGuarantee {

		ss, err := ledger.SignAndAddState(proposed, sk)
		if err != nil {
			return protocols.SideEffects{}, fmt.Errorf("error adding signed state: %w", err)
		}
		sideEffects := protocols.SideEffects{}
		messages := protocols.CreateSignedStateMessages(o.Id(), ss, ledger.MyIndex)
		sideEffects.MessagesToSend = append(sideEffects.MessagesToSend, messages...)
		return sideEffects, nil
	}

	return protocols.SideEffects{}, nil

}

// updateLedgerWithGuarantee updates the ledger channel funding to include the guarantee.
// If the user is the proposer a new ledger state will be created and signed.
// If the user is the follower then they will sign a ledger state proposal if it satisfies their expected guarantees.
func (o *Objective) updateLedgerWithGuarantee(ledgerConnection Connection, sk *[]byte) (protocols.SideEffects, error) {

	ledger := ledgerConnection.Channel // todo: #420 deprecate

	if ledger.IsProposer() { // If the user is the proposer craft a new proposal
		sideEffects, err := o.proposeLedgerUpdate(ledgerConnection, sk)
		if err != nil {
			return protocols.SideEffects{}, fmt.Errorf("error proposing ledger update: %w", err)
		} else {
			return sideEffects, nil
		}
	} else if _, ok := ledger.Proposed(); ok { // Otherwise if there is a proposal accept it if it satisfies the guarantee
		sideEffects, err := o.acceptLedgerUpdate(ledgerConnection, sk)
		if err != nil {
			return protocols.SideEffects{}, fmt.Errorf("error proposing ledger update: %w", err)
		} else {
			return sideEffects, nil
		}
	} else {
		return protocols.SideEffects{}, nil
	}

}

// ObjectiveRequest represents a request to create a new virtual funding objective.
type ObjectiveRequest struct {
	MyAddress         types.Address
	Intermediary      types.Address
	CounterParty      types.Address
	AppDefinition     types.Address
	AppData           types.Bytes
	ChallengeDuration *types.Uint256
	Outcome           outcome.Exit
	Nonce             int64
}

// Id returns the objective id for the request.
func (r ObjectiveRequest) Id() protocols.ObjectiveId {
	fixedPart := state.FixedPart{ChainId: big.NewInt(0), // TODO
		Participants:      []types.Address{r.MyAddress, r.Intermediary, r.CounterParty},
		ChannelNonce:      big.NewInt(r.Nonce),
		ChallengeDuration: r.ChallengeDuration}

	channelId, _ := fixedPart.ChannelId()
	return protocols.ObjectiveId(ObjectivePrefix + channelId.String())
}
