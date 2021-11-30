package directfund

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

const (
	WaitingForCompletePrefund  protocols.WaitingFor = "WaitingForCompletePrefund"
	WaitingForMyTurnToFund     protocols.WaitingFor = "WaitingForMyTurnToFund"
	WaitingForCompleteFunding  protocols.WaitingFor = "WaitingForCompleteFunding"
	WaitingForCompletePostFund protocols.WaitingFor = "WaitingForCompletePostFund"
	WaitingForNothing          protocols.WaitingFor = "WaitingForNothing" // Finished
)

func FundOnChainEffect(cId types.Destination, asset string, amount types.Funds) string {
	return "deposit" + amount.String() + "into" + cId.String()
}

var NoSideEffects = protocols.SideEffects{}

// errors
var ErrNotApproved = errors.New("objective not approved")

// DirectFundObjective is a cache of data computed by reading from the store. It stores (potentially) infinite data
type DirectFundObjective struct {
	Status    protocols.ObjectiveStatus
	channelId types.Destination

	participantIndex map[types.Address]uint // the index for each participant
	expectedStates   []state.State          // indexed by turn number

	myIndex uint // my index in the Participants array of the initial state

	preFundSigned []bool // indexed by participant.

	myDepositSafetyThreshold types.Funds // if the on chain holdings are equal to this amount it is safe for me to deposit
	myDepositTarget          types.Funds // I want to get the on chain holdings up to this much
	fullyFundedThreshold     types.Funds // if the on chain holdings are equal

	postFundSigned []bool // indexed by participant

	onChainHolding types.Funds
}

// New initiates a DirectFundingInitialState with data calculated from
// the supplied initialState and client address
func New(initialState state.State, myAddress types.Address) (DirectFundObjective, error) {
	if initialState.IsFinal {
		return DirectFundObjective{}, errors.New("attempted to initiate new direct-funding objective with IsFinal == true")
	}

	var init DirectFundObjective
	var err error

	init.Status = protocols.Unapproved
	init.channelId, err = initialState.ChannelId()
	if err != nil {
		return init, err
	}
	init.participantIndex = make(map[types.Address]uint)
	for i, v := range initialState.Participants {
		init.participantIndex[v] = uint(i)
	}
	init.expectedStates = make([]state.State, 2)
	init.expectedStates[0] = initialState

	init.expectedStates[1] = initialState.Clone()
	init.expectedStates[1].TurnNum = big.NewInt(1)

	for i, v := range initialState.Participants {
		if v == myAddress {
			init.myIndex = uint(i)
		}
	}

	myAllocatedAmount := initialState.Outcome.TotalAllocatedFor(
		types.AdddressToDestination(myAddress),
	)

	init.fullyFundedThreshold = initialState.Outcome.TotalAllocated()
	init.myDepositSafetyThreshold = initialState.Outcome.DepositSafetyThreshold(
		types.AdddressToDestination(myAddress),
	)
	init.myDepositTarget = init.myDepositSafetyThreshold.Add(myAllocatedAmount)

	init.preFundSigned = make([]bool, len(initialState.Participants))  // NOTE initialized to (false,false,...)
	init.postFundSigned = make([]bool, len(initialState.Participants)) // NOTE initialized to (false,false,...)
	init.onChainHolding = types.Funds{}

	return init, nil
}

// Public methods on the DirectFundingObjectiveState

func (s DirectFundObjective) Id() protocols.ObjectiveId {
	return protocols.ObjectiveId("DirectFunding-" + s.channelId.String())
}

func (s DirectFundObjective) Approve() protocols.Objective {
	updated := s.clone()
	// todo: consider case of s.Status == Rejected
	updated.Status = protocols.Approved

	return updated
}

func (s DirectFundObjective) Reject() protocols.Objective {
	updated := s.clone()
	updated.Status = protocols.Rejected
	return updated
}

// Update receives an ObjectiveEvent, applies all applicable event data to the DirectFundingObjectiveState,
// and returns the updated state
func (s DirectFundObjective) Update(event protocols.ObjectiveEvent) (protocols.Objective, error) {
	if s.channelId != event.ChannelId {
		return s, errors.New("event and objective channelIds do not match")
	}

	updated := s.clone()

	for _, sig := range event.Sigs {
		var turnNum int
		if updated.prefundComplete() {
			turnNum = 1
		} else {
			turnNum = 0
		}

		err := updated.applySignature(sig, turnNum)
		if err != nil {
			// If there was an error applying the signature, log it and swallow it
			// This is a conscious choice (to ignore signatures we don't expect)
			// Examples include faulty signatures, signatures by non-participants, signatures on unexpected states, etc
			fmt.Println(err)
		}
	}

	if event.Holdings != nil {
		updated.onChainHolding = event.Holdings
	}

	return updated, nil
}

