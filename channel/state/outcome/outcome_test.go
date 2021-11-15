package outcome

import (
	"bytes"
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
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

var zeroBytes = []byte{}
var testAllocations = Allocations{{
	Destination:    types.Destination(common.HexToHash("0x00000000000000000000000096f7123E3A80C9813eF50213ADEd0e4511CB820f")),
	Amount:         big.NewInt(1),
	AllocationType: 0,
	Metadata:       zeroBytes}}
var testExit = Exit{{Asset: common.HexToAddress("0x00"), Metadata: zeroBytes, Allocations: testAllocations}}

// copy-pasted from https://github.com/statechannels/exit-format/blob/201d4eb7554bac337a780cc8a640f6c45c3045a5/test/exit-format-ts.test.ts
var encodedExitReference, _ = hex.DecodeString("00000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000060000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000002000000000000000000000000096f7123e3a80c9813ef50213aded0e4511cb820f0000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000800000000000000000000000000000000000000000000000000000000000000000")

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
