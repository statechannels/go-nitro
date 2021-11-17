package types

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

var testData map[string]Funds = map[string]Funds{
	"blank": {},
	"zeros": {
		common.HexToAddress("0x00"): big.NewInt(0),
		common.HexToAddress("0x01"): big.NewInt(0),
		common.HexToAddress("0x02"): big.NewInt(0),
		common.HexToAddress("0x03"): big.NewInt(0),
		common.HexToAddress("0x0a"): big.NewInt(0),
		common.HexToAddress("0x0b"): big.NewInt(0),
		common.HexToAddress("0x0c"): big.NewInt(0),
	},
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
	"e": {
		common.HexToAddress("0x00"):  big.NewInt(1),
		common.HexToAddress("0xabc"): big.NewInt(0),
	},
	"f": {
		common.HexToAddress("0x00"): big.NewInt(2),
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
	// Check self-equality
	for _, f := range testData {
		equal := f.Equal(f)

		if !equal {
			t.Errorf("expected %s to equal %s, but it didn't", f, f)
		}
	}

	equalPairs := []fundsPair{
		{testData["a"], testData["b"]},
		{testData["a"], testData["e"]},
	}

	for _, p := range equalPairs {
		equal := p.a.Equal(p.b) && p.b.Equal(p.a)

		if !equal {
			t.Errorf("expected %s to equal %s, but it didn't", p.a, p.b)
		}
	}

	unequalPairs := []fundsPair{
		{testData["a"], testData["c"]},
		{testData["a"], testData["d"]},
		{testData["a"], testData["ab"]},
	}

	for _, p := range unequalPairs {
		equal := p.a.Equal(p.b)

		if equal {
			t.Errorf("expected %s to not equal %s, but it did", p.a, p.b)
		}
	}

}

type fundsPair struct {
	a Funds
	b Funds
}
