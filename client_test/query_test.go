package client_test

import (
	"math/big"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/client"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	"github.com/statechannels/go-nitro/internal/testdata"
	td "github.com/statechannels/go-nitro/internal/testdata"
	"github.com/statechannels/go-nitro/types"
)

const default_ledger_funding = 5_000_000

func TestLedgerLifecycle(t *testing.T) {

	logFile := "test_ledger_lifecycle.log"
	truncateLog(logFile)
	logDestination := newLogWriter(logFile)

	chain := chainservice.NewMockChain()
	chainServiceA := chainservice.NewMockChainService(chain, alice.Address())
	chainServiceI := chainservice.NewMockChainService(chain, irene.Address())
	broker := messageservice.NewBroker()

	aliceClient, _ := setupClient(alice.PrivateKey, chainServiceA, broker, logDestination, 0)
	ireneClient, _ := setupClient(irene.PrivateKey, chainServiceI, broker, logDestination, 0)

	// Set up an outcome that requires both participants to deposit
	outcome := testdata.Outcomes.Create(alice.Address(), irene.Address(), 7, 3, types.Address{})

	res := aliceClient.CreateLedgerChannel(irene.Address(), 0, outcome)
	ledgerId := res.ChannelId

	// Irene might not have received the objective yet so we only check alice
	checkLedgerChannel(t, ledgerId, outcome, client.Proposed, &aliceClient)

	waitTimeForCompletedObjectiveIds(t, &aliceClient, defaultTimeout, res.Id)
	waitTimeForCompletedObjectiveIds(t, &ireneClient, defaultTimeout, res.Id)

	checkLedgerChannel(t, ledgerId, outcome, client.Ready, &aliceClient, &ireneClient)

	closeId := aliceClient.CloseLedgerChannel(ledgerId)

	// Irene might not have received the objective yet so we only check alice
	checkLedgerChannel(t, ledgerId, outcome, client.Closing, &aliceClient)

	waitTimeForCompletedObjectiveIds(t, &aliceClient, defaultTimeout, closeId)
	waitTimeForCompletedObjectiveIds(t, &ireneClient, defaultTimeout, closeId)

	checkLedgerChannel(t, ledgerId, outcome, client.Complete, &aliceClient, &ireneClient)

}

func TestQueryPaymentChannels(t *testing.T) {

	// Setup logging
	logFile := "test_query.log"
	truncateLog(logFile)
	logDestination := newLogWriter(logFile)

	chain := chainservice.NewMockChain()
	chainServiceA := chainservice.NewMockChainService(chain, alice.Address())
	chainServiceB := chainservice.NewMockChainService(chain, bob.Address())
	chainServiceI := chainservice.NewMockChainService(chain, irene.Address())
	broker := messageservice.NewBroker()

	clientA, _ := setupClient(alice.PrivateKey, chainServiceA, broker, logDestination, 0)
	irene, _ := setupClient(irene.PrivateKey, chainServiceI, broker, logDestination, 0)
	clientB, _ := setupClient(bob.PrivateKey, chainServiceB, broker, logDestination, 0)
	ledgerAId := directlyFundALedgerChannel(t, clientA, irene, types.Address{})
	directlyFundALedgerChannel(t, clientB, irene, types.Address{})

	expectedLedgerA := client.LedgerChannelInfo{
		ID:     ledgerAId,
		Status: client.Ready,
		Balance: client.LedgerChannelBalance{
			AssetAddress:  types.Address{},
			Hub:           *irene.Address,
			Client:        *clientA.Address,
			ClientBalance: big.NewInt(default_ledger_funding),
			HubBalance:    big.NewInt(default_ledger_funding),
		}}

	fetchedLedgerA, err := clientA.GetLedgerChannel(ledgerAId)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(expectedLedgerA, fetchedLedgerA, cmp.AllowUnexported(big.Int{})); diff != "" {
		t.Fatalf("Query diff mismatch (-want +got):\n%s", diff)
	}

	id := clientA.CreateVirtualPaymentChannel(
		[]types.Address{*irene.Address},
		bob.Address(),
		0,
		td.Outcomes.Create(
			alice.Address(),
			bob.Address(),
			2,
			0,
			types.Address{},
		)).ChannelId

	res, err := clientA.GetPaymentChannel(id)
	if err != nil {
		t.Fatal(err)
	}

	expected := client.PaymentChannelInfo{
		ID:     id,
		Status: client.Proposed,
		Balance: client.PaymentChannelBalance{
			AssetAddress:   types.Address{},
			Payee:          *clientB.Address,
			Payer:          *clientA.Address,
			PaidSoFar:      big.NewInt(0),
			RemainingFunds: big.NewInt(2),
		}}

	if diff := cmp.Diff(expected, res, cmp.AllowUnexported(big.Int{})); diff != "" {
		t.Fatalf("Query diff mismatch (-want +got):\n%s", diff)
	}

}

// expectedInfo constructs a LedgerChannelInfo so we can easily compare it to the result of GetLedgerChannel
func expectedInfo(id types.Destination, outcome outcome.Exit, status client.ChannelStatus) client.LedgerChannelInfo {
	clientAdd, _ := outcome[0].Allocations[0].Destination.ToAddress()
	hubAdd, _ := outcome[0].Allocations[1].Destination.ToAddress()

	return client.LedgerChannelInfo{
		ID:     id,
		Status: status,
		Balance: client.LedgerChannelBalance{
			AssetAddress:  types.Address{},
			Hub:           hubAdd,
			Client:        clientAdd,
			ClientBalance: outcome[0].Allocations[0].Amount,
			HubBalance:    outcome[0].Allocations[1].Amount,
		}}
}

// checkLedgerChannel checks that the ledger channel has the expected outcome and status
// It will fail if the channel does not exist
func checkLedgerChannel(t *testing.T, ledgerId types.Destination, o outcome.Exit, status client.ChannelStatus, clients ...*client.Client) {

	for _, c := range clients {
		expected := expectedInfo(ledgerId, o, status)
		ledger, err := c.GetLedgerChannel(ledgerId)
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(expected, ledger, cmp.AllowUnexported(big.Int{})); diff != "" {
			t.Fatalf("Query diff mismatch (-want +got):\n%s", diff)
		}
	}
}
