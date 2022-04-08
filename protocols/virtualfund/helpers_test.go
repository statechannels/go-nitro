package virtualfund

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/types"
)

type actor struct {
	address     types.Address
	destination types.Destination
	privateKey  []byte
	role        uint
	name        string
}

////////////
// ACTORS //
////////////

var alice = actor{
	address:     common.HexToAddress(`0xD9995BAE12FEe327256FFec1e3184d492bD94C31`),
	destination: types.AddressToDestination(common.HexToAddress(`0xD9995BAE12FEe327256FFec1e3184d492bD94C31`)),
	privateKey:  common.Hex2Bytes(`7ab741b57e8d94dd7e1a29055646bafde7010f38a900f55bbd7647880faa6ee8`),
	role:        0,
	name:        "alice",
}

var p1 = actor{ // Aliases: The Hub, Irene
	address:     common.HexToAddress(`0xd4Fa489Eacc52BA59438993f37Be9fcC20090E39`),
	destination: types.AddressToDestination(common.HexToAddress(`0xd4Fa489Eacc52BA59438993f37Be9fcC20090E39`)),
	privateKey:  common.Hex2Bytes(`2030b463177db2da82908ef90fa55ddfcef56e8183caf60db464bc398e736e6f`),
	role:        1,
	name:        "p1",
}

var bob = actor{
	address:     common.HexToAddress(`0x760bf27cd45036a6C486802D30B5D90CfFBE31FE`),
	destination: types.AddressToDestination(common.HexToAddress(`0x760bf27cd45036a6C486802D30B5D90CfFBE31FE`)),
	privateKey:  common.Hex2Bytes(`62ecd49c4ccb41a70ad46532aed63cf815de15864bc415c87d507afd6a5e8da2`),
	role:        2,
	name:        "bob",
}

var allActors = []actor{alice, p1, bob}

// prepareConsensusChannel prepares a consensus channel with a consensus outcome
//  - allocating 6 to left
//  - allocating 4 to right
//  - including the given guarantees
func prepareConsensusChannel(role uint, left, right actor, guarantees ...consensus_channel.Guarantee) *consensus_channel.ConsensusChannel {
	fp := state.FixedPart{
		ChainId:           big.NewInt(9001),
		Participants:      []types.Address{left.address, right.address},
		ChannelNonce:      big.NewInt(0),
		AppDefinition:     types.Address{},
		ChallengeDuration: big.NewInt(45),
	}

	leftBal := consensus_channel.NewBalance(left.destination, big.NewInt(6))
	rightBal := consensus_channel.NewBalance(right.destination, big.NewInt(4))

	lo := *consensus_channel.NewLedgerOutcome(types.Address{}, leftBal, rightBal, guarantees)

	signedVars := consensus_channel.SignedVars{Vars: consensus_channel.Vars{Outcome: lo, TurnNum: 1}}
	leftSig, err := signedVars.Vars.AsState(fp).Sign(left.privateKey)
	if err != nil {
		panic(err)
	}
	rightSig, err := signedVars.Vars.AsState(fp).Sign(right.privateKey)
	if err != nil {
		panic(err)
	}
	sigs := [2]state.Signature{leftSig, rightSig}

	var cc consensus_channel.ConsensusChannel

	if role == 0 {
		cc, err = consensus_channel.NewLeaderChannel(fp, 1, lo, sigs)
	} else {
		cc, err = consensus_channel.NewFollowerChannel(fp, 1, lo, sigs)
	}
	if err != nil {
		panic(err)
	}

	return &cc
}
