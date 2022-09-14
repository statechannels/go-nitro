package virtualfund

import (
	"math/big"

	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/internal/testactors"
	"github.com/statechannels/go-nitro/types"
)

// prepareConsensusChannel prepares a consensus channel with a consensus outcome
//   - allocating 6 to leader
//   - allocating 4 to follower
//   - including the given guarantees
func prepareConsensusChannel(role uint, leader, follower testactors.Actor, guarantees ...consensus_channel.Guarantee) *consensus_channel.ConsensusChannel {
	return prepareConsensusChannelHelper(role, leader, follower, 6, 4, 1, guarantees...)
}

// consensusStateSignatures prepares a consensus channel with a consensus outcome and returns the signatures on the consensus state
func consensusStateSignatures(leader, follower testactors.Actor, guarantees ...consensus_channel.Guarantee) [2]state.Signature {
	return prepareConsensusChannelHelper(0, leader, follower, 0, 0, 2, guarantees...).Signatures()
}

func prepareConsensusChannelHelper(role uint, leader, follower testactors.Actor, leftBalance, rightBalance, turnNum int, guarantees ...consensus_channel.Guarantee) *consensus_channel.ConsensusChannel {
	fp := state.FixedPart{
		ChainId:           big.NewInt(9001),
		Participants:      []types.Address{leader.Address(), follower.Address()},
		ChannelNonce:      0,
		AppDefinition:     types.Address{},
		ChallengeDuration: big.NewInt(45),
	}

	leaderBal := consensus_channel.NewBalance(leader.Destination(), big.NewInt(int64(leftBalance)))
	followerBal := consensus_channel.NewBalance(follower.Destination(), big.NewInt(int64(rightBalance)))

	lo := *consensus_channel.NewLedgerOutcome(types.Address{}, leaderBal, followerBal, guarantees)

	signedVars := consensus_channel.SignedVars{Vars: consensus_channel.Vars{Outcome: lo, TurnNum: uint64(turnNum)}}
	leaderSig, err := signedVars.Vars.AsState(fp).Sign(leader.PrivateKey)
	if err != nil {
		panic(err)
	}
	followerSig, err := signedVars.Vars.AsState(fp).Sign(follower.PrivateKey)
	if err != nil {
		panic(err)
	}
	sigs := [2]state.Signature{leaderSig, followerSig}

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
