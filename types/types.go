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
