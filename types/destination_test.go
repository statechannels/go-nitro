package types

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestIsExternal(t *testing.T) {

	external := Destination{common.HexToHash("0x00000000000000000000000096f7123E3A80C9813eF50213ADEd0e4511CB820f")}
	internal := Destination{common.HexToHash("0x6f7123E3A80C9813eF50213A96f7123E3A80C9813eF50213ADEd0e4511CB820f")}

	if !external.IsExternal() {
		t.Errorf("Received bytes %x was declared internal, when it is external", external)
	}

	if internal.IsExternal() {
		t.Errorf("Received bytes %x was declared external, when it is internal", internal)
	}

}

func TestToAddress(t *testing.T) {
	referenceAddress := []Address{
		common.HexToAddress(`0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa`),
		common.HexToAddress(`0x0000000000000000000000000000000000000000`),
		common.HexToAddress(`0x96f7123E3A80C9813eF50213ADEd0e4511CB820f`),
	}
	areExternal := []Destination{
		{common.HexToHash("0x000000000000000000000000aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")},
		{common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000")},
		{common.HexToHash("0x00000000000000000000000096f7123E3A80C9813eF50213ADEd0e4511CB820f")},
	}

	areNotExternal := []Destination{
		{common.HexToHash("0x000000000000000000000001aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")},
		{common.HexToHash("0x100000000000000000000000aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")},
		{common.HexToHash("0x0000000000b0000000000000aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")},
	}

	for i, extAddress := range areExternal {
		convertedAddress, err := extAddress.ToAddress()
		if err != nil {
			t.Errorf("expected to convert %x to an external address, but failed", extAddress)
		}

		if convertedAddress != referenceAddress[i] {
			t.Errorf("expected %x to convert to %x, but it did not", extAddress, referenceAddress[i])
		}
	}

	for _, notExtAddress := range areNotExternal {
		if _, err := notExtAddress.ToAddress(); err == nil {
			t.Errorf("expected to fail when converting %x to an external address, but succeeded", notExtAddress)
		}
	}
}
