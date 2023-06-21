package rpc

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/statechannels/go-nitro/client/query"
	"github.com/statechannels/go-nitro/internal/safesync"
	"github.com/statechannels/go-nitro/payments"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/protocols/directdefund"
	"github.com/statechannels/go-nitro/protocols/directfund"
	"github.com/statechannels/go-nitro/protocols/virtualdefund"
	"github.com/statechannels/go-nitro/protocols/virtualfund"
	"github.com/statechannels/go-nitro/rpc/serde"
	"github.com/statechannels/go-nitro/rpc/transport"
	"github.com/statechannels/go-nitro/rpc/transport/ws"

	"github.com/statechannels/go-nitro/types"

	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/rand"
)

// RpcClient is a client for making nitro rpc calls
type RpcClient struct {
	transport             transport.Requester
	logger                zerolog.Logger
	completedObjectives   *safesync.Map[chan struct{}]
	ledgerChannelUpdates  *safesync.Map[chan query.LedgerChannelInfo]
	paymentChannelUpdates *safesync.Map[chan query.PaymentChannelInfo]
}

// response includes a payload or an error.
type response[T serde.ResponsePayload] struct {
	Payload T
	Error   error
}

// NewRpcClient creates a new RpcClient
func NewRpcClient(rpcServerUrl string, logger zerolog.Logger, trans transport.Requester) (*RpcClient, error) {
	c := &RpcClient{trans, logger, &safesync.Map[chan struct{}]{}, &safesync.Map[chan query.LedgerChannelInfo]{}, &safesync.Map[chan query.PaymentChannelInfo]{}}
	err := c.subscribeToNotifications()
	if err != nil {
		return nil, err
	}
	return c, nil
}

