// Package virtualfund implements an off-chain protocol to virtually fund a channel.
package virtualfund // import "github.com/statechannels/go-nitro/virtualfund"

import (
	"encoding/json"
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

// jsonConnection is a serialization-friendly struct representation
// of a Connection
type jsonConnection struct {
	Channel            types.Destination
	ExpectedGuarantees []assetGuarantee
}

// assetGuarantee is a serialization-friendly representation of
// map[asset]Allocation
type assetGuarantee struct {
	Asset     types.Address
	Guarantee outcome.Allocation
}

func (c Connection) MarshalJSON() ([]byte, error) {
	guarantees := []assetGuarantee{}
	for asset, guarantee := range c.ExpectedGuarantees {
		guarantees = append(guarantees, assetGuarantee{
			asset,
			guarantee,
		})
	}
	jsonC := jsonConnection{c.Channel.Id, guarantees}
	bytes, err := json.Marshal(jsonC)

	if err != nil {
		return []byte{}, err
	}

	return bytes, err
}

func (c *Connection) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	var jsonC jsonConnection
	err := json.Unmarshal(data, &jsonC)

	if err != nil {
		return err
	}

	c.Channel = channel.TwoPartyLedger{}
	c.Channel.Id = jsonC.Channel

	c.ExpectedGuarantees = make(map[types.Address]outcome.Allocation)

	for _, eg := range jsonC.ExpectedGuarantees {
		c.ExpectedGuarantees[eg.Asset] = eg.Guarantee
	}

	return nil
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

// jsonVirtualFundObjective is a replaces the VirtualFundObjective's channel pointers
// with the channel's Id, making jsonVirtualFundObjective suitable for serialization
type jsonVirtualFundObjective struct {
	Status protocols.ObjectiveStatus
	V      types.Destination

	ToMyLeft  []byte
	ToMyRight []byte

	N      uint
	MyRole uint

	A0 types.Funds
	B0 types.Funds

	RequestedLedgerUpdates bool
}

////////////////////////////////////////////////
// Public methods on the VirtualFundObjective //
////////////////////////////////////////////////

// New initiates a VirtualFundObjective.
func New(
	preApprove bool,
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

	var init VirtualFundObjective

	if preApprove {
		init.Status = protocols.Approved
	} else {
		init.Status = protocols.Unapproved
	}

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

// MarshalJSON returns a JSON representation of the VirtualFundObjective
//
// NOTE: Marshal -> Unmarshal is a lossy process. All channel data from
//       the virtual and ledger channels (other than Ids) is discarded
func (s VirtualFundObjective) MarshalJSON() ([]byte, error) {
	var left []byte
	var right []byte
	var err error

	if s.ToMyLeft == nil {
		left = []byte("null")
	} else {
		left, err = s.ToMyLeft.MarshalJSON()

		if err != nil {
			return nil, fmt.Errorf("error marshaling left channel of %v: %w", s, err)
		}
	}

	if s.ToMyRight == nil {
		right = []byte("null")
	} else {
		right, err = s.ToMyRight.MarshalJSON()

		if err != nil {
			return nil, fmt.Errorf("error marshaling right channel of %v: %w", s, err)
		}
	}

	jsonVFO := jsonVirtualFundObjective{
		s.Status,
		s.V.Id,
		left,
		right,
		s.n,
		s.MyRole,
		s.a0,
		s.b0,
		s.requestedLedgerUpdates,
	}
	return json.Marshal(jsonVFO)
}

// UnmarshalJSON populates the calling VirtualFundObjective with the
// json-encoded data
//
// NOTE: Marshal -> Unmarshal is a lossy process. All channel data from
//       the virtual and ledger channels (other than Ids) is discarded
func (vfo *VirtualFundObjective) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	var jsonVFO jsonVirtualFundObjective
	if err := json.Unmarshal(data, &jsonVFO); err != nil {
		return fmt.Errorf("failed to unmarshal the VirtualFundObjective: %w", err)
	}

	vfo.V = &channel.Channel{}
	vfo.V.Id = jsonVFO.V

	vfo.ToMyLeft = &Connection{}
	vfo.ToMyRight = &Connection{}
	if err := vfo.ToMyLeft.UnmarshalJSON(jsonVFO.ToMyLeft); err != nil {
		return fmt.Errorf("failed to unmarshal left ledger channel: %w", err)
	}
	if err := vfo.ToMyRight.UnmarshalJSON(jsonVFO.ToMyRight); err != nil {
		return fmt.Errorf("failed to unmarshal right ledger channel: %w", err)
	}

	vfo.Status = jsonVFO.Status
	vfo.n = jsonVFO.N
	vfo.MyRole = jsonVFO.MyRole
	vfo.a0 = jsonVFO.A0
	vfo.b0 = jsonVFO.B0
	vfo.requestedLedgerUpdates = jsonVFO.RequestedLedgerUpdates

	return nil
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
	return &updated
}

// Approve returns a rejected copy of the objective.
func (s VirtualFundObjective) Reject() protocols.Objective {
	updated := s.clone()
	updated.Status = protocols.Rejected
	return &updated
}

// Update receives an protocols.ObjectiveEvent, applies all applicable event data to the VirtualFundObjective,
// and returns the updated state.
func (s VirtualFundObjective) Update(event protocols.ObjectiveEvent) (protocols.Objective, error) {
	if s.Id() != event.ObjectiveId {
		return &s, fmt.Errorf("event and objective Ids do not match: %s and %s respectively", string(event.ObjectiveId), string(s.Id()))
	}

	updated := s.clone()

	var toMyLeftId types.Destination
	var toMyRightId types.Destination

	if s.MyRole != 0 {
		toMyLeftId = s.ToMyLeft.Channel.Id // Avoid this if it is nil
	}
	if s.MyRole != s.n+1 {
		toMyRightId = s.ToMyRight.Channel.Id // Avoid this if it is nil
	}

	for _, ss := range event.SignedStates {
		channelId, _ := ss.State().ChannelId() // TODO handle error
		switch channelId {
		case types.Destination{}:
			return &s, errors.New("null channel id") // catch this case to avoid a panic below -- because if Alice or Bob we allow a null channel.
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
			return &s, errors.New("event channelId out of scope of objective")
		}
	}
	return &updated, nil

}

// Crank inspects the extended state and declares a list of Effects to be executed
// It's like a state machine transition function where the finite / enumerable state is returned (computed from the extended state)
// rather than being independent of the extended state; and where there is only one type of event ("the crank") with no data on it at all.
func (s VirtualFundObjective) Crank(secretKey *[]byte) (protocols.Objective, protocols.SideEffects, protocols.WaitingFor, error) {
	updated := s.clone()

	// Input validation
	if updated.Status != protocols.Approved {
		return &updated, NoSideEffects, WaitingForNothing, ErrNotApproved
	}

	// Prefunding

	if !updated.V.PreFundSignedByMe() {
		sig, err := updated.V.PreFundState().Sign(*secretKey)
		if err != nil {
			return &s, NoSideEffects, WaitingForNothing, err
		}
		ss := state.NewSignedState(updated.V.PreFundState())
		err = ss.AddSignature(sig)
		if err != nil {
			return &s, NoSideEffects, WaitingForNothing, err
		}
		ok := updated.V.AddSignedState(ss)
		if !ok {
			return &s, NoSideEffects, WaitingForNothing, errors.New(`could not add prefund state`)
		}
		return &updated, NoSideEffects, WaitingForCompletePrefund, nil

	}

	if !updated.V.PreFundComplete() {
		return &updated, NoSideEffects, WaitingForCompletePrefund, nil
	}

	// Funding

	if !updated.requestedLedgerUpdates {
		updated.requestedLedgerUpdates = true
		return &updated, s.generateLedgerRequestSideEffects(), WaitingForCompleteFunding, nil
	}

	if !updated.fundingComplete() {
		return &updated, NoSideEffects, WaitingForCompleteFunding, nil
	}

	// Postfunding
	if !updated.V.PostFundSignedByMe() {
		sig, err := updated.V.PostFundState().Sign(*secretKey)
		if err != nil {
			return &s, NoSideEffects, WaitingForNothing, err
		}
		ss := state.NewSignedState(updated.V.PostFundState())
		err = ss.AddSignature(sig)
		if err != nil {
			return &s, NoSideEffects, WaitingForNothing, err
		}
		ok := updated.V.AddSignedState(ss)
		if !ok {
			return &s, NoSideEffects, WaitingForNothing, errors.New(`could not add postfund state`)
		}
		return &updated, NoSideEffects, WaitingForCompletePostFund, nil

	}

	if !updated.V.PostFundComplete() {
		return &updated, NoSideEffects, WaitingForCompletePostFund, nil
	}

	// Completion
	return &updated, NoSideEffects, WaitingForNothing, nil
}

func (s VirtualFundObjective) Channels() []types.Destination {
	ret := make([]types.Destination, 0, 3)
	ret = append(ret, s.V.Id)
	if s.MyRole != 0 {
		ret = append(ret, s.ToMyLeft.Channel.Id)
	}
	if s.MyRole != s.n+1 {
		ret = append(ret, s.ToMyRight.Channel.Id)
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
