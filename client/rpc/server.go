package rpc

import (
	"fmt"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	nitro "github.com/statechannels/go-nitro/client"
	"github.com/statechannels/go-nitro/network"
	netproto "github.com/statechannels/go-nitro/network/protocol"
	"github.com/statechannels/go-nitro/network/protocol/parser"
	"github.com/statechannels/go-nitro/network/serde"
	natstrans "github.com/statechannels/go-nitro/network/transport/nats"
)

// RpcServer handles nitro rpc requests and executes them on the nitro client
type RpcServer struct {
	nts    *network.NetworkService
	ns     *server.Server
	client *nitro.Client
}

func (rs *RpcServer) Url() string {
	return rs.ns.ClientURL()
}

func (rs *RpcServer) Close() {
	rs.ns.Shutdown()
	rs.nts.Close()
}

func NewRpcServer(nitroClient *nitro.Client) *RpcServer {

	opts := &server.Options{}
	ns, err := server.NewServer(opts)
	handleError(err)
	ns.Start()

	nc, err := nats.Connect(ns.ClientURL())
	handleError(err)

	trp := natstrans.NewNatsTransport(nc, []string{
		fmt.Sprintf("nitro.%s", network.DirectFundRequestMethod),
		fmt.Sprintf("nitro.%s", network.DirectDefundRequestMethod),
	})

	con, err := trp.PollConnection()
	handleError(err)

	nts := network.NewNetworkService(con, &serde.JsonRpc{})

	rs := &RpcServer{nts, ns, nitroClient}
	rs.registerHandlers()
	return rs
}

// registerHandlers registers the handlers for the rpc server
func (rs *RpcServer) registerHandlers() {
	rs.nts.RegisterRequestHandler(network.DirectFundRequestMethod, func(m *netproto.Message) {
		if len(m.Args) < 1 {
			panic("unexpected empty args for direct funding method")

		}

		for i := 0; i < len(m.Args); i++ {

			raw := m.Args[i].(map[string]interface{})
			req := parser.ParseDirectFundRequest(raw)
			objRes := req.Response(*rs.client.Address, nil)

			rs.client.CreateLedgerChannel(req.CounterParty, req.ChallengeDuration, req.Outcome)
			msg := netproto.NewMessage(netproto.TypeResponse, m.RequestId, network.DirectFundRequestMethod, []any{&objRes})

			rs.nts.SendMessage(msg)

		}
	})
}

// handleError "handles" an error by panicking
func handleError(err error) {

	if err != nil {
		panic(err)
	}
}
