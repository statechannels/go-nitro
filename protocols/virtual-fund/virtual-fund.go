package virtualfund

import (
	"errors"

	"github.com/statechannels/go-nitro/channel"
	"github.com/statechannels/go-nitro/channel/state"
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

	requestedLedgerUpdates bool // records that the ledger update side effects were previously generated (they may not have been executed yet)

	ParticipantIndex map[types.Address]uint // the index for each participant
	MyIndex          uint                   // my participant index in J

	PreFundSigned  []bool // indexed by participant. TODO should this be initialized with my own index showing true?
	PostFundSigned []bool // indexed by participant
}

// New initiates a VirtualFundObjective with data calculated from
// the supplied initialState and client address
func New(initialStateOfJ state.State, myAddress types.Address) (VirtualFundObjective, error) {
	var init VirtualFundObjective
	// TODO
	return init, nil
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
	if !updated.PreFundSigned[updated.MyIndex] {
		// todo sign the prefund
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
	if !updated.PostFundSigned[updated.MyIndex] {
		// todo: sign the postfund
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
		if !s.PreFundSigned[index] {
			return false
		}
	}
	return true
}

// postfundComplete returns true if all participants have signed a postfund state, as reflected by the extended state
func (s VirtualFundObjective) postfundComplete() bool {
	for _, index := range s.ParticipantIndex {
		if !s.PostFundSigned[index] {
			return false
		}
	}
	return true
}

// fundingComplete returns true if the appropriate ledger channel guarantees sufficient funds for J
func (s VirtualFundObjective) fundingComplete() bool {

	// Each peer commits to an update in L_{i-1} and L_i including the guarantees G_{i-1} and {G_i} respectively, and deducting b_0 from L_{I-1} and a_0 from L_i.
	// A = P_0 and B=P_n+1 are special cases. A only does the guarantee for L_0, and B only foes the guarantee for L_n.

	n := uint(len(s.L)) // n = numHops + 1 (the number of ledger channels)

	switch {
	case s.MyRole == 0: // Alice
		// return s.L[n].Affords()
	case s.MyRole < n+1: // Intermediary
	case s.MyRole == n+1: // Bob
	case s.MyRole > n+1: // Invalid

	}

	// if s.MyRole == n+1 { // If I'm Bob, or peer n+1
	// 	return s.L[n].GuaranteesFor(s.J.Id).IsNonZero() // TODO a proper check on each asset (against s.J.Total)
	// } else {
	// 	return s.L[s.MyRole].GuaranteesFor(s.J.Id).IsNonZero() // TODO a proper check on each asset (against s.J.Total)
	// }

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
