package client_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/statechannels/go-nitro/client"
	"github.com/statechannels/go-nitro/client/engine"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	p2pms "github.com/statechannels/go-nitro/client/engine/messageservice/p2p-message-service"
	"github.com/statechannels/go-nitro/client/engine/store"
	"github.com/statechannels/go-nitro/internal/testdata"
	"github.com/statechannels/go-nitro/rpc"
	"github.com/statechannels/go-nitro/rpc/transport"
	natstrans "github.com/statechannels/go-nitro/rpc/transport/nats"
	"github.com/statechannels/go-nitro/rpc/transport/wss"
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

	const connectionType transport.ConnectionType = "ws"
	rpcClientA, msgA, cleanupFnA := setupNitroNodeWithRPCClient(alice.PrivateKey, 3005, 4005, chainServiceA, logDestination, connectionType)
	rpcClientB, msgB, cleanupFnB := setupNitroNodeWithRPCClient(bob.PrivateKey, 3006, 4006, chainServiceB, logDestination, connectionType)
	rpcClientI, msgI, cleanupFnC := setupNitroNodeWithRPCClient(irene.PrivateKey, 3007, 4007, chainServiceI, logDestination, connectionType)

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

	rpcClientA.WaitForObjectiveCompletion(res.Id)
	rpcClientB.WaitForObjectiveCompletion(bobResponse.Id)
	rpcClientI.WaitForObjectiveCompletion(res.Id, bobResponse.Id)

	vRes := rpcClientA.CreateVirtual(
		[]types.Address{irene.Address()},
		bob.Address(),
		100,
		testdata.Outcomes.Create(alice.Address(), bob.Address(), 100, 100, types.Address{}))

	assert.Regexp(t, "VirtualFund.0x.*", vRes.Id)

	rpcClientA.WaitForObjectiveCompletion(vRes.Id)
	rpcClientB.WaitForObjectiveCompletion(vRes.Id)
	rpcClientI.WaitForObjectiveCompletion(vRes.Id)

	rpcClientA.Pay(vRes.ChannelId, 1)

	closeVId := rpcClientA.CloseVirtual(vRes.ChannelId)
	rpcClientA.WaitForObjectiveCompletion(closeVId)
	rpcClientB.WaitForObjectiveCompletion(closeVId)
	rpcClientI.WaitForObjectiveCompletion(closeVId)

	closeId := rpcClientA.CloseLedger(res.ChannelId)
	rpcClientA.WaitForObjectiveCompletion(closeId)
	rpcClientI.WaitForObjectiveCompletion(closeId)

}

// setupNitroNodeWithRPCClient is a helper function that spins up a Nitro Node RPC Server and returns an RPC client connected to it.
func setupNitroNodeWithRPCClient(
	pk []byte,
	msgPort int,
	rpcPort int,
	chain *chainservice.MockChainService,
	logDestination *os.File,
	connectionType transport.ConnectionType,
) (*rpc.RpcClient, *p2pms.P2PMessageService, func()) {
	messageservice := p2pms.NewMessageService("127.0.0.1", msgPort, pk)
	storeA := store.NewMemStore(pk)
	node := client.New(
		messageservice,
		chain,
		storeA,
		logDestination,
		&engine.PermissivePolicy{},
		nil)

	var serverConnection transport.Subscriber
	var clienConnection transport.Requester
	var err error
	switch connectionType {
	case "nats":
		serverConnection, err = natstrans.NewNatsConnectionAsServer(rpcPort)
		if err != nil {
			panic(err)
		}
		clienConnection, err = natstrans.NewNatsConnectionAsClient(serverConnection.Url())
		if err != nil {
			panic(err)
		}
	case "ws":
		serverConnection, err = wss.NewWebSocketConnectionAsServer(fmt.Sprint(rpcPort))
		if err != nil {
			panic(err)
		}
		clienConnection, err = wss.NewWebSocketConnectionAsClient(serverConnection.Url())
		if err != nil {
			panic(err)
		}
	default:
		err = fmt.Errorf("unknown connection type %v", connectionType)
		panic(err)
	}

	rpcServer, err := rpc.NewRpcServer(&node, createLogger(logDestination, node.Address.Hex(), "server"), serverConnection)
	if err != nil {
		panic(err)
	}
	rpcClient, err := rpc.NewRpcClient(rpcServer.Url(), alice.Address(), createLogger(logDestination, node.Address.Hex(), "client"), clienConnection)
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
