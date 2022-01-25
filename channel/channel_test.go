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

	_, err1 := New(s, 0)
	s.TurnNum = 0
	c, err2 := New(s, 0)

	testNew := func(t *testing.T) {
		if err1 == nil {
			t.Error(`expected error constructing with a non turnNum=0 state, but got none`)
		}
		if err2 != nil {
			t.Error(err2)
		}
	}

	testClone := func(t *testing.T) {
		r := c.Clone()
		if diff := cmp.Diff(r, c, cmp.Comparer(types.Equal)); diff != "" {
			t.Errorf("Clone: mismatch (-want +got):\n%s", diff)
		}

		r.latestSupportedStateTurnNum++
		if r.Equal(c) {
			t.Error("Clone: modifying the clone should not modify the original")
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
			t.Errorf(`incorrect PreFundState returned, got %v wanted %v`, c.PreFundState(), s)
		}
	}

	testPostFund := func(t *testing.T) {
		got, err1 := c.PostFundState().Hash()
		spf := s.Clone()
		spf.TurnNum = PostFundTurnNum
		want, err2 := spf.Hash()
		if err1 != nil {
			t.Error(err1)
		}
		if err2 != nil {
			t.Error(err2)
		}
		if got != want {
			t.Errorf(`incorrect PreFundState returned, got %v wanted %v`, c.PostFundState(), s)
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
		_, err1 := c.LatestSupportedState()
		if err1 == nil {
			t.Error(`c.LatestSupportedState(): expected an error since no state is yet supported, but got none`)
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

	testAddSignedState := func(t *testing.T) {
		// Begin testing the cases that are NOOPs returning false
		want := false

		got := c.AddSignedState(s, state.Signature{}) // note null signature
		if got != want {
			t.Error(`expected c.AddSignedState() to be false, but it was true`)
		}
		nonParticipantSignature, _ := s.Sign(common.Hex2Bytes(`2030b463177db2da82908ef90fa55ddfcef56e8183caf60db464bc398e736e6f`))
		got = c.AddSignedState(s, nonParticipantSignature) // note signature by non participant
		if got != want {
			t.Error(`expected c.AddSignedState() to be false, but it was true`)
		}
		alicePrivateKey := common.Hex2Bytes(`caab404f975b4620747174a75f08d98b4e5a7053b691b41bcfc0d839d48b7634`)
		v := state.State{ // TODO it would be terser to clone s and modify it -- but s.Clone() is broken https://github.com/statechannels/go-nitro/issues/96
			ChainId: s.ChainId,
			Participants: []types.Address{
				common.HexToAddress(`0xF5A1BB5607C9D079E46d1B3Dc33f257d937b43BD`),
			},
			ChannelNonce:      big.NewInt(37140676581),
			AppDefinition:     common.HexToAddress(`0x5e29E5Ab8EF33F050c7cc10B5a0456D975C5F88d`),
			ChallengeDuration: big.NewInt(60),
			AppData:           []byte{},
			Outcome:           state.TestOutcome,
			TurnNum:           5,
			IsFinal:           false,
		}
		v.ChannelNonce.Add(v.ChannelNonce, big.NewInt(1))
		aliceSignatureOnWrongState, _ := v.Sign(alicePrivateKey)
		got = c.AddSignedState(v, aliceSignatureOnWrongState) // note state from wrong channel
		if got != want {
			t.Error(`expected c.AddSignedState() to be false, but it was true`)
		}
		c.latestSupportedStateTurnNum = uint64(3)
		aliceSignatureOnCorrectState, _ := c.PostFundState().Sign(alicePrivateKey)
		got = c.AddSignedState(c.PostFundState(), aliceSignatureOnCorrectState) // note stale state
		if got != want {
			t.Error(`expected c.AddSignedState() to be false, but it was true`)
		}
		c.latestSupportedStateTurnNum = MaxTurnNum // Reset so there is no longer a supported state

		// Now test cases which update the Channel and return true
		want = true

		got = c.AddSignedState(c.PostFundState(), aliceSignatureOnCorrectState)
		if got != want {
			t.Error(`expected c.AddSignedState() to be true, but it was false`)
		}

		got2 := c.SignedStateForTurnNum[1]
		if got2.State.Outcome == nil || !got2.HasSignature(0) {
			t.Error(`state not added correctly`)
		}

		// Add Bob's signature and check that we now have a supported state
		bobPrivateKey := common.Hex2Bytes(`62ecd49c4ccb41a70ad46532aed63cf815de15864bc415c87d507afd6a5e8da2`)
		bobSignatureOnCorrectState, _ := c.PostFundState().Sign(bobPrivateKey)
		got = c.AddSignedState(c.PostFundState(), bobSignatureOnCorrectState)
		if got != want {
			t.Error(`expected c.AddSignedState() to be true, but it was false`)
		}
		got3 := c.latestSupportedStateTurnNum
		want3 := uint64(1)
		if got3 != want3 {
			t.Errorf(`expected c.latestSupportedStateTurnNum to be %v, but got %v`, want, got)
		}
		got4, err4 := c.LatestSupportedState()
		if err4 != nil {
			t.Error(err4)
		}
		if got4.TurnNum != want3 {
			t.Errorf(`expected LatestSupportedState with turnNum %v`, want3)
		}

	}

	t.Run(`TestNew`, testNew)
	t.Run(`TestClone`, testClone)
	t.Run(`TestPreFund`, testPreFund)
	t.Run(`TestPostFund`, testPostFund)
	t.Run(`TestPreFundSignedByMe`, testPreFundSignedByMe)
	t.Run(`TestPostFundSignedByMe`, testPostFundSignedByMe)
	t.Run(`TestPreFundComplete`, testPreFundComplete)
	t.Run(`TestPostFundComplete`, testPostFundComplete)
	t.Run(`TestLatestSupportedState`, testLatestSupportedState)
	t.Run(`TestTotal`, testTotal)
	t.Run(`TestAddSignedState`, testAddSignedState)

}
