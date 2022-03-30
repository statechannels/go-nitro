package testdata

import (
	"bytes"
	"math/big"
	"sort"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/types"
)

type outcomes struct {
	// Create returns a simple outcome {a: aBalance, b: bBalance} in the
	// zero-asset (chain-native token)
	Create func(a, b types.Address, aBalance, bBalance uint) outcome.Exit
	// CreateFromMap returns a simple outcome {addressOne: balanceOne ...} in the
	// zero-asset (chain-native token)
	CreateFromMap func(map[types.Address]uint) outcome.Exit
}

var Outcomes outcomes = outcomes{
	Create:        createOutcome,
	CreateFromMap: createOutcomeFromMap,
}

var chainId, _ = big.NewInt(0).SetString("9001", 10)
var someAppDefinition = common.HexToAddress(`0x5e29E5Ab8EF33F050c7cc10B5a0456D975C5F88d`)

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

var testVirtualState = state.State{
	ChainId: chainId,
	Participants: []types.Address{
		Actors.Alice.Address,
		Actors.Irene.Address,
		Actors.Bob.Address,
	},
	ChannelNonce:      big.NewInt(1234789),
	AppDefinition:     someAppDefinition,
	ChallengeDuration: big.NewInt(60),
	AppData:           []byte{},
	Outcome: Outcomes.CreateFromMap(map[types.Address]uint{
		Actors.Alice.Address: 6,
		Actors.Bob.Address:   4,
	}),
	TurnNum: 0,
	IsFinal: false,
}

var testState = state.State{
	ChainId: chainId,
	Participants: []types.Address{
		Actors.Alice.Address,
		Actors.Bob.Address,
	},
	ChannelNonce:      big.NewInt(37140676580),
	AppDefinition:     someAppDefinition,
	ChallengeDuration: big.NewInt(60),
	AppData:           []byte{},
	Outcome:           testOutcome,
	TurnNum:           5,
	IsFinal:           false,
}

func createLedgerState(client, hub types.Address, clientBalance, hubBalance uint) state.State {
	state := testState.Clone()
	state.Participants = []types.Address{
		client,
		hub,
	}
	state.Outcome = Outcomes.Create(client, hub, clientBalance, hubBalance)
	state.AppDefinition = types.Address{} // ledger channel running the consensus app
	state.TurnNum = 0                     // a requirement for channel.NewTwoPartyLedger

	return state
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

func createOutcomeFromMap(amounts map[types.Address]uint) outcome.Exit {
	var allocations []outcome.Allocation

	// Generate a list of addresses from the map
	addresses := make([]types.Address, 0, len(amounts))
	for address, _ := range amounts {
		addresses = append(addresses, address)
	}

	// Sort the addresses to ensure the order is deterministic
	sort.Slice(addresses, func(i, j int) bool {
		return bytes.Compare(addresses[i].Bytes(), addresses[j].Bytes()) < 0
	})

	// Create the allocations
	for _, address := range addresses {
		allocations = append(allocations, outcome.Allocation{
			Destination: types.AddressToDestination(address),
			Amount:      big.NewInt(int64((amounts[address]))),
		})
	}
	return outcome.Exit{outcome.SingleAssetExit{
		Allocations: allocations,
	}}
}
