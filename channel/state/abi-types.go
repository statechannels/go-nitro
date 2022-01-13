package state

import "github.com/ethereum/go-ethereum/accounts/abi"

// To encode objects as bytes, we need to construct an encoder, using abi.Arguments.
// An instance of abi.Arguments implements two functions relevant to us:
// - `Pack`, which packs go values for a given struct into bytes.
// - `unPack`, which unpacks bytes into go values
// To construct an abi.Arguments instance, we need to supply an array of "types", which are
// actually go values. The following types are used when encoding a state

// uint256 is the uint256 type for abi encoding
var uint256, _ = abi.NewType("uint256", "uint256", nil)

// bool is the bool type for abi encoding
var boolTy, _ = abi.NewType("bool", "bool", nil)

// destination is the bytes32 type for abi encoding
var destination, _ = abi.NewType("bytes32", "address", nil)

// bytes is the bytes type for abi encoding
var bytesTy, _ = abi.NewType("bytes", "bytes", nil)

// address is the address[] type for abi encoding
var addressArray, _ = abi.NewType("address[]", "address[]", nil)

// address is the address type for abi encoding
var address, _ = abi.NewType("address", "address", nil)
