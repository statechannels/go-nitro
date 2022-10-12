package chainservice

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"math/rand"
	"net/http"
	"time"

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
	f1Address                filecoinAddress.Address
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

	ecs := FevmChainService{chain, endpoint, f1Address, &na, naAddress, caAddress, vpaAddress, txSubmitter, make(chan Event, 10), logger}

	// err := fcs.subcribeToEvents() // TODO
	return &ecs, err
}

func (fcs *FevmChainService) rpcCall(method, params string, result interface{}) error {
	// requestBody := bytes.NewBuffer([]byte(`{ "jsonrpc": "2.0", "method": ` + method + `,"params": [` + params + `], "id":` + fmt.Sprint(rand.Intn(1000)) + `}`))
	requestBody := struct {
		Jsonrpc string   `json:"jsonrpc"`
		Method  string   `json:"method"`
		Params  []string `json:"params"`
		Id      int      `json:"id"`
	}{
		Jsonrpc: "2.0",
		Method:  method,
		Params:  []string{params},
		Id:      rand.Intn(1000),
	}
	encodedRequestBody, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}
	fmt.Println(string(encodedRequestBody))
	resp, err := http.Post(fcs.endpoint, "application/json", bytes.NewBuffer(encodedRequestBody))
	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)

	fmt.Println(string(body))

	if err != nil {
		return err
	}

	err = json.Unmarshal(body, result)
	if err != nil {
		return err
	}
	return nil
}

func (fcs *FevmChainService) filecoinNonce() (uint64, error) {
	result := struct {
		Result uint64 `json:"result"`
	}{Result: 999}
	err := fcs.rpcCall("Filecoin.MpoolGetNonce", fcs.f1Address.String(), &result)
	if err != nil {
		return 0, err
	}
	return result.Result, nil
}

func (fcs *FevmChainService) deployAdjudicator() error {

	nonce, err := fcs.filecoinNonce()
	if err != nil {
		return fmt.Errorf("could not get nonce %w", err)
	}

	hexBytecode, err := hex.DecodeString(NitroAdjudicator.NitroAdjudicatorMetaData.Bin[2:])
	if err != nil {
		log.Fatal(err)
	}
	// As of the "Iron" FVM release, it seems that the return value of things like eth_getBlockByNumber do not match the spec.
	// Linked to this (probably) https://github.com/filecoin-project/ref-fvm/issues/908
	// Since the geth ethClient calls out to eth_getBlockNumber and tries to deserialize the result to one including `logsBloom` parameter, the following command will not yet work:
	// naAddress, _, na, err := NitroAdjudicator.DeployNitroAdjudicator(txSubmitter, client)
	// https://ethereum.stackexchange.com/questions/107814/getting-current-base-fee-from-json-rpc
	signedTx, err := fcs.txSigner.Signer(
		fcs.txSigner.From,
		ethTypes.NewTx(&(ethTypes.DynamicFeeTx{
			// ChainID:   big.NewInt(31415),
			Nonce:     nonce,
			GasFeeCap: fcs.defaultTxOpts().GasFeeCap,
			Gas:       fcs.defaultTxOpts().GasLimit,
			Value:     big.NewInt(0),
			Data:      hexBytecode,
		})))
	if err != nil {
		return fmt.Errorf("could not sign tx %w", err)
	}

	fmt.Println("Sending tx")
	err = fcs.chain.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return fmt.Errorf("could not send tx %w", err)
	}
	fmt.Println("Tx sent")

	time.Sleep(10)
	isPending := true
	for isPending {
		tx, iP, err := fcs.chain.TransactionByHash(context.Background(), signedTx.Hash())
		if err != nil {
			return err
		}
		isPending = iP
		time.Sleep(10)
		fmt.Println(tx)
	}

	// naAddress, err := bind.WaitDeployed(context.Background(), fcs.chain, signedTx)
	// if err != nil {
	// 	return fmt.Errorf("could not wait for tx %w", err)
	// }
	fmt.Println("Deploy successful")
	// fcs.naAddress = naAddress
	// TODO populate na  on fcs
	return nil
}

// defaultTxOpts returns transaction options suitable for most transaction submissions
func (fcs *FevmChainService) defaultTxOpts() *bind.TransactOpts {
	return &bind.TransactOpts{
		From:      fcs.txSigner.From,
		Nonce:     fcs.txSigner.Nonce,
		Signer:    fcs.txSigner.Signer,
		GasFeeCap: big.NewInt(10_000_000_000),
		GasLimit:  1000000000, // BlockGasLimit / 10
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

func (fcs *FevmChainService) GetConsensusAppAddress() types.Address {
	return fcs.consensusAppAddress
}

func (fcs *FevmChainService) GetVirtualPaymentAppAddress() types.Address {
	return fcs.virtualPaymentAppAddress
}
