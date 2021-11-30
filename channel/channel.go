package channel

import (
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/types"
)

type Channel struct {
	Id             types.Destination
	OnChainFunding types.Funds

	state.FixedPart
	Support              []state.VariablePart
	LatestSupportedState state.State
}

func (c Channel) GuaranteesFor(channelId types.Destination) types.Funds {
	return types.Funds{} // TODO get this info from the Support
}

func (c Channel) Total() types.Funds {
	funds := types.Funds{}
	for _, sae := range c.LatestSupportedState.Outcome {
		funds[sae.Asset] = sae.Allocations.Total()
	}
	return funds
}
