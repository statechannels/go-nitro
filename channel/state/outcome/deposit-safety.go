package outcome

import (
	"math/big"

	"github.com/statechannels/go-nitro/types"
)

// DepositSafetyThreshold returns the amount of this asset that a user with
// the specified interests must see on-chain before the safe recoverability of
// their own deposts is guaranteed
func (s SingleAssetExit) DepositSafetyThreshold(interests []types.Destination) *big.Int {
	sum := big.NewInt(0)

	for _, allocation := range s.Allocations {

		for _, interest := range interests {
			if allocation.Destination == interest {
				// we have 'hit' a destination whose balances we are interested in protecting
				return sum
			}
		}

		sum.Add(sum, allocation.Amount)
	}

	return sum
}

// DepositSafetyThreshold returns the Funds that a user with the specified
// interests must see on-chain before the safe recoverability of their
// deposits is guaranteed
func (e Exit) DepositSafetyThreshold(interests ...types.Destination) types.Funds {
	threshold := types.Funds{}

	for _, assetExit := range e {
		threshold[assetExit.Asset] = assetExit.DepositSafetyThreshold(interests)
	}

	return threshold
}
