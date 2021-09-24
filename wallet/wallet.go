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
	CreateLedgerChannel(args CreateLedgerChannelArgs) (<-chan ChannelResult, error) // Returns a channel where the result will be sent
	CloseLedgerChannel()

	CreateVirtualChannel()
	CloseVirtualChannel()
}

// CapacityOracle exposes methods to query how many assets are locked into ledger channels, both by this wallet and by hubs
type CapacityOracle interface {
	// funds available to me in the ledger channel (that I can use to fund a new virtual channel)
	GetAvailableSpendCapacity(hub types.Address) (big.Int, error)

	// total user funds
	GetTotalSpendCapacity(intermediary types.Address) (big.Int, error)

	// funds available to the hub in the ledger channel (that they can use to collateralize a new virtual channel)
	GetAvailableReceiveCapacity(intermediary types.Address) (big.Int, error)
}

// CreateLedgerChannelArgs holds the data required for the CreateLedgerChannel API methods
type CreateLedgerChannelArgs struct {
	Asset             types.Address // Either the address of an ERC20 token OR the zero address (implying the native token).
	Hub               types.Address // The address corresponding to the channel signing key of a particular Nitro hub.
	MyBal             *big.Int      // My initial balance. Determines my spending capacity.
	HubBal            *big.Int      // The hub's initial balance. Determines my receive capacity.
	ChallengeDuration *big.Int      // The challenge timeout duration. Zero implies a default value will be chosen.
}

// AdjudicationStatus is an enumerated status for the adjudication of a channel. It mirrors on-chain storage.
type AdjudicationStatus int16

const (
	Challenge AdjudicationStatus = iota
	Open
	Finalized
)

// FundingStatus is an enumerated status for the funding of a channel. It is derived from either on-chain storage or the latest off-chain supported state.
type FundingStatus int16

const (
	NotYetSafeToFund FundingStatus = iota // There are insufficient funds to make my deposit safe
	SafeToDeposit                         // There are sufficient funds to make my deposit safe, but I have not yet deposited
	Deposited                             // I have already deposited, but the channel is not yet fully funded
	FullyFunded                           // The channel is fully funded on chain
	Exited                                // All of my funds have been exited from the channel on chain
)

type ChannelResult struct {
	ID                 types.Bytes32      // The channelId
	AdjudicationStatus AdjudicationStatus // The adjudication status of this channel.
	FundingStatus      FundingStatus      // The funding status of this channel.
}
