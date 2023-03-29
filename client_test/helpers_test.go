package client_test

import (
	"fmt"
	"io"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/go-cmp/cmp"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/client"
	"github.com/statechannels/go-nitro/client/engine"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	"github.com/statechannels/go-nitro/client/engine/store"
	"github.com/statechannels/go-nitro/client/query"
	"github.com/statechannels/go-nitro/crypto"
	"github.com/statechannels/go-nitro/internal/safesync"
	"github.com/statechannels/go-nitro/internal/testdata"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/rand"
	"github.com/statechannels/go-nitro/types"
	"github.com/tidwall/buntdb"
)

const TEST_CHAIN_ID = 1337

const defaultTimeout = 10 * time.Second

const DURABLE_STORE_FOLDER = "../data/client_test"

const ledgerChannelDeposit = 5_000_000

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

// setupClient is a helper function that contructs a client and returns the new client and its store.
func setupClient(pk []byte, chain chainservice.ChainService, msgBroker messageservice.Broker, logDestination io.Writer, meanMessageDelay time.Duration) (client.Client, store.Store) {
	myAddress := crypto.GetAddressFromSecretKeyBytes(pk)
	// TODO: Clean up test data folder?
	dataFolder := fmt.Sprintf("%s/%s/%d", DURABLE_STORE_FOLDER, myAddress.String(), rand.Uint64())
	messageservice := messageservice.NewTestMessageService(myAddress, msgBroker, meanMessageDelay)
	storeA := store.NewDurableStore(pk, dataFolder, buntdb.Config{})
	return client.New(messageservice, chain, storeA, logDestination, &engine.PermissivePolicy{}, nil), storeA
}

func closeClient(t *testing.T, client *client.Client) {
	err := client.Close()
	if err != nil {
		t.Fatal(err)
	}
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

// checkPaymentChannel checks that the ledger channel has the expected outcome and status
// It will fail if the channel does not exist
func checkPaymentChannel(t *testing.T, id types.Destination, o outcome.Exit, status query.ChannelStatus, clients ...*client.Client) {

	for _, c := range clients {
		expected := expectedPaymentInfo(id, o, status)
		ledger, err := c.GetPaymentChannel(id)
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(expected, ledger, cmp.AllowUnexported(big.Int{})); diff != "" {
			t.Fatalf("Payment channel diff mismatch (-want +got):\n%s", diff)
		}
	}
}

// expectedLedgerInfo constructs a LedgerChannelInfo so we can easily compare it to the result of GetLedgerChannel
func expectedLedgerInfo(id types.Destination, outcome outcome.Exit, status query.ChannelStatus) query.LedgerChannelInfo {
	clientAdd, _ := outcome[0].Allocations[0].Destination.ToAddress()
	hubAdd, _ := outcome[0].Allocations[1].Destination.ToAddress()

	return query.LedgerChannelInfo{
		ID:     id,
		Status: status,
		Balance: query.LedgerChannelBalance{
			AssetAddress:  types.Address{},
			Hub:           hubAdd,
			Client:        clientAdd,
			ClientBalance: outcome[0].Allocations[0].Amount,
			HubBalance:    outcome[0].Allocations[1].Amount,
		}}
}

// checkLedgerChannel checks that the ledger channel has the expected outcome and status
// It will fail if the channel does not exist
func checkLedgerChannel(t *testing.T, ledgerId types.Destination, o outcome.Exit, status query.ChannelStatus, clients ...*client.Client) {

	for _, c := range clients {
		expected := expectedLedgerInfo(ledgerId, o, status)
		ledger, err := c.GetLedgerChannel(ledgerId)
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(expected, ledger, cmp.AllowUnexported(big.Int{})); diff != "" {
			t.Fatalf("Ledger diff mismatch (-want +got):\n%s", diff)
		}
	}
}

func closeSimulatedChain(t *testing.T, chain chainservice.SimulatedChain) {
	if err := chain.Close(); err != nil {
		t.Fatal(err)
	}
}

func directlyDefundALedgerChannel(t *testing.T, alpha client.Client, beta client.Client, channelId types.Destination) {

	id := alpha.CloseLedgerChannel(channelId)
	waitTimeForCompletedObjectiveIds(t, &alpha, defaultTimeout, id)
	waitTimeForCompletedObjectiveIds(t, &beta, defaultTimeout, id)

}
func directlyFundALedgerChannel(t *testing.T, alpha client.Client, beta client.Client, asset common.Address) types.Destination {
	// Set up an outcome that requires both participants to deposit
	outcome := testdata.Outcomes.Create(*alpha.Address, *beta.Address, ledgerChannelDeposit, ledgerChannelDeposit, asset)

	response := alpha.CreateLedgerChannel(*beta.Address, 0, outcome)

	<-alpha.ObjectiveCompleteChan(response.Id)
	<-beta.ObjectiveCompleteChan(response.Id)

	return response.ChannelId
}
