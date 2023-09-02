package chainservice

import (
	"container/heap"
	"context"
	"fmt"
	"log/slog"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/internal/chain"
	"github.com/statechannels/go-nitro/internal/logging"
	NitroAdjudicator "github.com/statechannels/go-nitro/node/engine/chainservice/adjudicator"
	Token "github.com/statechannels/go-nitro/node/engine/chainservice/erc20"
	chainutils "github.com/statechannels/go-nitro/node/engine/chainservice/utils"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

var (
	naAbi, _                 = NitroAdjudicator.NitroAdjudicatorMetaData.GetAbi()
	concludedTopic           = naAbi.Events["Concluded"].ID
	allocationUpdatedTopic   = naAbi.Events["AllocationUpdated"].ID
	depositedTopic           = naAbi.Events["Deposited"].ID
	challengeRegisteredTopic = naAbi.Events["ChallengeRegistered"].ID
	challengeClearedTopic    = naAbi.Events["ChallengeCleared"].ID
)

var topicsToWatch = []common.Hash{
	allocationUpdatedTopic,
	concludedTopic,
	depositedTopic,
	challengeRegisteredTopic,
	challengeClearedTopic,
}

type ethChain interface {
	bind.ContractBackend
	ethereum.TransactionReader
	ethereum.ChainReader
	ChainID(ctx context.Context) (*big.Int, error)
}

// eventTracker holds on to events in memory and dispatches an event after required number of confirmations
type eventTracker struct {
	latestBlockNum uint64
	events         EventQueue
	mu             sync.Mutex
}

type EthChainService struct {
	chain                    ethChain
	na                       *NitroAdjudicator.NitroAdjudicator
	naAddress                common.Address
	consensusAppAddress      common.Address
	virtualPaymentAppAddress common.Address
	txSigner                 *bind.TransactOpts
	out                      chan Event
	logger                   *slog.Logger
	ctx                      context.Context
	cancel                   context.CancelFunc
	wg                       *sync.WaitGroup
	eventTracker             *eventTracker
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
// This has been reduced to 15 seconds to support local devnets with much shorter timeouts.
const RESUB_INTERVAL = 15 * time.Second

// REQUIRED_BLOCK_CONFIRMATIONS is how many blocks must be mined before an emitted event is processed
const REQUIRED_BLOCK_CONFIRMATIONS = 2

// NewEthChainService is a convenient wrapper around newEthChainService, which provides a simpler API
func NewEthChainService(chainOpts chain.ChainOpts) (*EthChainService, error) {
	if chainOpts.ChainPk == "" {
		return nil, fmt.Errorf("chainpk must be set")
	}
	if chainOpts.VpaAddress == chainOpts.CaAddress {
		return nil, fmt.Errorf("virtual payment app address and consensus app address cannot be the same: %s", chainOpts.VpaAddress.String())
	}

	ethClient, txSigner, err := chainutils.ConnectToChain(
		context.Background(),
		chainOpts.ChainUrl,
		chainOpts.ChainAuthToken,
		common.Hex2Bytes(chainOpts.ChainPk),
	)
	if err != nil {
		panic(err)
	}

	na, err := NitroAdjudicator.NewNitroAdjudicator(chainOpts.NaAddress, ethClient)
	if err != nil {
		panic(err)
	}

	return newEthChainService(ethClient, chainOpts.ChainStartBlock, na, chainOpts.NaAddress, chainOpts.CaAddress, chainOpts.VpaAddress, txSigner)
}

// newEthChainService constructs a chain service that submits transactions to a NitroAdjudicator
// and listens to events from an eventSource
func newEthChainService(chain ethChain, startBlock uint64, na *NitroAdjudicator.NitroAdjudicator,
	naAddress, caAddress, vpaAddress common.Address, txSigner *bind.TransactOpts,
) (*EthChainService, error) {
	ctx, cancelCtx := context.WithCancel(context.Background())

	logger := logging.LoggerWithAddress(slog.Default(), txSigner.From)

	eventQueue := EventQueue{}
	heap.Init(&eventQueue)
	tracker := &eventTracker{latestBlockNum: 0, events: eventQueue}

	// Use a buffered channel so we don't have to worry about blocking on writing to the channel.
	ecs := EthChainService{chain, na, naAddress, caAddress, vpaAddress, txSigner, make(chan Event, 10), logger, ctx, cancelCtx, &sync.WaitGroup{}, tracker}
	errChan, newBlockSub, newBlockChan, eventSub, eventChan, eventQuery, err := ecs.subscribeForLogs()
	if err != nil {
		return nil, err
	}

	// Prevent go routines from processing events before checkForMissedEvents completes
	ecs.eventTracker.mu.Lock()

	ecs.wg.Add(3)
	go ecs.listenForEventLogs(errChan, eventSub, eventChan, eventQuery)
	go ecs.listenForNewBlocks(errChan, newBlockSub, newBlockChan)
	go ecs.listenForErrors(errChan)

	// Search for any missed events emitted while this node was offline
	err = ecs.checkForMissedEvents(startBlock)
	if err != nil {
		return nil, err
	}
	ecs.eventTracker.mu.Unlock()

	return &ecs, nil
}

func (ecs *EthChainService) checkForMissedEvents(startBlock uint64) error {
	ecs.logger.Info("checking for missed chainservice events", "startBlock", startBlock)
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(startBlock)),
		ToBlock:   nil, // For the current block number
		Addresses: []common.Address{ecs.naAddress},
		Topics:    [][]common.Hash{topicsToWatch},
	}

	events, err := ecs.chain.FilterLogs(ecs.ctx, query)
	if err != nil {
		ecs.logger.Error("failed to retrieve old chain logs: %v", err)
		return err
	}
	ecs.logger.Info("found missed events", "numEvents", len(events))

	for _, event := range events {
		heap.Push(&ecs.eventTracker.events, event)
	}
	return nil
}

