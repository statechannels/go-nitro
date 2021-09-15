package outcome

import (
	"math/big"
	"testing"
)

func TestTransfer(t *testing.T) {

	initialHoldings := *big.NewInt(100)

	var initialAllocations = Allocations{{ // [{Alice: 2}]
		Destination:    "0x000000000000000000000000000000000000000000000000000000000000000a",
		Amount:         *big.NewInt(2),
		AllocationType: 0,
		Metadata:       "0x"}}

	var expectedNewAllocations = Allocations{{ // [{Alice: 0}]
		Destination:    "0x000000000000000000000000000000000000000000000000000000000000000a",
		Amount:         *big.NewInt(0),
		AllocationType: 0,
		Metadata:       "0x"}}

	var expectedExitAllocations = Allocations{{ // [{Alice: 2}]
		Destination:    "0x000000000000000000000000000000000000000000000000000000000000000a",
		Amount:         *big.NewInt(2),
		AllocationType: 0,
		Metadata:       "0x"}}

	got1, got2 := ComputeTransferEffectsAndInteractions(initialHoldings, initialAllocations, []uint{})
	want1 := expectedNewAllocations
	want2 := expectedExitAllocations

	if !got1.Equals(want1) {
		t.Errorf("got %+v, wanted %+v", got1, want1)
	}

	if !got2.Equals(want2) {
		t.Errorf("got %+v, wanted %+v", got2, want2)
	}

}
