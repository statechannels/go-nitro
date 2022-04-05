package consensus_channel

import (
	"fmt"

	"github.com/statechannels/go-nitro/channel/state"
)

var ErrNotLeader = fmt.Errorf("method may only be called by the channel leader")

// NewLeaderChannel constructs a new LeaderChannel
func NewLeaderChannel(fp state.FixedPart, turnNum uint64, outcome LedgerOutcome, signatures [2]state.Signature) (ConsensusChannel, error) {
	return newConsensusChannel(fp, Leader, turnNum, outcome, signatures)
}

// IsProposed returns whether or not the consensus state or any proposed state
// includes the given guarantee.
func (c *ConsensusChannel) IsProposed(g Guarantee) (bool, error) {
	if c.myIndex != Leader {
		return false, ErrNotLeader
	}
	latest, err := c.latestProposedVars()
	if err != nil {
		return false, err
	}

	return latest.Outcome.includes(g) && !c.Includes(g), nil
}

// Propose is called by the Leader and receives a proposal to add a guarantee,
// and generates and stores a SignedProposal in the queue, returning the
// resulting SignedProposal
//
// Note: the TurnNum on add is ignored; the correct turn number is computed
// and applied by c
func (c *ConsensusChannel) Propose(add Add, sk []byte) (SignedProposal, error) {
	if c.myIndex != Leader {
		return SignedProposal{}, ErrNotLeader
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

	signed := SignedProposal{Proposal: Proposal{toAdd: add}, Signature: signature}

	c.proposalQueue = append(c.proposalQueue, signed)
	return signed, nil
}

// UpdateConsensus is called by the Leader and iterates through
// the proposal queue until it finds the countersigned proposal.
//
// If this proposal was signed by the Follower:
//  - the consensus state is updated with the supplied proposal
//  - the proposal queue is trimmed
//
// If the countersupplied is stale (ie. proposal.TurnNum <= c.current.TurnNum) then
// their proposal is ignored.
//
// An error is returned if:
//  - the countersupplied proposal is not found
//  - or if it is found but not correctly by the Follower
func (c *ConsensusChannel) UpdateConsensus(countersigned SignedProposal) error {
	if c.myIndex != Leader {
		return ErrNotLeader
	}

	consensusCandidate := Vars{
		TurnNum: c.current.TurnNum,
		Outcome: c.current.Outcome.clone(),
	}
	if !countersigned.Proposal.isAddProposal() {
		// TODO: We'll need to expect other proposals in the future!
		return fmt.Errorf("unexpected proposal")
	}
	received := countersigned.Proposal.toAdd

	consensusTurnNum := received.turnNum

	if consensusTurnNum <= consensusCandidate.TurnNum {
		// We've already seen this proposal; return early
		return nil
	}

	for i, ourP := range c.proposalQueue {
		if !ourP.Proposal.isAddProposal() {
			// TODO: We'll need to expect other proposals in the future!
			return fmt.Errorf("unexpected proposal")
		}
		existing := ourP.Proposal.toAdd

		err := consensusCandidate.Add(existing)
		if err != nil {
			return err
		}

		if consensusCandidate.TurnNum == consensusTurnNum {
			signer, err := consensusCandidate.AsState(c.fp).RecoverSigner(countersigned.Signature)

			if err != nil {
				return fmt.Errorf("unable to recover signer: %w", err)
			}

			if signer != c.fp.Participants[Follower] {
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
