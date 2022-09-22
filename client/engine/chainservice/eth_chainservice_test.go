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

// TestEthChainService is a sanity check that uses an Ethereum testnet as opposed to a geth Simulated Backend.
// As the test makes network requests, the run time is variable and is typically longer than what is desired from a unit test.
// For the test to pass, a valid private key with testnet ETH as well as as an Infura API key are needed.
func TestEthChainService(t *testing.T) {
	t.Skip()
	// Add a valid private key with testnet Eth. DO NOT check into git.
	pkString := ""
	// Add a valid Infura API key. DO NOT check into git.
	apiKey := ""

	one := big.NewInt(1)
	tokenAddress := common.HexToAddress("0xFc5eeC0FC4c97fe6b6BDEd926f5947308ef0d922")

	pk, err := crypto.HexToECDSA(pkString)
	if err != nil {
		t.Fatal(err)
	}

	client, err := ethclient.Dial("wss://goerli.infura.io/ws/v3/" + apiKey)
	if err != nil {
		t.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	txSubmitter, err := bind.NewKeyedTransactorWithChainID(pk, big.NewInt(5))
	if err != nil {
		log.Fatal(err)
	}
	txSubmitter.GasPrice = gasPrice
	txSubmitter.GasLimit = uint64(300000) // in units

	naAddress := common.HexToAddress("0x52dfe327D871A85f1AB0252A6ac67DBBFb7A2A2F")
	caAddress := common.Address{}  // TODO use proper address
	vpaAddress := common.Address{} // TODO use proper address

	na, err := NitroAdjudicator.NewNitroAdjudicator(naAddress, client)
	if err != nil {
		t.Fatal(err)
	}

	cs, err := NewEthChainService(client, na, naAddress, caAddress, vpaAddress, txSubmitter, NoopLogger{})
	if err != nil {
		t.Fatal(err)
	}

	// Prepare test data to trigger EthChainService
	testDeposit := types.Funds{
		common.HexToAddress("0x00"): one,
		tokenAddress:                one,
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
