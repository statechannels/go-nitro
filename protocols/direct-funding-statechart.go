package protocols

import (
	"errors"
	"math/big"

	"github.com/statechannels/go-nitro/channel/state"
)

type DirectFundingExtendedState = DirectFundingObjectiveState // DirectFundingExtendedState contains the (potentially infinite) extended state of the Direct Funding machine.
// The extended state of this machine is a cache of a larger store of states and events stored by a Nitro wallet.
// This struct should be kept as shallow copyable (and this should be tested).

// DirectFundingProtocolState has both enumerable and extended state components.
type DirectFundingProtocolState struct {
	EnumerableState DirectFundingEnumerableState
	ExtendedState   DirectFundingExtendedState
}

// The event types for the state machine are enumerated.
type DirectFundingProtocolEventType int

const (
	PreFundReceived DirectFundingProtocolEventType = iota
	FundingUpdated
	PostFundReceived
)

// DirectFundingProtocolEvent has a type as well as other rich information (which may or may not be non nil).
type DirectFundingProtocolEvent struct {
	Type           DirectFundingProtocolEventType
	State          state.State
	Signature      state.Signature
	OnChainHolding *big.Int
}

// TODO can reducers be abstracted into an interface?

// NextState is the overall reducer / state transition function for the DirectFundingProtocol
func (s DirectFundingProtocolState) NextState(e DirectFundingProtocolEvent) (DirectFundingProtocolState, SideEffects, error) {
	if s.ExtendedState.Status != Approved {
		return s, NoSideEffects, ErrNotApproved
	}
	// it is better to switch on the state than on the event
	// https://dev.to/davidkpiano/you-don-t-need-a-library-for-state-machines-k7h
	switch s.EnumerableState {
	case WaitingForCompletePrefund:
		return s.nextStateFromWaitingForCompletePrefund(e)
	case WaitingForMyTurnToFund:
		return s.nextStateFromMyWaitingForMyTurnToFund(e)
	case WaitingForCompleteFunding:
		return s.nextStateFromWaitingForCompleteFunding(e)
	case WaitingForCompletePostFund:
		return s.nextStateFromWaitingForCompletePostFund(e)
	default:
		return s, NoSideEffects, nil
	}
}

// nextStateFromWaitingForCompletePrefund is a component of the overall DirectFundingProtocol reducer
// TODO when do we sign and send our own prefund state? When we construct the machine?
func (s DirectFundingProtocolState) nextStateFromWaitingForCompletePrefund(e DirectFundingProtocolEvent) (DirectFundingProtocolState, SideEffects, error) {
	if e.Type != PreFundReceived { // There's only one way out of this state
		return s, NoSideEffects, nil
	}
	newExtendedState := s.ExtendedState // Make a copy of the extended state because we anticipate needing to return an updated version

	signer, err := e.State.RecoverSigner(e.Signature)
	if err != nil {
		return s, NoSideEffects, err
	}

	signerIndex, present := newExtendedState.ParticipantIndex[signer]
	if !present {
		return s, NoSideEffects, errors.New(`signer is not a participant`)
	} else {
		newExtendedState.PreFundSigned[signerIndex] = true
	}

	if newExtendedState.PrefundComplete() {
		return DirectFundingProtocolState{WaitingForMyTurnToFund, newExtendedState}, NoSideEffects, nil
	} else {
		return DirectFundingProtocolState{WaitingForCompletePrefund, newExtendedState}, NoSideEffects, nil
	}
}

// nextStateFromFundingIncomplete is a component of the overall DirectFundingProtocol reducer
func (s DirectFundingProtocolState) nextStateFromMyWaitingForMyTurnToFund(e DirectFundingProtocolEvent) (DirectFundingProtocolState, SideEffects, error) {
	if e.Type != FundingUpdated { // There's only one way out of this state
		return s, NoSideEffects, nil
	}

	if gte(e.OnChainHolding, s.ExtendedState.MyDepositTarget) {
		// Can move to a new enumerable state
		return DirectFundingProtocolState{WaitingForCompleteFunding, s.ExtendedState}, NoSideEffects, nil
	}

	if s.ExtendedState.SafeToDeposit(e.OnChainHolding) {
		// Onlty here is it safe to deposit
		depositAmount := s.ExtendedState.AmountToDeposit(e.OnChainHolding)
		return s, SideEffects{FundOnChainEffect(s.ExtendedState.ChannelId, `ETH`, depositAmount)}, nil
	}

	return s, NoSideEffects, nil
}

// nextStateFromFundingIncomplete is a component of the overall DirectFundingProtocol reducer
func (s DirectFundingProtocolState) nextStateFromWaitingForCompleteFunding(e DirectFundingProtocolEvent) (DirectFundingProtocolState, SideEffects, error) {
	if e.Type != FundingUpdated { // There's only one way out of this state
		return s, NoSideEffects, nil
	}

	if s.ExtendedState.FundingComplete((e.OnChainHolding)) {
		// We make can progess to the next enumerable state
		return DirectFundingProtocolState{WaitingForCompletePostFund, s.ExtendedState}, SideEffects{SignPostFundEffect(s.ExtendedState.ChannelId)}, nil
	}

	return s, NoSideEffects, nil

}

