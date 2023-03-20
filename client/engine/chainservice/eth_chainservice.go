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
	"github.com/statechannels/go-nitro/internal/logging"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"

	"github.com/ethereum/go-ethereum/rpc"
)

var allocationUpdatedTopic = crypto.Keccak256Hash([]byte("AllocationUpdated(bytes32,uint256,uint256)"))
var concludedTopic = crypto.Keccak256Hash([]byte("Concluded(bytes32,uint48)"))
var depositedTopic = crypto.Keccak256Hash([]byte("Deposited(bytes32,address,uint256,uint256)"))
var challengeRegisteredTopic = crypto.Keccak256Hash([]byte("ChallengeRegistered(bytes32 indexed channelId, uint48 turnNumRecord, uint48 finalizesAt, bool isFinal, (address[],uint64,address,uint48) fixedPart, (((address,(uint8,bytes),(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[])[] proof, (((address,(uint8,bytes),(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[]) candidate)"))
var challengeClearedTopic = crypto.Keccak256Hash([]byte("ChallengeCleared(bytes32 indexed channelId, uint48 newTurnNumRecord)"))

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

// Since we only fetch events if there's a new block number
// it's safe to poll relatively frequently.
const EVENT_POLL_INTERVAL = 500 * time.Millisecond

// RESUB_INTERVAL is how often we resubscribe to log events.
// We do this to avoid https://github.com/ethereum/go-ethereum/issues/23845
// We use 2.5 minutes as the default filter timeout is 5 minutes.
// See https://github.com/ethereum/go-ethereum/blob/e14164d516600e9ac66f9060892e078f5c076229/eth/filters/filter_system.go#L43
const RESUB_INTERVAL = 2*time.Minute + 30*time.Second

// NewEthChainService constructs a chain service that submits transactions to a NitroAdjudicator
// and listens to events from an eventSource
func NewEthChainService(chain ethChain, na *NitroAdjudicator.NitroAdjudicator,
	naAddress, caAddress, vpaAddress common.Address, txSigner *bind.TransactOpts, logDestination io.Writer) (*EthChainService, error) {

	logging.ConfigureZeroLogger()

	logger := zerolog.New(logDestination).With().Timestamp().Str("txSigner", txSigner.From.String()[0:8]).Caller().Logger()
	ctx, cancelCtx := context.WithCancel(context.Background())

	// Use a buffered channel so we don't have to worry about blocking on writing to the channel.
	ecs := EthChainService{chain, na, naAddress, caAddress, vpaAddress, txSigner, make(chan Event, 10), logger, ctx, cancelCtx}

	if ecs.subscriptionsSupported() {
		logger.Printf("Notifications are supported by the chain. Using notifications to listen for events.")
		go ecs.subscribeForLogs()
	} else {
		logger.Printf("Notifications are NOT supported by the chain. Using polling to listen for events.")
		go ecs.pollForLogs()
	}
	return &ecs, nil
}

