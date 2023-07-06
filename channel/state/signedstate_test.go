package state

import (
	"encoding/json"
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/go-cmp/cmp"
)

func TestMergeWithDuplicateSignatures(t *testing.T) {
	// ss1 has only alice's signature
	ss1 := NewSignedState(TestState)
	sigA, _ := TestState.Sign(common.Hex2Bytes(`caab404f975b4620747174a75f08d98b4e5a7053b691b41bcfc0d839d48b7634`))
	_ = ss1.AddSignature(sigA)

	// ss2 has alice and bob's signatures
	sigB, _ := TestState.Sign(common.Hex2Bytes(`62ecd49c4ccb41a70ad46532aed63cf815de15864bc415c87d507afd6a5e8da2`))
	ss2 := NewSignedState(TestState)
	_ = ss2.AddSignature(sigA)
	_ = ss2.AddSignature(sigB)

	err := ss1.Merge(ss2)
	if err != nil {
		t.Error(err)
	}

	got := ss1
	want := SignedState{
		TestState,
		map[uint]Signature{
			0: sigA,
			1: sigB,
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("incorrect merge:\ngot\n\t%v,\nwanted\n\t%v", got, want)
	}
}

func TestMerge(t *testing.T) {
	ss1 := NewSignedState(TestState)
	sigA, _ := TestState.Sign(common.Hex2Bytes(`caab404f975b4620747174a75f08d98b4e5a7053b691b41bcfc0d839d48b7634`))
	_ = ss1.AddSignature(sigA)

	ss2 := NewSignedState(TestState)
	sigB, _ := TestState.Sign(common.Hex2Bytes(`62ecd49c4ccb41a70ad46532aed63cf815de15864bc415c87d507afd6a5e8da2`))
	_ = ss2.AddSignature(sigB)

	err := ss1.Merge(ss2)
	if err != nil {
		t.Error(err)
	}

	got := ss1
	want := SignedState{
		TestState,
		map[uint]Signature{
			0: sigA,
			1: sigB,
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf(`incorrect merge, got %v, wanted %v`, got, want)
	}
}

func TestJSON(t *testing.T) {
	ss1 := NewSignedState(TestState)
	sigA, _ := TestState.Sign(common.Hex2Bytes(`caab404f975b4620747174a75f08d98b4e5a7053b691b41bcfc0d839d48b7634`))
	_ = ss1.AddSignature(sigA)

	msgString := `{"State":{"Participants":["0xf5a1bb5607c9d079e46d1b3dc33f257d937b43bd","0x760bf27cd45036a6c486802d30b5d90cffbe31fe"],"ChannelNonce":37140676580,"AppDefinition":"0x5e29e5ab8ef33f050c7cc10b5a0456d975c5f88d","ChallengeDuration":60,"AppData":"","Outcome":[{"Asset":"0x0000000000000000000000000000000000000000","AssetMetadata":{"AssetType":0,"Metadata":""},"Allocations":[{"Destination":"0x000000000000000000000000f5a1bb5607c9d079e46d1b3dc33f257d937b43bd","Amount":5,"AllocationType":0,"Metadata":null},{"Destination":"0x000000000000000000000000ee18ff1575055691009aa246ae608132c57a422c","Amount":5,"AllocationType":0,"Metadata":null}]}],"TurnNum":5,"IsFinal":false},"Sigs":{"0":"0x2873c05a6ecebca3d2fde93ea3332b4423bbd2a60a973ec61b8abfe16294f69e59ddb5b57c79924b57db6d928a0f3ef657c6e0e9879d2bf387364ace6d3284fd1b"}}`

	t.Run(`TestMarshalJSON`, func(t *testing.T) {
		got, err := ss1.MarshalJSON()
		if err != nil {
			t.Error(err)
		}
		want := msgString
		if string(got) != want {
			t.Fatalf("incorrect MarshalJSON, got\n\t%v\nwanted\n\t%v", string(got), want)
		}
	})

	t.Run(`TestUnmarshalJSON`, func(t *testing.T) {
		got := SignedState{}
		err := json.Unmarshal([]byte(msgString), &got)
		if err != nil {
			t.Error(err)
		}
		want := ss1

		if !reflect.DeepEqual(got, ss1) {
			t.Errorf(`incorrect UnmarshalJSON, got %v, wanted %v`, got, want)
		}
	})
}

func TestSignedStateClone(t *testing.T) {
	compareStates := func(a, b SignedState) string {
		return cmp.Diff(a, b, cmp.AllowUnexported(a, big.Int{}))
	}

	ss1 := NewSignedState(TestState)
	sigA, _ := TestState.Sign(common.Hex2Bytes(`caab404f975b4620747174a75f08d98b4e5a7053b691b41bcfc0d839d48b7634`))
	_ = ss1.AddSignature(sigA)

	clone := ss1.Clone()

	if diff := compareStates(ss1, clone); diff != "" {
		t.Errorf("Clone: mismatch (-want +got):\n%s", diff)
	}
}

func TestSignatureGetters(t *testing.T) {
	ss := NewSignedState(TestState)
	sigA, _ := TestState.Sign(common.Hex2Bytes(`caab404f975b4620747174a75f08d98b4e5a7053b691b41bcfc0d839d48b7634`))
	_ = ss.AddSignature(sigA)

	got, err := ss.GetParticipantSignature(0)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, sigA) {
		t.Errorf("incorrect GetParticipantSignature, got %v, wanted %v", got, sigA)
	}

	if ss.HasAllSignatures() != false {
		t.Errorf("incorrect HasAllSignatures, expected false but got true")
	}

	expectedSigs := make([]Signature, len(ss.state.Participants))
	expectedSigs[0] = sigA
	gotSigs := ss.Signatures()

	if !reflect.DeepEqual(gotSigs, expectedSigs) {
		t.Errorf("incorrect Signatures, got %v, wanted %v", gotSigs, expectedSigs)
	}
}