// listenForErrors listens for errors on the error channel and attempts to handle them if they occur.
// TODO: Currently "handle" is panicking
func (ecs *EthChainService) listenForErrors(errChan <-chan error) {
	for {
		select {
		case <-ecs.ctx.Done():
			ecs.wg.Done()
			return
		case err := <-errChan:

			// Print to STDOUT in case we're using a noop logger
			fmt.Println(err)

			ecs.logger.Error("chain service error", "error", err)

			// Manually panic in case we're using a logger that doesn't call exit(1)
			panic(err)
		}
	}
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
				// TODO: wait for the Approve tx to be mined before continuing
			}
			holdings, err := ecs.na.Holdings(&bind.CallOpts{}, tokenAddress, tx.ChannelId())
			ecs.logger.Debug("existing holdings", "holdings", holdings)

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
		signedState := tx.SignedState.State()
		signatures := tx.SignedState.Signatures()
		nitroFixedPart := NitroAdjudicator.INitroTypesFixedPart(NitroAdjudicator.ConvertFixedPart(signedState.FixedPart()))
		nitroVariablePart := NitroAdjudicator.ConvertVariablePart(signedState.VariablePart())
		nitroSignatures := []NitroAdjudicator.INitroTypesSignature{NitroAdjudicator.ConvertSignature(signatures[0]), NitroAdjudicator.ConvertSignature(signatures[1])}

		candidate := NitroAdjudicator.INitroTypesSignedVariablePart{
			VariablePart: nitroVariablePart,
			Sigs:         nitroSignatures,
		}
		_, err := ecs.na.ConcludeAndTransferAllAssets(ecs.defaultTxOpts(), nitroFixedPart, candidate)
		return err
	case protocols.ChallengeTransaction:
		fp, candidate := NitroAdjudicator.ConvertSignedStateToFixedPartAndSignedVariablePart(tx.Candidate)
		proof := NitroAdjudicator.ConvertSignedStatesToProof(tx.Proof)
		challengerSig := NitroAdjudicator.ConvertSignature(tx.ChallengerSig)
		_, err := ecs.na.Challenge(ecs.defaultTxOpts(), fp, proof, candidate, challengerSig)
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
			ecs.logger.Debug("Processing Deposited event")
			nad, err := ecs.na.ParseDeposited(l)
			if err != nil {
				return fmt.Errorf("error in ParseDeposited: %w", err)
			}

			event := NewDepositedEvent(nad.Destination, l.BlockNumber, nad.Asset, nad.DestinationHoldings)
			ecs.out <- event

		case allocationUpdatedTopic:
			ecs.logger.Debug("Processing AllocationUpdated event")
			au, err := ecs.na.ParseAllocationUpdated(l)
			if err != nil {
				return fmt.Errorf("error in ParseAllocationUpdated: %w", err)
			}

			tx, pending, err := ecs.chain.TransactionByHash(context.Background(), l.TxHash)
			if pending {
				return fmt.Errorf("expected transaction to be part of the chain, but the transaction is pending")
			}
			if err != nil {
				return fmt.Errorf("error in TransactionByHash: %w", err)
			}

			assetAddress, err := assetAddressForIndex(ecs.na, tx, au.AssetIndex)
			if err != nil {
				return fmt.Errorf("error in assetAddressForIndex: %w", err)
			}
			ecs.logger.Debug("assetAddress", "assetAddress", assetAddress)

			event := NewAllocationUpdatedEvent(au.ChannelId, l.BlockNumber, assetAddress, au.FinalHoldings)
			ecs.out <- event

		case concludedTopic:
			ecs.logger.Debug("Processing Concluded event")
			ce, err := ecs.na.ParseConcluded(l)
			if err != nil {
				return fmt.Errorf("error in ParseConcluded: %w", err)
			}

			event := ConcludedEvent{commonEvent: commonEvent{channelID: ce.ChannelId, blockNum: l.BlockNumber}}
			ecs.out <- event

		case challengeRegisteredTopic:
			cr, err := ecs.na.ParseChallengeRegistered(l)
			if err != nil {
				return fmt.Errorf("error in ParseChallengeRegistered: %w", err)
			}
			event := NewChallengeRegisteredEvent(cr.ChannelId, l.BlockNumber, state.VariablePart{
				AppData: cr.Candidate.VariablePart.AppData,
				Outcome: NitroAdjudicator.ConvertBindingsExitToExit(cr.Candidate.VariablePart.Outcome),
				TurnNum: cr.Candidate.VariablePart.TurnNum.Uint64(),
				IsFinal: cr.Candidate.VariablePart.IsFinal,
			}, NitroAdjudicator.ConvertBindingsSignaturesToSignatures(cr.Candidate.Sigs))
			ecs.out <- event
		case challengeClearedTopic:
			ecs.logger.Info("Ignoring Challenge Cleared event")
		default:
			ecs.logger.Info("Ignoring unknown chain event topic", "topic", l.Topics[0].String())

		}
	}
	return nil
}

