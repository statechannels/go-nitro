package virtualfund

import (
	"math/big"

	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/internal/testactors"
	"github.com/statechannels/go-nitro/types"
)

// prepareConsensusChannel prepares a consensus channel with a consensus outcome
//   - allocating 5 to left participant
//   - allocating 5 to right participant
//   - including the given guarantees
func prepareConsensusChannel(role uint, leader, follower, leftActor testactors.Actor, guarantees ...consensus_channel.Guarantee) *consensus_channel.ConsensusChannel {
	return prepareConsensusChannelHelper(role, leader, follower, leftActor, 5, 5, 1, guarantees...)
}

// consensusStateSignatures prepares a consensus channel with a consensus outcome and returns the signatures on the consensus state
func consensusStateSignatures(leader, follower testactors.Actor, guarantees ...consensus_channel.Guarantee) [2]state.Signature {
	return prepareConsensusChannelHelper(0, leader, follower, leader, 0, 5, 2, guarantees...).Signatures()
}

func prepareConsensusChannelHelper(role uint, leader, follower, leftActor testactors.Actor, leftBalance, rightBalance, turnNum int, guarantees ...consensus_channel.Guarantee) *consensus_channel.ConsensusChannel {
	fp := state.FixedPart{
		Participants:      []types.Address{leader.Address(), follower.Address()},
		ChannelNonce:      0,
		AppDefinition:     types.Address{},
		ChallengeDuration: 45,
	}
	var leaderBal, followerBal consensus_channel.Balance

	if leader.Address() == leftActor.Address() {
		leaderBal = consensus_channel.NewBalance(leader.Destination(), big.NewInt(int64(leftBalance)))
		followerBal = consensus_channel.NewBalance(follower.Destination(), big.NewInt(int64(rightBalance)))
	} else {
		leaderBal = consensus_channel.NewBalance(leader.Destination(), big.NewInt(int64(rightBalance)))
		followerBal = consensus_channel.NewBalance(follower.Destination(), big.NewInt(int64(leftBalance)))
	}

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
