package consensus

import (
	"github.com/statechannels/go-nitro/channel"
)

// NOTE: example of direct channel app with no special logic

type ConsensusApp struct {
	//
}

func (a *ConsensusApp) Type() string {
	return "consensus"
}

func (a *ConsensusApp) HandleRequest(ch *channel.Channel, ty string, data interface{}) error {
	//

	return nil
}
