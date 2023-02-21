package client_test

import (
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/statechannels/go-nitro/client"
	"github.com/statechannels/go-nitro/client/engine"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	p2pms "github.com/statechannels/go-nitro/client/engine/messageservice/p2p-message-service"
	"github.com/statechannels/go-nitro/client/engine/store"
	"github.com/statechannels/go-nitro/internal/testdata"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/rpc"
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

func TestRpc(t *testing.T) {
	logFile := "test_rpc_client.log"
	truncateLog(logFile)
	logDestination := newLogWriter(logFile)

	chain := chainservice.NewMockChain()
	chainServiceA := chainservice.NewMockChainService(chain, alice.Address())
	chainServiceB := chainservice.NewMockChainService(chain, bob.Address())
	chainServiceI := chainservice.NewMockChainService(chain, irene.Address())

	rpcClientA, msgA, cleanupFnA := setupNitroNodeWithRPCClient(alice.PrivateKey, 3005, 4005, chainServiceA, logDestination)
	rpcClientB, msgB, cleanupFnB := setupNitroNodeWithRPCClient(bob.PrivateKey, 3006, 4006, chainServiceB, logDestination)
	rpcClientI, msgI, cleanupFnC := setupNitroNodeWithRPCClient(irene.PrivateKey, 3007, 4007, chainServiceI, logDestination)

	peers := []p2pms.PeerInfo{
		{Id: msgA.Id(), IpAddress: "127.0.0.1", Port: 3005, Address: alice.Address()},
		{Id: msgB.Id(), IpAddress: "127.0.0.1", Port: 3006, Address: bob.Address()},
		{Id: msgI.Id(), IpAddress: "127.0.0.1", Port: 3007, Address: irene.Address()},
	}
	// Connect nitro P2P message services
	msgA.AddPeers(peers)
	msgB.AddPeers(peers)
	msgI.AddPeers(peers)

	defer cleanupFnA()
	defer cleanupFnB()
	defer cleanupFnC()

	res := rpcClientA.CreateLedger(irene.Address(), 100, testdata.Outcomes.Create(alice.Address(), irene.Address(), 100, 100, types.Address{}))
	bobResponse := rpcClientB.CreateLedger(irene.Address(), 100, testdata.Outcomes.Create(bob.Address(), irene.Address(), 100, 100, types.Address{}))

	// Quick sanity check that we're getting a valid objective id
	assert.Regexp(t, "DirectFunding.0x.*", res.Id)

	waitForObjectiveCompletion(t, rpcClientA, res.Id)
	waitForObjectiveCompletion(t, rpcClientB, bobResponse.Id)
	waitForObjectiveCompletion(t, rpcClientI, res.Id, bobResponse.Id)

	vRes := rpcClientA.CreateVirtual(
		[]types.Address{irene.Address()},
		bob.Address(),
		100,
		testdata.Outcomes.Create(alice.Address(), bob.Address(), 100, 100, types.Address{}))

	assert.Regexp(t, "VirtualFund.0x.*", vRes.Id)

	waitForObjectiveCompletion(t, rpcClientA, vRes.Id)
	waitForObjectiveCompletion(t, rpcClientB, vRes.Id)
	waitForObjectiveCompletion(t, rpcClientI, vRes.Id)

	rpcClientA.Pay(vRes.ChannelId, 1)

	closeVId := rpcClientA.CloseVirtual(vRes.ChannelId)
	waitForObjectiveCompletion(t, rpcClientA, closeVId)
	waitForObjectiveCompletion(t, rpcClientB, closeVId)
	waitForObjectiveCompletion(t, rpcClientI, closeVId)

	closeId := rpcClientA.CloseLedger(res.ChannelId)
	waitForObjectiveCompletion(t, rpcClientA, closeId)
	waitForObjectiveCompletion(t, rpcClientI, closeId)

}

// setupNitroNodeWithRPCClient is a helper function that spins up a Nitro Node RPC Server and returns an RPC client connected to it.
func setupNitroNodeWithRPCClient(
	pk []byte,
	msgPort int,
	rpcPort int,
	chain *chainservice.MockChainService,
	logDestination *os.File,
) (*rpc.RpcClient, *p2pms.P2PMessageService, func()) {
	chainId, err := chain.GetChainId()
	if err != nil {
		panic(err)
	}
	messageservice := p2pms.NewMessageService("127.0.0.1", msgPort, pk)
	storeA := store.NewMemStore(pk)
	node := client.New(
		messageservice,
		chain,
		storeA,
		logDestination,
		&engine.PermissivePolicy{},
		nil)
	rpcServer := rpc.NewRpcServer(
		&node,
		chainId,
		createLogger(logDestination, node.Address.Hex(), "server"),
		rpcPort)
	rpcClient, err := rpc.NewRpcClient(rpcServer.Url(), alice.Address(), chainId, createLogger(logDestination, node.Address.Hex(), "client"))
	if err != nil {
		panic(err)
	}
	cleanupFn := func() {
		messageservice.Close()
		rpcClient.Close()
		rpcServer.Close()
	}
	return rpcClient, messageservice, cleanupFn
}

func waitForObjectiveCompletion(t *testing.T, client *rpc.RpcClient, objectiveIds ...protocols.ObjectiveId) {
	err := client.WaitForObjectiveCompletion(objectiveIds...)
	if err != nil {
		t.Error(err)
	}
}
