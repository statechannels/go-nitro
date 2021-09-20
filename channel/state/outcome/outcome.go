package outcome

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
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
	return a.Destination == b.Destination && a.AllocationType == b.AllocationType && a.Amount.Cmp(b.Amount) == 0
	// TODO a.Metadata.Equals(b.Metadata)
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

// Encode returns the rlp encoding of the Exit
func (e *Exit) Encode() (types.Bytes, error) {

	allocations := abi.ArgumentMarshaling{
		Name:         "Allocations",
		Type:         "tuple[]",
		InternalType: "struct ExitFormat.Allocation[]",
		Components: []abi.ArgumentMarshaling{
			{Name: "destination", Type: "bytes32", InternalType: "bytes32"},
			{Name: "amount", Type: "uint256", InternalType: "address"},
			{Name: "allocationType", Type: "uint8", InternalType: "address"},
			{Name: "metadata", Type: "bytes", InternalType: "address"},
		},
		Indexed: false,
	}
	exitTy, _ := abi.NewType("tuple[]", "struct ExitFormat.SingleAssetExit[]", []abi.ArgumentMarshaling{
		{Name: "asset", Type: "address", InternalType: "address"},
		{Name: "metadata", Type: "bytes", InternalType: "address"},
		allocations,
	})

	args2 := abi.Arguments{{Type: exitTy}}

	return args2.Pack(e)
}
