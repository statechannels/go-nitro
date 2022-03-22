package state

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-test/deep"
)

func TestSignedStateEqual(t *testing.T) {
	sigA, _ := TestState.Sign(common.Hex2Bytes(`caab404f975b4620747174a75f08d98b4e5a7053b691b41bcfc0d839d48b7634`))

	ss1 := NewSignedState(TestState)
	_ = ss1.AddSignature(sigA)
	ss2 := NewSignedState(TestState)
	_ = ss2.AddSignature(sigA)

	if !reflect.DeepEqual(ss1, ss2) {
		t.Errorf(`expected %v to Equal %v, but it did not`, ss1, ss2)
	}
}
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
		t.Errorf(`incorrect merge, got %v, wanted %v`, got, want)
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

	msgString := `{"State":{"ChainId":9001,"Participants":["0xf5a1bb5607c9d079e46d1b3dc33f257d937b43bd","0x760bf27cd45036a6c486802d30b5d90cffbe31fe"],"ChannelNonce":37140676580,"AppDefinition":"0x5e29e5ab8ef33f050c7cc10b5a0456d975c5f88d","ChallengeDuration":60,"AppData":"","Outcome":[{"Asset":"0x0000000000000000000000000000000000000000","Metadata":null,"Allocations":[{"Destination":[0,0,0,0,0,0,0,0,0,0,0,0,245,161,187,86,7,201,208,121,228,109,27,61,195,63,37,125,147,123,67,189],"Amount":5,"AllocationType":0,"Metadata":null},{"Destination":[0,0,0,0,0,0,0,0,0,0,0,0,238,24,255,21,117,5,86,145,0,154,162,70,174,96,129,50,197,122,66,44],"Amount":5,"AllocationType":0,"Metadata":null}]}],"TurnNum":5,"IsFinal":false},"Sigs":{"0":{"R":"cEs6/MbnAhAsoa8/c887N/MAfzaMQOi4HKgjpldAoFM=","S":"FAQK1MWY27BVpQQwFCoTUY4TMLedJO7Yb8vf8aepVYk=","V":0}}}`

	t.Run(`TestMarshalJSON`, func(t *testing.T) {
		got, err := ss1.MarshalJSON()
		if err != nil {
			t.Error(err)
		}
		want := msgString
		if string(got) != want {
			t.Errorf(`incorrect MarshalJSON, got %v, wanted %v`, string(got), want)
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
	ss1 := NewSignedState(TestState)
	sigA, _ := TestState.Sign(common.Hex2Bytes(`caab404f975b4620747174a75f08d98b4e5a7053b691b41bcfc0d839d48b7634`))
	_ = ss1.AddSignature(sigA)

	clone := ss1.Clone()

	if diff := deep.Equal(ss1, clone); diff != nil {
		t.Errorf("Clone: mismatch (-want +got):\n%s", diff)
	}

}
