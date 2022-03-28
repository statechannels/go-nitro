package consensus_channel

import "fmt"

type FollowerChannel struct {
	consensusChannel
}

func (c *FollowerChannel) SignNextProposal(expectedUpdate interface{}, pk []byte) error {
	if len(c.proposalQueue) == 0 {
		return fmt.Errorf("no proposals in the queue")
	}
	p, ok := c.proposalQueue[0].Proposal.(Add)
	if !ok {
		return fmt.Errorf("only Add proposal is supported for queued proposals")
	}
	expectedP, ok := expectedUpdate.(Add)
	if !ok {
		return fmt.Errorf("only Add proposal is supported for expected update")
	}

	if !p.Equal(expectedP) {
		return fmt.Errorf("expected proposal does not match first proposal in the queue")
	}

	vars := c.current.clone()

	err := vars.Add(p)
	if err != nil {
		return err
	}
	c.proposalQueue = c.proposalQueue[1:]

	return nil
}
