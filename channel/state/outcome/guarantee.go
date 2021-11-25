package outcome

import (
	"bytes"
	"encoding/json"
	"log"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/types"
)

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
