// Package virtualfund implements an off-chain protocol to virtually fund a channel.
package virtualfund // import "github.com/statechannels/go-nitro/virtualfund"

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"
	"reflect"

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

// Equal returns true if the Connection pointed to by the supplied pointer is deeply equal to the receiver.
func (c *Connection) Equal(d *Connection) bool {
	if c == nil && d == nil {
		return true
	}
	if !c.Channel.Equal(d.Channel) {
		return false
	}
	if !reflect.DeepEqual(c.ExpectedGuarantees, d.ExpectedGuarantees) {
		return false
	}
	return true

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

	requestedLedgerUpdates bool // records that the ledger update side effects were previously generated (they may not have been executed yet)
}

// NewObjective initiates an Objective.
func NewObjective(
	preApprove bool,
	initialStateOfV state.State,
	myAddress types.Address,
	ledgerChannelToMyLeft *channel.TwoPartyLedger,
	ledgerChannelToMyRight *channel.TwoPartyLedger,
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
			return Objective{}, err
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
			return Objective{}, err
		}
	}

	return init, nil
}

// Id returns the objective id.
func (o Objective) Id() protocols.ObjectiveId {
	return protocols.ObjectiveId("VirtualFund-" + o.V.Id.String())
}

// Approve returns an approved copy of the objective.
func (o Objective) Approve() protocols.Objective {
	updated := o.clone()
	// todo: consider case of s.Status == Rejected
	updated.Status = protocols.Approved
	return updated
}

// Approve returns a rejected copy of the objective.
func (o Objective) Reject() protocols.Objective {
	updated := o.clone()
	updated.Status = protocols.Rejected
	return updated
}

