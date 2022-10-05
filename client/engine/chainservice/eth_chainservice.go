package chainservice

import (
	"context"
	"fmt"
	"io"
	"log"
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

// RESUB_INTERVAL is how often we resubscribe to log events.
// We do this to avoid https://github.com/ethereum/go-ethereum/issues/23845
const RESUB_INTERVAL = 60 * time.Second

// NewEthChainService constructs a chain service that submits transactions to a NitroAdjudicator
// and listens to events from an eventSource
func NewEthChainService(chain ethChain, na *NitroAdjudicator.NitroAdjudicator,
	naAddress, caAddress, vpaAddress common.Address, txSigner *bind.TransactOpts, logDestination io.Writer) (*EthChainService, error) {
	logPrefix := "chainservice " + txSigner.From.String() + ": "
	logger := log.New(logDestination, logPrefix, log.Lmicroseconds|log.Lshortfile)
	// Use a buffered channel so we don't have to worry about blocking on writing to the channel.
	ecs := EthChainService{chain, na, naAddress, caAddress, vpaAddress, txSigner, make(chan Event, 10), logger}

	err := ecs.subcribeToEvents()
	return &ecs, err
}

// defaultTxOpts returns transaction options suitable for most transaction submissions
func (ecs *EthChainService) defaultTxOpts() *bind.TransactOpts {
	return &bind.TransactOpts{
		From:     ecs.txSigner.From,
		Nonce:    ecs.txSigner.Nonce,
		Signer:   ecs.txSigner.Signer,
		GasPrice: ecs.txSigner.GasPrice,
		GasLimit: ecs.txSigner.GasLimit,
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

func (ecs *EthChainService) subcribeToEvents() error {

	go ecs.listenForLogEvents()
	return nil
}

func (ecs *EthChainService) listenForLogEvents() {
	// Subsribe to Adjudicator events
	query := ethereum.FilterQuery{
		Addresses: []common.Address{ecs.naAddress},
	}
	logs := make(chan ethTypes.Log)
	sub, err := ecs.chain.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		panic(err)
	}
	for {
		select {
		case err := <-sub.Err():
			if err != nil {
				panic(err)
			}

			// If the error is nil then the subscription was closed and we need to re-subscribe.
			// This is a workaround for https://github.com/ethereum/go-ethereum/issues/23845
			var sErr error
			sub, sErr = ecs.chain.SubscribeFilterLogs(context.Background(), query, logs)
			if sErr != nil {
				panic(err)
			}
			ecs.logger.Println("resubscribed to filtered logs")

		case <-time.After(RESUB_INTERVAL):
			// Due to https://github.com/ethereum/go-ethereum/issues/23845 we can't rely on a long running subscription.
			// We unsub here and recreate the subscription in the next iteration of the select.
			sub.Unsubscribe()
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

func (ecs *EthChainService) GetConsensusAppAddress() types.Address {
	return ecs.consensusAppAddress
}

func (ecs *EthChainService) GetVirtualPaymentAppAddress() types.Address {
	return ecs.virtualPaymentAppAddress
}
