// Package client_test contains helpers and integration tests for go-nitro clients
package client_test // import "github.com/statechannels/go-nitro/client_test"

import (
	"io"
	"testing"

	"github.com/statechannels/go-nitro/client"
	"github.com/statechannels/go-nitro/client/engine"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	p2pms "github.com/statechannels/go-nitro/client/engine/messageservice/p2p-message-service"
	"github.com/statechannels/go-nitro/client/engine/store"
)

// setupClientWithP2PMessageService is a helper function that contructs a client and returns the new client and its store.
func setupClientWithP2PMessageService(pk []byte, port int, chain *chainservice.MockChainService, logDestination io.Writer) (client.Client, *p2pms.P2PMessageService) {

	messageservice := p2pms.NewMessageService("127.0.0.1", port, pk)
	storeA := store.NewMemStore(pk)
	return client.New(messageservice, chain, storeA, logDestination, &engine.PermissivePolicy{}, nil), messageservice
}

func TesP2PMessageService(t *testing.T) {

	// Setup logging
	logFile := "test_direct_fund_with_simple_tcp.log"
	truncateLog(logFile)
	logDestination := newLogWriter(logFile)

	chain := chainservice.NewMockChain()
	chainServiceA := chainservice.NewMockChainService(chain, alice.Address())
	chainServiceB := chainservice.NewMockChainService(chain, bob.Address())

	clientA, msA := setupClientWithP2PMessageService(alice.PrivateKey, 3005, chainServiceA, logDestination)
	clientB, msB := setupClientWithP2PMessageService(bob.PrivateKey, 3006, chainServiceB, logDestination)
	aPeer := p2pms.PeerInfo{Port: 3005, Id: msA.Id(), Address: alice.Address(), IpAddress: "127.0.0.1"}
	bPeer := p2pms.PeerInfo{Port: 3006, Id: msB.Id(), Address: bob.Address(), IpAddress: "127.0.0.1"}
	msA.AddPeers([]p2pms.PeerInfo{bPeer})
	msB.AddPeers([]p2pms.PeerInfo{aPeer})
	defer msA.Close()
	defer msB.Close()
	directlyFundALedgerChannel(t, clientA, clientB)

}
