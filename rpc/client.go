package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"

	"github.com/ethereum/go-ethereum/common"

	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/internal/logging"
	"github.com/statechannels/go-nitro/internal/safesync"
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
	"github.com/statechannels/go-nitro/rpc/transport/ws"
	"github.com/statechannels/go-nitro/types"
)

// RpcClientApi provides various functions to make RPC API calls to a nitro RPC server
type RpcClientApi interface {
	// Address returns the address of the nitro node
	Address() (common.Address, error)

	// CreateVoucher creates a voucher for the given channelId and amount and returns it.
	// It is the responsibility of the caller to send the voucher to the payee.
	CreateVoucher(chId types.Destination, amount uint64) (payments.Voucher, error)

	// ReceiveVoucher receives a voucher and adds it to the go-nitro store.
	// It returns the total amount received so far and the amount received from the voucher supplied.
	// It can be used to add a voucher that was sent outside of the go-nitro system.
	ReceiveVoucher(v payments.Voucher) (payments.ReceiveVoucherSummary, error)

	// GetPaymentChannel returns the payment channel information for the given channelId
	GetPaymentChannel(chId types.Destination) (query.PaymentChannelInfo, error)

	// CreatePaymentChannel creates a new virtual payment channel with the specified intermediaries, counterparty, ChallengeDuration, and outcome
	CreatePaymentChannel(intermediaries []types.Address, counterparty types.Address, ChallengeDuration uint32, outcome outcome.Exit) (virtualfund.ObjectiveResponse, error)

	// ClosePaymentChannel attempts to close the payment channel with the specified channelId
	ClosePaymentChannel(id types.Destination) (protocols.ObjectiveId, error)

	// GetLedgerChannel returns the ledger channel information for the given channelId
	GetLedgerChannel(id types.Destination) (query.LedgerChannelInfo, error)

	// GetAllLedgerChannels returns information about all ledger channels
	GetAllLedgerChannels() ([]query.LedgerChannelInfo, error)

	// GetPaymentChannelsByLedger returns all active payment channels for a given ledger channel
	GetPaymentChannelsByLedger(ledgerId types.Destination) ([]query.PaymentChannelInfo, error)

	// CreateLedgerChannel creates a new ledger channel with the specified counterparty, ChallengeDuration, and outcome
	CreateLedgerChannel(counterparty types.Address, ChallengeDuration uint32, outcome outcome.Exit) (directfund.ObjectiveResponse, error)

	// CloseLedgerChannel attempts to close the ledger channel with the specified channelId
	CloseLedgerChannel(id types.Destination) (protocols.ObjectiveId, error)

	// Pay uses the specified channel to pay the specified amount
	Pay(id types.Destination, amount uint64) (serde.PaymentRequest, error)

	// Close shuts down the RpcClient and closes the underlying transport
	Close() error

	// ObjectiveCompleteChan returns a channel that receives an empty struct when the objective with the given id is completed
	ObjectiveCompleteChan(id protocols.ObjectiveId) <-chan struct{}

	// LedgerChannelUpdatesChan returns a channel that receives ledger channel updates for the given ledger channel id
	LedgerChannelUpdatesChan(ledgerChannelId types.Destination) <-chan query.LedgerChannelInfo

	// PaymentChannelUpdatesChan returns a channel that receives payment channel updates for the given payment channel id
	PaymentChannelUpdatesChan(paymentChannelId types.Destination) <-chan query.PaymentChannelInfo
}

// rpcClient is the implementation
type rpcClient struct {
	transport             transport.Requester
	completedObjectives   *safesync.Map[chan struct{}]
	ledgerChannelUpdates  *safesync.Map[chan query.LedgerChannelInfo]
	paymentChannelUpdates *safesync.Map[chan query.PaymentChannelInfo]
	cancel                context.CancelFunc
	wg                    *sync.WaitGroup
	nodeAddress           common.Address
	logger                *slog.Logger
}

// response includes a payload or an error.
type response[T serde.ResponsePayload] struct {
	Payload T
	Error   error
}

// NewRpcClient creates a new RpcClient
func NewRpcClient(trans transport.Requester) (RpcClientApi, error) {
	ctx, cancel := context.WithCancel(context.Background())

	c := &rpcClient{trans, &safesync.Map[chan struct{}]{}, &safesync.Map[chan query.LedgerChannelInfo]{}, &safesync.Map[chan query.PaymentChannelInfo]{}, cancel, &sync.WaitGroup{}, common.Address{}, slog.Default()}

	// Retrieve the address and set it on the rpcClient
	res, err := waitForRequest[serde.NoPayloadRequest, common.Address](c, serde.GetAddressMethod, serde.NoPayloadRequest{})
	if err != nil {
		return nil, err
	}
	c.nodeAddress = res

	// Update the logger so we output the address
	c.logger = logging.LoggerWithAddress(c.logger, c.nodeAddress)

	notificationChan, err := c.transport.Subscribe()
	if err != nil {
		return nil, err
	}
	c.wg.Add(1)
	go c.subscribeToNotifications(ctx, notificationChan)

	return c, nil
}

