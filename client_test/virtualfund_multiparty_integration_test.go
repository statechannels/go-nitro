package client_test

import (
	"math/big"
	"math/rand"
	"testing"

	"github.com/statechannels/go-nitro/client"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/protocols/virtualfund"
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

	bobIds := createVirtualChannels(&clientAlice, bob, irene, 5)
	brianIds := createVirtualChannels(&clientAlice, brian, irene, 5)
	allIds := append(bobIds, brianIds...)
	waitTimeForCompletedObjectiveIds(t, &clientBob, defaultTimeout, bobIds...)
	waitTimeForCompletedObjectiveIds(t, &clientBrian, defaultTimeout, brianIds...)

	waitTimeForCompletedObjectiveIds(t, &clientAlice, defaultTimeout, allIds...)
	waitTimeForCompletedObjectiveIds(t, &clientIrene, defaultTimeout, allIds...)
}

// createVirtualChannels is a helper function to create virtual channels in bulk
func createVirtualChannels(client *client.Client, counterParty types.Address, intermediary types.Address, amount uint) []protocols.ObjectiveId {
	request := virtualfund.ObjectiveRequest{
		MyAddress:         *client.Address,
		CounterParty:      counterParty,
		Intermediary:      intermediary,
		Outcome:           createVirtualOutcome(*client.Address, counterParty),
		AppDefinition:     types.Address{},
		AppData:           types.Bytes{},
		ChallengeDuration: big.NewInt(0),
		Nonce:             rand.Int63(),
	}

	ids := []protocols.ObjectiveId{}
	for i := 0; i < int(amount); i++ {
		// Generate a unique nonce for each request
		request.Nonce = rand.Int63()

		id := client.CreateVirtualChannel(request)
		ids = append(ids, id)

	}
	return ids

}
