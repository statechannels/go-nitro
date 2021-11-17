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

func TestCanAfford(t *testing.T) {
	expectedButDidNot := "expected %s to afford %s, but it didn't"
	expectedNotButDid := "expected %s to not afford %s, but it did"

	// Check self-affordance, affordance of blank, affordance of zeros
	for _, f := range testData {
		canAfford := f.canAfford(f)

		if !canAfford {
			t.Errorf(expectedButDidNot, f, f)
		}

		canAffordBlank := f.canAfford(testData["blank"])

		if !canAffordBlank {
			t.Errorf(expectedButDidNot, f, testData["blank"])
		}

		canAffordZeros := f.canAfford(testData["zeros"])

		if !canAffordZeros {
			t.Errorf(expectedButDidNot, f, testData["zeros"])
		}
	}

	canAffordPairs := []fundsPair{
		{testData["a"], testData["b"]}, // equal funds
		{testData["b"], testData["a"]},
		{testData["f"], testData["a"]}, // more of single asset
		{testData["ab"], testData["a"]},
		{testData["ab"], testData["b"]},
		{testData["abc"], testData["b"]}, // mixed assets, enough of "relevant asset(s)"
		{testData["abcd"], testData["b"]},
		{testData["abcd"], testData["ab"]},
		{testData["abcd"], testData["abc"]},
	}

	for _, p := range canAffordPairs {
		canAfford := p.a.canAfford(p.b)

		if !canAfford {
			t.Errorf(expectedButDidNot, p.a, p.b)
		}
	}

	cannotAffordPairs := []fundsPair{
		{testData["zeros"], testData["a"]}, // zero cannot afford things
		{testData["blank"], testData["a"]}, // blank cannot afford things
		{testData["a"], testData["c"]},     // different assets
		{testData["c"], testData["a"]},
		{testData["a"], testData["f"]}, // less of single asset
		{testData["a"], testData["ab"]},
		{testData["b"], testData["ab"]},
		{testData["b"], testData["abc"]}, // mixed assets, less of "relevant asset(s)"
		{testData["a"], testData["abcd"]},
		{testData["ab"], testData["abcd"]},
		{testData["abc"], testData["abcd"]},
	}

	for _, p := range cannotAffordPairs {
		canAfford := p.a.canAfford(p.b)

		if canAfford {
			t.Errorf(expectedNotButDid, p.a, p.b)
		}
	}
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
		{testData["zeros"], testData["blank"]},
		{testData["a"], testData["b"]},
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
		{testData["blank"], testData["a"]},
		{testData["zeros"], testData["a"]},
		{testData["abc"], testData["abcd"]},
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
