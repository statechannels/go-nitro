package virtualfund

import (
	"errors"
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
	Channel            channel.Channel
	ExpectedGuarantees map[types.Address]outcome.Allocation
}

// VirtualFundObjective is a cache of data computed by reading from the store. It stores (potentially) infinite data
type VirtualFundObjective struct {
	Status protocols.ObjectiveStatus
	V      channel.Channel // this is J

	ToMyLeft  *Connection
	ToMyRight *Connection

	// n uint // TODO
	MyRole uint // index in the virtual funding protocol. 0 for Alice, n+1 for Bob. Otherwise, one of the intermediaries.

	a0 types.Funds // Initial balance for Alice
	b0 types.Funds // Initial balance for Bob

	requestedLedgerUpdates bool // records that the ledger update side effects were previously generated (they may not have been executed yet)

	ParticipantIndex map[types.Address]uint // the index for each participant
	MyIndex          uint                   // my participant index in J

	preFundSigned  []bool // indexed by participant. TODO should this be initialized with my own index showing true?
	postFundSigned []bool // indexed by participant
}

// New initiates a VirtualFundObjective with data calculated from
// the supplied initialState and client address
func New(
	initialStateOfV state.State,
	myAddress types.Address,
	myRole uint,
	ledgerChannelToMyLeft channel.Channel,
	ledgerChannelToMyRight channel.Channel,
) (VirtualFundObjective, error) {

	// TODO  validate that the Ledger cannels have isTwoPartyLedger=true

	var init VirtualFundObjective

	// Initialize channels
	init.V = channel.New(initialStateOfV, false, types.Destination{}, types.Destination{})

	n := uint(2) // TODO  uint(len(s.L)) // n = numHops + 1 (the number of ledger channels)

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

	init.MyRole = myRole // this should not be modified, so no need to make a new big.Int

	switch {
	case myRole == 0: // Alice
		init.ToMyRight = &Connection{}
		init.ToMyRight.Channel = ledgerChannelToMyRight
		init.ToMyRight.insertExpectedGuarantees(init.a0, init.b0, init.V.Id, init.ToMyRight.Channel.MyDestination, init.ToMyRight.Channel.TheirDestination)
	case myRole < n+1: // Intermediary
		init.ToMyRight = &Connection{}
		init.ToMyRight.Channel = ledgerChannelToMyRight
		init.ToMyRight.insertExpectedGuarantees(init.a0, init.b0, init.V.Id, init.ToMyRight.Channel.MyDestination, init.ToMyRight.Channel.TheirDestination)
		init.ToMyLeft = &Connection{}
		init.ToMyLeft.Channel = ledgerChannelToMyLeft
		init.ToMyLeft.insertExpectedGuarantees(init.a0, init.b0, init.V.Id, init.ToMyRight.Channel.TheirDestination, init.ToMyRight.Channel.MyDestination)
	case myRole == n+1: // Bob
		init.ToMyLeft = &Connection{}
		init.ToMyLeft.Channel = ledgerChannelToMyLeft
		init.ToMyLeft.insertExpectedGuarantees(init.a0, init.b0, init.V.Id, init.ToMyRight.Channel.TheirDestination, init.ToMyRight.Channel.MyDestination)
	default: // Invalid

	}

	init.preFundSigned = make([]bool, len(initialStateOfV.Participants))  // NOTE initialized to (false,false,...)
	init.postFundSigned = make([]bool, len(initialStateOfV.Participants)) // NOTE initialized to (false,false,...)

	// TODO
	return init, nil
}

// insertExpectedGuaranteesForLedgerChannel mutates the reciever Connection struct
func (connection *Connection) insertExpectedGuarantees(a0 types.Funds, b0 types.Funds, vId types.Destination, left types.Destination, right types.Destination) {
	expectedGuaranteesForLedgerChannel := make(map[types.Address]outcome.Allocation)
	metadata := outcome.GuaranteeMetadata{
		Left:  left,
		Right: right,
	}
	encodedGuarantee, _ := metadata.Encode() // TODO handle error
	for asset := range a0 {
		expectedGuaranteesForLedgerChannel[asset] = outcome.Allocation{
			Destination:    vId,
			Amount:         big.NewInt(0).Add(a0[asset], b0[asset]),
			AllocationType: outcome.GuaranteeAllocationType,
			Metadata:       encodedGuarantee,
		}
	}

	connection.ExpectedGuarantees = expectedGuaranteesForLedgerChannel
}

// Public methods on the VirtualFundObjective

// Id returns the objective id
func (s VirtualFundObjective) Id() protocols.ObjectiveId {
	return protocols.ObjectiveId("VirtualFundAsTerminal-" + s.V.Id.String())
}

// Approve returns an approved copy of the objective
func (s VirtualFundObjective) Approve() protocols.Objective {
	updated := s.clone()
	// todo: consider case of s.Status == Rejected
	updated.Status = protocols.Approved
	return updated
}

// Approve returns a rejected copy of the objective
func (s VirtualFundObjective) Reject() protocols.Objective {
	updated := s.clone()
	updated.Status = protocols.Rejected
	return updated
}

// Update receives an protocols.ObjectiveEvent, applies all applicable event data to the VirtualFundObjective,
// and returns the updated state
func (s VirtualFundObjective) Update(event protocols.ObjectiveEvent) (protocols.Objective, error) {

	if !s.inScope(event.ChannelId) {
		return s, errors.New("event channelId out of scope of objective")
	}

	updated := s.clone()

	// TODO

	return updated, nil
}

