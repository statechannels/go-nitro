package chainservice

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	filecoinAddress "github.com/filecoin-project/go-address"
	NitroAdjudicator "github.com/statechannels/go-nitro/client/engine/chainservice/adjudicator"
	"github.com/statechannels/go-nitro/types"
)

func TestFevm(t *testing.T) {

	nitroAddress := common.HexToAddress("0xFF000000000000000000000000000000000003fA")
	channelIdString := "0xd9b535b686bcae01a00da8767de21d8bfc9915d513833160e5f15044fb4a3643"

	channelId := types.Destination(common.HexToHash(channelIdString))

	client, err := ethclient.Dial(endpoint)
	if err != nil {
		t.Fatal(err)
	}

	na, err := NitroAdjudicator.NewNitroAdjudicator(nitroAddress, client)
	if err != nil {
		t.Fatal(err)
	}

	holdings, err := na.Holdings(&bind.CallOpts{}, types.Address{}, channelId)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(holdings)

	pk, err := crypto.HexToECDSA(pkString)
	if err != nil {
		t.Fatal(err)
	}
	txSubmitter, err := bind.NewKeyedTransactorWithChainID(pk, big.NewInt(int64(chainId)))
	if err != nil {
		t.Fatal(err)
	}

	f1Address, err := filecoinAddress.NewSecp256k1Address(crypto.FromECDSAPub(&pk.PublicKey))
	if err != nil {
		t.Fatalf("could not get address")
	}
	nonce, err := fvmNonce(f1Address)
	if err != nil {
		t.Fatalf("could not get nonce")
	}

	abi, err := NitroAdjudicator.NitroAdjudicatorMetaData.GetAbi()
	if err != nil {
		t.Fatal(err)
	}
	data, err := abi.Pack("deposit", common.Address{}, channelId, big.NewInt(0), big.NewInt(1))
	if err != nil {
		t.Fatalf("unable to abi encode: %v", err)
	}

	tx := ethTypes.DynamicFeeTx{
		To:        &nitroAddress,
		Nonce:     uint64(nonce),
		Gas:       1000000000,
		GasTipCap: big.NewInt(200761),
		GasFeeCap: big.NewInt(200_000_000_000),
		Value:     big.NewInt(1),
		Data:      data,
	}

	signedTx, err := txSubmitter.Signer(
		txSubmitter.From,
		ethTypes.NewTx(&tx))
	if err != nil {
		t.Fatalf("could not sign tx %v", err)
	}
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		t.Fatalf("unable to submit tx %v", err)
	}

	// Waiting for tx to be mined times out
	// receipt, err := bind.WaitMined(context.Background(), client, signedTx)
	// if err != nil {
	// 	t.Fatalf("could not wait for tx %v", err)
	// }
	// fmt.Println(receipt.Status)
}
