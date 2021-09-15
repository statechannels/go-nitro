package outcome

import (
	"math/big"
)

// types

// Allocation declares an Amount to be paid to a Destination.
type Allocation struct {
	Destination    string  // Either an ethereum address or an application-specific identifier
	Amount         big.Int // An amount of a particular asset
	AllocationType uint    // Directs calling code on how to interpret the allocation
	Metadata       string  // Custom metadata (optional field, can be zero bytes). This can be used flexibly by different protocols.
}

// Allocations is an array of type Allocation
type Allocations []Allocation

// SingleAssetExit declares an ordered list of Allocations for a single asset.
type SingleAssetExit struct {
	Asset       string // Either the zero address (implying the native token) or the address of an ERC20 contract
	Metadata    string // Can be used to encode arbitrary additional information that applies to all allocations.
	Allocations Allocations
}

// Exit is an ordered list of SingleAssetExits
type Exit []SingleAssetExit

// methods

// Equals returns true if the supplied Allocation matches the receiver Allocation, and false otherwise.
// Fields are compared with ==, except for big.Ints which are compared using Cmp
func (a Allocation) Equals(b Allocation) bool {
	return a.Destination == b.Destination && a.AllocationType == b.AllocationType && a.Metadata == b.Metadata && a.Amount.Cmp(&b.Amount) == 0
}

// Equals returns true if each of the supplied Allocations matches the receiver Allocation in the same position, and false otherwise.
func (a Allocations) Equals(b Allocations) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if !a[i].Equals(b[i]) {
			return false
		}
	}
	return true
}
