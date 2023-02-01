package outcome

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/statechannels/go-nitro/types"
)

type AssetMetadata struct {
	AssetType uint8
	Metadata  []byte
}

// SingleAssetExit declares an ordered list of Allocations for a single asset.
type SingleAssetExit struct {
	Asset         types.Address // Either the zero address (implying the native token) or the address of an ERC20 contract
	AssetMetadata AssetMetadata // Can be used to encode arbitrary additional information that applies to all allocations.
	Allocations   Allocations
}

// Equal returns true if the supplied SingleAssetExit is deeply equal to the receiver.
func (s SingleAssetExit) Equal(r SingleAssetExit) bool {
	return bytes.Equal(s.AssetMetadata.Metadata, r.AssetMetadata.Metadata) &&
		s.Asset == r.Asset &&
		s.Allocations.Equal(r.Allocations)

}

// Clone returns a deep clone of the receiver.
func (s SingleAssetExit) Clone() SingleAssetExit {
	return SingleAssetExit{
		Asset:         s.Asset,
		AssetMetadata: AssetMetadata{s.AssetMetadata.AssetType, s.AssetMetadata.Metadata},
		Allocations:   s.Allocations.Clone(),
	}

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

// Equal returns true if the supplied Exit is deeply equal to the receiver.
func (a Exit) Equal(b Exit) bool {
	if len(a) != len(b) {
		return false
	}
	for i, saeA := range a {
		if !saeA.Equal(b[i]) {
			return false
		}
	}
	return true
}

// Clone returns a deep clone of the receiver.
func (e Exit) Clone() Exit {
	clone := make(Exit, len(e))
	for i, sae := range e {
		clone[i] = sae.Clone()
	}
	return clone
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

// assetMetadataTy describes the shape of AssetMetadata such that github.com/ethereum/go-ethereum/accounts/abi can parse it
var assetMetadataTy = abi.ArgumentMarshaling{
	Name: "AssetMetadata",
	Type: "tuple",
	Components: []abi.ArgumentMarshaling{
		{Name: "AssetType", Type: "uint8"},
		{Name: "Metadata", Type: "bytes"},
	},
}

// ExitTy describes the shape of Exit such that github.com/ethereum/go-ethereum/accounts/abi can parse it
var ExitTy, _ = abi.NewType("tuple[]", "struct ExitFormat.SingleAssetExit[]", []abi.ArgumentMarshaling{
	{Name: "asset", Type: "address"},
	assetMetadataTy,
	allocationsTy,
})

type rawAssetMetaDataType struct {
	AssetType uint8
	Metadata  []uint8
}

// rawExitType is an alias to the type returned when using the github.com/ethereum/go-ethereum/accounts/abi Unpack method with exitTy
type rawExitType = []struct {
	Asset         common.Address       `json:"asset"`
	AssetMetadata rawAssetMetaDataType `json:"AssetMetadata"`
	Allocations   rawAllocationsType   `json:"Allocations"`
}

// convertToExit converts a rawExitType to an Exit
func convertToExit(r rawExitType) Exit {
	exit := make(Exit, len(r))
	for i, raw := range r {
		exit[i] = SingleAssetExit{
			Asset:         raw.Asset,
			AssetMetadata: AssetMetadata{raw.AssetMetadata.AssetType, raw.AssetMetadata.Metadata},
			Allocations:   convertToAllocations(raw.Allocations),
		}
	}

	return exit
}

// convertToAllocations converts a rawAllocationsType to an Allocations
func convertToAllocations(r rawAllocationsType) Allocations {
	allocations := make(Allocations, len(r))
	for i, raw := range r {
		allocations[i] = Allocation{
			Destination:    raw.Destination,
			AllocationType: AllocationType(raw.AllocationType),
			Amount:         raw.Amount,
			Metadata:       raw.Metadata,
		}
	}

	return allocations
}

// Encode returns the abi encoded Exit
func (e *Exit) Encode() (types.Bytes, error) {
	return abi.Arguments{{Type: ExitTy}}.Pack(e)
}

// Decode returns an Exit from an abi encoding
func Decode(data types.Bytes) (Exit, error) {
	unpacked, err := abi.Arguments{{Type: ExitTy}}.Unpack(data)
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

// easyExit is a more ergonomic data type which can be derived from an Exit
type easyExit map[common.Address]SingleAssetExit

// toEasyExit() convets an Exit into an easyExit.
//
// An EasyExit is a mapping from asset to SingleAssetExit, rather than an array.
// The conversion loses some information, because the position in the original array is not recorded in the map.
// The position has no semantic meaning, but does of course affect the hash of the exit.
// Furthermore, this transformation assumes there are *no* repeated entries.
// For these reasons, the transformation should be considered non-invertibile and used with care.
func (e Exit) toEasyExit() easyExit {
	easy := make(easyExit)
	for i := range e {
		easy[e[i].Asset] = e[i]
	}
	return easy
}

// Affords returns true if every allocation in the allocationMap can be afforded by the Exit, given the funds
//
// Both arguments are maps keyed by the same assets
func (e Exit) Affords(
	allocationMap map[common.Address]Allocation,
	funding types.Funds) bool {
	easyExit := e.toEasyExit()
	for asset := range allocationMap {
		x := funding[asset]
		if x == nil {
			return false
		}
		allocation := allocationMap[asset]
		if !easyExit[asset].Allocations.Affords(allocation, x) {
			return false
		}
	}
	return true

}

// DivertToGuarantee returns a new Exit, identical to the receiver but with
// (for each asset of the Exit)
// the leftDestination's amount reduced by leftFunds[asset],
// the rightDestination's amount reduced by rightAmount[asset],
// and a Guarantee appended for the guaranteeDestination.
// Where an asset is missing from leftFunds or rightFunds, it is treated as if the corresponding amount is zero.
func (e Exit) DivertToGuarantee(
	leftDestination types.Destination,
	rightDestination types.Destination,
	leftFunds types.Funds,
	rightFunds types.Funds,
	guaranteeDestination types.Destination,
) (Exit, error) {

	f := e.Clone()

	leftFundsClone := leftFunds.Clone()
	rightFundsClone := rightFunds.Clone()

	for i, sae := range f {
		asset := sae.Asset

		leftAmount, leftOk := leftFundsClone[asset]
		if !leftOk {
			leftAmount = big.NewInt(0)
		}
		rightAmount, rightOk := rightFundsClone[asset]
		if !rightOk {
			rightAmount = big.NewInt(0)
		}

		newAllocations, err := sae.Allocations.DivertToGuarantee(leftDestination, rightDestination, leftAmount, rightAmount, guaranteeDestination)

		if err != nil {
			return Exit{}, fmt.Errorf("could not divert to guarantee: %w", err)
		}
		f[i].Allocations = newAllocations
	}

	return f, nil
}
