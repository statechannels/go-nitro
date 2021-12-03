package outcome

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/types"
)

var testAllocations = Allocations{{
	Destination:    types.Destination(common.HexToHash("0x00000000000000000000000096f7123E3A80C9813eF50213ADEd0e4511CB820f")),
	Amount:         big.NewInt(1),
	AllocationType: 0,
	Metadata:       zeroBytes}}

var allocsX = Allocations{ // [{Alice: 2, Bob: 3}]
	{
		Destination:    alice,
		Amount:         big.NewInt(2),
		AllocationType: 0,
		Metadata:       make(types.Bytes, 0)},
	{
		Destination:    bob,
		Amount:         big.NewInt(3),
		AllocationType: 0,
		Metadata:       make(types.Bytes, 0)},
}

var allocsY = Allocations{ // [{Bob: 2, Alice: 1}]
	{
		Destination:    bob,
		Amount:         big.NewInt(2),
		AllocationType: 0,
		Metadata:       make(types.Bytes, 0)},
	{
		Destination:    alice,
		Amount:         big.NewInt(1),
		AllocationType: 0,
		Metadata:       make(types.Bytes, 0)},
}

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
