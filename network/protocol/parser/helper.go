package parser

import (
	"fmt"
	"reflect"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/statechannels/go-nitro/types"
)

func i2Uint8(v any) uint8 {
	switch v := v.(type) {
	case int64:
		return uint8(v)
	case float64:
		return uint8(v)
	case uint32:
		return uint8(v)
	case uint8:
		return v
	}
	panic(fmt.Sprintf("invalid type %s", reflect.TypeOf(v)))
}

func i2Uint32(v any) uint32 {
	switch v := v.(type) {
	case int64:
		return uint32(v)
	case float64:
		return uint32(v)
	case uint32:
		return v
	}
	panic(fmt.Sprintf("invalid type %s", reflect.TypeOf(v)))
}

func i2Uint64(v any) uint64 {
	switch v := v.(type) {
	case int64:
		return uint64(v)
	case float64:
		return uint64(v)
	case uint64:
		return v
	}
	panic(fmt.Sprintf("invalid type %s", reflect.TypeOf(v)))
}

func i2Uint256(v any) *types.Uint256 {
	switch v := v.(type) {
	case string:
		bigInt, ok := math.ParseBig256(v)
		if !ok {
			panic(fmt.Sprintf("parsing to bigint failed. val: %s", v))
		}
		return bigInt
	case float64:
		bigInt, ok := math.ParseBig256(fmt.Sprintf("%v", v))
		if !ok {
			panic(fmt.Sprintf("parsing to bigint failed. val: %v", v))
		}
		return bigInt
	}
	panic(fmt.Sprintf("invalid type %s", reflect.TypeOf(v)))
}

func toByteArray(v any) []byte {
	switch v := v.(type) {
	case []byte:
		return v
	default:
		var data []byte
		return data
	}
}

func hexesToAddresses(addressesArr []string) []types.Address {
	addresses := make([]types.Address, len(addressesArr))
	for i := 0; i < len(addresses); i++ {
		addresses[i] = common.HexToAddress(addressesArr[i])
	}

	return addresses
}
