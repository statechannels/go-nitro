package protocols

import (
	"math/big"

	"github.com/statechannels/go-nitro/channel/state"
)

// TODO this protocol does not specify how events are handled at all
// (it assumes that events are handled by pushing information into the store)

// Crank inspects the extended state and declares a list of Effects to be executed
// It's like a state machine transition function where the finite / enumerable state is returned (computed from the extended state)
// rather than being independent of the extended state; and where there is only one type of event ("the crank") with no data on it at all
func (s DirectFundingObjectiveState) Crank() (SideEffects, WaitingFor, error) {

	// Input validation
	if s.Status != Approved {
		return NoSideEffects, WaitingForNothing, ErrNotApproved
	}

	// Prefunding
	if !s.PreFundSigned[s.MyIndex] {
		return []string{SignPreFundEffect(s.ChannelId)}, WaitingForCompletePrefund, nil
	}
	if !s.PrefundComplete() {
		return NoSideEffects, WaitingForCompletePrefund, nil
	}

	// Funding
	fundingComplete := s.FundingComplete(s.OnChainHolding) // note all information stored in state (since there are no real events)
	// (contrast this with a FSM where we have the new on chain holding on the event)
	amountToDeposit := s.AmountToDeposit(s.OnChainHolding)
	safeToDeposit := s.SafeToDeposit(s.OnChainHolding)

	if !fundingComplete && !safeToDeposit {
		return []string{}, WaitingForMyTurnToFund, nil
	}

	if !fundingComplete && gt(amountToDeposit, zero) && safeToDeposit {
		var effects = make([]string, 0) // TODO loop over assets
		effects = append(effects, FundOnChainEffect(s.ChannelId, `eth`, amountToDeposit))
		if len(effects) > 0 {
			return effects, WaitingForCompleteFunding, nil
		}
	}

	if !fundingComplete {
		return NoSideEffects, WaitingForCompleteFunding, nil
	}

	// Postfunding
	if !s.PostFundSigned[s.MyIndex] {
		return []string{SignPostFundEffect(s.ChannelId)}, WaitingForCompletePostFund, nil
	}

	if !s.PostfundComplete() {
		return NoSideEffects, WaitingForCompletePostFund, nil
	}

	// Completion
	return []string{"Objective" + s.ChannelId.String() + "complete"}, WaitingForNothing, nil
}

func (s DirectFundingObjectiveState) Approve() (DirectFundingObjectiveState, error) {
	updated := s.Clone()
	// todo: consider case of s.Status == Rejected
	updated.Status = Approved

	return updated, nil
}

func InitializeDirectFundingObjectiveState(initialState state.State) (DirectFundingObjectiveState, error) {
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

	// init.MyIndex						// requires myAddress

	// init.MyDepositSafetyThreshold	// \
	// init.MyDepositTarget				//  -> further fcn inputs?
	// init.FullyFundedThreshold		// /

	init.PostFundSigned = make([]bool, len(initialState.Participants))
	init.OnChainHolding = big.NewInt(0)

	return init, nil
}

// SignatureReceived updates the objective's cache of which participants have signed which states
func (s DirectFundingObjectiveState) SignatureRecieved(signature state.Signature, turnNum int) (DirectFundingObjectiveState, error) {
	updated := s.Clone()

	signer, err := updated.ExpectedStates[turnNum].RecoverSigner(signature)
	index, ok := updated.ParticipantIndex[signer]

	if ok && err == nil {
		if turnNum == 0 {
			updated.PreFundSigned[index] = true
		}
	}

	return updated, nil
}

// todo: is this sufficient? Particularly: s has pointer members (*big.Int)
func (s DirectFundingObjectiveState) Clone() DirectFundingObjectiveState {
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
//     D2 --> |Yes| R3["my turn..."]
//     D2 --> |No| D3

//     D3{SafeToDeposit &<br> !FundingComplete}
//     D3 --> |Yes| R4[Deposit!]
//     D3 --> |No| D4

//     D4{!FundingComplete}
//     D4 --> |Yes| R5["complete funding..."]
//     D4 --> |No| D5

//     D5{ShouldSignPreFund}
//     D5 -->|Yes| R6[SignPostfund!]
//     D5 -->|No| D6

//     D6{!PostFundComplete}
//     D6 -->|Yes| R7["complete postfund..."]
//     D6 -->|No| R8

//     R8["finish"]
