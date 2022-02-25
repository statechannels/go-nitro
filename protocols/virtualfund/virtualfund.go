// Package virtualfund implements an off-chain protocol to virtually fund a channel.
package virtualfund // import "github.com/statechannels/go-nitro/virtualfund"

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"

	"github.com/statechannels/go-nitro/channel"
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

// errors
var ErrNotApproved = errors.New("objective not approved")

type Connection struct {
	Channel            *channel.TwoPartyLedger
	ExpectedGuarantees map[types.Address]outcome.Allocation
}

// VirtualFundObjective is a cache of data computed by reading from the store. It stores (potentially) infinite data.
type VirtualFundObjective struct {
	Status protocols.ObjectiveStatus
	V      *channel.SingleHopVirtualChannel

	ToMyLeft  *Connection
	ToMyRight *Connection

	n      uint // number of intermediaries
	MyRole uint // index in the virtual funding protocol. 0 for Alice, n for Bob. Otherwise, one of the intermediaries.

	a0 types.Funds // Initial balance for Alice
	b0 types.Funds // Initial balance for Bob

	requestedLedgerUpdates bool // records that the ledger update side effects were previously generated (they may not have been executed yet)
}

////////////////////////////////////////////////
// Public methods on the VirtualFundObjective //
////////////////////////////////////////////////

// New initiates a VirtualFundObjective.
func New(
	preApprove bool,
	initialStateOfV state.State,
	myAddress types.Address,
	ledgerChannelToMyLeft *channel.TwoPartyLedger,
	ledgerChannelToMyRight *channel.TwoPartyLedger,
) (VirtualFundObjective, error) {

	var init VirtualFundObjective

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
		return VirtualFundObjective{}, errors.New("not a participant in V")
	}

	// Initialize virtual channel
	v, err := channel.NewSingleHopVirtualChannel(initialStateOfV, init.MyRole)
	if err != nil {
		return VirtualFundObjective{}, err
	}

	init.V = v

	init.n = uint(len(initialStateOfV.Participants)) - 2 // NewSingleHopVirtualChannel will error unless there are at least 3 participants

	init.a0 = make(map[types.Address]*big.Int)
	init.b0 = make(map[types.Address]*big.Int)

	// Compute a0 and b0 from the initial state of J
	for i := range initialStateOfV.Outcome {
		asset := initialStateOfV.Outcome[i].Asset
		amount0 := initialStateOfV.Outcome[i].Allocations[0].Amount
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
		err = init.ToMyLeft.insertExpectedGuarantees(
			init.a0,
			init.b0,
			init.V.Id,
			types.AddressToDestination(init.V.Participants[init.MyRole-1]),
			types.AddressToDestination(init.V.Participants[init.MyRole]),
		)
		if err != nil {
			return VirtualFundObjective{}, err
		}
	}

	if !init.isBob() { // everyone other than Bob has a right-channel
		init.ToMyRight = &Connection{}
		init.ToMyRight.Channel = ledgerChannelToMyRight
		err = init.ToMyRight.insertExpectedGuarantees(
			init.a0,
			init.b0,
			init.V.Id,
			types.AddressToDestination(init.V.Participants[init.MyRole]),
			types.AddressToDestination(init.V.Participants[init.MyRole+1]),
		)
		if err != nil {
			return VirtualFundObjective{}, err
		}
	}

	return init, nil
}

// Id returns the objective id.
func (s VirtualFundObjective) Id() protocols.ObjectiveId {
	return protocols.ObjectiveId("VirtualFund-" + s.V.Id.String())
}

// Approve returns an approved copy of the objective.
func (s VirtualFundObjective) Approve() protocols.Objective {
	updated := s.clone()
	// todo: consider case of s.Status == Rejected
	updated.Status = protocols.Approved
	return updated
}

// Approve returns a rejected copy of the objective.
func (s VirtualFundObjective) Reject() protocols.Objective {
	updated := s.clone()
	updated.Status = protocols.Rejected
	return updated
}

