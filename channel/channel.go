package channel

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/types"
)

type SetupState struct {
	State    state.State
	Signed   bool
	Complete bool // A signature from each and every peer
}

type Channel struct {
	Id types.Destination

	PreFund  SetupState
	PostFund SetupState

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

	// if s.TurnNum != 0 return error // TODO

	c.PreFund = SetupState{
		s.Clone(),
		false, false,
	}
	c.PostFund = SetupState{
		s.Clone(),
		false, false,
	}
	c.PostFund.State.TurnNum = big.NewInt(1)

	c.Id, _ = c.PreFund.State.ChannelId() // TODO handle error

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

// AddSignedState adds a signed state to the Channel, updating the LatestSupportedState and Support if appropriate.
// Returns false and does not alter the channel if the state is "stale"
func (c Channel) AddSignedState(s state.State, sig state.Signature) bool {
	// TODO
	// If the turnNum is below that of the supported state, discard / error / return false
	// If it is greater than, keep it around in case it becomes supported in future
	// If it is equal to ... ? probably discard / error / return false

	// Check and update the latest supported state and proof
	return true
}

// AddSignedStates adds each signed state in the mapping
func (c Channel) AddSignedStates(mapping map[*state.State]state.Signature) {
	for state, sig := range mapping {
		c.AddSignedState(*state, sig)
	}
}
