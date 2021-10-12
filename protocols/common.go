package protocols

import (
	"math/big"

	"github.com/statechannels/go-nitro/types"
)

// A linear state machine with enumerated states.
// PreFundIncomplete => NotYetMyTurnToFund => FundingIncomplete => PostFundIncomplete => Finished
type DirectFundingEnumerableState int

const (
	WaitingForCompletePrefund DirectFundingEnumerableState = iota // 0
	WaitingForMyTurnToFund
	WaitingForCompleteFunding
	WaitingForCompletePostFund
	WaitingForNothing // Finished
)

// Effects to be declared. For now these are just strings. In future they may be more complex and type safe
type SideEffects []string

var NoSideEffects = []string{}

var zero = big.NewInt(0)

// ObjectiveState is a cache of data computed by reading from the store. It stores (potentially) infinite data
type ObjectiveState struct {
	ParticipantIndex map[types.Address]uint // the index for each participant

	PreFundSigned []bool // indexed by participant. TODO should this be initialized with my own index showing true?

	MyDepositSafetyThreshold *big.Int // if the on chain holdings are equal to this amount it is safe for me to deposit
	MyDepositTarget          *big.Int // I want to get the on chain holdings up to this much
	FullyFundedThreshold     *big.Int // if the on chain holdings are equal

	PostFundSigned []bool // indexed by participant
}

// Methods on the ObjectiveState

// PrefundComplete returns true if all participants have signed a prefund state, as reflected by the extended state
func (s ObjectiveState) PrefundComplete() bool {
	for _, index := range s.ParticipantIndex {
		if !s.PreFundSigned[index] {
			return false
		}
	}
	return true
}

// PostfundComplete returns true if all participants have signed a postfund state, as reflected by the extended state
func (s ObjectiveState) PostfundComplete() bool {
	for _, index := range s.ParticipantIndex {
		if !s.PostFundSigned[index] {
			return false
		}
	}
	return true
}
