package client_test

import (
	"testing"

	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	td "github.com/statechannels/go-nitro/internal/testdata"
)

// TestVirtualFundMultiParty tests the scenario where Alice creates virtual channels with Bob and Brian using Irene as the intermediary.
func TestVirtualFundMultiParty(t *testing.T) {

	logFile := "test_virtual_fund_multi_party.log"
	truncateLog(logFile)
	logDestination := newLogWriter(logFile)

	chain := chainservice.NewMockChain()
	chainServiceA := chainservice.NewMockChainService(chain, alice.Address())
	chainServiceBo := chainservice.NewMockChainService(chain, bob.Address())
	chainServiceBr := chainservice.NewMockChainService(chain, brian.Address())
	chainServiceI := chainservice.NewMockChainService(chain, irene.Address())
	broker := messageservice.NewBroker()

	clientAlice, _ := setupClient(alice.PrivateKey, chainServiceA, broker, logDestination, 0)
	clientBob, _ := setupClient(bob.PrivateKey, chainServiceBo, broker, logDestination, 0)
	clientBrian, _ := setupClient(brian.PrivateKey, chainServiceBr, broker, logDestination, 0)
	clientIrene, _ := setupClient(irene.PrivateKey, chainServiceI, broker, logDestination, 0)

	directlyFundALedgerChannel(t, clientAlice, clientIrene)
	directlyFundALedgerChannel(t, clientIrene, clientBob)
	directlyFundALedgerChannel(t, clientIrene, clientBrian)

	id := clientAlice.CreateVirtualPaymentChannel(irene.Address(), bob.Address(), 0, td.Outcomes.Create(
		alice.Address(),
		bob.Address(),
		1,
		1,
	)).Id

	id2 := clientAlice.CreateVirtualPaymentChannel(irene.Address(), brian.Address(), 0, td.Outcomes.Create(
		alice.Address(),
		brian.Address(),
		1,
		1,
	)).Id

	waitTimeForCompletedObjectiveIds(t, &clientBob, defaultTimeout, id)
	waitTimeForCompletedObjectiveIds(t, &clientBrian, defaultTimeout, id2)

	waitTimeForCompletedObjectiveIds(t, &clientAlice, defaultTimeout, id, id2)
	waitTimeForCompletedObjectiveIds(t, &clientIrene, defaultTimeout, id, id2)

}
