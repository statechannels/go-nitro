package state

import (
	"bytes"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/types"
)

func TestChannelId(t *testing.T) {

	chainId, _ := big.NewInt(0).SetString("9001", 10)

	channelPart := ChannelPart{
		chainId,
		[]types.Address{
			common.HexToAddress(`0x5e29E5Ab8EF33F050c7cc10B5a0456D975C5F88d`),
			common.HexToAddress(`0xEe18fF1575055691009aa246aE608132C57a422c`),
			common.HexToAddress(`0x95125c394F39bBa29178CAf5F0614EE80CBB1702`),
		}, big.NewInt(37140676580)}

	want := common.HexToHash(`0xb79270eb4cf4d11dcd5cf44c6337b1c9e4730dc3c26f022567bcf2eb63557a72`) // generated from our ts nitro-protocol package
	got, error := channelPart.ChannelId()

	if error != nil {
		t.Error(error)
	}
	if !bytes.Equal(want.Bytes(), got.Bytes()) {
		t.Errorf("Incorrect channel id. Got %x, wanted %x", got, want)
	}

}
