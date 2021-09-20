package outcome

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/types"
)

func TestEqualAllocations(t *testing.T) {

	var a1 = Allocations{{ // [{Alice: 2}]
		Destination:    common.HexToHash("0x0a"),
		Amount:         big.NewInt(2),
		AllocationType: 0,
		Metadata:       make(types.Bytes, 0)}}

	var a2 = Allocations{{ // [{Alice: 2}]
		Destination:    common.HexToHash("0x0a"),
		Amount:         big.NewInt(2),
		AllocationType: 0,
		Metadata:       make(types.Bytes, 0)}}

	if &a1 == &a2 {
		t.Errorf("expected distinct pointers, but got identical pointers")
	}

	if !a1.Equals(a2) {
		t.Errorf("expected equal Allocations, but got distinct Allocations")
	}

}

func TestExitEncode(t *testing.T) {

	zeroBytes := make(types.Bytes, 0)

	var a = Allocations{{
		Destination:    common.HexToHash("0x00000000000000000000000096f7123E3A80C9813eF50213ADEd0e4511CB820f"),
		Amount:         big.NewInt(1),
		AllocationType: 0,
		Metadata:       zeroBytes}}

	var exit = Exit{{Asset: common.HexToAddress("0x00"), Metadata: zeroBytes, Allocations: a}}
	var encodedExit, error = exit.Encode()

	var want, _ = hex.DecodeString("00000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000060000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000002000000000000000000000000096f7123e3a80c9813ef50213aded0e4511cb820f0000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000800000000000000000000000000000000000000000000000000000000000000000")
	// copy-pasted from https://github.com/statechannels/exit-format/blob/201d4eb7554bac337a780cc8a640f6c45c3045a5/test/exit-format-ts.test.ts
	if error != nil {
		t.Error(error)
	}

	if !encodedExit.Equals(want) {
		t.Errorf("incorrect encoding. Got %x, wanted %x", encodedExit, want)
	}

}