// nextStateFromPostfundIncomplete is a component of the overall DirectFundingProtocol reducer
func (s DirectFundingProtocolState) nextStateFromWaitingForCompletePostFund(e DirectFundingProtocolEvent) (DirectFundingProtocolState, SideEffects, error) {
	if e.Type != PostFundReceived { // There's only one way out of this state
		return s, NoSideEffects, nil
	}
	newExtendedState := s.ExtendedState // Make a copy of the extended state because we anticipate needing to return an updated version

	signer, err := e.State.RecoverSigner(e.Signature)
	if err != nil {
		return s, NoSideEffects, err
	}

	signerIndex, present := newExtendedState.ParticipantIndex[signer]
	if !present {
		return s, NoSideEffects, errors.New(`signer is not a participant`)
	} else {
		newExtendedState.PostFundSigned[signerIndex] = true
	}

	if newExtendedState.PostfundComplete() {
		return DirectFundingProtocolState{WaitingForNothing, newExtendedState}, NoSideEffects, nil
	} else {
		return DirectFundingProtocolState{WaitingForCompletePostFund, newExtendedState}, NoSideEffects, nil
	}
}

// mermaid diagram
// https://mermaid-js.github.io/mermaid-live-editor/edit#eyJjb2RlIjoic3RhdGVEaWFncmFtLXYyXG4gICAgV2FpdGluZ0ZvckNvbXBsZXRlUHJlZnVuZFxuXHRXYWl0aW5nRm9yTXlUdXJuVG9GdW5kXG5cdFdhaXRpbmdGb3JDb21wbGV0ZUZ1bmRpbmdcblx0V2FpdGluZ0ZvckNvbXBsZXRlUG9zdEZ1bmRcblxuXG4gICAgWypdIC0tPldhaXRpbmdGb3JDb21wbGV0ZVByZWZ1bmQgXG4gICAgV2FpdGluZ0ZvckNvbXBsZXRlUHJlZnVuZCAtLT4gV2FpdGluZ0Zvck15VHVyblRvRnVuZCA6IFByZUZ1bmRSZWNlaXZlZFxuICAgIFdhaXRpbmdGb3JDb21wbGV0ZVByZWZ1bmQgLS0-IFdhaXRpbmdGb3JDb21wbGV0ZVByZWZ1bmQgOiBQcmVGdW5kUmVjZWl2ZWRcbiAgICBXYWl0aW5nRm9yTXlUdXJuVG9GdW5kIC0tPiBXYWl0aW5nRm9yQ29tcGxldGVGdW5kaW5nIDogRnVuZGluZ1VwZGF0ZWRcbiAgICBXYWl0aW5nRm9yTXlUdXJuVG9GdW5kIC0tPiBXYWl0aW5nRm9yTXlUdXJuVG9GdW5kIDogRnVuZGluZ1VwZGF0ZWRcbiAgICBXYWl0aW5nRm9yQ29tcGxldGVGdW5kaW5nIC0tPiBXYWl0aW5nRm9yQ29tcGxldGVQb3N0RnVuZCA6IEZ1bmRpbmdVcGRhdGVkXG4gICAgV2FpdGluZ0ZvckNvbXBsZXRlRnVuZGluZyAtLT4gV2FpdGluZ0ZvckNvbXBsZXRlRnVuZGluZyA6IEZ1bmRpbmdVcGRhdGVkXG4gICAgV2FpdGluZ0ZvckNvbXBsZXRlUG9zdEZ1bmQgLS0-IFsqXSA6IFBvc3RGdW5kUmVjZWl2ZWRcblxuICAgICAgICAgICAgIiwibWVybWFpZCI6IntcbiAgXCJ0aGVtZVwiOiBcImRlZmF1bHRcIlxufSIsInVwZGF0ZUVkaXRvciI6ZmFsc2UsImF1dG9TeW5jIjp0cnVlLCJ1cGRhdGVEaWFncmFtIjpmYWxzZX0
// stateDiagram-v2
//     WaitingForCompletePrefund
// 	WaitingForMyTurnToFund
// 	WaitingForCompleteFunding
// 	WaitingForCompletePostFund

//     [*] -->WaitingForCompletePrefund
//     WaitingForCompletePrefund --> WaitingForMyTurnToFund : PreFundReceived
//     WaitingForCompletePrefund --> WaitingForCompletePrefund : PreFundReceived
//     WaitingForMyTurnToFund --> WaitingForCompleteFunding : FundingUpdated
//     WaitingForMyTurnToFund --> WaitingForMyTurnToFund : FundingUpdated
//     WaitingForCompleteFunding --> WaitingForCompletePostFund : FundingUpdated
//     WaitingForCompleteFunding --> WaitingForCompleteFunding : FundingUpdated
//     WaitingForCompletePostFund --> [*] : PostFundReceived