// Crank inspects the extended state and declares a list of Effects to be executed
// It's like a state machine transition function where the finite / enumerable state is returned (computed from the extended state)
// rather than being independent of the extended state; and where there is only one type of event ("the crank") with no data on it at all
func (s DirectFundObjective) Crank(secretKey *[]byte) (protocols.Objective, protocols.SideEffects, protocols.WaitingFor, error) {
	updated := s.clone()

	// Input validation
	if updated.Status != protocols.Approved {
		return updated, NoSideEffects, WaitingForNothing, ErrNotApproved
	}

	// Prefunding
	if !updated.preFundSigned[updated.myIndex] {
		// todo: {SignAndSendPreFundEffect(updated.ChannelId)} as SideEffects{}
		return updated, NoSideEffects, WaitingForCompletePrefund, nil
	}

	if !updated.prefundComplete() {
		return updated, NoSideEffects, WaitingForCompletePrefund, nil
	}

	// Funding
	fundingComplete := updated.fundingComplete() // note all information stored in state (since there are no real events)
	amountToDeposit := updated.amountToDeposit()
	safeToDeposit := updated.safeToDeposit()

	if !fundingComplete && !safeToDeposit {
		return updated, NoSideEffects, WaitingForMyTurnToFund, nil
	}

	if !fundingComplete && amountToDeposit.IsNonZero() && safeToDeposit {
		var effects = make([]string, 0) // TODO loop over assets
		effects = append(effects, FundOnChainEffect(updated.channelId, `eth`, amountToDeposit))
		if len(effects) > 0 {
			// todo: convert effects to SideEffects{} and return
			return updated, NoSideEffects, WaitingForCompleteFunding, nil
		}
	}

	if !fundingComplete {
		return updated, NoSideEffects, WaitingForCompleteFunding, nil
	}

	// Postfunding
	if !updated.postFundSigned[updated.myIndex] {
		// TODO sign the post fund state
		// TODO update updated.PostFundSigned[updated.MyIndex]
		// TODO prepare a message for peers with signature, return as SideEffects{}
		return updated, NoSideEffects, WaitingForCompletePostFund, nil
	}

	if !updated.postfundComplete() {
		return updated, NoSideEffects, WaitingForCompletePostFund, nil
	}

	// Completion
	return updated, NoSideEffects, WaitingForNothing, nil
}

//  Private methods on the DirectFundingObjectiveState

// prefundComplete returns true if all participants have signed a prefund state, as reflected by the extended state
func (s DirectFundObjective) prefundComplete() bool {
	for _, index := range s.participantIndex {
		if !s.preFundSigned[index] {
			return false
		}
	}
	return true
}

// postfundComplete returns true if all participants have signed a postfund state, as reflected by the extended state
func (s DirectFundObjective) postfundComplete() bool {
	for _, index := range s.participantIndex {
		if !s.postFundSigned[index] {
			return false
		}
	}
	return true
}

// fundingComplete returns true if the recorded OnChainHoldings are greater than or equal to the threshold for being fully funded.
func (s DirectFundObjective) fundingComplete() bool {
	for asset, threshold := range s.fullyFundedThreshold {
		chainHolding, ok := s.onChainHolding[asset]

		if !ok {
			return false
		}

		if types.Gt(threshold, chainHolding) {
			return false
		}
	}

	return true
}

// safeToDeposit returns true if the recorded OnChainHoldings are greater than or equal to the threshold for safety.
func (s DirectFundObjective) safeToDeposit() bool {
	for asset, safetyThreshold := range s.myDepositSafetyThreshold {
		chainHolding, ok := s.onChainHolding[asset]

		if !ok {
			return false
		}

		if types.Gt(safetyThreshold, chainHolding) {
			return false
		}
	}

	return true
}

// amountToDeposit computes the appropriate amount to deposit given the current recorded OnChainHoldings
func (s DirectFundObjective) amountToDeposit() types.Funds {
	deposits := make(types.Funds, len(s.onChainHolding))

	for asset, holding := range s.onChainHolding {
		deposits[asset] = big.NewInt(0).Sub(s.myDepositTarget[asset], holding)
	}

	return deposits
}

