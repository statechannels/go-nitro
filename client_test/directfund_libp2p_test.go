// Package client_test contains helpers and integration tests for go-nitro clients
package client_test // import "github.com/statechannels/go-nitro/client_test"

import (
	"io"
	"testing"

	"github.com/statechannels/go-nitro/client"
	"github.com/statechannels/go-nitro/client/engine"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	libp2pms "github.com/statechannels/go-nitro/client/engine/messageservice/lib-p2p-message-service"

	"github.com/statechannels/go-nitro/client/engine/store"
)

// // PeerInfo represents a peer libp2p message service.
// type PeerInfo struct {
// 	Port    int64
// 	Id      peer.ID
// 	Address types.Address
// 	// HostName is the hostname of the peer. Either an IP address or a DNS name.
// 	HostName string
// }

// setupClientWithLibP2p is a helper function that contructs a client and returns the new client and its store.
func setupClientWithLibP2p(pk []byte, port int, chain *chainservice.MockChainService, logDestination io.Writer) (client.Client, *libp2pms.P2PMessageService) {

	messageservice := libp2pms.NewMessageService("127.0.0.1", port, pk)
	storeA := store.NewMemStore(pk)
	return client.New(messageservice, chain, storeA, logDestination, &engine.PermissivePolicy{}, nil), messageservice
}

func TesLibP2PMessageService(t *testing.T) {

	// Setup logging
	logFile := "test_direct_fund_with_simple_tcp.log"
	truncateLog(logFile)
	logDestination := newLogWriter(logFile)

	chain := chainservice.NewMockChain()
	chainServiceA := chainservice.NewMockChainService(chain, alice.Address())
	chainServiceB := chainservice.NewMockChainService(chain, bob.Address())

	clientA, msA := setupClientWithLibP2p(alice.PrivateKey, 3005, chainServiceA, logDestination)
	clientB, msB := setupClientWithLibP2p(bob.PrivateKey, 3006, chainServiceB, logDestination)
	aPeer := libp2pms.PeerInfo{Port: 3005, Id: msA.Id(), Address: alice.Address(), IpAddress: "127.0.0.1"}
	bPeer := libp2pms.PeerInfo{Port: 3006, Id: msB.Id(), Address: bob.Address(), IpAddress: "127.0.0.1"}
	msA.AddPeers([]libp2pms.PeerInfo{bPeer})
	msB.AddPeers([]libp2pms.PeerInfo{aPeer})
	defer msA.Close()
	defer msB.Close()
	directlyFundALedgerChannel(t, clientA, clientB)

}
