package virtualfund

import (
	"math/big"

	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/internal/testactors"
	"github.com/statechannels/go-nitro/types"
)

// prepareConsensusChannel prepares a consensus channel with a consensus outcome
//  - allocating 6 to left
//  - allocating 4 to right
//  - including the given guarantees
func prepareConsensusChannel(role uint, left, right testactors.Actor, guarantees ...consensus_channel.Guarantee) *consensus_channel.ConsensusChannel {
	return prepareConsensusChannelHelper(role, left, right, 6, 4, 1, guarantees...)
}

// consensusStateSignatures prepares a consensus channel with a consensus outcome and returns the signatures on the consensus state
func consensusStateSignatures(left, right testactors.Actor, guarantees ...consensus_channel.Guarantee) [2]state.Signature {
	return prepareConsensusChannelHelper(0, left, right, 0, 0, 2, guarantees...).Signatures()
}

func prepareConsensusChannelHelper(role uint, left, right testactors.Actor, leftBalance, rightBalance, turnNum int, guarantees ...consensus_channel.Guarantee) *consensus_channel.ConsensusChannel {
	fp := state.FixedPart{
		ChainId:           big.NewInt(9001),
		Participants:      []types.Address{left.Address, right.Address},
		ChannelNonce:      big.NewInt(0),
		AppDefinition:     types.Address{},
		ChallengeDuration: big.NewInt(45),
	}

	leftBal := consensus_channel.NewBalance(left.Destination(), big.NewInt(int64(leftBalance)))
	rightBal := consensus_channel.NewBalance(right.Destination(), big.NewInt(int64(rightBalance)))

	lo := *consensus_channel.NewLedgerOutcome(types.Address{}, leftBal, rightBal, guarantees)

	signedVars := consensus_channel.SignedVars{Vars: consensus_channel.Vars{Outcome: lo, TurnNum: uint64(turnNum)}}
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
		cc, err = consensus_channel.NewLeaderChannel(fp, uint64(turnNum), lo, sigs)
	} else {
		cc, err = consensus_channel.NewFollowerChannel(fp, uint64(turnNum), lo, sigs)
	}
	if err != nil {
		panic(err)
	}

	return &cc
}
