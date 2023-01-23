package chainservice

import (
	"math/big"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSplitBlockRange(t *testing.T) {

	type bR struct {
		from uint
		to   uint
	}

	toBlockRange := func(br bR) blockRange {
		return blockRange{from: big.NewInt(int64(br.from)), to: big.NewInt(int64(br.to))}
	}

	toBlockRanges := func(brs []bR) []blockRange {
		blockRanges := make([]blockRange, len(brs))
		for i, br := range brs {
			blockRanges[i] = toBlockRange(br)
		}
		return blockRanges
	}

	type testCase struct {
		testBlockRange  bR
		testMaxInterval uint
		expectation     []bR
	}

	testCases := []testCase{
		{bR{from: 3, to: 7}, 2, []bR{{3, 5}, {6, 7}}},
		{bR{from: 0, to: 11282}, 5640, []bR{{0, 5640}, {5641, 11281}, {11282, 11282}}},
		{bR{from: 0, to: 1}, 100, []bR{{0, 1}}},
	}

	for i, tc := range testCases {
		if diff := cmp.Diff(
			splitBlockRange(
				toBlockRange(tc.testBlockRange),
				big.NewInt(int64(tc.testMaxInterval)),
			),
			toBlockRanges(tc.expectation),
			cmp.AllowUnexported(blockRange{}, big.Int{})); diff != "" {
			t.Fatalf("splitBlockRange output mismatch on test case %v. (-want +got):\n%s", i, diff)
		}
	}

}
