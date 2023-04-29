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
	logDestination := newLogWriter(logFile)

	chain := chainservice.NewMockChain()
	chainServiceA := chainservice.NewMockChainService(chain, ta.Alice.Address())
	chainServiceB := chainservice.NewMockChainService(chain, ta.Bob.Address())
	chainServiceI := chainservice.NewMockChainService(chain, ta.Irene.Address())

	rpcClientA, msgA, cleanupFnA := setupNitroNodeWithRPCClient(t, ta.Alice.PrivateKey, 3005, 4005, chainServiceA, logDestination, connectionType)
	rpcClientB, msgB, cleanupFnB := setupNitroNodeWithRPCClient(t, ta.Bob.PrivateKey, 3006, 4006, chainServiceB, logDestination, connectionType)
	rpcClientI, msgI, cleanupFnC := setupNitroNodeWithRPCClient(t, ta.Irene.PrivateKey, 3007, 4007, chainServiceI, logDestination, connectionType)
	waitForPeerInfoExchange(2, msgA, msgB, msgI)
	defer cleanupFnA()
	defer cleanupFnB()
	defer cleanupFnC()
	aliceLedgerOutcome := testdata.Outcomes.Create(ta.Alice.Address(), ta.Irene.Address(), 100, 100, types.Address{})
	bobLedgerOutcome := testdata.Outcomes.Create(ta.Bob.Address(), ta.Irene.Address(), 100, 100, types.Address{})
	res := rpcClientA.CreateLedger(ta.Irene.Address(), 100, aliceLedgerOutcome)
	bobResponse := rpcClientB.CreateLedger(ta.Irene.Address(), 100, bobLedgerOutcome)

	aliceLedgerNotifs := rpcClientA.LedgerChannelUpdatesChan(res.ChannelId)
	bobledgerNotifs := rpcClientB.LedgerChannelUpdatesChan(bobResponse.ChannelId)

	// Quick sanity check that we're getting a valid objective id
	assert.Regexp(t, "DirectFunding.0x.*", res.Id)

	<-rpcClientA.ObjectiveCompleteChan(res.Id)
	<-rpcClientB.ObjectiveCompleteChan(bobResponse.Id)
	<-rpcClientI.ObjectiveCompleteChan(res.Id)
	<-rpcClientI.ObjectiveCompleteChan(bobResponse.Id)

	checkNotification(t, expectedLedgerInfo(res.ChannelId, aliceLedgerOutcome, query.Proposed), aliceLedgerNotifs)
	checkNotification(t, expectedLedgerInfo(bobResponse.ChannelId, bobLedgerOutcome, query.Proposed), bobledgerNotifs)

	expectedAliceLedger := expectedLedgerInfo(res.ChannelId, aliceLedgerOutcome, query.Ready)
	checkQueryInfo(t, expectedAliceLedger, rpcClientA.GetLedgerChannel(res.ChannelId))

	expectedBobLedger := expectedLedgerInfo(bobResponse.ChannelId, bobLedgerOutcome, query.Ready)
	checkQueryInfo(t, expectedBobLedger, rpcClientB.GetLedgerChannel(bobResponse.ChannelId))

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

	// TODO: For some reason we don't get these notifications until after objective completed is read
	checkNotification(t, expectedAliceLedger, aliceLedgerNotifs)
	checkNotification(t, expectedBobLedger, bobledgerNotifs)

	aliceVirtualNotifs := rpcClientA.PaymentChannelUpdatesChan(vRes.ChannelId)
	bobVirtualNotifs := rpcClientB.PaymentChannelUpdatesChan(vRes.ChannelId)
	ireneVirtualNotifs := rpcClientI.PaymentChannelUpdatesChan(vRes.ChannelId)

	expectedVirtual := expectedPaymentInfo(vRes.ChannelId, initialOutcome, query.Proposed)
	// TODO: For some reason we don't get these notifications until after objective completed is ready
	checkNotification(t, expectedVirtual, aliceVirtualNotifs)
	checkNotification(t, expectedVirtual, bobVirtualNotifs)
	checkNotification(t, expectedVirtual, ireneVirtualNotifs)

	expectedVirtual = expectedPaymentInfo(vRes.ChannelId, initialOutcome, query.Ready)
	aliceVirtual := rpcClientA.GetVirtualChannel(vRes.ChannelId)
	checkQueryInfo(t, expectedVirtual, aliceVirtual)

	bobVirtual := rpcClientB.GetVirtualChannel(vRes.ChannelId)
	checkQueryInfo(t, expectedVirtual, bobVirtual)

	ireneVirtual := rpcClientI.GetVirtualChannel(vRes.ChannelId)
	checkQueryInfo(t, expectedVirtual, ireneVirtual)

	// TODO: I expect an event here, but it doesn't seem to be happening
	// checkNotification(t, expectedVirtual, aliceVirtualNotifs)
	// checkNotification(t, expectedVirtual, bobVirtualNotifs)
	// checkNotification(t, expectedVirtual, ireneVirtualNotifs)
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
	finalOutcome := testdata.Outcomes.Create(ta.Alice.Address(), ta.Bob.Address(), 99, 1, types.Address{})

	checkNotification(t, expectedPaymentInfo(vRes.ChannelId, finalOutcome, query.Ready), aliceVirtualNotifs)
	checkNotification(t, expectedPaymentInfo(vRes.ChannelId, finalOutcome, query.Ready), bobVirtualNotifs)
	// TODO Irene doesn't seem to get a notification here. Is that ok?
	// checkNotification(t, expectedPaymentInfo(vRes.ChannelId, finalOutcome, query.Ready), ireneVirtualNotifs)
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
		pk,
		logDestination)
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

func checkQueryInfo[T query.ChannelInfo](t *testing.T, expected T, fetched T) {
	if diff := cmp.Diff(expected, fetched, cmp.AllowUnexported(big.Int{})); diff != "" {
		panic(fmt.Errorf("Ledger diff mismatch (-want +got):\n%s", diff))
	}
}

func checkNotification[T query.ChannelInfo](t *testing.T, expected T, notifChan <-chan T) {
	notif := <-notifChan
	fmt.Printf("notif: %v\n", notif)
	if diff := cmp.Diff(expected, notif, cmp.AllowUnexported(big.Int{})); diff != "" {
		panic(fmt.Errorf("Notification diff mismatch (-want +got):\n%s", diff))
	}
}
