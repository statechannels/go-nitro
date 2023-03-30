package types

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/go-cmp/cmp"
)

func TestTextMarshalling(t *testing.T) {
	const destString = "0x000000000000000000000000f5a1bb5607c9d079e46d1b3dc33f257d937b43bd"
	dest := Destination(common.HexToHash("0x000000000000000000000000f5a1bb5607c9d079e46d1b3dc33f257d937b43bd"))

	// Test marshalling to text
	text, err := dest.MarshalText()
	if err != nil {
		t.Fatalf("failed to marshal destination to text: %+v", err)
	}
	if string(text) != destString {
		t.Fatalf("marshaled destination text %s does not match expected text %s", string(text), destString)
	}

	// Test unmarshalling from text
	var newDest Destination
	err = newDest.UnmarshalText([]byte(destString))
	if err != nil {
		t.Fatalf("failed to unmarshal destination to text: %+v", err)
	}
	if diff := cmp.Diff(newDest, dest); diff != "" {
		t.Fatalf("unmarshaled destination does not match expected destination:\n%s", diff)
	}
}
