package outcome

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/types"
)

var (
	guaranteeMetadata           = GuaranteeMetadata{Left: types.AddressToDestination(common.HexToAddress("0x0a")), Right: types.AddressToDestination(common.HexToAddress("0x0b"))}
	encodedGuaranteeMetadata, _ = hex.DecodeString("000000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000000b")
)

func TestGuaranteeMetadataEncode(t *testing.T) {
	encodedG, err := guaranteeMetadata.Encode()
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(encodedG, encodedGuaranteeMetadata) {
		t.Fatalf("incorrect encoding. Got %x, wanted %x", encodedG, encodedGuaranteeMetadata)
	}
}

func TestGuaranteeMetadataDecode(t *testing.T) {
	g, err := DecodeIntoGuaranteeMetadata(encodedGuaranteeMetadata)
	if err != nil {
		t.Error(err)
	}
	if g != guaranteeMetadata {
		t.Fatalf("incorrect encoding. Got %x, wanted %x", g, encodedGuaranteeMetadata)
	}
}
