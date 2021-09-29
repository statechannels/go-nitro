package state

import (
	"bytes"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/types"
)

var chainId, _ = big.NewInt(0).SetString("9001", 10)

var state = State{
	ChainId: chainId,
	Participants: []types.Address{
		common.HexToAddress(`0x5e29E5Ab8EF33F050c7cc10B5a0456D975C5F88d`),
		common.HexToAddress(`0xEe18fF1575055691009aa246aE608132C57a422c`),
		common.HexToAddress(`0x95125c394F39bBa29178CAf5F0614EE80CBB1702`),
	},
	ChannelNonce:      big.NewInt(37140676580),
	AppDefinition:     common.HexToAddress(`0x5e29E5Ab8EF33F050c7cc10B5a0456D975C5F88d`),
	ChallengeDuration: big.NewInt(60),
	AppData:           []byte{},
	Outcome:           outcome.Exit{},
	TurnNum:           big.NewInt(5),
	IsFinal:           false,
}

// The following constants are generated from our ts nitro-protocol package
var correctChannelId = common.HexToHash(`b79270eb4cf4d11dcd5cf44c6337b1c9e4730dc3c26f022567bcf2eb63557a72`)
var correctStateHash = common.HexToHash(`3e460a311caf589f1cf80036adfd092d05a30ff72f234918c5cdbc6c4333343a`)
var signerPrivateKey = common.Hex2Bytes(`caab404f975b4620747174a75f08d98b4e5a7053b691b41bcfc0d839d48b7634`)
var signerAddress = common.HexToAddress(`F5A1BB5607C9D079E46d1B3Dc33f257d937b43BD`)
var correctSignature = Signature{
	common.Hex2Bytes(`59d8e91bd182fb4d489bb2d76a6735d494d5bea24e4b51dd95c9d219293312d9`),
	common.Hex2Bytes(`32274a3cec23c31e0c073b3c071cf6e0c21260b0d292a10e6a04257a2d8e87fa`),
	byte(1), // ethers-js gives v:28, which is a legacy representation. and recoveryParam: 1 which corresponds to v here (i.e. it is the normalized version)
} // ethers "joinSignature" gives 0x59d8e91bd182fb4d489bb2d76a6735d494d5bea24e4b51dd95c9d219293312d932274a3cec23c31e0c073b3c071cf6e0c21260b0d292a10e6a04257a2d8e87fa1c

func TestChannelId(t *testing.T) {
	want := correctChannelId
	got, error := state.ChannelId()
	checkErrorAndTestForEqualBytes(t, error, "channelId", got.Bytes(), want.Bytes())
}

func TestHash(t *testing.T) {
	want := correctStateHash
	got, error := state.Hash()
	checkErrorAndTestForEqualBytes(t, error, "state hash", got.Bytes(), want.Bytes())
}

func TestSign(t *testing.T) {
	want_r, want_s, want_v := correctSignature.r, correctSignature.s, correctSignature.v
	got, error := state.Sign(signerPrivateKey)
	got_r, got_s, got_v := got.r, got.s, got.v

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
	got, error := state.RecoverSigner(correctSignature)
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
