package outcome

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/types"
)

func TestComputeTransferEffectsAndInteractions(t *testing.T) {
	initialHoldings := *big.NewInt(100)

	initialAllocations := Allocations{{ // [{Alice: 2}]
		Destination:    types.Destination(common.HexToHash("0x0a")),
		Amount:         big.NewInt(2),
		AllocationType: 0,
		Metadata:       make(types.Bytes, 0),
	}}

	expectedNewAllocations := Allocations{{ // [{Alice: 0}]
		Destination:    types.Destination(common.HexToHash("0x0a")),
		Amount:         big.NewInt(0),
		AllocationType: 0,
		Metadata:       make(types.Bytes, 0),
	}}

	expectedExitAllocations := Allocations{{ // [{Alice: 2}]
		Destination:    types.Destination(common.HexToHash("0x0a")),
		Amount:         big.NewInt(2),
		AllocationType: 0,
		Metadata:       make(types.Bytes, 0),
	}}

	got1, got2 := ComputeTransferEffectsAndInteractions(initialHoldings, initialAllocations, []uint{})
	want1 := expectedNewAllocations
	want2 := expectedExitAllocations

	if !got1.Equal(want1) {
		t.Fatalf("got %+v, wanted %+v", got1, want1)
	}

	if !got2.Equal(want2) {
		t.Fatalf("got %+v, wanted %+v", got2, want2)
	}
}
