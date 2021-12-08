package virtualfund

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/types"
)

// TODO bury these variables in a TestData struct?

// In general
// Alice = P_0 <=L_0=> P_1 <=L_1=> ... P_n <=L_n>= P_n+1 = Bob

// For these tests
// Alice <=L_0=> P_1 <=L_1=> Bob

////////////
// ACTORS //
////////////

var Alice = struct {
	address     types.Address
	destination types.Destination
	privateKey  []byte
}{
	address:     common.HexToAddress(`0xD9995BAE12FEe327256FFec1e3184d492bD94C31`),
	destination: types.AdddressToDestination(common.HexToAddress(`0xD9995BAE12FEe327256FFec1e3184d492bD94C31`)),
	privateKey:  common.Hex2Bytes(`7ab741b57e8d94dd7e1a29055646bafde7010f38a900f55bbd7647880faa6ee8`),
}

var P_1 = struct { // Aliases: The Hub, Irene
	address     types.Address
	destination types.Destination
	privateKey  []byte
}{
	address:     common.HexToAddress(`0xd4Fa489Eacc52BA59438993f37Be9fcC20090E39`),
	destination: types.AdddressToDestination(common.HexToAddress(`0xd4Fa489Eacc52BA59438993f37Be9fcC20090E39`)),
	privateKey:  common.Hex2Bytes(`2030b463177db2da82908ef90fa55ddfcef56e8183caf60db464bc398e736e6f`),
}

var Bob = struct {
	address     types.Address
	destination types.Destination
	privateKey  []byte
}{
	address:     common.HexToAddress(`0x760bf27cd45036a6C486802D30B5D90CfFBE31FE`),
	destination: types.AdddressToDestination(common.HexToAddress(`0x760bf27cd45036a6C486802D30B5D90CfFBE31FE`)),
	privateKey:  common.Hex2Bytes(`62ecd49c4ccb41a70ad46532aed63cf815de15864bc415c87d507afd6a5e8da2`),
}

/////////////////////
// VIRTUAL CHANNEL //
/////////////////////

// Virtual Channel
var VPreFund = state.State{
	ChainId:           big.NewInt(9001),
	Participants:      []types.Address{Alice.address, P_1.address, Bob.address}, // A single hop virtual channel
	ChannelNonce:      big.NewInt(0),
	AppDefinition:     types.Address{},
	ChallengeDuration: big.NewInt(45),
	AppData:           []byte{},
	Outcome: outcome.Exit{outcome.SingleAssetExit{
		Allocations: outcome.Allocations{
			outcome.Allocation{
				Destination: Alice.destination,
				Amount:      big.NewInt(5),
			},
			outcome.Allocation{
				Destination: Bob.destination,
				Amount:      big.NewInt(5),
			},
		},
	}},
	TurnNum: big.NewInt(0),
	IsFinal: false,
}

/////////////////////
// LEDGER CHANNELS //
/////////////////////

var L_0state = state.State{
	ChainId:           big.NewInt(9001),
	Participants:      []types.Address{Alice.address, P_1.address},
	ChannelNonce:      big.NewInt(0),
	AppDefinition:     types.Address{},
	ChallengeDuration: big.NewInt(45),
	AppData:           []byte{},
	Outcome: outcome.Exit{outcome.SingleAssetExit{
		Allocations: outcome.Allocations{
			outcome.Allocation{
				Destination: Alice.destination,
				Amount:      big.NewInt(5),
			},
			outcome.Allocation{
				Destination: P_1.destination,
				Amount:      big.NewInt(5),
			},
		},
	}},
	TurnNum: big.NewInt(0), // We use turnNum 0 so that we can use github.com/statechannels/go-nitro/channel.New().
	// It would be more realistic to have a higher TurnNum, but that would involve more boilerplate code.
	IsFinal: false,
}

var VId, _ = VPreFund.ChannelId()

var L_0guaranteemetadataencoded, _ = outcome.GuaranteeMetadata{
	Left:  Alice.destination,
	Right: P_1.destination,
}.Encode()

var L_0updatedstate = state.State{
	ChainId:           big.NewInt(9001),
	Participants:      []types.Address{Alice.address, P_1.address},
	ChannelNonce:      big.NewInt(0),
	AppDefinition:     types.Address{},
	ChallengeDuration: big.NewInt(45),
	AppData:           []byte{},
	Outcome: outcome.Exit{outcome.SingleAssetExit{
		Allocations: outcome.Allocations{
			outcome.Allocation{
				Destination: Alice.destination,
				Amount:      big.NewInt(0),
			},
			outcome.Allocation{
				Destination: P_1.destination,
				Amount:      big.NewInt(0),
			},
			outcome.Allocation{
				Destination:    VId,
				Amount:         big.NewInt(10),
				AllocationType: outcome.GuaranteeAllocationType,
				Metadata:       L_0guaranteemetadataencoded,
			},
		},
	}},
	TurnNum: big.NewInt(2), // This needs to be greater than the previous state else it will be rejected by Channel.AddSignedState
	IsFinal: false,
}
