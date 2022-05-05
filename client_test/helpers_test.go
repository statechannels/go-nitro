package client_test

import (
	"io"
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

	waitAndSendOn := func(completed map[protocols.ObjectiveId]bool, allDone chan interface{}) {

		// We continue to consume completed objective ids from the chan until all have been completed
		for got := range client.CompletedObjectives() {
			// Mark the objective as completed
			completed[got] = true

			// If all objectives are completed we can send the all done signal and return
			isDone := true
			for _, id := range ids {
				isDone = isDone && completed[id]
			}
			if isDone {
				allDone <- struct{}{}
				return

			}
		}

	}

	allDone := make(chan interface{})
	// Create a map to keep track of completed objectives
	completed := make(map[protocols.ObjectiveId]bool)

	go waitAndSendOn(completed, allDone)

	select {
	case <-time.After(timeout):
		incompleteIds := make([]protocols.ObjectiveId, 0)
		for _, id := range ids {
			isObjectiveDone := completed[id]
			if !isObjectiveDone {
				incompleteIds = append(incompleteIds, id)
			}
		}
		t.Fatalf("Objective ids %s failed to complete on client %s within %s", incompleteIds, client.Address, timeout)
	case <-allDone:
		return
	}
}

// setupClient is a helper function that contructs a client and returns the new client and its store.
func setupClient(pk []byte, chain chainservice.MockChain, msgBroker messageservice.Broker, logDestination io.Writer, meanMessageDelay time.Duration) (client.Client, store.Store) {
	myAddress := crypto.GetAddressFromSecretKeyBytes(pk)
	chain.Subscribe(myAddress)
	chainservice := chainservice.NewSimpleChainService(chain, myAddress)
	messageservice := messageservice.NewTestMessageService(myAddress, msgBroker, meanMessageDelay)
	storeA := store.NewMemStore(pk)
	return client.New(messageservice, chainservice, storeA, logDestination), storeA
}

// setupIntstrumentedClient is a helper function that contructs a client and returns the new client and its store.
//
// The client will be _instrumented_, which is useful in debugging. In particular it will output logs containing vector clock stamps into the supplied logDir directory.
func setupInstrumentedClient(pk []byte, chain chainservice.MockChain, msgBroker messageservice.Broker, logDestination io.Writer, meanMessageDelay time.Duration, logDir string, prettyName string) (client.Client, store.Store) {
	myAddress := crypto.GetAddressFromSecretKeyBytes(pk)
	chain.Subscribe(myAddress)
	chainservice := chainservice.NewSimpleChainService(chain, myAddress)
	messageservice := messageservice.NewVectorClockTestMessageService(myAddress, msgBroker, meanMessageDelay, logDir, prettyName)
	storeA := store.NewMemStore(pk)
	return client.New(messageservice, chainservice, storeA, logDestination), storeA
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
