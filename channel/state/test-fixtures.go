package state

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/types"
)

var chainId, _ = big.NewInt(0).SetString("9001", 10)

var TestState = State{
	ChainId: chainId,
	Participants: []types.Address{
		common.HexToAddress(`0xF5A1BB5607C9D079E46d1B3Dc33f257d937b43BD`), // private key caab404f975b4620747174a75f08d98b4e5a7053b691b41bcfc0d839d48b7634
		common.HexToAddress(`0xEe18fF1575055691009aa246aE608132C57a422c`),
		common.HexToAddress(`0x95125c394F39bBa29178CAf5F0614EE80CBB1702`),
	},
	ChannelNonce:      big.NewInt(37140676580),
	AppDefinition:     common.HexToAddress(`0x5e29E5Ab8EF33F050c7cc10B5a0456D975C5F88d`),
	ChallengeDuration: big.NewInt(60),
	AppData:           []byte{},
	Outcome: outcome.Exit{
		outcome.SingleAssetExit{
			Asset: types.Address{},
			Allocations: outcome.Allocations{
				outcome.Allocation{
					Destination: types.AdddresstoDestination(common.HexToAddress(`0xF5A1BB5607C9D079E46d1B3Dc33f257d937b43BD`)),
					Amount:      big.NewInt(5),
				},
				outcome.Allocation{
					Destination: types.AdddresstoDestination(common.HexToAddress(`0xEe18fF1575055691009aa246aE608132C57a422c`)),
					Amount:      big.NewInt(5),
				},
			},
		},
	},

	TurnNum: big.NewInt(5),
	IsFinal: false,
}
