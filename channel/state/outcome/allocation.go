package outcome

import (
	"bytes"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/statechannels/go-nitro/types"
)

type AllocationType uint8

const (
	NormalAllocationType AllocationType = iota
	GuaranteeAllocationType
)

// Allocation declares an Amount to be paid to a Destination.
type Allocation struct {
	Destination    types.Destination // Either an ethereum address or an application-specific identifier
	Amount         *types.Uint256    // An amount of a particular asset
	AllocationType AllocationType    // Directs calling code on how to interpret the allocation
	Metadata       []byte            // Custom metadata (optional field, can be zero bytes). This can be used flexibly by different protocols.
}

// Equal returns true if the supplied Allocation matches the receiver Allocation, and false otherwise.
// Fields are compared with ==, except for big.Ints which are compared using Cmp
func (a Allocation) Equal(b Allocation) bool {
	return a.Destination == b.Destination && a.AllocationType == b.AllocationType && a.Amount.Cmp(b.Amount) == 0 && bytes.Equal(a.Metadata, b.Metadata)
}

// Clone returns a deep copy of the receiver.
func (a Allocation) Clone() Allocation {
	return Allocation{
		Destination:    a.Destination,
		Amount:         big.NewInt(0).Set(a.Amount),
		AllocationType: a.AllocationType,
		Metadata:       a.Metadata,
	}
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

// Clone returns a deep copy of the receiver
func (a Allocations) Clone() Allocations {
	clone := make(Allocations, len(a))
	for i, allocation := range a {
		clone[i] = allocation.Clone()
	}
	return clone
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
func (a Allocations) Affords(given Allocation, funding *big.Int) bool {
	bigZero := big.NewInt(0)
	surplus := big.NewInt(0).Set(funding)
	for _, allocation := range a {

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

// rawAllocationsType is an alias to the type returned when using the github.com/ethereum/go-ethereum/accounts/abi Unpack method with allocationsTy
type rawAllocationsType = []struct {
	Destination    [32]byte `json:"destination"`
	Amount         *big.Int `json:"amount"`
	AllocationType uint8    `json:"allocationType"`
	Metadata       []uint8  `json:"metadata"`
}

// allocationsTy describes the shape of Allocations such that github.com/ethereum/go-ethereum/accounts/abi can parse it
var allocationsTy = abi.ArgumentMarshaling{
	Name: "Allocations",
	Type: "tuple[]",
	Components: []abi.ArgumentMarshaling{
		{Name: "destination", Type: "bytes32"},
		{Name: "amount", Type: "uint256"},
		{Name: "allocationType", Type: "uint8"},
		{Name: "metadata", Type: "bytes"},
	},
}

// DivertToGuarantee returns a new Allocations, identical to the receiver but with
// the leftDestination's amount reduced by leftAmount,
// the rightDestination's amount reduced by rightAmount,
// and a Guarantee appended for the guaranteeDestination
func (a Allocations) DivertToGuarantee(
	leftDestination types.Destination,
	rightDestination types.Destination,
	leftAmount *big.Int,
	rightAmount *big.Int,
	guaranteeDestination types.Destination,
) (Allocations, error) {

	if leftDestination == rightDestination {
		return Allocations{}, errors.New(`debtees must be distinct`)
	}

	newAllocations := make([]Allocation, 0, len(a)+1)
	for i, allocation := range a {
		newAllocations = append(newAllocations, allocation.Clone())
		switch newAllocations[i].Destination {
		case leftDestination:
			newAllocations[i].Amount.Sub(newAllocations[i].Amount, leftAmount)
		case rightDestination:
			newAllocations[i].Amount.Sub(newAllocations[i].Amount, rightAmount)
		}
		if types.Gt(big.NewInt(0), newAllocations[i].Amount) {
			return Allocations{}, errors.New(`insufficient funds`)
		}
	}
	encodedGuaranteeMetadata, err := GuaranteeMetadata{
		Left:  leftDestination,
		Right: rightDestination,
	}.Encode()

	if err != nil {
		return Allocations{}, errors.New(`error encoding guarantee`)
	}

	newAllocations = append(newAllocations, Allocation{
		Destination:    guaranteeDestination,
		Amount:         big.NewInt(0).Add(leftAmount, rightAmount),
		AllocationType: GuaranteeAllocationType,
		Metadata:       encodedGuaranteeMetadata,
	})

	return newAllocations, nil
}
