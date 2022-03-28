package consensus_channel

import (
	"fmt"

	"github.com/statechannels/go-nitro/channel/state"
)

var ErrNoProposals = fmt.Errorf("no proposals in the queue")
var ErrUnsupportedQueuedProposal = fmt.Errorf("only Add proposal is supported for queued proposals")
var ErrUnsupportedExpectedProposal = fmt.Errorf("only Add proposal is supported for expected update")
var ErrNonMatchingProposals = fmt.Errorf("expected proposal does not match first proposal in the queue")

type FollowerChannel struct {
	consensusChannel
}

func NewFollowerChannel(fp state.FixedPart, outcome LedgerOutcome, signatures [2]state.Signature) (FollowerChannel, error) {
	channel, err := newConsensusChannel(fp, leader, outcome, signatures)

	return FollowerChannel{consensusChannel: channel}, err
}

// ConsensusTurnNum returns the turn number of the current consensus state
func (c *FollowerChannel) ConsensusTurnNum() uint64 {
	return c.current.TurnNum
}

// Includes returns whether or not the consensus state includes the given guarantee
func (c *FollowerChannel) Includes(g Guarantee) bool {
	return c.current.Outcome.includes(g)
}

func (c *FollowerChannel) SignNextProposal(expectedUpdate interface{}, pk []byte) error {
	if len(c.proposalQueue) == 0 {
		return ErrNoProposals
	}
	p, ok := c.proposalQueue[0].Proposal.(Add)
	if !ok {
		return ErrUnsupportedQueuedProposal
	}
	expectedP, ok := expectedUpdate.(Add)
	if !ok {
		return ErrUnsupportedExpectedProposal
	}

	if !p.Equal(expectedP) {
		return ErrNonMatchingProposals
	}

	vars := Vars{
		TurnNum: c.current.TurnNum,
		Outcome: c.current.Outcome.clone(),
	}

	err := vars.Add(p)
	if err != nil {
		return err
	}
	c.proposalQueue = c.proposalQueue[1:]

	return nil
}
