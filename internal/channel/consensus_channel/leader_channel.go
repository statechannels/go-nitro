package consensus_channel

import (
	"fmt"

	"github.com/statechannels/go-nitro/channel/state"
)

var (
	ErrNotLeader              = fmt.Errorf("method may only be called by the channel leader")
	ErrProposalQueueExhausted = fmt.Errorf("proposal queue exhausted")
	ErrWrongSigner            = fmt.Errorf("proposal incorrectly signed")
)

// NewLeaderChannel constructs a new LeaderChannel
func NewLeaderChannel(fp state.FixedPart, turnNum uint64, outcome LedgerOutcome, signatures [2]state.Signature) (ConsensusChannel, error) {
	return newConsensusChannel(fp, Leader, turnNum, outcome, signatures)
}

// IsProposed returns true if a proposal in the queue would lead to g being included in the receiver's outcome, and false otherwise.
//
// Specific clarification: If the current outcome already includes g, IsProposed returns false.
func (c *ConsensusChannel) IsProposed(g Guarantee) (bool, error) {
	latest, err := c.latestProposedVars()
	if err != nil {
		return false, err
	}

	return latest.Outcome.includes(g) && !c.Includes(g), nil

}

// IsProposedNext returns if the next proposal in the queue would lead to g being included in the receiver's outcome, and false otherwise.
func (c *ConsensusChannel) IsProposedNext(g Guarantee) (bool, error) {
	vars := Vars{TurnNum: c.current.TurnNum, Outcome: c.current.Outcome.clone()}

	if len(c.proposalQueue) == 0 {
		return false, nil
	}

	p := c.proposalQueue[0]
	err := vars.HandleProposal(p.Proposal)
	if vars.TurnNum != p.TurnNum {
		return false, fmt.Errorf("proposal turn number %d does not match vars %d", p.TurnNum, vars.TurnNum)
	}

	if err != nil {
		return false, err
	}

	return vars.Outcome.includes(g) && !c.Includes(g), nil
}

// Propose is called by the Leader and receives a proposal to add a guarantee,
// and generates and stores a SignedProposal in the queue, returning the
// resulting SignedProposal
//
func (c *ConsensusChannel) Propose(proposal Proposal, sk []byte) (SignedProposal, error) {
	if c.MyIndex != Leader {
		return SignedProposal{}, ErrNotLeader
	}

	if proposal.ChannelID != c.Id {
		return SignedProposal{}, ErrIncorrectChannelID
	}
	vars, err := c.latestProposedVars()
	if err != nil {
		return SignedProposal{}, fmt.Errorf("unable to construct latest proposed vars: %w", err)
	}

	err = vars.HandleProposal(proposal)
	if err != nil {
		return SignedProposal{}, fmt.Errorf("propose could not add new state vars: %w", err)
	}

	signature, err := c.sign(vars, sk)
	if err != nil {
		return SignedProposal{}, fmt.Errorf("unable to sign state update: %f", err)
	}

	signed := SignedProposal{Proposal: proposal, Signature: signature, TurnNum: vars.TurnNum}

	c.proposalQueue = append(c.proposalQueue, signed)
	return signed, nil
}

// leaderReceive is called by the Leader and iterates through
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
//  - or if it is found but not correctly signed by the Follower
func (c *ConsensusChannel) leaderReceive(countersigned SignedProposal) error {
	if c.MyIndex != Leader {
		return ErrNotLeader
	}

	if err := c.validateProposalID(countersigned.Proposal); err != nil {
		return err
	}

	consensusCandidate := Vars{
		TurnNum: c.current.TurnNum,
		Outcome: c.current.Outcome.clone(),
	}
	consensusTurnNum := countersigned.TurnNum

	if consensusTurnNum <= consensusCandidate.TurnNum {
		// We've already seen this proposal; return early
		return nil
	}

	for i, ourP := range c.proposalQueue {

		err := consensusCandidate.HandleProposal(ourP.Proposal)
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
