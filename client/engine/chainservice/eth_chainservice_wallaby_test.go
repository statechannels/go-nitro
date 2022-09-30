package chainservice

import (
	"context"
	"log"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	NitroAdjudicator "github.com/statechannels/go-nitro/client/engine/chainservice/adjudicator"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

func TestEthChainServiceAgainstWallaby(t *testing.T) {
	// Add a valid private key with testnet Eth. DO NOT check into git.
	pkString := "6645aa9129061ccef190e1bb1e11319b3d716b3140eec27595d045dbd565733b" // or maybe 7b2254797065223a22736563703235366b31222c22507269766174654b6579223a225a6b57716b536b47484d37786b4f4737486845786d7a3178617a464137734a316c6442463239566c637a733d227d

	one := big.NewInt(1)

	pk, err := crypto.HexToECDSA(pkString)
	if err != nil {
		t.Fatal(err)
	}

	client, err := ethclient.Dial("https://wallaby.node.glif.io/rpc/v0")
	if err != nil {
		t.Fatal(err)
	}

	bn, err := client.BlockNumber(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(bn)

	// gasPrice, err := client.SuggestGasPrice(context.Background())
	// if err != nil {
	// 	log.Fatal(err)
	// }

	txSubmitter, err := bind.NewKeyedTransactorWithChainID(pk, big.NewInt(5))
	if err != nil {
		log.Fatal(err)
	}
	txSubmitter.GasPrice = big.NewInt(100) // gasPrice
	txSubmitter.GasLimit = uint64(300000)  // in units
	txSubmitter.GasFeeCap = big.NewInt(100)

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
