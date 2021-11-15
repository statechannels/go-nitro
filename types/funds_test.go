package types

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

var testData map[string]Funds = map[string]Funds{
	"a": {
		common.HexToAddress("0x00"): big.NewInt(1),
	},
	"b": {
		common.HexToAddress("0x00"): big.NewInt(1),
	},
	"c": {
		common.HexToAddress("0x01"): big.NewInt(1),
	},
	"d": {
		common.HexToAddress("0x02"): big.NewInt(1),
	},
	// manually calculated sums
	"ab": {
		common.HexToAddress("0x00"): big.NewInt(2),
	},
	"ac": {
		common.HexToAddress("0x00"): big.NewInt(1),
		common.HexToAddress("0x01"): big.NewInt(1),
	},
	"abc": {
		common.HexToAddress("0x00"): big.NewInt(2),
		common.HexToAddress("0x01"): big.NewInt(1),
	},
	"abcd": {
		common.HexToAddress("0x00"): big.NewInt(2),
		common.HexToAddress("0x01"): big.NewInt(1),
		common.HexToAddress("0x02"): big.NewInt(1),
	},
}

func TestSum(t *testing.T) {

}

func TestEqual(t *testing.T) {

}
