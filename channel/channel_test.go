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
	c, err2 := New(s, true, 0, state.TestOutcome[0].Allocations[0].Destination, state.TestOutcome[0].Allocations[1].Destination)

	testNew := func(t *testing.T) {
		if err1 == nil {
			t.Error(`expected error constructing with a non turnNum=0 state, but got none`)
		}
		if err2 != nil {
			t.Error(err2)
		}
	}

	testPreFund := func(t *testing.T) {
		got, err1 := c.PreFundState().Hash()
		want, err2 := s.Hash()
		if err1 != nil {
			t.Error(err1)
		}
		if err2 != nil {
			t.Error(err2)
		}
		if got != want {
			t.Errorf(`incorrect incorrect PreFundState returned, got %v wanted %v`, c.PreFundState(), s)
		}
	}

	testPostFund := func(t *testing.T) {
		got, err1 := c.PostFundState().Hash()
		spf := s.Clone()
		spf.TurnNum = big.NewInt(1)
		want, err2 := spf.Hash()
		if err1 != nil {
			t.Error(err1)
		}
		if err2 != nil {
			t.Error(err2)
		}
		if got != want {
			t.Errorf(`incorrect incorrect PreFundState returned, got %v wanted %v`, c.PostFundState(), s)
		}
	}

	t.Run(`TestNew`, testNew)
	t.Run(`TestPreFund`, testPreFund)
	t.Run(`TestPostFund`, testPostFund)
}