// applySignature updates the objective's cache of which participants have signed which states
func (s DirectFundObjective) applySignature(signature state.Signature, turnNum int) error {
	signer, err := s.expectedStates[turnNum].RecoverSigner(signature)
	if err != nil {
		return err
	}

	index, ok := s.participantIndex[signer]
	if !ok {
		return fmt.Errorf("signature received from unrecognized participant 0x%x", signer)
	}

	if turnNum == 0 {
		s.preFundSigned[index] = true
	} else if turnNum == 1 {
		s.postFundSigned[index] = true
	}
	return nil

}

// todo: is this sufficient? Particularly: s has pointer members (*big.Int)
func (s DirectFundObjective) clone() DirectFundObjective {
	return s
}

// mermaid diagram
// key:
// - effect!
// - waiting...
//
// https://mermaid-js.github.io/mermaid-live-editor/edit/#eyJjb2RlIjoiZ3JhcGggVERcbiAgICBTdGFydCAtLT4gQ3tJbnZhbGlkIElucHV0P31cbiAgICBDIC0tPnxZZXN8IEVbZXJyb3JdXG4gICAgQyAtLT58Tm98IEQwXG4gICAgXG4gICAgRDB7U2hvdWxkU2lnblByZUZ1bmR9XG4gICAgRDAgLS0-fFllc3wgUjFbU2lnblByZWZ1bmQhXVxuICAgIEQwIC0tPnxOb3wgRDFcbiAgICBcbiAgICBEMXtTYWZlVG9EZXBvc2l0ICY8YnI-ICFGdW5kaW5nQ29tcGxldGV9XG4gICAgRDEgLS0-IHxZZXN8IFIyW0Z1bmQgb24gY2hhaW4hXVxuICAgIEQxIC0tPiB8Tm98IEQyXG4gICAgXG4gICAgRDJ7IVNhZmVUb0RlcG9zaXQgJjxicj4gIUZ1bmRpbmdDb21wbGV0ZX1cbiAgICBEMiAtLT4gfFllc3wgUjNbXCJteSB0dXJuLi4uXCJdXG4gICAgRDIgLS0-IHxOb3wgRDNcblxuICAgIEQze1NhZmVUb0RlcG9zaXQgJjxicj4gIUZ1bmRpbmdDb21wbGV0ZX1cbiAgICBEMyAtLT4gfFllc3wgUjRbRGVwb3NpdCFdXG4gICAgRDMgLS0-IHxOb3wgRDRcblxuICAgIEQ0eyFGdW5kaW5nQ29tcGxldGV9XG4gICAgRDQgLS0-IHxZZXN8IFI1W1wiY29tcGxldGUgZnVuZGluZy4uLlwiXVxuICAgIEQ0IC0tPiB8Tm98IEQ1XG5cbiAgICBENXtTaG91bGRTaWduUHJlRnVuZH1cbiAgICBENSAtLT58WWVzfCBSNltTaWduUG9zdGZ1bmQhXVxuICAgIEQ1IC0tPnxOb3wgRDZcblxuICAgIEQ2eyFQb3N0RnVuZENvbXBsZXRlfVxuICAgIEQ2IC0tPnxZZXN8IFI3W1wiY29tcGxldGUgcG9zdGZ1bmQuLi5cIl1cbiAgICBENiAtLT58Tm98IFI4XG5cbiAgICBSOFtcImZpbmlzaFwiXVxuICAgIFxuXG5cbiIsIm1lcm1haWQiOiJ7fSIsInVwZGF0ZUVkaXRvciI6ZmFsc2UsImF1dG9TeW5jIjp0cnVlLCJ1cGRhdGVEaWFncmFtIjp0cnVlfQ//
// graph TD
//     Start --> C{Invalid Input?}
//     C -->|Yes| E[error]
//     C -->|No| D0

//     D0{ShouldSignPreFund}
//     D0 -->|Yes| R1[SignPrefund!]
//     D0 -->|No| D1

//     D1{SafeToDeposit &<br> !FundingComplete}
//     D1 --> |Yes| R2[Fund on chain!]
//     D1 --> |No| D2

//     D2{!SafeToDeposit &<br> !FundingComplete}
//     D2 --> |Yes| R3["wait my turn..."]
//     D2 --> |No| D3

//     D3{SafeToDeposit &<br> !FundingComplete}
//     D3 --> |Yes| R4[Deposit!]
//     D3 --> |No| D4

//     D4{!FundingComplete}
//     D4 --> |Yes| R5["wait for complete funding..."]
//     D4 --> |No| D5

//     D5{ShouldSignPostFund}
//     D5 -->|Yes| R6[SignPostfund!]
//     D5 -->|No| D6

//     D6{!PostFundComplete}
//     D6 -->|Yes| R7["wait for complete postfund..."]
//     D6 -->|No| R8

//     R8["finish"]
