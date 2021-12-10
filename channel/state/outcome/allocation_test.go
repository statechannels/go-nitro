package outcome

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/go-cmp/cmp"
	"github.com/statechannels/go-nitro/types"
)

func TestEqualAllocations(t *testing.T) {

	var a1 = Allocations{{ // [{Alice: 2}]
		Destination:    types.Destination(common.HexToHash("0x0a")),
		Amount:         big.NewInt(2),
		AllocationType: 0,
		Metadata:       make(types.Bytes, 0)}}

	var a2 = Allocations{{ // [{Alice: 2}]
		Destination:    types.Destination(common.HexToHash("0x0a")),
		Amount:         big.NewInt(2),
		AllocationType: 0,
		Metadata:       make(types.Bytes, 0)}}

	if &a1 == &a2 {
		t.Errorf("expected distinct pointers, but got identical pointers")
	}

	if !a1.Equal(a2) {
		t.Errorf("expected equal Allocations, but got distinct Allocations")
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
				t.Errorf(
					`Incorrect AffordFor: expected %v.Affords(%v,%v) to be %v, but got %v`,
					testcase.Allocations, testcase.GivenAllocation, testcase.Funding, testcase.Want, got)
			}
		})

	}

}

func TestDivertToGuarantee(t *testing.T) {

	aliceDestination := types.Destination(common.HexToHash("0x0a"))
	bobDestination := types.Destination(common.HexToHash("0x0b"))

	targetChannel := types.Destination(common.HexToHash("0xabc"))

	a := Allocations{
		{
			Destination:    aliceDestination,
			Amount:         big.NewInt(243),
			AllocationType: 0,
			Metadata:       make(types.Bytes, 0),
		},
		{
			Destination:    bobDestination,
			Amount:         big.NewInt(309),
			AllocationType: 0,
			Metadata:       make(types.Bytes, 0),
		},
	}

	got, err := a.DivertToGuarantee(aliceDestination, bobDestination, big.NewInt(5), big.NewInt(5), targetChannel)

	if err != nil {
		t.Error(err)
	}

	want := Allocations{
		{
			Destination:    aliceDestination,
			Amount:         big.NewInt(238),
			AllocationType: 0,
			Metadata:       make(types.Bytes, 0),
		},
		{
			Destination:    bobDestination,
			Amount:         big.NewInt(304),
			AllocationType: 0,
			Metadata:       make(types.Bytes, 0),
		},
		{
			Destination:    targetChannel,
			Amount:         big.NewInt(10),
			AllocationType: 1,
			Metadata:       append(aliceDestination.Bytes(), bobDestination.Bytes()...),
		},
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("TestNew: expectedGuarantee mismatch (-want +got):\n%s", diff)
	}
}