// subscriptionsSupported returns true if the node supports subscriptions for events.
// Otherwise returns false
func (ecs *EthChainService) subscriptionsSupported() bool {
	// This is slightly painful but seems like the only way to find out if notifications are supported
	// We attempt to subscribe (with a query that should never return a result) and check the error
	query := ethereum.FilterQuery{
		Addresses: []common.Address{ecs.naAddress},
		FromBlock: big.NewInt(0),
		ToBlock:   big.NewInt(0),
	}

	logs := make(chan ethTypes.Log, 1)
	sub, err := ecs.chain.SubscribeFilterLogs(context.Background(), query, logs)
	if err == rpc.ErrNotificationsUnsupported {
		return false
	}

	sub.Unsubscribe()
	return true
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

// fatalF is called to output a message and then panic, killing the chain service.
// It accepts a format string and arguments, as per fmt.Printf.
// If prints out the error to STDOUT, the logger and then exits the program.
func (ecs *EthChainService) fatalF(format string, v ...any) {

	// Print to STDOUT in case we're using a noop logger
	fmt.Println(fmt.Errorf(format, v...))

	ecs.logger.Fatal().Msgf(format, v...)

	// Manually panic in case we're using a logger that doesn't call exit(1)
	panic(fmt.Errorf(format, v...))

}

// dispatchChainEvents takes in a collection of event logs from the chain
// and dispatches events to the out channel
func (ecs *EthChainService) dispatchChainEvents(logs []ethTypes.Log) {
	for _, l := range logs {
		switch l.Topics[0] {
		case depositedTopic:
			nad, err := ecs.na.ParseDeposited(l)
			if err != nil {
				ecs.fatalF("error in ParseDeposited: %v", err)
			}

			event := NewDepositedEvent(nad.Destination, l.BlockNumber, nad.Asset, nad.AmountDeposited, nad.DestinationHoldings)
			ecs.out <- event
		case allocationUpdatedTopic:
			au, err := ecs.na.ParseAllocationUpdated(l)
			if err != nil {
				ecs.fatalF("error in ParseAllocationUpdated: %v", err)

			}

			tx, pending, err := ecs.chain.TransactionByHash(context.Background(), l.TxHash)
			if pending {
				ecs.fatalF("Expected transaction to be part of the chain, but the transaction is pending")

			}
			var assetAddress types.Address
			var amount *big.Int

			if err != nil {
				ecs.fatalF("error in TransactionByHash: %v", err)
			}

			assetAddress, amount, err = getChainHolding(ecs.na, tx, au)
			if err != nil {
				ecs.fatalF("error in getChainHoldings: %v", err)
			}

			event := NewAllocationUpdatedEvent(au.ChannelId, l.BlockNumber, assetAddress, amount)
			ecs.out <- event
		case concludedTopic:
			ce, err := ecs.na.ParseConcluded(l)
			if err != nil {
				ecs.fatalF("error in ParseConcluded: %v", err)

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

}

// getCurrentBlockNumber returns the current block number.
func (ecs *EthChainService) getCurrentBlockNum() *big.Int {
	h, err := ecs.chain.HeaderByNumber(context.Background(), nil)
	if err != nil {
		ecs.fatalF("headerByNumber failed: %w", err)
	}

	return h.Number
}

// subscribeForLogs subscribes for logs and pushes them to the out channel.
// It relies on notifications being supported by the chain node.
func (ecs *EthChainService) subscribeForLogs() {
	// Subscribe to Adjudicator events
	query := ethereum.FilterQuery{
		Addresses: []common.Address{ecs.naAddress},
	}
	logs := make(chan ethTypes.Log)
	sub, err := ecs.chain.SubscribeFilterLogs(ecs.ctx, query, logs)
	if err != nil {
		ecs.fatalF("subscribeFilterLogs failed: %w", err)
	}
	for {
		select {
		case <-ecs.ctx.Done():
			sub.Unsubscribe()
			return
		case err := <-sub.Err():
			if err != nil {
				ecs.fatalF("received error from the subscription channel: %w", err)
			}

			// If the error is nil then the subscription was closed and we need to re-subscribe.
			// This is a workaround for https://github.com/ethereum/go-ethereum/issues/23845
			var sErr error
			sub, sErr = ecs.chain.SubscribeFilterLogs(ecs.ctx, query, logs)
			if sErr != nil {
				ecs.fatalF("subscribeFilterLogs failed on resubscribe: %w", err)
			}
			ecs.logger.Print("resubscribed to filtered logs")

		case <-time.After(RESUB_INTERVAL):
			// Due to https://github.com/ethereum/go-ethereum/issues/23845 we can't rely on a long running subscription.
			// We unsub here and recreate the subscription in the next iteration of the select.
			sub.Unsubscribe()
		case chainEvent := <-logs:
			ecs.dispatchChainEvents([]ethTypes.Log{chainEvent})
		}
	}

}

type blockRange struct {
	from *big.Int
	to   *big.Int
}

// splitBlockRange takes a BlockRange and chunks it into a slice of BlockRanges, each having an interval no larger than the passed interval.
func splitBlockRange(total blockRange, maxInterval *big.Int) []blockRange {

	if total.from.Cmp(total.to) > 0 {
		panic(fmt.Sprintf("splitBlockRange: from > to. from = %v, to = %v", total.from, total.to))
	}

	slice := make([]blockRange, 0) // TODO precompute a capacity by dividing total interval by max interval

	start := big.NewInt(0).Set(total.from)
	for {
		finish := types.Min(total.to, big.NewInt(0).Add(start, maxInterval))
		slice = append(slice, blockRange{
			from: big.NewInt(0).Set(start),
			to:   big.NewInt(0).Set(finish),
		})
		if finish.Cmp(total.to) >= 0 {
			break
		}
		start = big.NewInt(0).Add(finish, big.NewInt(1))
	}

	return slice

}

// fetchLogsFromChain fetches logs from the chain from the given block number to the given block number.
// It splits the query into multiple queries if the range is too large.
func (ecs *EthChainService) fetchLogsFromChain(from *big.Int, to *big.Int) ([]ethTypes.Log, error) {

	logs := make([]ethTypes.Log, 0)

	blockRanges := splitBlockRange(blockRange{from, to}, big.NewInt(int64(MAX_QUERY_BLOCK_RANGE)))

	for _, bR := range blockRanges {

		query := ethereum.FilterQuery{
			Addresses: []common.Address{ecs.naAddress},
			FromBlock: bR.from,
			ToBlock:   bR.to,
		}

		fetchedLogs, err := ecs.chain.FilterLogs(context.Background(), query)

		if err != nil {
			return nil, err
		}

		logs = append(logs, fetchedLogs...)
	}
	return logs, nil

}

// pollForLogs periodically polls the chain client to check if there new events.
// It can function over a chain node that does not support notifications.
func (ecs *EthChainService) pollForLogs() {
	toBlock := ecs.getCurrentBlockNum()
	fetchedLogs, err := ecs.fetchLogsFromChain(big.NewInt(0), toBlock)
	if err != nil {
		ecs.fatalF("first fetchLogsFromChain failed: %w", err)
	}

	ecs.dispatchChainEvents(fetchedLogs)
	for {
		select {
		case <-ecs.ctx.Done():
			return
		case <-time.After(EVENT_POLL_INTERVAL):
			currentBlock := ecs.getCurrentBlockNum()

			if moreRecentBlockAvailable := currentBlock.Cmp(toBlock) > 0; moreRecentBlockAvailable {
				// The query includes the from and to blocks so we need to increment the from block to avoid duplicating events
				fromBlock := big.NewInt(0).Add(toBlock, big.NewInt(1))
				fetchedLogs, err = ecs.fetchLogsFromChain(fromBlock, currentBlock)
				toBlock.Set(currentBlock)

				if err != nil {
					ecs.fatalF("fetchLogsFromChain failed: %w", err)
				}

				ecs.dispatchChainEvents(fetchedLogs)
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

func (ecs *EthChainService) GetChainId() (*big.Int, error) {
	return ecs.chain.ChainID(context.Background())
}

func (ecs *EthChainService) Close() error {
	ecs.cancel()
	close(ecs.out)
	return nil
}
