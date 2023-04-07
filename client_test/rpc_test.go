package client_test

import (
	"fmt"
	"math/big"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/rs/zerolog"
	"github.com/statechannels/go-nitro/client"
	"github.com/statechannels/go-nitro/client/engine"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	p2pms "github.com/statechannels/go-nitro/client/engine/messageservice/p2p-message-service"
	"github.com/statechannels/go-nitro/client/engine/store"
	"github.com/statechannels/go-nitro/client/query"
	"github.com/statechannels/go-nitro/crypto"
	ta "github.com/statechannels/go-nitro/internal/testactors"
	"github.com/statechannels/go-nitro/internal/testdata"
	"github.com/statechannels/go-nitro/rpc"
	"github.com/statechannels/go-nitro/rpc/transport"
	natstrans "github.com/statechannels/go-nitro/rpc/transport/nats"
	"github.com/statechannels/go-nitro/rpc/transport/ws"
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

func TestRpcWithNats(t *testing.T) {
	executeRpcTest(t, "nats")
}

func TestRpcWithWebsockets(t *testing.T) {
	executeRpcTest(t, "ws")
}

func executeRpcTest(t *testing.T, connectionType transport.TransportType) {
	logFile := "test_rpc_client.log"
	truncateLog(logFile)
	logDestination := newLogWriter(logFile)

	chain := chainservice.NewMockChain()
	chainServiceA := chainservice.NewMockChainService(chain, ta.Alice.Address())
	chainServiceB := chainservice.NewMockChainService(chain, ta.Bob.Address())
	chainServiceI := chainservice.NewMockChainService(chain, ta.Irene.Address())

	rpcClientA, _, cleanupFnA := setupNitroNodeWithRPCClient(t, ta.Alice.PrivateKey, 3005, 4005, chainServiceA, logDestination, connectionType)
	rpcClientB, _, cleanupFnB := setupNitroNodeWithRPCClient(t, ta.Bob.PrivateKey, 3006, 4006, chainServiceB, logDestination, connectionType)
	rpcClientI, _, cleanupFnC := setupNitroNodeWithRPCClient(t, ta.Irene.PrivateKey, 3007, 4007, chainServiceI, logDestination, connectionType)

	defer cleanupFnA()
	defer cleanupFnB()
	defer cleanupFnC()
	aliceLedgerOutcome := testdata.Outcomes.Create(ta.Alice.Address(), ta.Irene.Address(), 100, 100, types.Address{})
	bobLedgerOutcome := testdata.Outcomes.Create(ta.Bob.Address(), ta.Irene.Address(), 100, 100, types.Address{})
	res := rpcClientA.CreateLedger(ta.Irene.Address(), 100, aliceLedgerOutcome)
	bobResponse := rpcClientB.CreateLedger(ta.Irene.Address(), 100, bobLedgerOutcome)

	// Quick sanity check that we're getting a valid objective id
	assert.Regexp(t, "DirectFunding.0x.*", res.Id)

	<-rpcClientA.ObjectiveCompleteChan(res.Id)
	<-rpcClientB.ObjectiveCompleteChan(bobResponse.Id)
	<-rpcClientI.ObjectiveCompleteChan(res.Id)
	<-rpcClientI.ObjectiveCompleteChan(bobResponse.Id)

	aliceLedger := rpcClientA.GetLedgerChannel(res.ChannelId)
	expectedAliceLedger := expectedLedgerInfo(res.ChannelId, aliceLedgerOutcome, query.Ready)
	if diff := cmp.Diff(expectedAliceLedger, aliceLedger, cmp.AllowUnexported(big.Int{})); diff != "" {
		t.Fatalf("Ledger diff mismatch (-want +got):\n%s", diff)
	}

	bobLedger := rpcClientB.GetLedgerChannel(bobResponse.ChannelId)
	expectedBobLedger := expectedLedgerInfo(bobResponse.ChannelId, bobLedgerOutcome, query.Ready)
	if diff := cmp.Diff(expectedBobLedger, bobLedger, cmp.AllowUnexported(big.Int{})); diff != "" {
		t.Fatalf("Ledger diff mismatch (-want +got):\n%s", diff)
	}

	initialOutcome := testdata.Outcomes.Create(ta.Alice.Address(), ta.Bob.Address(), 100, 0, types.Address{})
	vRes := rpcClientA.CreateVirtual(
		[]types.Address{ta.Irene.Address()},
		ta.Bob.Address(),
		100,
		initialOutcome)

	assert.Regexp(t, "VirtualFund.0x.*", vRes.Id)

	<-rpcClientA.ObjectiveCompleteChan(vRes.Id)
	<-rpcClientB.ObjectiveCompleteChan(vRes.Id)
	<-rpcClientI.ObjectiveCompleteChan(vRes.Id)

	expectedVirtual := expectedPaymentInfo(vRes.ChannelId, initialOutcome, query.Ready)
	aliceVirtual := rpcClientA.GetVirtualChannel(vRes.ChannelId)
	if diff := cmp.Diff(expectedVirtual, aliceVirtual, cmp.AllowUnexported(big.Int{})); diff != "" {
		t.Fatalf("Virtual diff mismatch for alice (-want +got):\n%s", diff)
	}
	bobVirtual := rpcClientB.GetVirtualChannel(vRes.ChannelId)
	if diff := cmp.Diff(expectedVirtual, bobVirtual, cmp.AllowUnexported(big.Int{})); diff != "" {
		t.Fatalf("Virtual diff mismatch for bob (-want +got):\n%s", diff)
	}

	ireneVirtual := rpcClientI.GetVirtualChannel(vRes.ChannelId)
	if diff := cmp.Diff(expectedVirtual, ireneVirtual, cmp.AllowUnexported(big.Int{})); diff != "" {
		t.Fatalf("Virtual diff mismatch for irene (-want +got):\n%s", diff)
	}

	rpcClientA.Pay(vRes.ChannelId, 1)

	closeVId := rpcClientA.CloseVirtual(vRes.ChannelId)
	<-rpcClientA.ObjectiveCompleteChan(closeVId)
	<-rpcClientB.ObjectiveCompleteChan(closeVId)
	<-rpcClientI.ObjectiveCompleteChan(closeVId)

	closeId := rpcClientA.CloseLedger(res.ChannelId)
	<-rpcClientA.ObjectiveCompleteChan(closeId)
	<-rpcClientI.ObjectiveCompleteChan(closeId)

	closeIdB := rpcClientB.CloseLedger(bobResponse.ChannelId)
	<-rpcClientB.ObjectiveCompleteChan(closeIdB)
	<-rpcClientI.ObjectiveCompleteChan(closeIdB)
}

// setupNitroNodeWithRPCClient is a helper function that spins up a Nitro Node RPC Server and returns an RPC client connected to it.
func setupNitroNodeWithRPCClient(
	t *testing.T,
	pk []byte,
	msgPort int,
	rpcPort int,
	chain *chainservice.MockChainService,
	logDestination *os.File,
	connectionType transport.TransportType,
) (*rpc.RpcClient, *p2pms.P2PMessageService, func()) {
	messageservice := p2pms.NewMessageService("127.0.0.1",
		msgPort,
		crypto.GetAddressFromSecretKeyBytes(pk),
		pk)
	storeA := store.NewMemStore(pk)
	node := client.New(
		messageservice,
		chain,
		storeA,
		logDestination,
		&engine.PermissivePolicy{},
		nil)

	var serverConnection transport.Responder
	var clienConnection transport.Requester
	var err error
	switch connectionType {
	case "nats":
		serverConnection, err = natstrans.NewNatsTransportAsServer(rpcPort)
		if err != nil {
			panic(err)
		}
		clienConnection, err = natstrans.NewNatsTransportAsClient(serverConnection.Url())
		if err != nil {
			panic(err)
		}
	case "ws":
		serverConnection, err = ws.NewWebSocketTransportAsServer(fmt.Sprint(rpcPort))
		if err != nil {
			panic(err)
		}
		clienConnection, err = ws.NewWebSocketTransportAsClient(serverConnection.Url())
		if err != nil {
			panic(err)
		}
	default:
		err = fmt.Errorf("unknown connection type %v", connectionType)
		panic(err)
	}

	logger := createLogger(logDestination, node.Address.Hex(), "server")
	rpcServer, err := rpc.NewRpcServer(&node, &logger, serverConnection)
	if err != nil {
		panic(err)
	}
	rpcClient, err := rpc.NewRpcClient(rpcServer.Url(), ta.Alice.Address(), createLogger(logDestination, node.Address.Hex(), "client"), clienConnection)
	if err != nil {
		panic(err)
	}
	cleanupFn := func() {
		rpcClient.Close()
		rpcServer.Close()
	}
	return rpcClient, messageservice, cleanupFn
}
