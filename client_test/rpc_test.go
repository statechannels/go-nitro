package client_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	p2pms "github.com/statechannels/go-nitro/client/engine/messageservice/p2p-message-service"
	"github.com/statechannels/go-nitro/client/rpc"
	"github.com/statechannels/go-nitro/internal/testdata"
	"github.com/statechannels/go-nitro/types"
	"github.com/stretchr/testify/assert"
)

func createLogger(logDestination *os.File, clientName, rpcRole string) zerolog.Logger {
	return zerolog.New(logDestination).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Str("client", clientName).
		Str("rpc", rpcRole).
		Str("scope", "").
		Logger()
}

func TestRpcClient(t *testing.T) {
	logDestination := newLogWriter("test_rpc_client.log")

	chain := chainservice.NewMockChain()
	chainServiceA := chainservice.NewMockChainService(chain, alice.Address())
	chainServiceB := chainservice.NewMockChainService(chain, bob.Address())
	chainId, err := chainServiceA.GetChainId()
	if err != nil {
		t.Fatal(err)
	}

	clientA, msgA := setupClientWithP2PMessageService(alice.PrivateKey, 3005, chainServiceA, createLogger(logDestination, "alice", ""))
	clientB, msgB := setupClientWithP2PMessageService(bob.PrivateKey, 3006, chainServiceB, createLogger(logDestination, "bob", ""))
	peers := []p2pms.PeerInfo{
		{Id: msgA.Id(), IpAddress: "127.0.0.1", Port: 3005, Address: alice.Address()},
		{Id: msgB.Id(), IpAddress: "127.0.0.1", Port: 3006, Address: bob.Address()},
	}
	// Connect nitro P2P message services
	msgA.AddPeers(peers)
	msgB.AddPeers(peers)

	defer msgA.Close()
	defer msgB.Close()

	rpcServerA := rpc.NewRpcServer(&clientA, chainId, createLogger(logDestination, "alice", "server"))
	rpcClientA := rpc.NewRpcClient(rpcServerA.Url(), alice.Address(), chainId, createLogger(logDestination, "alice", "client"))
	defer rpcServerA.Close()
	defer rpcClientA.Close()
	testOutcome := testdata.Outcomes.Create(alice.Address(), bob.Address(), 100, 100, types.Address{})

	res := rpcClientA.CreateLedger(bob.Address(), 100, testOutcome)

	// Quick sanity check that we're getting a valid objective id
	assert.Regexp(t, "DirectFunding.0x.*", res.Id)

	fmt.Printf("CreateLedger response: %v+\n", res)

	waitTimeForCompletedObjectiveIds(t, &clientA, defaultTimeout, res.Id)
	waitTimeForCompletedObjectiveIds(t, &clientB, defaultTimeout, res.Id)

	closeId := rpcClientA.CloseLedger(res.ChannelId)

	waitTimeForCompletedObjectiveIds(t, &clientA, defaultTimeout, closeId)
	waitTimeForCompletedObjectiveIds(t, &clientB, defaultTimeout, closeId)

}
