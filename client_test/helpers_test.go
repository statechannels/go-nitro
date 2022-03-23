package client_test

import (
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/statechannels/go-nitro/client"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	"github.com/statechannels/go-nitro/client/engine/store"
	"github.com/statechannels/go-nitro/crypto"
	"github.com/statechannels/go-nitro/protocols"
)

const defaultTimeout = time.Second

// waitWithTimeoutForCompletedObjectiveIds waits up to the given timeout for completed objectives and returns when the all objective ids provided have been completed.
// If the timeout lapses and the objectives have not all completed, the parent test will be failed.
func waitTimeForCompletedObjectiveIds(t *testing.T, client *client.Client, timeout time.Duration, ids ...protocols.ObjectiveId) {
	waitAndSendOn := func(allDone chan interface{}) {
		waitForCompletedObjectiveIds(client, ids...)
		allDone <- struct{}{}
	}
	allDone := make(chan interface{})
	go waitAndSendOn(allDone)

	select {
	case <-time.After(timeout):
		t.Fatalf("Objective ids %s failed to complete in one second on client %s", ids, client.Address)
	case <-allDone:
		return
	}
}

// waitForCompletedObjectiveIds waits for completed objectives and returns when the all objective ids provided have been completed.
func waitForCompletedObjectiveIds(client *client.Client, ids ...protocols.ObjectiveId) {
	// Create a map of all objective ids to wait for and set to false
	completed := make(map[protocols.ObjectiveId]bool)
	for _, id := range ids {
		completed[id] = false
	}
	// We continue to consume completed objective ids from the chan until all have been completed
	for got := range client.CompletedObjectives() {
		// Mark the objective as completed
		completed[got] = true

		// If all objectives are completed we can return
		isDone := true
		for _, objectiveCompleted := range completed {
			isDone = isDone && objectiveCompleted
		}
		if isDone {
			return
		}
	}
}

// setupClient is a helper function that contructs a client and returns the new client and message service.
func setupClient(pk []byte, chain chainservice.MockChain, msgBroker messageservice.Broker, logFilename string) client.Client {
	myAddress := crypto.GetAddressFromSecretKeyBytes(pk)
	chain.Subscribe(myAddress)
	chainservice := chainservice.NewSimpleChainService(chain, myAddress)
	messageservice := messageservice.NewTestMessageService(myAddress, msgBroker)
	storeA := store.NewMockStore(pk)
	logDestination := newLogWriter(logFilename)
	return client.New(messageservice, chainservice, storeA, logDestination)
}

func truncateLog(logFile string) {
	logDestination := newLogWriter(logFile)

	err := logDestination.Truncate(0)
	if err != nil {
		log.Fatal(err)
	}
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
