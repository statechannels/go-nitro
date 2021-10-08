package channel

import (
	"math/big"

	"github.com/statechannels/go-nitro/channel/state"
)

type Channel interface {
	Id() string
	LatestSupportedState() state.State

	// Prefunding
	ShouldSignPreFund() bool // true IFF it is safe to do so AND I have not already done so
	IsPrefundComplete() bool

	// funding
	AmountToDeposit(asset string, currentAmountHeld big.Int) (big.Int, bool) // the bool is true if it safe to deposit
	IsFundingComplete(asset string, currentAmountHeld big.Int) bool

	// Postfunding
	ShouldSignPostFund() bool // true IFF it is safe to do so AND I have not already done so
	IsPostFundComplete() bool
}
