package rpc

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	nitro "github.com/statechannels/go-nitro/client"
	"github.com/statechannels/go-nitro/network"
	"github.com/statechannels/go-nitro/network/serde"
	natstrans "github.com/statechannels/go-nitro/network/transport/nats"
	"github.com/statechannels/go-nitro/protocols/directfund"
)

// RpcServer handles nitro rpc requests and executes them on the nitro client
type RpcServer struct {
	nts     *network.NetworkService
	ns      *server.Server
	client  *nitro.Client
	chainId *big.Int
}

func (rs *RpcServer) Url() string {
	return rs.ns.ClientURL()
}

func (rs *RpcServer) Close() {
	rs.ns.Shutdown()
	rs.nts.Close()
}

func NewRpcServer(nitroClient *nitro.Client, chainId *big.Int, logger zerolog.Logger) *RpcServer {

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

	nts := network.NewNetworkService(con)
	nts.Logger = logger
	rs := &RpcServer{nts, ns, nitroClient, chainId}
	rs.registerHandlers()
	return rs
}

// registerHandlers registers the handlers for the rpc server
func (rs *RpcServer) registerHandlers() {
	rs.nts.RegisterRequestHandler(network.DirectFundRequestMethod, func(data []byte) {
		rs.nts.Logger.Trace().Msgf("Rpc server received request: %+v", data)

		rpcRequest := serde.JsonRpcDirectFundRequest{}
		err := json.Unmarshal(data, &rpcRequest)
		if err != nil {
			panic("could not unmarshal direct fund objective request")
		}

		// todo: objective request is redefined so that it has a valid objectiveStarted channel.
		// 	Should find a better way to accomplish this.
		objectiveRequestWithChan := directfund.NewObjectiveRequest(
			rpcRequest.ObjectiveRequest.CounterParty,
			rpcRequest.ObjectiveRequest.ChallengeDuration,
			rpcRequest.ObjectiveRequest.Outcome,
			rpcRequest.ObjectiveRequest.Nonce,
			rpcRequest.ObjectiveRequest.AppDefinition,
		)

		rs.client.IncomingObjectiveRequests() <- objectiveRequestWithChan

		objRes := rpcRequest.ObjectiveRequest.Response(*rs.client.Address, rs.chainId)
		msg := serde.NewDirectFundResponseMessage(rpcRequest.Id, objRes)
		messageData, err := json.Marshal(msg)
		if err != nil {
			panic("Could not marshal direct fund response message")
		}

		rs.nts.SendMessage(network.DirectFundRequestMethod, messageData)
	})
}

// handleError "handles" an error by panicking
func handleError(err error) {

	if err != nil {
		panic(err)
	}
}
