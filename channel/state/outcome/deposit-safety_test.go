package outcome

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/statechannels/go-nitro/types"
)

func TestDepositSafetyThreshold(t *testing.T) {
	testCases := []struct {
		Exit        Exit
		Participant types.Destination
		Want        types.Funds
	}{
		{e, alice, types.Funds{
			types.Address{}:    big.NewInt(0),
			types.Address{123}: big.NewInt(2),
		}},
		{e, bob, types.Funds{
			types.Address{}:    big.NewInt(2),
			types.Address{123}: big.NewInt(0),
		}},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprint("Case ", i), func(t *testing.T) {
			got := testCase.Exit.DepositSafetyThreshold(testCase.Participant)
			if !got.Equal(testCase.Want) {
				t.Fatalf("Expected safety threshold for participant %v on exit %v to be %v, but got %v",
					testCase.Participant, testCase.Exit, testCase.Want, got)
			}
		})
	}
}
