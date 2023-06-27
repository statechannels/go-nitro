package query

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/statechannels/go-nitro/types"
)

type ChannelStatus string

// TODO: Think through statuses
const (
	Proposed ChannelStatus = "Proposed"
	Open     ChannelStatus = "Open"
	Closing  ChannelStatus = "Closing"
	Complete ChannelStatus = "Complete"
)

// PaymentChannelBalance contains the balance of a uni-directional payment channel
type PaymentChannelBalance struct {
	AssetAddress   types.Address
	Payee          types.Address
	Payer          types.Address
	PaidSoFar      *hexutil.Big
	RemainingFunds *hexutil.Big
}

// PaymentChannelInfo contains balance and status info about a payment channel
type PaymentChannelInfo struct {
	ID      types.Destination
	Status  ChannelStatus
	Balance PaymentChannelBalance
}

// LedgerChannelInfo contains balance and status info about a ledger channel
type LedgerChannelInfo struct {
	ID      types.Destination
	Status  ChannelStatus
	Balance LedgerChannelBalance
}

// LedgerChannelBalance contains the balance of a ledger channel
type LedgerChannelBalance struct {
	AssetAddress types.Address
	Me           types.Address
	Them         types.Address
	MyBalance    *hexutil.Big
	TheirBalance *hexutil.Big
}

// Equal returns true if the other LedgerChannelBalance is equal to this one
func (lcb LedgerChannelBalance) Equal(other LedgerChannelBalance) bool {
	return lcb.AssetAddress == other.AssetAddress &&
		lcb.Them == other.Them &&
		lcb.Me == other.Me &&
		lcb.TheirBalance.ToInt().Cmp(other.TheirBalance.ToInt()) == 0 &&
		lcb.MyBalance.ToInt().Cmp(other.MyBalance.ToInt()) == 0
}

// BalanceOf returns the balance of the given address in the channel, and an
// error if the address is not a participant.
func (lcb LedgerChannelBalance) BalanceOf(a types.Address) (*hexutil.Big, error) {
	if a == lcb.Me {
		return lcb.MyBalance, nil
	} else if a == lcb.Them {
		return lcb.TheirBalance, nil
	} else {
		return nil, fmt.Errorf("%s is not a participant", a)
	}
}

// Equal returns true if the other LedgerChannelInfo is equal to this one
func (li LedgerChannelInfo) Equal(other LedgerChannelInfo) bool {
	return li.ID == other.ID && li.Status == other.Status && li.Balance.Equal(other.Balance)
}

// Equal returns true if the other PaymentChannelInfo is equal to this one
func (pci PaymentChannelInfo) Equal(other PaymentChannelInfo) bool {
	return pci.ID == other.ID && pci.Status == other.Status && pci.Balance.Equal(other.Balance)
}

// Equal returns true if the other PaymentChannelBalance is equal to this one
func (pcb PaymentChannelBalance) Equal(other PaymentChannelBalance) bool {
	return pcb.AssetAddress == other.AssetAddress &&
		pcb.Payee == other.Payee &&
		pcb.Payer == other.Payer &&
		pcb.PaidSoFar.ToInt().Cmp(other.PaidSoFar.ToInt()) == 0 &&
		pcb.RemainingFunds.ToInt().Cmp(other.RemainingFunds.ToInt()) == 0
}
