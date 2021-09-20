package outcome

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/types"
)

func TestComputeTransferEffectsAndInteractions(t *testing.T) {

	initialHoldings := *big.NewInt(100)

	var initialAllocations = Allocations{{ // [{Alice: 2}]
		Destination:    common.HexToHash("0x0a"),
		Amount:         big.NewInt(2),
		AllocationType: 0,
		Metadata:       make(types.Bytes, 0)}}

	var expectedNewAllocations = Allocations{{ // [{Alice: 0}]
		Destination:    common.HexToHash("0x0a"),
		Amount:         big.NewInt(0),
		AllocationType: 0,
		Metadata:       make(types.Bytes, 0)}}

	var expectedExitAllocations = Allocations{{ // [{Alice: 2}]
		Destination:    common.HexToHash("0x0a"),
		Amount:         big.NewInt(2),
		AllocationType: 0,
		Metadata:       make(types.Bytes, 0)}}

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
