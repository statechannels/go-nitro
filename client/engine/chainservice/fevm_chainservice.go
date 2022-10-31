package chainservice

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
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
	NitroAdjudicator "github.com/statechannels/go-nitro/client/engine/chainservice/adjudicator"
	"github.com/statechannels/go-nitro/client/engine/store/safesync"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

const endpoint = "https://wallaby.node.glif.io/rpc/v0"
const chainId = 31415

func ethNonce(ethAddress common.Address) (int64, error) {
	type nonceResultTy struct {
		Result string `json:"result"`
	}
	var responseBody nonceResultTy
	err := rpcCall("eth_getTransactionCount", `["`+ethAddress.Hex()+`", "latest"]`, &responseBody)
	if err != nil {
		return 0, err
	}
	// TODO we need to catch errors by inspecting the responseBody for an error

	r, ok := big.NewInt(0).SetString(responseBody.Result[2:], 16)
	if !ok {
		return 0, errors.New("could not set bigint to string")
	}
	return r.Int64() + 1, nil

}

func latestBlockNum() (uint64, error) {
	responseBody := make(map[string]interface{})

	err := rpcCall("eth_blockNumber", `[]`, &responseBody)
	if err != nil {
		return 0, err
	}
	blockNumHex := responseBody["result"].(string)
	blockNum, success := new(big.Int).SetString(blockNumHex[2:], 16)
	if !success {
		log.Fatalf("Unable to convert block number %+v", blockNumHex)
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
	pk              *ecdsa.PrivateKey
}

func NewFevmChainService(pk *ecdsa.PrivateKey) ChainService {
	nitroAddress := common.HexToAddress("0xcf6E2F189DBDcfaC875097337957121060f38C2a")
	client, err := ethclient.Dial(endpoint)
	if err != nil {
		log.Fatal(err)
	}
	na, err := NitroAdjudicator.NewNitroAdjudicator(nitroAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	cs := &FevmChainService{na, nitroAddress, client, make(chan Event, 100), safesync.Map[watchDepositInfo]{}, pk}
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
			holdings, err := cs.na.Holdings(&bind.CallOpts{From: crypto.PubkeyToAddress(cs.pk.PublicKey)}, tokenAddress, tx.ChannelId())
			if err != nil {
				return err
			}

			txSubmitter, err := bind.NewKeyedTransactorWithChainID(cs.pk, big.NewInt(int64(chainId)))
			if err != nil {
				log.Fatal(err)
			}

			ethAddress := crypto.PubkeyToAddress(cs.pk.PublicKey)

			if err != nil {
				log.Fatalf("could not get address")
			}
			nonce, err := ethNonce(ethAddress)
			if err != nil {
				log.Fatal(err)
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
		return fmt.Errorf("withdraw transaction is not supported")

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
	cs.watchedChannels.Store(channelId.String(),
		watchDepositInfo{channelId: channelId, fundingTarget: total, ourFundingTarget: deposit, largestHeld: big.NewInt(0)})
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
				currentHoldings, err := cs.na.Holdings(&bind.CallOpts{From: crypto.PubkeyToAddress(cs.pk.PublicKey)}, info.asset, info.channelId)
				if err != nil {
					panic(err)
				}
				// Only send an event if the amount on chain has gone up.
				if currentHoldings.Cmp(info.largestHeld) > 0 {
					event := NewDepositedEvent(info.channelId, latestBlock, info.asset, info.ourFundingTarget, currentHoldings)
					cs.out <- event
					info.largestHeld.Set(currentHoldings)
				}
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