// Update receives an protocols.ObjectiveEvent, applies all applicable event data to the VirtualFundObjective,
// and returns the updated state.
func (s VirtualFundObjective) Update(event protocols.ObjectiveEvent) (protocols.Objective, error) {
	if s.Id() != event.ObjectiveId {
		return s, fmt.Errorf("event and objective Ids do not match: %s and %s respectively", string(event.ObjectiveId), string(s.Id()))
	}

	updated := s.clone()

	var toMyLeftId types.Destination
	var toMyRightId types.Destination

	if !s.isAlice() {
		toMyLeftId = s.ToMyLeft.Channel.Id // Avoid this if it is nil
	}
	if !s.isBob() {
		toMyRightId = s.ToMyRight.Channel.Id // Avoid this if it is nil
	}

	for _, ss := range event.SignedStates {
		channelId, _ := ss.State().ChannelId() // TODO handle error
		switch channelId {
		case types.Destination{}:
			return s, errors.New("null channel id") // catch this case to avoid a panic below -- because if Alice or Bob we allow a null channel.
		case s.V.Id:
			updated.V.AddSignedState(ss)
			// We expect pre and post fund state signatures.
		case toMyLeftId:
			updated.ToMyLeft.Channel.AddSignedState(ss)
			// We expect a countersigned state including an outcome with expected guarantee. We don't know the exact statehash, though.
		case toMyRightId:
			updated.ToMyRight.Channel.AddSignedState(ss)
			// We expect a countersigned state including an outcome with expected guarantee. We don't know the exact statehash, though.
		default:
			return s, errors.New("event channelId out of scope of objective")
		}
	}
	return updated, nil

}

// Crank inspects the extended state and declares a list of Effects to be executed
// It's like a state machine transition function where the finite / enumerable state is returned (computed from the extended state)
// rather than being independent of the extended state; and where there is only one type of event ("the crank") with no data on it at all.
func (s VirtualFundObjective) Crank(secretKey *[]byte) (protocols.Objective, protocols.SideEffects, protocols.WaitingFor, []protocols.LedgerRequest, error) {
	updated := s.clone()

	sideEffects := protocols.SideEffects{}
	ledgerRequests := []protocols.LedgerRequest{}
	// Input validation
	if updated.Status != protocols.Approved {
		return updated, sideEffects, WaitingForNothing, []protocols.LedgerRequest{}, ErrNotApproved
	}

	// Prefunding

	if !updated.V.PreFundSignedByMe() {
		ss, err := updated.V.SignAndAddPrefund(secretKey)
		if err != nil {
			return s, protocols.SideEffects{}, WaitingForNothing, []protocols.LedgerRequest{}, err
		}
		messages := protocols.CreateSignedStateMessages(s.Id(), ss, s.V.MyIndex)
		sideEffects.MessagesToSend = append(sideEffects.MessagesToSend, messages...)
	}

	if !updated.V.PreFundComplete() {
		return updated, sideEffects, WaitingForCompletePrefund, ledgerRequests, nil
	}

	// Funding

	if !updated.requestedLedgerUpdates {
		updated.requestedLedgerUpdates = true
		ledgerRequests = append(ledgerRequests, s.generateLedgerRequestSideEffects()...)
	}

	if !updated.fundingComplete() {
		return updated, sideEffects, WaitingForCompleteFunding, ledgerRequests, nil
	}

	// Postfunding
	if !updated.V.PostFundSignedByMe() {
		ss, err := updated.V.SignAndAddPostfund(secretKey)
		if err != nil {
			return s, protocols.SideEffects{}, WaitingForNothing, []protocols.LedgerRequest{}, err
		}
		messages := protocols.CreateSignedStateMessages(s.Id(), ss, s.V.MyIndex)
		sideEffects.MessagesToSend = append(sideEffects.MessagesToSend, messages...)
	}

	if !updated.V.PostFundComplete() {
		return updated, sideEffects, WaitingForCompletePostFund, ledgerRequests, nil
	}

	// Completion
	return updated, sideEffects, WaitingForNothing, ledgerRequests, nil
}

func (s VirtualFundObjective) Channels() []*channel.Channel {
	ret := make([]*channel.Channel, 0, 3)
	ret = append(ret, &s.V.Channel)
	if !s.isAlice() {
		ret = append(ret, &s.ToMyLeft.Channel.Channel)
	}
	if !s.isBob() {
		ret = append(ret, &s.ToMyRight.Channel.Channel)
	}
	return ret
}

//////////////////////////////////////////////////
//  Private methods on the VirtualFundObjective //
//////////////////////////////////////////////////

// insertExpectedGuaranteesForLedgerChannel mutates the reciever Connection struct.
func (connection *Connection) insertExpectedGuarantees(a0 types.Funds, b0 types.Funds, vId types.Destination, left types.Destination, right types.Destination) error {
	expectedGuaranteesForLedgerChannel := make(map[types.Address]outcome.Allocation)
	metadata := outcome.GuaranteeMetadata{
		Left:  left,
		Right: right,
	}
	encodedGuarantee, err := metadata.Encode()
	if err != nil {
		return err
	}

	channelFunds := a0.Add(b0)

	for asset, amount := range channelFunds {
		expectedGuaranteesForLedgerChannel[asset] = outcome.Allocation{
			Destination:    vId,
			Amount:         amount,
			AllocationType: outcome.GuaranteeAllocationType,
			Metadata:       encodedGuarantee,
		}
	}
	connection.ExpectedGuarantees = expectedGuaranteesForLedgerChannel
	return nil
}

