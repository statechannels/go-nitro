package outcome

import (
	"bytes"
	"encoding/json"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
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

// Exit is an ordered list of SingleAssetExits
type Exit []SingleAssetExit

func (a Exit) Equal(b Exit) bool {
	if len(a) != len(b) {
		return false
	}
	for i, saeA := range a {
		saeB := b[i]
		if !bytes.Equal(saeA.Metadata, saeB.Metadata) ||
			saeA.Asset != saeB.Asset ||
			!saeA.Allocations.Equal(saeB.Allocations) {
			return false
		}
	}
	return true
}

// TotalAllocated returns the sum of all Funds that are allocated by the outcome.
//
// NOTE that these Funds are potentially different from a channel's capacity to
// pay out a given set of allocations, which is limited by the channel's holdings
func (e Exit) TotalAllocated() types.Funds {
	fullValue := types.Funds{}

	for _, assetExit := range e {
		fullValue[assetExit.Asset] = assetExit.TotalAllocated()
	}

	return fullValue
}

// TotalAllocatedFor returns the total amount allocated to the given dest (regardless of AllocationType)
func (e Exit) TotalAllocatedFor(dest types.Destination) types.Funds {
	total := types.Funds{}

	for _, assetAllocation := range e {
		total[assetAllocation.Asset] = assetAllocation.TotalAllocatedFor(dest)
	}

	return total
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

// exitTy describes the shape of Exit such that github.com/ethereum/go-ethereum/accounts/abi can parse it
var exitTy, _ = abi.NewType("tuple[]", "struct ExitFormat.SingleAssetExit[]", []abi.ArgumentMarshaling{
	{Name: "asset", Type: "address"},
	{Name: "metadata", Type: "bytes"},
	allocationsTy,
})

// rawAllocationsType is an alias to the type returned when using the github.com/ethereum/go-ethereum/accounts/abi Unpack method with allocationsTy
type rawAllocationsType = []struct {
	Destination    [32]byte `json:"destination"`
	Amount         *big.Int `json:"amount"`
	AllocationType uint8    `json:"allocationType"`
	Metadata       []uint8  `json:"metadata"`
}

// rawExitType is an alias to the type returned when using the github.com/ethereum/go-ethereum/accounts/abi Unpack method with exitTy
type rawExitType = []struct {
	Asset       common.Address     `json:"asset"`
	Metadata    []uint8            `json:"metadata"`
	Allocations rawAllocationsType `json:"Allocations"`
}

// convertToExit converts a rawExitType to an Exit
func convertToExit(r rawExitType) Exit {
	var exit Exit
	j, err := json.Marshal(r)

	if err != nil {
		log.Fatal(`error marshalling`)
	}

	err = json.Unmarshal(j, &exit)

	if err != nil {
		log.Fatal(`error unmarshalling`, err)
	}

	return exit
}

// Encode returns the abi encoded Exit
func (e *Exit) Encode() (types.Bytes, error) {
	return abi.Arguments{{Type: exitTy}}.Pack(e)
}

// Decode returns an Exit from an abi encoding
func Decode(data types.Bytes) (Exit, error) {
	unpacked, err := abi.Arguments{{Type: exitTy}}.Unpack(data)
	if err != nil {
		return nil, err
	}
	return convertToExit(unpacked[0].(rawExitType)), nil
}

// Hash returns the keccak256 hash of the Exit
func (e *Exit) Hash() (types.Bytes32, error) {
	if encoded, err := e.Encode(); err == nil {
		return crypto.Keccak256Hash(encoded), nil
	} else {
		return types.Bytes32{}, err
	}
}
