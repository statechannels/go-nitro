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

	want := common.HexToHash(`0xb79270eb4cf4d11dcd5cf44c6337b1c9e4730dc3c26f022567bcf2eb63557a72`) // generated from our ts nitro-protocol package
	got, error := state.ChannelId()

	if error != nil {
		t.Error(error)
	}
	if !bytes.Equal(want.Bytes(), got.Bytes()) {
		t.Errorf("Incorrect channel id. Got %x, wanted %x", got, want)
	}
}

func TestHash(t *testing.T) {

	want := common.HexToHash(`0x3e460a311caf589f1cf80036adfd092d05a30ff72f234918c5cdbc6c4333343a`) // generated from our ts nitro-protocol package
	got, error := state.Hash()

	if error != nil {
		t.Error(error)
	}
	if !bytes.Equal(want.Bytes(), got.Bytes()) {
		t.Errorf("Incorrect state hash. Got %x, wanted %x", got, want)
	}
}