// NewHttpRpcClient creates a new rpcClient using an http transport
func NewHttpRpcClient(rpcServerUrl string) (RpcClientApi, error) {
	transport, err := ws.NewWebSocketTransportAsClient(rpcServerUrl)
	if err != nil {
		return nil, err
	}
	return NewRpcClient(transport)
}

// Address returns the address of the the nitro node
func (rc *rpcClient) Address() (common.Address, error) {
	return rc.nodeAddress, nil
}

// CreateVoucher creates a voucher for the given channelId and amount and returns it.
// It is the responsibility of the caller to send the voucher to the payee.
func (rc *rpcClient) CreateVoucher(chId types.Destination, amount uint64) (payments.Voucher, error) {
	req := serde.PaymentRequest{Channel: chId, Amount: amount}
	return waitForRequest[serde.PaymentRequest, payments.Voucher](rc, serde.CreateVoucherRequestMethod, req)
}

// ReceiveVoucher receives a voucher and adds it to the go-nitro store.
// It returns the total amount received so far and the amount received from the voucher supplied.
// It can be used to add a voucher that was sent outside of the go-nitro system.
func (rc *rpcClient) ReceiveVoucher(v payments.Voucher) (payments.ReceiveVoucherSummary, error) {
	return waitForRequest[payments.Voucher, payments.ReceiveVoucherSummary](rc, serde.ReceiveVoucherRequestMethod, v)
}

func (rc *rpcClient) GetPaymentChannel(chId types.Destination) (query.PaymentChannelInfo, error) {
	req := serde.GetPaymentChannelRequest{Id: chId}

	return waitForRequest[serde.GetPaymentChannelRequest, query.PaymentChannelInfo](rc, serde.GetPaymentChannelRequestMethod, req)
}

// CreatePaymentChannel creates a new virtual payment channel
func (rc *rpcClient) CreatePaymentChannel(intermediaries []types.Address, counterparty types.Address, ChallengeDuration uint32, outcome outcome.Exit) (virtualfund.ObjectiveResponse, error) {
	objReq := virtualfund.NewObjectiveRequest(
		intermediaries,
		counterparty,
		100,
		outcome,
		rand.Uint64(),
		common.Address{})

	return waitForRequest[virtualfund.ObjectiveRequest, virtualfund.ObjectiveResponse](rc, serde.CreatePaymentChannelRequestMethod, objReq)
}

// ClosePaymentChannel attempts to close the payment channel with supplied id
func (rc *rpcClient) ClosePaymentChannel(id types.Destination) (protocols.ObjectiveId, error) {
	objReq := virtualdefund.NewObjectiveRequest(
		id)

	return waitForRequest[virtualdefund.ObjectiveRequest, protocols.ObjectiveId](rc, serde.ClosePaymentChannelRequestMethod, objReq)
}

func (rc *rpcClient) GetLedgerChannel(id types.Destination) (query.LedgerChannelInfo, error) {
	req := serde.GetLedgerChannelRequest{Id: id}

	return waitForRequest[serde.GetLedgerChannelRequest, query.LedgerChannelInfo](rc, serde.GetLedgerChannelRequestMethod, req)
}

// GetAllLedgerChannels returns all ledger channels
func (rc *rpcClient) GetAllLedgerChannels() ([]query.LedgerChannelInfo, error) {
	return waitForRequest[serde.NoPayloadRequest, []query.LedgerChannelInfo](rc, serde.GetAllLedgerChannelsMethod, struct{}{})
}

// GetPaymentChannelsByLedger returns all active payment channels for a given ledger channel
func (rc *rpcClient) GetPaymentChannelsByLedger(ledgerId types.Destination) ([]query.PaymentChannelInfo, error) {
	return waitForRequest[serde.GetPaymentChannelsByLedgerRequest, []query.PaymentChannelInfo](rc, serde.GetPaymentChannelsByLedgerMethod, serde.GetPaymentChannelsByLedgerRequest{LedgerId: ledgerId})
}

// CreateLedger creates a new ledger channel
func (rc *rpcClient) CreateLedgerChannel(counterparty types.Address, ChallengeDuration uint32, outcome outcome.Exit) (directfund.ObjectiveResponse, error) {
	objReq := directfund.NewObjectiveRequest(
		counterparty,
		100,
		outcome,
		rand.Uint64(),
		common.Address{})

	return waitForRequest[directfund.ObjectiveRequest, directfund.ObjectiveResponse](rc, serde.CreateLedgerChannelRequestMethod, objReq)
}

// CloseLedger closes a ledger channel
func (rc *rpcClient) CloseLedgerChannel(id types.Destination) (protocols.ObjectiveId, error) {
	objReq := directdefund.NewObjectiveRequest(id)

	return waitForRequest[directdefund.ObjectiveRequest, protocols.ObjectiveId](rc, serde.CloseLedgerChannelRequestMethod, objReq)
}

// Pay uses the specified channel to pay the specified amount
func (rc *rpcClient) Pay(id types.Destination, amount uint64) (serde.PaymentRequest, error) {
	pReq := serde.PaymentRequest{Amount: amount, Channel: id}
	return waitForRequest[serde.PaymentRequest, serde.PaymentRequest](rc, serde.PayRequestMethod, pReq)
}

