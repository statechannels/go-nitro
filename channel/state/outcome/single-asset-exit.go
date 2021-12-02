package outcome

import (
	"math/big"

	"github.com/statechannels/go-nitro/types"
)

// SingleAssetExit declares an ordered list of Allocations for a single asset.
type SingleAssetExit struct {
	Asset       types.Address // Either the zero address (implying the native token) or the address of an ERC20 contract
	Metadata    []byte        // Can be used to encode arbitrary additional information that applies to all allocations.
	Allocations Allocations
}

// TotalAllocated returns the toal amount allocated, summed across all destinations (regardless of AllocationType)
func (sae SingleAssetExit) TotalAllocated() *big.Int {
	return sae.Allocations.Total()
}

// TotalAllocatedFor returns the total amount allocated for the specific destination
func (sae SingleAssetExit) TotalAllocatedFor(dest types.Destination) *big.Int {
	return sae.Allocations.TotalFor(dest)
}
