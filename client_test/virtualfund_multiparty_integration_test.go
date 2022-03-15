package client_test

import (
	"math/big"
	"testing"

	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	"github.com/statechannels/go-nitro/types"
)

// TestMultiPartyVirtualFundIntegration tests the scenario where Alice creates virtual channels with Bob and Brian using Irene as the intermediary.
func TestMultiPartyVirtualFundIntegration(t *testing.T) {

	logFile := "virtualfund_multiparty_client_test.log"
	truncateLog(logFile)

	chain := chainservice.NewMockChain()
	broker := messageservice.NewBroker()

	clientAlice := setupClient(aliceKey, chain, broker, logFile)
	clientBob := setupClient(bobKey, chain, broker, logFile)
	clientBrian := setupClient(brianKey, chain, broker, logFile)
	clientIrene := setupClient(ireneKey, chain, broker, logFile)

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
