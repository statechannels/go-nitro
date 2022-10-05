package chainservice

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/big"
	"math/rand"
	"net/http"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	filecoinAddress "github.com/filecoin-project/go-address"
	NitroAdjudicator "github.com/statechannels/go-nitro/client/engine/chainservice/adjudicator"
	Token "github.com/statechannels/go-nitro/client/engine/chainservice/erc20"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

type FevmChainService struct {
	chain                    *ethclient.Client
	endpoint                 string
	na                       *NitroAdjudicator.NitroAdjudicator
	naAddress                common.Address
	consensusAppAddress      common.Address
	virtualPaymentAppAddress common.Address
	txSigner                 *bind.TransactOpts
	out                      chan Event
	logger                   *log.Logger
}

// NewFevmChainService constructs a chain service that points at a FEVM (Filecoin Ethereum Virtual Machine) endpoint. It deploys and submits transactions to a NitroAdjudicator
// and listens to events from an eventSource
func NewFevmChainService(endpoint string, pkString string, logDestination io.Writer) (*FevmChainService, error) {

	// hardcoded chain id
	chainId := big.NewInt(31415)

	chain, err := ethclient.Dial(endpoint)
	if err != nil {
		return &FevmChainService{}, err
	}

	pk, err := crypto.HexToECDSA(pkString)

	if err != nil {
		return &FevmChainService{}, err
	}

	address := crypto.PubkeyToAddress(pk.PublicKey)
	logPrefix := "chainservice " + address.String() + ": "
	logger := log.New(logDestination, logPrefix, log.Lmicroseconds|log.Lshortfile)

	f1Address, err := filecoinAddress.NewSecp256k1Address(crypto.FromECDSAPub(&pk.PublicKey))
	if err != nil {
		return &FevmChainService{}, err
	}
	logger.Println("filecoin address is ", f1Address)

	txSubmitter, err := bind.NewKeyedTransactorWithChainID(pk, chainId)

	// TODO these are stubbed but will be deployed in future
	na := NitroAdjudicator.NitroAdjudicator{}
	naAddress := types.Address{}
	caAddress := types.Address{}
	vpaAddress := types.Address{}

	ecs := FevmChainService{chain, endpoint, &na, naAddress, caAddress, vpaAddress, txSubmitter, make(chan Event, 10), logger}

	// err := fcs.subcribeToEvents() // TODO
	return &ecs, err
}

func (fcs *FevmChainService) rpcCall(method, params string, result interface{}) error {
	resp, err := http.Post(fcs.endpoint, "application/json", bytes.NewBuffer([]byte(`{ "jsonrpc": "2.0", "method": `+method+`,"params": [`+params+`], "id":`+fmt.Sprint(rand.Intn(1000))+`}`)))
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	type responseTy2 struct {
		Result string `json:"result"`
	}
	err = json.Unmarshal(body, &result)
	return nil
}

func (fcs *FevmChainService) filecoinNonce() (int64, error) {
	result := struct {
		Result int64 `json:"result"`
	}{}
	err := fcs.rpcCall("Filecoin.MpoolGetNonce", "", result)
	if err != nil {
		return 0, err
	}
	return result.Result, nil
}

// defaultTxOpts returns transaction options suitable for most transaction submissions
func (fcs *FevmChainService) defaultTxOpts() *bind.TransactOpts {
	return &bind.TransactOpts{
		From:     fcs.txSigner.From,
		Nonce:    fcs.txSigner.Nonce,
		Signer:   fcs.txSigner.Signer,
		GasPrice: fcs.txSigner.GasPrice,
		GasLimit: fcs.txSigner.GasLimit,
	}
}

