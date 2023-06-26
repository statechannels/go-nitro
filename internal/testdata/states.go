package testdata

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/internal/testactors"
	"github.com/statechannels/go-nitro/types"
)

// SimpleItem is an ergonomic way to create a simple allocation item, without having to create a big.Int
type SimpleItem struct {
	Dest   types.Destination
	Amount int64
}

type outcomes struct {
	// Create returns a simple outcome {a: aBalance, b: bBalance} with the supplied
	// erc20 token address as the asset. A zero address implies the
	// zero-asset (chain-native token)
	Create func(a, b types.Address, aBalance, bBalance uint, token common.Address) outcome.Exit
	// CreateLongOutcome returns a simple outcome {addressOne: balanceOne ...} in the
	// zero-asset (chain-native token)
	// The outcome can be of arbitrary length, and is formed in order that the SimpleItems are provided
	CreateLongOutcome func(items ...SimpleItem) outcome.Exit
}

var Outcomes outcomes = outcomes{
	Create:            createOutcome,
	CreateLongOutcome: createLongOutcome,
}

var someAppDefinition = common.HexToAddress(`0x5e29E5Ab8EF33F050c7cc10B5a0456D975C5F88d`)

var testOutcome = createLongOutcome(
	SimpleItem{testactors.Alice.Destination(), 6},
	SimpleItem{testactors.Bob.Destination(), 6},
)

var testVirtualState = state.State{
	Participants: []types.Address{
		testactors.Alice.Address(),
		testactors.Irene.Address(),
		testactors.Bob.Address(),
	},
	ChannelNonce:      1234789,
	AppDefinition:     someAppDefinition,
	ChallengeDuration: 60,
	AppData:           []byte{},
	Outcome: Outcomes.CreateLongOutcome(
		SimpleItem{testactors.Alice.Destination(), 6},
		SimpleItem{testactors.Bob.Destination(), 4},
	),
	TurnNum: 0,
	IsFinal: false,
}

var testState = state.State{
	Participants: []types.Address{
		testactors.Alice.Address(),
		testactors.Bob.Address(),
	},
	ChannelNonce:      37140676580,
	AppDefinition:     someAppDefinition,
	ChallengeDuration: 60,
	AppData:           []byte{},
	Outcome:           testOutcome,
	TurnNum:           5,
	IsFinal:           false,
}

func createLedgerState(leader, follower types.Address, clientBalance, hubBalance uint) state.State {
	state := testState.Clone()
	state.Participants = []types.Address{
		leader,
		follower,
	}
	state.Outcome = Outcomes.Create(leader, follower, clientBalance, hubBalance, common.Address{})
	state.AppDefinition = types.Address{} // ledger channel running the consensus app
	state.TurnNum = 0

	return state
}

// createOutcome is a helper function to create a two-actor outcome
func createOutcome(first types.Address, second types.Address, x, y uint, asset common.Address) outcome.Exit {
	return outcome.Exit{outcome.SingleAssetExit{
		Asset: asset,
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

// createLongOutcome creates an outcome of arbitrary length, in the order of the provided SimpleItems
func createLongOutcome(items ...SimpleItem) outcome.Exit {
	sae := outcome.SingleAssetExit{}
	for _, i := range items {
		a := outcome.Allocation{
			Destination: i.Dest,
			Amount:      big.NewInt(i.Amount),
		}
		sae.Allocations = append(sae.Allocations, a)
	}

	return outcome.Exit{sae}
}
