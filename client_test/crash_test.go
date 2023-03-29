// Package client_test contains helpers and integration tests for go-nitro clients
package client_test // import "github.com/statechannels/go-nitro/client_test"

import (
	"fmt"
	"os"
	"testing"

	"github.com/statechannels/go-nitro/client"
	"github.com/statechannels/go-nitro/client/engine"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	"github.com/statechannels/go-nitro/client/engine/store"
	"github.com/statechannels/go-nitro/rand"
	"github.com/statechannels/go-nitro/types"
	"github.com/tidwall/buntdb"

	ta "github.com/statechannels/go-nitro/internal/testactors"
)

func TestCrashTolerance(t *testing.T) {

	// Setup logging
	logFile := "test_crash_tolerance.log"
	truncateLog(logFile)
	logDestination := newLogWriter(logFile)

	// Setup chain service
	sim, bindings, ethAccounts, err := chainservice.SetupSimulatedBackend(3)
	defer closeSimulatedChain(t, sim)
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

	dataFolder := fmt.Sprintf("../data/%d", rand.Uint64())
	defer os.RemoveAll(dataFolder)

	// Client setup
	storeA := store.NewDurableStore(ta.Alice.PrivateKey, dataFolder, buntdb.Config{SyncPolicy: buntdb.Always})
	messageserviceA := messageservice.NewTestMessageService(ta.Alice.Address(), broker, 0)
	clientA := client.New(messageserviceA, chainA, storeA, logDestination, &engine.PermissivePolicy{}, nil)

	clientB, _ := setupClient(ta.Bob.PrivateKey, chainB, broker, logDestination, 0)
	defer closeClient(t, &clientB)
	// End Client setup

	// test successful condition for setup / teadown of unused ledger channel
	{
		channelId := directlyFundALedgerChannel(t, clientA, clientB, types.Address{})

		closeClient(t, &clientA)
		anotherMessageserviceA := messageservice.NewTestMessageService(ta.Alice.Address(), broker, 0)
		anotherChainA, err := chainservice.NewSimulatedBackendChainService(sim, bindings, ethAccounts[0], logDestination)
		anotherStoreA := store.NewDurableStore(ta.Alice.PrivateKey, dataFolder, buntdb.Config{SyncPolicy: buntdb.Always})
		if err != nil {
			t.Fatal(err)
		}
		anotherClientA := client.New(
			anotherMessageserviceA,
			anotherChainA,
			anotherStoreA, logDestination, &engine.PermissivePolicy{}, nil)
		defer closeClient(t, &anotherClientA)

		directlyDefundALedgerChannel(t, anotherClientA, clientB, channelId)

	}

}
