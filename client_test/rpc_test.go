package client_test

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
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
	"github.com/statechannels/go-nitro/protocols/directfund"
	"github.com/statechannels/go-nitro/protocols/virtualfund"
	"github.com/statechannels/go-nitro/rpc"
	"github.com/statechannels/go-nitro/rpc/transport"
	natstrans "github.com/statechannels/go-nitro/rpc/transport/nats"
	"github.com/statechannels/go-nitro/rpc/transport/ws"
	"github.com/statechannels/go-nitro/types"
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
	executeNRpcTest(t, "nats", 2)
	executeNRpcTest(t, "nats", 3)
	executeNRpcTest(t, "nats", 4)
}

func TestRpcWithWebsockets(t *testing.T) {
	executeNRpcTest(t, "ws", 2)
	executeNRpcTest(t, "ws", 3)
	executeNRpcTest(t, "ws", 4)
}

func executeNRpcTest(t *testing.T, connectionType transport.TransportType, n int) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Test panicked: %v", r)
			t.FailNow()
		}
	}()

	if n < 2 {
		t.Errorf("n must be at least 2: alice and bob")
		return
	}

	//////////////////////
	// Setup
	//////////////////////

	t.Logf("Starting test with %d clients", n)
	logFile := fmt.Sprintf("test_%d_rpc_clients_over_%s.log", n, connectionType)
	logDestination := newLogWriter(logFile)
	defer logDestination.Close()

	chain := chainservice.NewMockChain()
	defer chain.Close()

	// create n actors
	actors := make([]ta.Actor, n)
	for i := 0; i < n; i++ {
		sk := `000000000000000000000000000000000000000000000000000000000000000` + strconv.Itoa(i+1)
		actors[i] = ta.Actor{
			PrivateKey: common.Hex2Bytes(sk),
		}
	}
	t.Logf("%d actors created", n)

	chainServices := make([]*chainservice.MockChainService, n)
	for i := 0; i < n; i++ {
		chainServices[i] = chainservice.NewMockChainService(chain, actors[i].Address())
	}

	clients := make([]*rpc.RpcClient, n)
	msgServices := make([]*p2pms.P2PMessageService, n)

	for i := 0; i < n; i++ {
		rpcClient, msg, cleanup := setupNitroNodeWithRPCClient(t, actors[i].PrivateKey, 3005+i, 4005+i, chainServices[i], logDestination, connectionType)
		clients[i] = rpcClient
		msgServices[i] = msg
		defer cleanup()
	}
	t.Logf("%d Clients created", n)

	waitForPeerInfoExchange(msgServices...)
	t.Logf("Peer exchange complete")

	// create n-1 ledger channels
	ledgerChannels := make([]directfund.ObjectiveResponse, n-1)
	for i := 0; i < n-1; i++ {
		outcome := simpleOutcome(actors[i].Address(), actors[i+1].Address(), 100, 100)
		ledgerChannels[i] = clients[i].CreateLedgerChannel(actors[i+1].Address(), 100, outcome)

		if !directfund.IsDirectFundObjective(ledgerChannels[i].Id) {
			t.Errorf("expected direct fund objective, got %s", ledgerChannels[i].Id)
		}
	}
	// wait for the ledger channels to be ready for each client
	for i, client := range clients {
		if i != 0 { // not alice
			<-client.ObjectiveCompleteChan(ledgerChannels[i-1].Id) // left channel
		}
		if i != n-1 { // not bob
			<-client.ObjectiveCompleteChan(ledgerChannels[i].Id) // right channel
		}
	}
	t.Log("Ledger channels created")

	// assert existence & reporting of expected ledger channels
	for i, client := range clients {
		if i != 0 {
			leftLC := ledgerChannels[i-1]
			expectedLeftLC := createLedgerInfo(leftLC.ChannelId, simpleOutcome(actors[i-1].Address(), actors[i].Address(), 100, 100), query.Open)
			checkQueryInfo(t, expectedLeftLC, client.GetLedgerChannel(leftLC.ChannelId))
		}
		if i != n-1 {
			rightLC := ledgerChannels[i]
			expectedRightLC := createLedgerInfo(rightLC.ChannelId, simpleOutcome(actors[i].Address(), actors[i+1].Address(), 100, 100), query.Open)
			checkQueryInfo(t, expectedRightLC, client.GetLedgerChannel(rightLC.ChannelId))
		}
	}

	//////////////////////////////////////////////////////////////////
	// create virtual channel, execute payment, close virtual channel
	//////////////////////////////////////////////////////////////////

	intermediaries := make([]types.Address, len(actors)-2)
	for i, actor := range actors[1 : len(actors)-1] {
		intermediaries[i] = actor.Address()
	}

	alice := actors[0]
	aliceClient := clients[0]
	bob := actors[n-1]
	bobClient := clients[n-1]
	aliceLedger := ledgerChannels[0]
	bobLedger := ledgerChannels[n-2]

	initialOutcome := simpleOutcome(actors[0].Address(), actors[n-1].Address(), 100, 0)

	vabCreateResponse := aliceClient.CreatePaymentChannel(
		intermediaries,
		bob.Address(),
		100,
		initialOutcome,
	)
	expectedVirtualChannel := createPaychInfo(
		vabCreateResponse.ChannelId,
		initialOutcome,
		query.Open,
	)

	// wait for the virtual channel to be ready, and
	// assert correct reporting from query api
	for i, client := range clients {
		<-client.ObjectiveCompleteChan(vabCreateResponse.Id)
		channelInfo := client.GetPaymentChannel(vabCreateResponse.ChannelId)
		checkQueryInfo(t, expectedVirtualChannel, channelInfo)
		if i != 0 {
			checkQueryInfoCollection(t, expectedVirtualChannel, 1,
				client.GetPaymentChannelsByLedger(ledgerChannels[i-1].ChannelId))
		}
		if i != n-1 {
			checkQueryInfoCollection(t, expectedVirtualChannel, 1,
				client.GetPaymentChannelsByLedger(ledgerChannels[i].ChannelId))
		}
	}

	if !virtualfund.IsVirtualFundObjective(vabCreateResponse.Id) {
		t.Errorf("expected virtual fund objective, got %s", vabCreateResponse.Id)
	}

	aliceClient.Pay(vabCreateResponse.ChannelId, 1)

	vabClosure := aliceClient.ClosePaymentChannel(vabCreateResponse.ChannelId)
	for _, client := range clients {
		<-client.ObjectiveCompleteChan(vabClosure)
	}

	laiClosure := aliceClient.CloseLedgerChannel(aliceLedger.ChannelId)
	<-aliceClient.ObjectiveCompleteChan(laiClosure)

	if n != 2 { // for n=2, alice and bob share a ledger, which should only be closed once.
		libClosure := bobClient.CloseLedgerChannel(bobLedger.ChannelId)
		<-bobClient.ObjectiveCompleteChan(libClosure)
	}

	//////////////////////////
	// perform wrap-up checks
	//////////////////////////

	for i, client := range clients {
		if i != 0 {
			leftLC := ledgerChannels[i-1]
			vcCount := len(client.GetPaymentChannelsByLedger(leftLC.ChannelId))
			if vcCount != 0 {
				t.Errorf("expected no virtual channels in ledger channel %s, got %d", leftLC.ChannelId, vcCount)
			}
		}
		if i != n-1 {
			rightLC := ledgerChannels[i]
			vcCount := len(client.GetPaymentChannelsByLedger(rightLC.ChannelId))
			if vcCount != 0 {
				t.Errorf("expected no virtual channels in ledger channel %s, got %d", rightLC.ChannelId, vcCount)
			}
		}
	}

	aliceLedgerNotifs := aliceClient.LedgerChannelUpdatesChan(ledgerChannels[0].ChannelId)
	expectedAliceLedgerNotifs := createLedgerStory(
		aliceLedger.ChannelId, alice.Address(), actors[1].Address(), // actor[1] is the first intermediary - can be Bob if n=2 (0-hop)
		[]channelStatusShorthand{
			{100, 100, query.Proposed},
			{100, 100, query.Open},
			{0, 100, query.Open},  // alice's balance forwarded to the guarantee for the virtual channel
			{99, 101, query.Open}, // returns to alice & actors[1] after closure
			{99, 101, query.Closing},
			{99, 101, query.Complete},
		},
	)
	checkNotifications(t, "aliceLedger", expectedAliceLedgerNotifs, []query.LedgerChannelInfo{}, aliceLedgerNotifs, defaultTimeout)

	bobLedgerNotifs := bobClient.LedgerChannelUpdatesChan(bobLedger.ChannelId)
	expectedBobLedgerNotifs := createLedgerStory(
		bobLedger.ChannelId, actors[n-2].Address(), bob.Address(),
		[]channelStatusShorthand{
			{100, 100, query.Proposed},
			{100, 100, query.Open},
			{0, 100, query.Open},
			{99, 101, query.Open},
			{99, 101, query.Complete},
		},
	)
	if n != 2 { // bob does not trigger a ledger-channel close if n=2 - alice does
		expectedBobLedgerNotifs = append(expectedBobLedgerNotifs,
			createLedgerInfo(bobLedger.ChannelId, simpleOutcome(actors[n-2].Address(), bob.Address(), 99, 101), query.Closing),
		)
	}
	checkNotifications(t, "bobLedger", expectedBobLedgerNotifs, []query.LedgerChannelInfo{}, bobLedgerNotifs, defaultTimeout)

	requiredVCNotifs := createPaychStory(
		vabCreateResponse.ChannelId, alice.Address(), bob.Address(),
		[]channelStatusShorthand{
			{100, 0, query.Proposed},
			{100, 0, query.Open},
			{99, 1, query.Complete},
		},
	)
	optionalVCNotifs := createPaychStory(
		vabCreateResponse.ChannelId, alice.Address(), bob.Address(),
		[]channelStatusShorthand{
			{99, 1, query.Closing},
			// TODO: Sometimes we see a closing notification with the original balance.
			// See https://github.com/statechannels/go-nitro/issues/1306
			{99, 1, query.Open},
			{100, 0, query.Closing},
		},
	)

	aliceVirtualNotifs := aliceClient.PaymentChannelUpdatesChan(vabCreateResponse.ChannelId)
	checkNotifications(t, "aliceVirtual", requiredVCNotifs, optionalVCNotifs, aliceVirtualNotifs, defaultTimeout)
	bobVirtualNotifs := bobClient.PaymentChannelUpdatesChan(vabCreateResponse.ChannelId)
	checkNotifications(t, "bobVirtual", requiredVCNotifs, optionalVCNotifs, bobVirtualNotifs, defaultTimeout)
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
	var clientConnection transport.Requester
	var err error
	switch connectionType {
	case "nats":
		serverConnection, err = natstrans.NewNatsTransportAsServer(rpcPort)
		if err != nil {
			panic(err)
		}
		clientConnection, err = natstrans.NewNatsTransportAsClient(serverConnection.Url())
		if err != nil {
			panic(err)
		}
	case "ws":
		serverConnection, err = ws.NewWebSocketTransportAsServer(fmt.Sprint(rpcPort))
		if err != nil {
			panic(err)
		}
		clientConnection, err = ws.NewWebSocketTransportAsClient(serverConnection.Url())
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
	rpcClient, err := rpc.NewRpcClient(rpcServer.Url(), createLogger(logDestination, node.Address.Hex(), "client"), clientConnection)
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
		t.Errorf("Channel query info diff mismatch (-want +got):\n%s", diff)
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
		t.Fatalf("did not find info %v in channel infos: %v", expected, fetched)
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
func checkNotifications[T channelInfo](t *testing.T, client string, required []T, optional []T, notifChan <-chan T, timeout time.Duration) {
	// This is map containing both required and optional notifications.
	// We use the json representation of the notification as the key and a boolean as the value.
	// The boolean value is true if the notification is required and false if it is optional.
	// When a notification is received it is removed from acceptableNotifications
	acceptableNotifications := make(map[string]bool)
	unexpectedNotifications := make(map[string]bool)
	logUnexpected := func() {
		for notif := range unexpectedNotifications {
			t.Logf("%s received unexpected notification: %v", client, notif)
		}
	}

	for _, r := range required {
		acceptableNotifications[marshalToJson(t, r)] = true
	}
	for _, o := range optional {
		acceptableNotifications[marshalToJson(t, o)] = false
	}

	for !areRequiredComplete(acceptableNotifications) {
		select {
		case info := <-notifChan:

			notifJSON := marshalToJson(t, info)
			t.Logf("%s received %v+", client, info)

			// Check that the notification is a required or optional one.
			_, isExpected := acceptableNotifications[notifJSON]

			if isExpected {
				// To signal we received a notification we delete it from the map
				delete(acceptableNotifications, notifJSON)
			} else {
				unexpectedNotifications[notifJSON] = true
			}

		case <-time.After(timeout):
			logUnexpected()
			t.Fatalf("%s timed out waiting for notification(s): \n%v", client, incompleteRequired(acceptableNotifications))
		}
	}
	if len(unexpectedNotifications) > 0 {
		logUnexpected()
		t.FailNow()
	}
}

// incompleteRequired returns a debug string listing
// required notifications that have not been received.
func incompleteRequired(notifs map[string]bool) string {
	required := ""
	for k, isRequired := range notifs {
		if isRequired {
			required += k + "\n"
		}
	}
	return required
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