func (ecs *EthChainService) listenForEventLogs(errorChan chan<- error, eventSub ethereum.Subscription, eventChan chan ethTypes.Log, eventQuery ethereum.FilterQuery) {
out:
	for {
		select {
		case <-ecs.ctx.Done():
			eventSub.Unsubscribe()
			ecs.wg.Done()
			return

		case err := <-eventSub.Err():
			if err != nil {
				errorChan <- fmt.Errorf("received error from event subscription channel: %w", err)
				break out
			}

			// If the error is nil then the subscription was closed and we need to re-subscribe.
			// This is a workaround for https://github.com/ethereum/go-ethereum/issues/23845
			var sErr error
			eventSub, sErr = ecs.chain.SubscribeFilterLogs(ecs.ctx, eventQuery, eventChan)
			if sErr != nil {
				errorChan <- fmt.Errorf("subscribeFilterLogs failed on resubscribe: %w", err)
				break out
			}
			ecs.logger.Log(context.Background(), logging.LevelTrace, "resubscribed to filtered event logs")

		case <-time.After(RESUB_INTERVAL):
			// Due to https://github.com/ethereum/go-ethereum/issues/23845 we can't rely on a long running subscription.
			// We unsub here and recreate the subscription in the next iteration of the select.
			eventSub.Unsubscribe()

		case chainEvent := <-eventChan:
			ecs.logger.Debug("queueing new chainEvent", "block-num", chainEvent.BlockNumber)
			ecs.updateEventTracker(errorChan, nil, &chainEvent)
		}
	}
}

