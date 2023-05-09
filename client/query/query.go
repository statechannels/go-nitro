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

// GetVirtualFundObjective returns the virtual fund objective for the given channel if it exists.
func GetVirtualFundObjective(id types.Destination, store store.Store) (*virtualfund.Objective, bool) {
	// This is slightly awkward but if the virtual defunding objective is complete it won't come back if we query by channel id
	// We manually construct the objective id and query by that
	virtualFundId := protocols.ObjectiveId(virtualfund.ObjectivePrefix + id.String())
	o, err := store.GetObjectiveById(virtualFundId)
	if err != nil {
		return nil, false
	}
	return o.(*virtualfund.Objective), true
}

// GetVoucherBalance returns the amount paid and remaining for a given channel based on vouchers received.
// If not vouchers are received for the channel, it returns 0 for paid and remaining.
func GetVoucherBalance(id types.Destination, vm *payments.VoucherManager) (paid, remaining *big.Int, err error) {
	paid, remaining = big.NewInt(0), big.NewInt(0)

	if noVouchers := !vm.ChannelRegistered(id); noVouchers {
		return
	}
	paid, err = vm.Paid(id)
	if err != nil {
		return nil, nil, err
	}

	remaining, err = vm.Remaining(id)
	if err != nil {
		return nil, nil, err
	}

	return paid, remaining, nil
}

// GetPaymentChannelInfo returns the PaymentChannelInfo for the given channel
// It does this by querying the provided store and voucher manager
func GetPaymentChannelInfo(id types.Destination, store store.Store, vm *payments.VoucherManager) (PaymentChannelInfo, error) {
	// Otherwise we can just check the store
	c, channelFound := store.GetChannelById(id)

	if channelFound {

		paid, remaining, err := GetVoucherBalance(id, vm)
		if err != nil {
			return PaymentChannelInfo{}, err
		}
		o, _ := GetVirtualFundObjective(id, store)
		return ConstructPaymentInfo(c, o, paid, remaining)
	}
	return PaymentChannelInfo{}, fmt.Errorf("could not find channel with id %v", id)
}

func GetAllLedgerChannels(store store.Store, consensusAppDefinition types.Address) ([]LedgerChannelInfo, error) {
	toReturn := []LedgerChannelInfo{}

	allConsensus, err := store.GetAllConsensusChannels()
	if err != nil {
		return []LedgerChannelInfo{}, err
	}
	for _, con := range allConsensus {
		toReturn = append(toReturn, ConstructLedgerInfoFromConsensus(con))
	}
	allChannels := store.GetChannelsByAppDefinition(consensusAppDefinition)
	for _, c := range allChannels {
		toReturn = append(toReturn, ConstructLedgerInfoFromChannel(c))
	}
	return toReturn, nil
}

func GetPaymentChannelsByLedger(ledgerId types.Destination, store store.Store, vm *payments.VoucherManager) ([]PaymentChannelInfo, error) {
	// If a ledger channel is actively funding payment channels it must be in the form of a consensus channel
	con, err := store.GetConsensusChannelById(ledgerId)
	if err != nil {
		return []PaymentChannelInfo{}, fmt.Errorf("could not find any payment channels funded by %s: %w", ledgerId, err)
	}

	toQuery := con.ConsensusVars().Outcome.FundingTargets()

	paymentChannels, err := store.GetChannelsByIds(toQuery)
	if err != nil {
		return []PaymentChannelInfo{}, fmt.Errorf("could not query the store about ids %v: %w", toQuery, err)
	}
	objectives, err := store.GetObjectiveByChannelIds(toQuery)
	if err != nil {
		return []PaymentChannelInfo{}, fmt.Errorf("could not query the store about ids %v: %w", toQuery, err)
	}

	toReturn := []PaymentChannelInfo{}
	for _, p := range paymentChannels {
		paid, remaining, err := GetVoucherBalance(p.Id, vm)
		o := objectives[p.Id]

		if err != nil {
			return []PaymentChannelInfo{}, err
		}

		vfo, _ := o.(*virtualfund.Objective)
		info, err := ConstructPaymentInfo(p, vfo, paid, remaining)
		if err != nil {
			return []PaymentChannelInfo{}, err
		}
		toReturn = append(toReturn, info)
	}
	return toReturn, nil
}

// GetLedgerChannelInfo returns the LedgerChannelInfo for the given channel
// It does this by querying the provided store
func GetLedgerChannelInfo(id types.Destination, store store.Store) (LedgerChannelInfo, error) {
	c, ok := store.GetChannelById(id)
	if ok {
		return ConstructLedgerInfoFromChannel(c), nil
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

func ConstructLedgerInfoFromChannel(c *channel.Channel) LedgerChannelInfo {
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

func ConstructPaymentInfo(c *channel.Channel, vfo *virtualfund.Objective, paid, remaining *big.Int) (PaymentChannelInfo, error) {
	status := getStatusFromChannel(c)

	if vfo != nil && vfo.Status == protocols.Completed {
		// This means intermediaries may not have a fully signed postfund state even though the channel is "ready"
		// To determine the the correct status we check the status of the virtual fund objective

		status = Ready
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
