package client_test

import (
	"io"
	"math/big"
	"testing"

	"github.com/statechannels/go-nitro/client"
	"github.com/statechannels/go-nitro/client/engine"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	p2pms "github.com/statechannels/go-nitro/client/engine/messageservice/p2p-message-service"
	"github.com/statechannels/go-nitro/client/engine/store"
	"github.com/statechannels/go-nitro/crypto"
	td "github.com/statechannels/go-nitro/internal/testdata"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

// setupClientWithP2PMessageService is a helper function that contructs a client and returns the new client and its store.
func setupClientWithP2PMessageService(pk []byte, port int, chain *chainservice.MockChainService, logDestination io.Writer) (client.Client, *p2pms.P2PMessageService) {

	messageservice := p2pms.NewMessageService(
		"127.0.0.1",
		port,
		crypto.GetAddressFromSecretKeyBytes(pk),
		generateMessageKey(pk))
	storeA := store.NewMemStore(pk)
	return client.New(messageservice, chain, storeA, logDestination, &engine.PermissivePolicy{}, nil), messageservice
}

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
	defer closeClient(t, &clientA)
	clientB, msgB := setupClientWithP2PMessageService(bob.PrivateKey, 3006, chainServiceB, logDestination)
	defer closeClient(t, &clientB)
	clientI, msgI := setupClientWithP2PMessageService(irene.PrivateKey, 3007, chainServiceI, logDestination)
	defer closeClient(t, &clientI)
	peers := []p2pms.PeerInfo{
		{Id: msgA.Id(), IpAddress: "127.0.0.1", Port: 3005, Address: alice.Address()},
		{Id: msgB.Id(), IpAddress: "127.0.0.1", Port: 3006, Address: bob.Address()},
		{Id: msgI.Id(), IpAddress: "127.0.0.1", Port: 3007, Address: irene.Address()},
	}

	msgA.AddPeers(peers)
	msgB.AddPeers(peers)
	msgI.AddPeers(peers)

	directlyFundALedgerChannel(t, clientA, clientI, types.Address{})
	directlyFundALedgerChannel(t, clientI, clientB, types.Address{})
	outcome := td.Outcomes.Create(alice.Address(), bob.Address(), 100, 100, types.Address{})
	r := clientA.CreateVirtualPaymentChannel(
		[]types.Address{irene.Address()},
		bob.Address(),
		0,
		outcome,
	)

	ids := []protocols.ObjectiveId{r.Id}

	waitTimeForCompletedObjectiveIds(t, &clientA, defaultTimeout, ids...)
	waitTimeForCompletedObjectiveIds(t, &clientB, defaultTimeout, ids...)
	waitTimeForCompletedObjectiveIds(t, &clientI, defaultTimeout, ids...)
	clientA.Pay(r.ChannelId, big.NewInt(5))

	expected := BasicVoucherInfo{big.NewInt(5), r.ChannelId}
	waitTimeForReceivedVoucher(t, &clientB, defaultTimeout, expected)

}
