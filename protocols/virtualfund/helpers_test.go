package virtualfund

import (
	"math/big"

	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/types"
)

type actor struct {
	address     types.Address
	destination types.Destination
	privateKey  []byte
	role        uint
}

func prepareConsensusChannel(role uint, left, right actor) *consensus_channel.ConsensusChannel {
	fp := state.FixedPart{
		ChainId:           big.NewInt(9001),
		Participants:      []types.Address{left.address, right.address},
		ChannelNonce:      big.NewInt(0),
		AppDefinition:     types.Address{},
		ChallengeDuration: big.NewInt(45),
	}

	leftBal := consensus_channel.NewBalance(left.destination, big.NewInt(6))
	rightBal := consensus_channel.NewBalance(right.destination, big.NewInt(4))

	lo := *consensus_channel.NewLedgerOutcome(types.Address{}, leftBal, rightBal, []consensus_channel.Guarantee{})

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
