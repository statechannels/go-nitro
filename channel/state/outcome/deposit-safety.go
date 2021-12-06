package outcome

import (
	"math/big"

	"github.com/statechannels/go-nitro/types"
)

// DepositSafetyThreshold returns the amount of this asset that a user with
// the specified interest must see on-chain before the safe recoverability of
// their own deposits is guaranteed
func (s SingleAssetExit) DepositSafetyThreshold(interest types.Destination) *big.Int {
	sum := big.NewInt(0)

	for _, allocation := range s.Allocations {

		if allocation.Destination == interest {
			// we have 'hit' the destination whose balances we are interested in protecting
			return sum
		}

		sum.Add(sum, allocation.Amount)
	}

	return sum
}

// DepositSafetyThreshold returns the Funds that a user with the specified
// interest must see on-chain before the safe recoverability of their
// deposits is guaranteed
func (e Exit) DepositSafetyThreshold(interest types.Destination) types.Funds {
	threshold := types.Funds{}

	for _, assetExit := range e {
		threshold[assetExit.Asset] = assetExit.DepositSafetyThreshold(interest)
	}

	return threshold
}
