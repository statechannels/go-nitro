package client_test

import (
	"log"
	"math/big"
	"os"
	"testing"

	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/types"
)

// TestMultiPartyVirtualFundIntegration tests the scenario where Alice creates virtual channels with Bob and Brian using Irene as the intermediary.
func TestMultiPartyVirtualFundIntegration(t *testing.T) {
	t.Skip()

	// Set up logging
	logDestination, err := os.OpenFile("virtualfund_multiparty_client_test.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	// Reset log destination file
	err = logDestination.Truncate(0)
	if err != nil {
		log.Fatal(err)
	}

	chain := chainservice.NewMockChain([]types.Address{alice, bob, irene, brian})

	clientAlice, aliceMS := setupClient(aliceKey, chain, logDestination)
	clientBob, bobMS := setupClient(bobKey, chain, logDestination)
	clientBrian, brianMS := setupClient(brianKey, chain, logDestination)
	clientIrene, ireneMS := setupClient(ireneKey, chain, logDestination)

	connectMessageServices(aliceMS, bobMS, ireneMS, brianMS)

	directlyFundALedgerChannel(t, clientAlice, clientIrene)
	directlyFundALedgerChannel(t, clientIrene, clientBob)
	directlyFundALedgerChannel(t, clientIrene, clientBrian)

	id := clientAlice.CreateVirtualChannel(bob, irene, types.Address{}, types.Bytes{}, createVirtualOutcome(alice, bob), big.NewInt(0))
	id2 := clientAlice.CreateVirtualChannel(brian, irene, types.Address{}, types.Bytes{}, createVirtualOutcome(alice, brian), big.NewInt(0))

	waitTimeForCompletedObjectiveIds(t, &clientBob, defaultTimeout, id)
	waitTimeForCompletedObjectiveIds(t, &clientBrian, defaultTimeout, id2)

	waitTimeForCompletedObjectiveIds(t, &clientAlice, defaultTimeout, id, id2)
	waitTimeForCompletedObjectiveIds(t, &clientIrene, defaultTimeout, id, id2)

}
