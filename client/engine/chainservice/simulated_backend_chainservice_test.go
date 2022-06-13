package chainservice

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/go-cmp/cmp"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

func TestSimulatedBackendChainService(t *testing.T) {
	sim, na, naAddress, ethAccounts, err := SetupSimulatedBackend(1)
	if err != nil {
		t.Fatal(err)
	}

	cs := NewSimulatedBackendChainService(sim, sim, na, naAddress, ethAccounts[0])

	// Prepare test data to trigger EthChainService
	testDeposit := types.Funds{
		common.HexToAddress("0x00"): big.NewInt(1),
	}
	channelID := types.Destination(common.HexToHash(`4ebd366d014a173765ba1e50f284c179ade31f20441bec41664712aac6cc461d`))
	testTx := protocols.ChainTransaction{
		ChannelId: channelID,
		Deposit:   testDeposit,
		Type:      protocols.DepositTransactionType,
	}

	out := cs.SubscribeToEvents(ethAccounts[0].From)
	// Submit transactiom
	cs.SendTransaction(testTx)

	// Check that the recieved event matches the expected event
	receivedEvent := <-out
	expectedEvent := DepositedEvent{CommonEvent: CommonEvent{channelID: channelID, BlockNum: 2}, Holdings: testDeposit}
	if diff := cmp.Diff(expectedEvent, receivedEvent, cmp.AllowUnexported(CommonEvent{})); diff != "" {
		t.Fatalf("Received event did not match expectation; (-want +got):\n%s", diff)
	}

	sim.Close()
}