func (ecs *EthChainService) listenForNewBlocks(errorChan chan<- error, newBlockSub ethereum.Subscription, newBlockChan chan *ethTypes.Header) {
out:
	for {
		select {
		case <-ecs.ctx.Done():
			newBlockSub.Unsubscribe()
			ecs.wg.Done()
			return

		case err := <-newBlockSub.Err():
			if err != nil {
				errorChan <- fmt.Errorf("received error from new block subscription channel: %w", err)
				break out
			}

			// If the error is nil then the subscription was closed and we need to re-subscribe.
			// This is a workaround for https://github.com/ethereum/go-ethereum/issues/23845
			var sErr error
			newBlockSub, sErr = ecs.chain.SubscribeNewHead(ecs.ctx, newBlockChan)
			if sErr != nil {
				errorChan <- fmt.Errorf("subscribeNewHead failed on resubscribe: %w", err)
				break out
			}
			ecs.logger.Log(context.Background(), logging.LevelTrace, "resubscribed to new blocks")

		case newBlock := <-newBlockChan:
			newBlockNum := newBlock.Number.Uint64()
			ecs.logger.Log(context.Background(), logging.LevelTrace, "detected new block", "block-num", newBlockNum)
			ecs.updateEventTracker(errorChan, &newBlockNum, nil)
		}
	}
}

// updateEventTracker accepts a new block number and/or new event and dispatches a chain event if there are enough block confirmations
func (ecs *EthChainService) updateEventTracker(errorChan chan<- error, blockNumber *uint64, chainEvent *ethTypes.Log) {
	// lock the mutex for the shortest amount of time. The mutex only need to be locked to update the eventTracker data structure
	ecs.eventTracker.mu.Lock()

	if blockNumber != nil && *blockNumber > ecs.eventTracker.latestBlockNum {
		ecs.eventTracker.latestBlockNum = *blockNumber
	}
	if chainEvent != nil {
		heap.Push(&ecs.eventTracker.events, *chainEvent)
		ecs.logger.Debug("event added to queue", "updated-queue-length", ecs.eventTracker.events.Len())
	}

	eventsToDispatch := []ethTypes.Log{}
	for ecs.eventTracker.events.Len() > 0 && ecs.eventTracker.latestBlockNum >= (ecs.eventTracker.events)[0].BlockNumber+REQUIRED_BLOCK_CONFIRMATIONS {
		chainEvent := heap.Pop(&ecs.eventTracker.events).(ethTypes.Log)
		eventsToDispatch = append(eventsToDispatch, chainEvent)
		ecs.logger.Debug("event popped from queue", "updated-queue-length", ecs.eventTracker.events.Len())

	}
	ecs.eventTracker.mu.Unlock()

	err := ecs.dispatchChainEvents(eventsToDispatch)
	if err != nil {
		errorChan <- fmt.Errorf("failed dispatchChainEvents: %w", err)
		return
	}
}

// subscribeForLogs subscribes for logs and pushes them to the out channel.
// It relies on notifications being supported by the chain node.
func (ecs *EthChainService) subscribeForLogs() (chan error, ethereum.Subscription, chan *ethTypes.Header, ethereum.Subscription, chan ethTypes.Log, ethereum.FilterQuery, error) {
	// Subscribe to Adjudicator events
	eventQuery := ethereum.FilterQuery{
		Addresses: []common.Address{ecs.naAddress},
		Topics:    [][]common.Hash{topicsToWatch},
	}
	eventChan := make(chan ethTypes.Log)
	eventSub, err := ecs.chain.SubscribeFilterLogs(ecs.ctx, eventQuery, eventChan)
	if err != nil {
		return nil, nil, nil, nil, nil, ethereum.FilterQuery{}, fmt.Errorf("subscribeFilterLogs failed: %w", err)
	}
	errorChan := make(chan error)

	newBlockChan := make(chan *ethTypes.Header)
	newBlockSub, err := ecs.chain.SubscribeNewHead(ecs.ctx, newBlockChan)
	if err != nil {
		return nil, nil, nil, nil, nil, ethereum.FilterQuery{}, fmt.Errorf("subscribeNewHead failed: %w", err)
	}

	return errorChan, newBlockSub, newBlockChan, eventSub, eventChan, eventQuery, nil
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
	ecs.wg.Wait()
	return nil
}
