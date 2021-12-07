package outcome

import (
	"bytes"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/types"
)

func TestEqualExits(t *testing.T) {
	var e1 = Exit{SingleAssetExit{
		Asset:    common.HexToAddress("0x00"),
		Metadata: make(types.Bytes, 0),
		Allocations: Allocations{{
			Destination:    types.Destination(common.HexToHash("0x0a")),
			Amount:         big.NewInt(2),
			AllocationType: 0,
			Metadata:       make(types.Bytes, 0)}},
	}}

	// equal to e1
	var e2 = Exit{SingleAssetExit{
		Asset:    common.HexToAddress("0x00"),
		Metadata: make(types.Bytes, 0),
		Allocations: Allocations{{
			Destination:    types.Destination(common.HexToHash("0x0a")),
			Amount:         big.NewInt(2),
			AllocationType: 0,
			Metadata:       make(types.Bytes, 0)}},
	}}

	if &e1 == &e2 {
		t.Error("expected distinct pointers, but got idendical pointers")
	}

	if !e1.Equal(e2) {
		t.Error("expected equal Exits, but got distinct Exits")
	}

	// each equal to e1 except in one aspect
	var distinctExits []Exit = []Exit{
		{SingleAssetExit{
			Asset:    common.HexToAddress("0x01"), // distinct Asset
			Metadata: make(types.Bytes, 0),
			Allocations: Allocations{{
				Destination:    types.Destination(common.HexToHash("0x0a")),
				Amount:         big.NewInt(2),
				AllocationType: 0,
				Metadata:       make(types.Bytes, 0)}},
		}},
		{SingleAssetExit{
			Asset:    common.HexToAddress("0x00"),
			Metadata: []byte{1}, // distinct metadata
			Allocations: Allocations{{
				Destination:    types.Destination(common.HexToHash("0x0a")),
				Amount:         big.NewInt(2),
				AllocationType: 0,
				Metadata:       make(types.Bytes, 0)}},
		}},
		{SingleAssetExit{
			Asset:    common.HexToAddress("0x00"),
			Metadata: make(types.Bytes, 0),
			Allocations: Allocations{{
				Destination:    types.Destination(common.HexToHash("0x0b")), // distinct destination
				Amount:         big.NewInt(2),
				AllocationType: 0,
				Metadata:       make(types.Bytes, 0)}},
		}},
		{SingleAssetExit{
			Asset:    common.HexToAddress("0x00"),
			Metadata: make(types.Bytes, 0),
			Allocations: Allocations{{
				Destination:    types.Destination(common.HexToHash("0x0a")),
				Amount:         big.NewInt(3), // distinct amount
				AllocationType: 0,
				Metadata:       make(types.Bytes, 0)}},
		}},
		{SingleAssetExit{
			Asset:    common.HexToAddress("0x00"),
			Metadata: make(types.Bytes, 0),
			Allocations: Allocations{{
				Destination:    types.Destination(common.HexToHash("0x0a")),
				Amount:         big.NewInt(2),
				AllocationType: 1, // distinct allocationType
				Metadata:       make(types.Bytes, 0)}},
		}},
		{SingleAssetExit{
			Asset:    common.HexToAddress("0x00"),
			Metadata: make(types.Bytes, 0),
			Allocations: Allocations{{
				Destination:    types.Destination(common.HexToHash("0x0a")),
				Amount:         big.NewInt(2),
				AllocationType: 0,
				Metadata:       []byte{1}}}, // distinct metadata
		}},
	}

	for _, v := range distinctExits {
		if e1.Equal(v) {
			t.Error("expected distinct Exits but found them equal")
		}
	}
}

func TestExitEncode(t *testing.T) {
	var encodedExit, err = testExit.Encode()

	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(encodedExit, encodedExitReference) {
		t.Errorf("incorrect encoding. Got %x, wanted %x", encodedExit, encodedExitReference)
	}
}

func TestExitDecode(t *testing.T) {
	var decodedExit, err = Decode(encodedExitReference)
	if err != nil {
		t.Error(err)
	}

	if !testExit.Equal(decodedExit) {
		t.Error("decoded exit does not match expectation")
	}
}

func TestTotal(t *testing.T) {

	total := allocsX.Total()
	if total.Cmp(big.NewInt(5)) != 0 {
		t.Errorf(`Expected total to be 5, got %v`, total)
	}
}

func TestTotalAllocated(t *testing.T) {
	want := types.Funds{
		types.Address{}:    big.NewInt(5),
		types.Address{123}: big.NewInt(3),
	}

	got := e.TotalAllocated()

	if !got.Equal(want) {
		t.Errorf("Expected %v.TotalAllocated() to equal %v, but it was %v",
			e, want, got)
	}
}

func TestTotalFor(t *testing.T) {
	testCases := []struct {
		Exit        Exit
		Participant types.Destination
		Want        types.Funds
	}{
		{e, alice, types.Funds{
			types.Address{}:    big.NewInt(2),
			types.Address{123}: big.NewInt(1),
		}},
		{e, bob, types.Funds{
			types.Address{}:    big.NewInt(3),
			types.Address{123}: big.NewInt(2),
		}},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprint("Case ", i), func(t *testing.T) {
			got := testCase.Exit.TotalAllocatedFor(testCase.Participant)
			if !got.Equal(testCase.Want) {
				t.Errorf("Expected TotalAllocatedFor for participant %v on exit %v to be %v, but got %v",
					testCase.Participant, testCase.Exit, testCase.Want, got)
			}
		})
	}
}

func TestExitAffords(t *testing.T) {

	allocationMap := map[types.Address]Allocation{
		{}: testExit[0].Allocations[0],
	}

	got := testExit.Affords(allocationMap, types.Funds{}) // This should not panic
	want := false
	if !(got == want) {
		t.Error(`Affords: expected test exit to not afford the given allocation with nil funds`)
	}
}
