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
	got, ok := ms.GetObjectiveById("404")
	if ok {
		t.Errorf("expected not to find the 404 objective, but found %v", got)
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

	ms.SetObjective(testObj)
	got, ok = ms.GetObjectiveById(testObj.Id())

	if !ok {
		t.Errorf("expected to find the inserted objective, but didn't")
	}
	if got.Id() != testObj.Id() {
		t.Errorf("expected to retrieve same objective Id as was passed in, but didn't")
	}
}

func TestGetObjectiveByChannelId(t *testing.T) {
	// todo
}

// BUG(geoknee)
func TestGetChannelSecretKey(t *testing.T) {
	// from state/test-fixtures.go
	sk := common.Hex2Bytes("caab404f975b4620747174a75f08d98b4e5a7053b691b41bcfc0d839d48b7634")
	pk := common.HexToAddress("0xF5A1BB5607C9D079E46d1B3Dc33f257d937b43BD")

	ms := NewMockStore(sk)
	key := ms.GetChannelSecretKey()

	signedMsg, _ := state.SignEthereumMessage([]byte("asdfasdf"), *key)

	recoveredSigner, _ := state.RecoverEthereumMessageSigner([]byte("asdfasdf"), signedMsg)

	if recoveredSigner != pk {
		t.Errorf("expected to recover %x, but got %x", pk, recoveredSigner)
	}
}
