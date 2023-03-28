package integration_test

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/client"
	"github.com/statechannels/go-nitro/client/engine"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	p2pms "github.com/statechannels/go-nitro/client/engine/messageservice/p2p-message-service"
	"github.com/statechannels/go-nitro/client/engine/store"
	"github.com/statechannels/go-nitro/internal/safesync"
	"github.com/statechannels/go-nitro/internal/testactors"
	"github.com/statechannels/go-nitro/internal/testdata"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
	"github.com/tidwall/buntdb"
)

func getActor(name testactors.ActorName) testactors.Actor {
	switch name {
	case testactors.AliceName:
		return testactors.Alice
	case testactors.BobName:
		return testactors.Bob
	case testactors.IreneName:
		return testactors.Irene
	case testactors.BrianName:
		return testactors.Brian
	default:
		panic("Unknown actor")
	}
}

func getTestParticipant(tc TestRun, name testactors.ActorName) TestParticipant {
	for _, p := range tc.Participants {
		if p.Name == name {
			return p
		}
	}
	panic("Unknown participant")
}

func getActorIndex(name testactors.ActorName) int {
	switch name {
	case testactors.AliceName:
		return 0
	case testactors.BobName:
		return 1
	case testactors.IreneName:
		return 2
	case testactors.BrianName:
		return 3
	default:
		panic("Unknown actor")
	}
}

func setupMessageService(tc TestRun, actorName testactors.ActorName, si sharedInra) messageservice.MessageService {
	actor := getActor(actorName)
	switch tc.MessageService {
	case TestMessageService:
		return messageservice.NewTestMessageService(actor.Address(), *si.broker, tc.MessageDelay)
	case P2PMessageService:
		return p2pms.NewMessageService("127.0.0.1", int(actor.Port), actor.PrivateKey)
	default:
		panic("Unknown message service")
	}
}

func setupChainService(tc TestRun, actorName testactors.ActorName, si sharedInra) chainservice.ChainService {
	a := getActor(actorName)
	switch tc.Chain {
	case MockChain:
		return chainservice.NewMockChainService(si.mockChain, a.Address())
	case SimulatedChain:
		logDestination := newLogWriter(tc.LogName)
		// TODO: USE unique eth address
		cs, err := chainservice.NewSimulatedBackendChainService(*si.simulatedChain, *si.bindings, si.ethAccounts[getActorIndex(actorName)], logDestination)
		if err != nil {
			panic(err)
		}
		return cs
	default:
		panic("Unknown chain service")
	}
}

func setupStore(tc TestRun, actorName testactors.ActorName, si sharedInra) store.Store {
	me := getActor(actorName)
	storeType := getTestParticipant(tc, actorName).StoreType
	switch getTestParticipant(tc, actorName).StoreType {
	case "MemStore":
		return store.NewMemStore(me.PrivateKey)
	case "PersistStore":
		dataFolder := fmt.Sprintf("%s/%s/%d%d", STORE_TEST_DATA_FOLDER, me.Address().String(), rand.Uint64(), time.Now().UnixNano())
		return store.NewPersistStore(me.PrivateKey, dataFolder, buntdb.Config{})
	default:
		panic(fmt.Sprintf("Unknown store type %s", storeType))
	}
}

func setupIntegrationClient(tc TestRun, actorName testactors.ActorName, si sharedInra) client.Client {

	messageService := setupMessageService(tc, actorName, si)
	cs := setupChainService(tc, actorName, si)
	store := setupStore(tc, actorName, si)
	return client.New(messageService, cs, store, newLogWriter(tc.LogName), &engine.PermissivePolicy{}, nil)
}

func newLogWriter(logFile string) *os.File {
	err := os.MkdirAll("../artifacts", os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	filename := filepath.Join("../artifacts", logFile)
	logDestination, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	if err != nil {
		log.Fatal(err)
	}

	return logDestination
}

func setupLedgerChannel(t *testing.T, alpha client.Client, beta client.Client, asset common.Address) types.Destination {
	// Set up an outcome that requires both participants to deposit
	outcome := testdata.Outcomes.Create(*alpha.Address, *beta.Address, ledgerChannelDeposit, ledgerChannelDeposit, asset)

	response := alpha.CreateLedgerChannel(*beta.Address, 0, outcome)

	<-alpha.ObjectiveCompleteChan(response.Id)
	<-beta.ObjectiveCompleteChan(response.Id)

	return response.ChannelId
}

func closeLedgerChannel(t *testing.T, alpha client.Client, beta client.Client, channelId types.Destination) {
	response := alpha.CloseLedgerChannel(channelId)

	<-alpha.ObjectiveCompleteChan(response)
	<-beta.ObjectiveCompleteChan(response)
}

// waitWithTimeoutForCompletedObjectiveIds waits up to the given timeout for completed objectives and returns when the all objective ids provided have been completed.
// If the timeout lapses and the objectives have not all completed, the parent test will be failed.
func waitTimeForCompletedObjectiveIds(t *testing.T, client *client.Client, timeout time.Duration, ids ...protocols.ObjectiveId) {

	incomplete := safesync.Map[<-chan struct{}]{}

	var wg sync.WaitGroup

	for _, id := range ids {
		incomplete.Store(string(id), client.ObjectiveCompleteChan(id))
		wg.Add(1)
	}

	incomplete.Range(
		func(id string, ch <-chan struct{}) bool {
			go func() {
				<-ch
				incomplete.Delete(string(id))
				wg.Done()
			}()
			return true
		})

	allDone := make(chan struct{})

	go func() {
		wg.Wait()
		allDone <- struct{}{}
	}()

	select {
	case <-time.After(timeout):
		incompleteIds := make([]string, 0)
		incomplete.Range(func(key string, value <-chan struct{}) bool {
			incompleteIds = append(incompleteIds, key)
			return true
		})
		t.Fatalf("Objective ids %s failed to complete on client %s within %s", incompleteIds, client.Address, timeout)
	case <-allDone:
		return
	}
}
