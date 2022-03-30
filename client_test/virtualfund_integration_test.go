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

func TestVirtualFundIntegration(t *testing.T) {

	// Setup logging
	logDestination := &bytes.Buffer{}
	t.Cleanup(flushToFileCleanupFn(logDestination, "virtualfund_client_test.log"))

	chain := chainservice.NewMockChain()
	broker := messageservice.NewBroker()

	clientA := setupClient(alice.PrivateKey, chain, broker, logDestination, 0)
	clientB := setupClient(bob.PrivateKey, chain, broker, logDestination, 0)
	clientI := setupClient(irene.PrivateKey, chain, broker, logDestination, 0)

	directlyFundALedgerChannel(t, clientA, clientI)
	directlyFundALedgerChannel(t, clientI, clientB)

	outcome := td.Outcomes.Create(alice.Address, bob.Address, 1, 1)
	request := virtualfund.ObjectiveRequest{
		MyAddress:         alice.Address,
		CounterParty:      bob.Address,
		Intermediary:      irene.Address,
		Outcome:           outcome,
		AppDefinition:     types.Address{},
		AppData:           types.Bytes{},
		ChallengeDuration: big.NewInt(0),
		Nonce:             rand.Int63(),
	}
	id := clientA.CreateVirtualChannel(request)

	waitTimeForCompletedObjectiveIds(t, &clientA, defaultTimeout, id)
	waitTimeForCompletedObjectiveIds(t, &clientB, defaultTimeout, id)
	waitTimeForCompletedObjectiveIds(t, &clientI, defaultTimeout, id)

}
