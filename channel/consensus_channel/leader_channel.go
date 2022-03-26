package consensus_channel

import (
	"fmt"

	"github.com/statechannels/go-nitro/channel/state"
)

// LeaderChannel is used by a leader's virtualfund objective to make and receive ledger updates
type LeaderChannel struct {
	consensusChannel
}

// NewLeaderChannel constructs a new LeaderChannel
func NewLeaderChannel(fp state.FixedPart, outcome LedgerOutcome, signatures [2]state.Signature) (LeaderChannel, error) {
	channel, err := newConsensusChannel(fp, leader, outcome, signatures)

	return LeaderChannel{consensusChannel: channel}, err
}

// ConsensusTurnNum returns the turn number of the current consensus state
func (c *LeaderChannel) ConsensusTurnNum() uint64 {
	return c.current.TurnNum
}

// IsProposed returns whether or not the consensus state or any proposed state
// includes the given guarantee
func (c *LeaderChannel) IsProposed(g Guarantee) (bool, error) {
	latest, err := c.latestProposedVars()
	if err != nil {
		return false, err
	}

	return latest.Outcome.includes(g), nil
}

// Includes returns whether or not the consensus state includes the given guarantee
func (c *LeaderChannel) Includes(g Guarantee) bool {
	return c.current.Outcome.includes(g)
}

// Propose receives a proposal to add a guarantee, and generates and stores a SignedProposal in
// the queue, returning the resulting SignedProposal
// Note: the TurnNum on add is ignored; the correct turn number is computed by c
func (c *LeaderChannel) Propose(add Add, sk []byte) (SignedProposal, error) {
	if c.myIndex != leader {
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
