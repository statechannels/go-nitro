package client

import (
	"math/big"

	"github.com/statechannels/go-nitro/types"
)

type ChannelStatus string

const Proposed ChannelStatus = "Proposed"
const Ready ChannelStatus = "Ready"
const Closing ChannelStatus = "Closing"
const Complete ChannelStatus = "Complete"

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
