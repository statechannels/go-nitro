package channel

import (
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/go-cmp/cmp"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/helpers"
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
		if helpers.HasShallowCopy(r, c) {
			t.Fatal("Clone has shallow copy")
		}

		if diff := cmp.Diff(*r, *c, cmp.Comparer(types.Equal)); diff != "" {
			t.Fatalf("Clone: mismatch (-want +got):\n%s", diff)
		}

		r.latestSupportedStateTurnNum++
		if r.Equal(*c) {
			t.Error("Clone: modifying the clone should not modify the original")
		}

		r.Participants[0] = common.HexToAddress("0x0000000000000000000000000000000000000001")
		if r.Participants[0] == c.Participants[0] {
			t.Error("Clone: modifying the clone should not modify the original")
		}

		var nilChannel *Channel
		clone := nilChannel.Clone()
		if clone != nil {
			t.Fatal("Tried to clone a Channel via a nil pointer, but got something not nil")
		}

		// TODO This is just testing hasShallowCopy so it should be moved to the helpers package
		r = c.Clone()
		// Modify our clone so it is a shallow copy refering to the same map values
		c.SignedStateForTurnNum = r.SignedStateForTurnNum
		if isShallow := helpers.HasShallowCopy(r, c); !isShallow {
			t.Fatal("Expected isShallowCopy to return true")
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
			t.Fatalf(`incorrect PreFundState returned, got %v wanted %v`, c.PreFundState(), s)
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
			t.Fatalf(`incorrect PreFundState returned, got %v wanted %v`, c.PostFundState(), s)
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

	testLatestSignedState := func(t *testing.T) {
		_, err := c.LatestSignedState()
		if errors.Is(err, errors.New("No states are signed")) {
			t.Error(`c.LatestSignedState(): expected an empty SingedState since no state is yet signed`)
		}
	}

	testTotal := func(t *testing.T) {
		got := c.Total()
		want := types.Funds{
			common.Address{}: big.NewInt(10),
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Fatalf("TestCrank: side effects mismatch (-want +got):\n%s", diff)
		}
	}

	alicePrivateKey := common.Hex2Bytes(`caab404f975b4620747174a75f08d98b4e5a7053b691b41bcfc0d839d48b7634`)
	bobPrivateKey := common.Hex2Bytes(`62ecd49c4ccb41a70ad46532aed63cf815de15864bc415c87d507afd6a5e8da2`)
	testAddSignedState := func(t *testing.T) {
		myC, _ := New(s, 0)
		ss := state.NewSignedState(s)
		sigA, err := ss.State().Sign(alicePrivateKey)
		if err != nil {
			t.Error(err)
		}
		sigB, err := ss.State().Sign(bobPrivateKey)
		if err != nil {
			t.Error(err)
		}
		err = ss.AddSignature(sigA)
		if err != nil {
			t.Error(err)
		}
		err = ss.AddSignature(sigB)
		if err != nil {
			t.Error(err)
		}
		if ok := myC.AddSignedState(ss); !ok {
			t.Error("AddSignedState returned false")
		}

		// It should properly update the latestSupportedStateNum
		if myC.latestSupportedStateTurnNum != 0 {
			t.Fatalf("Expected latestSupportedStateTurnNum of 0 but got %d", myC.latestSupportedStateTurnNum)

		}
		// verify the signatures
		expectedSigs := []state.Signature{sigA, sigB}
		for i := range myC.Participants {
			gotSig, err := myC.SignedStateForTurnNum[s.TurnNum].GetParticipantSignature(uint(i))
			if err != nil {
				panic(err)
			}
			wantSig := expectedSigs[i]
			if !gotSig.Equal(wantSig) {
				t.Fatalf("Expected to find signature %x at index 0, but got %x", wantSig, gotSig)
			}
		}
	}

	testAddSignedStates := func(t *testing.T) {
		myC, _ := New(s, 0)

		ss := state.NewSignedState(s)
		sigA, err := ss.State().Sign(alicePrivateKey)
		if err != nil {
			t.Error(err)
		}
		sigB, err := ss.State().Sign(bobPrivateKey)
		if err != nil {
			t.Error(err)
		}
		err = ss.AddSignature(sigA)
		if err != nil {
			t.Error(err)
		}
		err = ss.AddSignature(sigB)
		if err != nil {
			t.Error(err)
		}
		if ok := myC.AddSignedStates([]state.SignedState{ss}); !ok {
			t.Error("AddSignedStates returned false")
		}

		// It should properly update the latestSupportedStateNum
		got := myC.latestSupportedStateTurnNum
		if got != 0 {
			t.Fatalf("Expected latestSupportedStateTurnNum of 0 but got %d", got)
		}

		// verify the signatures
		expectedSigs := []state.Signature{sigA, sigB}
		for i := range myC.Participants {
			gotSig, err := myC.SignedStateForTurnNum[s.TurnNum].GetParticipantSignature(uint(i))
			if err != nil {
				panic(err)
			}
			wantSig := expectedSigs[i]
			if !gotSig.Equal(wantSig) {
				t.Fatalf("Expected to find signature %x at index 0, but got %x", wantSig, gotSig)
			}
		}

	}

	testAddStateWithSignature := func(t *testing.T) {
		// Begin testing the cases that are NOOPs returning false
		want := false
		got := c.AddStateWithSignature(s, state.Signature{})
		if got != want {
			t.Error(`expected c.AddSignedState() to be false, but it was true`)
		}
		nonParticipantSignature, _ := s.Sign(common.Hex2Bytes(`2030b463177db2da82908ef90fa55ddfcef56e8183caf60db464bc398e736e6f`))
		got = c.AddStateWithSignature(s, nonParticipantSignature) // note signature by non participant
		if got != want {
			t.Error(`expected c.AddSignedState() to be false, but it was true`)
		}

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
		got = c.AddStateWithSignature(v, aliceSignatureOnWrongState) // note state from wrong channel
		if got != want {
			t.Error(`expected c.AddSignedState() to be false, but it was true`)
		}
		c.latestSupportedStateTurnNum = uint64(3)
		aliceSignatureOnCorrectState, _ := c.PostFundState().Sign(alicePrivateKey)
		got = c.AddStateWithSignature(c.PostFundState(), aliceSignatureOnCorrectState) // note stale state
		if got != want {
			t.Error(`expected c.AddSignedState() to be false, but it was true`)
		}
		c.latestSupportedStateTurnNum = MaxTurnNum // Reset so there is no longer a supported state

		// Now test cases which update the Channel and return true
		want = true

		got = c.AddStateWithSignature(c.PostFundState(), aliceSignatureOnCorrectState)
		if got != want {
			t.Error(`expected c.AddSignedState() to be true, but it was false`)
		}

		// Check whether latestSignedState is correct
		latestSignedState, err := c.LatestSignedState()
		if err != nil {
			t.Error(err)
		}
		expectedSignedState := state.NewSignedState(c.PostFundState())
		err = expectedSignedState.Sign(&alicePrivateKey)
		if err != nil {
			t.Error(err)
		}
		if diff := cmp.Diff(expectedSignedState, latestSignedState, cmp.Comparer(types.Equal)); diff != "" {
			t.Fatalf("LatestSignedState: mismatch (-want +got):\n%s", diff)
		}

		got2 := c.SignedStateForTurnNum[1]
		if got2.State().Outcome == nil || !got2.HasSignatureForParticipant(0) {
			t.Error(`state not added correctly`)
		}

		// Add Bob's signature and check that we now have a supported state
		bobPrivateKey := common.Hex2Bytes(`62ecd49c4ccb41a70ad46532aed63cf815de15864bc415c87d507afd6a5e8da2`)
		bobSignatureOnCorrectState, _ := c.PostFundState().Sign(bobPrivateKey)
		got = c.AddStateWithSignature(c.PostFundState(), bobSignatureOnCorrectState)
		if got != want {
			t.Error(`expected c.AddSignedState() to be true, but it was false`)
		}
		got3 := c.latestSupportedStateTurnNum
		want3 := uint64(1)
		if got3 != want3 {
			t.Fatalf(`expected c.latestSupportedStateTurnNum to be %v, but got %v`, want, got)
		}
		got4, err4 := c.LatestSupportedState()
		if err4 != nil {
			t.Error(err4)
		}
		if got4.TurnNum != want3 {
			t.Fatalf(`expected LatestSupportedState with turnNum %v`, want3)
		}

		// Check whether latestSignedState is correct
		latestSignedState, err = c.LatestSignedState()
		if err != nil {
			t.Error(err)
		}
		err = expectedSignedState.Sign(&bobPrivateKey)
		if err != nil {
			t.Error(err)
		}
		if diff := cmp.Diff(latestSignedState, expectedSignedState, cmp.Comparer(types.Equal)); diff != "" {
			t.Fatalf("LatestSignedState: mismatch (-want +got):\n%s", diff)
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
	t.Run(`TestLatestSignedState`, testLatestSignedState)
	t.Run(`TestTotal`, testTotal)
	t.Run(`TestAddStateWithSignature`, testAddStateWithSignature)
	t.Run(`TestAddSignedStates`, testAddSignedStates)
	t.Run(`TestAddSignedState`, testAddSignedState)

}

func TestTwoPartyLedger(t *testing.T) {
	s := state.TestState.Clone()
	s.TurnNum = 0
	testClone := func(t *testing.T) {
		r, err := NewTwoPartyLedger(s, 0)
		if err != nil {
			t.Fatal(err)
		}
		c := r.Clone()
		if diff := cmp.Diff(*r, *c, cmp.Comparer(types.Equal)); diff != "" {
			t.Fatalf("Clone: mismatch (-want +got):\n%s", diff)
		}

		r.latestSupportedStateTurnNum++
		if r.Channel.Equal(c.Channel) {
			t.Error("Clone: modifying the clone should not modify the original")
		}

		r.Participants[0] = common.HexToAddress("0x0000000000000000000000000000000000000001")
		if r.Participants[0] == c.Participants[0] {
			t.Error("Clone: modifying the clone should not modify the original")
		}

		var nilTwoPartyLedger *TwoPartyLedger
		clone := nilTwoPartyLedger.Clone()
		if clone != nil {
			t.Fatal("Tried to clone a TwoPartyLedger via a nil pointer, but got something not nil")
		}

	}

	t.Run(`TestClone`, testClone)
}

func TestSingleHopVirtualChannel(t *testing.T) {
	s := state.TestState.Clone()
	s.Participants = append(s.Participants, s.Participants[0]) // ensure three participants
	s.TurnNum = 0
	testClone := func(t *testing.T) {
		r, err := NewSingleHopVirtualChannel(s, 0)
		if err != nil {
			t.Fatal(err)
		}
		c := r.Clone()
		if diff := cmp.Diff(*r, *c, cmp.Comparer(types.Equal)); diff != "" {
			t.Fatalf("Clone: mismatch (-want +got):\n%s", diff)
		}

		r.latestSupportedStateTurnNum++
		if r.Channel.Equal(c.Channel) {
			t.Error("Clone: modifying the clone should not modify the original")
		}

		r.Participants[0] = common.HexToAddress("0x0000000000000000000000000000000000000001")
		if r.Participants[0] == c.Participants[0] {
			t.Error("Clone: modifying the clone should not modify the original")
		}

		var nilChannel *SingleHopVirtualChannel
		clone := nilChannel.Clone()
		if clone != nil {
			t.Fatal("Tried to clone a Channel via a nil pointer, but got something not nil")
		}

	}
	t.Run(`TestClone`, testClone)
}
