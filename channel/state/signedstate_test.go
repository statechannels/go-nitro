package state

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestSignedStateEqual(t *testing.T) {
	sigA, _ := TestState.Sign(common.Hex2Bytes(`caab404f975b4620747174a75f08d98b4e5a7053b691b41bcfc0d839d48b7634`))

	ss1 := NewSignedState(TestState)
	ss1.AddSignature(sigA)
	ss2 := NewSignedState(TestState)
	ss2.AddSignature(sigA)

	if !ss1.Equal(ss2) {
		t.Errorf(`expected %v to Equal %v, but it did not`, ss1, ss2)
	}
}

func TestMerge(t *testing.T) {

	ss1 := NewSignedState(TestState)
	sigA, _ := TestState.Sign(common.Hex2Bytes(`caab404f975b4620747174a75f08d98b4e5a7053b691b41bcfc0d839d48b7634`))
	ss1.AddSignature(sigA)

	ss2 := NewSignedState(TestState)
	sigB, _ := TestState.Sign(common.Hex2Bytes(`62ecd49c4ccb41a70ad46532aed63cf815de15864bc415c87d507afd6a5e8da2`))
	ss2.AddSignature(sigB)

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

	if !got.Equal(want) {
		t.Errorf(`incorrect merge, got %v, wanted %v`, got, want)
	}

}

// func TestMarshalJSON(t *testing.T) {
