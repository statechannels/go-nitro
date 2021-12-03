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

// VirtualFundObjective is a cache of data computed by reading from the store. It stores (potentially) infinite data
type VirtualFundObjective struct {
	Status protocols.ObjectiveStatus
	J      channel.Channel          // this is J
	L      map[uint]channel.Channel // this will contain 1 or 2 Ledger channels. For example L_0 if Alice (i=0). L_0 and L_1 for the first intermediary (i=1). L_n+1 for Bob (i=n+1)
	MyRole uint                     // index in the virtual funding protocol. 0 for Alice, n+1 for Bob. Otherwise, one of the intermediaries.

	a0 types.Funds // Initial balance for Alice
	b0 types.Funds // Initial balance for Bob

	ExpectedGuarantees map[uint]map[types.Address]outcome.Allocation // For each ledger channel, for each asset -- the expected guarantee that diverts funds from L_i to V,

	requestedLedgerUpdates bool // records that the ledger update side effects were previously generated (they may not have been executed yet)

	ParticipantIndex map[types.Address]uint // the index for each participant
	MyIndex          uint                   // my participant index in J

	preFundSigned  []bool // indexed by participant. TODO should this be initialized with my own index showing true?
	postFundSigned []bool // indexed by participant
}

// New initiates a VirtualFundObjective with data calculated from
// the supplied initialState and client address
func New(
	initialStateOfJ state.State,
	myAddress types.Address,
	myRole uint,
	ledgerChannelToMyLeft channel.Channel,
	ledgerChannelToMyRight channel.Channel,
) (VirtualFundObjective, error) {

	// TODO  validate that the Ledger cannels have isTwoPartyLedger=true

	var init VirtualFundObjective

	// Initialize channels
	init.J = channel.New(initialStateOfJ, false, types.Destination{}, types.Destination{})
	init.L = make(map[uint]channel.Channel)

	n := uint(2) // TODO  uint(len(s.L)) // n = numHops + 1 (the number of ledger channels)

	init.a0 = make(map[types.Address]*big.Int)
	init.b0 = make(map[types.Address]*big.Int)
	// Compute a0 and b0 from the initial state of J
	for i := range initialStateOfJ.Outcome {
		asset := initialStateOfJ.Outcome[i].Asset
		amount0 := initialStateOfJ.Outcome[i].Allocations[0].Amount
		amount1 := initialStateOfJ.Outcome[i].Allocations[1].Amount
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
	init.ExpectedGuarantees = make(map[uint]map[types.Address]outcome.Allocation)

	switch {
	case myRole == 0: // Alice
		init.L[0] = ledgerChannelToMyRight
		init.insertExpectedGuaranteesForLedgerChannel(0, init.L[0].MyDestination, init.L[0].TheirDestination) // This Ledger channel is on *my right*, so I am on the *left* of it
	case myRole < n+1: // Intermediary
		init.L[myRole] = ledgerChannelToMyRight
		init.L[myRole-1] = ledgerChannelToMyLeft
		init.insertExpectedGuaranteesForLedgerChannel(myRole, init.L[myRole].MyDestination, init.L[myRole].TheirDestination)       // This Ledger channel is on *my right*, so I am on the *left* of it
		init.insertExpectedGuaranteesForLedgerChannel(myRole-1, init.L[myRole-1].TheirDestination, init.L[myRole-1].MyDestination) // This Ledger channel is on *my left*, so I am on the *right* of it
	case myRole == n+1: // Bob
		init.L[myRole-1] = ledgerChannelToMyLeft
		init.insertExpectedGuaranteesForLedgerChannel(n, init.L[myRole-1].TheirDestination, init.L[myRole-1].MyDestination) // This Ledger channel is on *my left*, so I am on the *right* of it
	default: // Invalid

	}

	init.preFundSigned = make([]bool, len(initialStateOfJ.Participants))  // NOTE initialized to (false,false,...)
	init.postFundSigned = make([]bool, len(initialStateOfJ.Participants)) // NOTE initialized to (false,false,...)

	// TODO
	return init, nil
}

// insertExpectedGuaranteesForLedgerChannel mutates the VirtualFundObjective
func (init *VirtualFundObjective) insertExpectedGuaranteesForLedgerChannel(i uint, left types.Destination, right types.Destination) {
	expectedGuaranteesForLedgerChannel := make(map[types.Address]outcome.Allocation)
	metadata := outcome.GuaranteeMetadata{
		Left:  left,
		Right: right,
	}
	encodedGuarantee, _ := metadata.Encode() // TODO handle error
	for asset := range init.a0 {
		expectedGuaranteesForLedgerChannel[asset] = outcome.Allocation{
			Destination:    init.J.Id,
			Amount:         big.NewInt(0).Add(init.a0[asset], init.b0[asset]),
			AllocationType: outcome.GuaranteeAllocationType,
			Metadata:       encodedGuarantee,
		}
	}

	init.ExpectedGuarantees[i] = expectedGuaranteesForLedgerChannel
}

// Public methods on the VirtualFundObjective

// Id returns the objective id
func (s VirtualFundObjective) Id() protocols.ObjectiveId {
	return protocols.ObjectiveId("VirtualFundAsTerminal-" + s.J.Id.String())
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
		sig, _ := updated.J.PreFund.Sign(*secretKey)     // TODO handle error
		updated.J.AddSignedState(updated.J.PreFund, sig) // TODO handle return value (or not)
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
		sig, _ := updated.J.PostFund.Sign(*secretKey)     // TODO handle error
		updated.J.AddSignedState(updated.J.PostFund, sig) // TODO handle return value (or not)
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

	n := uint(len(s.L)) // n = numHops + 1 (the number of ledger channels)

	switch {
	case s.MyRole == 0: // Alice
		return s.ledgerChannelAffordsExpectedGuarantees(0)
	case s.MyRole < n+1: // Intermediary
		return s.ledgerChannelAffordsExpectedGuarantees(s.MyRole-1) && s.ledgerChannelAffordsExpectedGuarantees(s.MyRole)
	case s.MyRole == n+1: // Bob
		return s.ledgerChannelAffordsExpectedGuarantees(n)
	default: // Invalid
		return false
	}

}

// ledgerChannelAffordsExpectedGuarantees returns true if the ledger channel with given index i affords the expected guarantees for V
func (s VirtualFundObjective) ledgerChannelAffordsExpectedGuarantees(i uint) bool {
	return s.L[i].Affords(s.ExpectedGuarantees[i], s.L[i].OnChainFunding)
}

// generateLedgerRequestSideEffects generates the appropriate side effects, which (when executed and countersigned) will update 1 or 2 ledger channels to guarantee the joint channel
func (s VirtualFundObjective) generateLedgerRequestSideEffects() protocols.SideEffects {
	sideEffects := protocols.SideEffects{}
	sideEffects.LedgerRequests = make([]protocols.LedgerRequest, 2)
	if s.MyRole > 0 { // Not Alice
		sideEffects.LedgerRequests = append(sideEffects.LedgerRequests,
			protocols.LedgerRequest{
				LedgerId:    s.L[s.MyRole-1].Id,
				Destination: s.J.Id,
				Amount:      s.J.Total(),
				Guarantee:   []types.Address{s.J.FixedPart.Participants[s.MyRole-1], s.J.FixedPart.Participants[s.MyRole]},
			})
	}
	n := uint(len(s.L)) // n = numHops + 1 (the number of ledger channels)
	if s.MyRole < n {   // Not Bob
		sideEffects.LedgerRequests = append(sideEffects.LedgerRequests,
			protocols.LedgerRequest{
				LedgerId:    s.L[s.MyRole].Id,
				Destination: s.J.Id,
				Amount:      s.J.Total(),
				Guarantee:   []types.Address{s.J.FixedPart.Participants[s.MyRole], s.J.FixedPart.Participants[s.MyRole+1]},
			})
	}
	return sideEffects
}

// inScope returns true if the supplied channelId is the joint channel or one of the ledger channels. Can be used to filter out events that don't concern these channels.
func (s VirtualFundObjective) inScope(channelId types.Destination) bool {
	if channelId == s.J.Id {
		return true
	}
	for _, channel := range s.L {
		if channelId == channel.Id {
			return true
		}
	}

	return false
}

// todo: is this sufficient? Particularly: s has pointer members (*big.Int)
func (s VirtualFundObjective) clone() VirtualFundObjective {
	return s
}
