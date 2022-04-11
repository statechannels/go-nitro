package virtualfund

import (
	"math/big"

	con_chan "github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/internal/testactors"
	"github.com/statechannels/go-nitro/types"
)

// prepareConsensusChannel prepares a consensus channel with a consensus outcome
//  - allocating 6 to left
//  - allocating 4 to right
//  - including the given guarantees
func prepareConsensusChannel(role uint, left, right testactors.Actor, guarantees ...con_chan.Guarantee) *con_chan.ConsensusChannel {
	fp := state.FixedPart{
		ChainId:           big.NewInt(9001),
		Participants:      []types.Address{left.Address, right.Address},
		ChannelNonce:      big.NewInt(0),
		AppDefinition:     types.Address{},
		ChallengeDuration: big.NewInt(45),
	}

	leftBal := con_chan.NewBalance(left.Destination(), big.NewInt(6))
	rightBal := con_chan.NewBalance(right.Destination(), big.NewInt(4))

	lo := *con_chan.NewLedgerOutcome(types.Address{}, leftBal, rightBal, guarantees)

	signedVars := con_chan.SignedVars{Vars: con_chan.Vars{Outcome: lo, TurnNum: 1}}
	leftSig, err := signedVars.Vars.AsState(fp).Sign(left.PrivateKey)
	if err != nil {
		panic(err)
	}
	rightSig, err := signedVars.Vars.AsState(fp).Sign(right.PrivateKey)
	if err != nil {
		panic(err)
	}
	sigs := [2]state.Signature{leftSig, rightSig}

	var cc con_chan.ConsensusChannel

	if role == 0 {
		cc, err = con_chan.NewLeaderChannel(fp, 1, lo, sigs)
	} else {
		cc, err = con_chan.NewFollowerChannel(fp, 1, lo, sigs)
	}
	if err != nil {
		panic(err)
	}

	return &cc
}
