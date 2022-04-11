package state

import (
	"bytes"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/go-cmp/cmp"
	"github.com/statechannels/go-nitro/types"
)

// The following constants are generated from our ts nitro-protocol package
var correctChannelId = common.HexToHash(`252e8d12b641d8dd3cde74d299047963c99ce52bded2fbd38afa8e07b76a9306`)
var correctStateHash = common.HexToHash(`867a304a6c24962522085dc09c09b47413897a13700e94190c4e3320af0a3779`)
var signerPrivateKey = common.Hex2Bytes(`caab404f975b4620747174a75f08d98b4e5a7053b691b41bcfc0d839d48b7634`)
var signerAddress = common.HexToAddress(`F5A1BB5607C9D079E46d1B3Dc33f257d937b43BD`)
var correctSignature = Signature{
	R: common.Hex2Bytes(`704b3afcc6e702102ca1af3f73cf3b37f3007f368c40e8b81ca823a65740a053`),
	S: common.Hex2Bytes(`14040ad4c598dbb055a50430142a13518e1330b79d24eed86fcbdff1a7a95589`),
	V: byte(0),
}

func TestCloneSignature(t *testing.T) {

	toCopy := Signature{
		R: common.Hex2Bytes(`704b3afcc6e702102ca1af3f73cf3b37f3007f368c40e8b81ca823a65740a053`),
		S: common.Hex2Bytes(`14040ad4c598dbb055a50430142a13518e1330b79d24eed86fcbdff1a7a95589`),
		V: byte(0),
	}

	got := CloneSignature(toCopy)

	// mutate the signature we cloned so we catch a shallow copy
	toCopy.R = common.Hex2Bytes(`0`)
	toCopy.S = common.Hex2Bytes(`0`)
	toCopy.V = byte(1)

	if !bytes.Equal(common.Hex2Bytes(`704b3afcc6e702102ca1af3f73cf3b37f3007f368c40e8b81ca823a65740a053`), got.R) {
		t.Fatalf("Incorrect r param in signature. Got %x, wanted %x", got.R, common.Hex2Bytes(`704b3afcc6e702102ca1af3f73cf3b37f3007f368c40e8b81ca823a65740a053`))
	}
	if !bytes.Equal(common.Hex2Bytes(`14040ad4c598dbb055a50430142a13518e1330b79d24eed86fcbdff1a7a95589`), got.S) {
		t.Fatalf("Incorrect s param in signature. Got %x, wanted %x", got.S, common.Hex2Bytes(`14040ad4c598dbb055a50430142a13518e1330b79d24eed86fcbdff1a7a95589`))
	}
	if byte(0) != got.V {
		t.Fatalf("Incorrect v param in signature. Got %x, wanted %x", got.V, byte(0))
	}
}

func TestChannelId(t *testing.T) {
	want := correctChannelId
	got, error := TestState.ChannelId()
	checkErrorAndTestForEqualBytes(t, error, "channelId", got.Bytes(), want.Bytes())
}

func TestHash(t *testing.T) {
	want := correctStateHash
	got, error := TestState.Hash()
	checkErrorAndTestForEqualBytes(t, error, "state hash", got.Bytes(), want.Bytes())
}

func TestSign(t *testing.T) {
	want_r, want_s, want_v := correctSignature.R, correctSignature.S, correctSignature.V
	got, error := TestState.Sign(signerPrivateKey)
	got_r, got_s, got_v := got.R, got.S, got.V

	if error != nil {
		t.Error(error)
	}
	if !bytes.Equal(want_r, got_r) {
		t.Fatalf("Incorrect r param in signature. Got %x, wanted %x", got_r, want_r)
	}
	if !bytes.Equal(want_s, got_s) {
		t.Fatalf("Incorrect s param in signature. Got %x, wanted %x", got_s, want_s)
	}
	if want_v != got_v {
		t.Fatalf("Incorrect v param in signature. Got %x, wanted %x", got_v, want_v)
	}
}

func TestEqualParticipants(t *testing.T) {
	sameParticipants := []types.Address{
		common.HexToAddress(`0xF5A1BB5607C9D079E46d1B3Dc33f257d937b43BD`),
		common.HexToAddress(`0x760bf27cd45036a6C486802D30B5D90CfFBE31FE`),
	}
	differentParticipants := []types.Address{
		common.HexToAddress(`0x760bf27cd45036a6C486802D30B5D90CfFBE31FE`),
		common.HexToAddress(`0xF5A1BB5607C9D079E46d1B3Dc33f257d937b43BD`),
	}

	if equalParticipants(sameParticipants, TestState.Participants) != true {
		t.Error(`expected equal participants`)
	}

	if equalParticipants(sameParticipants, differentParticipants) == true {
		t.Error(`expected unequal participants`)
	}
}

func TestEqual(t *testing.T) {
	want := State{
		ChainId: chainId,
		Participants: []types.Address{
			common.HexToAddress(`0xF5A1BB5607C9D079E46d1B3Dc33f257d937b43BD`), // private key caab404f975b4620747174a75f08d98b4e5a7053b691b41bcfc0d839d48b7634
			common.HexToAddress(`0x760bf27cd45036a6C486802D30B5D90CfFBE31FE`), // private key 62ecd49c4ccb41a70ad46532aed63cf815de15864bc415c87d507afd6a5e8da2
		},
		ChannelNonce:      big.NewInt(37140676580),
		AppDefinition:     common.HexToAddress(`0x5e29E5Ab8EF33F050c7cc10B5a0456D975C5F88d`),
		ChallengeDuration: big.NewInt(60),
		AppData:           []byte{},
		Outcome:           TestOutcome,
		TurnNum:           5,
		IsFinal:           false,
	}

	got := TestState

	if !got.Equal(want) {
		t.Fatalf(`expected %v to equal %v, but it did not`, got, want)
	}

	want.IsFinal = true

	if got.Equal(want) {
		t.Fatalf(`expected %v to not equal %v, but it did`, got, want)
	}

}

func TestClone(t *testing.T) {

	clone := TestState.Clone()

	if diff := cmp.Diff(TestState, clone); diff != "" {
		t.Fatalf("Clone: mismatch (-want +got):\n%s", diff)
	}

	clone.ChannelNonce.Add(clone.ChannelNonce, big.NewInt(1))
	clone.Outcome[0].Allocations[0].Amount.Add(clone.ChannelNonce, big.NewInt(1))

	if clone.Equal(TestState) {
		t.Fatalf(`expected %v to not equal %v, but it did`, clone, TestState)
	}

	if TestState.ChannelNonce.Cmp(big.NewInt(37140676580)) != 0 || TestState.Outcome[0].Allocations[0].Amount.Cmp(big.NewInt(5)) != 0 {
		t.Fatalf(`State.Clone(): original is modified when clone is modified `)
	}
}

func TestRecoverSigner(t *testing.T) {
	got, error := TestState.RecoverSigner(correctSignature)
	want := signerAddress
	checkErrorAndTestForEqualBytes(t, error, "signer recovered", got.Bytes(), want.Bytes())
}

func checkErrorAndTestForEqualBytes(t *testing.T, err error, descriptor string, got []byte, want []byte) {
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(want, got) {
		t.Fatalf("Incorrect "+descriptor+". Got %x, wanted %x", got, want)
	}
}
