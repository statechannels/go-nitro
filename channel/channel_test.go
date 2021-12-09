package channel

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/go-cmp/cmp"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/types"
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

	testPreFundSignedByMe := func(t *testing.T) {
		got := c.PreFundSignedByMe()
		want := false
		if got != want {
			t.Error(`expected c.PreFundSignedByMe() to be false, but it is true`)
		}

	}

	testPostFundSignedByMe := func(t *testing.T) {
		got := c.PostFundSignedByMe()
		want := false
		if got != want {
			t.Error(`expected c.PostFundSignedByMe() to be false, but it is true`)
		}

	}

	testPreFundComplete := func(t *testing.T) {
		got := c.PreFundComplete()
		want := false
		if got != want {
			t.Error(`expected c.PreFundComplete() to be false, but it is true`)
		}

	}

	testPostFundComplete := func(t *testing.T) {
		got := c.PostFundComplete()
		want := false
		if got != want {
			t.Error(`expected c.PostFundComplete() to be false, but it is true`)
		}
	}

	testLatestSupportedState := func(t *testing.T) {
		got, err1 := c.LatestSupportedState().Hash()
		want, err2 := s.Hash()
		if err1 != nil {
			t.Error(err1)
		}
		if err2 != nil {
			t.Error(err2)
		}
		if got != want {
			t.Errorf(`incorrect LatestSupportedState returned, got %v wanted %v`, c.LatestSupportedState(), s)
		}
	}

	testTotal := func(t *testing.T) {
		got := c.Total()
		want := types.Funds{
			common.Address{}: big.NewInt(10),
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("TestCrank: side effects mismatch (-want +got):\n%s", diff)
		}
	}

	t.Run(`TestNew`, testNew)
	t.Run(`TestPreFund`, testPreFund)
	t.Run(`TestPostFund`, testPostFund)
	t.Run(`TestPreFundSignedByMe`, testPreFundSignedByMe)
	t.Run(`TestPostFundSignedByMe`, testPostFundSignedByMe)
	t.Run(`TestPreFundComplete`, testPreFundComplete)
	t.Run(`TestPostFundComplete`, testPostFundComplete)
	t.Run(`TestLatestSupportedState`, testLatestSupportedState)
	t.Run(`TestTotal`, testTotal)
}
