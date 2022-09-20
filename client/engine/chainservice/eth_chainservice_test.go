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
	t.Skip("This depends on a running chain instance with deployed contracts")
	// This is a funded key on the test hardhat network
	// See https://github.com/statechannels/hardhat-docker
	pkString := "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

	one := big.NewInt(1)
	// TODO: Deploy the token contract
	// tokenAddress := common.HexToAddress("0xFc5eeC0FC4c97fe6b6BDEd926f5947308ef0d922")

	pk, err := crypto.HexToECDSA(pkString)
	if err != nil {
		t.Fatal(err)
	}

	// This assumes a local chain instance is running (like hardhat)
	client, err := ethclient.Dial("ws://0.0.0.0:8545/")
	if err != nil {
		t.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	txSubmitter, err := bind.NewKeyedTransactorWithChainID(pk, big.NewInt(1337))
	if err != nil {
		log.Fatal(err)
	}
	txSubmitter.GasPrice = gasPrice
	txSubmitter.GasLimit = uint64(300000) // in units

	// This the address the adjudicator contract is deployed to by our docker hardhat instance
	// see https://github.com/statechannels/hardhat-docker
	naAddress := common.HexToAddress("0x5fbdb2315678afecb367f032d93f642f64180aa3")
	na, err := NitroAdjudicator.NewNitroAdjudicator(naAddress, client)
	if err != nil {
		t.Fatal(err)
	}

	cs, err := NewEthChainService(client, na, naAddress, common.Address{}, txSubmitter, NoopLogger{})
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
	// Submit transaction
	err = cs.SendTransaction(testTx)
	if err != nil {
		t.Fatal(err)
	}

	// Check that the recieved events matches the expected event
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

	if len(testDeposit) != 0 {
		t.Fatalf("Mismatch between the deposit transaction and the received events")
	}
}
