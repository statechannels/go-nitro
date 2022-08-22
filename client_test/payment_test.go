package client_test

import (
	"math/big"
	"math/rand"
	"testing"

	"github.com/statechannels/go-nitro/client/engine/chainservice"
	td "github.com/statechannels/go-nitro/internal/testdata"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/protocols/virtualfund"
	"github.com/statechannels/go-nitro/types"
)

func TestPayments(t *testing.T) {

	// Setup logging
	logFile := "test_payments.log"
	truncateLog(logFile)
	logDestination := newLogWriter(logFile)

	chain := chainservice.NewMockChain()
	chainServiceA := chainservice.NewMockChainService(chain, alice.Address())
	chainServiceB := chainservice.NewMockChainService(chain, bob.Address())
	chainServiceI := chainservice.NewMockChainService(chain, irene.Address())

	peers := map[types.Address]string{
		alice.Address(): "localhost:3005",
		bob.Address():   "localhost:3006",
		irene.Address(): "localhost:3007",
	}

	clientA, msgA := setupClientWithSimpleTCP(alice.PrivateKey, chainServiceA, peers, logDestination, 0)
	clientB, msgB := setupClientWithSimpleTCP(bob.PrivateKey, chainServiceB, peers, logDestination, 0)
	clientI, msgI := setupClientWithSimpleTCP(irene.PrivateKey, chainServiceI, peers, logDestination, 0)
	defer msgA.Close()
	defer msgB.Close()
	defer msgI.Close()

	directlyFundALedgerChannel(t, clientA, clientI)
	directlyFundALedgerChannel(t, clientI, clientB)
	outcome := td.Outcomes.Create(alice.Address(), bob.Address(), 100, 100)
	request := virtualfund.ObjectiveRequest{

		CounterParty:      bob.Address(),
		Intermediary:      irene.Address(),
		Outcome:           outcome,
		AppDefinition:     types.Address{},
		AppData:           types.Bytes{},
		ChallengeDuration: big.NewInt(0),
		Nonce:             rand.Int63(),
	}

	r := clientA.CreateVirtualChannel(request)

	ids := []protocols.ObjectiveId{r.Id}

	waitTimeForCompletedObjectiveIds(t, &clientA, defaultTimeout, ids...)
	waitTimeForCompletedObjectiveIds(t, &clientB, defaultTimeout, ids...)
	waitTimeForCompletedObjectiveIds(t, &clientI, defaultTimeout, ids...)
	clientA.Pay(r.ChannelId, big.NewInt(5))

	payment := <-clientB.ReceivedVouchers()

	if payment.Amount().Cmp(big.NewInt(5)) != 0 {

		t.Errorf("Expected payment amount to be 5, got %v", payment.Amount())
	}

}
