package outcome

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/statechannels/go-nitro/types"
)

// A Guarantee is an Allocation with AllocationType == GuaranteeAllocationType and Metadata = encode(GuaranteeMetaData)
type GuaranteeMetadata struct {
	// The peer who plays the role of Alice (peer 0)
	Left types.Destination
	// The peer who plays the role of Bob (peer n+1, where n=len(intermediaries)
	Right types.Destination
}

// guaranteeMetadataTy describes the shape of GuaranteeMetadata, so that the abi encoder knows how to encode it
var guaranteeMetadataTy, _ = abi.NewType("tuple", "struct GuaranteeMetadata", []abi.ArgumentMarshaling{
	{Name: "Left", Type: "bytes32"},
	{Name: "Right", Type: "bytes32"},
})

// Encode returns the abi.encoded GuaranteeMetadata (suitable for packing in an Allocation.Metadata field)
func (m GuaranteeMetadata) Encode() (types.Bytes, error) {
	return abi.Arguments{{Type: guaranteeMetadataTy}}.Pack(m)
}

// rawGuaranteeMetadataType is an alias to the type returned when using the github.com/ethereum/go-ethereum/accounts/abi Unpack method with guaranteeMetadataTy
type rawGuaranteeMetadataType = struct {
	Left  [32]uint8 "json:\"Left\""
	Right [32]uint8 "json:\"Right\""
}

// convertToGuaranteeMetadata converts a rawGuaranteeMetadataType to a GuaranteeMetadata
func convertToGuaranteeMetadata(r rawGuaranteeMetadataType) GuaranteeMetadata {
	guaranteeMetadata := GuaranteeMetadata{Left: types.Destination(r.Left), Right: types.Destination(r.Right)}
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
