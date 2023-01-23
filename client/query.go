package client

import (
	"math/big"

	"github.com/statechannels/go-nitro/channel"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/client/engine/store"
	"github.com/statechannels/go-nitro/types"
)

type ChannelStatus string

const Proposed ChannelStatus = "Proposed"
const Ready ChannelStatus = "Ready"
const Closing ChannelStatus = "Closing"
const Complete ChannelStatus = "Complete"

type PaymentChannelBalance struct {
	AssetAddress   types.Address
	Payee          types.Address
	Payer          types.Address
	PaidSoFar      *big.Int
	RemainingFunds *big.Int
}
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

func getPaymentChannelInfo(channel *channel.Channel) PaymentChannelInfo {
	return PaymentChannelInfo{
		ID:      channel.Id,
		Status:  getStatusFromChannel(channel),
		Balance: getPaymentChannelBalance(channel),
	}
}

type LedgerChannelInfo struct {
	ID      types.Destination
	Status  ChannelStatus
	Balance LedgerChannelBalance
}
type LedgerChannelBalance struct {
	AssetAddress  types.Address
	Hub           types.Address
	Client        types.Address
	HubBalance    *big.Int
	ClientBalance *big.Int
}

func getLatest(channel *channel.Channel) state.State {
	if channel.HasSupportedState() {
		supported, _ := channel.LatestSupportedState()
		return supported
	}
	return channel.PreFundState()
}
func getBalanceFromState(latest state.State) LedgerChannelBalance {

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

func GetLedgerChannelInfo(id types.Destination, store store.Store) (LedgerChannelInfo, error) {
	c, ok := store.GetChannelById(id)
	if ok {

		return LedgerChannelInfo{
			ID:      c.Id,
			Status:  getStatusFromChannel(c),
			Balance: getBalanceFromState(getLatest(c)),
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
		Balance: getBalanceFromState(latest),
	}, nil

}

type NitroQuery interface {
	// GetPaymentChannelByReceiver returns the first payment found channel with the given receiver.
	// If no payment channel exists nil is returned.
	GetPaymentChannelByReceiver(receiver types.Address) PaymentChannelInfo
	// GetLedgerChannelByHub returns the first ledger channel found with the given hub.
	// If no ledger channel exists nil is returned.
	GetLedgerChannelByHub(Hub types.Address) LedgerChannelInfo
	// GetLedgerChannel returnns the ledger channel for the given id.
	// If no ledger channel exists nil is returned.
	GetLedgerChannel(id types.Destination) LedgerChannelInfo
	// GetPaymentChannel returnns the ledger channel for the given id.
	// If no ledger channel exists nil is returned.
	GetPaymentChannel(id types.Destination) PaymentChannelInfo
}
