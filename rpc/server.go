package rpc

import (
	"encoding/json"
	"math/big"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/rs/zerolog"
	nitro "github.com/statechannels/go-nitro/client"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/protocols/directdefund"
	"github.com/statechannels/go-nitro/protocols/directfund"
	"github.com/statechannels/go-nitro/protocols/virtualdefund"
	"github.com/statechannels/go-nitro/protocols/virtualfund"
	"github.com/statechannels/go-nitro/rpc/serde"
	"github.com/statechannels/go-nitro/rpc/transport"
	"github.com/statechannels/go-nitro/rpc/transport/nats/wss"
)

// RpcServer handles nitro rpc requests and executes them on the nitro client
type RpcServer struct {
	connection transport.Subscriber
	ns         *server.Server // TODO we don't want this nats-specific thing in here
	client     *nitro.Client
	chainId    *big.Int
	logger     zerolog.Logger
}

func (rs *RpcServer) Url() string {
	return "127.0.0.1:1234"
}

func (rs *RpcServer) Close() {
	rs.connection.Close()
	rs.ns.Shutdown()
}

func NewRpcServer(nitroClient *nitro.Client, chainId *big.Int, logger zerolog.Logger) *RpcServer {

	ws := wss.NewWebSocketConnectionAsServer("1234")

	rs := &RpcServer{ws, nil, nitroClient, chainId, logger}
	err := rs.registerHandlers()
	if err != nil {
		panic(err)
	}

	return rs
}

// registerHandlers registers the handlers for the rpc server
func (rs *RpcServer) registerHandlers() error {
	err := subcribeToRequest(rs, serde.DirectFundRequestMethod, func(obj directfund.ObjectiveRequest) directfund.ObjectiveResponse {
		return rs.client.CreateLedgerChannel(obj.CounterParty, obj.ChallengeDuration, obj.Outcome)
	})

	if err != nil {
		return err
	}

	err = subcribeToRequest(rs, serde.DirectDefundRequestMethod, func(obj directdefund.ObjectiveRequest) protocols.ObjectiveId {
		return rs.client.CloseLedgerChannel(obj.ChannelId)
	})

	if err != nil {
		return err
	}

	err = subcribeToRequest(rs, serde.VirtualFundRequestMethod, func(obj virtualfund.ObjectiveRequest) virtualfund.ObjectiveResponse {
		return rs.client.CreateVirtualPaymentChannel(obj.Intermediaries, obj.CounterParty, obj.ChallengeDuration, obj.Outcome)
	})

	if err != nil {
		return err
	}

	err = subcribeToRequest(rs, serde.VirtualDefundRequestMethod, func(obj virtualdefund.ObjectiveRequest) protocols.ObjectiveId {
		return rs.client.CloseVirtualChannel(obj.ChannelId)
	})

	if err != nil {
		return err
	}

	err = subcribeToRequest(rs, serde.PayRequestMethod, func(payReq serde.PaymentRequest) serde.PaymentRequest {
		rs.client.Pay(payReq.Channel, big.NewInt(int64(payReq.Amount)))
		return payReq
	})

	return err
}

func subcribeToRequest[T serde.RequestPayload, U serde.ResponsePayload](rs *RpcServer, method serde.RequestMethod, processPayload func(T) U) error {
	return rs.connection.Subscribe(method, func(data []byte) []byte {
		rs.logger.Trace().Msgf("Rpc server received request: %+v", data)
		rpcRequest := serde.JsonRpcRequest[T]{}
		err := json.Unmarshal(data, &rpcRequest)
		if err != nil {
			panic("could not unmarshal objective request")
		}
		obj := rpcRequest.Params
		objResponse := processPayload(obj)

		msg := serde.NewJsonRpcResponse(rpcRequest.Id, objResponse)
		messageData, err := json.Marshal(msg)
		if err != nil {
			panic("Could not marshal response message")
		}

		return messageData
	})
}
