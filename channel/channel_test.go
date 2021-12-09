package channel

import (
	"math/big"
	"testing"

	"github.com/statechannels/go-nitro/channel/state"
)

func TestChannel(t *testing.T) {
	s := state.TestState.Clone()
	_, err1 := New(s, true, 0, state.TestOutcome[0].Allocations[0].Destination, state.TestOutcome[0].Allocations[1].Destination)
	s.TurnNum = big.NewInt(0)
	_, err2 := New(s, true, 0, state.TestOutcome[0].Allocations[0].Destination, state.TestOutcome[0].Allocations[1].Destination)

	testNew := func(t *testing.T) {
		if err1 == nil {
			t.Error(`expected error constructing with a non turnNum=0 state, but got none`)
		}
		if err2 != nil {
			t.Error(err2)
		}
	}

	t.Run(`TestNew`, testNew)
}
