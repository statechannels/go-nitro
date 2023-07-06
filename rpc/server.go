package rpc

import (
	"context"
	"encoding/json"
	"math/big"
	"sync"

	"github.com/rs/zerolog"
	nitro "github.com/statechannels/go-nitro/node"
	"github.com/statechannels/go-nitro/node/query"
	"github.com/statechannels/go-nitro/payments"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/protocols/directdefund"
	"github.com/statechannels/go-nitro/protocols/directfund"
	"github.com/statechannels/go-nitro/protocols/virtualdefund"
	"github.com/statechannels/go-nitro/protocols/virtualfund"
	"github.com/statechannels/go-nitro/rand"
	"github.com/statechannels/go-nitro/rpc/serde"
	"github.com/statechannels/go-nitro/rpc/transport"
	"github.com/statechannels/go-nitro/types"
)

// RpcServer handles nitro rpc requests and executes them on the nitro node
type RpcServer struct {
	transport transport.Responder
	node      *nitro.Node
	logger    *zerolog.Logger
	cancel    context.CancelFunc
	wg        *sync.WaitGroup
}

func (rs *RpcServer) Url() string {
	return rs.transport.Url()
}

func (rs *RpcServer) Address() *types.Address {
	return rs.node.Address
}

func (rs *RpcServer) Close() error {
	rs.cancel()
	rs.wg.Wait()

	rs.transport.Close()
	return rs.node.Close()
}

