package protocols

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/types"
)

var chainId, _ = big.NewInt(0).SetString("9001", 10)

var initialState = state.State{ // todo factor into a test fixture
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

func TestNew(t *testing.T) {
	s, _ := NewDirectFundingObjectiveState(initialState, initialState.Participants[0])
	fmt.Println(s)
}
