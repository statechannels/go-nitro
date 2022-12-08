package pingpong

import (
	"fmt"
	"testing"
	"time"

	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	"github.com/statechannels/go-nitro/internal/testactors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var alice = testactors.Alice
var bob = testactors.Bob

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
	_, _, _ = storeB, messageServiceA, messageServiceB

	pingPongA := NewPingPongApp(clientA.GetEngine(), alice.Address())
	pingPongB := NewPingPongApp(clientB.GetEngine(), bob.Address())

	clientA.GetAppManager().RegisterApp(pingPongA.App)
	clientB.GetAppManager().RegisterApp(pingPongB.App)

	chId := directlyFundALedgerChannel(t, clientA, clientB)
	c, err := storeA.GetConsensusChannelById(chId)
	require.NoError(t, err)

	done := make(chan int64, 16)
	err = pingPongA.Ping(c, done)
	require.NoError(t, err)

	timeout := time.After(10 * time.Millisecond)
	select {
	case <-timeout:
		t.Fatal("timeout")
	case rtt := <-done:
		assert.LessOrEqual(t, rtt, int64(10_000_000))
		fmt.Printf("Pong received, rtt: %d ns\n", rtt)
	}
}
