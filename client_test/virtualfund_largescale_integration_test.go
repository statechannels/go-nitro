package client_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"math/rand"
	"os"
	"path"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/statechannels/go-nitro/client"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	"github.com/statechannels/go-nitro/client/engine/store"
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

	prettyPrintDict := make(map[string]string)

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
	broker := messageservice.NewBroker()

	// Setup singleton (instrumented) clients
	retrievalProvider, retrievalProviderStore := setupClient(bob.PrivateKey, chain, broker, logDestination, 0)
	prettyPrintDict[retrievalProvider.Address.String()] = "RP"

	paymentHub, paymentHubStore := setupClient(irene.PrivateKey, chain, broker, logDestination, 0)
	prettyPrintDict[paymentHub.Address.String()] = "PH"

	// Connect RP to PH
	lID := directlyFundALedgerChannel(t, retrievalProvider, paymentHub)
	prettyPrintDict[lID.String()] = "L"

	// Setup a number of RCs, each with a ledger connection to PH
	retrievalClients := make([]client.Client, numRetrievalClients)
	rcStores := make([]store.Store, numRetrievalClients)
	for i := range retrievalClients {
		secretKey, _ := nc.GeneratePrivateKeyAndAddress()
		retrievalClients[i], rcStores[i] = setupClient(secretKey, chain, broker, logDestination, 0)
		prettyPrintDict[retrievalClients[i].Address.String()] = "RC" + fmt.Sprint(i)
		lID := directlyFundALedgerChannel(t, retrievalClients[i], paymentHub)
		prettyPrintDict[lID.String()] = "L" + fmt.Sprint(i)
	}

	// Switch to instrumented clients
	chain = chainservice.NewMockChain()
	broker = messageservice.NewBroker()
	retrievalProvider = client.New(
		messageservice.NewVectorClockTestMessageService(bob.Address(), broker, 0, vectorClockLogDir),
		chainservice.NewSimpleChainService(&chain, bob.Address()),
		retrievalProviderStore,
		logDestination,
	)
	paymentHub = client.New(
		messageservice.NewVectorClockTestMessageService(irene.Address(), broker, 0, vectorClockLogDir),
		chainservice.NewSimpleChainService(&chain, irene.Address()),
		paymentHubStore,
		logDestination,
	)
	for i := range retrievalClients {
		retrievalClients[i] =
			client.New(
				messageservice.NewVectorClockTestMessageService(*retrievalClients[i].Address, broker, 0, vectorClockLogDir),
				chainservice.NewSimpleChainService(&chain, *retrievalClients[i].Address),
				rcStores[i],
				logDestination,
			)

	}

	// All Retrieval Clients try to start a virtual channel with the retrievalProvider, through the Payment Hub
	for i, client := range retrievalClients {
		go createVirtualChannelWithRetrievalProvider(client, retrievalProvider, prettyPrintDict, i)
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
	prettify(t, vectorClockLogDir, "shiviz.log", prettyPrintDict)

}

// combineLogs runs the GoVector CLI utility
func combineLogs(t *testing.T, logDir string, combinedLogsFilename string, RP client.Client, PH client.Client, retrievalClients []client.Client) {
	_, filename, _, _ := runtime.Caller(1)
	logDir = path.Join(path.Dir(filename), logDir)

	output := `(?<host>\S*) (?<clock>{.*})\n(?<event>.*)`
	output += "\n"
	output += "\n"

	// we want RP leftmost
	// we want PH immediately to the right of RP
	input, err := ioutil.ReadFile(path.Join(logDir, RP.Address.String()+"-Log.txt"))
	if err != nil {
		t.Fatal(err)
	}
	output += string(input)
	output += "\n"
	input, err = ioutil.ReadFile(path.Join(logDir, PH.Address.String()+"-Log.txt"))
	if err != nil {
		t.Fatal(err)
	}
	output += string(input)
	output += "\n"
	for i := range retrievalClients {
		input, err := ioutil.ReadFile(path.Join(logDir, retrievalClients[i].Address.String()+"-Log.txt"))
		if err != nil {
			t.Fatal(err)
		}
		output += string(input)
	}

	if err = ioutil.WriteFile(path.Join(logDir, combinedLogsFilename), []byte(output), 0666); err != nil {
		t.Fatal(err)
	}

}

// prettify replaces addresses and destinations using the supplied prettyPrintDict
func prettify(t *testing.T, logDir string, combinedLogsFilename string, prettyPrintDict map[string]string) {
	input, err := ioutil.ReadFile(path.Join(logDir, combinedLogsFilename))
	if err != nil {
		t.Fatal(err)
	}

	output := string(input)
	for key := range prettyPrintDict {
		output = strings.Replace(string(output), key, prettyPrintDict[key], -1)
	}

	if err = ioutil.WriteFile(path.Join(logDir, combinedLogsFilename)+"_pretty", []byte(output), 0666); err != nil {
		t.Fatal(err)
	}
}

func createVirtualChannelWithRetrievalProvider(c client.Client, retrievalProvider client.Client, prettyPrintDict map[string]string, i int) {
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
	prettyPrintDict[c.CreateVirtualChannel(withRetrievalProvider).ChannelId.String()] = "V" + fmt.Sprint(i)
}
