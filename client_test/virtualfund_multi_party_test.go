package client_test

import (
	"testing"

	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	td "github.com/statechannels/go-nitro/internal/testdata"
	"github.com/statechannels/go-nitro/types"
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
	defer closeClient(t, &clientAlice)
	clientBob, _ := setupClient(bob.PrivateKey, chainServiceBo, broker, logDestination, 0)
	defer closeClient(t, &clientBob)
	clientBrian, _ := setupClient(brian.PrivateKey, chainServiceBr, broker, logDestination, 0)
	defer closeClient(t, &clientBrian)
	clientIrene, _ := setupClient(irene.PrivateKey, chainServiceI, broker, logDestination, 0)
	defer closeClient(t, &clientIrene)

	directlyFundALedgerChannel(t, clientAlice, clientIrene, types.Address{})
	directlyFundALedgerChannel(t, clientIrene, clientBob, types.Address{})
	directlyFundALedgerChannel(t, clientIrene, clientBrian, types.Address{})

	id := clientAlice.CreateVirtualPaymentChannel(
		[]types.Address{irene.Address()},
		bob.Address(),
		0,
		td.Outcomes.Create(
			alice.Address(),
			bob.Address(),
			1,
			1,
			types.Address{},
		)).Id

	id2 := clientAlice.CreateVirtualPaymentChannel(
		[]types.Address{irene.Address()},
		brian.Address(),
		0,
		td.Outcomes.Create(
			alice.Address(),
			brian.Address(),
			1,
			1,
			types.Address{},
		)).Id

	waitTimeForCompletedObjectiveIds(t, &clientBob, defaultTimeout, id)
	waitTimeForCompletedObjectiveIds(t, &clientBrian, defaultTimeout, id2)

	waitTimeForCompletedObjectiveIds(t, &clientAlice, defaultTimeout, id, id2)
	waitTimeForCompletedObjectiveIds(t, &clientIrene, defaultTimeout, id, id2)

}
