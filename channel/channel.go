package channel

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
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

// Affords returns true if, for each asset keying the input variables, the channel can afford the allocation given the funding.
// The decision is made based on the latest supported state of the channel.
//
// Both arguments are maps keyed by the same asset
func (c Channel) Affords(
	allocationMap map[common.Address]outcome.Allocation,
	fundingMap types.Funds) bool {
	return c.LatestSupportedState.Outcome.Affords(allocationMap, fundingMap)
}
