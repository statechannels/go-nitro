package parser

import (
	"fmt"
	"reflect"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/statechannels/go-nitro/types"
)

func i2Uint8(v any) uint8 {
	switch v.(type) {
	case int64:
		return uint8(v.(int64))
	case float64:
		return uint8(v.(float64))
	case uint32:
		return uint8(v.(uint32))
	case uint8:
		return v.(uint8)
	}
	panic(fmt.Sprintf("invalid type %s", reflect.TypeOf(v)))
}

func i2Uint32(v any) uint32 {
	switch v.(type) {
	case int64:
		return uint32(v.(int64))
	case float64:
		return uint32(v.(float64))
	case uint32:
		return v.(uint32)
	}
	panic(fmt.Sprintf("invalid type %s", reflect.TypeOf(v)))
}

func i2Uint64(v any) uint64 {
	switch v.(type) {
	case int64:
		return uint64(v.(int64))
	case float64:
		return uint64(v.(float64))
	case uint64:
		return v.(uint64)
	}
	panic(fmt.Sprintf("invalid type %s", reflect.TypeOf(v)))
}

func i2Uint256(v any) *types.Uint256 {
	switch v.(type) {
	case string:
		bigInt, ok := math.ParseBig256(v.(string))
		if !ok {
			panic(fmt.Sprintf("parsing to bigint failed. val: %s", v.(string)))
		}
		return bigInt
	case float64:
		bigInt, ok := math.ParseBig256(fmt.Sprintf("%v", v.(float64)))
		if !ok {
			panic(fmt.Sprintf("parsing to bigint failed. val: %v", v.(float64)))
		}
		return bigInt
	}
	panic(fmt.Sprintf("invalid type %s", reflect.TypeOf(v)))
}

func toByteArray(v any) []byte {
	switch v.(type) {
	case []byte:
		return v.([]byte)
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
