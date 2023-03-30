package types

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestIsExternal(t *testing.T) {
	external := Destination(common.HexToHash("0x00000000000000000000000096f7123E3A80C9813eF50213ADEd0e4511CB820f"))
	internal := Destination(common.HexToHash("0x6f7123E3A80C9813eF50213A96f7123E3A80C9813eF50213ADEd0e4511CB820f"))

	if !external.IsExternal() {
		t.Fatalf("Received bytes %x was declared internal, when it is external", external)
	}

	if internal.IsExternal() {
		t.Fatalf("Received bytes %x was declared external, when it is internal", internal)
	}
}

var referenceAddress = []Address{
	common.HexToAddress(`0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa`),
	common.HexToAddress(`0x0000000000000000000000000000000000000000`),
	common.HexToAddress(`0x96f7123E3A80C9813eF50213ADEd0e4511CB820f`),
}

var areExternal = []Destination{
	Destination(common.HexToHash("0x000000000000000000000000aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")),
	Destination(common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000")),
	Destination(common.HexToHash("0x00000000000000000000000096f7123E3A80C9813eF50213ADEd0e4511CB820f")),
}

func TestToAddress(t *testing.T) {
	for i, extAddress := range areExternal {
		convertedAddress, err := extAddress.ToAddress()
		if err != nil {
			t.Fatalf("expected to convert %x to an external address, but failed", extAddress)
		}

		if convertedAddress != referenceAddress[i] {
			t.Fatalf("expected %x to convert to %x, but it did not", extAddress, referenceAddress[i])
		}
	}

	areNotExternal := []Destination{
		Destination(common.HexToHash("0x000000000000000000000001aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")),
		Destination(common.HexToHash("0x100000000000000000000000aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")),
		Destination(common.HexToHash("0x0000000000b0000000000000aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")),
	}

	for _, notExtAddress := range areNotExternal {
		if _, err := notExtAddress.ToAddress(); err == nil {
			t.Fatalf("expected to fail when converting %x to an external address, but succeeded", notExtAddress)
		}
	}
}

func TestToDestination(t *testing.T) {
	for i, refAddress := range referenceAddress {
		convertedAddress := AddressToDestination(refAddress)

		if convertedAddress != areExternal[i] {
			t.Fatalf("expected %x to convert to %x, but it did not", refAddress, areExternal[i])
		}
	}
}
