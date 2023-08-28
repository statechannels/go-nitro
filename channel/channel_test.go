package channel

import (
	"encoding/json"
	"errors"
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/go-cmp/cmp"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/internal/testhelpers"
	"github.com/statechannels/go-nitro/types"
)

func TestChannel(t *testing.T) {
	compareChannels := func(a, b *Channel) string {
		return cmp.Diff(*a, *b, cmp.AllowUnexported(*a, big.Int{}, state.SignedState{}, OffChainData{}, OnChainData{}))
	}

	compareStates := func(a, b state.SignedState) string {
		return cmp.Diff(a, b, cmp.AllowUnexported(a, big.Int{}))
	}

	s := state.TestState.Clone()

	c, err2 := New(s, 0)

	testNew := func(t *testing.T) {
		if err2 != nil {
			t.Error(err2)
		}
	}

	testClone := func(t *testing.T) {
		r := c.Clone()

		if diff := compareChannels(c, r); diff != "" {
			t.Errorf("Clone: mismatch (-want +got):\n%s", diff)
		}

		r.OffChain.LatestSupportedStateTurnNum++
		if reflect.DeepEqual(r, *c) {
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
		if myC.OffChain.LatestSupportedStateTurnNum != s.TurnNum {
			t.Fatalf("Expected latestSupportedStateTurnNum of %d but got %d", s.TurnNum, myC.OffChain.LatestSupportedStateTurnNum)
		}
		// verify the signatures
		expectedSigs := []state.Signature{sigA, sigB}
		for i := range myC.Participants {
			gotSig, err := myC.OffChain.SignedStateForTurnNum[s.TurnNum].GetParticipantSignature(uint(i))
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
			Participants: []types.Address{
				common.HexToAddress(`0xF5A1BB5607C9D079E46d1B3Dc33f257d937b43BD`),
			},
			ChannelNonce:      37140676581,
			AppDefinition:     common.HexToAddress(`0x5e29E5Ab8EF33F050c7cc10B5a0456D975C5F88d`),
			ChallengeDuration: 60,
			AppData:           []byte{},
			Outcome:           state.TestOutcome,
			TurnNum:           5,
			IsFinal:           false,
		}
		v.ChannelNonce += 1
		aliceSignatureOnWrongState, _ := v.Sign(alicePrivateKey)
		got = c.AddStateWithSignature(v, aliceSignatureOnWrongState) // note state from wrong channel
		if got != want {
			t.Error(`expected c.AddSignedState() to be false, but it was true`)
		}
		c.OffChain.LatestSupportedStateTurnNum = uint64(3)
		aliceSignatureOnCorrectState, _ := c.PostFundState().Sign(alicePrivateKey)
		got = c.AddStateWithSignature(c.PostFundState(), aliceSignatureOnCorrectState) // note stale state
		if got != want {
			t.Error(`expected c.AddSignedState() to be false, but it was true`)
		}
		c.OffChain.LatestSupportedStateTurnNum = MaxTurnNum // Reset so there is no longer a supported state

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
		testhelpers.SignState(&expectedSignedState, &alicePrivateKey)

		if diff := cmp.Diff(expectedSignedState, latestSignedState, cmp.AllowUnexported(expectedSignedState)); diff != "" {
			t.Errorf("LatestSignedState: mismatch (-want +got):\n%s", diff)
		}

		got2 := c.OffChain.SignedStateForTurnNum[1]
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
		got3 := c.OffChain.LatestSupportedStateTurnNum
		want3 := uint64(1)
		if got3 != want3 {
			t.Fatalf(`expected c.LatestSupportedStateTurnNum to be %v, but got %v`, want, got)
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
		testhelpers.SignState(&expectedSignedState, &bobPrivateKey)

		if diff := compareStates(latestSignedState, expectedSignedState); diff != "" {
			t.Errorf("LatestSignedState: mismatch (-want +got):\n%s", diff)
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
	t.Run(`TestAddSignedState`, testAddSignedState)
}

func TestVirtualChannel(t *testing.T) {
	compareChannels := func(a, b *VirtualChannel) string {
		return cmp.Diff(*a, *b, cmp.AllowUnexported(*a, big.Int{}, state.SignedState{}, Channel{}, OnChainData{}, OffChainData{}))
	}

	s := state.TestState.Clone()
	s.Participants = append(s.Participants, s.Participants[0]) // ensure three participants
	s.TurnNum = 0
	testClone := func(t *testing.T) {
		r, err := NewVirtualChannel(s, 0)
		if err != nil {
			t.Fatal(err)
		}
		c := r.Clone()
		if diff := compareChannels(r, c); diff != "" {
			t.Errorf("Clone: mismatch (-want +got):\n%s", diff)
		}

		r.OffChain.LatestSupportedStateTurnNum++
		if reflect.DeepEqual(r.Channel, c.Channel) {
			t.Error("Clone: modifying the clone should not modify the original")
		}

		r.Participants[0] = common.HexToAddress("0x0000000000000000000000000000000000000001")
		if r.Participants[0] == c.Participants[0] {
			t.Error("Clone: modifying the clone should not modify the original")
		}

		var nilChannel *VirtualChannel
		clone := nilChannel.Clone()
		if clone != nil {
			t.Fatal("Tried to clone a Channel via a nil pointer, but got something not nil")
		}
	}

	t.Run(`TestClone SingleHop`, testClone)
	s.Participants = append(s.Participants, s.Participants[1]) // add a fourth participant
	t.Run(`TestClone DoubleHop`, testClone)
}

func TestSerde(t *testing.T) {
	ss := state.NewSignedState(state.TestState)
	signedStateForTurnNum := make(map[uint64]state.SignedState)
	signedStateForTurnNum[0] = ss

	someChannel := Channel{
		Id:        types.Destination{1},
		MyIndex:   1,
		FixedPart: state.TestState.FixedPart(),
		OffChain: OffChainData{
			SignedStateForTurnNum:       signedStateForTurnNum,
			LatestSupportedStateTurnNum: 2,
		},
		OnChain: OnChainData{
			Holdings:  types.Funds{},
			StateHash: common.Hash{},
			Outcome:   outcome.Exit{},
		},
	}

	someChannelJSON := `{"Id":"0x0100000000000000000000000000000000000000000000000000000000000000","MyIndex":1,"Participants":["0xf5a1bb5607c9d079e46d1b3dc33f257d937b43bd","0x760bf27cd45036a6c486802d30b5d90cffbe31fe"],"ChannelNonce":37140676580,"AppDefinition":"0x5e29e5ab8ef33f050c7cc10b5a0456d975c5f88d","ChallengeDuration":60,"OnChain":{"Holdings":{},"Outcome":[],"StateHash":"0x0000000000000000000000000000000000000000000000000000000000000000"},"OffChain":{"SignedStateForTurnNum":{"0":{"State":{"Participants":["0xf5a1bb5607c9d079e46d1b3dc33f257d937b43bd","0x760bf27cd45036a6c486802d30b5d90cffbe31fe"],"ChannelNonce":37140676580,"AppDefinition":"0x5e29e5ab8ef33f050c7cc10b5a0456d975c5f88d","ChallengeDuration":60,"AppData":"","Outcome":[{"Asset":"0x0000000000000000000000000000000000000000","AssetMetadata":{"AssetType":0,"Metadata":""},"Allocations":[{"Destination":"0x000000000000000000000000f5a1bb5607c9d079e46d1b3dc33f257d937b43bd","Amount":5,"AllocationType":0,"Metadata":null},{"Destination":"0x000000000000000000000000ee18ff1575055691009aa246ae608132c57a422c","Amount":5,"AllocationType":0,"Metadata":null}]}],"TurnNum":5,"IsFinal":false},"Sigs":{}}},"LatestSupportedStateTurnNum":2}}`

	// Marshalling
	got, err := json.Marshal(someChannel)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(string(got), someChannelJSON); diff != "" {
		t.Fatalf("incorrect json marshaling (-want +got):\n%s", diff)
	}

	// Unmarshalling
	var c Channel
	err = json.Unmarshal([]byte(someChannelJSON), &c)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(c, someChannel, cmp.AllowUnexported(state.SignedState{})); diff != "" {
		t.Fatalf("incorrect json unmarshaling (-want +got):\n%s", diff)
	}
}