// Crank inspects the extended state and declares a list of Effects to be executed
// It's like a state machine transition function where the finite / enumerable state is returned (computed from the extended state)
// rather than being independent of the extended state; and where there is only one type of event ("the crank") with no data on it at all
func (s VirtualFundObjective) Crank(secretKey *[]byte) (protocols.Objective, protocols.SideEffects, protocols.WaitingFor, error) {
	updated := s.clone()

	// Input validation
	if updated.Status != protocols.Approved {
		return updated, NoSideEffects, WaitingForNothing, ErrNotApproved
	}

	// TODO could perform checks on s.L (should only have 1 or 2 channels in there)
	// Prefunding
	if !updated.preFundSigned[updated.MyIndex] {
		sig, _ := updated.V.PreFund.Sign(*secretKey)     // TODO handle error
		updated.V.AddSignedState(updated.V.PreFund, sig) // TODO handle return value (or not)
		updated.preFundSigned[updated.MyIndex] = true
		return updated, NoSideEffects, WaitingForCompletePrefund, nil
	}

	if !updated.prefundComplete() {
		return updated, NoSideEffects, WaitingForCompletePrefund, nil
	}

	// Funding

	if !updated.requestedLedgerUpdates {
		updated.requestedLedgerUpdates = true
		return updated, s.generateLedgerRequestSideEffects(), WaitingForCompleteFunding, nil
	}

	if !updated.fundingComplete() {
		return updated, NoSideEffects, protocols.WaitingForCompleteFunding, nil
	}

	// Postfunding
	if !updated.postFundSigned[updated.MyIndex] {
		sig, _ := updated.V.PostFund.Sign(*secretKey)     // TODO handle error
		updated.V.AddSignedState(updated.V.PostFund, sig) // TODO handle return value (or not)
		updated.postFundSigned[updated.MyIndex] = true
		return updated, NoSideEffects, WaitingForCompletePostFund, nil
	}

	if !updated.postfundComplete() {
		return updated, NoSideEffects, WaitingForCompletePostFund, nil
	}

	// Completion
	return updated, NoSideEffects, WaitingForNothing, nil
}

//  Private methods on the VirtualFundObjective

// prefundComplete returns true if all participants have signed a prefund state, as reflected by the extended state
func (s VirtualFundObjective) prefundComplete() bool {
	for _, index := range s.ParticipantIndex {
		if !s.preFundSigned[index] {
			return false
		}
	}
	return true
}

// postfundComplete returns true if all participants have signed a postfund state, as reflected by the extended state
func (s VirtualFundObjective) postfundComplete() bool {
	for _, index := range s.ParticipantIndex {
		if !s.postFundSigned[index] {
			return false
		}
	}
	return true
}

// fundingComplete returns true if the appropriate ledger channel guarantees sufficient funds for J
func (s VirtualFundObjective) fundingComplete() bool {

	// Each peer commits to an update in L_{i-1} and L_i including the guarantees G_{i-1} and {G_i} respectively, and deducting b_0 from L_{I-1} and a_0 from L_i.
	// A = P_0 and B=P_n+1 are special cases. A only does the guarantee for L_0 (deducting a0), and B only foes the guarantee for L_n (deducting b0).

	n := uint(2) // n = numHops + 1 (the number of ledger channels)

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

// TODO
func (connection *Connection) ledgerChannelAffordsExpectedGuarantees() bool {
	return connection.Channel.Affords(connection.ExpectedGuarantees, connection.Channel.OnChainFunding)
}

// generateLedgerRequestSideEffects generates the appropriate side effects, which (when executed and countersigned) will update 1 or 2 ledger channels to guarantee the joint channel
func (s VirtualFundObjective) generateLedgerRequestSideEffects() protocols.SideEffects {
	sideEffects := protocols.SideEffects{}
	sideEffects.LedgerRequests = make([]protocols.LedgerRequest, 0)
	if s.MyRole > 0 { // Not Alice
		sideEffects.LedgerRequests = append(sideEffects.LedgerRequests,
			protocols.LedgerRequest{
				LedgerId:    s.ToMyLeft.Channel.Id,
				Destination: s.V.Id,
				Amount:      s.V.Total(),
				Left:        s.ToMyLeft.Channel.TheirDestination,
				Right:       s.ToMyLeft.Channel.MyDestination,
			})
	}
	n := uint(2)      // n = numHops + 1 (the number of ledger channels)
	if s.MyRole < n { // Not Bob
		sideEffects.LedgerRequests = append(sideEffects.LedgerRequests,
			protocols.LedgerRequest{
				LedgerId:    s.ToMyRight.Channel.Id,
				Destination: s.V.Id,
				Amount:      s.V.Total(),
				Left:        s.ToMyRight.Channel.MyDestination,
				Right:       s.ToMyRight.Channel.TheirDestination,
			})
	}
	return sideEffects
}

// inScope returns true if the supplied channelId is the joint channel or one of the ledger channels. Can be used to filter out events that don't concern these channels.
func (s VirtualFundObjective) inScope(channelId types.Destination) bool {

	switch channelId {
	case s.V.Id:
		return true
	case s.ToMyLeft.Channel.Id:
		return true
	case s.ToMyRight.Channel.Id:
		return true
	}

	return false
}

// todo: is this sufficient? Particularly: s has pointer members (*big.Int)
func (s VirtualFundObjective) clone() VirtualFundObjective {
	return s
}