// NewHttpRpcClient creates a new RpcClient using an http transport
func NewHttpRpcClient(rpcServerUrl string) (*RpcClient, error) {
	transport, err := ws.NewWebSocketTransportAsClient(rpcServerUrl)
	if err != nil {
		return nil, err
	}
	c := &RpcClient{transport, zerolog.New(os.Stdout), &safesync.Map[chan struct{}]{}, &safesync.Map[chan query.LedgerChannelInfo]{}, &safesync.Map[chan query.PaymentChannelInfo]{}}
	err = c.subscribeToNotifications()
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (rc *RpcClient) GetVirtualChannel(id types.Destination) query.PaymentChannelInfo {
	req := serde.GetPaymentChannelRequest{Id: id}

	return waitForRequest[serde.GetPaymentChannelRequest, query.PaymentChannelInfo](rc, serde.GetPaymentChannelRequestMethod, req)
}

// CreateLedger creates a new ledger channel
func (rc *RpcClient) CreateVirtual(intermediaries []types.Address, counterparty types.Address, ChallengeDuration uint32, outcome outcome.Exit) virtualfund.ObjectiveResponse {
	objReq := virtualfund.NewObjectiveRequest(
		intermediaries,
		counterparty,
		100,
		outcome,
		rand.Uint64(),
		common.Address{})

	return waitForRequest[virtualfund.ObjectiveRequest, virtualfund.ObjectiveResponse](rc, serde.VirtualFundRequestMethod, objReq)
}

// Pay uses the specified channel to pay the specified amount
func (rc *RpcClient) Pay(id types.Destination, amount uint64) {
	pReq := serde.PaymentRequest{Amount: amount, Channel: id}

	waitForRequest[serde.PaymentRequest, serde.PaymentRequest](rc, serde.PayRequestMethod, pReq)
}

// CreatePayment creates a voucher that can be redeemed against the given
// channel id for a gain of the specified amount.
//
// The returned voucher must be sent to the counterparty, who can redeem it.
func (rc *RpcClient) CreatePayment(id types.Destination, amount uint64) payments.Voucher {
	pReq := serde.PaymentRequest{Amount: amount, Channel: id}

	return waitForRequest[serde.PaymentRequest, payments.Voucher](rc, serde.CreatePaymentMethod, pReq)
}

// ReceivePayment receives a voucher and forwards it to go-nitro.
func (rc *RpcClient) ReceivePayment(voucher payments.Voucher) query.PaymentChannelPaymentReceipt {
	return waitForRequest[serde.ReceivePaymentRequest, query.PaymentChannelPaymentReceipt](rc, serde.ReceiveVoucherRequestMethod, voucher)
}

// CloseVirtual closes a virtual channel
func (rc *RpcClient) CloseVirtual(id types.Destination) protocols.ObjectiveId {
	objReq := virtualdefund.NewObjectiveRequest(
		id)

	return waitForRequest[virtualdefund.ObjectiveRequest, protocols.ObjectiveId](rc, serde.VirtualDefundRequestMethod, objReq)
}

func (rc *RpcClient) GetLedgerChannel(id types.Destination) query.LedgerChannelInfo {
	req := serde.GetLedgerChannelRequest{Id: id}

	return waitForRequest[serde.GetLedgerChannelRequest, query.LedgerChannelInfo](rc, serde.GetLedgerChannelRequestMethod, req)
}

// GetAllLedgerChannels returns all ledger channels
func (rc *RpcClient) GetAllLedgerChannels() []query.LedgerChannelInfo {
	return waitForRequest[serde.NoPayloadRequest, []query.LedgerChannelInfo](rc, serde.GetAllLedgerChannelsMethod, struct{}{})
}

// GetPaymentChannelsByLedger returns all active payment channels for a given ledger channel
func (rc *RpcClient) GetPaymentChannelsByLedger(ledgerId types.Destination) []query.PaymentChannelInfo {
	return waitForRequest[serde.GetPaymentChannelsByLedgerRequest, []query.PaymentChannelInfo](rc, serde.GetPaymentChannelsByLedgerMethod, serde.GetPaymentChannelsByLedgerRequest{LedgerId: ledgerId})
}

// CreateLedger creates a new ledger channel
func (rc *RpcClient) CreateLedger(counterparty types.Address, ChallengeDuration uint32, outcome outcome.Exit) directfund.ObjectiveResponse {
	objReq := directfund.NewObjectiveRequest(
		counterparty,
		100,
		outcome,
		rand.Uint64(),
		common.Address{})

	return waitForRequest[directfund.ObjectiveRequest, directfund.ObjectiveResponse](rc, serde.DirectFundRequestMethod, objReq)
}

// CloseLedger closes a ledger channel
func (rc *RpcClient) CloseLedger(id types.Destination) protocols.ObjectiveId {
	objReq := directdefund.NewObjectiveRequest(id)

	return waitForRequest[directdefund.ObjectiveRequest, protocols.ObjectiveId](rc, serde.DirectDefundRequestMethod, objReq)
}

func (rc *RpcClient) Close() {
	rc.transport.Close()
}

func (rc *RpcClient) subscribeToNotifications() error {
	notificationChan, err := rc.transport.Subscribe()
	rc.logger.Trace().Msg("Subscribed to notifications")
	go func() {
		for data := range notificationChan {
			rc.logger.Trace().Bytes("data", data).Msg("Received notification")
			method, err := getNotificationMethod(data)
			if err != nil {
				panic(err)
			}
			switch method {
			case serde.ObjectiveCompleted:
				rpcRequest := serde.JsonRpcRequest[protocols.ObjectiveId]{}
				err := json.Unmarshal(data, &rpcRequest)
				if err != nil {
					panic(err)
				}
				c, _ := rc.completedObjectives.LoadOrStore(string(rpcRequest.Params), make(chan struct{}))
				close(c)
			case serde.LedgerChannelUpdated:
				rpcRequest := serde.JsonRpcRequest[query.LedgerChannelInfo]{}
				err := json.Unmarshal(data, &rpcRequest)
				if err != nil {
					panic(err)
				}
				c, _ := rc.ledgerChannelUpdates.LoadOrStore(string(rpcRequest.Params.ID.String()), make(chan query.LedgerChannelInfo, 100))
				c <- rpcRequest.Params

			case serde.PaymentChannelUpdated:
				rpcRequest := serde.JsonRpcRequest[query.PaymentChannelInfo]{}
				err := json.Unmarshal(data, &rpcRequest)
				if err != nil {
					panic(err)
				}
				c, _ := rc.paymentChannelUpdates.LoadOrStore(string(rpcRequest.Params.ID.String()), make(chan query.PaymentChannelInfo, 100))
				c <- rpcRequest.Params
			}

		}
	}()
	return err
}

func waitForRequest[T serde.RequestPayload, U serde.ResponsePayload](rc *RpcClient, method serde.RequestMethod, requestData T) U {
	resChan, err := sendRPCRequest[T, U](method, requestData, rc.transport, rc.logger)
	if err != nil {
		panic(err)
	}

	res := <-resChan
	if res.Error != nil {
		panic(res.Error)
	}

	return res.Payload
}

// ObjectiveCompleteChan returns a chan that receives an empty struct when the objective with given id is completed
func (rc *RpcClient) ObjectiveCompleteChan(id protocols.ObjectiveId) <-chan struct{} {
	c, _ := rc.completedObjectives.LoadOrStore(string(id), make(chan struct{}))
	return c
}

// LedgerChannelUpdatesChan returns a chan that receives ledger channel updates.
func (rc *RpcClient) LedgerChannelUpdatesChan(ledgerChannelId types.Destination) <-chan query.LedgerChannelInfo {
	c, _ := rc.ledgerChannelUpdates.LoadOrStore(string(ledgerChannelId.String()), make(chan query.LedgerChannelInfo, 100))
	return c
}

// PaymentChannelUpdatesChan returns a chan that receives payment channel updates.
func (rc *RpcClient) PaymentChannelUpdatesChan(paymentChannelId types.Destination) <-chan query.PaymentChannelInfo {
	c, _ := rc.paymentChannelUpdates.LoadOrStore(string(paymentChannelId.String()), make(chan query.PaymentChannelInfo, 100))
	return c
}

// sendRPCRequest uses the supplied transport and payload to send a non-blocking JSONRPC request.
// It returns a channel that sends a response payload. If the request fails to send, an error is returned.
func sendRPCRequest[T serde.RequestPayload, U serde.ResponsePayload](method serde.RequestMethod, request T, trans transport.Requester, logger zerolog.Logger) (<-chan response[U], error) {
	returnChan := make(chan response[U], 1)

	requestId := rand.Uint64()
	message := serde.NewJsonRpcRequest(requestId, method, request)
	data, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	logger.Trace().
		Str("method", string(method)).
		Msg("sent message")

	go func() {
		responseData, err := trans.Request(data)
		if err != nil {
			returnChan <- response[U]{Error: err}
		}

		logger.Trace().Msgf("Rpc client received response: %+v", string(responseData))

		jsonResponse := serde.JsonRpcResponse[U]{}
		err = json.Unmarshal(responseData, &jsonResponse)
		if err != nil {
			returnChan <- response[U]{Error: err}
		}

		returnChan <- response[U]{jsonResponse.Result, nil}
	}()
	return returnChan, nil
}

// getNotificationMethod parses the raw notification and returns the notification method
func getNotificationMethod(raw []byte) (serde.NotificationMethod, error) {
	var notif map[string]interface{}

	err := json.Unmarshal(raw, &notif)
	if err != nil {
		return "", err
	}

	method, ok := notif["method"].(string)
	if !ok {
		return "", fmt.Errorf("method not found in notification")
	}
	return serde.NotificationMethod(method), nil
}
