// Package client_test contains helpers and integration tests for go-nitro clients
package client_test // import "github.com/statechannels/go-nitro/client_test"

import (
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/statechannels/go-nitro/client"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	"github.com/statechannels/go-nitro/client/engine/store"
	"github.com/statechannels/go-nitro/internal/testdata"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/protocols/virtualfund"
	"github.com/statechannels/go-nitro/types"
)

func directlyDefundALedgerChannel(t *testing.T, alpha client.Client, beta client.Client, channelId types.Destination) {

	id := alpha.CloseLedgerChannel(channelId)
	waitTimeForCompletedObjectiveIds(t, &alpha, defaultTimeout, id)
	waitTimeForCompletedObjectiveIds(t, &beta, defaultTimeout, id)

}
func TestDirectDefund(t *testing.T) {

	// Setup logging
	logFile := "test_direct_defund.log"
	truncateLog(logFile)
	logDestination := newLogWriter(logFile)

	// Setup chain service
	sim, bindings, ethAccounts, err := chainservice.SetupSimulatedBackend(3)
	if err != nil {
		t.Fatal(err)
	}

	chainA, err := chainservice.NewSimulatedBackendChainService(sim, bindings, ethAccounts[0], logDestination)
	if err != nil {
		t.Fatal(err)
	}
	chainI, err := chainservice.NewSimulatedBackendChainService(sim, bindings, ethAccounts[1], logDestination)
	if err != nil {
		t.Fatal(err)
	}
	chainB, err := chainservice.NewSimulatedBackendChainService(sim, bindings, ethAccounts[2], logDestination)
	if err != nil {
		t.Fatal(err)
	}
	// End chain service setup

	broker := messageservice.NewBroker()

	// Client setup
	clientA, storeA := setupClient(alice.PrivateKey, chainA, broker, logDestination, 0)
	clientI, _ := setupClient(irene.PrivateKey, chainI, broker, logDestination, 0)
	clientB, storeB := setupClient(bob.PrivateKey, chainB, broker, logDestination, 0)
	// End Client setup

	// test successful condition for setup / teadown of unused ledger channel
	{
		channelId := directlyFundALedgerChannel(t, clientA, clientB)
		directlyDefundALedgerChannel(t, clientA, clientB, channelId)

		// Ensure that we no longer have a consensus channel in the store
		// And that we have a regular Channel instead
		for _, clientStore := range []store.Store{storeA, storeB} {

			// Ensure that we have a regular channel in the store
			// And that we no longer have a consensus channel in the store
			c, channelInStore := clientStore.GetChannelById(channelId)
			_, err := clientStore.GetConsensusChannelById(channelId)
			if !channelInStore {
				t.Fatalf("expected a Channel to have been created")
			}
			if consensusChannelStillInStore := (err == nil); consensusChannelStillInStore {
				t.Fatalf("Expected ConsensusChannel to have been destroyed in %v's store, but it was not", clientStore.GetAddress())
			}

			if c.OnChainFunding.IsNonZero() {
				t.Fatal("Expected zero on chain funding, but got nonzero")
			}

		}
	}

	// test failure of teardown of a ledger channel currently funding a virtual channel
	{
		ddfoTarget := directlyFundALedgerChannel(t, clientA, clientI)
		directlyFundALedgerChannel(t, clientI, clientB)

		// create & virtual channel between A and B through I
		outcome := testdata.Outcomes.Create(alice.Address(), bob.Address(), 1, 1)
		request := virtualfund.ObjectiveRequest{
			CounterParty:      bob.Address(),
			Intermediary:      irene.Address(),
			Outcome:           outcome,
			AppDefinition:     types.Address{},
			AppData:           types.Bytes{},
			ChallengeDuration: big.NewInt(0),
			Nonce:             rand.Int63(),
		}
		response := clientA.CreateVirtualChannel(request)

		waitTimeForCompletedObjectiveIds(t, &clientA, defaultTimeout, response.Id)

		clientA.CloseLedgerChannel(ddfoTarget)

		select {
		case <-time.After(time.Second * 10):
			t.Fatalf("expected ddfo on active ledger channel to fail, but no failures occurred within 10 seconds")
		case rejected := <-clientA.FailedObjectives():
			if rejected != protocols.ObjectiveId("DirectDefunding-"+ddfoTarget.String()) {
				t.Errorf("expected ddfo on active ledger channel to fail, but it didn't")
			}
		}

	}
}
