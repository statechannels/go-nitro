package consensus_channel

import (
	"fmt"

	"github.com/statechannels/go-nitro/channel/state"
)

type LeaderChannel struct {
	LeaderInterface
	ConsensusChannel
}

func NewLeaderChannel(fp state.FixedPart, outcome LedgerOutcome, signatures [2]state.Signature) (LeaderChannel, error) {
	channel, err := NewConsensusChannel(fp, leader, outcome, signatures)

	return LeaderChannel{ConsensusChannel: channel}, err
}

type LeaderInterface interface {
	Propose(add Add, sk []byte) (SignedProposal, error)
}

// Propose receives a proposal to add a guarantee, and generates and stores a SignedProposal in
// the queue, returning the resulting SignedProposal
// Note: the TurnNum on add is ignored; the correct turn number is computed by c
func (c *LeaderChannel) Propose(add Add, sk []byte) (SignedProposal, error) {
	if c.MyIndex != leader {
		return SignedProposal{}, fmt.Errorf("only proposer can call Add")
	}

	vars, err := c.latestProposedVars()
	if err != nil {
		return SignedProposal{}, fmt.Errorf("unable to construct latest proposed vars: %w", err)
	}

	add.turnNum = vars.TurnNum + 1

	err = vars.Add(add)
	if err != nil {
		return SignedProposal{}, fmt.Errorf("propose could not add new state vars: %w", err)
	}

	signature, err := c.sign(vars, sk)
	if err != nil {
		return SignedProposal{}, fmt.Errorf("unable to sign state update: %f", err)
	}

	signed := SignedProposal{Proposal: add, Signature: signature}

	c.proposalQueue = append(c.proposalQueue, signed)
	return signed, nil
}
