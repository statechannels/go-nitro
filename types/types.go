// Package types defines common types
package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// An ethereum address (20 bytes)
type Address = common.Address

// An ethereum hash (32 bytes)
type Bytes32 = common.Hash

// An arbitrary length byte slice
type Bytes []byte

// Compares two Bytes for equality
func (a Bytes) Equals(b Bytes) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != (b[i]) {
			return false
		}
	}
	return true
}

// We use a big.Int to represent Solidity's uint256
type Uint256 = big.Int

// Min returns the minimum of the supplied integers as a pointer
func Min(a *Uint256, b *Uint256) *Uint256 {
	switch a.Cmp(b) {
	case -1:
		return a
	default:
		return b
	}
}
