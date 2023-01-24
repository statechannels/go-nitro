package client

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/client"
	"github.com/statechannels/go-nitro/client/engine"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	p2pms "github.com/statechannels/go-nitro/client/engine/messageservice/p2p-message-service"
	"github.com/statechannels/go-nitro/client/engine/store"
	"github.com/statechannels/go-nitro/internal/testactors"
	"github.com/statechannels/go-nitro/internal/testdata"
	"github.com/statechannels/go-nitro/network"
	netproto "github.com/statechannels/go-nitro/network/protocol"
	"github.com/statechannels/go-nitro/network/protocol/parser"
	"github.com/statechannels/go-nitro/network/serde"
	nats2 "github.com/statechannels/go-nitro/network/transport/nats"
	"github.com/statechannels/go-nitro/protocols/directdefund"
	"github.com/statechannels/go-nitro/protocols/directfund"
	"github.com/statechannels/go-nitro/types"
	"github.com/stretchr/testify/assert"
)

func setupClientWithP2PMessageService(pk []byte, port int, chain *chainservice.MockChainService, logDestination io.Writer) (client.Client, *p2pms.P2PMessageService) {
	messageservice := p2pms.NewMessageService("127.0.0.1", port, pk)
	storeA := store.NewMemStore(pk)

	return client.New(messageservice, chain, storeA, logDestination, &engine.PermissivePolicy{}, nil), messageservice
}

