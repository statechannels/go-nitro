package chainservice

import (
	"context"
	"fmt"
	"io"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog"
	NitroAdjudicator "github.com/statechannels/go-nitro/client/engine/chainservice/adjudicator"
	Token "github.com/statechannels/go-nitro/client/engine/chainservice/erc20"
	chainutils "github.com/statechannels/go-nitro/client/engine/chainservice/utils"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

var (
	allocationUpdatedTopic   = crypto.Keccak256Hash([]byte("AllocationUpdated(bytes32,uint256,uint256)"))
	concludedTopic           = crypto.Keccak256Hash([]byte("Concluded(bytes32,uint48)"))
	depositedTopic           = crypto.Keccak256Hash([]byte("Deposited(bytes32,address,uint256,uint256)"))
	challengeRegisteredTopic = crypto.Keccak256Hash([]byte("ChallengeRegistered(bytes32 indexed channelId, uint48 turnNumRecord, uint48 finalizesAt, bool isFinal, (address[],uint64,address,uint48) fixedPart, (((address,(uint8,bytes),(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[])[] proof, (((address,(uint8,bytes),(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[]) candidate)"))
	challengeClearedTopic    = crypto.Keccak256Hash([]byte("ChallengeCleared(bytes32 indexed channelId, uint48 newTurnNumRecord)"))
)

type ethChain interface {
	bind.ContractBackend
	ethereum.TransactionReader
	ChainID(ctx context.Context) (*big.Int, error)
}

type EthChainService struct {
	chain                    ethChain
	na                       *NitroAdjudicator.NitroAdjudicator
	naAddress                common.Address
	consensusAppAddress      common.Address
	virtualPaymentAppAddress common.Address
	txSigner                 *bind.TransactOpts
	out                      chan Event
	logger                   zerolog.Logger
	ctx                      context.Context
	cancel                   context.CancelFunc
}

// MAX_QUERY_BLOCK_RANGE is the maximum range of blocks we query for events at once.
// Most json-rpc nodes restrict the amount of blocks you can search.
// For example Wallaby supports a maximum range of 2880
// See https://github.com/Zondax/rosetta-filecoin/blob/b395b3e04401be26c6cdf6a419e14ce85e2f7331/tools/wallaby/files/config.toml#L243
const MAX_QUERY_BLOCK_RANGE = 2000

// RESUB_INTERVAL is how often we resubscribe to log events.
// We do this to avoid https://github.com/ethereum/go-ethereum/issues/23845
// We use 2.5 minutes as the default filter timeout is 5 minutes.
// See https://github.com/ethereum/go-ethereum/blob/e14164d516600e9ac66f9060892e078f5c076229/eth/filters/filter_system.go#L43
const RESUB_INTERVAL = 2*time.Minute + 30*time.Second

// NewEthChainService is a convenient wrapper around NewEthChainService, which provides a simpler API
func NewEthChainService(chainUrl, chainPk string, naAddress, caAddress, vpaAddress common.Address, logDestination io.Writer) (*EthChainService, error) {
	if vpaAddress == caAddress {
		return nil, fmt.Errorf("virtual payment app address and consensus app address cannot be the same: %s", vpaAddress.String())
	}
	ethClient, txSigner, err := chainutils.ConnectToChain(context.Background(), chainUrl, common.Hex2Bytes(chainPk))
	if err != nil {
		panic(err)
	}

	na, err := NitroAdjudicator.NewNitroAdjudicator(naAddress, ethClient)
	if err != nil {
		panic(err)
	}

	return newEthChainService(ethClient, na, naAddress, caAddress, vpaAddress, txSigner, logDestination)
}

// NewEthChainService constructs a chain service that submits transactions to a NitroAdjudicator
// and listens to events from an eventSource
func newEthChainService(chain ethChain, na *NitroAdjudicator.NitroAdjudicator,
	naAddress, caAddress, vpaAddress common.Address, txSigner *bind.TransactOpts, logDestination io.Writer,
) (*EthChainService, error) {
	logger := zerolog.New(logDestination).With().Timestamp().Str("txSigner", txSigner.From.String()[0:8]).Caller().Logger()
	ctx, cancelCtx := context.WithCancel(context.Background())

	// Use a buffered channel so we don't have to worry about blocking on writing to the channel.
	ecs := EthChainService{chain, na, naAddress, caAddress, vpaAddress, txSigner, make(chan Event, 10), logger, ctx, cancelCtx}
	errChan, err := ecs.subscribeForLogs()
	// TODO: Return error from chain service instead of panicking
	go func() {
		for err := range errChan {
			// Print to STDOUT in case we're using a noop logger
			fmt.Println(err)
			ecs.logger.Fatal().Err(err)
			// Manually panic in case we're using a logger that doesn't call exit(1)
			panic(err)
		}
	}()
	if err != nil {
		return nil, err
	}

	return &ecs, nil
}

// defaultTxOpts returns transaction options suitable for most transaction submissions
func (ecs *EthChainService) defaultTxOpts() *bind.TransactOpts {
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

		candidate := NitroAdjudicator.INitroTypesSignedVariablePart{
			VariablePart: nitroVariablePart,
			Sigs:         nitroSignatures,
		}
		_, err := ecs.na.ConcludeAndTransferAllAssets(ecs.defaultTxOpts(), nitroFixedPart, candidate)
		return err

	default:
		return fmt.Errorf("unexpected transaction type %T", tx)
	}
}

