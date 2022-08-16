package client_test

import (
	"encoding/json"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"path"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/statechannels/go-nitro/client"
	"github.com/statechannels/go-nitro/client/engine"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	"github.com/statechannels/go-nitro/client/engine/store"
	"github.com/statechannels/go-nitro/client/engine/store/safesync"
	nc "github.com/statechannels/go-nitro/crypto"
	td "github.com/statechannels/go-nitro/internal/testdata"
	"github.com/statechannels/go-nitro/protocols/virtualfund"
	"github.com/statechannels/go-nitro/types"
)

// TestLargeScaleVirtualFundIntegration may be used to test a "large scale" payment channel newtork.
// It uses terminology from the Filecoin retrieval market:
// It spins up one retrieval provider, one payment hub and several (a configurable number of) retrieval clients.
// The clients are instrumented and emit vector clock logs, which are combined into an output file ../artifacts/shiviz.log at the end of the test run.
// The output shiviz.log can be pasted into https://bestchai.bitbucket.io/shiviz/ to visualize the messages which are sent.
func TestLargeScaleVirtualFundIntegration(t *testing.T) {

	prettyPrintDict := safesync.Map[string]{}

	// Increase numRetrievalClients to simulate multiple retrieval clients all wanting to pay the same retrieval provider through the same hub
	const numRetrievalClients = 3

	// Set a directory for each client to store vector clock logs, and cleat it.
	vectorClockLogDir := "../artifacts/vectorclock"
	os.RemoveAll(vectorClockLogDir)

	// Setup regular logging
	logFile := "largescale_client_test.log"
	truncateLog(logFile)
	logDestination := newLogWriter(logFile)

	// Setup central services
	chain := chainservice.NewMockChain()
	chainServiceB := chainservice.NewMockChainService(chain, bob.Address())
	chainServiceI := chainservice.NewMockChainService(chain, irene.Address())
	broker := messageservice.NewBroker()

	// Setup singleton (instrumented) clients
	retrievalProvider, retrievalProviderStore := setupClient(bob.PrivateKey, chainServiceB, broker, logDestination, 0)
	prettyPrintDict.Store(retrievalProvider.Address.String(), "RP")

	paymentHub, paymentHubStore := setupClient(irene.PrivateKey, chainServiceI, broker, logDestination, 0)
	prettyPrintDict.Store(paymentHub.Address.String(), "PH")

	// Connect RP to PH
	lID := directlyFundALedgerChannel(t, retrievalProvider, paymentHub)
	prettyPrintDict.Store(lID.String(), "L")

	// Setup a number of RCs, each with a ledger connection to PH
	retrievalClients := make([]client.Client, numRetrievalClients)
	rcStores := make([]store.Store, numRetrievalClients)
	for i := range retrievalClients {
		secretKey, address := nc.GeneratePrivateKeyAndAddress()
		chainService := chainservice.NewMockChainService(chain, address)
		retrievalClients[i], rcStores[i] = setupClient(secretKey, chainService, broker, logDestination, 0)
		prettyPrintDict.Store(retrievalClients[i].Address.String(), "RC"+fmt.Sprint(i))
		lID := directlyFundALedgerChannel(t, retrievalClients[i], paymentHub)
		prettyPrintDict.Store(lID.String(), "L"+fmt.Sprint(i))
	}

	// Switch to instrumented clients
	chainService := chainservice.NewMockChainService(chain, bob.Address())
	broker = messageservice.NewBroker()
	retrievalProvider = client.New(
		messageservice.NewVectorClockTestMessageService(bob.Address(), broker, 0, vectorClockLogDir),
		chainService,
		retrievalProviderStore,
		logDestination,
		&engine.PermissivePolicy{},
		nil,
	)
	chainService = chainservice.NewMockChainService(chain, irene.Address())
	paymentHub = client.New(
		messageservice.NewVectorClockTestMessageService(irene.Address(), broker, 0, vectorClockLogDir),
		chainService,
		paymentHubStore,
		logDestination,
		&engine.PermissivePolicy{},
		nil,
	)
	for i := range retrievalClients {
		chainService = chainservice.NewMockChainService(chain, *retrievalClients[i].Address)
		retrievalClients[i] =
			client.New(
				messageservice.NewVectorClockTestMessageService(*retrievalClients[i].Address, broker, 0, vectorClockLogDir),
				chainService,
				rcStores[i],
				logDestination,
				&engine.PermissivePolicy{},
				nil,
			)
	}

	// All Retrieval Clients try to start a virtual channel with the retrievalProvider, through the Payment Hub
	for i, client := range retrievalClients {
		go createVirtualChannelWithRetrievalProvider(client, retrievalProvider, &prettyPrintDict, i)
	}

	// HACK: wait a second for stuff to happen (be better to wait for objectives to finish)
	<-time.After(1 * time.Second)

	// Write a snapshot of the RPs ledger channel to the logs (for interest)
	retrievalProviderHubConnection, _ := retrievalProviderStore.GetConsensusChannel(*paymentHub.Address)
	finalOutcome, _ := json.Marshal(retrievalProviderHubConnection.SupportedSignedState().State().Outcome)
	_, _ = logDestination.Write(finalOutcome)

	// Combine vector clock logs together, ready for input to the visualizer
	combineLogs(t, vectorClockLogDir, "shiviz.log", retrievalProvider, paymentHub, retrievalClients)

	// prettify log
	prettify(t, vectorClockLogDir, "shiviz.log", &prettyPrintDict)

}