// User real blockchain simulated_backend_service
func TestRunRpcNats(t *testing.T) {
	logger := zerolog.New(zerolog.ConsoleWriter{
		Out:           os.Stdout,
		TimeFormat:    time.RFC3339,
		PartsOrder:    []string{"time", "level", "caller", "client", "scope", "message"},
		FieldsExclude: []string{"time", "level", "caller", "message", "client", "scope"},
	}).
		Level(zerolog.InfoLevel).
		With().
		Timestamp().
		Str("client", "").
		Str("scope", "").
		Logger()

	opts := &server.Options{}
	ns, err := server.NewServer(opts)

	assert.NoError(t, err)
	ns.Start()

	nc, err := nats.Connect(ns.ClientURL())
	chain := chainservice.NewMockChain()

	alice := testactors.Alice
	bob := testactors.Bob

	chainServiceA := chainservice.NewMockChainService(chain, alice.Address())
	chainServiceB := chainservice.NewMockChainService(chain, bob.Address())

	trp := nats2.NewNatsTransport(nc, []string{fmt.Sprintf("nitro.%s", network.DirectFundRequestMethod), "nitro.test-topic"})

	clientA, msgA := setupClientWithP2PMessageService(alice.PrivateKey, 3005, chainServiceA, logger)
	clientB, msgB := setupClientWithP2PMessageService(bob.PrivateKey, 3006, chainServiceB, logger)

	con, err := trp.PollConnection()
	if err != nil {
		assert.NoError(t, err)
	}

	nts := network.NewNetworkService(con, &serde.JsonRpc{})
	nts.RegisterResponseHandler(network.DirectFundRequestMethod, func(m *netproto.Message) {
		logger.Info().Msgf("Objective updated: %v", *m)
	})

	nts.RegisterErrorHandler(network.DirectFundRequestMethod, func(m *netproto.Message) {
		logger.Error().Msgf("Objective failed: %v", *m)
	})

	peers := []p2pms.PeerInfo{
		{Id: msgA.Id(), IpAddress: "127.0.0.1", Port: 3005, Address: alice.Address()},
		{Id: msgB.Id(), IpAddress: "127.0.0.1", Port: 3006, Address: bob.Address()},
	}

	// Connect nitro P2P message services
	msgA.AddPeers(peers)
	msgB.AddPeers(peers)

	defer msgA.Close()
	defer msgB.Close()

	trpA := nats2.NewNatsTransport(nc, []string{
		fmt.Sprintf("nitro.%s", network.DirectFundRequestMethod),
		fmt.Sprintf("nitro.%s", network.DirectDefundRequestMethod),
	})
	conA, err := trpA.PollConnection()
	if err != nil {
		logger.Fatal().Msg(err.Error())
	}
	ntsA := network.NewNetworkService(conA, &serde.JsonRpc{})

	objReq := directfund.NewObjectiveRequest(
		alice.Address(),
		100, testdata.Outcomes.Create(alice.Address(), bob.Address(), 100, 100, types.Address{}), rand.Uint64(), bob.Address())

	//&directfund.ObjectiveRequest{
	//	CounterParty:      alice.Address(),
	//	ChallengeDuration: 100,
	//	Outcome:           testdata.Outcomes.Create(alice.Address(), bob.Address(), 100, 100, types.Address{}),
	//	AppDefinition:     chainServiceA.GetConsensusAppAddress(),
	//	// Appdata implicitly zero
	//	Nonce: rand.Uint64(),
	//}

	ntsA.RegisterRequestHandler(network.DirectFundRequestMethod, func(m *netproto.Message) {
		if len(m.Args) < 1 {
			logger.Fatal().Msg("unexpected empty args for direct funding method")
			return
		}

		for i := 0; i < len(m.Args); i++ {
			res := m.Args[i].(map[string]interface{})
			req := parser.ParseDirectFundRequest(res)

			assert.Equal(t, req.CounterParty, alice.Address())
			assert.Equal(t, req.AppDefinition, chainServiceA.GetConsensusAppAddress())

			objRes := req.Response(alice.Address(), nil)

			nts.SendMessage(netproto.NewMessage(netproto.TypeResponse, m.RequestId, network.DirectFundRequestMethod, []any{&objRes}))
			clientA.Engine.ObjectiveRequestsFromAPI <- req
		}
	})

	ntsA.RegisterRequestHandler(network.DirectDefundRequestMethod, func(m *netproto.Message) {
		if len(m.Args) < 1 {
			logger.Fatal().Msg("unexpected empty args for direct defunding method")
			return
		}

		for i := 0; i < len(m.Args); i++ {
			res := m.Args[i].(map[string]any)
			req := parser.ParseDirectDefundRequest(res)
			clientB.Engine.ObjectiveRequestsFromAPI <- req
		}
	})

	nts.SendMessage(netproto.NewMessage(netproto.TypeRequest, rand.Uint64(), network.DirectFundRequestMethod, []any{&objReq}))
	go func() {
		time.Sleep(100)
		fixedPart := state.FixedPart{
			Participants:      []types.Address{alice.Address(), bob.Address()},
			ChannelNonce:      objReq.Nonce,
			AppDefinition:     objReq.AppDefinition,
			ChallengeDuration: objReq.ChallengeDuration,
		}
		req := directdefund.ObjectiveRequest{ChannelId: fixedPart.ChannelId()}
		nts.SendMessage(netproto.NewMessage(netproto.TypeRequest, rand.Uint64(), network.DirectDefundRequestMethod, []any{&req}))
	}()
	<-time.After(1 * time.Second)
}