func (rc *rpcClient) Close() error {
	rc.cancel()
	rc.wg.Wait()
	return rc.transport.Close()
}

func (rc *rpcClient) subscribeToNotifications(ctx context.Context, notificationChan <-chan []byte) {
	rc.logger.Debug("Subscribed to notifications")
	for {
		select {
		case <-ctx.Done():
			rc.wg.Done()
			return
		case data := <-notificationChan:

			method, err := getNotificationMethod(data)
			if err != nil {
				panic(err)
			}
			switch method {
			case serde.ObjectiveCompleted:
				rpcRequest := serde.JsonRpcSpecificRequest[protocols.ObjectiveId]{}
				err := json.Unmarshal(data, &rpcRequest)
				rc.logger.Debug("Received notification", "method", method, "data", rpcRequest)
				if err != nil {
					panic(err)
				}
				c, _ := rc.completedObjectives.LoadOrStore(string(rpcRequest.Params), make(chan struct{}))
				close(c)
			case serde.LedgerChannelUpdated:
				rpcRequest := serde.JsonRpcSpecificRequest[query.LedgerChannelInfo]{}
				err := json.Unmarshal(data, &rpcRequest)
				rc.logger.Debug("Received notification", "method", method, "data", rpcRequest)
				if err != nil {
					panic(err)
				}
				c, _ := rc.ledgerChannelUpdates.LoadOrStore(string(rpcRequest.Params.ID.String()), make(chan query.LedgerChannelInfo, 100))
				c <- rpcRequest.Params

			case serde.PaymentChannelUpdated:
				rpcRequest := serde.JsonRpcSpecificRequest[query.PaymentChannelInfo]{}
				err := json.Unmarshal(data, &rpcRequest)
				rc.logger.Debug("Received notification", "method", method, "data", rpcRequest)
				if err != nil {
					panic(err)
				}
				c, _ := rc.paymentChannelUpdates.LoadOrStore(string(rpcRequest.Params.ID.String()), make(chan query.PaymentChannelInfo, 100))
				c <- rpcRequest.Params
			}

		}
	}
}

// ObjectiveCompleteChan returns a chan that receives an empty struct when the objective with given id is completed
func (rc *rpcClient) ObjectiveCompleteChan(id protocols.ObjectiveId) <-chan struct{} {
	c, _ := rc.completedObjectives.LoadOrStore(string(id), make(chan struct{}))
	return c
}

// LedgerChannelUpdatesChan returns a chan that receives ledger channel updates.
func (rc *rpcClient) LedgerChannelUpdatesChan(ledgerChannelId types.Destination) <-chan query.LedgerChannelInfo {
	c, _ := rc.ledgerChannelUpdates.LoadOrStore(string(ledgerChannelId.String()), make(chan query.LedgerChannelInfo, 100))
	return c
}

// PaymentChannelUpdatesChan returns a chan that receives payment channel updates.
func (rc *rpcClient) PaymentChannelUpdatesChan(paymentChannelId types.Destination) <-chan query.PaymentChannelInfo {
	c, _ := rc.paymentChannelUpdates.LoadOrStore(string(paymentChannelId.String()), make(chan query.PaymentChannelInfo, 100))
	return c
}

func waitForRequest[T serde.RequestPayload, U serde.ResponsePayload](rc *rpcClient, method serde.RequestMethod, requestData T) (U, error) {
	rc.wg.Add(1)
	defer rc.wg.Done()

	res, err := sendRequest[T, U](rc.transport, method, requestData, rc.logger, rc.wg)
	if err != nil {
		panic(err)
	}

	return res.Payload, res.Error
}

// sendRequest uses the supplied transport and payload to send a JSONRPC request.
//   - Returns an error if:
//     [1] the request fails to send
//     [2] the response cannot be parsed
//   - Otherwise, returns the JSONRPC server's response
func sendRequest[T serde.RequestPayload, U serde.ResponsePayload](trans transport.Requester, method serde.RequestMethod, reqPayload T, logger *slog.Logger, wg *sync.WaitGroup) (response[U], error) {
	requestId := rand.Uint64()
	message := serde.NewJsonRpcSpecificRequest(requestId, method, reqPayload)
	data, err := json.Marshal(message)
	if err != nil {
		return response[U]{}, err
	}

	logger.Debug("sent message", "method", string(method))

	responseData, err := trans.Request(data)
	if err != nil {
		return response[U]{}, err
	}

	// First check if there is an error present in the jsonrpc response
	jsonResponse := serde.JsonRpcGeneralResponse{}
	err = json.Unmarshal(responseData, &jsonResponse)
	if err != nil {
		return response[U]{}, err
	} else if jsonResponse.Error != (serde.JsonRpcError{}) {
		return response[U]{Error: jsonResponse.Error}, nil
	}

	// Now convert response.Result into the specific type for this request, and return that
	successResponse := serde.JsonRpcSuccessResponse[U]{}
	err = json.Unmarshal(responseData, &successResponse)
	if err != nil {
		return response[U]{}, err
	}
	return response[U]{Payload: successResponse.Result}, nil
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
