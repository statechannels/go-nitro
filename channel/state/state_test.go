package state

import (
	"bytes"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/types"
)

// The following constants are generated from our ts nitro-protocol package
var correctChannelId = common.HexToHash(`4ebd366d014a173765ba1e50f284c179ade31f20441bec41664712aac6cc461d`)
var correctStateHash = common.HexToHash(`c8d5eae9ca84647bafc1bd26a7058a230cd45cb3bf21b37b6330053f4e3ebd0e`)
var signerPrivateKey = common.Hex2Bytes(`caab404f975b4620747174a75f08d98b4e5a7053b691b41bcfc0d839d48b7634`)
var signerAddress = common.HexToAddress(`F5A1BB5607C9D079E46d1B3Dc33f257d937b43BD`)
var correctSignature = Signature{
	common.Hex2Bytes(`b3b69fbfbdcb3100d6e5758c5661d0d793bc227716d16fd6235ccd588cae2849`),
	common.Hex2Bytes(`500969f691a848245910e9ac7688bbc28198b6a6e723299751bda6234bff77f3`),
	byte(1),
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
		t.Errorf("Incorrect r param in signature. Got %x, wanted %x", got_r, want_r)
	}
	if !bytes.Equal(want_s, got_s) {
		t.Errorf("Incorrect s param in signature. Got %x, wanted %x", got_s, want_s)
	}
	if want_v != got_v {
		t.Errorf("Incorrect v param in signature. Got %x, wanted %x", got_v, want_v)
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
		TurnNum:           big.NewInt(5),
		IsFinal:           false,
	}

	got := TestState

	if !got.Equal(want) {
		t.Errorf(`expected %v to equal %v, but it did not`, got, want)
	}

	want.IsFinal = true

	if got.Equal(want) {
		t.Errorf(`expected %v to not equal %v, but it did`, got, want)
	}

}

func TestClone(t *testing.T) {

	clone := TestState.Clone()

	if !clone.Equal(TestState) {
		t.Errorf(`expected %v to equal %v, but it did not`, clone, TestState)
	}

	clone.ChannelNonce.Add(clone.ChannelNonce, big.NewInt(1))

	if clone.Equal(TestState) {
		t.Errorf(`expected %v to not equal %v, but it did`, clone, TestState)
	}

	if TestState.ChannelNonce.Cmp(big.NewInt(2)) == 0 {
		t.Errorf(`original.Clone() is modified when original is modified `)
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
		t.Errorf("Incorrect "+descriptor+". Got %x, wanted %x", got, want)
	}
}
