package query

import (
	"math/big"

	"github.com/statechannels/go-nitro/channel"
	"github.com/statechannels/go-nitro/types"
)

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
	Status  channel.ChannelStatus
	Balance PaymentChannelBalance
}

// LedgerChannelInfo contains balance and status info about a ledger channel
type LedgerChannelInfo struct {
	ID      types.Destination
	Status  channel.ChannelStatus
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

// Equal returns true if the other LedgerChannelBalance is equal to this one
func (lcb LedgerChannelBalance) Equal(other LedgerChannelBalance) bool {
	return lcb.AssetAddress == other.AssetAddress &&
		lcb.Hub == other.Hub &&
		lcb.Client == other.Client &&
		lcb.HubBalance.Cmp(other.HubBalance) == 0 &&
		lcb.ClientBalance.Cmp(other.ClientBalance) == 0
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
		pcb.PaidSoFar.Cmp(other.PaidSoFar) == 0 &&
		pcb.RemainingFunds.Cmp(other.RemainingFunds) == 0
}
