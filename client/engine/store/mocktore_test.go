package store

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/channel/state"
	directfund "github.com/statechannels/go-nitro/protocols/direct-fund"
	"github.com/statechannels/go-nitro/types"
)

func TestNewMockStore(t *testing.T) {
	NewMockStore([]byte{'a', 'b', 'c'})
}

func TestSetGetObjective(t *testing.T) {
	ms := NewMockStore([]byte{})
        id := "404"
	got, ok := ms.GetObjectiveById(id)
	if ok {
		t.Errorf("expected not to find the %s objective, but found %v", id, got)
	}

	ts := state.TestState
	ts.TurnNum = 0

	testObj, _ := directfund.New(
		ts,
		ts.Participants[0],
		false,
		types.AddressToDestination(ts.Participants[0]),
		types.AddressToDestination(ts.Participants[1]),
	)

	if err := ms.SetObjective(testObj); err != nil {
		t.Errorf("error setting objective %v: %s", testObj, err.Error())
	}

	got, ok = ms.GetObjectiveById(testObj.Id())

	if !ok {
		t.Errorf("expected to find the inserted objective, but didn't")
	}
	if got.Id() != testObj.Id() {
		t.Errorf("expected to retrieve same objective Id as was passed in, but didn't")
	}
}

func TestGetObjectiveByChannelId(t *testing.T) {
	ms := NewMockStore([]byte{})

	ts := state.TestState
	ts.TurnNum = 0

	testObj, _ := directfund.New(
		ts,
		ts.Participants[0],
		false,
		types.AddressToDestination(ts.Participants[0]),
		types.AddressToDestination(ts.Participants[1]),
	)

	if err := ms.SetObjective(testObj); err != nil {
		t.Errorf("error setting objective %v: %s", testObj, err.Error())
	}

	got, ok := ms.GetObjectiveByChannelId(testObj.C.Id)

	if !ok {
		t.Errorf("expected to find the inserted objective, but didn't")
	}
	if got.Id() != testObj.Id() {
		t.Errorf("expected to retrieve same objective Id as was passed in, but didn't")
	}
}

func TestGetChannelSecretKey(t *testing.T) {
	// from state/test-fixtures.go
	sk := common.Hex2Bytes("caab404f975b4620747174a75f08d98b4e5a7053b691b41bcfc0d839d48b7634")
	pk := common.HexToAddress("0xF5A1BB5607C9D079E46d1B3Dc33f257d937b43BD")

	ms := NewMockStore(sk)
	key := ms.GetChannelSecretKey()

	msg := []byte("sign this")

	signedMsg, _ := state.SignEthereumMessage(msg, *key)

	recoveredSigner, _ := state.RecoverEthereumMessageSigner(msg, signedMsg)

	if recoveredSigner != pk {
		t.Errorf("expected to recover %x, but got %x", pk, recoveredSigner)
	}
}