// SendTransaction sends the transaction and blocks until it has been submitted.
func (fcs *FevmChainService) SendTransaction(tx protocols.ChainTransaction) error {
	switch tx := tx.(type) {
	case protocols.DepositTransaction:
		for tokenAddress, amount := range tx.Deposit {
			txOpts := fcs.defaultTxOpts()
			ethTokenAddress := common.Address{}
			if tokenAddress == ethTokenAddress {
				txOpts.Value = amount
			} else {
				tokenTransactor, err := Token.NewTokenTransactor(tokenAddress, fcs.chain)
				if err != nil {
					return err
				}
				_, err = tokenTransactor.Approve(fcs.defaultTxOpts(), fcs.naAddress, amount)
				if err != nil {
					return err
				}
			}
			holdings, err := fcs.na.Holdings(&bind.CallOpts{}, tokenAddress, tx.ChannelId())
			if err != nil {
				return err
			}

			_, err = fcs.na.Deposit(txOpts, tokenAddress, tx.ChannelId(), holdings, amount)
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
		_, err := fcs.na.ConcludeAndTransferAllAssets(fcs.defaultTxOpts(), nitroFixedPart, proof, candidate)
		return err

	default:
		return fmt.Errorf("unexpected transaction type %T", tx)
	}
}

func (fcs *FevmChainService) subcribeToEvents() error {
	// Subsribe to Adjudicator events
	query := ethereum.FilterQuery{
		Addresses: []common.Address{fcs.naAddress},
	}
	logs := make(chan ethTypes.Log)
	sub, err := fcs.chain.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		return err
	}
	go fcs.listenForLogEvents(sub, logs)
	return nil
}

func (fcs *FevmChainService) listenForLogEvents(sub ethereum.Subscription, logs chan ethTypes.Log) {
	for {
		select {
		case err := <-sub.Err():
			// TODO should we try resubscribing to chain events
			fcs.logger.Printf("event subscription error: %v", err)
		case chainEvent := <-logs:
			switch chainEvent.Topics[0] {
			case depositedTopic:
				nad, err := fcs.na.ParseDeposited(chainEvent)
				if err != nil {
					fcs.logger.Printf("error in ParseDeposited: %v", err)
				}

				event := NewDepositedEvent(nad.Destination, chainEvent.BlockNumber, nad.Asset, nad.AmountDeposited, nad.DestinationHoldings)
				fcs.out <- event
			case allocationUpdatedTopic:
				au, err := fcs.na.ParseAllocationUpdated(chainEvent)
				if err != nil {
					fcs.logger.Printf("error in ParseAllocationUpdated: %v", err)
				}

				tx, pending, err := fcs.chain.TransactionByHash(context.Background(), chainEvent.TxHash)
				if pending {
					fcs.logger.Printf("Expected transacion to be part of the chain, but the transaction is pending")
				}
				if err != nil {
					fcs.logger.Printf("error in TransactoinByHash: %v", err)
				}

				assetAddress, amount, err := getChainHolding(fcs.na, tx, au)
				if err != nil {
					fcs.logger.Printf("error in getChainHoldings: %v", err)
				}
				event := NewAllocationUpdatedEvent(au.ChannelId, chainEvent.BlockNumber, assetAddress, amount)
				fcs.out <- event
			case concludedTopic:
				ce, err := fcs.na.ParseConcluded(chainEvent)
				if err != nil {
					fcs.logger.Printf("error in ParseConcluded: %v", err)
				}

				event := ConcludedEvent{commonEvent: commonEvent{channelID: ce.ChannelId, BlockNum: chainEvent.BlockNumber}}
				fcs.out <- event
			default:
				fcs.logger.Printf("Unknown chain event")
			}
		}
	}
}

// EventFeed returns the out chan, and narrows the type so that external consumers may only receive on it.
func (fcs *FevmChainService) EventFeed() <-chan Event {
	return fcs.out
}

func (fcs *FevmChainService) GetConsensusAppAddress() types.Address {
	return fcs.consensusAppAddress
}

func (fcs *FevmChainService) GetVirtualPaymentAppAddress() types.Address {
	return fcs.virtualPaymentAppAddress
}