// fundingComplete returns true if the appropriate ledger channel guarantees sufficient funds for J
func (s VirtualFundObjective) fundingComplete() bool {

	// Each peer commits to an update in L_{i-1} and L_i including the guarantees G_{i-1} and {G_i} respectively, and deducting b_0 from L_{I-1} and a_0 from L_i.
	// A = P_0 and B=P_n are special cases. A only does the guarantee for L_0 (deducting a0), and B only foes the guarantee for L_n (deducting b0).

	switch {
	case s.isAlice(): // Alice
		return s.ToMyRight.ledgerChannelAffordsExpectedGuarantees()
	default: // Intermediary
		return s.ToMyRight.ledgerChannelAffordsExpectedGuarantees() && s.ToMyLeft.ledgerChannelAffordsExpectedGuarantees()
	case s.isBob(): // Bob
		return s.ToMyLeft.ledgerChannelAffordsExpectedGuarantees()
	}

}

// ledgerChannelAffordsExpectedGuarantees returns true if, for the channel inside the connection, and for each asset keying the input variables, the channel can afford the allocation given the funding.
// The decision is made based on the latest supported state of the channel.
//
// Both arguments are maps keyed by the same asset.
func (connection *Connection) ledgerChannelAffordsExpectedGuarantees() bool {
	return connection.Channel.Affords(connection.ExpectedGuarantees, connection.Channel.OnChainFunding)
}

// generateLedgerRequestSideEffects generates the appropriate side effects, which (when executed and countersigned) will update 1 or 2 ledger channels to guarantee the joint channel.
func (s VirtualFundObjective) generateLedgerRequestSideEffects() []protocols.LedgerRequest {

	requests := make([]protocols.LedgerRequest, 0)

	leftAmount := s.V.LeftAmount()
	rightAmount := s.V.RightAmount()

	if !s.isAlice() {
		requests = append(requests,
			protocols.LedgerRequest{
				ObjectiveId: s.Id(),
				LedgerId:    s.ToMyLeft.Channel.Id,
				Destination: s.V.Id,
				Left:        types.AddressToDestination(s.V.Participants[s.MyRole-1]),
				LeftAmount:  leftAmount,
				Right:       types.AddressToDestination(s.V.Participants[s.MyRole]),
				RightAmount: rightAmount,
			})
	}
	if !s.isBob() {
		requests = append(requests,
			protocols.LedgerRequest{
				ObjectiveId: s.Id(),
				LedgerId:    s.ToMyRight.Channel.Id,
				Destination: s.V.Id,
				Left:        types.AddressToDestination(s.V.Participants[s.MyRole]),
				LeftAmount:  leftAmount,
				Right:       types.AddressToDestination(s.V.Participants[s.MyRole+1]),
				RightAmount: rightAmount,
			})
	}
	return requests
}

// Clone returns a deep copy of the receiver
func (s *VirtualFundObjective) clone() VirtualFundObjective {
	clone := VirtualFundObjective{}
	clone.Status = s.Status
	vClone := s.V.Clone()
	clone.V = vClone

	if s.ToMyLeft != nil {
		lClone := s.ToMyLeft.Channel.Clone()
		clone.ToMyLeft = &Connection{
			Channel:            lClone,
			ExpectedGuarantees: s.ToMyLeft.ExpectedGuarantees,
		}
	}

	if s.ToMyRight != nil {
		rClone := s.ToMyRight.Channel.Clone()
		clone.ToMyRight = &Connection{
			Channel:            rClone,
			ExpectedGuarantees: s.ToMyRight.ExpectedGuarantees,
		}
	}

	clone.n = s.n
	clone.MyRole = s.MyRole

	clone.a0 = s.a0
	clone.b0 = s.b0

	clone.requestedLedgerUpdates = s.requestedLedgerUpdates

	return clone
}

// isAlice returns true if the reciever represents participant 0 in the virtualfund protocol.
func (s *VirtualFundObjective) isAlice() bool {
	return s.MyRole == 0
}

// isBob returns true if the reciever represents participant n+1 in the virtualfund protocol.
func (s *VirtualFundObjective) isBob() bool {
	return s.MyRole == s.n+1
}
