package protocols

import (
	"errors"
	"math/big"

	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/types"
)

// A linear state machine with enumerated states.
// PreFundIncomplete => NotYetMyTurnToFund => FundingIncomplete => PostFundIncomplete => Finished
type DirectFundingEnumerableState int

const (
	WaitingForCompletePrefund  = "WaitingForCompletePrefund"
	WaitingForMyTurnToFund     = "WaitingForMyTurnToFund"
	WaitingForCompleteFunding  = "WaitingForCompleteFunding"
	WaitingForCompletePostFund = "WaitingForCompletePostFund"
	WaitingForNothing          = "WaitingForNothing" // Finished
)

func SignPreFundEffect(cId types.Bytes32) string {
	return "sign Prefundsetup for" + cId.String()
}
func SignPostFundEffect(cId types.Bytes32) string {
	return "sign Postfundsetup for" + cId.String()
}
func FundOnChainEffect(cId types.Bytes32, asset string, amount types.Holdings) string {
	return "deposit" + amount.String() + "into" + cId.String()
}

var NoSideEffects = SideEffects{}

var zero = big.NewInt(0)

type ObjectiveStatus int8

const (
	Unapproved ObjectiveStatus = iota
	Approved
	Rejected
)

// DirectFundingObjectiveState is a cache of data computed by reading from the store. It stores (potentially) infinite data
type DirectFundingObjectiveState struct {
	Status    ObjectiveStatus
	ChannelId types.Bytes32

	ParticipantIndex map[types.Address]uint // the index for each participant
	ExpectedStates   []state.State          // indexed by turn number

	MyIndex uint // my participant index

	PreFundSigned []bool // indexed by participant. TODO should this be initialized with my own index showing true?

	MyDepositSafetyThreshold types.Holdings // if the on chain holdings are equal to this amount it is safe for me to deposit
	MyDepositTarget          types.Holdings // I want to get the on chain holdings up to this much
	FullyFundedThreshold     types.Holdings // if the on chain holdings are equal

	PostFundSigned []bool // indexed by participant

	OnChainHolding types.Holdings
}

// Methods on the ObjectiveState

// PrefundComplete returns true if all participants have signed a prefund state, as reflected by the extended state
func (s DirectFundingObjectiveState) PrefundComplete() bool {
	for _, index := range s.ParticipantIndex {
		if !s.PreFundSigned[index] {
			return false
		}
	}
	return true
}

// PostfundComplete returns true if all participants have signed a postfund state, as reflected by the extended state
func (s DirectFundingObjectiveState) PostfundComplete() bool {
	for _, index := range s.ParticipantIndex {
		if !s.PostFundSigned[index] {
			return false
		}
	}
	return true
}

// FundingComplete returns true if the supplied onChainHoldings are greater than or equal to the threshold for being fully funded.
func (s DirectFundingObjectiveState) FundingComplete(onChainHoldings types.Holdings) bool {
	for asset, threshold := range s.FullyFundedThreshold {
		chainHolding, ok := onChainHoldings[asset]

		if !ok {
			return false
		}

		if gt(threshold, chainHolding) {
			return false
		}
	}

	return true
}

// SafeToDeposit returns true if the supplied onChainHoldings are greater than or equal to the threshold for safety.
func (s DirectFundingObjectiveState) SafeToDeposit(onChainHoldings types.Holdings) bool {
	for asset, safetyThreshold := range s.MyDepositSafetyThreshold {
		chainHolding, ok := onChainHoldings[asset]

		if !ok {
			return false
		}

		if gt(safetyThreshold, chainHolding) {
			return false
		}
	}

	return true
}

// AmountToDeposit computes the appropriate amount to deposit using the supplied onChainHolding
func (s DirectFundingObjectiveState) AmountToDeposit(onChainHoldings types.Holdings) types.Holdings {
	deposits := make(types.Holdings, len(onChainHoldings))

	for asset, holding := range onChainHoldings {
		deposits[asset] = big.NewInt(0).Sub(s.MyDepositTarget[asset], holding)
	}

	return deposits
}

func gte(a *big.Int, b *big.Int) bool {
	return a.Cmp(b) > -1
}

func gt(a *big.Int, b *big.Int) bool {
	return a.Cmp(b) > 0
}

// errors
var ErrNotApproved = errors.New("objective not approved")
