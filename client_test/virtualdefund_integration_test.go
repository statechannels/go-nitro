// Package client_test contains helpers and integration tests for go-nitro clients
package client_test // import "github.com/statechannels/go-nitro/client_test"

import (
	"bytes"
	"math/big"
	"testing"

	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
)

func TestVirtualDefundIntegration(t *testing.T) {

	// Setup logging
	logDestination := &bytes.Buffer{}
	t.Cleanup(flushToFileCleanupFn(logDestination, "virtualdefund_client_test.log"))

	chain := chainservice.NewMockChain()
	broker := messageservice.NewBroker()

	clientA, _ := setupClient(alice.PrivateKey, chain, broker, logDestination, 0)
	clientB, _ := setupClient(bob.PrivateKey, chain, broker, logDestination, 0)
	clientI, _ := setupClient(irene.PrivateKey, chain, broker, logDestination, 0)

	cIds := openVirtualChannels(t, clientA, clientB, clientI, 2)

	paidToBob := big.NewInt(1)

	id, id2 := clientA.CloseVirtualChannel(cIds[0], paidToBob), clientA.CloseVirtualChannel(cIds[1], paidToBob)

	waitTimeForCompletedObjectiveIds(t, &clientA, defaultTimeout, id, id2)
	waitTimeForCompletedObjectiveIds(t, &clientB, defaultTimeout, id, id2)
	waitTimeForCompletedObjectiveIds(t, &clientI, defaultTimeout, id, id2)

}
