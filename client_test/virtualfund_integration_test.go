package client_test

import (
	"math/big"
	"testing"

	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/types"
)

func TestVirtualFundIntegration(t *testing.T) {

	// Set up logging
	logFile := "virtualfund_client_test.log"
	truncateLog(logFile)

	chain := chainservice.NewMockChain()

	clientA, messageserviceA := setupClient(aliceKey, chain, logFile)
	clientB, messageserviceB := setupClient(bobKey, chain, logFile)
	clientI, messageserviceI := setupClient(ireneKey, chain, logFile)

	connectMessageServices(messageserviceA, messageserviceB, messageserviceI)

	directlyFundALedgerChannel(t, clientA, clientI)
	directlyFundALedgerChannel(t, clientI, clientB)

	outcome := createVirtualOutcome(alice, bob)

	id := clientA.CreateVirtualChannel(bob, irene, types.Address{}, types.Bytes{}, outcome, big.NewInt(0))

	waitTimeForCompletedObjectiveIds(t, &clientA, defaultTimeout, id)
	waitTimeForCompletedObjectiveIds(t, &clientB, defaultTimeout, id)
	waitTimeForCompletedObjectiveIds(t, &clientI, defaultTimeout, id)

}
