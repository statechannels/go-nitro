package protocols

import (
	"errors"
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

func SignPreFundEffect(cId types.Bytes32) string {
	return "sign Prefundsetup for" + cId.String()
}
func SignPostFundEffect(cId types.Bytes32) string {
	return "sign Postfundsetup for" + cId.String()
}
func FundOnChainEffect(cId types.Bytes32, asset string, amount *big.Int) string {
	return "deposit" + amount.Text(64) + "into" + cId.String()
}

var NoSideEffects = []string{}

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

	MyIndex uint // my participant index

	PreFundSigned []bool // indexed by participant. TODO should this be initialized with my own index showing true?

	MyDepositSafetyThreshold *big.Int // if the on chain holdings are equal to this amount it is safe for me to deposit
	MyDepositTarget          *big.Int // I want to get the on chain holdings up to this much
	FullyFundedThreshold     *big.Int // if the on chain holdings are equal

	PostFundSigned []bool // indexed by participant

	OnChainHolding *big.Int
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

// FundingComplete returns true if the supplied onChainHolding is greater than or equal to the threshold for being fully funded.
func (s DirectFundingObjectiveState) FundingComplete(onChainHolding *big.Int) bool {
	return gte(onChainHolding, s.FullyFundedThreshold)
}

// SafeToDeposit returns true if the supplied onChainHolding is greater than or equal to the threshold for safety.
func (s DirectFundingObjectiveState) SafeToDeposit(onChainHolding *big.Int) bool {
	return gte(onChainHolding, s.FullyFundedThreshold)
}

// AmountToDeposit computes the appropriate amount to deposit using the supplied onChainHolding
func (s DirectFundingObjectiveState) AmountToDeposit(onChainHolding *big.Int) *big.Int {
	return big.NewInt(0).Sub(s.MyDepositTarget, onChainHolding)
}

func gte(a *big.Int, b *big.Int) bool {
	return a.Cmp(b) > -1
}

func gt(a *big.Int, b *big.Int) bool {
	return a.Cmp(b) > 0
}

// errors
var ErrNotApproved = errors.New("objective not approved")
