// Package wallet contains the interfaces for the components of a nitro wallet
package wallet // import "github.com/statechannels/go-nitro/wallet"

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
	// CreateLedgerChannel Requests a new ledger channel with a given hub and initial balances. Returns a go channel where the result will be sent.
	CreateLedgerChannel(args CreateLedgerChannelArgs) (<-chan ChannelResult, error)
	// CloseLedgerChannel Requests an existing ledger channel with given ID be closed and defunded on chain. Returns a go channel where the result will be sent.
	CloseLedgerChannel(ID types.Bytes32) (<-chan ChannelResult, error)

	// CreateVirtualChannel Requests a new virtual channel through a given hub and initial balances. Returns a go channel where the result will be sent.
	CreateVirtualChannel(args CreateVirtualChannelArgs) (<-chan ChannelResult, error)
	// CloseVirtualChannel Requests an existing virtual channel with given ID be closed and defunded off chain. Returns a go channel where the result will be sent.
	CloseVirtualChannel(ID types.Bytes32) (<-chan ChannelResult, error)

	BytecodeCacher

	ApproveObjective(ID string)
}

// CreateLedgerChannelArgs holds the data required for the CreateLedgerChannel API method
type CreateLedgerChannelArgs struct {
	Asset             types.Address // Either the address of an ERC20 token OR the zero address (implying the native token).
	Hub               types.Address // The address corresponding to the channel signing key of a particular Nitro hub.
	MyBal             *big.Int      // My initial balance. Determines my spending capacity.
	HubBal            *big.Int      // The hub's initial balance. Determines my receive capacity.
	ChallengeDuration *big.Int      // The challenge timeout duration. Zero implies a default value will be chosen.
}

// CreateVirtualChannelArgs holds the data required for the CreateVirtualChannel API method
type CreateVirtualChannelArgs struct {
	Asset             types.Address // Either the address of an ERC20 token OR the zero address (implying the native token).
	Hub               types.Address // The address corresponding to the channel signing key of a particular Nitro hub.
	Counterparty      types.Address // The address corresponding to the channel signing key of a particular Nitro participant connected to the same hub.
	MyBal             *big.Int      // My initial balance. Determines my spending capacity.
	CounterpartyBal   *big.Int      // The hub's initial balance. Determines my receive capacity.
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

// CapacityOracle exposes methods to query how many assets are locked into ledger channels, both by this wallet and by the hub's.
type CapacityOracle interface {
	// GetCapacities returns information about the current outcome of the ledger channel with the given hub, in a mapping keyed by the asset.
	GetCapacities(hub types.Address) (map[types.Address]LedgerChannelCapacites, error)
}

type LedgerChannelCapacites struct {
	FreeSendable   *big.Int // Funds available to me in the ledger channel (that I can use to fund a new virtual channel).
	FreeReceivable *big.Int // Funds available to the hub in the ledger channel (that they can use to collateralize a new virtual channel).

	LockedForMe  *big.Int // Funds locked in virtual channels that will return to me.
	LockedForHub *big.Int // Funds locked in virtual channels that will not return to me.
}