//func initNats() *nats.Conn {
//	opts := &server.Options{}
//	ns, _ := server.NewServer(opts)
//
//	nc, _ := nats.Connect(ns.ClientURL())
//
//	return nc
//}
//
//func initLogger() zerolog.Logger {
//	logger := zerolog.New(zerolog.ConsoleWriter{
//		Out:           os.Stdout,
//		TimeFormat:    time.RFC3339,
//		PartsOrder:    []string{"time", "level", "caller", "client", "scope", "message"},
//		FieldsExclude: []string{"time", "level", "caller", "message", "client", "scope"},
//	}).
//		Level(zerolog.InfoLevel).
//		With().
//		Timestamp().
//		Str("client", "").
//		Str("scope", "").
//		Logger()
//
//	return logger
//}
//
//func TestFundDefundFlow(t *testing.T) {
//	nc := initNats()
//	logger := initLogger()
//	//broker := messageservice.NewBroker()
//
//	bob := testactors.Bob
//	alice := testactors.Alice
//
//	// Tried yo use simulated blockchain but no luck
//	//sim, bindings, _, err := chainservice.SetupSimulatedBackend(0)
//	chain := chainservice.NewMockChain()
//	chainServiceA := chainservice.NewMockChainService(chain, bob.Address())
//	chainServiceB := chainservice.NewMockChainService(chain, alice.Address())
//
//	clientA, _ := setupClientWithP2PMessageService(bob.PrivateKey, 3005, chainServiceA, logger)
//	clientB, _ := setupClientWithP2PMessageService(alice.PrivateKey, 3006, chainServiceB, logger)
//
//	trp := nats2.NewNatsTransport(nc, []string{
//		fmt.Sprintf("nitro.%s", directfund.DirectFundRequestMethod),
//		fmt.Sprintf("nitro.%s", directdefund.DirectDefundRequestMethod),
//	})
//	natsConn, err := trp.PollConnection()
//	assert.NoError(t, err, "we should be able to poll connection")
//
//	// initialize our network service with nats connection
//	networkService := network.NewNetworkService(natsConn, &serde.MsgPack{})
//	defer networkService.Close()
//
//	// define messages
//	objReq := &directfund.ObjectiveRequest{
//		CounterParty:      bob.Address(),
//		ChallengeDuration: 100,
//		Outcome:           testdata.Outcomes.Create(bob.Address(), alice.Address(), 100, 200, types.Address{}),
//		// Not too sure if this is right
//		AppDefinition: chainServiceA.GetConsensusAppAddress(),
//		// Appdata implicitly zero
//		Nonce: rand.Uint64(),
//	}
//	channelId := getChannelIdFromFundObjectiveRequest(objReq, []types.Address{bob.Address(), alice.Address()})
//
//	networkService.RegisterRequestHandler(directdefund.DirectDefundRequestMethod, func(m *netproto.Message) {
//		if len(m.Args) < 1 {
//			return
//		}
//
//		for i := 0; i < len(m.Args); i++ {
//			res := m.Args[i].(map[string]interface{})
//			req := directfund.NewDirectFundObjectiveRequest(res)
//
//			assert.Equal(t, req.CounterParty, bob.Address())
//			assert.Equal(t, req.AppDefinition, chainServiceA.GetConsensusAppAddress())
//
//			clientA.Engine.ObjectiveRequestsFromAPI <- req
//			clientB.Engine.ObjectiveRequestsFromAPI <- req
//
//			engineEvent := clientA.Engine.ToApi()
//			println(engineEvent)
//			response := directfund.ObjectiveResponse{
//				Id:        protocols.ObjectiveId(req.CounterParty.String() + fmt.Sprint(req.Nonce)),
//				ChannelId: channelId,
//			}
//			networkService.SendMessage(response.ToResponseMessage(m.RequestId))
//		}
//	})
//
//	networkService.RegisterResponseHandler(directfund.DirectFundRequestMethod, func(m *netproto.Message) {
//
//	})
//
//	networkService.RegisterRequestHandler(directdefund.DirectDefundRequestMethod, func(m *netproto.Message) {
//		if len(m.Args) < 1 {
//			return
//		}
//	})
//}
//
//func getChannelIdFromFundObjectiveRequest(req *directfund.ObjectiveRequest, participants []types.Address) types.Destination {
//	fixedPart := state.FixedPart{
//		ChainId:           state.TestState.ChainId,
//		Participants:      participants,
//		ChannelNonce:      req.Nonce,
//		AppDefinition:     req.AppDefinition,
//		ChallengeDuration: req.ChallengeDuration,
//	}
//
//	return fixedPart.ChannelId()
//}
