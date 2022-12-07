package pingpong

import (
	"testing"
	"time"

	"github.com/statechannels/go-nitro/app"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	"github.com/statechannels/go-nitro/internal/testactors"
	"github.com/stretchr/testify/require"
)

var alice = testactors.Alice
var bob = testactors.Bob

// Tests that the MarginApp complies with the app.App interface
func TestMarginAppType(t *testing.T) {
	var _ app.App = (*PingPongApp)(nil)
}

func TestFundingMethod(t *testing.T) {
	// Setup logging
	logFile := "test_app_ping.log"
	TruncateLog(logFile)
	logDestination := NewLogWriter(logFile)

	// Setup chain service
	sim, bindings, ethAccounts, err := chainservice.SetupSimulatedBackend(2)
	require.NoError(t, err)
	chainA, err := chainservice.NewSimulatedBackendChainService(sim, bindings, ethAccounts[0], logDestination)
	require.NoError(t, err)

	chainB, err := chainservice.NewSimulatedBackendChainService(sim, bindings, ethAccounts[1], logDestination)
	require.NoError(t, err)

	broker := messageservice.NewBroker()

	clientA, storeA, messageServiceA := SetupClient(alice.PrivateKey, chainA, broker, logDestination, 0)
	clientB, storeB, messageServiceB := SetupClient(bob.PrivateKey, chainB, broker, logDestination, 0)

	pingPongA := NewPingPongApp(clientA.GetEngine(), alice.Address())
	pingPongB := NewPingPongApp(clientB.GetEngine(), bob.Address())

	clientA.GetAppManager().RegisterApp(pingPongA)
	clientB.GetAppManager().RegisterApp(pingPongB)

	chId := directlyFundALedgerChannel(t, clientA, clientB)
	c, err := storeA.GetConsensusChannelById(chId)
	require.NoError(t, err)

	err = pingPongA.Ping(c)
	require.NoError(t, err)

	// Use sleep to wait for the message to be processed
	time.Sleep(1 * time.Second)

	_, _ = clientA, clientB
	_, _ = storeA, storeB
	_, _ = messageServiceA, messageServiceB
}
