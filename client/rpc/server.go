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
	"github.com/statechannels/go-nitro/protocols/directdefund"
	"github.com/statechannels/go-nitro/protocols/directfund"
	"github.com/statechannels/go-nitro/protocols/virtualdefund"
	"github.com/statechannels/go-nitro/protocols/virtualfund"
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

	trp := natstrans.NewNatsTransport(nc, getTopics())

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
	rs.nts.Subscribe(fmt.Sprintf("nitro.%s", serde.DirectFundRequestMethod), func(data []byte) []byte {

		rs.nts.Logger.Trace().Msgf("Rpc server received request: %+v", data)

		rpcRequest := serde.JsonRpcRequest[directfund.ObjectiveRequest]{}
		err := json.Unmarshal(data, &rpcRequest)
		if err != nil {
			panic("could not unmarshal direct fund objective request")
		}

		// todo: objective request is redefined so that it has a valid objectiveStarted channel.
		// 	Should find a better way to accomplish this.
		objectiveRequestWithChan := directfund.NewObjectiveRequest(
			rpcRequest.Params.CounterParty,
			rpcRequest.Params.ChallengeDuration,
			rpcRequest.Params.Outcome,
			rpcRequest.Params.Nonce,
			rpcRequest.Params.AppDefinition,
		)

		rs.client.IncomingObjectiveRequests() <- objectiveRequestWithChan

		objRes := rpcRequest.Params.Response(*rs.client.Address, rs.chainId)
		msg := serde.NewJsonRpcResponse(rpcRequest.Id, objRes)
		messageData, err := json.Marshal(msg)
		if err != nil {
			panic("Could not marshal direct fund response message")
		}

		return messageData
	})

	rs.nts.Subscribe(fmt.Sprintf("nitro.%s", serde.DirectDefundRequestMethod), func(data []byte) []byte {
		rs.nts.Logger.Trace().Msgf("Rpc server received request: %+v", data)

		rpcRequest := serde.JsonRpcRequest[directdefund.ObjectiveRequest]{}
		err := json.Unmarshal(data, &rpcRequest)
		if err != nil {
			panic("could not unmarshal direct fund objective request")
		}

		// todo: objective request is redefined so that it has a valid objectiveStarted channel.
		// 	Should find a better way to accomplish this.
		objectiveRequestWithChan := directdefund.NewObjectiveRequest(
			rpcRequest.Params.ChannelId,
		)

		rs.client.IncomingObjectiveRequests() <- objectiveRequestWithChan

		objRes := rpcRequest.Params.Id(*rs.client.Address, rs.chainId)
		msg := serde.NewJsonRpcResponse(rpcRequest.Id, objRes)
		messageData, err := json.Marshal(msg)
		if err != nil {
			panic("Could not marshal direct fund response message")
		}

		return messageData
	})

	rs.nts.Subscribe(fmt.Sprintf("nitro.%s", serde.VirtualFundRequestMethod), func(data []byte) []byte {

		rs.nts.Logger.Trace().Msgf("Rpc server received request: %+v", data)

		rpcRequest := serde.JsonRpcRequest[virtualfund.ObjectiveRequest]{}
		err := json.Unmarshal(data, &rpcRequest)
		if err != nil {
			panic("could not unmarshal virtual fund objective request")
		}

		// todo: objective request is redefined so that it has a valid objectiveStarted channel.
		// 	Should find a better way to accomplish this.
		objectiveRequestWithChan := virtualfund.NewObjectiveRequest(
			rpcRequest.Params.Intermediaries,
			rpcRequest.Params.CounterParty,
			rpcRequest.Params.ChallengeDuration,
			rpcRequest.Params.Outcome,
			rpcRequest.Params.Nonce,
			rpcRequest.Params.AppDefinition,
		)

		rs.client.IncomingObjectiveRequests() <- objectiveRequestWithChan

		objRes := rpcRequest.Params.Response(*rs.client.Address)
		msg := serde.NewJsonRpcResponse(rpcRequest.Id, objRes)
		messageData, err := json.Marshal(msg)
		if err != nil {
			panic("Could not marshal direct fund response message")
		}

		return messageData
	})

	rs.nts.Subscribe(fmt.Sprintf("nitro.%s", serde.VirtualDefundRequestMethod), func(data []byte) []byte {
		rs.nts.Logger.Trace().Msgf("Rpc server received request: %+v", data)

		rpcRequest := serde.JsonRpcRequest[virtualdefund.ObjectiveRequest]{}
		err := json.Unmarshal(data, &rpcRequest)
		if err != nil {
			panic("could not unmarshal virtual defund objective request")
		}

		// todo: objective request is redefined so that it has a valid objectiveStarted channel.
		// 	Should find a better way to accomplish this.
		objectiveRequestWithChan := virtualdefund.NewObjectiveRequest(
			rpcRequest.Params.ChannelId,
		)

		rs.client.IncomingObjectiveRequests() <- objectiveRequestWithChan

		objRes := rpcRequest.Params.Id(*rs.client.Address, rs.chainId)
		msg := serde.NewJsonRpcResponse(rpcRequest.Id, objRes)
		messageData, err := json.Marshal(msg)
		if err != nil {
			panic("Could not marshal direct fund response message")
		}

		return messageData
	})

	rs.nts.Subscribe(fmt.Sprintf("nitro.%s", serde.PayRequestMethod), func(data []byte) []byte {
		rs.nts.Logger.Trace().Msgf("Rpc server received request: %+v", data)

		rpcRequest := serde.JsonRpcRequest[serde.PaymentRequest]{}
		err := json.Unmarshal(data, &rpcRequest)
		if err != nil {
			panic("could not unmarshal pay objective request")
		}

		rs.client.Pay(rpcRequest.Params.Channel, big.NewInt(int64(rpcRequest.Params.Amount)))

		// TODO: What should we return here? A voucher?
		msg := serde.NewJsonRpcResponse(rpcRequest.Id, rpcRequest.Params)
		messageData, err := json.Marshal(msg)
		if err != nil {
			panic("Could not marshal pay response message")
		}

		return messageData
	})
}

// handleError "handles" an error by panicking
func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
