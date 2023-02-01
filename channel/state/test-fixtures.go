package state

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/types"
)

var TestOutcome = outcome.Exit{
	outcome.SingleAssetExit{
		Asset: types.Address{},
		AssetMetadata: outcome.AssetMetadata{
			AssetType: 0,
			Metadata:  []byte{0x00},
		},
		Allocations: outcome.Allocations{
			outcome.Allocation{
				Destination: types.AddressToDestination(common.HexToAddress(`0xF5A1BB5607C9D079E46d1B3Dc33f257d937b43BD`)),
				Amount:      big.NewInt(5),
			},
			outcome.Allocation{
				Destination: types.AddressToDestination(common.HexToAddress(`0xEe18fF1575055691009aa246aE608132C57a422c`)),
				Amount:      big.NewInt(5),
			},
		},
	},
}

var TestState = State{
	Participants: []types.Address{
		common.HexToAddress(`0xF5A1BB5607C9D079E46d1B3Dc33f257d937b43BD`), // private key caab404f975b4620747174a75f08d98b4e5a7053b691b41bcfc0d839d48b7634
		common.HexToAddress(`0x760bf27cd45036a6C486802D30B5D90CfFBE31FE`), // private key 62ecd49c4ccb41a70ad46532aed63cf815de15864bc415c87d507afd6a5e8da2
	},
	ChannelNonce:      37140676580,
	AppDefinition:     common.HexToAddress(`0x5e29E5Ab8EF33F050c7cc10B5a0456D975C5F88d`),
	ChallengeDuration: 60,
	AppData:           []byte{},
	Outcome:           TestOutcome,
	TurnNum:           5,
	IsFinal:           false,
}
