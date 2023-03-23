// Package client_test contains helpers and integration tests for go-nitro clients
package client_test // import "github.com/statechannels/go-nitro/client_test"

import (
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/statechannels/go-nitro/client"
	"github.com/statechannels/go-nitro/client/engine"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	"github.com/statechannels/go-nitro/client/engine/store"
	"github.com/statechannels/go-nitro/internal/testactors"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/tidwall/buntdb"
)

const STORE_TEST_DATA_FOLDER = "../data/store_test"

type StoreType string

const MemStore StoreType = "MemStore"
const PersistStore StoreType = "PersistStore"

func TestStoreImplementations(t *testing.T) {

	// Clean up all the test data we create at the end of the test
	defer os.RemoveAll(STORE_TEST_DATA_FOLDER)

	cases := []struct {
		StoreA StoreType
		StoreB StoreType
		StoreI StoreType
	}{
		{"MemStore", "MemStore", "MemStore"},
		{"MemStore", "PersistStore", "MemStore"},
		{"PersistStore", "PersistStore", "PersistStore"},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("Alice %s,Bob %s,Irene %s", tc.StoreA, tc.StoreB, tc.StoreI), func(t *testing.T) {

			logFile := "test_store_compatibility.log"
			truncateLog(logFile)
			logDestination := newLogWriter(logFile)

			broker := messageservice.NewBroker()
			messageserviceA := messageservice.NewTestMessageService(alice.Address(), broker, 0)
			messageserviceB := messageservice.NewTestMessageService(bob.Address(), broker, 0)
			messageserviceI := messageservice.NewTestMessageService(irene.Address(), broker, 0)

			storeA := openStore(alice, tc.StoreA)
			storeB := openStore(bob, tc.StoreB)
			storeI := openStore(irene, tc.StoreI)

			// Setup chain service
			sim, bindings, ethAccounts, err := chainservice.SetupSimulatedBackend(3)
			checkError(t, err)
			defer sim.Close()
			chainA, err := chainservice.NewSimulatedBackendChainService(sim, bindings, ethAccounts[0], logDestination)
			checkError(t, err)
			chainB, err := chainservice.NewSimulatedBackendChainService(sim, bindings, ethAccounts[2], logDestination)
			checkError(t, err)
			chainI, err := chainservice.NewSimulatedBackendChainService(sim, bindings, ethAccounts[1], logDestination)
			checkError(t, err)

			clientA := client.New(messageserviceA, chainA, storeA, logDestination, &engine.PermissivePolicy{}, nil)
			clientI := client.New(messageserviceI, chainI, storeI, logDestination, &engine.PermissivePolicy{}, nil)
			clientB := client.New(messageserviceB, chainB, storeB, logDestination, &engine.PermissivePolicy{}, nil)

			defer clientA.Close()
			defer clientB.Close()
			defer clientI.Close()

			cIds := openVirtualChannels(t, clientA, clientB, clientI, 3)
			for i := 0; i < len(cIds); i++ {
				clientA.Pay(cIds[i], big.NewInt(int64(1)))
			}

			ids := make([]protocols.ObjectiveId, len(cIds))
			for i := 0; i < len(cIds); i++ {
				// alternative who is responsible for closing the channel
				switch i % 3 {
				case 0:
					ids[i] = clientA.CloseVirtualChannel(cIds[i])
				case 1:
					ids[i] = clientB.CloseVirtualChannel(cIds[i])
				case 2:
					ids[i] = clientI.CloseVirtualChannel(cIds[i])
				}

			}
			waitTimeForCompletedObjectiveIds(t, &clientA, defaultTimeout, ids...)
			waitTimeForCompletedObjectiveIds(t, &clientB, defaultTimeout, ids...)
			waitTimeForCompletedObjectiveIds(t, &clientI, defaultTimeout, ids...)

		})
	}
}

func openStore(me testactors.Actor, storeType StoreType) store.Store {

	switch storeType {
	case "MemStore":
		return store.NewMemStore(me.PrivateKey)
	case "PersistStore":
		dataFolder := fmt.Sprintf("%s/%s/%d%d", STORE_TEST_DATA_FOLDER, me.Address().String(), rand.Uint64(), time.Now().UnixNano())

		return store.NewPersistStore(me.PrivateKey, dataFolder, buntdb.Config{})
	default:
		panic(fmt.Sprintf("Unknown store type %s", storeType))
	}

}

func checkError(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}
