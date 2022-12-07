package pingpong

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/statechannels/go-nitro/client"
	"github.com/statechannels/go-nitro/client/engine"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	"github.com/statechannels/go-nitro/client/engine/store"
	"github.com/statechannels/go-nitro/crypto"
	"github.com/statechannels/go-nitro/internal/testdata"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

func TruncateLog(logFile string) {
	logDestination := NewLogWriter(logFile)

	err := logDestination.Truncate(0)
	if err != nil {
		log.Fatal(err)
	}
}

func NewLogWriter(logFile string) *os.File {
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

// SetupClient is a helper function that contructs a client and returns the new client and its store.
func SetupClient(pk []byte, chain chainservice.ChainService, msgBroker messageservice.Broker, logDestination io.Writer, meanMessageDelay time.Duration) (client.Client, store.Store, *messageservice.TestMessageService) {
	myAddress := crypto.GetAddressFromSecretKeyBytes(pk)
	messageservice := messageservice.NewTestMessageService(myAddress, msgBroker, meanMessageDelay)
	storeA := store.NewMemStore(pk)
	return client.New(messageservice, chain, storeA, logDestination, &engine.PermissivePolicy{}, nil), storeA, &messageservice
}

const ledgerChannelDeposit = 5_000_000
const defaultTimeout = 10 * time.Second

func directlyFundALedgerChannel(t *testing.T, alpha client.Client, beta client.Client) types.Destination {
	// Set up an outcome that requires both participants to deposit
	outcome := testdata.Outcomes.Create(*alpha.Address, *beta.Address, ledgerChannelDeposit, ledgerChannelDeposit)
	response := alpha.CreateLedgerChannel(*beta.Address, 0, outcome)

	waitTimeForCompletedObjectiveIds(t, &alpha, defaultTimeout, response.Id)
	waitTimeForCompletedObjectiveIds(t, &beta, defaultTimeout, response.Id)
	return response.ChannelId
}

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
