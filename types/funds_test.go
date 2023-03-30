package types

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/go-cmp/cmp"
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
	equalPairs := []fundsPair{
		{Sum(testData["a"],
			testData["b"]), testData["ab"]},
		{Sum(testData["a"],
			testData["b"],
			testData["c"]), testData["abc"]},
		{Sum(testData["a"],
			testData["b"],
			testData["c"],
			testData["d"]), testData["abcd"]},
	}

	// f == Sum(f, zeros), f == Sum(f, blank)
	for _, f := range testData {
		equalPairs = append(equalPairs, fundsPair{f, Sum(f, testData["zeros"])})
		equalPairs = append(equalPairs, fundsPair{f, Sum(f, testData["blank"])})
	}

	for i, p := range equalPairs {
		if !p.a.Equal(p.b) {
			t.Fatalf("test_sum_%d: expected %s to equal %s, but it did not", i, p.a, p.b)
		}
	}
}

func TestAdd(t *testing.T) {
	equalPairs := []fundsPair{
		{
			testData["a"].Add(testData["b"]),
			testData["ab"],
		},
		{
			testData["a"].Add(testData["b"], testData["c"]),
			testData["abc"],
		},
		{
			testData["a"].Add(testData["b"]).Add(testData["c"]).Add(testData["d"]),
			testData["abcd"],
		},
	}

	// f == f.Add(zeros), f == f.Add(blanks)
	for _, f := range testData {
		equalPairs = append(equalPairs, fundsPair{f, f.Add(testData["zeros"])})
		equalPairs = append(equalPairs, fundsPair{f, f.Add(testData["blank"])})
	}

	for i, p := range equalPairs {
		if !p.a.Equal(p.b) {
			t.Fatalf("test_sum_%d: expected %s to equal %s, but it did not", i, p.a, p.b)
		}
	}
}

func TestCanAfford(t *testing.T) {
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

	// Check self-affordance, affordance of blank, affordance of zeros
	for _, f := range testData {
		canAffordPairs = append(canAffordPairs, fundsPair{f, f})
		canAffordPairs = append(canAffordPairs, fundsPair{f, testData["blank"]})
		canAffordPairs = append(canAffordPairs, fundsPair{f, testData["zeros"]})
	}

	for _, p := range canAffordPairs {
		canAfford := p.a.canAfford(p.b)

		if !canAfford {
			t.Fatalf("expected %s to afford %s, but it didn't", p.a, p.b)
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
			t.Fatalf("expected %s to not afford %s, but it did", p.a, p.b)
		}
	}
}

func TestEqual(t *testing.T) {
	equalPairs := []fundsPair{
		{testData["zeros"], testData["blank"]},
		{testData["a"], testData["b"]},
	}

	// Check self-equality
	for _, f := range testData {
		equalPairs = append(equalPairs, fundsPair{f, f})
	}

	for _, p := range equalPairs {
		equal := p.a.Equal(p.b) && p.b.Equal(p.a)

		if !equal {
			t.Fatalf("expected %s to equal %s, but it didn't", p.a, p.b)
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
			t.Fatalf("expected %s to not equal %s, but it did", p.a, p.b)
		}
	}
}

type fundsPair struct {
	a Funds
	b Funds
}

func TestFundsClone(t *testing.T) {
	f := testData["a"]
	clone := f.Clone()

	if diff := cmp.Diff(f, clone); diff != "" {
		t.Fatalf("Clone: mismatch (-want +got):\n%s", diff)
	}
}
