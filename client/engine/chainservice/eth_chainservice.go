package chainservice

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	NitroAdjudicator "github.com/statechannels/go-nitro/client/engine/chainservice/adjudicator"
	Token "github.com/statechannels/go-nitro/client/engine/chainservice/erc20"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

var allocationUpdatedTopic = crypto.Keccak256Hash([]byte("AllocationUpdated(bytes32,uint256,uint256)"))
var concludedTopic = crypto.Keccak256Hash([]byte("Concluded(bytes32,uint48)"))
var depositedTopic = crypto.Keccak256Hash([]byte("Deposited(bytes32,address,uint256,uint256)"))

type ethChain interface {
	bind.ContractBackend
	ethereum.TransactionReader
	BlockNumber(ctx context.Context) (uint64, error)
}

type EthChainService struct {
	chain                    ethChain
	na                       *NitroAdjudicator.NitroAdjudicator
	naAddress                common.Address
	consensusAppAddress      common.Address
	virtualPaymentAppAddress common.Address
	txSigner                 *bind.TransactOpts
	out                      chan Event
	logger                   *log.Logger
}

// Since we only fetch events if there's a new block number
// it's safe to poll relatively frequently.
const EVENT_POLL_INTERVAL = 500 * time.Millisecond

// NewEthChainService constructs a chain service that submits transactions to a NitroAdjudicator
// and listens to events from an eventSource
func NewEthChainService(chain ethChain, na *NitroAdjudicator.NitroAdjudicator,
	naAddress, caAddress, vpaAddress common.Address, txSigner *bind.TransactOpts, logDestination io.Writer) (*EthChainService, error) {
	logPrefix := "chainservice " + txSigner.From.String() + ": "
	logger := log.New(logDestination, logPrefix, log.Lmicroseconds|log.Lshortfile)
	// Use a buffered channel so we don't have to worry about blocking on writing to the channel.
	ecs := EthChainService{chain, na, naAddress, caAddress, vpaAddress, txSigner, make(chan Event, 10), logger}

	go ecs.listenForLogEvents(context.Background())
	return &ecs, nil
}

// defaultTxOpts returns transaction options suitable for most transaction submissions
func (ecs *EthChainService) defaultTxOpts() *bind.TransactOpts {
	// FEVM only supports type 2 transactions
	return &bind.TransactOpts{
		From:      ecs.txSigner.From,
		Nonce:     ecs.txSigner.Nonce,
		Signer:    ecs.txSigner.Signer,
		GasFeeCap: ecs.txSigner.GasFeeCap,
		GasTipCap: ecs.txSigner.GasTipCap,
		GasLimit:  ecs.txSigner.GasLimit,
		GasPrice:  ecs.txSigner.GasPrice,
	}
}

// SendTransaction sends the transaction and blocks until it has been submitted.
func (ecs *EthChainService) SendTransaction(tx protocols.ChainTransaction) error {
	switch tx := tx.(type) {
	case protocols.DepositTransaction:
		for tokenAddress, amount := range tx.Deposit {
			txOpts := ecs.defaultTxOpts()
			ethTokenAddress := common.Address{}
			if tokenAddress == ethTokenAddress {
				txOpts.Value = amount
			} else {
				tokenTransactor, err := Token.NewTokenTransactor(tokenAddress, ecs.chain)
				if err != nil {
					return err
				}
				_, err = tokenTransactor.Approve(ecs.defaultTxOpts(), ecs.naAddress, amount)
				if err != nil {
					return err
				}
			}
			holdings, err := ecs.na.Holdings(&bind.CallOpts{}, tokenAddress, tx.ChannelId())
			if err != nil {
				return err
			}

			_, err = ecs.na.Deposit(txOpts, tokenAddress, tx.ChannelId(), holdings, amount)
			if err != nil {
				return err
			}
		}
		return nil
	case protocols.WithdrawAllTransaction:
		state := tx.SignedState.State()
		signatures := tx.SignedState.Signatures()
		nitroFixedPart := NitroAdjudicator.INitroTypesFixedPart(NitroAdjudicator.ConvertFixedPart(state.FixedPart()))
		nitroVariablePart := NitroAdjudicator.ConvertVariablePart(state.VariablePart())
		nitroSignatures := []NitroAdjudicator.INitroTypesSignature{NitroAdjudicator.ConvertSignature(signatures[0]), NitroAdjudicator.ConvertSignature(signatures[1])}
		proof := make([]NitroAdjudicator.INitroTypesSignedVariablePart, 0)
		candidate := NitroAdjudicator.INitroTypesSignedVariablePart{
			VariablePart: nitroVariablePart,
			Sigs:         nitroSignatures,
		}
		_, err := ecs.na.ConcludeAndTransferAllAssets(ecs.defaultTxOpts(), nitroFixedPart, proof, candidate)
		return err

	default:
		return fmt.Errorf("unexpected transaction type %T", tx)
	}
}

