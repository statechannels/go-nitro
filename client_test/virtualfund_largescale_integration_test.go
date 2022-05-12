package client_test

import (
	"encoding/json"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"runtime"
	"testing"
	"time"

	"github.com/statechannels/go-nitro/client"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	"github.com/statechannels/go-nitro/client/engine/store"
	nc "github.com/statechannels/go-nitro/crypto"
	td "github.com/statechannels/go-nitro/internal/testdata"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/protocols/virtualfund"
	"github.com/statechannels/go-nitro/types"
)

// TestLargeScaleVirtualFundIntegration may be used to test a "large scale" payment channel newtork.
// It uses terminology from the Filecoin retrieval market:
// It spins up one retrieval provider, one payment hub and several (a configurable number of) retrieval clients.
// The clients are instrumented and emit vector clock logs, which are combined into an output file ../artifacts/shiviz.log at the end of the test run.
// The output shiviz.log can be pasted into https://bestchai.bitbucket.io/shiviz/ to visualize the messages which are sent.
func TestLargeScaleVirtualFundIntegration(t *testing.T) {

	// t.Skip() // This test is skipped because it requires an external dependency to run.
	// go install github.com/DistributedClocks/GoVector@latest
	// You may need to add GOPATH/bin to your PATH

	// Increase numRetrievalClients to simulate multiple retrieval clients all wanting to pay the same retrieval provider through the same hub
	const numRetrievalClients = 1

	// Set a directory for each client to store vector clock logs, and cleat it.
	vectorClockLogDir := "../artifacts/vectorclock"
	os.RemoveAll(vectorClockLogDir)

	// Setup regular logging
	logFile := "largescale_client_test.log"
	truncateLog(logFile)
	logDestination := newLogWriter(logFile)

	// Setup central services
	chain := chainservice.NewMockChain()
	broker := messageservice.NewBroker()

	// Setup singleton (instrumented) clients
	retrievalProvider, retrievalProviderStore := setupClient(bob.PrivateKey, chain, broker, logDestination, 0)
	paymentHub, paymentHubStore := setupClient(irene.PrivateKey, chain, broker, logDestination, 0)

	// Connect RP to PH
	directlyFundALedgerChannel(t, retrievalProvider, paymentHub)

	// Setup a number of RCs, each with a ledger connection to PH
	retrievalClients := make([]client.Client, numRetrievalClients)
	rcStores := make([]store.Store, numRetrievalClients)
	for i := range retrievalClients {
		secretKey, _ := nc.GeneratePrivateKeyAndAddress()
		retrievalClients[i], rcStores[i] = setupClient(secretKey, chain, broker, logDestination, 0)
		directlyFundALedgerChannel(t, retrievalClients[i], paymentHub)
	}

	// Switch to instrumented clients
	chain = chainservice.NewMockChain()
	broker = messageservice.NewBroker()
	retrievalProvider = client.New(
		messageservice.NewVectorClockTestMessageService(bob.Address(), broker, 0, vectorClockLogDir, "RP"),
		chainservice.NewSimpleChainService(&chain, bob.Address()),
		retrievalProviderStore,
		logDestination,
	)
	paymentHub = client.New(
		messageservice.NewVectorClockTestMessageService(irene.Address(), broker, 0, vectorClockLogDir, "PH"),
		chainservice.NewSimpleChainService(&chain, irene.Address()),
		paymentHubStore,
		logDestination,
	)
	for i := range retrievalClients {
		retrievalClients[i] =
			client.New(
				messageservice.NewVectorClockTestMessageService(*retrievalClients[i].Address, broker, 0, vectorClockLogDir, "RC"+fmt.Sprint(i)),
				chainservice.NewSimpleChainService(&chain, *retrievalClients[i].Address),
				rcStores[i],
				logDestination,
			)

	}

	// All Retrieval Clients try to start a virtual channel with the retrievalProvider, through the Payment Hub
	for _, client := range retrievalClients {
		go createVirtualChannelWithRetrievalProvider(client, retrievalProvider)
	}

	// HACK: wait a second for stuff to happen (be better to wait for objectives to finish)
	<-time.After(1 * time.Second)

	// Write a snapshot of the RPs ledger channel to the logs (for interest)
	retrievalProviderHubConnection, _ := retrievalProviderStore.GetConsensusChannel(*paymentHub.Address)
	finalOutcome, _ := json.Marshal(retrievalProviderHubConnection.SupportedSignedState().State().Outcome)
	_, _ = logDestination.Write(finalOutcome)

	// Combine vector clock logs together, ready for input to the visualizer
	combineLogs(t, vectorClockLogDir, "shiviz.log")

}

// combineLogs runs the GoVector CLI utility
func combineLogs(t *testing.T, logDir string, combinedLogsFilename string) {
	_, filename, _, _ := runtime.Caller(1)
	logDir = path.Join(path.Dir(filename), logDir)
	// NOTE: you may need to add GOPATH to PATH
	_, err := exec.Command("GoVector", "--log_type", "shiviz", "--log_dir", logDir, "--outfile", path.Join(logDir, combinedLogsFilename)).Output()
	if err != nil {
		t.Fatal(err)
	}
}

func createVirtualChannelWithRetrievalProvider(c client.Client, retrievalProvider client.Client) protocols.ObjectiveId {
	withRetrievalProvider := virtualfund.ObjectiveRequest{
		MyAddress:    *c.Address,
		CounterParty: *retrievalProvider.Address,
		Intermediary: irene.Address(),
		Outcome: td.Outcomes.Create(
			*c.Address,
			*retrievalProvider.Address,
			1,
			1,
		),
		AppDefinition:     types.Address{},
		AppData:           types.Bytes{},
		ChallengeDuration: big.NewInt(0),
		Nonce:             rand.Int63(),
	}
	return c.CreateVirtualChannel(withRetrievalProvider).Id
}
