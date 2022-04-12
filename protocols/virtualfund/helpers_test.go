package virtualfund

import (
	"math/big"

	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/internal/testactors"
	"github.com/statechannels/go-nitro/types"
)

func prepareConsensusChannel(role uint, left, right testactors.Actor) *consensus_channel.ConsensusChannel {
	fp := state.FixedPart{
		ChainId:           big.NewInt(9001),
		Participants:      []types.Address{left.Address, right.Address},
		ChannelNonce:      big.NewInt(0),
		AppDefinition:     types.Address{},
		ChallengeDuration: big.NewInt(45),
	}

	leftBal := consensus_channel.NewBalance(left.Destination(), big.NewInt(6))
	rightBal := consensus_channel.NewBalance(right.Destination(), big.NewInt(4))

	lo := *consensus_channel.NewLedgerOutcome(types.Address{}, leftBal, rightBal, []consensus_channel.Guarantee{})

	signedVars := consensus_channel.SignedVars{Vars: consensus_channel.Vars{Outcome: lo, TurnNum: 1}}
	leftSig, err := signedVars.Vars.AsState(fp).Sign(left.PrivateKey)
	if err != nil {
		panic(err)
	}
	rightSig, err := signedVars.Vars.AsState(fp).Sign(right.PrivateKey)
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
