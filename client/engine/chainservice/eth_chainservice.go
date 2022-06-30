package chainservice

import (
	"context"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	NitroAdjudicator "github.com/statechannels/go-nitro/client/engine/chainservice/adjudicator"
	"github.com/statechannels/go-nitro/client/engine/store/safesync"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

var allocationUpdatedTopic = crypto.Keccak256Hash([]byte("AllocationUpdated(bytes32,uint256,uint256)"))
var concludedTopic = crypto.Keccak256Hash([]byte("Concluded(bytes32,uint48)"))
var depositedTopic = crypto.Keccak256Hash([]byte("Deposited(bytes32,address,uint256,uint256)"))

type ethChain interface {
	SubscribeFilterLogs(ctx context.Context, query ethereum.FilterQuery, ch chan<- ethTypes.Log) (ethereum.Subscription, error)
	TransactionByHash(ctx context.Context, txHash common.Hash) (*ethTypes.Transaction, bool, error)
}

type EthChainService struct {
	ChainServiceBase
	na       *NitroAdjudicator.NitroAdjudicator
	txSigner *bind.TransactOpts
}

// NewEthChainService constructs a chain service that submits transactions to a NitroAdjudicator
// and listens to events from an eventSource
func NewEthChainService(chain ethChain, na *NitroAdjudicator.NitroAdjudicator, naAddress common.Address, txSigner *bind.TransactOpts) *EthChainService {
	ecs := EthChainService{ChainServiceBase: newChainServiceBase()}
	ecs.out = safesync.Map[chan Event]{}
	ecs.na = na
	ecs.txSigner = txSigner

	go ecs.listenForLogEvents(na, naAddress, chain)

	return &ecs
}

// defaultTxOpts returns transaction options suitable for most transaction submissions
func (ecs *EthChainService) defaultTxOpts() *bind.TransactOpts {
	return &bind.TransactOpts{
		From:     ecs.txSigner.From,
		Nonce:    ecs.txSigner.Nonce,
		Signer:   ecs.txSigner.Signer,
		GasPrice: big.NewInt(10000000000),
	}
}

// SendTransaction sends the transaction and blocks until it has been submitted.
func (ecs *EthChainService) SendTransaction(tx protocols.ChainTransaction) []*ethTypes.Transaction {
	switch tx := tx.(type) {
	case protocols.DepositTransaction:
		ethTxs := []*ethTypes.Transaction{}
		for tokenAddress, amount := range tx.Deposit {
			txOpts := ecs.defaultTxOpts()
			ethTokenAddress := common.Address{}
			if tokenAddress == ethTokenAddress {
				txOpts.Value = amount
			}
			holdings, err := ecs.na.Holdings(&bind.CallOpts{}, tokenAddress, tx.ChannelId())
			if err != nil {
				panic(err)
			}

			ethTx, err := ecs.na.Deposit(txOpts, tokenAddress, tx.ChannelId(), holdings, amount)
			if err != nil {
				panic(err)
			}
			ethTxs = append(ethTxs, ethTx)
		}
		return ethTxs
	case protocols.WithdrawAllTransaction:
		state := tx.SignedState.State()
		signatures := tx.SignedState.Signatures()
		nitroFixedPart := NitroAdjudicator.IForceMoveFixedPart(state.FixedPart())
		nitroVariablePart := NitroAdjudicator.ConvertVariablePart(state.VariablePart())
		nitroSignatures := []NitroAdjudicator.IForceMoveSignature{NitroAdjudicator.ConvertSignature(signatures[0]), NitroAdjudicator.ConvertSignature(signatures[1])}

		ethTx, err := ecs.na.ConcludeAndTransferAllAssets(ecs.defaultTxOpts(), nitroFixedPart, nitroVariablePart, 1, []uint8{0, 0}, nitroSignatures)
		if err != nil {
			panic(err)
		}
		return []*ethTypes.Transaction{ethTx}

	default:
		panic("unexpected chain transaction")
	}
}

func (ecs *EthChainService) listenForLogEvents(na *NitroAdjudicator.NitroAdjudicator, naAddress common.Address, chain ethChain) {
	query := ethereum.FilterQuery{
		Addresses: []common.Address{naAddress},
	}
	logs := make(chan ethTypes.Log)
	sub, err := chain.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatal(err)
	}
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case chainEvent := <-logs:
			switch chainEvent.Topics[0] {
			case depositedTopic:
				nad, err := na.ParseDeposited(chainEvent)
				if err != nil {
					log.Fatal(err)
				}

				event := DepositedEvent{
					CommonEvent: CommonEvent{
						channelID: nad.Destination,
						BlockNum:  chainEvent.BlockNumber,
					},
					Asset:           nad.Asset,
					AmountDeposited: nad.AmountDeposited,
					NowHeld:         nad.DestinationHoldings,
				}
				ecs.broadcast(event)
			case allocationUpdatedTopic:
				au, err := na.ParseAllocationUpdated(chainEvent)
				if err != nil {
					panic(err)
				}

				tx, pending, err := chain.TransactionByHash(context.Background(), chainEvent.TxHash)
				if pending || err != nil {
					panic("Expected transacion to be part of the chain")
				}

				assetAddress, amount, err := getChainHolding(na, tx, au)
				if err != nil {
					panic(err)
				}
				event := AllocationUpdatedEvent{CommonEvent: CommonEvent{channelID: au.ChannelId, BlockNum: chainEvent.BlockNumber}, AssetAddress: assetAddress, AssetAmount: amount}
				ecs.broadcast(event)
			case concludedTopic:
				ce, err := na.ParseConcluded(chainEvent)
				if err != nil {
					panic(err)
				}

				event := ConcludedEvent{CommonEvent: CommonEvent{channelID: ce.ChannelId, BlockNum: chainEvent.BlockNumber}}
				ecs.broadcast(event)
			default:
				panic("Unknown chain event")
			}
		}
	}
}
