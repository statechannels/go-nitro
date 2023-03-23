package rpc

import (
	"encoding/json"
	"math/big"

	"github.com/rs/zerolog"
	nitro "github.com/statechannels/go-nitro/client"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/protocols/directdefund"
	"github.com/statechannels/go-nitro/protocols/directfund"
	"github.com/statechannels/go-nitro/protocols/virtualdefund"
	"github.com/statechannels/go-nitro/protocols/virtualfund"
	"github.com/statechannels/go-nitro/rand"
	"github.com/statechannels/go-nitro/rpc/serde"
	"github.com/statechannels/go-nitro/rpc/transport"
)

// RpcServer handles nitro rpc requests and executes them on the nitro client
type RpcServer struct {
	transport transport.Responder
	client    *nitro.Client
	logger    *zerolog.Logger
}

func (rs *RpcServer) Url() string {
	return rs.transport.Url()
}

func (rs *RpcServer) Close() {
	rs.client.Close()
	rs.transport.Close()
}

func NewRpcServer(nitroClient *nitro.Client, logger *zerolog.Logger, trans transport.Responder) (*RpcServer, error) {
	rs := &RpcServer{trans, nitroClient, logger}
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

		if !json.Valid(requestData) {
			return marshalResponse(parseError, rs.logger)
		}

		validationResult := validateRequest(requestData, rs.logger)
		if validationResult.Error != nil {
			return validationResult.Error
		}

		switch serde.RequestMethod(validationResult.Method) {
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
			responseErr := methodNotFoundError
			responseErr.Id = validationResult.Id
			return marshalResponse(responseErr, rs.logger)
		}
	}
	return rs.transport.RegisterRequestHandler(subscriber)
}

func processRequest[T serde.RequestPayload, U serde.ResponsePayload](rs *RpcServer, requestData []byte, processPayload func(T) U) []byte {
	rpcRequest := serde.JsonRpcRequest[T]{}
	// todo: unmarshal will fail only when the requestData is not valid json.
	// At the moment, there is no validation that the required fields are populated in the request.
	err := json.Unmarshal(requestData, &rpcRequest)
	if err != nil {
		return marshalResponse(unexpectedRequestUnmarshalError2, rs.logger)
	}
	obj := rpcRequest.Params
	objResponse := processPayload(obj)
	response := serde.NewJsonRpcResponse(rpcRequest.Id, objResponse)
	return marshalResponse(response, rs.logger)
}

// Marshal and return response data
func marshalResponse(response any, log *zerolog.Logger) []byte {
	responseData, err := json.Marshal(response)
	if err != nil {
		log.Panic().Err(err).Msg("Could not marshal response")
	}
	return responseData
}

type validationResult struct {
	Error  []byte
	Method string
	Id     uint64
}

func validateRequest(requestData []byte, logger *zerolog.Logger) validationResult {
	var request map[string]interface{}
	vr := validationResult{}
	err := json.Unmarshal(requestData, &request)
	if err != nil {
		vr.Error = marshalResponse(unexpectedRequestUnmarshalError, logger)
		return vr
	}

	// jsonrpc spec says id can be a string, number.
	// We only support numbers: https://github.com/statechannels/go-nitro/issues/1160
	// When golang unmarshals JSON into an interface value, float64 is used for numbers.
	requestId := request["id"]
	fRequestId, ok := requestId.(float64)
	if !ok {
		vr.Error = marshalResponse(invalidRequestError, logger)
		return vr
	}

	if fRequestId != float64(uint64(fRequestId)) {
		vr.Error = marshalResponse(invalidRequestError, logger)
		return vr
	}
	vr.Id = uint64(fRequestId)

	sJsonrpc, ok := request["jsonrpc"].(string)
	if !ok || sJsonrpc != "2.0" {
		requestError := invalidRequestError
		requestError.Id = vr.Id
		vr.Error = marshalResponse(requestError, logger)
		return vr
	}

	_, ok = request["method"].(string)
	if !ok {
		requestError := invalidRequestError
		requestError.Id = vr.Id
		vr.Error = marshalResponse(requestError, logger)
		return vr
	}
	vr.Method = request["method"].(string)

	return vr
}

func (rs *RpcServer) sendNotifications() {
	go func() {
		for completedObjective := range rs.client.CompletedObjectives() {
			rs.logger.Trace().Msgf("Sending notification: %+v", completedObjective)
			request := serde.NewJsonRpcRequest(rand.GetRandGenerator().Uint64(), serde.ObjectiveCompleted, completedObjective)
			data, err := json.Marshal(request)
			if err != nil {
				panic(err)
			}
			err = rs.transport.Notify(data)
			if err != nil {
				panic(err)
			}
		}
	}()
}
