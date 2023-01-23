package client

import (
	"math/big"

	"github.com/statechannels/go-nitro/channel"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/client/engine/store"
	"github.com/statechannels/go-nitro/types"
)

type ChannelStatus string

// TODO: Think through statuses
const Proposed ChannelStatus = "Proposed"
const Ready ChannelStatus = "Ready"
const Closing ChannelStatus = "Closing"
const Complete ChannelStatus = "Complete"

// PaymentChannelBalance contains the balance of a uni-directional payment channel
type PaymentChannelBalance struct {
	AssetAddress   types.Address
	Payee          types.Address
	Payer          types.Address
	PaidSoFar      *big.Int
	RemainingFunds *big.Int
}

// PaymentChannelInfo contains balance and status info about a payment channel
type PaymentChannelInfo struct {
	ID      types.Destination
	Status  ChannelStatus
	Balance PaymentChannelBalance
}

// getStatusFromChannel returns the status of the channel
func getStatusFromChannel(c *channel.Channel) ChannelStatus {
	if c.FinalSignedByMe() {
		if c.FinalCompleted() {
			return Complete
		}
		return Closing
	}
	if !c.PostFundComplete() {
		return Proposed
	}

	return Ready
}

func getPaymentChannelBalance(c *channel.Channel) PaymentChannelBalance {

	latest := c.PreFundState()
	if c.HasSupportedState() {

		supported, _ := c.LatestSupportedState()
		latest = supported
	}

	numParticipants := len(latest.Participants)
	// TODO: We assume single asset outcomes
	outcome := latest.Outcome[0]
	asset := outcome.Asset
	payer := latest.Participants[0]
	payee := latest.Participants[numParticipants-1]
	paidSoFar := outcome.Allocations[1].Amount
	remaining := outcome.Allocations[0].Amount
	return PaymentChannelBalance{
		AssetAddress:   asset,
		Payer:          payer,
		Payee:          payee,
		PaidSoFar:      paidSoFar,
		RemainingFunds: remaining,
	}
}

// getPaymentChannelInfo returns the PaymentChannelInfo for the given channel
func getPaymentChannelInfo(channel *channel.Channel) PaymentChannelInfo {
	return PaymentChannelInfo{
		ID:      channel.Id,
		Status:  getStatusFromChannel(channel),
		Balance: getPaymentChannelBalance(channel),
	}
}

// LedgerChannelInfo contains balance and status info about a ledger channel
type LedgerChannelInfo struct {
	ID      types.Destination
	Status  ChannelStatus
	Balance LedgerChannelBalance
}

// LedgerChannelBalance contains the balance of a ledger channel
type LedgerChannelBalance struct {
	AssetAddress  types.Address
	Hub           types.Address
	Client        types.Address
	HubBalance    *big.Int
	ClientBalance *big.Int
}

// getLatestSupported returns the latest supported state of the channel or the prefund state if no supported state exists
func getLatestSupported(channel *channel.Channel) state.State {
	if channel.HasSupportedState() {
		supported, _ := channel.LatestSupportedState()
		return supported
	}
	return channel.PreFundState()
}

// getLedgerBalanceFromState returns the balance of the ledger channel from the given state
func getLedgerBalanceFromState(latest state.State) LedgerChannelBalance {

	// TODO: We assume single asset outcomes
	outcome := latest.Outcome[0]
	asset := outcome.Asset
	client := latest.Participants[0]
	clientBalance := outcome.Allocations[0].Amount
	hub := latest.Participants[1]
	hubBalance := outcome.Allocations[1].Amount

	return LedgerChannelBalance{
		AssetAddress:  asset,
		Hub:           hub,
		Client:        client,
		HubBalance:    hubBalance,
		ClientBalance: clientBalance,
	}
}

// getLedgerChannelInfo returns the LedgerChannelInfo for the given channel
// It does this by querying the provided store
func getLedgerChannelInfo(id types.Destination, store store.Store) (LedgerChannelInfo, error) {
	c, ok := store.GetChannelById(id)
	if ok {

		return LedgerChannelInfo{
			ID:      c.Id,
			Status:  getStatusFromChannel(c),
			Balance: getLedgerBalanceFromState(getLatestSupported(c)),
		}, nil
	}

	con, err := store.GetConsensusChannelById(id)
	if err != nil {
		return LedgerChannelInfo{}, err
	}

	latest := con.ConsensusVars().AsState(con.FixedPart())
	return LedgerChannelInfo{
		ID:      con.Id,
		Status:  Ready,
		Balance: getLedgerBalanceFromState(latest),
	}, nil

}
