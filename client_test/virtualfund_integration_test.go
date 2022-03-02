package client_test

import (
	"log"
	"math/big"
	"os"
	"testing"

	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/types"
)

func TestVirtualFundIntegration(t *testing.T) {

	// Set up logging
	logDestination, err := os.OpenFile("virtualfund_client_test.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	// Reset log destination file
	err = logDestination.Truncate(0)
	if err != nil {
		log.Fatal(err)
	}

	chain := chainservice.NewMockChain([]types.Address{alice, bob, irene})

	clientA, messageserviceA := setupClient(aliceKey, chain, logDestination)
	clientB, messageserviceB := setupClient(bobKey, chain, logDestination)
	clientI, messageserviceI := setupClient(ireneKey, chain, logDestination)

	connectMessageServices(messageserviceA, messageserviceB, messageserviceI)

	directlyFundALedgerChannel(clientA, clientI)
	directlyFundALedgerChannel(clientI, clientB)

	outcome := createVirtualOutcome(alice, bob)

	id := clientA.CreateVirtualChannel(bob, irene, types.Address{}, types.Bytes{}, outcome, big.NewInt(0))

	waitForCompletedObjectiveId(id, &clientA)
	waitForCompletedObjectiveId(id, &clientB)
	waitForCompletedObjectiveId(id, &clientI)

}
