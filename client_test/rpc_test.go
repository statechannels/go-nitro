package client_test

import (
	"encoding/json"
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

// func executeNRpcTest(t)

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
	waitForPeerInfoExchange(msgA, msgB, msgI)
	defer cleanupFnA()
	defer cleanupFnB()
	defer cleanupFnC()

	aliceLedgerOutcome := testdata.Outcomes.Create(ta.Alice.Address(), ta.Irene.Address(), 100, 100, types.Address{})
	bobLedgerOutcome := testdata.Outcomes.Create(ta.Bob.Address(), ta.Irene.Address(), 100, 100, types.Address{})

	laiRes := rpcClientA.CreateLedger(ta.Irene.Address(), 100, aliceLedgerOutcome)
	lbiRes := rpcClientB.CreateLedger(ta.Irene.Address(), 100, bobLedgerOutcome)

	aliceLedgerNotifs := rpcClientA.LedgerChannelUpdatesChan(laiRes.ChannelId)
	bobLedgerNotifs := rpcClientB.LedgerChannelUpdatesChan(lbiRes.ChannelId)

	// Quick sanity check that we're getting a valid objective id
	assert.Regexp(t, "DirectFunding.0x.*", laiRes.Id)

	<-rpcClientA.ObjectiveCompleteChan(laiRes.Id)
	<-rpcClientB.ObjectiveCompleteChan(lbiRes.Id)
	<-rpcClientI.ObjectiveCompleteChan(laiRes.Id)
	<-rpcClientI.ObjectiveCompleteChan(lbiRes.Id)

	expectedAliceLedger := createLedgerInfo(laiRes.ChannelId, aliceLedgerOutcome, query.Open)
	checkQueryInfo(t, expectedAliceLedger, rpcClientA.GetLedgerChannel(laiRes.ChannelId))
	checkQueryInfoCollection(t, expectedAliceLedger, 1, rpcClientA.GetAllLedgerChannels())

	expectedBobLedger := createLedgerInfo(lbiRes.ChannelId, bobLedgerOutcome, query.Open)
	checkQueryInfo(t, expectedBobLedger, rpcClientB.GetLedgerChannel(lbiRes.ChannelId))
	checkQueryInfoCollection(t, expectedBobLedger, 1, rpcClientB.GetAllLedgerChannels())

	initialOutcome := testdata.Outcomes.Create(ta.Alice.Address(), ta.Bob.Address(), 100, 0, types.Address{})
	vabRes := rpcClientA.CreateVirtual(
		[]types.Address{ta.Irene.Address()},
		ta.Bob.Address(),
		100,
		initialOutcome)
	aliceVirtualNotifs := rpcClientA.PaymentChannelUpdatesChan(vabRes.ChannelId)
	bobVirtualNotifs := rpcClientB.PaymentChannelUpdatesChan(vabRes.ChannelId)
	assert.Regexp(t, "VirtualFund.0x.*", vabRes.Id)

	<-rpcClientA.ObjectiveCompleteChan(vabRes.Id)
	<-rpcClientB.ObjectiveCompleteChan(vabRes.Id)
	<-rpcClientI.ObjectiveCompleteChan(vabRes.Id)

	expectedVirtual := expectedPaymentInfo(vabRes.ChannelId, initialOutcome, query.Open)
	aliceVab := rpcClientA.GetVirtualChannel(vabRes.ChannelId)
	checkQueryInfo(t, expectedVirtual, aliceVab)
	checkQueryInfoCollection(t, expectedVirtual, 1, rpcClientA.GetPaymentChannelsByLedger(laiRes.ChannelId))

	bobVab := rpcClientB.GetVirtualChannel(vabRes.ChannelId)
	checkQueryInfo(t, expectedVirtual, bobVab)
	checkQueryInfoCollection(t, expectedVirtual, 1, rpcClientB.GetPaymentChannelsByLedger(lbiRes.ChannelId))

	ireneVab := rpcClientI.GetVirtualChannel(vabRes.ChannelId)
	checkQueryInfo(t, expectedVirtual, ireneVab)
	checkQueryInfoCollection(t, expectedVirtual, 1, rpcClientI.GetPaymentChannelsByLedger(lbiRes.ChannelId))
	checkQueryInfoCollection(t, expectedVirtual, 1, rpcClientI.GetPaymentChannelsByLedger(laiRes.ChannelId))
	rpcClientA.Pay(vabRes.ChannelId, 1)

	closeVId := rpcClientA.CloseVirtual(vabRes.ChannelId)
	<-rpcClientA.ObjectiveCompleteChan(closeVId)
	<-rpcClientB.ObjectiveCompleteChan(closeVId)
	<-rpcClientI.ObjectiveCompleteChan(closeVId)

	closeId := rpcClientA.CloseLedger(laiRes.ChannelId)
	<-rpcClientA.ObjectiveCompleteChan(closeId)
	<-rpcClientI.ObjectiveCompleteChan(closeId)

	closeIdB := rpcClientB.CloseLedger(lbiRes.ChannelId)
	<-rpcClientB.ObjectiveCompleteChan(closeIdB)
	<-rpcClientI.ObjectiveCompleteChan(closeIdB)

	if len(rpcClientA.GetPaymentChannelsByLedger(laiRes.ChannelId)) != 0 {
		t.Error("Alice should not have any payment channels open")
	}
	if len(rpcClientB.GetPaymentChannelsByLedger(lbiRes.ChannelId)) != 0 {
		t.Error("Bob should not have any payment channels open")
	}
	if len(rpcClientI.GetPaymentChannelsByLedger(laiRes.ChannelId)) != 0 {
		t.Error("Irene should not have any payment channels open")
	}
	if len(rpcClientI.GetPaymentChannelsByLedger(lbiRes.ChannelId)) != 0 {
		t.Error("Irene should not have any payment channels open")
	}

	expectedAliceLedgerNotifs := createLedgerStory(
		laiRes.ChannelId, ta.Alice.Address(), ta.Irene.Address(),
		[]ledgerStatusShorthand{
			{100, 100, query.Proposed},
			{100, 100, query.Open},
			{0, 100, query.Open},
			{99, 101, query.Open},
			{99, 101, query.Closing},
			{99, 101, query.Complete},
		},
	)

	checkNotifications(t, expectedAliceLedgerNotifs, []query.LedgerChannelInfo{}, aliceLedgerNotifs, defaultTimeout)

	expectedBobLedgerNotifs := createLedgerStory(
		lbiRes.ChannelId, ta.Bob.Address(), ta.Irene.Address(),
		[]ledgerStatusShorthand{
			{100, 100, query.Proposed},
			{100, 100, query.Open},
			{100, 0, query.Open},
			{101, 99, query.Open},
			{101, 99, query.Closing},
			{101, 99, query.Complete},
		},
	)

	checkNotifications(t, expectedBobLedgerNotifs, []query.LedgerChannelInfo{}, bobLedgerNotifs, defaultTimeout)

	requiredVirtualNotifs := []query.PaymentChannelInfo{
		createPaychInfo(vabRes.ChannelId, simpleOutcome(ta.Alice.Address(), ta.Bob.Address(), 100, 0), query.Proposed),
		createPaychInfo(vabRes.ChannelId, simpleOutcome(ta.Alice.Address(), ta.Bob.Address(), 100, 0), query.Open),
		createPaychInfo(vabRes.ChannelId, simpleOutcome(ta.Alice.Address(), ta.Bob.Address(), 99, 1), query.Open),
		createPaychInfo(vabRes.ChannelId, simpleOutcome(ta.Alice.Address(), ta.Bob.Address(), 99, 1), query.Complete),
	}
	optionalVirtualNotifs := []query.PaymentChannelInfo{
		createPaychInfo(vabRes.ChannelId, simpleOutcome(ta.Alice.Address(), ta.Bob.Address(), 99, 1), query.Closing),
	}
	checkNotifications(t, requiredVirtualNotifs, optionalVirtualNotifs, aliceVirtualNotifs, defaultTimeout)

	checkNotifications(t, requiredVirtualNotifs, optionalVirtualNotifs, bobVirtualNotifs, defaultTimeout)
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
		true,
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
		panic(fmt.Errorf("Channel query info diff mismatch (-want +got):\n%s", diff))
	}
}

