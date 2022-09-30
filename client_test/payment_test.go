package client_test

import (
	"math/big"
	"testing"

	"github.com/statechannels/go-nitro/client/engine/chainservice"
	p2pms "github.com/statechannels/go-nitro/client/engine/messageservice/p2p-message-service"
	td "github.com/statechannels/go-nitro/internal/testdata"
	"github.com/statechannels/go-nitro/protocols"
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

	clientA, msgA := setupClientWithP2PMessageService(alice.PrivateKey, 3005, chainServiceA, logDestination)
	clientB, msgB := setupClientWithP2PMessageService(bob.PrivateKey, 3006, chainServiceB, logDestination)
	clientI, msgI := setupClientWithP2PMessageService(irene.PrivateKey, 3007, chainServiceI, logDestination)
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
	outcome := td.Outcomes.Create(alice.Address(), bob.Address(), 100, 100)
	r := clientA.CreateVirtualPaymentChannel(irene.Address(), bob.Address(), 0, outcome)

	ids := []protocols.ObjectiveId{r.Id}

	waitTimeForCompletedObjectiveIds(t, &clientA, defaultTimeout, ids...)
	waitTimeForCompletedObjectiveIds(t, &clientB, defaultTimeout, ids...)
	waitTimeForCompletedObjectiveIds(t, &clientI, defaultTimeout, ids...)
	clientA.Pay(r.ChannelId, big.NewInt(5))

	expected := BasicVoucherInfo{big.NewInt(5), r.ChannelId}
	waitTimeForReceivedVoucher(t, &clientB, defaultTimeout, expected)

}
