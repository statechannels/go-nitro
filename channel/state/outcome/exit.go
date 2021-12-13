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

// Clone returns a deep clone of the reciever.
func (e Exit) Clone() Exit {
	clone := make(Exit, len(e))
	for i, sae := range e {
		clone[i] = SingleAssetExit{
			Asset:       sae.Asset,
			Allocations: sae.Allocations.Clone(),
		}
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

// exitTy describes the shape of Exit such that github.com/ethereum/go-ethereum/accounts/abi can parse it
var exitTy, _ = abi.NewType("tuple[]", "struct ExitFormat.SingleAssetExit[]", []abi.ArgumentMarshaling{
	{Name: "asset", Type: "address"},
	{Name: "metadata", Type: "bytes"},
	allocationsTy,
})

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
