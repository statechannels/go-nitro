package outcome

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/statechannels/go-nitro/types"
)

// An HTLC is an Allocation with AllocationType == HTLCAllocationType and Metadata = encode(HTLCMetaData)
type HTLCMetadata struct {
	// The peer who has made the payment
	Payer types.Address
	// The peer who will receive the payment
	Payee types.Address

	// The block number at which the HTLC will expire
	ExpirationBlock uint64

	// The hash whose preimage will unlock the HTLC
	Hash types.Bytes32
}

// guaranteeMetadataTy describes the shape of HTLCMetadata, so that the abi encoder knows how to encode it
var htlcMetadataTy, _ = abi.NewType("tuple", "struct GuaranteeMetadata", []abi.ArgumentMarshaling{
	{Name: "Payer", Type: "address"},
	{Name: "Right", Type: "address"},
	{Name: "ExpirationBlock", Type: "uint64"},
	{Name: "Hash", Type: "bytes32"},
})

// Encode returns the abi.encoded HTLCMetadata (suitable for packing in an Allocation.Metadata field)
func (h HTLCMetadata) Encode() (types.Bytes, error) {
	return abi.Arguments{{Type: htlcMetadataTy}}.Pack(h)
}

// rawHTLCMetadataType is an alias to the type returned when using the github.com/ethereum/go-ethereum/accounts/abi Unpack method with guaranteeMetadataTy
type rawHTLCMetadataType = struct {
	Payer           [20]uint8
	Payee           [20]uint8
	ExpirationBlock uint64
	Hash            [32]uint8
}

// convertToHTLCMetadata converts a rawGuaranteeMetadataType to a GuaranteeMetadata
func convertToHTLCMetadata(r rawHTLCMetadataType) HTLCMetadata {
	htlcMetadata := HTLCMetadata{
		Payer:           r.Payer,
		Payee:           r.Payee,
		ExpirationBlock: r.ExpirationBlock,
		Hash:            r.Hash,
	}
	return htlcMetadata
}

// Decode returns a GuaranteeMetaData from an abi encoding
func DecodeIntoHTLCMetadata(m []byte) (HTLCMetadata, error) {
	unpacked, err := abi.Arguments{{Type: htlcMetadataTy}}.Unpack(m)
	if err != nil {
		return HTLCMetadata{}, err
	}
	return convertToHTLCMetadata(unpacked[0].(rawHTLCMetadataType)), nil
}
