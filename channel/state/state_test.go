package state

import (
	"bytes"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

// The following constants are generated from our ts nitro-protocol package
var correctChannelId = common.HexToHash(`cd578811171bb0291c7dc59081b9a910960b78001dece1954c00da9d8c12a925`)
var correctStateHash = common.HexToHash(`e77a7a870087be40371b8c9b9f50a7c0e410387ce7ea8a22c9c61ea837277420`)
var signerPrivateKey = common.Hex2Bytes(`caab404f975b4620747174a75f08d98b4e5a7053b691b41bcfc0d839d48b7634`)
var signerAddress = common.HexToAddress(`F5A1BB5607C9D079E46d1B3Dc33f257d937b43BD`)
var correctSignature = Signature{
	common.Hex2Bytes(`ae94c67295f65a8de360e80aa50cba3988810239178f03d54e4dd9f95860c182`),
	common.Hex2Bytes(`143fd9bfc0438313fbf50fcaad82be25f1a9d536bf7efe951f20ced3f721694c`),
	byte(0),
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
