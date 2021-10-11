package protocols

import (
	"errors"
	"math/big"

	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/types"
)

// A linear state machine with enumerated states.
// PreFundIncomplete => NotYetMyTurnToFund => FundingIncomplete => PostFundIncomplete
type DirectFundingEnumerableState int

const (
	PreFundIncomplete DirectFundingEnumerableState = iota // 0
	NotYetMyTurnToFund
	FundingIncomplete
	PostFundIncomplete
)

// DirectFundingExtendedState contains the (potentially infinite) extended state of the Direct Funding machine.
// The extended state of this machine is a cache of a larger store of states and events stored by a Nitro wallet.
// This struct should be kept as shallow copyable (and this should be tested).
type DirectFundingExtendedState struct {
	ParticipantIndex map[types.Address]uint // the index for each participant
	PreFundSigned    map[uint]bool          // indexed by participant
	DepositCovered   map[uint]*big.Int      // indexed by participant
	PostFundSigned   map[uint]bool          // indexed by participant
}

// PrefunComplete returns true if all participants have signed a prefund state, as reflected by the extended state
func (s DirectFundingExtendedState) PrefundComplete() bool {
	for _, index := range s.ParticipantIndex {
		if !s.PreFundSigned[index] {
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
)

// DirectFundingProtocolEvent has a type as well as other rich information (which may or may not be non nil).
type DirectFundingProtocolEvent struct {
	Type           DirectFundingProtocolEventType
	State          state.State
	Signature      state.Signature
	OnChainHolding *big.Int
}

// NextState is the overall reducer / state transition function for the DirectFundingProtocol
func (s DirectFundingProtocolState) NextState(e DirectFundingProtocolEvent) (DirectFundingProtocolState, error) {
	// it is better to switch on the state than on the event
	// https://dev.to/davidkpiano/you-don-t-need-a-library-for-state-machines-k7h
	switch s.EnumerableState {
	case PreFundIncomplete:
		return s.nextStateFromPrefundIncomplete(e)
	case NotYetMyTurnToFund:
		fallthrough // TODO
	case FundingIncomplete:
		fallthrough // TODO
	case PostFundIncomplete:
		fallthrough // TODO
	default:
		return s, nil
	}
}

// nextStateFromPrefundIncomplete is a component of the overall DirectFundingProtocol reducer
func (s DirectFundingProtocolState) nextStateFromPrefundIncomplete(e DirectFundingProtocolEvent) (DirectFundingProtocolState, error) {
	if e.Type != PreFundReceived { // There's only one way out of this state
		return s, nil
	}
	newExtendedState := s.ExtendedState // Make a copy of the extended state

	signer, err := e.State.RecoverSigner(e.Signature)
	if err != nil {
		return s, err
	}

	signerIndex, present := newExtendedState.ParticipantIndex[signer]
	if !present {
		return s, errors.New(`signer is not a participant`)
	} else {
		newExtendedState.PreFundSigned[signerIndex] = true
	}

	if newExtendedState.PrefundComplete() {
		return DirectFundingProtocolState{NotYetMyTurnToFund, newExtendedState}, nil
	} else {
		return DirectFundingProtocolState{PostFundIncomplete, newExtendedState}, nil
	}

}
