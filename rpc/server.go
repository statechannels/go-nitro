package rpc

import (
	"encoding/json"
	"fmt"
	"math/big"
	"math/rand"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	nitro "github.com/statechannels/go-nitro/client"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/protocols/directdefund"
	"github.com/statechannels/go-nitro/protocols/directfund"
	"github.com/statechannels/go-nitro/protocols/virtualdefund"
	"github.com/statechannels/go-nitro/protocols/virtualfund"
	"github.com/statechannels/go-nitro/rpc/serde"
	"github.com/statechannels/go-nitro/rpc/transport"
	natstrans "github.com/statechannels/go-nitro/rpc/transport/nats"
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

func NewRpcServer(nitroClient *nitro.Client, chainId *big.Int, logger zerolog.Logger, rpcPort int, connectionType transport.ConnectionType) (*RpcServer, error) {
	var con transport.Subscriber
	var err error
	switch connectionType {
	case transport.Nats:
		opts := &server.Options{Port: rpcPort}
		ns, err := server.NewServer(opts)
		if err != nil {
			return nil, err
		}
		ns.Start()

		nc, err := nats.Connect(ns.ClientURL())
		if err != nil {
			return nil, err
		}
		con = natstrans.NewNatsConnection(nc)
	case transport.Ws:
		con = wss.NewWebSocketConnectionAsServer(fmt.Sprint(rpcPort))
	default:
		return nil, fmt.Errorf("unknown connection type %v", connectionType)
	}

	rs := &RpcServer{con, nil, nitroClient, chainId, logger, fmt.Sprint(rpcPort)}
	rs.sendNotifications()
	err = rs.registerHandlers()
	if err != nil {
		return nil, err
	}

	return rs, nil
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
			rs.logger.Trace().Msgf("Sending notification: %+v", completedObjective)
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
		rs.logger.Trace().Msgf("Rpc server received request: %+v", string(data))
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
