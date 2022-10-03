package chainservice

import (
	"bytes"
	"context"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/filecoin-project/go-address"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	NitroAdjudicator "github.com/statechannels/go-nitro/client/engine/chainservice/adjudicator"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
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

	one := big.NewInt(1)

	pk, err := crypto.HexToECDSA(pkString)
	if err != nil {
		t.Fatal(err)
	}
	pubKey := crypto.FromECDSA(pk)
	secp256k1Address, err := address.NewSecp256k1Address(pubKey)
	t.Log(secp256k1Address)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer([]byte(`{ "jsonrpc": "2.0", "method": "eth_blockNumber","params": [], "id":67}`)))

	// 	curl --location --request POST 'https://wallaby.node.glif.io/rpc/v0' \
	//   --header 'Content-Type: application/json' \
	//   --data-raw '{
	//   "jsonrpc":"2.0",
	//   "method":"eth_blockNumber",
	//   "params":[],
	//   "id":67
	//   }'

	// resp, err := http.Post(endpoint, "application/json", strings.NewReader(string(data)))

	if err != nil {
		t.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(body))

	// gasPrice, err := client.SuggestGasPrice(context.Background())
	// if err != nil {
	// 	log.Fatal(err)
	// }

	txSubmitter, err := bind.NewKeyedTransactorWithChainID(pk, big.NewInt(5))
	if err != nil {
		log.Fatal(err)
	}

	txSubmitter.GasLimit = uint64(300000) // in units
	txSubmitter.GasFeeCap = big.NewInt(100)

	txSubmitter.Nonce = big.NewInt(1)

	// nonce, err := client.NonceAt(context.Background(), txSubmitter.From, nil)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// txSubmitter.Nonce = big.NewInt(int64(nonce))

	// getting an error on the below: failed to retrieve account nonce: actor not found
	// possible solution:
	//   const f1addr = fa.newSecp256k1Address(pubKey).toString();
	//   const priorityFee = await callRpc("eth_maxPriorityFeePerGas");
	//   const nonce = await callRpc("Filecoin.MpoolGetNonce", [f1addr]); we need to call filecoin mpoolgetnonce with the filecoin address

	naAddress, _, na, err := NitroAdjudicator.DeployNitroAdjudicator(txSubmitter, client)

	if err != nil {
		t.Fatal(err)
	}

	caAddress := common.Address{}  // TODO use proper address
	vpaAddress := common.Address{} // TODO use proper address

	cs, err := NewEthChainService(client, na, naAddress, caAddress, vpaAddress, txSubmitter, NoopLogger{})
	if err != nil {
		t.Fatal(err)
	}

	// Prepare test data to trigger EthChainService
	testDeposit := types.Funds{
		common.HexToAddress("0x00"): one,
	}

	channelID := types.Destination(common.HexToHash(`4ebd366d014a173765ba1e50f284c179ade31f20441bec41664712aac6cc461d`))
	testTx := protocols.NewDepositTransaction(channelID, testDeposit)

	out := cs.EventFeed()
	// Submit transactiom
	err = cs.SendTransaction(testTx)
	if err != nil {
		t.Fatal(err)
	}

	// Check that the recieved events matches the expected event
	for i := 0; i < 2; i++ {
		receivedEvent := <-out
		dEvent := receivedEvent.(DepositedEvent)
		expectedEvent := NewDepositedEvent(channelID, 2, dEvent.AssetAddress, testDeposit[dEvent.AssetAddress], testDeposit[dEvent.AssetAddress])
		// TODO to validate BlockNum and NowHeld values, chain state prior to transaction must be inspected
		ignoreBlockNum := cmpopts.IgnoreFields(commonEvent{}, "BlockNum")
		ignoreNowHeld := cmpopts.IgnoreFields(DepositedEvent{}, "NowHeld")

		if diff := cmp.Diff(expectedEvent, dEvent, cmp.AllowUnexported(DepositedEvent{}, commonEvent{}, big.Int{}), ignoreBlockNum, ignoreNowHeld); diff != "" {
			t.Fatalf("Received event did not match expectation; (-want +got):\n%s", diff)
		}
		delete(testDeposit, dEvent.AssetAddress)
	}

	if len(testDeposit) != 0 {
		t.Fatalf("Mismatch between the deposit transaction and the received events")
	}
}