// dispatchChainEvents takes in a collection of event logs from the chain
// and dispatches events to the out channel
func (ecs *EthChainService) dispatchChainEvents(logs []ethTypes.Log) error {
	for _, l := range logs {
		switch l.Topics[0] {
		case depositedTopic:
			nad, err := ecs.na.ParseDeposited(l)
			if err != nil {
				return fmt.Errorf("error in ParseDeposited: %w", err)
			}

			event := NewDepositedEvent(nad.Destination, l.BlockNumber, nad.Asset, nad.AmountDeposited, nad.DestinationHoldings)
			ecs.out <- event
		case allocationUpdatedTopic:
			au, err := ecs.na.ParseAllocationUpdated(l)
			if err != nil {
				return fmt.Errorf("error in ParseAllocationUpdated: %w", err)
			}

			tx, pending, err := ecs.chain.TransactionByHash(context.Background(), l.TxHash)
			if pending {
				return fmt.Errorf("expected transaction to be part of the chain, but the transaction is pending")
			}
			var assetAddress types.Address
			var amount *big.Int

			if err != nil {
				return fmt.Errorf("error in TransactionByHash: %w", err)
			}

			assetAddress, amount, err = getChainHolding(ecs.na, tx, au)
			if err != nil {
				return fmt.Errorf("error in getChainHoldings: %w", err)
			}

			event := NewAllocationUpdatedEvent(au.ChannelId, l.BlockNumber, assetAddress, amount)
			ecs.out <- event
		case concludedTopic:
			ce, err := ecs.na.ParseConcluded(l)
			if err != nil {
				return fmt.Errorf("error in ParseConcluded: %w", err)
			}

			event := ConcludedEvent{commonEvent: commonEvent{channelID: ce.ChannelId, BlockNum: l.BlockNumber}}
			ecs.out <- event

		case challengeRegisteredTopic:
			ecs.logger.Info().Msg("Ignoring Challenge Registered event")
		case challengeClearedTopic:
			ecs.logger.Info().Msg("Ignoring Challenge Cleared event")
		default:
			ecs.logger.Info().Str("topic", l.Topics[0].String()).Msg("Ignoring unknown chain event topic")
		}
	}
	return nil
}

// subscribeForLogs subscribes for logs and pushes them to the out channel.
// It relies on notifications being supported by the chain node.
func (ecs *EthChainService) subscribeForLogs() (<-chan error, error) {
	// Subscribe to Adjudicator events
	query := ethereum.FilterQuery{
		Addresses: []common.Address{ecs.naAddress},
	}
	logs := make(chan ethTypes.Log)
	sub, err := ecs.chain.SubscribeFilterLogs(ecs.ctx, query, logs)
	if err != nil {
		return nil, fmt.Errorf("subscribeFilterLogs failed: %w", err)
	}
	errorChan := make(chan error)
	// Must be in a goroutine to not block chain service constructor
	go func() {
	out:
		for {
			select {
			case <-ecs.ctx.Done():
				sub.Unsubscribe()
				return
			case err := <-sub.Err():
				if err != nil {
					errorChan <- fmt.Errorf("received error from the subscription channel: %w", err)
					break out
				}

				// If the error is nil then the subscription was closed and we need to re-subscribe.
				// This is a workaround for https://github.com/ethereum/go-ethereum/issues/23845
				var sErr error
				sub, sErr = ecs.chain.SubscribeFilterLogs(ecs.ctx, query, logs)
				if sErr != nil {
					errorChan <- fmt.Errorf("subscribeFilterLogs failed on resubscribe: %w", err)
					break out
				}
				ecs.logger.Print("resubscribed to filtered logs")

			case <-time.After(RESUB_INTERVAL):
				// Due to https://github.com/ethereum/go-ethereum/issues/23845 we can't rely on a long running subscription.
				// We unsub here and recreate the subscription in the next iteration of the select.
				sub.Unsubscribe()
			case chainEvent := <-logs:
				err = ecs.dispatchChainEvents([]ethTypes.Log{chainEvent})
				if err != nil {
					errorChan <- fmt.Errorf("error in dispatchChainEvents: %w", err)
					break out
				}
			}
		}
	}()

	return errorChan, nil
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

func (ecs *EthChainService) GetChainId() (*big.Int, error) {
	return ecs.chain.ChainID(context.Background())
}

func (ecs *EthChainService) Close() error {
	ecs.cancel()
	return nil
}