// newRpcServerWithoutNotifications creates a new rpc server without notifications enabled
func newRpcServerWithoutNotifications(nitroNode *nitro.Node, logger *zerolog.Logger, trans transport.Responder) (*RpcServer, error) {
	rs := &RpcServer{trans, nitroNode, logger, func() {}, &sync.WaitGroup{}}

	err := rs.registerHandlers()
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func NewRpcServer(nitroNode *nitro.Node, logger *zerolog.Logger, trans transport.Responder) (*RpcServer, error) {
	ctx, cancel := context.WithCancel(context.Background())
	rs := &RpcServer{trans, nitroNode, logger, cancel, &sync.WaitGroup{}}

	rs.wg.Add(1)
	go rs.sendNotifications(ctx)
	err := rs.registerHandlers()
	if err != nil {
		return nil, err
	}

	return rs, nil
}

// registerHandlers registers the handlers for the rpc server
func (rs *RpcServer) registerHandlers() (err error) {
	handlerV1 := func(requestData []byte) []byte {
		rs.logger.Trace().Msgf("Rpc server received request: %+v", string(requestData))

		if !json.Valid(requestData) {
			rs.logger.Error().Msg("request is not valid json")
			errRes := serde.NewJsonRpcErrorResponse(0, serde.ParseError)
			return marshalResponse(errRes, rs.logger)
		}

		jsonrpcReq, errRes := validateJsonrpcRequest(requestData, rs.logger)
		if errRes != nil {
			rs.logger.Error().Msg("could not validate jsonrpc request")
			return errRes
		}

		switch serde.RequestMethod(jsonrpcReq.Method) {
		case serde.CreateVoucherRequestMethod:
			return processRequest(rs, requestData, func(req serde.PaymentRequest) (payments.Voucher, error) {
				v, err := rs.node.CreateVoucher(req.Channel, big.NewInt(int64(req.Amount)))
				if err != nil {
					return payments.Voucher{}, err
				}
				return v, nil
			})
		case serde.ReceiveVoucherRequestMethod:
			return processRequest(rs, requestData, func(req payments.Voucher) (serde.ReceiveVoucherResponse, error) {
				total, delta, err := rs.node.ReceiveVoucher(req)
				if err != nil {
					return serde.ReceiveVoucherResponse{}, err
				}
				r := serde.ReceiveVoucherResponse{
					Total: total,
					Delta: delta,
				}
				return r, nil
			})
		case serde.GetAddressMethod:
			return processRequest(rs, requestData, func(req serde.NoPayloadRequest) (string, error) {
				return rs.node.Address.Hex(), nil
			})
		case serde.VersionMethod:
			return processRequest(rs, requestData, func(req serde.NoPayloadRequest) (string, error) {
				return rs.node.Version(), nil
			})
		case serde.CreateLedgerChannelRequestMethod:
			return processRequest(rs, requestData, func(req directfund.ObjectiveRequest) (directfund.ObjectiveResponse, error) {
				return rs.node.CreateLedgerChannel(req.CounterParty, req.ChallengeDuration, req.Outcome)
			})
		case serde.CloseLedgerChannelRequestMethod:
			return processRequest(rs, requestData, func(req directdefund.ObjectiveRequest) (protocols.ObjectiveId, error) {
				return rs.node.CloseLedgerChannel(req.ChannelId)
			})
		case serde.CreatePaymentChannelRequestMethod:
			return processRequest(rs, requestData, func(req virtualfund.ObjectiveRequest) (virtualfund.ObjectiveResponse, error) {
				return rs.node.CreatePaymentChannel(req.Intermediaries, req.CounterParty, req.ChallengeDuration, req.Outcome)
			})
		case serde.ClosePaymentChannelRequestMethod:
			return processRequest(rs, requestData, func(req virtualdefund.ObjectiveRequest) (protocols.ObjectiveId, error) {
				return rs.node.ClosePaymentChannel(req.ChannelId)
			})
		case serde.PayRequestMethod:
			return processRequest(rs, requestData, func(req serde.PaymentRequest) (serde.PaymentRequest, error) {
				if err := serde.ValidatePaymentRequest(req); err != nil {
					return serde.PaymentRequest{}, err
				}
				rs.node.Pay(req.Channel, big.NewInt(int64(req.Amount)))
				return req, nil
			})
		case serde.GetPaymentChannelRequestMethod:
			return processRequest(rs, requestData, func(req serde.GetPaymentChannelRequest) (query.PaymentChannelInfo, error) {
				if err := serde.ValidateGetPaymentChannelRequest(req); err != nil {
					return query.PaymentChannelInfo{}, err
				}
				return rs.node.GetPaymentChannel(req.Id)
			})
		case serde.GetLedgerChannelRequestMethod:
			return processRequest(rs, requestData, func(req serde.GetLedgerChannelRequest) (query.LedgerChannelInfo, error) {
				return rs.node.GetLedgerChannel(req.Id)
			})
		case serde.GetAllLedgerChannelsMethod:
			return processRequest(rs, requestData, func(req serde.NoPayloadRequest) ([]query.LedgerChannelInfo, error) {
				return rs.node.GetAllLedgerChannels()
			})
		case serde.GetPaymentChannelsByLedgerMethod:
			return processRequest(rs, requestData, func(req serde.GetPaymentChannelsByLedgerRequest) ([]query.PaymentChannelInfo, error) {
				if err := serde.ValidateGetPaymentChannelsByLedgerRequest(req); err != nil {
					return []query.PaymentChannelInfo{}, err
				}
				return rs.node.GetPaymentChannelsByLedger(req.LedgerId)
			})
		default:
			errRes := serde.NewJsonRpcErrorResponse(jsonrpcReq.Id, serde.MethodNotFoundError)
			return marshalResponse(errRes, rs.logger)
		}
	}

	err = rs.transport.RegisterRequestHandler("v1", handlerV1)
	return err
}

func processRequest[T serde.RequestPayload, U serde.ResponsePayload](rs *RpcServer, requestData []byte, processPayload func(T) (U, error)) []byte {
	rpcRequest := serde.JsonRpcSpecificRequest[T]{}
	// This unmarshal will fail only when the requestData is not valid json.
	// Request-specific params validation is optionally performed as part of the processPayload function
	err := json.Unmarshal(requestData, &rpcRequest)
	if err != nil {
		response := serde.NewJsonRpcErrorResponse(rpcRequest.Id, serde.ParamsUnmarshalError)
		return marshalResponse(response, rs.logger)
	}

	payload := rpcRequest.Params
	processedResponse, err := processPayload(payload)
	if err != nil {
		responseErr := serde.InternalServerError // default error

		if jsonErr, ok := err.(serde.JsonRpcError); ok {
			// overwrite defaults if error object contains jsonrpc error fields
			responseErr.Code = jsonErr.Code
			responseErr.Message = jsonErr.Message
		}

		response := serde.NewJsonRpcErrorResponse(rpcRequest.Id, responseErr)
		return marshalResponse(response, rs.logger)
	}

	response := serde.NewJsonRpcResponse(rpcRequest.Id, processedResponse)
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

func validateJsonrpcRequest(requestData []byte, logger *zerolog.Logger) (serde.JsonRpcGeneralRequest, []byte) {
	var request map[string]interface{}
	vr := serde.JsonRpcGeneralRequest{}
	err := json.Unmarshal(requestData, &request)
	if err != nil {
		errRes := serde.NewJsonRpcErrorResponse(0, serde.RequestUnmarshalError)
		return serde.JsonRpcGeneralRequest{}, marshalResponse(errRes, logger)
	}

	// jsonrpc spec says id can be a string, number.
	// We only support numbers: https://github.com/statechannels/go-nitro/issues/1160
	// When golang unmarshals JSON into an interface value, float64 is used for numbers.
	requestId := request["id"]
	fRequestId, ok := requestId.(float64)
	if !ok || fRequestId != float64(uint64(fRequestId)) {
		errRes := serde.NewJsonRpcErrorResponse(0, serde.InvalidRequestError)
		return serde.JsonRpcGeneralRequest{}, marshalResponse(errRes, logger)
	}
	vr.Id = uint64(fRequestId)

	sJsonrpc, ok := request["jsonrpc"].(string)
	if !ok || sJsonrpc != "2.0" {
		errRes := serde.NewJsonRpcErrorResponse(vr.Id, serde.InvalidRequestError)
		return serde.JsonRpcGeneralRequest{}, marshalResponse(errRes, logger)
	}

	sMethod, ok := request["method"].(string)
	if !ok {
		errRes := serde.NewJsonRpcErrorResponse(vr.Id, serde.InvalidRequestError)
		return serde.JsonRpcGeneralRequest{}, marshalResponse(errRes, logger)
	}
	vr.Method = sMethod

	return vr, nil
}

func (rs *RpcServer) sendNotifications(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			rs.wg.Done()
			return

		case completedObjective, ok := <-rs.node.CompletedObjectives():
			if !ok {
				rs.logger.Warn().Msg("CompletedObjectives channel closed, exiting sendNotifications")
				return
			}
			err := sendNotification(rs, serde.ObjectiveCompleted, completedObjective)
			if err != nil {
				panic(err)
			}
		case ledgerInfo, ok := <-rs.node.LedgerUpdates():
			if !ok {
				rs.logger.Warn().Msg("LedgerUpdates channel closed, exiting sendNotifications")
				return
			}
			err := sendNotification(rs, serde.LedgerChannelUpdated, ledgerInfo)
			if err != nil {
				panic(err)
			}
		case paymentInfo, ok := <-rs.node.PaymentUpdates():
			if !ok {
				rs.logger.Warn().Msg("PaymentUpdates channel closed, exiting sendNotifications")
				return
			}
			err := sendNotification(rs, serde.PaymentChannelUpdated, paymentInfo)
			if err != nil {
				panic(err)
			}
		}
	}
}

func sendNotification[T serde.NotificationMethod, U serde.NotificationPayload](rs *RpcServer, method T, payload U) error {
	rs.logger.Trace().Msgf("Sending notification: %+v", payload)
	request := serde.NewJsonRpcSpecificRequest(rand.Uint64(), method, payload)
	data, err := json.Marshal(request)
	if err != nil {
		return err
	}
	return rs.transport.Notify(data)
}
