package client_test

import (
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

	logFile := "virtualfund_multiparty_client_test.log"
	truncateLog(logFile)

	chain := chainservice.NewMockChain()
	broker := messageservice.NewBroker()

	clientAlice := setupClient(alice.PrivateKey, chain, broker, logFile, 0)
	clientBob := setupClient(bob.PrivateKey, chain, broker, logFile, 0)
	clientBrian := setupClient(brian.PrivateKey, chain, broker, logFile, 0)
	clientIrene := setupClient(irene.PrivateKey, chain, broker, logFile, 0)

	directlyFundALedgerChannel(t, clientAlice, clientIrene)
	directlyFundALedgerChannel(t, clientIrene, clientBob)
	directlyFundALedgerChannel(t, clientIrene, clientBrian)
	withBobRequest := virtualfund.ObjectiveRequest{
		MyAddress:    alice.Address,
		CounterParty: bob.Address,
		Intermediary: irene.Address,
		Outcome: td.Outcomes.Create(
			alice.Address,
			bob.Address,
			1,
			1,
		),
		AppDefinition:     types.Address{},
		AppData:           types.Bytes{},
		ChallengeDuration: big.NewInt(0),
		Nonce:             rand.Int63(),
	}
	withBrianRequest := virtualfund.ObjectiveRequest{
		MyAddress:    alice.Address,
		CounterParty: brian.Address,
		Intermediary: irene.Address,
		Outcome: td.Outcomes.Create(
			alice.Address,
			brian.Address,
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
