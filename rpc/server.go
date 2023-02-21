package rpc

import (
	"encoding/json"
	"fmt"
	"math/big"
	"math/rand"

	"github.com/rs/zerolog"
	nitro "github.com/statechannels/go-nitro/client"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/protocols/directdefund"
	"github.com/statechannels/go-nitro/protocols/directfund"
	"github.com/statechannels/go-nitro/protocols/virtualdefund"
	"github.com/statechannels/go-nitro/protocols/virtualfund"
	"github.com/statechannels/go-nitro/rpc/serde"
	"github.com/statechannels/go-nitro/rpc/transport"
)

// RpcServer handles nitro rpc requests and executes them on the nitro client
type RpcServer struct {
	connection transport.Responder
	client     *nitro.Client
	logger     zerolog.Logger
}

func (rs *RpcServer) Url() string {
	return rs.connection.Url()
}

func (rs *RpcServer) Close() {
	rs.connection.Close()
}

func NewRpcServer(nitroClient *nitro.Client, logger zerolog.Logger, connection transport.Responder) (*RpcServer, error) {
	rs := &RpcServer{connection, nitroClient, logger}
	rs.sendNotifications()
	err := rs.registerHandlers()
	if err != nil {
		return nil, err
	}

	return rs, nil
}

// registerHandlers registers the handlers for the rpc server
func (rs *RpcServer) registerHandlers() error {
	subscriber := func(requestData []byte) []byte {
		rs.logger.Trace().Msgf("Rpc server received request: %+v", string(requestData))
		requestJson := serde.AnyJsonRpcRequest{}
		err := json.Unmarshal(requestData, &requestJson)
		if err != nil {
			panic(err)
		}

		switch serde.RequestMethod(requestJson.Method) {
		case serde.DirectFundRequestMethod:
			return processRequest(rs, requestData, func(obj directfund.ObjectiveRequest) directfund.ObjectiveResponse {
				return rs.client.CreateLedgerChannel(obj.CounterParty, obj.ChallengeDuration, obj.Outcome)
			})
		case serde.DirectDefundRequestMethod:
			return processRequest(rs, requestData, func(obj directdefund.ObjectiveRequest) protocols.ObjectiveId {
				return rs.client.CloseLedgerChannel(obj.ChannelId)
			})
		case serde.VirtualFundRequestMethod:
			return processRequest(rs, requestData, func(obj virtualfund.ObjectiveRequest) virtualfund.ObjectiveResponse {
				return rs.client.CreateVirtualPaymentChannel(obj.Intermediaries, obj.CounterParty, obj.ChallengeDuration, obj.Outcome)
			})
		case serde.VirtualDefundRequestMethod:
			return processRequest(rs, requestData, func(obj virtualdefund.ObjectiveRequest) protocols.ObjectiveId {
				return rs.client.CloseVirtualChannel(obj.ChannelId)
			})
		case serde.PayRequestMethod:
			return processRequest(rs, requestData, func(obj serde.PaymentRequest) serde.PaymentRequest {
				rs.client.Pay(obj.Channel, big.NewInt(int64(obj.Amount)))
				return obj
			})
		default:
			panic(fmt.Errorf("unknown method: %s", requestJson.Method))
		}
	}
	return rs.connection.Respond(subscriber)
}

func processRequest[T serde.RequestPayload, U serde.ResponsePayload](rs *RpcServer, requestData []byte, processPayload func(T) U) []byte {
	rs.logger.Trace().Msgf("Rpc server received request: %+v", requestData)
	rpcRequest := serde.JsonRpcRequest[T]{}
	err := json.Unmarshal(requestData, &rpcRequest)
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
			err = rs.connection.Notify(data)
			if err != nil {
				panic(err)
			}
		}
	}()
}
