// Package types defines common types.
package types // import "github.com/statechannels/go-nitro/types"

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// An ethereum address (20 bytes)
type Address = common.Address

// An ethereum hash (32 bytes)
type Bytes32 = common.Hash

// Destination represents a payable address in go-nitro. In a state channel network,
// payable address are either:
//   - Internal: a 32-byte nitro channel ID, or
//   - External: a blockchain account or contract address, left-padded with 0s
type Destination Bytes32

// An arbitrary length byte slice
type Bytes []byte

// We use a big.Int to represent Solidity's uint256
type Uint256 = big.Int

// A {tokenAddress: amount} map. Address 0 represents a chain's native token (ETH, FIL, etc)
type Funds map[common.Address]*big.Int

// ConstError is a const-friendly error type.
type ConstError string

func (c ConstError) Error() string { return string(c) }
