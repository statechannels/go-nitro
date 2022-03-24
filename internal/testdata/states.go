package testdata

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/types"
)

type outcomes struct {
	// Create returns a simple outcome {a: aBalance, b: bBalance} in the
	// zero-asset (chain-native token)
	Create func(a, b types.Address, aBalance, bBalance uint) outcome.Exit
}

var Outcomes outcomes = outcomes{
	Create: createOutcome,
}

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

// createOutcome is a helper function to create a two-actor outcome
func createOutcome(first types.Address, second types.Address, x, y uint) outcome.Exit {

	return outcome.Exit{outcome.SingleAssetExit{
		Allocations: outcome.Allocations{
			outcome.Allocation{
				Destination: types.AddressToDestination(first),
				Amount:      big.NewInt(int64(x)),
			},
			outcome.Allocation{
				Destination: types.AddressToDestination(second),
				Amount:      big.NewInt(int64(y)),
			},
		},
	}}
}
