package outcome

import (
	"bytes"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/statechannels/go-nitro/types"
)

// Allocation declares an Amount to be paid to a Destination.
type Allocation struct {
	Destination    types.Bytes32  // Either an ethereum address or an application-specific identifier
	Amount         *types.Uint256 // An amount of a particular asset
	AllocationType uint8          // Directs calling code on how to interpret the allocation
	Metadata       []byte         // Custom metadata (optional field, can be zero bytes). This can be used flexibly by different protocols.
}

// Equals returns true if the supplied Allocation matches the receiver Allocation, and false otherwise.
// Fields are compared with ==, except for big.Ints which are compared using Cmp
func (a Allocation) Equals(b Allocation) bool {
	return a.Destination == b.Destination && a.AllocationType == b.AllocationType && a.Amount.Cmp(b.Amount) == 0 && bytes.Equal(a.Metadata, b.Metadata)
}

// Allocations is an array of type Allocation
type Allocations []Allocation

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

// SingleAssetExit declares an ordered list of Allocations for a single asset.
type SingleAssetExit struct {
	Asset       types.Address // Either the zero address (implying the native token) or the address of an ERC20 contract
	Metadata    []byte        // Can be used to encode arbitrary additional information that applies to all allocations.
	Allocations Allocations
}

// Exit is an ordered list of SingleAssetExits
type Exit []SingleAssetExit

func (a Exit) Equals(b Exit) bool {
	if len(a) != len(b) {
		return false
	}
	for i, saeA := range a {
		saeB := b[i]
		if !bytes.Equal(saeA.Metadata, saeB.Metadata) ||
			saeA.Asset != saeB.Asset ||
			!saeA.Allocations.Equals(saeB.Allocations) {
			return false
		}
	}
	return true
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

// Encode returns the abi encoded Exit
func (e *Exit) Encode() (types.Bytes, error) {
	return abi.Arguments{{Type: exitTy}}.Pack(e)
}

type rawAllocationsType []struct {
	Destination    [32]uint8 "json:\"destination\""
	Amount         *big.Int  "json:\"amount\""
	AllocationType uint8     "json:\"allocationType\""
	Metadata       []uint8   "json:\"metadata\""
}

type rawExitType []struct {
	Asset       common.Address     "json:\"asset\""
	Metadata    []uint8            "json:\"metadata\""
	Allocations rawAllocationsType "json:\"Allocations\""
}

func convertToAllocations(r rawAllocationsType) Allocations {
	allocation := make(Allocations, len(r))
	for i, a := range r {
		allocation[i] = Allocation{Destination: a.Destination, Amount: a.Amount, Metadata: a.Metadata, AllocationType: a.AllocationType}
	}
	return allocation
}

func convertToExit(r rawExitType) Exit {
	exit := make(Exit, len(r))

	for i, sae := range r {
		exit[i] = SingleAssetExit{Asset: sae.Asset, Metadata: sae.Metadata, Allocations: convertToAllocations(sae.Allocations)}
	}
	return exit
}

// Decode returns an Exit from an abi encoding
func Decode(data types.Bytes) (Exit, error) {
	unpacked, _ := abi.Arguments{{Type: exitTy}}.Unpack(data)
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