// Update receives an protocols.ObjectiveEvent, applies all applicable event data to the VirtualFundObjective,
// and returns the updated state.
func (o Objective) Update(event protocols.ObjectiveEvent) (protocols.Objective, error) {
	if o.Id() != event.ObjectiveId {
		return o, fmt.Errorf("event and objective Ids do not match: %s and %s respectively", string(event.ObjectiveId), string(o.Id()))
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

	for _, ss := range event.SignedStates {
		channelId, _ := ss.State().ChannelId() // TODO handle error
		switch channelId {
		case types.Destination{}:
			return o, errors.New("null channel id") // catch this case to avoid a panic below -- because if Alice or Bob we allow a null channel.
		case o.V.Id:
			updated.V.AddSignedState(ss)
			// We expect pre and post fund state signatures.
		case toMyLeftId:
			updated.ToMyLeft.Channel.AddSignedState(ss)
			// We expect a countersigned state including an outcome with expected guarantee. We don't know the exact statehash, though.
		case toMyRightId:
			updated.ToMyRight.Channel.AddSignedState(ss)
			// We expect a countersigned state including an outcome with expected guarantee. We don't know the exact statehash, though.
		default:
			return o, errors.New("event channelId out of scope of objective")
		}
	}
	return updated, nil

}

// Crank inspects the extended state and declares a list of Effects to be executed
// It's like a state machine transition function where the finite / enumerable state is returned (computed from the extended state)
// rather than being independent of the extended state; and where there is only one type of event ("the crank") with no data on it at all.
func (o Objective) Crank(secretKey *[]byte) (protocols.Objective, protocols.SideEffects, protocols.WaitingFor, []protocols.GuaranteeRequest, error) {
	updated := o.clone()

	sideEffects := protocols.SideEffects{}
	ledgerRequests := []protocols.GuaranteeRequest{}
	// Input validation
	if updated.Status != protocols.Approved {
		return updated, sideEffects, WaitingForNothing, []protocols.GuaranteeRequest{}, ErrNotApproved
	}

	// Prefunding

	if !updated.V.PreFundSignedByMe() {
		ss, err := updated.V.SignAndAddPrefund(secretKey)
		if err != nil {
			return o, protocols.SideEffects{}, WaitingForNothing, []protocols.GuaranteeRequest{}, err
		}
		messages := protocols.CreateSignedStateMessages(o.Id(), ss, o.V.MyIndex)
		sideEffects.MessagesToSend = append(sideEffects.MessagesToSend, messages...)
	}

	if !updated.V.PreFundComplete() {
		return updated, sideEffects, WaitingForCompletePrefund, ledgerRequests, nil
	}

	// Funding

	if !updated.requestedLedgerUpdates {
		updated.requestedLedgerUpdates = true
		ledgerRequests = append(ledgerRequests, o.generateGuaranteeRequestSideEffects()...)
	}

	if !updated.fundingComplete() {
		return updated, sideEffects, WaitingForCompleteFunding, ledgerRequests, nil
	}

	// Postfunding
	if !updated.V.PostFundSignedByMe() {
		ss, err := updated.V.SignAndAddPostfund(secretKey)
		if err != nil {
			return o, protocols.SideEffects{}, WaitingForNothing, []protocols.GuaranteeRequest{}, err
		}
		messages := protocols.CreateSignedStateMessages(o.Id(), ss, o.V.MyIndex)
		sideEffects.MessagesToSend = append(sideEffects.MessagesToSend, messages...)
	}

	if !updated.V.PostFundComplete() {
		return updated, sideEffects, WaitingForCompletePostFund, ledgerRequests, nil
	}

	// Completion
	return updated, sideEffects, WaitingForNothing, ledgerRequests, nil
}

func (o Objective) Channels() []*channel.Channel {
	ret := make([]*channel.Channel, 0, 3)
	ret = append(ret, &o.V.Channel)
	if !o.isAlice() {
		ret = append(ret, &o.ToMyLeft.Channel.Channel)
	}
	if !o.isBob() {
		ret = append(ret, &o.ToMyRight.Channel.Channel)
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

// ledgerChannelAffordsExpectedGuarantees returns true if, for the channel inside the connection, and for each asset keying the input variables, the channel can afford the allocation given the funding.
// The decision is made based on the latest supported state of the channel.
//
// Both arguments are maps keyed by the same asset.
func (connection *Connection) ledgerChannelAffordsExpectedGuarantees() bool {
	return connection.Channel.Affords(connection.ExpectedGuarantees, connection.Channel.OnChainFunding)
}

// generateGuaranteeRequestSideEffects generates the appropriate side effects, which (when executed and countersigned) will update 1 or 2 ledger channels to guarantee the joint channel.
func (o Objective) generateGuaranteeRequestSideEffects() []protocols.GuaranteeRequest {

	requests := make([]protocols.GuaranteeRequest, 0)

	leftAmount := o.V.LeftAmount()
	rightAmount := o.V.RightAmount()

	if !o.isAlice() {
		requests = append(requests,
			protocols.GuaranteeRequest{
				ObjectiveId: o.Id(),
				LedgerId:    o.ToMyLeft.Channel.Id,
				Destination: o.V.Id,
				Left:        types.AddressToDestination(o.V.Participants[o.MyRole-1]),
				LeftAmount:  leftAmount,
				Right:       types.AddressToDestination(o.V.Participants[o.MyRole]),
				RightAmount: rightAmount,
			})
	}
	if !o.isBob() {
		requests = append(requests,
			protocols.GuaranteeRequest{
				ObjectiveId: o.Id(),
				LedgerId:    o.ToMyRight.Channel.Id,
				Destination: o.V.Id,
				Left:        types.AddressToDestination(o.V.Participants[o.MyRole]),
				LeftAmount:  leftAmount,
				Right:       types.AddressToDestination(o.V.Participants[o.MyRole+1]),
				RightAmount: rightAmount,
			})
	}
	return requests
}

// Equal returns true if the supplied DirectFundObjective is deeply equal to the receiver.
func (o Objective) Equal(r Objective) bool {
	return o.Status == r.Status &&
		o.V.Equal(r.V) &&
		o.ToMyLeft.Equal(r.ToMyLeft) &&
		o.ToMyRight.Equal(r.ToMyRight) &&
		o.n == r.n &&
		o.MyRole == r.MyRole &&
		o.a0.Equal(r.a0) &&
		o.b0.Equal(r.b0) &&
		o.requestedLedgerUpdates == r.requestedLedgerUpdates
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
			Channel:            lClone,
			ExpectedGuarantees: o.ToMyLeft.ExpectedGuarantees,
		}
	}

	if o.ToMyRight != nil {
		rClone := o.ToMyRight.Channel.Clone()
		clone.ToMyRight = &Connection{
			Channel:            rClone,
			ExpectedGuarantees: o.ToMyRight.ExpectedGuarantees,
		}
	}

	clone.n = o.n
	clone.MyRole = o.MyRole

	clone.a0 = o.a0
	clone.b0 = o.b0

	clone.requestedLedgerUpdates = o.requestedLedgerUpdates

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
