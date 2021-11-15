package protocols

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/types"
)

const (
	WaitingForCompletePrefund  WaitingFor = "WaitingForCompletePrefund"
	WaitingForMyTurnToFund     WaitingFor = "WaitingForMyTurnToFund"
	WaitingForCompleteFunding  WaitingFor = "WaitingForCompleteFunding"
	WaitingForCompletePostFund WaitingFor = "WaitingForCompletePostFund"
	WaitingForNothing          WaitingFor = "WaitingForNothing" // Finished
)

func FundOnChainEffect(cId types.Destination, asset string, amount types.Funds) string {
	return "deposit" + amount.String() + "into" + cId.String()
}

var NoSideEffects = SideEffects{}

// errors
var ErrNotApproved = errors.New("objective not approved")

// DirectFundingObjectiveState is a cache of data computed by reading from the store. It stores (potentially) infinite data
type DirectFundingObjectiveState struct {
	Status    ObjectiveStatus
	ChannelId types.Destination

	ParticipantIndex map[types.Address]uint // the index for each participant
	ExpectedStates   []state.State          // indexed by turn number

	MyIndex uint // my participant index

	PreFundSigned []bool // indexed by participant.

	MyDepositSafetyThreshold types.Funds // if the on chain holdings are equal to this amount it is safe for me to deposit
	MyDepositTarget          types.Funds // I want to get the on chain holdings up to this much
	FullyFundedThreshold     types.Funds // if the on chain holdings are equal

	PostFundSigned []bool // indexed by participant

	OnChainHolding types.Funds
}

// NewDirectFundingObjectiveState initiates a DirectFundingInitialState with data calculated from
// the supplied initialState and client address
func NewDirectFundingObjectiveState(initialState state.State, myAddress types.Address) (DirectFundingObjectiveState, error) {
	var init DirectFundingObjectiveState
	var err error

	init.Status = Unapproved
	init.ChannelId, err = initialState.ChannelId()
	if err != nil {
		return init, err
	}
	for i, v := range initialState.Participants {
		init.ParticipantIndex[v] = uint(i)
	}

	init.ExpectedStates[0] = initialState

	fixed := initialState.FixedPart()
	init.ExpectedStates[1].ChainId = fixed.ChainId
	init.ExpectedStates[1].Participants = fixed.Participants
	init.ExpectedStates[1].ChannelNonce = fixed.ChannelNonce
	init.ExpectedStates[1].AppDefinition = fixed.AppDefinition
	init.ExpectedStates[1].ChallengeDuration = fixed.ChallengeDuration

	init.ExpectedStates[1].Outcome = initialState.Outcome
	init.ExpectedStates[1].AppData = initialState.AppData
	init.ExpectedStates[1].IsFinal = false
	init.ExpectedStates[1].TurnNum = big.NewInt(1)

	for i, v := range initialState.Participants {
		if v == myAddress { // todo: myAddress should really be something akin to myInterests, which could include internal destinations
			init.MyIndex = uint(i)
		}
	}

	if channelOutcome, err := outcome.Decode(initialState.VariablePart().EncodedOutcome); err == nil {
		init.FullyFundedThreshold = types.Funds{}

		for _, assetExit := range channelOutcome {
			assetAddress := assetExit.Asset
			sum := big.NewInt(0)
			threshold := big.NewInt(0)
			myShare := big.NewInt(0)

			for i, allocation := range assetExit.Allocations {
				sum = sum.Add(sum, allocation.Amount)

				if i < int(init.MyIndex) {
					threshold = threshold.Add(threshold, allocation.Amount)
				} else if i == int(init.MyIndex) {
					myShare = myShare.Add(myShare, allocation.Amount)
				}
			}

			init.FullyFundedThreshold[assetAddress] = sum
			init.MyDepositSafetyThreshold[assetAddress] = threshold
			init.MyDepositTarget[assetAddress] = myShare.Add(myShare, threshold)
		}
	}

	init.PreFundSigned = make([]bool, len(initialState.Participants))  // NOTE initialized to (false,false,...)
	init.PostFundSigned = make([]bool, len(initialState.Participants)) // NOTE initialized to (false,false,...)
	init.OnChainHolding = types.Funds{}

	return init, nil
}

// Public methods on the DirectFundingObjectiveState

func (s DirectFundingObjectiveState) Id() ObjectiveId {
	return ObjectiveId("DirectFunding-" + s.ChannelId.String())
}

