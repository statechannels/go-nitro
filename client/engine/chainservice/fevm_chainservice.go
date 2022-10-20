package chainservice

import (
	"bytes"
	"context"
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
	"github.com/statechannels/go-nitro/client/engine/store/safesync"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

const endpoint = "https://wallaby.node.glif.io/rpc/v0"
const pkString = "716b7161580785bc96a4344eb52d23131aea0caf42a52dcf9f8aee9eef9dc3cd"
const chainId = 31415

func fvmNonce(f1Address filecoinAddress.Address) (int64, error) {
	type nonceResultTy struct {
		Result float64 `json:"result"`
	}
	var responseBody nonceResultTy
	err := rpcCall("Filecoin.MpoolGetNonce", `["`+f1Address.String()+`"]`, &responseBody)
	if err != nil {
		return 0, err
	}
	return int64(responseBody.Result), nil
}

func latestBlockNum() (uint64, error) {
	type blockResultTy struct {
		Number string `json:"result"`
	}
	var responseBody blockResultTy
	err := rpcCall("eth_getBlockByNumber", `["latest", "false"]`, &responseBody)
	if err != nil {
		return 0, err
	}
	blockNum, success := new(big.Int).SetString(responseBody.Number, 16)
	if !success {
		log.Fatalf("Unable to convert block number %+v", responseBody)
	}
	return blockNum.Uint64(), nil
}

func rpcCall(method, params string, result interface{}) error {
	reqBody := `{"jsonrpc": "2.0", "method": "` + method + `","params":` + params + `, "id":` + fmt.Sprint(rand.Intn(1000)) + `}`

	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer([]byte(reqBody)))
	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}
	return nil
}

type FevmChainService struct {
	na              *NitroAdjudicator.NitroAdjudicator
	naAddress       types.Address
	ethClient       *ethclient.Client
	out             chan Event
	watchedChannels safesync.Map[watchDepositInfo]
}

func NewFevmChainService() ChainService {
	nitroAddress := common.HexToAddress("0xFF000000000000000000000000000000000003fA")
	client, err := ethclient.Dial(endpoint)
	if err != nil {
		log.Fatal(err)
	}
	na, err := NitroAdjudicator.NewNitroAdjudicator(nitroAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	cs := &FevmChainService{na, nitroAddress, client, make(chan Event, 100), safesync.Map[watchDepositInfo]{}}
	go cs.pollChain(context.Background())
	return cs
}

func (cs *FevmChainService) EventFeed() <-chan Event {
	return cs.out
}

func (cs *FevmChainService) SendTransaction(tx protocols.ChainTransaction) error {
	switch tx := tx.(type) {
	case protocols.DepositTransaction:
		for tokenAddress, amount := range tx.Deposit {
			ethTokenAddress := common.Address{}
			if tokenAddress != ethTokenAddress {
				return fmt.Errorf("erc20 tokens are not supported")
			}
			holdings, err := cs.na.Holdings(&bind.CallOpts{}, tokenAddress, tx.ChannelId())
			if err != nil {
				return err
			}

			pk, err := crypto.HexToECDSA(pkString)
			if err != nil {
				log.Fatal(err)
			}
			txSubmitter, err := bind.NewKeyedTransactorWithChainID(pk, big.NewInt(int64(chainId)))
			if err != nil {
				log.Fatal(err)
			}

			f1Address, err := filecoinAddress.NewSecp256k1Address(crypto.FromECDSAPub(&pk.PublicKey))
			if err != nil {
				log.Fatalf("could not get address")
			}
			nonce, err := fvmNonce(f1Address)
			if err != nil {
				log.Fatalf("could not get nonce")
			}

			abi, err := NitroAdjudicator.NitroAdjudicatorMetaData.GetAbi()
			if err != nil {
				log.Fatal(err)
			}
			data, err := abi.Pack("deposit", common.Address{}, tx.ChannelId(), holdings, amount)
			if err != nil {
				log.Fatalf("unable to abi encode: %v", err)
			}

			chainTx := ethTypes.DynamicFeeTx{
				To:        &cs.naAddress,
				Nonce:     uint64(nonce),
				Gas:       1000000000,
				GasTipCap: big.NewInt(200761),
				GasFeeCap: big.NewInt(200_000_000_000),
				Value:     amount,
				Data:      data,
			}

			signedTx, err := txSubmitter.Signer(
				txSubmitter.From,
				ethTypes.NewTx(&chainTx))
			if err != nil {
				log.Fatal("could not sign tx %w", err)
			}
			err = cs.ethClient.SendTransaction(context.Background(), signedTx)
			if err != nil {
				log.Fatal("unable to submit tx %w", err)
			}
		}
		return nil
	case protocols.WithdrawAllTransaction:
		return fmt.Errorf("Withdraw transaction is not supported")

	default:
		return fmt.Errorf("unexpected transaction type %T", tx)
	}
}

func (cs *FevmChainService) GetConsensusAppAddress() types.Address {
	return common.Address{}
}

func (cs *FevmChainService) GetVirtualPaymentAppAddress() types.Address {
	return common.Address{}
}

func (cs *FevmChainService) Monitor(channelId types.Destination, ourDeposit, expectedTotal types.Funds) {
	deposit := ourDeposit[common.Address{}]
	total := expectedTotal[common.Address{}]
	cs.watchedChannels.Store(channelId.String(), watchDepositInfo{channelId: channelId, fundingTarget: total, ourFundingTarget: deposit})
}

// pollChain periodically polls the chain for holdings changes.
func (cs *FevmChainService) pollChain(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(POLL_INTERVAL):
			latestBlock, err := latestBlockNum()
			if err != nil {
				panic(err)
			}

			completed := make([]string, 0)
			// Range over all open deposit infos and check if the holdings have been updated.
			cs.watchedChannels.Range(func(key string, info watchDepositInfo) bool {
				currentHoldings, err := cs.na.Holdings(&bind.CallOpts{}, info.asset, info.channelId)
				if err != nil {
					panic(err)
				}
				event := NewDepositedEvent(info.channelId, latestBlock, info.asset, info.ourFundingTarget, currentHoldings)
				log.Printf("Sending deposited event. latest block %v, currentHoldings %v", latestBlock, currentHoldings)
				cs.out <- event
				// We only want to remove the channel if the deposit is fully complete.
				if currentHoldings.Cmp(info.fundingTarget) >= 0 {
					completed = append(completed, key)
				}

				return true
			})

			// Remove all completed tx infos.
			for _, key := range completed {
				cs.watchedChannels.Delete(key)
			}
		}
	}
}
