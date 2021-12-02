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
// interests must see on-chain before the safe recoverability of their
// deposits is guaranteed
func (e Exit) DepositSafetyThreshold(interests ...types.Destination) types.Funds {
	threshold := types.Funds{}

	for _, assetExit := range e {
		threshold[assetExit.Asset] = assetExit.DepositSafetyThreshold(interests)
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
