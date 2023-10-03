package chainservice

import (
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
	"github.com/statechannels/go-nitro/internal/logging"
	NitroAdjudicator "github.com/statechannels/go-nitro/node/engine/chainservice/adjudicator"
	Token "github.com/statechannels/go-nitro/node/engine/chainservice/erc20"
	chainutils "github.com/statechannels/go-nitro/node/engine/chainservice/utils"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

type ChainOpts struct {
	ChainUrl        string
	ChainStartBlock uint64
	ChainAuthToken  string
	ChainPk         string
	NaAddress       common.Address
	VpaAddress      common.Address
	CaAddress       common.Address
}

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

const (
	MIN_BACKOFF_TIME = 1 * time.Second
	MAX_BACKOFF_TIME = 5 * time.Minute
)

type ethChain interface {
	bind.ContractBackend
	ethereum.TransactionReader
	ethereum.ChainReader
	ChainID(ctx context.Context) (*big.Int, error)
}

// eventTracker holds on to events in memory and dispatches an event after required number of confirmations
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
	eventSub                 ethereum.Subscription
	newBlockSub              ethereum.Subscription
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

// MAX_EPOCHS is the maximum range of old epochs we can query with a single "FilterLogs" request
// This is a restriction enforced by the rpc provider
const MAX_EPOCHS = 60480

// NewEthChainService is a convenient wrapper around newEthChainService, which provides a simpler API
func NewEthChainService(chainOpts ChainOpts) (ChainService, error) {
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
	tracker := NewEventTracker(startBlock)

	// Use a buffered channel so we don't have to worry about blocking on writing to the channel.
	ecs := EthChainService{chain, na, naAddress, caAddress, vpaAddress, txSigner, make(chan Event, 10), logger, ctx, cancelCtx, &sync.WaitGroup{}, tracker, nil, nil}
	errChan, newBlockChan, eventChan, eventQuery, err := ecs.subscribeForLogs()
	if err != nil {
		return nil, err
	}

	// Prevent go routines from processing events before checkForMissedEvents completes
	ecs.eventTracker.mu.Lock()
	defer ecs.eventTracker.mu.Unlock()

	ecs.wg.Add(3)
	go ecs.listenForEventLogs(errChan, eventChan, eventQuery)
	go ecs.listenForNewBlocks(errChan, newBlockChan)
	go ecs.listenForErrors(errChan)

	// Search for any missed events emitted while this node was offline
	err = ecs.checkForMissedEvents(startBlock)
	if err != nil {
		return nil, err
	}

	return &ecs, nil
}

func (ecs *EthChainService) checkForMissedEvents(startBlock uint64) error {
	// Fetch the latest block
	latestBlock, err := ecs.chain.BlockByNumber(ecs.ctx, nil)
	if err != nil {
		return err
	}

	latestBlockNum := latestBlock.NumberU64()
	ecs.logger.Info("checking for missed chain events", "startBlock", startBlock, "currentBlock", latestBlockNum)

	// Loop through in chunks of MAX_EPOCHS
	for currentStart := startBlock; currentStart <= latestBlockNum; {
		currentEnd := currentStart + MAX_EPOCHS
		if currentEnd > latestBlockNum {
			currentEnd = latestBlockNum
		}

		// Create a query for the current chunk
		query := ethereum.FilterQuery{
			FromBlock: big.NewInt(int64(currentStart)),
			ToBlock:   big.NewInt(int64(currentEnd)),
			Addresses: []common.Address{ecs.naAddress},
			Topics:    [][]common.Hash{topicsToWatch},
		}

		// Fetch logs for the current chunk
		missedEvents, err := ecs.chain.FilterLogs(ecs.ctx, query)
		if err != nil {
			ecs.logger.Error("failed to retrieve old chain logs. " + err.Error())
			errorMsg := "*** To avoid this error, consider increasing the chainstartblock value in your configuration before restarting the node."
			errorMsg += " Note that this may cause your node to miss chain events emitted prior to the chainstartblock."
			ecs.logger.Error(errorMsg)
			return err
		}
		ecs.logger.Info("finished checking for missed chain events in range", "fromBlock", currentStart, "toBlock", currentEnd, "numMissedEvents", len(missedEvents))

		for _, event := range missedEvents {
			ecs.eventTracker.Push(event)
		}

		currentStart = currentEnd + 1 // Move to the next chunk
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
			ecs.logger.Error("chain service error", "error", err)
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

			event := NewDepositedEvent(nad.Destination, l.BlockNumber, l.TxIndex, nad.Asset, nad.DestinationHoldings)
			ecs.out <- event

		case allocationUpdatedTopic:
			ecs.logger.Debug("Processing AllocationUpdated event")
			au, err := ecs.na.ParseAllocationUpdated(l)
			if err != nil {
				return fmt.Errorf("error in ParseAllocationUpdated: %w", err)
			}

			tx, pending, err := ecs.chain.TransactionByHash(ecs.ctx, l.TxHash)
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

			event := NewAllocationUpdatedEvent(au.ChannelId, l.BlockNumber, l.TxIndex, assetAddress, au.FinalHoldings)
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
			event := NewChallengeRegisteredEvent(cr.ChannelId, l.BlockNumber, l.TxIndex, state.VariablePart{
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

func (ecs *EthChainService) listenForEventLogs(errorChan chan<- error, eventChan chan ethTypes.Log, eventQuery ethereum.FilterQuery) {
	for {
		select {
		case <-ecs.ctx.Done():
			ecs.eventSub.Unsubscribe()
			ecs.wg.Done()
			return

		case err := <-ecs.eventSub.Err():
			ecs.eventTracker.mu.Lock()
			defer ecs.eventTracker.mu.Unlock()

			latestBlockNum := ecs.eventTracker.latestBlockNum

			if err != nil {
				ecs.logger.Warn("error in chain event subscription: " + err.Error())
				ecs.eventSub.Unsubscribe()
			} else {
				ecs.logger.Warn("chain event subscription closed")
			}

			// Use exponential backoff loop to attempt to re-establish subscription
			for backoffTime := MIN_BACKOFF_TIME; backoffTime < MAX_BACKOFF_TIME; backoffTime *= 2 {
				eventSub, err := ecs.chain.SubscribeFilterLogs(ecs.ctx, eventQuery, eventChan)
				if err != nil {
					ecs.logger.Warn("failed to resubscribe to chain events, retrying", "backoffTime", backoffTime)
					time.Sleep(backoffTime)
					continue
				}

				ecs.eventSub = eventSub
				ecs.logger.Debug("resubscribed to chain events")
				err = ecs.checkForMissedEvents(latestBlockNum)
				if err != nil {
					errorChan <- fmt.Errorf("subscribeFilterLogs failed during checkForMissedEvents: " + err.Error())
					return
				}

				break
			}

			ecs.logger.Error("subscribeFilterLogs failed to resubscribe")
			errorChan <- fmt.Errorf("subscribeFilterLogs failed to resubscribe")
			return

		case <-time.After(RESUB_INTERVAL):
			// Due to https://github.com/ethereum/go-ethereum/issues/23845 we can't rely on a long running subscription.
			// We unsub here and recreate the subscription in the next iteration of the select.
			ecs.eventSub.Unsubscribe()

		case chainEvent := <-eventChan:
			ecs.logger.Debug("queueing new chainEvent", "block-num", chainEvent.BlockNumber)
			ecs.updateEventTracker(errorChan, nil, &chainEvent)
		}
	}
}

func (ecs *EthChainService) listenForNewBlocks(errorChan chan<- error, newBlockChan chan *ethTypes.Header) {
	for {
		select {
		case <-ecs.ctx.Done():
			ecs.newBlockSub.Unsubscribe()
			ecs.wg.Done()
			return

		case err := <-ecs.newBlockSub.Err():
			if err != nil {
				ecs.logger.Warn("error in chain new block subscription: " + err.Error())
				ecs.newBlockSub.Unsubscribe()
			} else {
				ecs.logger.Warn("chain new block subscription closed")
			}

			// Use exponential backoff loop to attempt to re-establish subscription
			for backoffTime := MIN_BACKOFF_TIME; backoffTime < MAX_BACKOFF_TIME; backoffTime *= 2 {
				newBlockSub, err := ecs.chain.SubscribeNewHead(ecs.ctx, newBlockChan)
				if err != nil {
					errorChan <- fmt.Errorf("subscribeNewHead failed to resubscribe: %w", err)
					time.Sleep(backoffTime)
					continue
				}

				ecs.newBlockSub = newBlockSub
				ecs.logger.Debug("resubscribed to chain new blocks")
				break
			}

			ecs.logger.Error("subscribeNewHead failed to resubscribe")
			errorChan <- fmt.Errorf("subscribeNewHead failed to resubscribe")
			return

		case newBlock := <-newBlockChan:
			newBlockNum := newBlock.Number.Uint64()
			ecs.logger.Log(ecs.ctx, logging.LevelTrace, "detected new block", "block-num", newBlockNum)
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
		ecs.eventTracker.Push(*chainEvent)
		ecs.logger.Debug("event added to queue", "updated-queue-length", ecs.eventTracker.events.Len())
	}

	eventsToDispatch := []ethTypes.Log{}
	for ecs.eventTracker.events.Len() > 0 && ecs.eventTracker.latestBlockNum >= (ecs.eventTracker.events)[0].BlockNumber+REQUIRED_BLOCK_CONFIRMATIONS {
		chainEvent := ecs.eventTracker.Pop()
		ecs.logger.Debug("event popped from queue", "updated-queue-length", ecs.eventTracker.events.Len())

		// Ensure event & associated tx is still in the chain before adding to eventsToDispatch
		oldBlock, err := ecs.chain.BlockByNumber(context.Background(), new(big.Int).SetUint64(chainEvent.BlockNumber))
		if err != nil {
			ecs.logger.Error("failed to fetch block: %v", err)
			errorChan <- fmt.Errorf("failed to fetch block: %v", err)
			return
		}

		if oldBlock.Hash() != chainEvent.BlockHash {
			ecs.logger.Warn("dropping event because its block is no longer in the chain (possible re-org)", "blockNumber", chainEvent.BlockNumber, "blockHash", chainEvent.BlockHash)
			continue
		}

		eventsToDispatch = append(eventsToDispatch, chainEvent)
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
func (ecs *EthChainService) subscribeForLogs() (chan error, chan *ethTypes.Header, chan ethTypes.Log, ethereum.FilterQuery, error) {
	// Subscribe to Adjudicator events
	eventQuery := ethereum.FilterQuery{
		Addresses: []common.Address{ecs.naAddress},
		Topics:    [][]common.Hash{topicsToWatch},
	}
	eventChan := make(chan ethTypes.Log)
	eventSub, err := ecs.chain.SubscribeFilterLogs(ecs.ctx, eventQuery, eventChan)
	if err != nil {
		return nil, nil, nil, ethereum.FilterQuery{}, fmt.Errorf("subscribeFilterLogs failed: %w", err)
	}
	ecs.eventSub = eventSub
	errorChan := make(chan error)

	newBlockChan := make(chan *ethTypes.Header)
	newBlockSub, err := ecs.chain.SubscribeNewHead(ecs.ctx, newBlockChan)
	if err != nil {
		return nil, nil, nil, ethereum.FilterQuery{}, fmt.Errorf("subscribeNewHead failed: %w", err)
	}
	ecs.newBlockSub = newBlockSub

	return errorChan, newBlockChan, eventChan, eventQuery, nil
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
	return ecs.chain.ChainID(ecs.ctx)
}

func (ecs *EthChainService) GetLastConfirmedBlockNum() uint64 {
	var confirmedBlockNum uint64

	ecs.eventTracker.mu.Lock()
	defer ecs.eventTracker.mu.Unlock()

	// Check for potential underflow
	if ecs.eventTracker.latestBlockNum >= REQUIRED_BLOCK_CONFIRMATIONS {
		confirmedBlockNum = ecs.eventTracker.latestBlockNum - REQUIRED_BLOCK_CONFIRMATIONS
	} else {
		confirmedBlockNum = 0
	}

	return confirmedBlockNum
}

func (ecs *EthChainService) Close() error {
	ecs.cancel()
	ecs.wg.Wait()
	return nil
}
