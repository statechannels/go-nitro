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

func TestPaymentChannelLifecycle(t *testing.T) {

	// Setup logging
	logFile := "test_payment_lifecycle.log"
	truncateLog(logFile)
	logDestination := newLogWriter(logFile)

	chain := chainservice.NewMockChain()
	chainServiceA := chainservice.NewMockChainService(chain, alice.Address())
	chainServiceB := chainservice.NewMockChainService(chain, bob.Address())
	chainServiceI := chainservice.NewMockChainService(chain, irene.Address())
	broker := messageservice.NewBroker()

	aliceClient, _ := setupClient(alice.PrivateKey, chainServiceA, broker, logDestination, 0)
	ireneClient, _ := setupClient(irene.PrivateKey, chainServiceI, broker, logDestination, 0)
	bobClient, _ := setupClient(bob.PrivateKey, chainServiceB, broker, logDestination, 0)

	directlyFundALedgerChannel(t, aliceClient, ireneClient, types.Address{})
	directlyFundALedgerChannel(t, bobClient, ireneClient, types.Address{})

	o := td.Outcomes.Create(
		alice.Address(),
		bob.Address(),
		2,
		0,
		types.Address{},
	)

	res := aliceClient.CreateVirtualPaymentChannel(
		[]types.Address{*ireneClient.Address},
		bob.Address(),
		0,
		td.Outcomes.Create(
			alice.Address(),
			bob.Address(),
			2,
			0,
			types.Address{},
		))

	checkPaymentChannel(t, res.ChannelId, o, client.Proposed, &aliceClient)

	waitTimeForCompletedObjectiveIds(t, &aliceClient, defaultTimeout, res.Id)
	waitTimeForCompletedObjectiveIds(t, &bobClient, defaultTimeout, res.Id)
	waitTimeForCompletedObjectiveIds(t, &ireneClient, defaultTimeout, res.Id)

	// TODO Irene will return proposed because she doesn't bother listening for the post fund setup
	checkPaymentChannel(t, res.ChannelId, o, client.Ready, &aliceClient, &bobClient)

	aliceClient.Pay(res.ChannelId, big.NewInt(1))

	updatedOutcome := td.Outcomes.Create(alice.Address(),
		bob.Address(),
		1,
		1,
		types.Address{})
	checkPaymentChannel(t, res.ChannelId, updatedOutcome, client.Ready, &aliceClient, &bobClient)

	closeId := aliceClient.CloseVirtualChannel(res.ChannelId)

	checkPaymentChannel(t, res.ChannelId, updatedOutcome, client.Closing, &aliceClient)

	waitTimeForCompletedObjectiveIds(t, &aliceClient, defaultTimeout, closeId)
	waitTimeForCompletedObjectiveIds(t, &bobClient, defaultTimeout, closeId)
	waitTimeForCompletedObjectiveIds(t, &ireneClient, defaultTimeout, closeId)

	checkPaymentChannel(t, res.ChannelId, updatedOutcome, client.Complete, &aliceClient, &bobClient, &ireneClient)
}

// expectedPaymentInfo constructs a LedgerChannelInfo so we can easily compare it to the result of GetPaymentChannel
func expectedPaymentInfo(id types.Destination, outcome outcome.Exit, status client.ChannelStatus) client.PaymentChannelInfo {
	payer, _ := outcome[0].Allocations[0].Destination.ToAddress()
	payee, _ := outcome[0].Allocations[1].Destination.ToAddress()

	return client.PaymentChannelInfo{
		ID:     id,
		Status: status,
		Balance: client.PaymentChannelBalance{
			AssetAddress:   types.Address{},
			Payee:          payee,
			Payer:          payer,
			RemainingFunds: outcome[0].Allocations[0].Amount,
			PaidSoFar:      outcome[0].Allocations[1].Amount,
		}}
}

// checkPaymentChannel checks that the ledger channel has the expected outcome and status
// It will fail if the channel does not exist
func checkPaymentChannel(t *testing.T, id types.Destination, o outcome.Exit, status client.ChannelStatus, clients ...*client.Client) {

	for _, c := range clients {
		expected := expectedPaymentInfo(id, o, status)
		ledger, err := c.GetPaymentChannel(id)
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(expected, ledger, cmp.AllowUnexported(big.Int{})); diff != "" {
			t.Fatalf("Payment channel diff mismatch (-want +got):\n%s", diff)
		}
	}
}

// expectedLedgerInfo constructs a LedgerChannelInfo so we can easily compare it to the result of GetLedgerChannel
func expectedLedgerInfo(id types.Destination, outcome outcome.Exit, status client.ChannelStatus) client.LedgerChannelInfo {
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
		expected := expectedLedgerInfo(ledgerId, o, status)
		ledger, err := c.GetLedgerChannel(ledgerId)
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(expected, ledger, cmp.AllowUnexported(big.Int{})); diff != "" {
			t.Fatalf("Ledger diff mismatch (-want +got):\n%s", diff)
		}
	}
}
