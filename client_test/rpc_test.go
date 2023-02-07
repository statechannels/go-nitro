package client_test

import (
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
	chainServiceI := chainservice.NewMockChainService(chain, irene.Address())

	clientA, msgA := setupClientWithP2PMessageService(alice.PrivateKey, 3005, chainServiceA, createLogger(logDestination, "alice", ""))
	clientB, msgB := setupClientWithP2PMessageService(bob.PrivateKey, 3006, chainServiceB, createLogger(logDestination, "bob", ""))
	clientI, msgI := setupClientWithP2PMessageService(irene.PrivateKey, 3007, chainServiceI, createLogger(logDestination, "irene", ""))
	peers := []p2pms.PeerInfo{
		{Id: msgA.Id(), IpAddress: "127.0.0.1", Port: 3005, Address: alice.Address()},
		{Id: msgB.Id(), IpAddress: "127.0.0.1", Port: 3006, Address: bob.Address()},
		{Id: msgI.Id(), IpAddress: "127.0.0.1", Port: 3007, Address: irene.Address()},
	}
	// Connect nitro P2P message services
	msgA.AddPeers(peers)
	msgB.AddPeers(peers)
	msgI.AddPeers(peers)

	defer msgA.Close()
	defer msgB.Close()
	defer msgI.Close()

	rpcServerA := rpc.NewRpcServer(&clientA, chainId, createLogger(logDestination, "alice", "server"))
	rpcClientA := rpc.NewRpcClient(rpcServerA.Url(), alice.Address(), chainId, createLogger(logDestination, "alice", "client"))
	defer rpcServerA.Close()
	defer rpcClientA.Close()

	res := rpcClientA.CreateLedger(irene.Address(), 100, testdata.Outcomes.Create(alice.Address(), irene.Address(), 100, 100, types.Address{}))

	// Quick sanity check that we're getting a valid objective id
	assert.Regexp(t, "DirectFunding.0x.*", res.Id)
	waitTimeForCompletedObjectiveIds(t, &clientA, defaultTimeout, res.Id)
	waitTimeForCompletedObjectiveIds(t, &clientI, defaultTimeout, res.Id)

	bobResponse := clientB.CreateLedgerChannel(irene.Address(), 100, testdata.Outcomes.Create(bob.Address(), irene.Address(), 100, 100, types.Address{}))

	waitTimeForCompletedObjectiveIds(t, &clientB, defaultTimeout, bobResponse.Id)
	waitTimeForCompletedObjectiveIds(t, &clientI, defaultTimeout, bobResponse.Id)

	vRes := rpcClientA.CreateVirtual(
		[]types.Address{irene.Address()},
		bob.Address(),
		100,
		testdata.Outcomes.Create(alice.Address(), bob.Address(), 100, 100, types.Address{}))

	assert.Regexp(t, "VirtualFund.0x.*", vRes.Id)
	waitTimeForCompletedObjectiveIds(t, &clientA, defaultTimeout, vRes.Id)
	waitTimeForCompletedObjectiveIds(t, &clientB, defaultTimeout, vRes.Id)
	waitTimeForCompletedObjectiveIds(t, &clientI, defaultTimeout, vRes.Id)
	rpcClientA.Pay(vRes.ChannelId, 1)

	closeVId := rpcClientA.CloseVirtual(vRes.ChannelId)
	waitTimeForCompletedObjectiveIds(t, &clientA, defaultTimeout, closeVId)
	waitTimeForCompletedObjectiveIds(t, &clientB, defaultTimeout, closeVId)
	waitTimeForCompletedObjectiveIds(t, &clientI, defaultTimeout, closeVId)

	closeId := rpcClientA.CloseLedger(res.ChannelId)
	waitTimeForCompletedObjectiveIds(t, &clientA, defaultTimeout, closeId)
	waitTimeForCompletedObjectiveIds(t, &clientI, defaultTimeout, closeId)

}
