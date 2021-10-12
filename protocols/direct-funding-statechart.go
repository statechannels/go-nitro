package protocols

import (
	"errors"
	"math/big"

	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/types"
)

// DirectFundingExtendedState contains the (potentially infinite) extended state of the Direct Funding machine.
// The extended state of this machine is a cache of a larger store of states and events stored by a Nitro wallet.
// This struct should be kept as shallow copyable (and this should be tested).
type DirectFundingExtendedState struct {
	ParticipantIndex map[types.Address]uint // the index for each participant
	PreFundSigned    []bool                 // indexed by participant. TODO should this be initialized with my own index showing true?

	MyDepositSafetyThreshold *big.Int // if the on chain holdings are equal to this amount it is safe for me to deposit
	MyDepositTarget          *big.Int // I want to get the on chain holdings up to this much
	FullyFundedThreshold     *big.Int // if the on chain holdings are equal

	PostFundSigned []bool // indexed by participant
}

// PrefundComplete returns true if all participants have signed a prefund state, as reflected by the extended state
func (s DirectFundingExtendedState) PrefundComplete() bool {
	for _, index := range s.ParticipantIndex {
		if !s.PreFundSigned[index] {
			return false
		}
	}
	return true
}

// PostfundComplete returns true if all participants have signed a postfund state, as reflected by the extended state
func (s DirectFundingExtendedState) PostfundComplete() bool {
	for _, index := range s.ParticipantIndex {
		if !s.PostFundSigned[index] {
			return false
		}
	}
	return true
}

// DirectFundingProtocol state has both enumerable and extended state components.
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

func PrepareDepositTransaction(amount *big.Int) string {
	// TODO a proper implementation
	return `deposit` + amount.String()
}

// TODO can reducers be abstracted into an interface?

// NextState is the overall reducer / state transition function for the DirectFundingProtocol
func (s DirectFundingProtocolState) NextState(e DirectFundingProtocolEvent) (DirectFundingProtocolState, SideEffects, error) {
	// it is better to switch on the state than on the event
	// https://dev.to/davidkpiano/you-don-t-need-a-library-for-state-machines-k7h
	switch s.EnumerableState {
	case WaitingForCompletePrefund:
		return s.nextStateFromWaitingForCompletePrefund(e)
	case WaitingForCompleteFunding:
		return s.nextStateFromFundingIncomplete(e)
	case WaitingForCompletePostFund:
		return s.nextStateFromPostfundIncomplete(e)
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
		return DirectFundingProtocolState{WaitingForCompleteFunding, newExtendedState}, NoSideEffects, nil
	} else {
		return DirectFundingProtocolState{WaitingForCompletePostFund, newExtendedState}, NoSideEffects, nil
	}
}

// nextStateFromFundingIncomplete is a component of the overall DirectFundingProtocol reducer
func (s DirectFundingProtocolState) nextStateFromFundingIncomplete(e DirectFundingProtocolEvent) (DirectFundingProtocolState, SideEffects, error) {
	if e.Type != FundingUpdated { // There's only one way out of this state
		return s, NoSideEffects, nil
	}

	if e.OnChainHolding.Cmp(s.ExtendedState.FullyFundedThreshold) > -1 {
		// We make can progess to the next enumerable state
		return DirectFundingProtocolState{WaitingForCompletePostFund, s.ExtendedState}, NoSideEffects, nil
	}

	// We aren't fully funded

	if e.OnChainHolding.Cmp(s.ExtendedState.MyDepositTarget) > -1 {
		// Don't need to do anything but wait.
		return s, NoSideEffects, nil
	}

	// We haven't yet hit my deposit target

	if e.OnChainHolding.Cmp(s.ExtendedState.MyDepositSafetyThreshold) > -1 {
		depositAmount := big.NewInt(0).Sub(s.ExtendedState.MyDepositTarget, e.OnChainHolding)
		// TODO declare a side effect to deposit depositAmount
		return s, SideEffects{PrepareDepositTransaction(depositAmount)}, nil
	}

	// It isn't yet safe for me to fund

	return s, NoSideEffects, nil

}

// nextStateFromPostfundIncomplete is a component of the overall DirectFundingProtocol reducer
func (s DirectFundingProtocolState) nextStateFromPostfundIncomplete(e DirectFundingProtocolEvent) (DirectFundingProtocolState, SideEffects, error) {
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
