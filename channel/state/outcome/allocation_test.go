package outcome

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/go-cmp/cmp"
	"github.com/statechannels/go-nitro/types"
)

func TestEqualAllocations(t *testing.T) {
	a1 := Allocations{{ // [{Alice: 2}]
		Destination:    types.Destination(common.HexToHash("0x0a")),
		Amount:         big.NewInt(2),
		AllocationType: 0,
		Metadata:       make(types.Bytes, 0),
	}}

	a2 := Allocations{{ // [{Alice: 2}]
		Destination:    types.Destination(common.HexToHash("0x0a")),
		Amount:         big.NewInt(2),
		AllocationType: 0,
		Metadata:       make(types.Bytes, 0),
	}}

	if &a1 == &a2 {
		t.Fatalf("expected distinct pointers, but got identical pointers")
	}

	if !a1.Equal(a2) {
		t.Fatalf("expected equal Allocations, but got distinct Allocations")
	}
}

func TestAffords(t *testing.T) {
	testCases := map[string]struct {
		Allocations     Allocations
		GivenAllocation Allocation
		Funding         *big.Int
		Want            bool
	}{
		"case 0": {allocsX, allocsX[0], big.NewInt(3), true},
		"case 1": {allocsX, allocsX[0], big.NewInt(2), true},
		"case 2": {allocsX, allocsX[0], big.NewInt(1), false},
		"case 3": {allocsX, allocsX[1], big.NewInt(6), true},
		"case 4": {allocsX, allocsX[1], big.NewInt(5), true},
		"case 5": {allocsX, allocsX[1], big.NewInt(4), false},
		"case 6": {allocsX, allocsX[1], big.NewInt(2), false},
		"case 7": {allocsX, Allocation{}, big.NewInt(2), false},
	}

	for name, testcase := range testCases {
		t.Run(name, func(t *testing.T) {
			got := testcase.Allocations.Affords(testcase.GivenAllocation, testcase.Funding)
			if got != testcase.Want {
				t.Fatalf(
					`Incorrect AffordFor: expected %v.Affords(%v,%v) to be %v, but got %v`,
					testcase.Allocations, testcase.GivenAllocation, testcase.Funding, testcase.Want, got)
			}
		})
	}
}

func TestAllocationClone(t *testing.T) {
	a := Allocation{ // [{Alice: 2}]
		Destination:    types.Destination(common.HexToHash("0x0a")),
		Amount:         big.NewInt(2),
		AllocationType: 0,
		Metadata:       make(types.Bytes, 0),
	}

	clone := a.Clone()

	if diff := cmp.Diff(a, clone); diff != "" {
		t.Fatalf("Clone: mismatch (-want +got):\n%s", diff)
	}
}

func TestAllocationsClone(t *testing.T) {
	as := Allocations{
		{ // [{Alice: 2}]
			Destination:    types.Destination(common.HexToHash("0x0a")),
			Amount:         big.NewInt(2),
			AllocationType: 0,
			Metadata:       make(types.Bytes, 0),
		},
		{ // [{Bob: 3}]
			Destination:    types.Destination(common.HexToHash("0x0b")),
			Amount:         big.NewInt(3),
			AllocationType: 0,
			Metadata:       make(types.Bytes, 0),
		},
	}

	clone := as.Clone()

	if diff := cmp.Diff(as, clone); diff != "" {
		t.Fatalf("Clone: mismatch (-want +got):\n%s", diff)
	}
}
