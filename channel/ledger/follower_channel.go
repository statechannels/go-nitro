package ledger

import (
	"fmt"

	"github.com/statechannels/go-nitro/channel/state"
)

var ErrNotFollower = fmt.Errorf("method may only be called by channel follower")
var ErrNoProposals = fmt.Errorf("no proposals in the queue")
var ErrUnsupportedQueuedProposal = fmt.Errorf("only Add proposal is supported for queued proposals")
var ErrUnsupportedExpectedProposal = fmt.Errorf("only Add proposal is supported for expected update")
var ErrNonMatchingProposals = fmt.Errorf("expected proposal does not match first proposal in the queue")
var ErrInvalidProposalSignature = fmt.Errorf("invalid signature for proposal")
var ErrInvalidTurnNum = fmt.Errorf("the proposal turn number is not the next turn number")

// NewFollowerChannel constructs a new FollowerChannel
func NewFollowerChannel(fp state.FixedPart, turnNum uint64, outcome LedgerOutcome, signatures [2]state.Signature) (LedgerChannel, error) {
	return newLedgerChannel(fp, Follower, turnNum, outcome, signatures)
}

// SignNextProposal is called by the follower and inspects whether the
// expected proposal matches the first proposal in the queue. If so,
// the proposal is removed from the queue and integrated into the channel state.
func (c *LedgerChannel) SignNextProposal(expectedProposal Proposal, sk []byte) (SignedProposal, error) {
	if c.MyIndex != Follower {
		return SignedProposal{}, ErrNotFollower
	}

	if err := c.validateProposalID(expectedProposal); err != nil {
		return SignedProposal{}, err
	}

	if len(c.proposalQueue) == 0 {
		return SignedProposal{}, ErrNoProposals
	}

	p := c.proposalQueue[0].Proposal

	if !p.Equal(&expectedProposal) {
		return SignedProposal{}, ErrNonMatchingProposals
	}

	// vars are cloned and modified instead of modified in place to simplify recovering from error
	vars := Vars{
		TurnNum: c.current.TurnNum,
		Outcome: c.current.Outcome.clone(),
	}
	err := vars.HandleProposal(p)
	if err != nil {
		return SignedProposal{}, err
	}

	signature, err := c.sign(vars, sk)
	if err != nil {
		return SignedProposal{}, fmt.Errorf("unable to sign state update: %f", err)
	}

	signed := c.proposalQueue[0]
	c.current = SignedVars{
		Vars:       vars,
		Signatures: [2]state.Signature{signed.Signature, signature},
	}
	c.proposalQueue = c.proposalQueue[1:]

	return SignedProposal{signature, signed.Proposal}, nil
}

// followerReceive is called by the follower to validate a proposal from the leader and add it to the proposal queue
func (c *LedgerChannel) followerReceive(p SignedProposal) error {
	if c.MyIndex != Follower {
		return ErrNotFollower
	}

	if err := c.validateProposalID(p.Proposal); err != nil {
		return err
	}

	// Get the latest proposal vars we have
	vars, err := c.latestProposedVars()
	if err != nil {
		return fmt.Errorf("could not generate the current proposal: %w", err)
	}

	if p.Proposal.TurnNum() != vars.TurnNum+1 {
		return ErrInvalidTurnNum
	}
	// Add the incoming proposal to the vars
	err = vars.HandleProposal(p.Proposal)
	if err != nil {
		return fmt.Errorf("receive could not add new state vars: %w", err)
	}

	// Validate the signature
	signer, err := c.recoverSigner(vars, p.Signature)
	if err != nil {
		return fmt.Errorf("receive could not recover signature: %w", err)
	}
	if signer != c.Leader() {
		return ErrInvalidProposalSignature
	}

	// Update the proposal queue
	c.proposalQueue = append(c.proposalQueue, p)

	return nil
}
