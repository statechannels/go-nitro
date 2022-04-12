package virtualfund

import (
	"math/big"

	"github.com/statechannels/go-nitro/channel/consensus_channel"
	con_chan "github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/internal/testactors"
	"github.com/statechannels/go-nitro/types"
)

type CChanConfig struct {
	left       testactors.Actor
	right      testactors.Actor
	leftBal    int64
	rightBal   int64
	leader     bool
	guarantees []con_chan.Guarantee
	props      []con_chan.Proposal
}

// prepareConsensusChannel prepares a consensus channel with a consensus outcome
//  - allocating a default amount of 6 to cfg.left
//  - allocating a default amount of 4 to cfg.right
//  - including the given guarantees
//  - ensuring that the props are signed and stored by the consensus channel
// Note: The props that are passed in cfg don't need to be set up correctly!
// The correct turn number and channelId will be set in order to ensure that
// the resulting queue is in a valid state
func prepareConsensusChannel(cfg CChanConfig) *con_chan.ConsensusChannel {
	leftBal := cfg.leftBal
	if leftBal == 0 {
		leftBal = 6
	}
	rightBal := cfg.rightBal
	if rightBal == 0 {
		rightBal = 4
	}

	initialOutcome := func() con_chan.LedgerOutcome {
		left := con_chan.NewBalance(cfg.left.Destination(), big.NewInt(leftBal))
		right := con_chan.NewBalance(cfg.right.Destination(), big.NewInt(rightBal))

		return *con_chan.NewLedgerOutcome(types.Address{}, left, right, cfg.guarantees)

	}

	participants := [2]types.Address{
		cfg.left.Address, cfg.right.Address,
	}
	fp := state.FixedPart{
		Participants:      participants[:],
		ChainId:           big.NewInt(0),
		ChannelNonce:      big.NewInt(9001),
		ChallengeDuration: big.NewInt(100),
	}

	startingTurnNum := uint64(1)

	initialVars := con_chan.Vars{TurnNum: uint64(startingTurnNum), Outcome: initialOutcome()}
	asState := initialVars.AsState(fp)
	leftSig, _ := asState.Sign(cfg.left.PrivateKey)
	rightSig, _ := asState.Sign(cfg.right.PrivateKey)
	sigs := [2]state.Signature{leftSig, rightSig}

	var c con_chan.ConsensusChannel
	var err error
	if cfg.leader {
		c, err = con_chan.NewLeaderChannel(fp, 1, initialOutcome(), sigs)
	} else {
		c, err = con_chan.NewFollowerChannel(fp, 1, initialOutcome(), sigs)
	}
	if err != nil {
		panic(err)
	}

	// make sure the proposals are all valid for this channel
	turnNum := startingTurnNum
	vars := con_chan.Vars{TurnNum: startingTurnNum, Outcome: initialOutcome()}
	for _, p := range cfg.props {
		turnNum += 1
		var correctedProp con_chan.Proposal
		switch p.Type() {
		case consensus_channel.AddProposal:
			correctedProp = con_chan.NewAddProposal(c.Id, turnNum, p.ToAdd.Guarantee, p.ToAdd.LeftDeposit)
		case consensus_channel.RemoveProposal:
			correctedProp = con_chan.NewRemoveProposal(c.Id, turnNum, p.ToRemove.Target, p.ToRemove.LeftAmount, p.ToRemove.RightAmount)
		}

		if cfg.leader {
			// Call Propose with the corrected proposal
			_, err := c.Propose(correctedProp, cfg.left.PrivateKey)
			if err != nil {
				panic(err)
			}
		} else {
			// Compute the correct signature on the corrected prop
			err = vars.HandleProposal(correctedProp)
			if err != nil {
				panic(err)
			}
			s := vars.AsState(fp)
			sig, err := s.Sign(cfg.left.PrivateKey)
			if err != nil {
				panic(err)
			}
			sp := con_chan.SignedProposal{Proposal: correctedProp, Signature: sig}

			// Receive the signed prop
			err = c.Receive(sp)
			if err != nil {
				panic(err)
			}
		}
	}

	return &c
}
