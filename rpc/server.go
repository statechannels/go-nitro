package rpc

import (
	"encoding/json"
	"math/big"
	"math/rand"

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
	"github.com/statechannels/go-nitro/rpc/transport/wss"
)

// RpcServer handles nitro rpc requests and executes them on the nitro client
type RpcServer struct {
	connection transport.Subscriber
	ns         *server.Server // TODO we don't want this nats-specific thing in here
	client     *nitro.Client
	chainId    *big.Int
	logger     zerolog.Logger
	port       string
}

func (rs *RpcServer) Url() string {
	return "ws://127.0.0.1:" + rs.port
}

func (rs *RpcServer) Close() {
	rs.connection.Close()
	rs.ns.Shutdown()
}

func NewRpcServer(port string, nitroClient *nitro.Client, chainId *big.Int, logger zerolog.Logger) *RpcServer {

	ws := wss.NewWebSocketConnectionAsServer(port)

	rs := &RpcServer{ws, nil, nitroClient, chainId, logger, port}
	rs.sendNotifications()
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

func (rs *RpcServer) sendNotifications() {
	go func() {
		for completedObjective := range rs.client.CompletedObjectives() {
			request := serde.NewJsonRpcRequest(rand.Uint64(), serde.ObjectiveCompleted, completedObjective)
			data, err := json.Marshal(request)
			if err != nil {
				panic(err)
			}
			err = rs.connection.Notify(serde.ObjectiveCompleted, data)
			if err != nil {
				panic(err)
			}
		}
	}()
}

func subcribeToRequest[T serde.RequestPayload, U serde.ResponsePayload](rs *RpcServer, method serde.RequestMethod, processPayload func(T) U) error {
	return rs.connection.Respond(method, func(data []byte) []byte {
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
