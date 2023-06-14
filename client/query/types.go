package query

import (
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

type PaymentReceiptStatus string

const (
	PRSreceived        PaymentReceiptStatus = "Received"
	PRSmisaddressed    PaymentReceiptStatus = "Misaddressed"
	PRSchannelNotFound PaymentReceiptStatus = "ChannelNotFound"
	PRSincorrectSigner PaymentReceiptStatus = "IncorrectSigner"
	PRSengineError     PaymentReceiptStatus = "EngineError"
)

// PaymentChannelBalance contains the balance of a uni-directional payment channel
type PaymentChannelBalance struct {
	AssetAddress   types.Address
	Payee          types.Address
	Payer          types.Address
	PaidSoFar      *hexutil.Big
	RemainingFunds *hexutil.Big
}

// PaymentChannelPaymentReceipt contains the net-result of applying a voucher
// against a payment channel, or an error status if the voucher was invalid.
type PaymentChannelPaymentReceipt struct {
	ID             types.Destination
	AmountReceived *hexutil.Big
	Status         PaymentReceiptStatus
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
	AssetAddress  types.Address
	Hub           types.Address
	Client        types.Address
	HubBalance    *hexutil.Big
	ClientBalance *hexutil.Big
}

// Equal returns true if the other LedgerChannelBalance is equal to this one
func (lcb LedgerChannelBalance) Equal(other LedgerChannelBalance) bool {
	return lcb.AssetAddress == other.AssetAddress &&
		lcb.Hub == other.Hub &&
		lcb.Client == other.Client &&
		lcb.HubBalance.ToInt().Cmp(other.HubBalance.ToInt()) == 0 &&
		lcb.ClientBalance.ToInt().Cmp(other.ClientBalance.ToInt()) == 0
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
