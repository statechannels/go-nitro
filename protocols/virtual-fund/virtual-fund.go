package virtualfund

import (
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

var NoSideEffects = protocols.SideEffects{}

// errors
var ErrNotApproved = errors.New("objective not approved")

type Connection struct {
	Channel            channel.TwoPartyLedger
	ExpectedGuarantees map[types.Address]outcome.Allocation
}

// VirtualFundObjective is a cache of data computed by reading from the store. It stores (potentially) infinite data.
type VirtualFundObjective struct {
	Status protocols.ObjectiveStatus
	V      *channel.Channel

	ToMyLeft  *Connection
	ToMyRight *Connection

	n      uint // number of ledger channels (num_hops + 1)
	MyRole uint // index in the virtual funding protocol. 0 for Alice, n+1 for Bob. Otherwise, one of the intermediaries.

	a0 types.Funds // Initial balance for Alice
	b0 types.Funds // Initial balance for Bob

	requestedLedgerUpdates bool // records that the ledger update side effects were previously generated (they may not have been executed yet)
}

////////////////////////////////////////////////
// Public methods on the VirtualFundObjective //
////////////////////////////////////////////////

// New initiates a VirtualFundObjective.
func New(
	initialStateOfV state.State,
	myAddress types.Address,
	n uint, // number of ledger channels (num_hops + 1)
	myRole uint,
	ledgerChannelToMyLeft channel.TwoPartyLedger,
	ledgerChannelToMyRight channel.TwoPartyLedger,
) (VirtualFundObjective, error) {
	// role and ledger-channel checks
	if myRole > n+1 {
		return VirtualFundObjective{}, fmt.Errorf(`invalid role <%d> specified in %d-hop virtual-fund objective`,
			myRole, n-1)
	}

	switch myRole {
	case 0: // Alice
		if !ledgerChannelToMyRight.IsTwoPartyLedger {
			return VirtualFundObjective{}, errors.New(`alice's right-channel is not a two-party ledger channel`)
		}
	case n + 1: // Bob
		if !ledgerChannelToMyLeft.IsTwoPartyLedger {
			return VirtualFundObjective{}, errors.New(`bob's left-channel is not a two-party ledger channel`)
		}
	default: // Intermediary
		if !(ledgerChannelToMyLeft.IsTwoPartyLedger && ledgerChannelToMyRight.IsTwoPartyLedger) {
			return VirtualFundObjective{}, errors.New(`supplied channels are not two party ledger channels`)
		}
	}

	var init VirtualFundObjective

	// Initialize virtual channel
	v, err := channel.New(initialStateOfV, myRole)
	if err != nil {
		return VirtualFundObjective{}, err
	}

	init.V = &v
	init.n = n
	init.MyRole = myRole
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
	if myRole != 0 { // everyone other than Alice has a left-channel
		init.ToMyLeft = &Connection{}
		init.ToMyLeft.Channel = ledgerChannelToMyLeft
		err = init.ToMyLeft.insertExpectedGuarantees(
			init.a0,
			init.b0,
			init.V.Id,
			init.ToMyLeft.Channel.TheirDestination(),
			init.ToMyLeft.Channel.MyDestination())
		if err != nil {
			return VirtualFundObjective{}, err
		}
	}

	if myRole != n+1 { // everyone other than Bob has a right-channel
		init.ToMyRight = &Connection{}
		init.ToMyRight.Channel = ledgerChannelToMyRight
		err = init.ToMyRight.insertExpectedGuarantees(
			init.a0,
			init.b0,
			init.V.Id,
			init.ToMyRight.Channel.MyDestination(),
			init.ToMyRight.Channel.TheirDestination())

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

	updated := s.clone()

	var toMyLeftId types.Destination
	var toMyRightId types.Destination

	if s.MyRole != 0 {
		toMyLeftId = s.ToMyLeft.Channel.Id // Avoid this if it is nil
	}
	if s.MyRole != s.n+1 {
		toMyRightId = s.ToMyRight.Channel.Id // Avoid this if it is nil
	}

	switch event.ChannelId {
	case types.Destination{}:
		return s, errors.New("null channel id") // catch this case to avoid a panic below -- because if Alice or Bob we allow a null channel.
	case s.V.Id:
		updated.V.AddSignedStates(event.Sigs)
		// We expect pre and post fund state signatures.
	case toMyLeftId:
		updated.ToMyLeft.Channel.AddSignedStates(event.Sigs)
		// We expect a countersigned state including an outcome with expected guarantee. We don't know the exact statehash, though.
	case toMyRightId:
		updated.ToMyRight.Channel.AddSignedStates(event.Sigs)
		// We expect a countersigned state including an outcome with expected guarantee. We don't know the exact statehash, though.
	default:
		return s, errors.New("event channelId out of scope of objective")
	}
	return updated, nil

}

// Crank inspects the extended state and declares a list of Effects to be executed
// It's like a state machine transition function where the finite / enumerable state is returned (computed from the extended state)
// rather than being independent of the extended state; and where there is only one type of event ("the crank") with no data on it at all.
func (s VirtualFundObjective) Crank(secretKey *[]byte) (protocols.Objective, protocols.SideEffects, protocols.WaitingFor, error) {
	updated := s.clone()

	// Input validation
	if updated.Status != protocols.Approved {
		return updated, NoSideEffects, WaitingForNothing, ErrNotApproved
	}

	// Prefunding

	if !updated.V.PreFundSignedByMe() {
		sig, err := updated.V.PreFundState().Sign(*secretKey)
		if err != nil {
			return s, NoSideEffects, WaitingForNothing, err
		}
		ok := updated.V.AddSignedState(updated.V.PreFundState(), sig)
		if !ok {
			return s, NoSideEffects, WaitingForNothing, errors.New(`could not add prefund state`)
		}
		return updated, NoSideEffects, WaitingForCompletePrefund, nil

	}

	if !updated.V.PreFundComplete() {
		return updated, NoSideEffects, WaitingForCompletePrefund, nil
	}

	// Funding

	if !updated.requestedLedgerUpdates {
		updated.requestedLedgerUpdates = true
		return updated, s.generateLedgerRequestSideEffects(), WaitingForCompleteFunding, nil
	}

	if !updated.fundingComplete() {
		return updated, NoSideEffects, WaitingForCompleteFunding, nil
	}

	// Postfunding
	if !updated.V.PostFundSignedByMe() {
		sig, err := updated.V.PostFundState().Sign(*secretKey)
		if err != nil {
			return s, NoSideEffects, WaitingForNothing, err
		}
		ok := updated.V.AddSignedState(updated.V.PostFundState(), sig)
		if !ok {
			return s, NoSideEffects, WaitingForNothing, errors.New(`could not add postfund state`)
		}
		return updated, NoSideEffects, WaitingForCompletePostFund, nil

	}

	if !updated.V.PostFundComplete() {
		return updated, NoSideEffects, WaitingForCompletePostFund, nil
	}

	// Completion
	return updated, NoSideEffects, WaitingForNothing, nil
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
	// A = P_0 and B=P_n+1 are special cases. A only does the guarantee for L_0 (deducting a0), and B only foes the guarantee for L_n (deducting b0).

	n := s.n

	switch {
	case s.MyRole == 0: // Alice
		return s.ToMyRight.ledgerChannelAffordsExpectedGuarantees()
	case s.MyRole < n+1: // Intermediary
		return s.ToMyRight.ledgerChannelAffordsExpectedGuarantees() && s.ToMyLeft.ledgerChannelAffordsExpectedGuarantees()
	case s.MyRole == n+1: // Bob
		return s.ToMyLeft.ledgerChannelAffordsExpectedGuarantees()
	default: // Invalid
		return false
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
func (s VirtualFundObjective) generateLedgerRequestSideEffects() protocols.SideEffects {
	sideEffects := protocols.SideEffects{}
	sideEffects.LedgerRequests = make([]protocols.LedgerRequest, 0)
	if s.MyRole > 0 { // Not Alice
		sideEffects.LedgerRequests = append(sideEffects.LedgerRequests,
			protocols.LedgerRequest{
				LedgerId:    s.ToMyLeft.Channel.Id,
				Destination: s.V.Id,
				Amount:      s.V.Total(),
				Left:        s.ToMyLeft.Channel.TheirDestination(),
				Right:       s.ToMyLeft.Channel.MyDestination(),
			})
	}
	n := s.n
	if s.MyRole < n { // Not Bob
		sideEffects.LedgerRequests = append(sideEffects.LedgerRequests,
			protocols.LedgerRequest{
				LedgerId:    s.ToMyRight.Channel.Id,
				Destination: s.V.Id,
				Amount:      s.V.Total(),
				Left:        s.ToMyRight.Channel.MyDestination(),
				Right:       s.ToMyRight.Channel.TheirDestination(),
			})
	}
	return sideEffects
}

// Clone returns a deep copy of the receiver
func (s *VirtualFundObjective) clone() VirtualFundObjective {
	clone := VirtualFundObjective{}
	clone.Status = s.Status
	vClone := s.V.Clone()
	clone.V = &vClone

	if s.ToMyLeft != nil {
		clone.ToMyLeft = &Connection{
			Channel:            s.ToMyLeft.Channel.Clone(),
			ExpectedGuarantees: s.ToMyLeft.ExpectedGuarantees,
		}
	}

	if s.ToMyRight != nil {
		clone.ToMyRight = &Connection{
			Channel:            s.ToMyRight.Channel.Clone(),
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
