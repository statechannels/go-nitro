// Package client_test contains helpers and integration tests for go-nitro clients
package client_test // import "github.com/statechannels/go-nitro/client_test"

import (
	"io"
	"testing"
	"time"

	"github.com/statechannels/go-nitro/client"
	"github.com/statechannels/go-nitro/client/engine"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	simpletcp "github.com/statechannels/go-nitro/client/engine/messageservice/simple-tcp"
	"github.com/statechannels/go-nitro/client/engine/store"
	"github.com/statechannels/go-nitro/crypto"
	"github.com/statechannels/go-nitro/types"
)

// setupClientWithSimpleTCP is a helper function that contructs a client and returns the new client and its store.
func setupClientWithSimpleTCP(pk []byte, chain *chainservice.MockChainService, peers map[types.Address]string, logDestination io.Writer, meanMessageDelay time.Duration) (client.Client, *simpletcp.SimpleTCPMessageService, store.Store) {
	myAddress := crypto.GetAddressFromSecretKeyBytes(pk)
	messageservice := simpletcp.NewSimpleTCPMessageService(peers[myAddress], peers)
	store := store.NewMemStore(pk)
	return client.New(messageservice, chain, store, logDestination, &engine.PermissivePolicy{}, nil), messageservice, store
}

func TestSimpleTCPMessageService(t *testing.T) {

	// Setup logging
	logFile := "test_direct_fund_with_simple_tcp.log"
	truncateLog(logFile)
	logDestination := newLogWriter(logFile)

	chain := chainservice.NewMockChain()
	chainServiceA := chainservice.NewMockChainService(chain, alice.Address())
	chainServiceB := chainservice.NewMockChainService(chain, bob.Address())

	peers := map[types.Address]string{
		alice.Address(): "localhost:3005",
		bob.Address():   "localhost:3006",
	}
	clientA, msA, _ := setupClientWithSimpleTCP(alice.PrivateKey, chainServiceA, peers, logDestination, 0)
	clientB, msB, _ := setupClientWithSimpleTCP(bob.PrivateKey, chainServiceB, peers, logDestination, 0)
	defer msA.Close()
	defer msB.Close()
	directlyFundALedgerChannel(t, clientA, clientB)

}
