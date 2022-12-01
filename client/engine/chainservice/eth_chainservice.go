package chainservice

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/big"
	"strings"
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

const FEVM_PARSE_ERROR = "json: cannot unmarshal hex string \"0x\" into Go struct field txJSON.v of type *hexutil.Big"

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

// fatalError is called when a goroutine encounters an unrecoverable error.
// If prints out the error to STDOUT, the logger and then exits the program.
func (ecs *EthChainService) fatalError(format string, v ...any) {

	// Print to STDOUT in case we're using a noop logger
	fmt.Println(fmt.Errorf(format, v...))
	// FatalF prints to the logger then calls exit(1)
	ecs.logger.Fatalf(format, v...)

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
				ecs.fatalError("error in ParseDeposited: %v", err)
			}

			event := NewDepositedEvent(nad.Destination, l.BlockNumber, nad.Asset, nad.AmountDeposited, nad.DestinationHoldings)
			fmt.Printf("DISPATCHING EVENT\n%+v\n", event)
			ecs.out <- event
		case allocationUpdatedTopic:
			au, err := ecs.na.ParseAllocationUpdated(l)
			if err != nil {
				ecs.fatalError("error in ParseAllocationUpdated: %v", err)

			}

			tx, pending, err := ecs.chain.TransactionByHash(context.Background(), l.TxHash)
			if pending {
				ecs.fatalError("Expected transaction to be part of the chain, but the transaction is pending")

			}
			var assetAddress types.Address
			var amount *big.Int
			switch {

			// TODO: Workaround for https://github.com/filecoin-project/ref-fvm/issues/1158
			case err != nil && strings.Contains(err.Error(), FEVM_PARSE_ERROR):
				ecs.logger.Printf("WARNING: Cannot parse transaction, assuming asset address of 0x00...")
				fmt.Printf("WARNING: Cannot parse transaction, assuming asset address of 0x00...")
				assetAddress := types.Address{}
				amount, err = getAssetHoldings(ecs.na, assetAddress, new(big.Int).SetUint64(au.Raw.BlockNumber), au.ChannelId)
				if err != nil {
					ecs.fatalError("error in getAssetHoldings: %v", err)
				}

			case err != nil:
				ecs.fatalError("error in TransactionByHash: %v", err)

			default:
				assetAddress, amount, err = getChainHolding(ecs.na, tx, au)
				if err != nil {
					ecs.fatalError("error in getChainHoldings: %v", err)
				}
			}

			event := NewAllocationUpdatedEvent(au.ChannelId, l.BlockNumber, assetAddress, amount)
			fmt.Printf("DISPATCHING EVENT\n%+v\n", event)
			ecs.out <- event
		case concludedTopic:
			ce, err := ecs.na.ParseConcluded(l)
			if err != nil {
				ecs.fatalError("error in ParseConcluded: %v", err)

			}

			event := ConcludedEvent{commonEvent: commonEvent{channelID: ce.ChannelId, BlockNum: l.BlockNumber}}
			fmt.Printf("DISPATCHING EVENT\n%+v\n", event)
			ecs.out <- event

		default:
			ecs.fatalError("Unknown chain event")
		}
	}

}

// getCurrentBlockNumber returns the current block number.
func (ecs *EthChainService) getCurrentBlockNum() *big.Int {
	h, err := ecs.chain.HeaderByNumber(context.Background(), nil)
	if err != nil {
		panic(err)
	}

	return h.Number
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
				// The query includes the from and to blocks so we need to increment the from block to avoid duplicating events
				query.FromBlock = big.NewInt(0).Add(query.ToBlock, big.NewInt(1))
				query.ToBlock = big.NewInt(0).Set(currentBlock)

				fetchedLogs, err := ecs.chain.FilterLogs(context.Background(), query)
				fmt.Printf("Polling from %d to %d found %d logs\n", query.FromBlock, query.ToBlock, len(fetchedLogs))
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
