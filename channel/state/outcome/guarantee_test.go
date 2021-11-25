package outcome

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

var guaranteeMetadata = GuaranteeMetadata{Left: common.HexToAddress("0x0a"), Right: common.HexToAddress("0x0b")}
var encodedGuaranteeMetadata, _ = hex.DecodeString("000000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000000b")

func TestGuaranteeMetadataEncode(t *testing.T) {
	encodedG, err := guaranteeMetadata.Encode()
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(encodedG, encodedGuaranteeMetadata) {
		t.Errorf("incorrect encoding. Got %x, wanted %x", encodedG, encodedGuaranteeMetadata)
	}
}

func TestGuaranteeMetadataDecode(t *testing.T) {
	g, err := DecodeIntoGuaranteeMetadata(encodedGuaranteeMetadata)
	if err != nil {
		t.Error(err)
	}
	if !g.Equal(guaranteeMetadata) {
		t.Errorf("incorrect encoding. Got %x, wanted %x", g, encodedGuaranteeMetadata)
	}
}
