package client_test

import (
	"testing"

	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/types"
)

func TestVirtualFundWithSimpleTCPMessageService(t *testing.T) {

	// Setup logging
	logFile := "test_virtual_fund_with_simple_tcp.log"
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

	ids := createVirtualChannels(clientA, bob.Address(), irene.Address(), 5)
	waitTimeForCompletedObjectiveIds(t, &clientA, defaultTimeout, ids...)
	waitTimeForCompletedObjectiveIds(t, &clientB, defaultTimeout, ids...)
	waitTimeForCompletedObjectiveIds(t, &clientI, defaultTimeout, ids...)
}
