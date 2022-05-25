package client_test

import (
	"io"
	"testing"

	"github.com/statechannels/go-nitro/client"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	simplep2p "github.com/statechannels/go-nitro/client/engine/messageservice/simple-p2p"

	"github.com/statechannels/go-nitro/client/engine/store"
	"github.com/statechannels/go-nitro/crypto"
	"github.com/statechannels/go-nitro/types"
)

// setupClientWithSimpleP2P is a helper function that contructs a client and returns the new client and its store.
func setupClientWithSimpleP2P(pk []byte, chain chainservice.MockChain, peers map[types.Address]simplep2p.PeerInfo, logDestination io.Writer) (client.Client, *simplep2p.P2PMessageService) {
	myAddress := crypto.GetAddressFromSecretKeyBytes(pk)
	chainservice := chainservice.NewSimpleChainService(&chain, myAddress)
	messageservice := simplep2p.NewP2PMessageService(peers[myAddress], peers)
	storeA := store.NewMemStore(pk)
	return client.New(messageservice, chainservice, storeA, logDestination), messageservice
}

func TestVirtualFundWithSimpleP2PMessageService(t *testing.T) {

	// Setup logging
	logFile := "p2p.log"
	truncateLog(logFile)
	logDestination := newLogWriter(logFile)

	chain := chainservice.NewMockChain()

	peers := map[types.Address]simplep2p.PeerInfo{
		alice.Address(): simplep2p.GeneratePeerInfo(alice.Address(), 3010),
		bob.Address():   simplep2p.GeneratePeerInfo(bob.Address(), 3011),
		irene.Address(): simplep2p.GeneratePeerInfo(irene.Address(), 3012),
	}

	clientA, msgA := setupClientWithSimpleP2P(alice.PrivateKey, chain, peers, logDestination)
	clientB, msgB := setupClientWithSimpleP2P(bob.PrivateKey, chain, peers, logDestination)
	clientI, msgI := setupClientWithSimpleP2P(irene.PrivateKey, chain, peers, logDestination)
	msgA.DialPeers()
	msgB.DialPeers()
	msgI.DialPeers()
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
