package client_test

import (
	"fmt"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/rs/zerolog"
	"github.com/statechannels/go-nitro/channel/state/outcome"
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

func simpleOutcome(a, b types.Address, aBalance, bBalance uint) outcome.Exit {
	return testdata.Outcomes.Create(a, b, aBalance, bBalance, types.Address{})
}

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
	bobLedgerNotifs := rpcClientB.LedgerChannelUpdatesChan(bobResponse.ChannelId)

	// Quick sanity check that we're getting a valid objective id
	assert.Regexp(t, "DirectFunding.0x.*", res.Id)

	<-rpcClientA.ObjectiveCompleteChan(res.Id)
	<-rpcClientB.ObjectiveCompleteChan(bobResponse.Id)
	<-rpcClientI.ObjectiveCompleteChan(res.Id)
	<-rpcClientI.ObjectiveCompleteChan(bobResponse.Id)

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
	aliceVirtualNotifs := rpcClientA.PaymentChannelUpdatesChan(vRes.ChannelId)

	assert.Regexp(t, "VirtualFund.0x.*", vRes.Id)

	<-rpcClientA.ObjectiveCompleteChan(vRes.Id)
	<-rpcClientB.ObjectiveCompleteChan(vRes.Id)
	<-rpcClientI.ObjectiveCompleteChan(vRes.Id)

	expectedVirtual := expectedPaymentInfo(vRes.ChannelId, initialOutcome, query.Ready)
	aliceVirtual := rpcClientA.GetVirtualChannel(vRes.ChannelId)
	checkQueryInfo(t, expectedVirtual, aliceVirtual)
	bobVirtual := rpcClientB.GetVirtualChannel(vRes.ChannelId)
	checkQueryInfo(t, expectedVirtual, bobVirtual)
	ireneVirtual := rpcClientI.GetVirtualChannel(vRes.ChannelId)
	checkQueryInfo(t, expectedVirtual, ireneVirtual)

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

	expectedAliceLedgerNotifs := []query.LedgerChannelInfo{
		expectedLedgerInfo(res.ChannelId, simpleOutcome(ta.Alice.Address(), ta.Irene.Address(), 100, 100), query.Proposed),
		expectedLedgerInfo(res.ChannelId, simpleOutcome(ta.Alice.Address(), ta.Irene.Address(), 100, 100), query.Ready),
		expectedLedgerInfo(res.ChannelId, simpleOutcome(ta.Alice.Address(), ta.Irene.Address(), 0, 100), query.Ready),
		expectedLedgerInfo(res.ChannelId, simpleOutcome(ta.Alice.Address(), ta.Irene.Address(), 99, 101), query.Ready),
		expectedLedgerInfo(res.ChannelId, simpleOutcome(ta.Alice.Address(), ta.Irene.Address(), 99, 101), query.Closing),
		expectedLedgerInfo(res.ChannelId, simpleOutcome(ta.Alice.Address(), ta.Irene.Address(), 99, 101), query.Complete),
	}
	checkNotifications(t, expectedAliceLedgerNotifs, aliceLedgerNotifs, defaultTimeout)

	expectedBobLedgerNotifs := []query.LedgerChannelInfo{
		expectedLedgerInfo(bobResponse.ChannelId, simpleOutcome(ta.Bob.Address(), ta.Irene.Address(), 100, 100), query.Proposed),
		expectedLedgerInfo(bobResponse.ChannelId, simpleOutcome(ta.Bob.Address(), ta.Irene.Address(), 100, 100), query.Ready),
		expectedLedgerInfo(bobResponse.ChannelId, simpleOutcome(ta.Bob.Address(), ta.Irene.Address(), 100, 0), query.Ready),
		expectedLedgerInfo(bobResponse.ChannelId, simpleOutcome(ta.Bob.Address(), ta.Irene.Address(), 101, 99), query.Ready),
		expectedLedgerInfo(bobResponse.ChannelId, simpleOutcome(ta.Bob.Address(), ta.Irene.Address(), 101, 99), query.Closing),
		expectedLedgerInfo(bobResponse.ChannelId, simpleOutcome(ta.Bob.Address(), ta.Irene.Address(), 101, 99), query.Complete),
	}
	checkNotifications(t, expectedBobLedgerNotifs, bobLedgerNotifs, defaultTimeout)

	expectedVirtualNotifs := []query.PaymentChannelInfo{
		expectedPaymentInfo(vRes.ChannelId, simpleOutcome(ta.Alice.Address(), ta.Bob.Address(), 100, 0), query.Proposed),
		expectedPaymentInfo(vRes.ChannelId, simpleOutcome(ta.Alice.Address(), ta.Bob.Address(), 100, 0), query.Ready),
		expectedPaymentInfo(vRes.ChannelId, simpleOutcome(ta.Alice.Address(), ta.Bob.Address(), 99, 1), query.Ready),
		expectedPaymentInfo(vRes.ChannelId, simpleOutcome(ta.Alice.Address(), ta.Bob.Address(), 99, 1), query.Closing),
		expectedPaymentInfo(vRes.ChannelId, simpleOutcome(ta.Alice.Address(), ta.Bob.Address(), 99, 1), query.Complete),
	}
	checkNotifications(t, expectedVirtualNotifs, aliceVirtualNotifs, defaultTimeout)
	// TODO: Since we don't know exactly when bob receives and starts on the virtual channel
	// it's possible we could miss the first notification so for now we skip the check
	// checkNotifications(t, expectedVirtualNotifs, bobVirtualNotifs, defaultTimeout)
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

type channelInfo interface {
	query.LedgerChannelInfo | query.PaymentChannelInfo
}

func checkQueryInfo[T channelInfo](t *testing.T, expected T, fetched T) {
	if diff := cmp.Diff(expected, fetched, cmp.AllowUnexported(big.Int{})); diff != "" {
		panic(fmt.Errorf("Ledger diff mismatch (-want +got):\n%s", diff))
	}
}

// checkNotifications checks that the expected notifications are received on the notifChan.
// Due to the async nature of RPC notifications (and how quickly are clients communicate), the order of the notifications is not guaranteed.
// This function checks that all the expected notifications are received, but not in any particular order.
func checkNotifications[T channelInfo](t *testing.T, expected []T, notifChan <-chan T, timeout time.Duration) {
	fetched := make([]T, len(expected))
	for i := range expected {
		select {
		case info := <-notifChan:
			fetched[i] = info
		case <-time.After(timeout):
			t.Fatalf("Timed out waiting for notification.\n Fetched:%+v\n", fetched)
		}
	}
	for _, expected := range expected {
		found := false
		for _, fetched := range fetched {
			if (cmp.Equal(expected, fetched, cmp.AllowUnexported(big.Int{}))) {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("Expected notification not found: %v in fetched.\nFetched %+v\n", expected, fetched)
		}

	}
}
