package protocols

import (
	"errors"

	"github.com/statechannels/go-nitro/channel"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

// TODO use round 1/2/3 terminology?
const (
	WaitingForCompletePrefund  protocols.WaitingFor = "WaitingForCompletePrefund"
	WaitingForCompleteFunding  protocols.WaitingFor = "WaitingForCompleteFunding"
	WaitingForCompletePostFund protocols.WaitingFor = "WaitingForCompletePostFund"
	WaitingForNothing          protocols.WaitingFor = "WaitingForNothing" // Finished
)

var NoSideEffects = protocols.SideEffects{}

// errors
var ErrNotApproved = errors.New("objective not approved")

// VirtualFundingObjectiveState is a cache of data computed by reading from the store. It stores (potentially) infinite data
type VirtualFundingObjectiveState struct {
	Status protocols.ObjectiveStatus
	J      channel.Channel   // this is J
	L      []channel.Channel // this is {L_i} (a vector of Ledger channels)
	MyRole uint              // 0 for Alice, n+1 for Bob. Otherwise, one of the intermediaries.

	requestedLedgerUpdates bool // records that the ledger update side effects were previously generated (they may not have been executed yet)

	ParticipantIndex map[types.Address]uint // the index for each participant
	ExpectedStates   []state.State          // indexed by turn number

	MyIndex uint // my participant index in J

	PreFundSigned []bool // indexed by participant. TODO should this be initialized with my own index showing true?

	PostFundSigned []bool // indexed by participant
}

// NewVirtualFundingObjectiveState initiates a DirectFundingInitialState with data calculated from
// the supplied initialState and client address
func NewVirtualFundingObjectiveState(initialState state.State, myAddress types.Address) (VirtualFundingObjectiveState, error) {
	var init VirtualFundingObjectiveState
	// TODO
	return init, nil
}

// Public methods on the VirtualFundingObjectiveState

func (s VirtualFundingObjectiveState) Id() protocols.ObjectiveId {
	return protocols.ObjectiveId("VirtualFundingAsTerminal-" + s.J.Id.String())
}

func (s VirtualFundingObjectiveState) Approve() protocols.Objective {
	updated := s.clone()
	// todo: consider case of s.Status == Rejected
	updated.Status = protocols.Approved
	return updated
}

func (s VirtualFundingObjectiveState) Reject() protocols.Objective {
	updated := s.clone()
	updated.Status = protocols.Rejected
	return updated
}

// Update receives an protocols.ObjectiveEvent, applies all applicable event data to the VirtualFundingObjectiveState,
// and returns the updated state
func (s VirtualFundingObjectiveState) Update(event protocols.ObjectiveEvent) (protocols.Objective, error) {

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
func (s VirtualFundingObjectiveState) Crank(secretKey *[]byte) (protocols.Objective, protocols.SideEffects, protocols.WaitingFor, error) {
	updated := s.clone()

	// Input validation
	if updated.Status != protocols.Approved {
		return updated, NoSideEffects, WaitingForNothing, ErrNotApproved
	}

	// Prefunding
	if !updated.PreFundSigned[updated.MyIndex] {

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

	fundingComplete := updated.fundingComplete()
	if !fundingComplete {
		return updated, NoSideEffects, protocols.WaitingForCompleteFunding, nil
	}
	// Postfunding
	if !updated.PostFundSigned[updated.MyIndex] {
		// todo: []string{SignPostFundEffect(updated.ChannelId)} as SideEffects{}
		return updated, NoSideEffects, WaitingForCompletePostFund, nil
	}

	if !updated.postfundComplete() {
		return updated, NoSideEffects, WaitingForCompletePostFund, nil
	}

	// Completion
	// todo: []string{"protocols.Objective" + s.ChannelId.String() + "complete"} as SideEffects{}
	return updated, NoSideEffects, WaitingForNothing, nil
}

//  Private methods on the VirtualFundingObjectiveState

// prefundComplete returns true if all participants have signed a prefund state, as reflected by the extended state
func (s VirtualFundingObjectiveState) prefundComplete() bool {
	for _, index := range s.ParticipantIndex {
		if !s.PreFundSigned[index] {
			return false
		}
	}
	return true
}

// postfundComplete returns true if all participants have signed a postfund state, as reflected by the extended state
func (s VirtualFundingObjectiveState) postfundComplete() bool {
	for _, index := range s.ParticipantIndex {
		if !s.PostFundSigned[index] {
			return false
		}
	}
	return true
}

// fundingComplete returns true if the appropriate ledger channel guarantees sufficient funds for J
func (s VirtualFundingObjectiveState) fundingComplete() bool {

	n := uint(len(s.L)) // n = numHops + 1 (the number of ledger channels)

	if s.MyRole == n+1 { // If I'm Bob, or peer n+1
		return s.L[n].GuaranteesFor(s.J.Id).IsNonZero() // TODO a proper check
	} else {
		return s.L[s.MyRole].GuaranteesFor(s.J.Id).IsNonZero() // TODO a proper check
	}

}

func (s VirtualFundingObjectiveState) generateLedgerRequestSideEffects() protocols.SideEffects {
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

func (s VirtualFundingObjectiveState) inScope(channelId types.Bytes32) bool {
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
func (s VirtualFundingObjectiveState) clone() VirtualFundingObjectiveState {
	return s
}
