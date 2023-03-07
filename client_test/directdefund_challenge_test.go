// Package client_test contains helpers and integration tests for go-nitro clients
package client_test // import "github.com/statechannels/go-nitro/client_test"

import (
	"testing"

	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	"github.com/statechannels/go-nitro/types"
)

func TestDirectDefundChallenge(t *testing.T) {

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

	chainB, err := chainservice.NewSimulatedBackendChainService(sim, bindings, ethAccounts[2], logDestination)
	if err != nil {
		t.Fatal(err)
	}
	// End chain service setup

	broker := messageservice.NewBroker()

	// Client setup
	clientA, _ := setupClient(alice.PrivateKey, chainA, broker, logDestination, 0)

	clientB, _ := setupClient(bob.PrivateKey, chainB, broker, logDestination, 0)
	// End Client setup

	channelId := directlyFundALedgerChannel(t, clientA, clientB, types.Address{})

	clientB.Stop() // Stopping client B prevents the cooperative close of the channel.
	// clientA should detect this and revert to challenging the channel on chain.

	id := clientA.CloseLedgerChannel(channelId)
	waitTimeForCompletedObjectiveIds(t, &clientA, defaultTimeout, id)
	waitTimeForCompletedObjectiveIds(t, &clientB, defaultTimeout, id)

}
