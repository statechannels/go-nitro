package client_test

import (
	"math/big"
	"testing"

	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	"github.com/statechannels/go-nitro/types"
)

func TestVirtualFundIntegration(t *testing.T) {

	// Set up logging
	logFile := "virtualfund_client_test.log"
	truncateLog(logFile)

	chain := chainservice.NewMockChain()
	broker := messageservice.NewBroker()

	clientA := setupClient(aliceKey, chain, broker, logFile)
	clientB := setupClient(bobKey, chain, broker, logFile)
	clientI := setupClient(ireneKey, chain, broker, logFile)

	directlyFundALedgerChannel(t, clientA, clientI)
	directlyFundALedgerChannel(t, clientI, clientB)

	outcome := createVirtualOutcome(alice, bob)

	id := clientA.CreateVirtualChannel(bob, irene, types.Address{}, types.Bytes{}, outcome, big.NewInt(0))

	waitTimeForCompletedObjectiveIds(t, &clientA, defaultTimeout, id)
	waitTimeForCompletedObjectiveIds(t, &clientB, defaultTimeout, id)
	waitTimeForCompletedObjectiveIds(t, &clientI, defaultTimeout, id)

}
