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

// IsProposed returns whether or not the consensus state or any proposed state
// includes the given guarantee.
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

// UpdateConsensus iterates through the proposal queue until it finds the countersigned proposal.
// If this proposal was signed by the Follower:
// - the consensus state is updated with the supplied proposal
// - the proposal queue is trimmed
//
// If the countersupplied is stale (ie. proposal.TurnNum <= c.current.TurnNum) then
// their proposal is ignored.
//
// An error is returned if:
// - the countersupplied proposal is not found
// - or if it is found but not correctly by the Follower
func (c *LeaderChannel) UpdateConsensus(countersigned SignedProposal) error {
	consensusCandidate := Vars{
		TurnNum: c.current.TurnNum,
		Outcome: c.current.Outcome.clone(),
	}

	received, ok := countersigned.Proposal.(Add)
	if !ok {
		// TODO: We'll need to expect other proposals in the future!
		return fmt.Errorf("unexpected proposal")
	}
	consensusTurnNum := received.turnNum

	if consensusTurnNum <= consensusCandidate.TurnNum {
		// We've already seen this proposal; return early
		return nil
	}

	for i, ourP := range c.proposalQueue {
		existing, ok := ourP.Proposal.(Add)
		if !ok {
			// TODO: We'll need to expect other proposals in the future!
			return fmt.Errorf("unexpected proposal")
		}

		err := consensusCandidate.Add(existing)
		if err != nil {
			return err
		}

		if consensusCandidate.TurnNum == consensusTurnNum {
			signer, err := consensusCandidate.asState(c.fp).RecoverSigner(countersigned.Signature)

			if err != nil {
				return fmt.Errorf("unable to recover signer: %w", err)
			}

			if signer != c.fp.Participants[follower] {
				return ErrWrongSigner
			}

			mySig := ourP.Signature
			c.current = SignedVars{
				Vars:       consensusCandidate,
				Signatures: [2]state.Signature{mySig, countersigned.Signature},
			}

			c.proposalQueue = c.proposalQueue[i+1:]

			return nil
		}
	}

	return ErrProposalQueueExhausted
}

var ErrProposalQueueExhausted = fmt.Errorf("proposal queue exhausted")
var ErrWrongSigner = fmt.Errorf("proposal incorrectly signed")