func checkQueryInfoCollection[T channelInfo](t *testing.T, expected T, expectedLength int, fetched []T) {
	if len(fetched) != expectedLength {
		t.Fatalf("expected %d channel infos, got %d", expectedLength, len(fetched))
	}
	found := false
	for _, fetched := range fetched {
		if cmp.Equal(expected, fetched, cmp.AllowUnexported(big.Int{})) {
			found = true
			break
		}
	}
	if !found {
		panic(fmt.Errorf("did not find info %v in channel infos: %v", expected, fetched))
	}
}

// marshalToJson marshals the given object to json and returns the string representation.
func marshalToJson[T channelInfo](t *testing.T, info T) string {
	jsonBytes, err := json.Marshal(info)
	if err != nil {
		t.Fatal(err)
	}
	return string(jsonBytes)
}

// checkNotifications checks that notifications are received on the notifChan.
//
// required specifies the notifications that must be received. checkNotifications will fail
// if any of these notifications are not received.
//
// optional specifies notifications that may be received. checkNotifications will not fail
// if any of these notifications are not received.
//
// If a notification is received that is neither in required or optional, checkNotifications will fail.
func checkNotifications[T channelInfo](t *testing.T, required []T, optional []T, notifChan <-chan T, timeout time.Duration) {
	// This is map containing both required and optional notifications.
	// We use the json representation of the notification as the key and a boolean as the value.
	// The boolean value is true if the notification is required and false if it is optional.
	// When a notification is received it is removed from acceptableNotifications
	acceptableNotifications := make(map[string]bool)

	for _, r := range required {
		acceptableNotifications[marshalToJson(t, r)] = true
	}
	for _, o := range optional {
		acceptableNotifications[marshalToJson(t, o)] = false
	}

	for !areRequiredComplete(acceptableNotifications) {
		select {
		case info := <-notifChan:

			js := marshalToJson(t, info)

			// Check that the notification is a required or optional one.
			_, found := acceptableNotifications[js]
			if !found {
				t.Fatalf("Received unexpected notification: %v", info)
			}
			// To signal we received a notification we delete it from the map
			delete(acceptableNotifications, js)

		case <-time.After(timeout):
			t.Fatalf("Timed out waiting for notification.\n")
		}
	}
}

// areRequiredComplete checks if all the required notifications have been received.
// It does this by checking that there are no members of the map that are true.
func areRequiredComplete(notifs map[string]bool) bool {
	for _, isRequired := range notifs {
		if isRequired {
			return false
		}
	}
	return true
}
