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

	IsTwoPartyLedger bool
	MyDestination    types.Destination
	TheirDestination types.Destination // must be nonzero if a two party ledger channel
}

// New constructs a new channel from the supplied state
func New(s state.State, isTwoPartyLedger bool, myDestination types.Destination, theirDestination types.Destination) Channel {
	c := Channel{}

	c.OnChainFunding = make(types.Funds)

	c.LatestSupportedState = s.Clone()
	c.FixedPart = c.LatestSupportedState.FixedPart()

	c.Support = make([]state.VariablePart, 0)
	c.MyDestination = myDestination
	c.TheirDestination = theirDestination
	c.IsTwoPartyLedger = isTwoPartyLedger
	return c
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
