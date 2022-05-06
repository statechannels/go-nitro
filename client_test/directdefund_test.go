// Package client_test contains helpers and integration tests for go-nitro clients
package client_test // import "github.com/statechannels/go-nitro/client_test"

import (
	"testing"
	"time"

	"github.com/statechannels/go-nitro/client"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	"github.com/statechannels/go-nitro/types"
)

func directlyDefundALedgerChannel(t *testing.T, alpha client.Client, beta client.Client, channelId types.Destination) {

	id := alpha.CloseDirectChannel(channelId)
	waitTimeForCompletedObjectiveIds(t, &alpha, defaultTimeout, id)
	waitTimeForCompletedObjectiveIds(t, &beta, defaultTimeout, id)

}
func TestDirectDefund(t *testing.T) {

	// Setup logging
	logFile := "test_direct_defund.log"
	truncateLog(logFile)
	logDestination := newLogWriter(logFile)

	chain := chainservice.NewMockChain()
	broker := messageservice.NewBroker()

	clientA, _ := setupClient(alice.PrivateKey, chain, broker, logDestination, time.Millisecond*100)
	clientB, _ := setupClient(bob.PrivateKey, chain, broker, logDestination, time.Millisecond*100)

	numOfChannels := 1
	cIds := make([]types.Destination, numOfChannels)
	for i := 0; i < numOfChannels; i++ {
		cIds[i] = directlyFundALedgerChannel(t, clientA, clientB)

	}

	for i := 0; i < numOfChannels; i++ {
		directlyDefundALedgerChannel(t, clientA, clientB, cIds[i])

	}

}
