package protocols

import (
	"errors"
	"math/big"

	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/types"
)

// A linear state machine
// PreFundIncomplete => NotYetMyTurnToFund => FundingIncomplete => PostFundIncomplete
type DirectFundingEnumerableState int

const (
	PreFundIncomplete DirectFundingEnumerableState = iota // 0
	NotYetMyTurnToFund
	FundingIncomplete
	PostFundIncomplete
)

// This should be shallow copyable
type DirectFundingExtendedState struct {
	ParticipantIndex map[types.Address]uint // the index for each participant
	PreFundSigned    map[uint]bool          // indexed by participant
	DepositCovered   map[uint]*big.Int      // indexed by participant
	PostFundSigned   map[uint]bool          // indexed by participant
}

func (s DirectFundingExtendedState) PrefundComplete() bool {
	for _, index := range s.ParticipantIndex {
		if !s.PreFundSigned[index] {
			return false
		}
	}
	return true
}

type DirectFundingProtocolState struct {
	EnumerableState DirectFundingEnumerableState
	ExtendedState   DirectFundingExtendedState
}

// The events for the state machine
type DirectFundingProtocolEventType int

const (
	PreFundReceived DirectFundingProtocolEventType = iota
)

type DirectFundingProtocolEvent struct {
	Type           DirectFundingProtocolEventType
	State          state.State
	Signature      state.Signature
	OnChainHolding *big.Int
}

// TODO we need to also have a context to turn this into a state chart

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
