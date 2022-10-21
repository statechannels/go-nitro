package chainservice

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

func TestFevmChainService(t *testing.T) {
	one := big.NewInt(1)

	cs := NewFevmChainService()

	// Prepare test data to trigger EthChainService
	testDeposit := types.Funds{
		common.HexToAddress("0x00"): one,
	}
	channelID := types.Destination(common.HexToHash(`4ebd366d014a173765ba1e50f284c179ade31f20441bec41664712aac6cc4611`))
	testTx := protocols.NewDepositTransaction(channelID, testDeposit)

	out := cs.EventFeed()
	cs.Monitor(channelID, testDeposit, testDeposit)

	// Submit transactiom
	err := cs.SendTransaction(testTx)
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
}
