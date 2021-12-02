package outcome

import (
	"bytes"
	"math/big"

	"github.com/statechannels/go-nitro/types"
)

// Allocation declares an Amount to be paid to a Destination.
type Allocation struct {
	Destination    types.Destination // Either an ethereum address or an application-specific identifier
	Amount         *types.Uint256    // An amount of a particular asset
	AllocationType uint8             // Directs calling code on how to interpret the allocation
	Metadata       []byte            // Custom metadata (optional field, can be zero bytes). This can be used flexibly by different protocols.
}

// TODO AllocationType should be an enum?
const NormalAllocationType = uint8(0)
const GuaranteeAllocationType = uint8(1)

// Equal returns true if the supplied Allocation matches the receiver Allocation, and false otherwise.
// Fields are compared with ==, except for big.Ints which are compared using Cmp
func (a Allocation) Equal(b Allocation) bool {
	return a.Destination == b.Destination && a.AllocationType == b.AllocationType && a.Amount.Cmp(b.Amount) == 0 && bytes.Equal(a.Metadata, b.Metadata)
}

// Allocations is an array of type Allocation
type Allocations []Allocation

// Equal returns true if each of the supplied Allocations matches the receiver Allocation in the same position, and false otherwise.
func (a Allocations) Equal(b Allocations) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if !a[i].Equal(b[i]) {
			return false
		}
	}
	return true
}

// Total returns the toal amount allocated, summed across all destinations (regardless of AllocationType)
func (a Allocations) Total() *big.Int {
	total := big.NewInt(0)
	for _, allocation := range a {
		total.Add(total, allocation.Amount)
	}
	return total
}

// TotalFor returns the total amount allocated to the given dest (regardless of AllocationType)
func (a Allocations) TotalFor(dest types.Destination) *big.Int {
	total := big.NewInt(0)
	for _, allocation := range a {
		if allocation.Destination == dest {
			total.Add(total, allocation.Amount)
		}
	}
	return total
}

// Affords returns true if the allocations can afford the given allocation given the input funding, false otherwise.
//
// To afford the given allocation, the allocations must include something equal-in-value to it,
// as well as having sufficient funds left over for it after reserving funds from the input funding for all allocations with higher priority.
// Note that "equal-in-value" implies the same allocation type and metadata (if any).
func (allocations Allocations) Affords(given Allocation, funding *big.Int) bool {
	bigZero := big.NewInt(0)
	surplus := big.NewInt(0).Set(funding)
	for _, allocation := range allocations {

		if allocation.Equal(given) {
			return surplus.Cmp(given.Amount) >= 0
		}

		surplus.Sub(surplus, allocation.Amount)

		if surplus.Cmp(bigZero) != 1 {
			break // no funds remain for further allocations
		}

	}
	return false
}
