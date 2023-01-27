package client_test

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	p2pms "github.com/statechannels/go-nitro/client/engine/messageservice/p2p-message-service"
	"github.com/statechannels/go-nitro/internal/testactors"
	"github.com/statechannels/go-nitro/internal/testdata"
	"github.com/statechannels/go-nitro/network"
	netproto "github.com/statechannels/go-nitro/network/protocol"
	"github.com/statechannels/go-nitro/network/protocol/parser"
	"github.com/statechannels/go-nitro/network/serde"
	nats2 "github.com/statechannels/go-nitro/network/transport/nats"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/protocols/directfund"
	"github.com/statechannels/go-nitro/types"
	"github.com/stretchr/testify/assert"
)

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
	if err != nil {
		t.Error(err)
	}
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

	dfObjResponseChan := make(chan directfund.ObjectiveResponse, 1)
	ntsA.RegisterRequestHandler(network.DirectFundRequestMethod, func(m *netproto.Message) {
		if len(m.Args) < 1 {
			logger.Fatal().Msg("unexpected empty args for direct funding method")
			return
		}

		for i := 0; i < len(m.Args); i++ {
			res := m.Args[i].(map[string]interface{})
			req := parser.ParseDirectFundRequest(res)

			assert.Equal(t, req.CounterParty, bob.Address())
			assert.Equal(t, req.AppDefinition, chainServiceA.GetConsensusAppAddress())

			objRes := req.Response(alice.Address(), nil)

			nts.SendMessage(netproto.NewMessage(netproto.TypeResponse, m.RequestId, network.DirectFundRequestMethod, []any{&objRes}))
			dfObjResponseChan <- clientA.CreateLedgerChannel(req.CounterParty, req.ChallengeDuration, req.Outcome)
		}
	})

	// todo: add this test back
	// ntsA.RegisterRequestHandler(network.DirectDefundRequestMethod, func(m *netproto.Message) {
	// 	if len(m.Args) < 1 {
	// 		logger.Fatal().Msg("unexpected empty args for direct defunding method")
	// 		return
	// 	}

	// 	for i := 0; i < len(m.Args); i++ {
	// 		res := m.Args[i].(map[string]any)
	// 		req := parser.ParseDirectDefundRequest(res)
	// 		clientB.CloseLedgerChannel(req.ChannelId)
	// 	}
	// })

	objReq := directfund.NewObjectiveRequest(
		bob.Address(),
		100, testdata.Outcomes.Create(alice.Address(), bob.Address(), 100, 100, types.Address{}), rand.Uint64(), common.Address{})
	nts.SendMessage(netproto.NewMessage(netproto.TypeRequest, rand.Uint64(), network.DirectFundRequestMethod, []any{&objReq}))

	dfObjResponse := <-dfObjResponseChan
	ids := []protocols.ObjectiveId{dfObjResponse.Id}
	waitTimeForCompletedObjectiveIds(t, &clientA, defaultTimeout, ids...)
	waitTimeForCompletedObjectiveIds(t, &clientB, defaultTimeout, ids...)

	// todo: add me back
	// req := directdefund.ObjectiveRequest{ChannelId: dfObjResponse.ChannelId}
	// nts.SendMessage(netproto.NewMessage(netproto.TypeRequest, rand.Uint64(), network.DirectDefundRequestMethod, []any{&req}))
}
