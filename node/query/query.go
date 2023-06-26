package query

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/statechannels/go-nitro/channel"
	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/node/engine/store"
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
	return Open
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
		PaidSoFar:      (*hexutil.Big)(paidSoFar),
		RemainingFunds: (*hexutil.Big)(remaining),
	}
}

// getLatestSupportedOrPreFund returns the latest supported state of the channel
// or the prefund state if no supported state exists
func getLatestSupportedOrPreFund(channel *channel.Channel) (state.State, error) {
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
		HubBalance:    (*hexutil.Big)(hubBalance),
		ClientBalance: (*hexutil.Big)(clientBalance),
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
	if (id == types.Destination{}) {
		err := types.InvalidParamsError
		err.Message = "a valid channel id must be provided"
		return PaymentChannelInfo{}, err
	}
	// Otherwise we can just check the store
	c, channelFound := store.GetChannelById(id)

	if channelFound {
		paid, remaining, err := GetVoucherBalance(id, vm)
		if err != nil {
			return PaymentChannelInfo{}, err
		}

		return ConstructPaymentInfo(c, paid, remaining)
	}
	err := types.InvalidParamsError
	err.Message = fmt.Sprintf("Could not find channel with id %v", id)
	return PaymentChannelInfo{}, err
}

// GetAllLedgerChannels returns a `LedgerChannelInfo` for each ledger channel in the store.
func GetAllLedgerChannels(store store.Store, consensusAppDefinition types.Address) ([]LedgerChannelInfo, error) {
	toReturn := []LedgerChannelInfo{}

	allConsensus, err := store.GetAllConsensusChannels()
	if err != nil {
		return []LedgerChannelInfo{}, err
	}
	for _, con := range allConsensus {
		toReturn = append(toReturn, ConstructLedgerInfoFromConsensus(con))
	}
	allChannels, err := store.GetChannelsByAppDefinition(consensusAppDefinition)
	if err != nil {
		return []LedgerChannelInfo{}, err
	}
	for _, c := range allChannels {
		l, err := ConstructLedgerInfoFromChannel(c)
		if err != nil {
			return []LedgerChannelInfo{}, err
		}
		toReturn = append(toReturn, l)
	}
	return toReturn, nil
}

// GetPaymentChannelsByLedger returns a `PaymentChannelInfo` for each active payment channel funded by the given ledger channel.
func GetPaymentChannelsByLedger(ledgerId types.Destination, s store.Store, vm *payments.VoucherManager) ([]PaymentChannelInfo, error) {
	// If a ledger channel is actively funding payment channels it must be in the form of a consensus channel
	con, err := s.GetConsensusChannelById(ledgerId)
	// If the ledger channel is not a consensus channel we know that there are no payment channels funded by it
	if errors.Is(err, store.ErrNoSuchChannel) {
		return []PaymentChannelInfo{}, nil
	}
	if err != nil {
		return []PaymentChannelInfo{}, fmt.Errorf("could not find any payment channels funded by %s: %w", ledgerId, err)
	}

	toQuery := con.ConsensusVars().Outcome.FundingTargets()

	paymentChannels, err := s.GetChannelsByIds(toQuery)
	if err != nil {
		return []PaymentChannelInfo{}, fmt.Errorf("could not query the store about ids %v: %w", toQuery, err)
	}

	toReturn := []PaymentChannelInfo{}
	for _, p := range paymentChannels {
		paid, remaining, err := GetVoucherBalance(p.Id, vm)
		if err != nil {
			return []PaymentChannelInfo{}, err
		}

		info, err := ConstructPaymentInfo(p, paid, remaining)
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
		return ConstructLedgerInfoFromChannel(c)
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
		Status:  Open,
		Balance: getLedgerBalanceFromState(latest),
	}
}

func ConstructLedgerInfoFromChannel(c *channel.Channel) (LedgerChannelInfo, error) {
	latest, err := getLatestSupportedOrPreFund(c)
	if err != nil {
		return LedgerChannelInfo{}, err
	}
	return LedgerChannelInfo{
		ID:      c.Id,
		Status:  getStatusFromChannel(c),
		Balance: getLedgerBalanceFromState(latest),
	}, nil
}

func ConstructPaymentInfo(c *channel.Channel, paid, remaining *big.Int) (PaymentChannelInfo, error) {
	status := getStatusFromChannel(c)
	// ADR 0009 allows for intermediaries to exit the protocol before receiving all signed post funds
	// So for intermediaries we return Open once they have signed their post fund state
	amIntermediary := c.MyIndex != 0 && c.MyIndex != uint(len(c.Participants)-1)
	if amIntermediary && c.PostFundSignedByMe() {
		status = Open
	}

	latest, err := getLatestSupportedOrPreFund(c)
	if err != nil {
		return PaymentChannelInfo{}, err
	}
	balance := getPaymentChannelBalance(c.Participants, latest.Outcome)

	balance.PaidSoFar.ToInt().Set(paid)

	balance.RemainingFunds.ToInt().Set(remaining)

	return PaymentChannelInfo{
		ID:      c.Id,
		Status:  status,
		Balance: balance,
	}, nil
}
