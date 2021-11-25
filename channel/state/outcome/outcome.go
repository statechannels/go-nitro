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

// A Guarantee is an Allocation with AllocationType == GuaranteeAllocationType and Metadata = encode(GuaranteeMetaData)
type GuaranteeMetadata struct {
	Left  types.Address // The peer who plays the role of Alice (peer 0)
	Right types.Address // The peer who plays the role of Bob (peer n+1, where n=len(peers))
}

// guaranteeMetadataTy describes the shape of GuaranteeMetadata, so that the abi encoder knows how to encode it
var guaranteeMetadataTy, _ = abi.NewType("tuple", "struct GuaranteeMetadata", []abi.ArgumentMarshaling{
	{Name: "Left", Type: "address"},
	{Name: "Right", Type: "address"},
})

// Encode returns the abi.encoded GuaranteeMetadata (suitable for packing in an Allocation.Metadata field)
func (m GuaranteeMetadata) Encode() (types.Bytes, error) {
	return abi.Arguments{{Type: guaranteeMetadataTy}}.Pack(m)
}

// Equal returns true if the reciever has identically valued fields to the supplied object
func (m GuaranteeMetadata) Equal(n GuaranteeMetadata) bool {
	return bytes.Equal(m.Left.Bytes(), n.Left.Bytes()) && bytes.Equal(m.Right.Bytes(), n.Right.Bytes())

}

// rawGuaranteeMetadataType is an alias to the type returned when using the github.com/ethereum/go-ethereum/accounts/abi Unpack method with guaranteeMetadataTy
type rawGuaranteeMetadataType = struct {
	Left  common.Address "json:\"Left\""
	Right common.Address "json:\"Right\""
}

// convertToGuaranteeMetadata converts a rawGuaranteeMetadataType to a GuaranteeMetadata
func convertToGuaranteeMetadata(r rawGuaranteeMetadataType) GuaranteeMetadata {
	var guaranteeMetadata GuaranteeMetadata
	j, err := json.Marshal(r)

	if err != nil {
		log.Fatal(`error marshalling`)
	}

	err = json.Unmarshal(j, &guaranteeMetadata)

	if err != nil {
		log.Fatal(`error unmarshalling`, err)
	}

	return guaranteeMetadata
}

// Decode returns a GuaranteeMetaData from an abi encoding
func DecodeIntoGuaranteeMetadata(m []byte) (GuaranteeMetadata, error) {
	unpacked, err := abi.Arguments{{Type: guaranteeMetadataTy}}.Unpack(m)
	if err != nil {
		return GuaranteeMetadata{}, err
	}
	return convertToGuaranteeMetadata(unpacked[0].(rawGuaranteeMetadataType)), nil
}

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
