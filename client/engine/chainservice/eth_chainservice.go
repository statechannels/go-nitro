package chainservice

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	NitroAdjudicator "github.com/statechannels/go-nitro/client/engine/chainservice/adjudicator"
	Token "github.com/statechannels/go-nitro/client/engine/chainservice/erc20"
	"github.com/statechannels/go-nitro/client/engine/store/safesync"
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
	chainServiceBase
	chain               ethChain
	na                  *NitroAdjudicator.NitroAdjudicator
	naAddress           common.Address
	consensusAppAddress common.Address
	txSigner            *bind.TransactOpts
}

// NewEthChainService constructs a chain service that submits transactions to a NitroAdjudicator
// and listens to events from an eventSource
func NewEthChainService(chain ethChain, na *NitroAdjudicator.NitroAdjudicator, naAddress common.Address, caAddress common.Address, txSigner *bind.TransactOpts) *EthChainService {
	ecs := EthChainService{chainServiceBase: newChainServiceBase()}
	ecs.out = safesync.Map[chan Event]{}
	ecs.chain = chain
	ecs.na = na
	ecs.naAddress = naAddress
	ecs.consensusAppAddress = caAddress
	ecs.txSigner = txSigner

	go ecs.listenForLogEvents()

	return &ecs
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
		ethTxs := []*ethTypes.Transaction{}
		for tokenAddress, amount := range tx.Deposit {
			txOpts := ecs.defaultTxOpts()
			ethTokenAddress := common.Address{}
			if tokenAddress == ethTokenAddress {
				txOpts.Value = amount
			} else {
				tokenTransactor, err := Token.NewTokenTransactor(tokenAddress, ecs.chain)
				if err != nil {
					panic(err)
				}
				_, err = tokenTransactor.Approve(ecs.defaultTxOpts(), ecs.naAddress, amount)
				if err != nil {
					panic(err)
				}
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
			SignedBy:     big.NewInt(0b11),
		}}
		_, err := ecs.na.ConcludeAndTransferAllAssets(ecs.defaultTxOpts(), nitroFixedPart, nitroSignedVariableParts)
		if err != nil {
			panic(err)
		}
		return nil

	default:
		return fmt.Errorf("unexpected transaction type %T", tx)
	}
}

func (ecs *EthChainService) listenForLogEvents() {
	query := ethereum.FilterQuery{
		Addresses: []common.Address{ecs.naAddress},
	}
	logs := make(chan ethTypes.Log)
	sub, err := ecs.chain.SubscribeFilterLogs(context.Background(), query, logs)
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
				nad, err := ecs.na.ParseDeposited(chainEvent)
				if err != nil {
					log.Fatal(err)
				}

				event := NewDepositedEvent(nad.Destination, chainEvent.BlockNumber, nad.Asset, nad.AmountDeposited, nad.DestinationHoldings)
				ecs.broadcast(event)
			case allocationUpdatedTopic:
				au, err := ecs.na.ParseAllocationUpdated(chainEvent)
				if err != nil {
					panic(err)
				}

				tx, pending, err := ecs.chain.TransactionByHash(context.Background(), chainEvent.TxHash)
				if pending || err != nil {
					panic("Expected transacion to be part of the chain")
				}

				assetAddress, amount, err := getChainHolding(ecs.na, tx, au)
				if err != nil {
					panic(err)
				}
				event := NewAllocationUpdatedEvent(au.ChannelId, chainEvent.BlockNumber, assetAddress, amount)
				ecs.broadcast(event)
			case concludedTopic:
				ce, err := ecs.na.ParseConcluded(chainEvent)
				if err != nil {
					panic(err)
				}

				event := ConcludedEvent{commonEvent: commonEvent{channelID: ce.ChannelId, BlockNum: chainEvent.BlockNumber}}
				ecs.broadcast(event)
			default:
				panic("Unknown chain event")
			}
		}
	}
}
