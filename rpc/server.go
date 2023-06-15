package rpc

import (
	"bytes"
	"context"
	"encoding/json"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/rs/zerolog"
	nitro "github.com/statechannels/go-nitro/client"
	"github.com/statechannels/go-nitro/client/query"
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

// RpcServer handles nitro rpc requests and executes them on the nitro client
type RpcServer struct {
	transport transport.Responder
	client    *nitro.Client
	logger    *zerolog.Logger
	cancel    context.CancelFunc
	wg        *sync.WaitGroup
}

func (rs *RpcServer) Url() string {
	return rs.transport.Url()
}

func (rs *RpcServer) Close() error {
	rs.cancel()
	rs.wg.Wait()

	rs.transport.Close()
	return rs.client.Close()
}

// newRpcServerWithoutNotifications creates a new rpc server without notifications enabled
func newRpcServerWithoutNotifications(nitroClient *nitro.Client, logger *zerolog.Logger, trans transport.Responder) (*RpcServer, error) {
	rs := &RpcServer{trans, nitroClient, logger, func() {}, &sync.WaitGroup{}}

	err := rs.registerHandlers()
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func NewRpcServer(nitroClient *nitro.Client, logger *zerolog.Logger, trans transport.Responder) (*RpcServer, error) {
	ctx, cancel := context.WithCancel(context.Background())
	rs := &RpcServer{trans, nitroClient, logger, cancel, &sync.WaitGroup{}}

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
			return marshalResponse(types.ParseError, rs.logger)
		}

		validationResult := validateRequest(requestData, rs.logger)
		if validationResult.Error != nil {
			return validationResult.Error
		}

		switch serde.RequestMethod(validationResult.Method) {
		case serde.GetAddressMethod:
			return processRequest(rs, requestData, func(T serde.NoPayloadRequest) (string, error) {
				return rs.client.Address.Hex(), nil
			})
		case serde.VersionMethod:
			return processRequest(rs, requestData, func(T serde.NoPayloadRequest) (string, error) {
				return rs.client.Version(), nil
			})
		case serde.DirectFundRequestMethod:
			return processRequest(rs, requestData, func(obj directfund.ObjectiveRequest) (directfund.ObjectiveResponse, error) {
				return rs.client.CreateLedgerChannel(obj.CounterParty, obj.ChallengeDuration, obj.Outcome), nil
			})
		case serde.DirectDefundRequestMethod:
			return processRequest(rs, requestData, func(obj directdefund.ObjectiveRequest) (protocols.ObjectiveId, error) {
				return rs.client.CloseLedgerChannel(obj.ChannelId), nil
			})
		case serde.VirtualFundRequestMethod:
			return processRequest(rs, requestData, func(obj virtualfund.ObjectiveRequest) (virtualfund.ObjectiveResponse, error) {
				return rs.client.CreateVirtualPaymentChannel(obj.Intermediaries, obj.CounterParty, obj.ChallengeDuration, obj.Outcome), nil
			})
		case serde.VirtualDefundRequestMethod:
			return processRequest(rs, requestData, func(obj virtualdefund.ObjectiveRequest) (protocols.ObjectiveId, error) {
				return rs.client.CloseVirtualChannel(obj.ChannelId), nil
			})
		case serde.PayRequestMethod:
			return processRequest(rs, requestData, func(obj serde.PaymentRequest) (serde.PaymentRequest, error) {
				rs.client.Pay(obj.Channel, big.NewInt(int64(obj.Amount)))
				return obj, nil
			})
		case serde.CreatePaymentMethod:
			return processRequest(rs, requestData, func(obj serde.PaymentRequest) (payments.Voucher, error) {
				v, err := rs.client.CreatePayment(obj.Channel, big.NewInt(int64(obj.Amount)))
				if err != nil {
					return payments.Voucher{}, err
				}
				return <-v, nil
			})
		case serde.ReceiveVoucherRequestMethod:
			return processRequest(rs, requestData, func(v serde.ReceivePaymentRequest) (query.PaymentChannelPaymentReceipt, error) {
				pc, err := rs.client.GetPaymentChannel(v.ChannelId)
				if err != nil {
					return query.PaymentChannelPaymentReceipt{
						Status: query.PRSchannelNotFound,
					}, err
				}

				me := rs.client.Address

				if !bytes.Equal(me.Bytes(), pc.Balance.Payee.Bytes()) {
					return query.PaymentChannelPaymentReceipt{
						ID:     v.ChannelId,
						Status: query.PRSmisaddressed,
					}, err
				}

				signer, err := v.RecoverSigner()

				if !bytes.Equal(signer.Bytes(), pc.Balance.Payer.Bytes()) || err != nil {
					return query.PaymentChannelPaymentReceipt{
						ID:     v.ChannelId,
						Status: query.PRSincorrectSigner,
					}, err
				}

				amountReceived := big.NewInt(0).Sub(v.Amount, (*big.Int)(pc.Balance.PaidSoFar))
				affords := (*big.Int)(pc.Balance.RemainingFunds).Cmp(amountReceived) >= 0
				if !affords {
					return query.PaymentChannelPaymentReceipt{
						ID:     v.ChannelId,
						Status: query.PRSinsufficientFunds,
					}, err
				}

				// construct the "message" for engine consumption. Engine will use msgService
				// to "send" to itself and then process as any other voucher.
				voucherMessage := protocols.Message{
					To:       *me,
					From:     signer,
					Payments: []payments.Voucher{v},
				}
				rs.client.PushMessage(voucherMessage)

				// return an *optimistic* appraisal of amount received. Engine processing could feasibly
				// fail, but above checks are pretty comprehensive
				return query.PaymentChannelPaymentReceipt{
					ID:             v.ChannelId,
					AmountReceived: (*hexutil.Big)(amountReceived),
					Status:         query.PRSreceived,
				}, nil
			})
		case serde.GetPaymentChannelRequestMethod:
			return processRequest(rs, requestData, func(r serde.GetPaymentChannelRequest) (query.PaymentChannelInfo, error) {
				return rs.client.GetPaymentChannel(r.Id)
			})
		case serde.GetLedgerChannelRequestMethod:
			return processRequest(rs, requestData, func(r serde.GetLedgerChannelRequest) (query.LedgerChannelInfo, error) {
				return rs.client.GetLedgerChannel(r.Id)
			})
		case serde.GetAllLedgerChannelsMethod:
			return processRequest(rs, requestData, func(r serde.NoPayloadRequest) ([]query.LedgerChannelInfo, error) {
				return rs.client.GetAllLedgerChannels()
			})
		case serde.GetPaymentChannelsByLedgerMethod:
			return processRequest(rs, requestData, func(r serde.GetPaymentChannelsByLedgerRequest) ([]query.PaymentChannelInfo, error) {
				return rs.client.GetPaymentChannelsByLedger(r.LedgerId)
			})
		default:
			responseErr := types.MethodNotFoundError
			responseErr.Id = validationResult.Id
			return marshalResponse(responseErr, rs.logger)
		}
	}

	err = rs.transport.RegisterRequestHandler("v1", handlerV1)
	return err
}

