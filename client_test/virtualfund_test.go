package client_test

import (
	"math/big"
	"math/rand"
	"testing"

	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	"github.com/statechannels/go-nitro/internal/protocols/virtualfund"
	td "github.com/statechannels/go-nitro/internal/testdata"
	"github.com/statechannels/go-nitro/internal/types"
)

func TestVirtualFund(t *testing.T) {

	// Setup logging
	logFile := "test_virtual_fund.log"
	truncateLog(logFile)
	logDestination := newLogWriter(logFile)

	chain := chainservice.NewMockChain()
	broker := messageservice.NewBroker()

	clientA, _ := setupClient(alice.PrivateKey, chain, broker, logDestination, 0)
	clientB, _ := setupClient(bob.PrivateKey, chain, broker, logDestination, 0)
	clientI, _ := setupClient(irene.PrivateKey, chain, broker, logDestination, 0)

	directlyFundALedgerChannel(t, clientA, clientI)
	directlyFundALedgerChannel(t, clientI, clientB)

	outcome := td.Outcomes.Create(alice.Address(), bob.Address(), 1, 1)
	request := virtualfund.ObjectiveRequest{
		MyAddress:         alice.Address(),
		CounterParty:      bob.Address(),
		Intermediary:      irene.Address(),
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
