package policy

import (
	"github.com/statechannels/go-nitro/payments"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/protocols/directfund"
	"github.com/statechannels/go-nitro/protocols/virtualfund"
	"github.com/statechannels/go-nitro/types"
)

func NewLedgerChannelMaxSpendPolicy(me types.Address, maxAmount types.Funds) PolicyMaker {
	return &maxSpendPolicy{me: me, maxAmount: maxAmount, channelType: ledger}
}

func NewPaymentChannelMaxSpendPolicy(me types.Address, maxAmount types.Funds) PolicyMaker {
	return &maxSpendPolicy{me: me, maxAmount: maxAmount, channelType: payment}
}

type channelType string

const (
	ledger  channelType = "ledger"
	payment channelType = "payment"
)

type maxSpendPolicy struct {
	me          types.Address
	maxAmount   types.Funds
	channelType channelType
}

func (mp *maxSpendPolicy) ShouldApprove(o protocols.Objective) bool {
	switch obj := o.(type) {
	case *virtualfund.Objective:
		if mp.channelType == payment && obj.MyRole == payments.PAYER_INDEX {
			myAmount := obj.V.PreFundState().Outcome.TotalAllocatedFor(types.AddressToDestination(mp.me))

			if !underMaxAmount(myAmount, mp.maxAmount) {
				return false
			}
		}
	case *directfund.Objective:
		if mp.channelType == ledger {
			myAmount := obj.C.PreFundState().Outcome.TotalAllocatedFor(types.AddressToDestination(mp.me))
			if !underMaxAmount(myAmount, mp.maxAmount) {
				return false
			}
		}

	}

	return true
}

func underMaxAmount(myAmount types.Funds, maxAmount types.Funds) bool {
	for asset, amt := range myAmount {
		if amt.Cmp(maxAmount[asset]) > 0 {
			return false
		}
	}
	return true
}
