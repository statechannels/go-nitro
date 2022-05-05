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
	const numRetrievalClients = 1

	logDir := "../artifacts/vectorclock"

	os.RemoveAll(logDir)

	logFile := "largescale_client_test.log"

	truncateLog(logFile)
	logDestination := newLogWriter(logFile)

	chain := chainservice.NewMockChain()
	broker := messageservice.NewBroker()

	retrievalProvider, retrievalProviderStore := setupInstrumentedClient(bob.PrivateKey, chain, broker, logDestination, 0, logDir, "RP")
	paymentHub, _ := setupInstrumentedClient(irene.PrivateKey, chain, broker, logDestination, 0, logDir, "PH")

	directlyFundALedgerChannel(t, retrievalProvider, paymentHub)

	retrievalClients := make([]client.Client, numRetrievalClients)
	for i := range retrievalClients {
		secretKey, _ := nc.GeneratePrivateKeyAndAddress()
		retrievalClients[i], _ = setupInstrumentedClient(secretKey, chain, broker, logDestination, 0, logDir, "RC"+fmt.Sprint(i))
		directlyFundALedgerChannel(t, retrievalClients[i], paymentHub)
		go createVirtualChannelWithRetrievalProvider(retrievalClients[i], retrievalProvider)
	}

	<-time.After(1 * time.Second)

	retrievalProviderHubConnection, _ := retrievalProviderStore.GetConsensusChannel(*paymentHub.Address)

	finalOutcome, _ := json.Marshal(retrievalProviderHubConnection.SupportedSignedState().State().Outcome)

	logDestination.Write(finalOutcome)

	combineLogs(t, logDir, "shiviz.log")

}

func combineLogs(t *testing.T, logDir string, combinedLogsFilename string) {
	_, filename, _, _ := runtime.Caller(1)
	logDir = path.Join(path.Dir(filename), logDir)
	// NOTE: you may need to add GOPATH to PATH
	output, err := exec.Command("GoVector", "--log_type", "shiviz", "--log_dir", logDir, "--outfile", path.Join(logDir, combinedLogsFilename)).Output()
	t.Log(string(output), err)
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
	return c.CreateVirtualChannel(withRetrievalProvider)
}