func processRequest[T serde.RequestPayload, U serde.ResponsePayload](rs *RpcServer, requestData []byte, processPayload func(T) (U, error)) []byte {
	rpcRequest := serde.JsonRpcRequest[T]{}
	// todo: unmarshal will fail only when the requestData is not valid json.
	// At the moment, there is no validation that the required fields are populated in the request.
	err := json.Unmarshal(requestData, &rpcRequest)
	if err != nil {
		return marshalResponse(types.UnexpectedRequestUnmarshalError2, rs.logger)
	}
	obj := rpcRequest.Params
	objResponse, err := processPayload(obj)
	if err != nil {
		responseErr := types.JsonRpcError{
			Code:    types.InternalServerError.Code, // default error code
			Message: err.Error(),
			Id:      rpcRequest.Id,
		}

		if jsonErr, ok := err.(types.JsonRpcError); ok {
			responseErr.Code = jsonErr.Code // overwrite default if error object contains Code field
		}

		return marshalResponse(responseErr, rs.logger)
	}

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
		vr.Error = marshalResponse(types.UnexpectedRequestUnmarshalError, logger)
		return vr
	}

	// jsonrpc spec says id can be a string, number.
	// We only support numbers: https://github.com/statechannels/go-nitro/issues/1160
	// When golang unmarshals JSON into an interface value, float64 is used for numbers.
	requestId := request["id"]
	fRequestId, ok := requestId.(float64)
	if !ok {
		vr.Error = marshalResponse(types.InvalidRequestError, logger)
		return vr
	}

	if fRequestId != float64(uint64(fRequestId)) {
		vr.Error = marshalResponse(types.InvalidRequestError, logger)
		return vr
	}
	vr.Id = uint64(fRequestId)

	sJsonrpc, ok := request["jsonrpc"].(string)
	if !ok || sJsonrpc != "2.0" {
		requestError := types.InvalidRequestError
		requestError.Id = vr.Id
		vr.Error = marshalResponse(requestError, logger)
		return vr
	}

	_, ok = request["method"].(string)
	if !ok {
		requestError := types.InvalidRequestError
		requestError.Id = vr.Id
		vr.Error = marshalResponse(requestError, logger)
		return vr
	}
	vr.Method = request["method"].(string)

	return vr
}

func (rs *RpcServer) sendNotifications(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			rs.wg.Done()
			return
		case completedObjective, ok := <-rs.client.CompletedObjectives():
			if !ok {
				rs.logger.Warn().Msg("CompletedObjectives channel closed, exiting sendNotifications")
				return
			}
			err := sendNotification(rs, serde.ObjectiveCompleted, completedObjective)
			if err != nil {
				panic(err)
			}
		case ledgerInfo, ok := <-rs.client.LedgerUpdates():
			if !ok {
				rs.logger.Warn().Msg("LedgerUpdates channel closed, exiting sendNotifications")
				return
			}
			err := sendNotification(rs, serde.LedgerChannelUpdated, ledgerInfo)
			if err != nil {
				panic(err)
			}
		case paymentInfo, ok := <-rs.client.PaymentUpdates():
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
	request := serde.NewJsonRpcRequest(rand.Uint64(), method, payload)
	data, err := json.Marshal(request)
	if err != nil {
		return err
	}
	return rs.transport.Notify(data)
}
