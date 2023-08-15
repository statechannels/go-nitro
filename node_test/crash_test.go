package node_test // import "github.com/statechannels/go-nitro/node_test"

import (
	"fmt"
	"os"
	"testing"

	"github.com/statechannels/go-nitro/internal/logging"
	ta "github.com/statechannels/go-nitro/internal/testactors"
	"github.com/statechannels/go-nitro/node"
	"github.com/statechannels/go-nitro/node/engine"
	"github.com/statechannels/go-nitro/node/engine/chainservice"
	"github.com/statechannels/go-nitro/node/engine/messageservice"
	"github.com/statechannels/go-nitro/node/engine/store"
	"github.com/statechannels/go-nitro/rand"
	"github.com/statechannels/go-nitro/types"
	"github.com/tidwall/buntdb"
)

func TestCrashTolerance(t *testing.T) {
	// Setup logging
	logFile := "test_crash_tolerance.log"
	logDestination := logging.NewLogWriter("../artifacts", logFile)

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
	storeA, err := store.NewDurableStore(ta.Alice.PrivateKey, dataFolder, buntdb.Config{SyncPolicy: buntdb.Always})
	if err != nil {
		t.Fatal(err)
	}
	messageserviceA := messageservice.NewTestMessageService(ta.Alice.Address(), broker, 0)
	nodeA := node.New(messageserviceA, chainA, storeA, logDestination, &engine.PermissivePolicy{}, nil)

	nodeB, _ := setupNode(ta.Bob.PrivateKey, chainB, broker, logDestination, 0)
	defer closeNode(t, &nodeB)
	// End Client setup

	t.Log("Node setup complete")

	// test successful condition for setup / teardown of unused ledger channel
	{
		channelId := openLedgerChannel(t, nodeA, nodeB, types.Address{})

		closeNode(t, &nodeA)
		anotherMessageserviceA := messageservice.NewTestMessageService(ta.Alice.Address(), broker, 0)
		anotherChainA, err := chainservice.NewSimulatedBackendChainService(sim, bindings, ethAccounts[0], logDestination)
		if err != nil {
			t.Fatal(err)
		}
		anotherStoreA, err := store.NewDurableStore(ta.Alice.PrivateKey, dataFolder, buntdb.Config{SyncPolicy: buntdb.Always})
		if err != nil {
			t.Fatal(err)
		}
		anotherClientA := node.New(
			anotherMessageserviceA,
			anotherChainA,
			anotherStoreA, logDestination, &engine.PermissivePolicy{}, nil)
		defer closeNode(t, &anotherClientA)

		closeLedgerChannel(t, anotherClientA, nodeB, channelId)

	}
}