func (s DirectFundingObjectiveState) Approve() Objective {
	updated := s.clone()
	// todo: consider case of s.Status == Rejected
	updated.Status = Approved

	return updated
}

func (s DirectFundingObjectiveState) Reject() Objective {
	updated := s.clone()
	updated.Status = Rejected
	return updated
}

// Update receives an ObjectiveEvent, applies all applicable event data to the DirectFundingObjectiveState,
// and returns the updated state
func (s DirectFundingObjectiveState) Update(event ObjectiveEvent) (Objective, error) {
	if s.ChannelId != event.ChannelId {
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
			return s, err
		}
	}

	if event.Holdings != nil {
		updated.OnChainHolding = event.Holdings
	}

	return updated, nil
}

// Crank inspects the extended state and declares a list of Effects to be executed
// It's like a state machine transition function where the finite / enumerable state is returned (computed from the extended state)
// rather than being independent of the extended state; and where there is only one type of event ("the crank") with no data on it at all
func (s DirectFundingObjectiveState) Crank(secretKey *[]byte) (Objective, SideEffects, WaitingFor, error) {
	updated := s.clone()

	// Input validation
	if updated.Status != Approved {
		return updated, NoSideEffects, WaitingForNothing, ErrNotApproved
	}

	// Prefunding
	if !updated.PreFundSigned[updated.MyIndex] {
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
		effects = append(effects, FundOnChainEffect(updated.ChannelId, `eth`, amountToDeposit))
		if len(effects) > 0 {
			// todo: effects as SideEffects{}
			return updated, NoSideEffects, WaitingForCompleteFunding, nil
		}
	}

	if !fundingComplete {
		return updated, NoSideEffects, WaitingForCompleteFunding, nil
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
	// todo: []string{"Objective" + s.ChannelId.String() + "complete"} as SideEffects{}
	return updated, NoSideEffects, WaitingForNothing, nil
}

//  Private methods on the DirectFundingObjectiveState

// prefundComplete returns true if all participants have signed a prefund state, as reflected by the extended state
func (s DirectFundingObjectiveState) prefundComplete() bool {
	for _, index := range s.ParticipantIndex {
		if !s.PreFundSigned[index] {
			return false
		}
	}
	return true
}

// postfundComplete returns true if all participants have signed a postfund state, as reflected by the extended state
func (s DirectFundingObjectiveState) postfundComplete() bool {
	for _, index := range s.ParticipantIndex {
		if !s.PostFundSigned[index] {
			return false
		}
	}
	return true
}

// fundingComplete returns true if the recorded OnChainHoldings are greater than or equal to the threshold for being fully funded.
func (s DirectFundingObjectiveState) fundingComplete() bool {
	for asset, threshold := range s.FullyFundedThreshold {
		chainHolding, ok := s.OnChainHolding[asset]

		if !ok {
			return false
		}

		if gt(threshold, chainHolding) {
			return false
		}
	}

	return true
}

// safeToDeposit returns true if the recorded OnChainHoldings are greater than or equal to the threshold for safety.
func (s DirectFundingObjectiveState) safeToDeposit() bool {
	for asset, safetyThreshold := range s.MyDepositSafetyThreshold {
		chainHolding, ok := s.OnChainHolding[asset]

		if !ok {
			return false
		}

		if gt(safetyThreshold, chainHolding) {
			return false
		}
	}

	return true
}

// amountToDeposit computes the appropriate amount to deposit given the current recorded OnChainHoldings
func (s DirectFundingObjectiveState) amountToDeposit() types.Funds {
	deposits := make(types.Funds, len(s.OnChainHolding))

	for asset, holding := range s.OnChainHolding {
		deposits[asset] = big.NewInt(0).Sub(s.MyDepositTarget[asset], holding)
	}

	return deposits
}

// SignatureReceived updates the objective's cache of which participants have signed which states
func (s DirectFundingObjectiveState) applySignature(signature state.Signature, turnNum int) error {
	signer, err := s.ExpectedStates[turnNum].RecoverSigner(signature)
	if err != nil {
		return err
	}

	index, ok := s.ParticipantIndex[signer]
	if !ok {
		return fmt.Errorf("signature recieved from unrecognized participant 0x%x", signer)
	}

	if turnNum == 0 {
		s.PreFundSigned[index] = true
	} else if turnNum == 1 {
		s.PostFundSigned[index] = true
	}

	return nil
}

// todo: is this sufficient? Particularly: s has pointer members (*big.Int)
func (s DirectFundingObjectiveState) clone() DirectFundingObjectiveState {
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
