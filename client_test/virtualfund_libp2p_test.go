package client_test

import (
	"testing"

	"github.com/statechannels/go-nitro/client/engine/chainservice"
	p2pms "github.com/statechannels/go-nitro/client/engine/messageservice/p2p-message-service"
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

	clientA, msgA := setupClientWithLibP2p(alice.PrivateKey, 3005, chainServiceA, logDestination)
	clientB, msgB := setupClientWithLibP2p(bob.PrivateKey, 3006, chainServiceB, logDestination)
	clientI, msgI := setupClientWithLibP2p(irene.PrivateKey, 3007, chainServiceI, logDestination)

	peers := []p2pms.PeerInfo{
		{Id: msgA.Id(), IpAddress: "127.0.0.1", Port: 3005, Address: alice.Address()},
		{Id: msgB.Id(), IpAddress: "127.0.0.1", Port: 3006, Address: bob.Address()},
		{Id: msgI.Id(), IpAddress: "127.0.0.1", Port: 3007, Address: irene.Address()},
	}

	msgA.AddPeers(peers)
	msgB.AddPeers(peers)
	msgI.AddPeers(peers)

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
