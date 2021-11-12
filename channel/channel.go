package channel

import (
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/types"
)

type Channel struct {
	Id             types.Bytes32
	OnChainFunding types.Funds

	state.FixedPart
	Support []state.VariablePart
}

func (c Channel) GuaranteesFor(channelId types.Destination) types.Funds {
	return types.Funds{} // TODO get this info from the Support
}

func (c Channel) Total() types.Funds {
	return types.Funds{} // TODO get this info from the Support
}
