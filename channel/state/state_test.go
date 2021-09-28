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

func TestChannelId(t *testing.T) {

	want := common.HexToHash(`b79270eb4cf4d11dcd5cf44c6337b1c9e4730dc3c26f022567bcf2eb63557a72`) // generated from our ts nitro-protocol package
	got, error := state.ChannelId()

	if error != nil {
		t.Error(error)
	}
	if !bytes.Equal(want.Bytes(), got.Bytes()) {
		t.Errorf("Incorrect channel id. Got %x, wanted %x", got, want)
	}
}

func TestHash(t *testing.T) {

	want := common.HexToHash(`3e460a311caf589f1cf80036adfd092d05a30ff72f234918c5cdbc6c4333343a`) // generated from our ts nitro-protocol package
	got, error := state.Hash()

	if error != nil {
		t.Error(error)
	}
	if !bytes.Equal(want.Bytes(), got.Bytes()) {
		t.Errorf("Incorrect state hash. Got %x, wanted %x", got, want)
	}
}

func TestSign(t *testing.T) {
	privateKey := common.Hex2Bytes(`caab404f975b4620747174a75f08d98b4e5a7053b691b41bcfc0d839d48b7634`)

	want_r := common.Hex2Bytes(`59d8e91bd182fb4d489bb2d76a6735d494d5bea24e4b51dd95c9d219293312d9`)
	want_s := common.Hex2Bytes(`32274a3cec23c31e0c073b3c071cf6e0c21260b0d292a10e6a04257a2d8e87fa`)
	want_v := byte(1) // ethers-js gives v:28, which is a legacy representation. and recoveryParam: 1 which corresponds to v here (i.e. it is the normalized version)
	// ethers "joinSignature" gives 0x59d8e91bd182fb4d489bb2d76a6735d494d5bea24e4b51dd95c9d219293312d932274a3cec23c31e0c073b3c071cf6e0c21260b0d292a10e6a04257a2d8e87fa1c

	got, error := state.Sign(privateKey)
	got_r, got_s, got_v := SplitSignature(got)

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