// combineLogs combines the logs into one file, and deletes them
func combineLogs(t *testing.T, logDir string, combinedLogsFilename string, RP client.Client, PH client.Client, retrievalClients []client.Client) {
	_, filename, _, _ := runtime.Caller(1)
	logDir = path.Join(path.Dir(filename), logDir)

	output := `(?<host>\S*) (?<clock>{.*})\n(?<event>.*)`
	output += "\n"
	output += "\n"

	// we want RP leftmost
	// we want PH immediately to the right of RP
	input, err := os.ReadFile(path.Join(logDir, RP.Address.String()+"-Log.txt"))
	if err != nil {
		t.Fatal(err)
	}
	os.Remove(path.Join(logDir, RP.Address.String()+"-Log.txt"))
	output += string(input)
	output += "\n"
	input, err = os.ReadFile(path.Join(logDir, PH.Address.String()+"-Log.txt"))
	if err != nil {
		t.Fatal(err)
	}
	os.Remove(path.Join(logDir, PH.Address.String()+"-Log.txt"))
	output += string(input)
	output += "\n"
	for i := range retrievalClients {
		input, err := os.ReadFile(path.Join(logDir, retrievalClients[i].Address.String()+"-Log.txt"))
		if err != nil {
			t.Fatal(err)
		}
		os.Remove(path.Join(logDir, retrievalClients[i].Address.String()+"-Log.txt"))
		output += string(input)
	}

	if err = os.WriteFile(path.Join(logDir, combinedLogsFilename), []byte(output), 0666); err != nil {
		t.Fatal(err)
	}

}

// prettify replaces addresses and destinations using the supplied prettyPrintDict
func prettify(t *testing.T, logDir string, combinedLogsFilename string, prettyPrintDict *safesync.Map[string]) {
	input, err := os.ReadFile(path.Join(logDir, combinedLogsFilename))
	if err != nil {
		t.Fatal(err)
	}

	output := string(input)
	replaceFun := func(key string, value string) bool {
		output = strings.Replace(string(output), key, value, -1)
		return true
	}
	prettyPrintDict.Range(replaceFun)

	if err = os.WriteFile(path.Join(logDir, combinedLogsFilename)+"_pretty", []byte(output), 0666); err != nil {
		t.Fatal(err)
	}
}

func createVirtualChannelWithRetrievalProvider(c client.Client, retrievalProvider client.Client, prettyPrintDict *safesync.Map[string], i int) {
	withRetrievalProvider := virtualfund.ObjectiveRequest{
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
	prettyPrintDict.Store(c.CreateVirtualChannel(withRetrievalProvider).ChannelId.String(), "V"+fmt.Sprint(i))
}
