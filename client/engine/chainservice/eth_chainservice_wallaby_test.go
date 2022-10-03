package chainservice

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/filecoin-project/go-address"
	NitroAdjudicator "github.com/statechannels/go-nitro/client/engine/chainservice/adjudicator"
)

func TestEthChainServiceAgainstWallaby(t *testing.T) {

	endpoint := "https://wallaby.node.glif.io/rpc/v0"

	client, err := ethclient.Dial(endpoint)
	if err != nil {
		t.Fatal(err)
	}

	bn, err := client.BlockNumber(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(bn)

	// Add a valid private key with testnet Eth. DO NOT check into git.
	pkString := "6645aa9129061ccef190e1bb1e11319b3d716b3140eec27595d045dbd565733b"

	// one := big.NewInt(1)

	pk, err := crypto.HexToECDSA(pkString)
	if err != nil {
		t.Fatal(err)
	}
	f1Address, err := address.NewSecp256k1Address(crypto.FromECDSAPub(&pk.PublicKey))

	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer([]byte(`{ "jsonrpc": "2.0", "method": "Filecoin.MpoolGetNonce","params": ["`+f1Address.String()+`"], "id":67}`)))

	if err != nil {
		t.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		t.Fatal(err)
	}

	type responseTy struct {
		jsonrpc string
		result  int64
		id      int64
	}

	responseBody := responseTy{}
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		t.Fatal(err)
	}

	nonce := responseBody.result

	txSubmitter, err := bind.NewKeyedTransactorWithChainID(pk, big.NewInt(31415))
	if err != nil {
		log.Fatal(err)
	}

	txSubmitter.GasLimit = uint64(47755863) // in units
	txSubmitter.GasFeeCap = big.NewInt(100)

	txSubmitter.Nonce = big.NewInt(nonce + 1)

	// As of the "Iron" FVM release, it seems that the return value of things like eth_getBlockByNumber do not match the spec.
	// Linked to this (probably) https://github.com/filecoin-project/ref-fvm/issues/908
	// Since the geth ethClient calls out to eth_getBlockNumber and tries to deserialize the result to one including `logsBloom` parameter, the following command will not yet work:
	// naAddress, _, na, err := NitroAdjudicator.DeployNitroAdjudicator(txSubmitter, client)

	signedTx, err := txSubmitter.Signer(txSubmitter.From,
		gethTypes.NewTx(&(gethTypes.DynamicFeeTx{
			// ChainID:   big.NewInt(31415),
			Nonce:     txSubmitter.Nonce.Uint64(),
			GasTipCap: txSubmitter.GasTipCap,
			GasFeeCap: txSubmitter.GasFeeCap,
			Gas:       txSubmitter.GasLimit,
			Value:     big.NewInt(0),
			Data:      []byte(NitroAdjudicator.NitroAdjudicatorMetaData.Bin),
		})))
	if err != nil {
		t.Fatal(err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		t.Fatal(err)
	}

	receipt, err := client.TransactionReceipt(context.Background(), signedTx.Hash())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(receipt)

	// caAddress := common.Address{}  // TODO use proper address
	// vpaAddress := common.Address{} // TODO use proper address

	// cs, err := NewEthChainService(client, na, naAddress, caAddress, vpaAddress, txSubmitter, NoopLogger{})
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// // Prepare test data to trigger EthChainService
	// testDeposit := types.Funds{
	// 	common.HexToAddress("0x00"): one,
	// }

	// channelID := types.Destination(common.HexToHash(`4ebd366d014a173765ba1e50f284c179ade31f20441bec41664712aac6cc461d`))
	// testTx := protocols.NewDepositTransaction(channelID, testDeposit)

	// out := cs.EventFeed()
	// // Submit transactiom
	// err = cs.SendTransaction(testTx)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// // Check that the recieved events matches the expected event
	// for i := 0; i < 2; i++ {
	// 	receivedEvent := <-out
	// 	dEvent := receivedEvent.(DepositedEvent)
	// 	expectedEvent := NewDepositedEvent(channelID, 2, dEvent.AssetAddress, testDeposit[dEvent.AssetAddress], testDeposit[dEvent.AssetAddress])
	// 	// TODO to validate BlockNum and NowHeld values, chain state prior to transaction must be inspected
	// 	ignoreBlockNum := cmpopts.IgnoreFields(commonEvent{}, "BlockNum")
	// 	ignoreNowHeld := cmpopts.IgnoreFields(DepositedEvent{}, "NowHeld")

	// 	if diff := cmp.Diff(expectedEvent, dEvent, cmp.AllowUnexported(DepositedEvent{}, commonEvent{}, big.Int{}), ignoreBlockNum, ignoreNowHeld); diff != "" {
	// 		t.Fatalf("Received event did not match expectation; (-want +got):\n%s", diff)
	// 	}
	// 	delete(testDeposit, dEvent.AssetAddress)
	// }

	// if len(testDeposit) != 0 {
	// 	t.Fatalf("Mismatch between the deposit transaction and the received events")
	// }
}