// dispatchChainEvents takes in a collection of event logs from the chain
// and dispatches events to the out channel
func (ecs *EthChainService) dispatchChainEvents(logs []ethTypes.Log) {
	for _, l := range logs {
		switch l.Topics[0] {
		case depositedTopic:
			nad, err := ecs.na.ParseDeposited(l)
			if err != nil {
				ecs.logger.Fatalf("error in ParseDeposited: %v", err)
			}

			event := NewDepositedEvent(nad.Destination, l.BlockNumber, nad.Asset, nad.AmountDeposited, nad.DestinationHoldings)
			ecs.out <- event
		case allocationUpdatedTopic:
			au, err := ecs.na.ParseAllocationUpdated(l)
			if err != nil {
				ecs.logger.Fatalf("error in ParseAllocationUpdated: %v", err)

			}

			tx, pending, err := ecs.chain.TransactionByHash(context.Background(), l.TxHash)
			if pending {
				ecs.logger.Fatalf("Expected transacion to be part of the chain, but the transaction is pending")

			}
			if err != nil {
				ecs.logger.Fatalf("error in TransactoinByHash: %v", err)

			}

			assetAddress, amount, err := getChainHolding(ecs.na, tx, au)
			if err != nil {
				ecs.logger.Fatalf("error in getChainHoldings: %v", err)

			}
			event := NewAllocationUpdatedEvent(au.ChannelId, l.BlockNumber, assetAddress, amount)
			ecs.out <- event
		case concludedTopic:
			ce, err := ecs.na.ParseConcluded(l)
			if err != nil {
				ecs.logger.Fatalf("error in ParseConcluded: %v", err)

			}

			event := ConcludedEvent{commonEvent: commonEvent{channelID: ce.ChannelId, BlockNum: l.BlockNumber}}
			ecs.out <- event

		default:
			ecs.logger.Printf("Unknown chain event")
		}
	}

}

// getCurrentBlockNumber returns the current block number.
func (ecs *EthChainService) getCurrentBlockNum() *big.Int {
	// TODO: We specifically call BlockNumber as HeaderByNumber can fail
	// see https://github.com/filecoin-project/ref-fvm/issues/1135
	blockNum, err := ecs.chain.BlockNumber(context.Background())
	if err != nil {
		panic(err)
	}

	return big.NewInt(int64(blockNum))
}

// listenForLogEvents periodically polls the chain client to check if there new events.
// If so it dispatches events to
func (ecs *EthChainService) listenForLogEvents(ctx context.Context) {

	// The initial query we want to get all events since the contract was deployed to the current block.
	query := ethereum.FilterQuery{
		Addresses: []common.Address{ecs.naAddress},
		FromBlock: nil,
		ToBlock:   ecs.getCurrentBlockNum(),
	}

	fetchedLogs, err := ecs.chain.FilterLogs(context.Background(), query)

	if err != nil {
		panic(err)
	}

	ecs.dispatchChainEvents(fetchedLogs)
	for {
		select {
		case <-time.After(EVENT_POLL_INTERVAL):
			currentBlock := ecs.getCurrentBlockNum()
			if moreRecentBlockAvailable := currentBlock.Cmp(query.ToBlock) > 0; moreRecentBlockAvailable {
				// We update to query to be between our previous block number and the latest block number
				query.FromBlock = big.NewInt(0).Set(query.ToBlock)
				query.ToBlock = big.NewInt(0).Set(currentBlock)

				fetchedLogs, err := ecs.chain.FilterLogs(context.Background(), query)

				if err != nil {
					panic(err)
				}

				ecs.dispatchChainEvents(fetchedLogs)
			}
		case <-ctx.Done():
			return
		}
	}

}

// EventFeed returns the out chan, and narrows the type so that external consumers may only receive on it.
func (ecs *EthChainService) EventFeed() <-chan Event {
	return ecs.out
}

func (ecs *EthChainService) GetConsensusAppAddress() types.Address {
	return ecs.consensusAppAddress
}

func (ecs *EthChainService) GetVirtualPaymentAppAddress() types.Address {
	return ecs.virtualPaymentAppAddress
}
