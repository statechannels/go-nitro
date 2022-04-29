package client_test

import (
	"bytes"
	"math/big"
	"math/rand"
	"testing"

	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	td "github.com/statechannels/go-nitro/internal/testdata"
	"github.com/statechannels/go-nitro/protocols/virtualfund"
	"github.com/statechannels/go-nitro/types"
)

// TestMultiPartyVirtualFundIntegration tests the scenario where Alice creates virtual channels with Bob and Brian using Irene as the intermediary.
func TestMultiPartyVirtualFundIntegration(t *testing.T) {

	logDestination := &bytes.Buffer{}
	t.Cleanup(flushToFileCleanupFn(logDestination, "virtualfund_multiparty_client_test.log"))

	chain := chainservice.NewMockChain()
	broker := messageservice.NewBroker()

	clientAlice, _ := setupClient(alice.PrivateKey, chain, broker, logDestination, 0)
	clientBob, _ := setupClient(bob.PrivateKey, chain, broker, logDestination, 0)
	clientBrian, _ := setupClient(brian.PrivateKey, chain, broker, logDestination, 0)
	clientIrene, _ := setupClient(irene.PrivateKey, chain, broker, logDestination, 0)

	directlyFundALedgerChannel(t, clientAlice, clientIrene)
	directlyFundALedgerChannel(t, clientIrene, clientBob)
	directlyFundALedgerChannel(t, clientIrene, clientBrian)
	withBobRequest := virtualfund.ObjectiveRequest{
		MyAddress:    alice.Address(),
		CounterParty: bob.Address(),
		Intermediary: irene.Address(),
		Outcome: td.Outcomes.Create(
			alice.Address(),
			bob.Address(),
			1,
			1,
		),
		AppDefinition:     types.Address{},
		AppData:           types.Bytes{},
		ChallengeDuration: big.NewInt(0),
		Nonce:             rand.Int63(),
	}
	withBrianRequest := virtualfund.ObjectiveRequest{
		MyAddress:    alice.Address(),
		CounterParty: brian.Address(),
		Intermediary: irene.Address(),
		Outcome: td.Outcomes.Create(
			alice.Address(),
			brian.Address(),
			1,
			1,
		),
		AppDefinition:     types.Address{},
		AppData:           types.Bytes{},
		ChallengeDuration: big.NewInt(0),
		Nonce:             rand.Int63(),
	}
	id := clientAlice.CreateVirtualChannel(withBobRequest)
	id2 := clientAlice.CreateVirtualChannel(withBrianRequest)

	waitTimeForCompletedObjectiveIds(t, &clientBob, defaultTimeout, id)
	waitTimeForCompletedObjectiveIds(t, &clientBrian, defaultTimeout, id2)

	waitTimeForCompletedObjectiveIds(t, &clientAlice, defaultTimeout, id, id2)
	waitTimeForCompletedObjectiveIds(t, &clientIrene, defaultTimeout, id, id2)

}
