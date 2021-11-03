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

// We use a big.Int to represent Solidity's uint256
type Uint256 = big.Int

// A [token address]: amount map. Address 0 represents a chain's native token (ETH, FIL, etc)
type Holdings map[common.Address]*big.Int

// IsNonZero returns true if the Holdings structure has any non-zero asset
func (h Holdings) IsNonZero() bool {
	for _, v := range h {
		if v.Cmp(big.NewInt(0)) == 1 {
			return true
		}
	}
	return false
}

// String returns a bracket-separaged list of assets: {[0x0a,0x01][0x0b,0x01]}
func (h Holdings) String() string {
	if len(h) == 0 {
		return "{}"
	}
	var s string = "{"
	for asset, amount := range h {
		s += "[" + asset.Hex() + "," + amount.Text(64) + "]"
	}
	s = s + "}"
	return s
}

// Add sums all assets and returns the result. Does not modify the calling Holdings object
func (h Holdings) Add(a ...Holdings) Holdings {
	a = append(a, h)
	return Sum(a...)
}

func Sum(a ...Holdings) Holdings {
	sum := Holdings{}

	for _, holdings := range a {
		for asset, amount := range holdings {
			sum[asset] = sum[asset].Add(sum[asset], amount)
		}
	}

	return sum
}
