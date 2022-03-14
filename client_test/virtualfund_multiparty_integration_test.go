package client_test

import (
	"math/big"
	"testing"

	"github.com/statechannels/go-nitro/client"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

// TestMultiPartyVirtualFundIntegration tests the scenario where Alice creates virtual channels with Bob and Brian using Irene as the intermediary.
func TestMultiPartyVirtualFundIntegration(t *testing.T) {

	logFile := "virtualfund_multiparty_client_test.log"
	truncateLog(logFile)

	chain := chainservice.NewMockChain([]types.Address{alice, bob, irene, brian})

	clientAlice, aliceMS := setupClient(aliceKey, chain, logFile)
	clientBob, bobMS := setupClient(bobKey, chain, logFile)
	clientBrian, brianMS := setupClient(brianKey, chain, logFile)
	clientIrene, ireneMS := setupClient(ireneKey, chain, logFile)

	connectMessageServices(aliceMS, bobMS, ireneMS, brianMS)

	directlyFundALedgerChannel(t, clientAlice, clientIrene)
	directlyFundALedgerChannel(t, clientIrene, clientBob)
	directlyFundALedgerChannel(t, clientIrene, clientBrian)

	bobIds := createVirtualChannels(&clientAlice, bob, irene, 3)
	brianIds := createVirtualChannels(&clientAlice, brian, irene, 3)
	allIds := append(bobIds, brianIds...)
	waitTimeForCompletedObjectiveIds(t, &clientBob, defaultTimeout, bobIds...)
	waitTimeForCompletedObjectiveIds(t, &clientBrian, defaultTimeout, brianIds...)

	waitTimeForCompletedObjectiveIds(t, &clientAlice, defaultTimeout, allIds...)
	waitTimeForCompletedObjectiveIds(t, &clientIrene, defaultTimeout, allIds...)
}

// createVirtualChannels is a helper function to create virtual channels in bulk
func createVirtualChannels(client *client.Client, counterParty types.Address, intermediary types.Address, amount uint) []protocols.ObjectiveId {
	ids := []protocols.ObjectiveId{}
	for i := 0; i < int(amount); i++ {
		id := client.CreateVirtualChannel(counterParty, irene, types.Address{}, types.Bytes{}, createVirtualOutcome(alice, bob), big.NewInt(0))
		ids = append(ids, id)

	}
	return ids

}
