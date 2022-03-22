package client_test

import (
	"math/big"
	"math/rand"
	"testing"

	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	"github.com/statechannels/go-nitro/protocols/virtualfund"
	"github.com/statechannels/go-nitro/types"
)

func TestVirtualFundIntegration(t *testing.T) {

	// Set up logging
	logFile := "virtualfund_client_test.log"
	truncateLog(logFile)

	chain := chainservice.NewMockChain()
	broker := messageservice.NewBroker()

	clientA := setupClient(alice.PrivateKey, chain, broker, logFile)
	clientB := setupClient(bob.PrivateKey, chain, broker, logFile)
	clientI := setupClient(irene.PrivateKey, chain, broker, logFile)

	directlyFundALedgerChannel(t, clientA, clientI)
	directlyFundALedgerChannel(t, clientI, clientB)

	outcome := createVirtualOutcome(alice.Address, bob.Address)
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
