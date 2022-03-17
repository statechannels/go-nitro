package store

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/go-cmp/cmp"
	"github.com/statechannels/go-nitro/channel/state"
	nc "github.com/statechannels/go-nitro/crypto"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/protocols/directfund"
)

func TestNewMockStore(t *testing.T) {
	sk := common.Hex2Bytes(`2af069c584758f9ec47c4224a8becc1983f28acfbe837bd7710b70f9fc6d5e44`)
	NewMockStore(sk)
}

func TestSetGetObjective(t *testing.T) {
	sk := common.Hex2Bytes(`2af069c584758f9ec47c4224a8becc1983f28acfbe837bd7710b70f9fc6d5e44`)

	ms := NewMockStore(sk)

	id := protocols.ObjectiveId("404")
	got, err := ms.GetObjectiveById(id)
	if err == nil {
		t.Errorf("expected not to find the %s objective, but found %v", id, got)
	}

	ts := state.TestState
	ts.TurnNum = 0
	request := directfund.ObjectiveRequest{
		MyAddress:         ts.Participants[0],
		CounterParty:      ts.Participants[1],
		AppData:           ts.AppData,
		AppDefinition:     ts.AppDefinition,
		ChallengeDuration: ts.ChallengeDuration,
		Nonce:             ts.ChannelNonce.Int64(),
		Outcome:           ts.Outcome,
	}
	testObj, _ := directfund.NewObjective(request, false)

	if err := ms.SetObjective(&testObj); err != nil {
		t.Errorf("error setting objective %v: %s", testObj, err.Error())
	}

	got, err = ms.GetObjectiveById(testObj.Id())

	if err != nil {
		t.Errorf("expected to find the inserted objective, but didn't: %s", err)
	}

	if got.Id() != testObj.Id() {
		t.Errorf("expected to retrieve same objective Id as was passed in, but didn't")
	}
}

func TestGetObjectiveByChannelId(t *testing.T) {
	sk := common.Hex2Bytes(`2af069c584758f9ec47c4224a8becc1983f28acfbe837bd7710b70f9fc6d5e44`)

	ms := NewMockStore(sk)

	ts := state.TestState
	ts.TurnNum = 0
	request := directfund.ObjectiveRequest{
		MyAddress:         ts.Participants[0],
		CounterParty:      ts.Participants[1],
		AppData:           ts.AppData,
		AppDefinition:     ts.AppDefinition,
		ChallengeDuration: ts.ChallengeDuration,
		Nonce:             ts.ChannelNonce.Int64(),
		Outcome:           ts.Outcome,
	}
	testObj, _ := directfund.NewObjective(request, false)

	if err := ms.SetObjective(&testObj); err != nil {
		t.Errorf("error setting objective %v: %s", testObj, err.Error())
	}

	got, ok := ms.GetObjectiveByChannelId(testObj.C.Id)

	if !ok {
		t.Errorf("expected to find the inserted objective, but didn't")
	}
	if got.Id() != testObj.Id() {
		t.Errorf("expected to retrieve same objective Id as was passed in, but didn't")
	}
	if diff := cmp.Diff(got, &testObj); diff != "" {
		t.Errorf("expected no diff between set and retrieved objective, but found:\n%s", diff)
	}
}

func TestGetChannelSecretKey(t *testing.T) {
	// from state/test-fixtures.go
	sk := common.Hex2Bytes("caab404f975b4620747174a75f08d98b4e5a7053b691b41bcfc0d839d48b7634")
	pk := common.HexToAddress("0xF5A1BB5607C9D079E46d1B3Dc33f257d937b43BD")

	ms := NewMockStore(sk)
	key := ms.GetChannelSecretKey()

	msg := []byte("sign this")

	signedMsg, _ := nc.SignEthereumMessage(msg, *key)
	recoveredSigner, _ := nc.RecoverEthereumMessageSigner(msg, signedMsg)

	if recoveredSigner != pk {
		t.Errorf("expected to recover %x, but got %x", pk, recoveredSigner)
	}
}
