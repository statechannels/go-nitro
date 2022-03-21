package td

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/types"
)

var chainId, _ = big.NewInt(0).SetString("9001", 10)

var testOutcome = outcome.Exit{
	outcome.SingleAssetExit{
		Asset: types.Address{}, // the native token of the chain
		Allocations: outcome.Allocations{
			outcome.Allocation{
				Destination: Actors.Alice.Destination(),
				Amount:      big.NewInt(6),
			},
			outcome.Allocation{
				Destination: Actors.Bob.Destination(),
				Amount:      big.NewInt(4),
			},
		},
	},
}

var testState = state.State{
	ChainId: chainId,
	Participants: []types.Address{
		Actors.Alice.Address,
		Actors.Bob.Address,
	},
	ChannelNonce:      big.NewInt(37140676580),
	AppDefinition:     common.HexToAddress(`0x5e29E5Ab8EF33F050c7cc10B5a0456D975C5F88d`),
	ChallengeDuration: big.NewInt(60),
	AppData:           []byte{},
	Outcome:           testOutcome,
	TurnNum:           5,
	IsFinal:           false,
}
