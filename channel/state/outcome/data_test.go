package outcome

import (
	"encoding/hex"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/types"
)

var alice = types.Destination(common.HexToHash("0x0a"))
var bob = types.Destination(common.HexToHash("0x0b"))

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

var e Exit = Exit{
	{
		Asset:       types.Address{}, // eth, fil, etc.
		Metadata:    zeroBytes,
		Allocations: allocsX,
	},
	{
		Asset:       types.Address{123}, // some token
		Metadata:    zeroBytes,
		Allocations: allocsY,
	},
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

// Exported doesn't show up in the package's godoc
var Exported = 1
