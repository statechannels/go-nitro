// Package wallet contains the interfaces for the components of a nitro wallet
package wallet

import (
	"math/big"

	"github.com/statechannels/go-nitro/types"
)

// Wallet exposes a simple API to the consuming program
type Wallet interface {
	ChannelManager
	CapacityOracle
}

// ChannelManager accepts requests to open and close channels
type ChannelManager interface {
	CreateLedgerChannel()
	CreateVirtualChannel()
	CloseVirtualChannel()
}

// CapacityOracle exposes methods to query how many assets are locked into ledger channels, both by this wallet and by hubs
type CapacityOracle interface {
	// funds available to me in the ledger channel (that I can use to fund a new virtual channel)
	GetAvailableSpendCapacity(intermediary types.Address) (big.Int, error)

	// total user funds
	GetTotalSpendCapacity(intermediary types.Address) (big.Int, error)

	// funds available to the hub in the ledger channel (that they can use to collateralize a new virtual channel)
	GetAvailableReceiveCapacity(intermediary types.Address) (big.Int, error)
}
