package protocols

import "math/big"

// A linear state machine with enumerated states.
// PreFundIncomplete => NotYetMyTurnToFund => FundingIncomplete => PostFundIncomplete
type DirectFundingEnumerableState int

const (
	PreFundIncomplete DirectFundingEnumerableState = iota // 0
	FundingIncomplete
	PostFundIncomplete
	Finished
)

// Effects to be declared. For now these are just strings. In future they may be more complex
type SideEffects []string

var NoSideEffects = []string{}

var zero = big.NewInt(0)
