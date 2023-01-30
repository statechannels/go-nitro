package client_test

import (
	"fmt"
	"testing"

	"github.com/rs/zerolog"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	p2pms "github.com/statechannels/go-nitro/client/engine/messageservice/p2p-message-service"
	"github.com/statechannels/go-nitro/client/rpc"
	"github.com/statechannels/go-nitro/internal/testactors"
	"github.com/statechannels/go-nitro/internal/testdata"
	"github.com/statechannels/go-nitro/types"
	"github.com/stretchr/testify/assert"
)

func TestRpcClient(t *testing.T) {
	logDestination := newLogWriter("test_rpc_client.log")

	logger := zerolog.New(logDestination).
		Level(zerolog.InfoLevel).
		With().
		Timestamp().
		Str("client", "").
		Str("scope", "").
		Logger()

	chain := chainservice.NewMockChain()
	chainServiceA := chainservice.NewMockChainService(chain, alice.Address())
	chainServiceB := chainservice.NewMockChainService(chain, bob.Address())

	clientA, msgA := setupClientWithP2PMessageService(alice.PrivateKey, 3005, chainServiceA, logger)
	_, msgB := setupClientWithP2PMessageService(bob.PrivateKey, 3006, chainServiceB, logger)
	peers := []p2pms.PeerInfo{
		{Id: msgA.Id(), IpAddress: "127.0.0.1", Port: 3005, Address: alice.Address()},
		{Id: msgB.Id(), IpAddress: "127.0.0.1", Port: 3006, Address: bob.Address()},
	}
	// Connect nitro P2P message services
	msgA.AddPeers(peers)
	msgB.AddPeers(peers)

	defer msgA.Close()
	defer msgB.Close()

	alice := testactors.Alice
	bob := testactors.Bob

	rpcServerA := rpc.NewRpcServer(&clientA)
	rpcClientA := rpc.NewRpcClient(rpcServerA.Url(), alice.Address(), clientA.ChainId)
	defer rpcServerA.Close()
	defer rpcClientA.Close()
	testOutcome := testdata.Outcomes.Create(alice.Address(), bob.Address(), 100, 100, types.Address{})

	res := rpcClientA.CreateLedger(bob.Address(), 100, testOutcome)

	// Quick sanity check that we're getting a valid objective id
	assert.Regexp(t, "DirectFunding.0x.*", res.Id)

	fmt.Printf("CreateLedger response: %v+\n", res)

}
