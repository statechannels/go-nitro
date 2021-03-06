package chainservice

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	NitroAdjudicator "github.com/statechannels/go-nitro/client/engine/chainservice/adjudicator"
	Token "github.com/statechannels/go-nitro/client/engine/chainservice/erc20"
	"github.com/statechannels/go-nitro/protocols"
)

var allocationUpdatedTopic = crypto.Keccak256Hash([]byte("AllocationUpdated(bytes32,uint256,uint256)"))
var concludedTopic = crypto.Keccak256Hash([]byte("Concluded(bytes32,uint48)"))
var depositedTopic = crypto.Keccak256Hash([]byte("Deposited(bytes32,address,uint256,uint256)"))

type ethChain interface {
	bind.ContractBackend
	ethereum.TransactionReader
}

type EthChainService struct {
	chain               ethChain
	na                  *NitroAdjudicator.NitroAdjudicator
	naAddress           common.Address
	consensusAppAddress common.Address
	txSigner            *bind.TransactOpts
	out                 chan Event
	logger              *log.Logger
}

// NewEthChainService constructs a chain service that submits transactions to a NitroAdjudicator
// and listens to events from an eventSource
func NewEthChainService(chain ethChain, na *NitroAdjudicator.NitroAdjudicator,
	naAddress common.Address, caAddress common.Address, txSigner *bind.TransactOpts, logDestination io.Writer) (*EthChainService, error) {
	logPrefix := "chainservice " + txSigner.From.String() + ": "
	logger := log.New(logDestination, logPrefix, log.Lmicroseconds|log.Lshortfile)
	// Use a buffered channel so we don't have to worry about blocking on writing to the channel.
	ecs := EthChainService{chain, na, naAddress, caAddress, txSigner, make(chan Event, 10), logger}

	err := ecs.subcribeToEvents()
	return &ecs, err
}

// defaultTxOpts returns transaction options suitable for most transaction submissions
func (ecs *EthChainService) defaultTxOpts() *bind.TransactOpts {
	return &bind.TransactOpts{
		From:   ecs.txSigner.From,
		Nonce:  ecs.txSigner.Nonce,
		Signer: ecs.txSigner.Signer,
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
		nitroFixedPart := NitroAdjudicator.INitroTypesFixedPart(state.FixedPart())
		nitroVariablePart := NitroAdjudicator.ConvertVariablePart(state.VariablePart())
		nitroSignatures := []NitroAdjudicator.INitroTypesSignature{NitroAdjudicator.ConvertSignature(signatures[0]), NitroAdjudicator.ConvertSignature(signatures[1])}
		nitroSignedVariableParts := []NitroAdjudicator.INitroTypesSignedVariablePart{{
			VariablePart: nitroVariablePart,
			Sigs:         nitroSignatures,
		}}
		_, err := ecs.na.ConcludeAndTransferAllAssets(ecs.defaultTxOpts(), nitroFixedPart, nitroSignedVariableParts)
		return err

	default:
		return fmt.Errorf("unexpected transaction type %T", tx)
	}
}

func (ecs *EthChainService) subcribeToEvents() error {
	// Subsribe to Adjudicator events
	query := ethereum.FilterQuery{
		Addresses: []common.Address{ecs.naAddress},
	}
	logs := make(chan ethTypes.Log)
	sub, err := ecs.chain.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		return err
	}
	go ecs.listenForLogEvents(sub, logs)
	return nil
}

func (ecs *EthChainService) listenForLogEvents(sub ethereum.Subscription, logs chan ethTypes.Log) {
	for {
		select {
		case err := <-sub.Err():
			// TODO should we try resubscribing to chain events
			ecs.logger.Printf("event subscription error: %v", err)
		case chainEvent := <-logs:
			switch chainEvent.Topics[0] {
			case depositedTopic:
				nad, err := ecs.na.ParseDeposited(chainEvent)
				if err != nil {
					ecs.logger.Printf("error in ParseDeposited: %v", err)
				}

				event := NewDepositedEvent(nad.Destination, chainEvent.BlockNumber, nad.Asset, nad.AmountDeposited, nad.DestinationHoldings)
				ecs.out <- event
			case allocationUpdatedTopic:
				au, err := ecs.na.ParseAllocationUpdated(chainEvent)
				if err != nil {
					ecs.logger.Printf("error in ParseAllocationUpdated: %v", err)
				}

				tx, pending, err := ecs.chain.TransactionByHash(context.Background(), chainEvent.TxHash)
				if pending {
					ecs.logger.Printf("Expected transacion to be part of the chain, but the transaction is pending")
				}
				if err != nil {
					ecs.logger.Printf("error in TransactoinByHash: %v", err)
				}

				assetAddress, amount, err := getChainHolding(ecs.na, tx, au)
				if err != nil {
					ecs.logger.Printf("error in getChainHoldings: %v", err)
				}
				event := NewAllocationUpdatedEvent(au.ChannelId, chainEvent.BlockNumber, assetAddress, amount)
				ecs.out <- event
			case concludedTopic:
				ce, err := ecs.na.ParseConcluded(chainEvent)
				if err != nil {
					ecs.logger.Printf("error in ParseConcluded: %v", err)
				}

				event := ConcludedEvent{commonEvent: commonEvent{channelID: ce.ChannelId, BlockNum: chainEvent.BlockNumber}}
				ecs.out <- event
			default:
				ecs.logger.Printf("Unknown chain event")
			}
		}
	}
}

// EventFeed returns the out chan, and narrows the type so that external consumers may only receive on it.
func (ecs *EthChainService) EventFeed() <-chan Event {
	return ecs.out
}
