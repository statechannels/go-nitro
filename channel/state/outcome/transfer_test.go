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
		Destination:    types.Destination(common.HexToHash("0x0a")),
		Amount:         big.NewInt(2),
		AllocationType: 0,
		Metadata:       make(types.Bytes, 0)}}

	var expectedNewAllocations = Allocations{{ // [{Alice: 0}]
		Destination:    types.Destination(common.HexToHash("0x0a")),
		Amount:         big.NewInt(0),
		AllocationType: 0,
		Metadata:       make(types.Bytes, 0)}}

	var expectedExitAllocations = Allocations{{ // [{Alice: 2}]
		Destination:    types.Destination(common.HexToHash("0x0a")),
		Amount:         big.NewInt(2),
		AllocationType: 0,
		Metadata:       make(types.Bytes, 0)}}

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

// Run with go test -fuzz FuzzTransfer
func FuzzTransfer(f *testing.F) {

	var fuzzTarget = func(t *testing.T, initialHoldings uint, destination string, amount uint, metadata []byte) {
		initialAllocations := Allocations{{ // [{Alice: 2}]
			Destination:    types.Destination(common.HexToHash("destination")),
			Amount:         big.NewInt(int64(amount)),
			AllocationType: 0,
			Metadata:       metadata,
		}}

		// Simply fuzz the target without inspecing the return values
		// TODO test some basic invariant via the return values?
		got1, got2 := ComputeTransferEffectsAndInteractions(*big.NewInt(int64(initialHoldings)), initialAllocations, []uint{})
		if types.Gt(got1.Total(), initialAllocations.Total()) {
			t.Fatal("new allocations allocates more than initial allocations")
		}
		if types.Gt(got2.Total(), initialAllocations.Total()) {
			t.Fatal("exit allocations allocates more than initial allocations")
		}

	}

	f.Add(uint(7), "0x0a", uint(4), []byte{})
	f.Fuzz(fuzzTarget)
}
