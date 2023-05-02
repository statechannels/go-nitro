package query

import (
	"fmt"
	"math/big"

	"github.com/statechannels/go-nitro/channel"
	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/client/engine/store"
	"github.com/statechannels/go-nitro/payments"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/protocols/virtualfund"
	"github.com/statechannels/go-nitro/types"
)

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

// getPaymentChannelBalance generates a PaymentChannelBalance from the given participants and outcome
func getPaymentChannelBalance(participants []types.Address, outcome outcome.Exit) PaymentChannelBalance {
	numParticipants := len(participants)
	// TODO: We assume single asset outcomes
	sao := outcome[0]
	asset := sao.Asset
	payer := participants[0]
	payee := participants[numParticipants-1]
	paidSoFar := big.NewInt(0).Set(sao.Allocations[1].Amount)
	remaining := big.NewInt(0).Set(sao.Allocations[0].Amount)
	return PaymentChannelBalance{
		AssetAddress:   asset,
		Payer:          payer,
		Payee:          payee,
		PaidSoFar:      paidSoFar,
		RemainingFunds: remaining,
	}
}

// getLatestSupported returns the latest supported state of the channel
// or the prefund state if no supported state exists
func getLatestSupported(channel *channel.Channel) (state.State, error) {
	if channel.HasSupportedState() {
		return channel.LatestSupportedState()
	}
	return channel.PreFundState(), nil
}

// getLedgerBalanceFromState returns the balance of the ledger channel from the given state
func getLedgerBalanceFromState(latest state.State) LedgerChannelBalance {
	// TODO: We assume single asset outcomes
	outcome := latest.Outcome[0]
	asset := outcome.Asset
	client := latest.Participants[0]
	clientBalance := big.NewInt(0).Set(outcome.Allocations[0].Amount)
	hub := latest.Participants[1]
	hubBalance := big.NewInt(0).Set(outcome.Allocations[1].Amount)

	return LedgerChannelBalance{
		AssetAddress:  asset,
		Hub:           hub,
		Client:        client,
		HubBalance:    hubBalance,
		ClientBalance: clientBalance,
	}
}

func isVirtualFundObjective(id types.Destination, store store.Store) (*virtualfund.Objective, bool) {
	// This is slightly awkward but if the virtual defunding objective is complete it won't come back if we query by channel id
	// We manually construct the objective id and query by that
	virtualFundId := protocols.ObjectiveId(virtualfund.ObjectivePrefix + id.String())
	o, err := store.GetObjectiveById(virtualFundId)
	if err != nil {
		return nil, false
	}
	return o.(*virtualfund.Objective), true
}

// GetPaymentChannelInfo returns the PaymentChannelInfo for the given channel
// It does this by querying the provided store and voucher manager
func GetPaymentChannelInfo(id types.Destination, store store.Store, vm *payments.VoucherManager) (PaymentChannelInfo, error) {
	// Otherwise we can just check the store
	c, channelFound := store.GetChannelById(id)

	if channelFound {
		status := getStatusFromChannel(c)

		// This means intermediaries may not have a fully signed postfund state even though the channel is "ready"
		// To determine the the correct status we check the status of the virtual fund objective
		fund, isVirtualFund := isVirtualFundObjective(id, store)
		if status == Proposed && isVirtualFund && fund.Status == protocols.Completed {
			status = Ready
		}

		latest, err := getLatestSupported(c)
		if err != nil {
			return PaymentChannelInfo{}, err
		}
		balance := getPaymentChannelBalance(c.Participants, latest.Outcome)

		// If we have received vouchers we want to update the channel balance to reflect the vouchers
		if hasVouchers := vm.ChannelRegistered(id); status == Ready && hasVouchers {

			paid, err := vm.Paid(id)
			if err != nil {
				return PaymentChannelInfo{}, err
			}
			balance.PaidSoFar.Set(paid)

			remaining, err := vm.Remaining(id)
			if err != nil {
				return PaymentChannelInfo{}, err
			}
			balance.RemainingFunds.Set(remaining)
		}

		return PaymentChannelInfo{
			ID:      id,
			Status:  status,
			Balance: balance,
		}, nil
	}
	return PaymentChannelInfo{}, fmt.Errorf("could not find channel with id %v", id)
}

// GetLedgerChannelInfo returns the LedgerChannelInfo for the given channel
// It does this by querying the provided store
func GetLedgerChannelInfo(id types.Destination, store store.Store) (LedgerChannelInfo, error) {
	c, ok := store.GetChannelById(id)
	if ok {
		return ConstructLedgerFromChannel(c), nil
	}

	con, err := store.GetConsensusChannelById(id)
	if err != nil {
		return LedgerChannelInfo{}, err
	}

	return ConstructLedgerInfoFromConsensus(con), nil
}

func ConstructLedgerInfoFromConsensus(con *consensus_channel.ConsensusChannel) LedgerChannelInfo {
	latest := con.ConsensusVars().AsState(con.FixedPart())
	return LedgerChannelInfo{
		ID:      con.Id,
		Status:  Ready,
		Balance: getLedgerBalanceFromState(latest),
	}
}

func ConstructLedgerFromChannel(c *channel.Channel) LedgerChannelInfo {
	latest, err := getLatestSupported(c)
	if err != nil {
		panic(err)
	}
	return LedgerChannelInfo{
		ID:      c.Id,
		Status:  getStatusFromChannel(c),
		Balance: getLedgerBalanceFromState(latest),
	}
}

func ConstructPaymentInfo(c *channel.Channel, objective protocols.Objective, paid, remaining *big.Int) (PaymentChannelInfo, error) {
	status := getStatusFromChannel(c)

	if objective != nil {
		// This means intermediaries may not have a fully signed postfund state even though the channel is "ready"
		// To determine the the correct status we check the status of the virtual fund objective
		fund, isVirtualFund := objective.(*virtualfund.Objective)
		if status == Proposed && isVirtualFund && fund.Status == protocols.Completed {
			status = Ready
		}
	}
	latest, err := getLatestSupported(c)
	if err != nil {
		return PaymentChannelInfo{}, err
	}
	balance := getPaymentChannelBalance(c.Participants, latest.Outcome)

	balance.PaidSoFar.Set(paid)

	balance.RemainingFunds.Set(remaining)

	return PaymentChannelInfo{
		ID:      c.Id,
		Status:  status,
		Balance: balance,
	}, nil
}
