// Package client_test contains helpers and integration tests for go-nitro clients
package client_test // import "github.com/statechannels/go-nitro/client_test"

import (
	"math/big"
	"testing"

	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	"github.com/statechannels/go-nitro/protocols"
)

func TestVirtualDefundIntegration(t *testing.T) {

	// Setup logging
	logFile := "test_virtual_defund.log"
	truncateLog(logFile)
	logDestination := newLogWriter(logFile)

	chain := chainservice.NewMockChain()
	broker := messageservice.NewBroker()

	clientA, _ := setupClient(alice.PrivateKey, chain, broker, logDestination, 0)
	clientB, _ := setupClient(bob.PrivateKey, chain, broker, logDestination, 0)
	clientI, _ := setupClient(irene.PrivateKey, chain, broker, logDestination, 0)

	// TODO: This test only supports defunding 1 virtual channel due to https://github.com/statechannels/go-nitro/issues/637
	cIds := openVirtualChannels(t, clientA, clientB, clientI, 1)

	paidToBob := big.NewInt(1)

	ids := make([]protocols.ObjectiveId, len(cIds))
	for i := 0; i < len(cIds); i++ {
		ids[i] = clientA.CloseVirtualChannel(cIds[i], paidToBob)

	}
	waitTimeForCompletedObjectiveIds(t, &clientA, defaultTimeout, ids...)
	waitTimeForCompletedObjectiveIds(t, &clientB, defaultTimeout, ids...)
	waitTimeForCompletedObjectiveIds(t, &clientI, defaultTimeout, ids...)

}
