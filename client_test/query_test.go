package client_test

import (
	"math/big"
	"testing"

	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	"github.com/statechannels/go-nitro/client/query"
	"github.com/statechannels/go-nitro/internal/testdata"

	"github.com/statechannels/go-nitro/types"
)

func TestQueryLedgerChannel(t *testing.T) {

	logFile := "test_query_ledger_channel.log"
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

	// It is possible the objective completes for Alice before we query it
	// so the status could be either Proposed or Ready
	// Irene may not have received the objective yet so we only check Alice
	ledger, err := aliceClient.GetLedgerChannel(ledgerId)
	if err != nil {
		t.Fatal(err)
	}
	if ledger.Status != query.Proposed && ledger.Status != query.Ready {
		t.Fatalf("Expected status to be Proposed or Ready but got %v", ledger.Status)
	}
	waitTimeForCompletedObjectiveIds(t, &aliceClient, defaultTimeout, res.Id)
	waitTimeForCompletedObjectiveIds(t, &ireneClient, defaultTimeout, res.Id)

	checkLedgerChannel(t, ledgerId, outcome, query.Ready, &aliceClient, &ireneClient)

	closeId := aliceClient.CloseLedgerChannel(ledgerId)

	// Irene might not have received the objective yet so we only check alice
	checkLedgerChannel(t, ledgerId, outcome, query.Closing, &aliceClient)

	waitTimeForCompletedObjectiveIds(t, &aliceClient, defaultTimeout, closeId)
	waitTimeForCompletedObjectiveIds(t, &ireneClient, defaultTimeout, closeId)

	checkLedgerChannel(t, ledgerId, outcome, query.Complete, &aliceClient, &ireneClient)

}

func TestQueryPaymentChannel(t *testing.T) {

	// Setup logging
	logFile := "test_query_payment_channel.log"
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

	o := testdata.Outcomes.Create(
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
		testdata.Outcomes.Create(
			alice.Address(),
			bob.Address(),
			2,
			0,
			types.Address{},
		))

	checkPaymentChannel(t, res.ChannelId, o, query.Proposed, &aliceClient)

	waitTimeForCompletedObjectiveIds(t, &aliceClient, defaultTimeout, res.Id)
	waitTimeForCompletedObjectiveIds(t, &bobClient, defaultTimeout, res.Id)
	waitTimeForCompletedObjectiveIds(t, &ireneClient, defaultTimeout, res.Id)

	// TODO Irene will return proposed because she doesn't bother listening for the post fund setup
	checkPaymentChannel(t, res.ChannelId, o, query.Ready, &aliceClient, &bobClient)

	aliceClient.Pay(res.ChannelId, big.NewInt(1))
	<-bobClient.ReceivedVouchers()
	updatedOutcome := testdata.Outcomes.Create(alice.Address(),
		bob.Address(),
		1,
		1,
		types.Address{})

	checkPaymentChannel(t, res.ChannelId, updatedOutcome, query.Ready, &aliceClient, &bobClient)

	closeId := aliceClient.CloseVirtualChannel(res.ChannelId)

	checkPaymentChannel(t, res.ChannelId, updatedOutcome, query.Closing, &aliceClient)

	waitTimeForCompletedObjectiveIds(t, &aliceClient, defaultTimeout, closeId)
	waitTimeForCompletedObjectiveIds(t, &bobClient, defaultTimeout, closeId)
	waitTimeForCompletedObjectiveIds(t, &ireneClient, defaultTimeout, closeId)

	checkPaymentChannel(t, res.ChannelId, updatedOutcome, query.Complete, &aliceClient, &bobClient, &ireneClient)
}

// expectedPaymentInfo constructs a LedgerChannelInfo so we can easily compare it to the result of GetPaymentChannel
func expectedPaymentInfo(id types.Destination, outcome outcome.Exit, status query.ChannelStatus) query.PaymentChannelInfo {
	payer, _ := outcome[0].Allocations[0].Destination.ToAddress()
	payee, _ := outcome[0].Allocations[1].Destination.ToAddress()

	return query.PaymentChannelInfo{
		ID:     id,
		Status: status,
		Balance: query.PaymentChannelBalance{
			AssetAddress:   types.Address{},
			Payee:          payee,
			Payer:          payer,
			RemainingFunds: outcome[0].Allocations[0].Amount,
			PaidSoFar:      outcome[0].Allocations[1].Amount,
		}}
}
