package rpc

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/statechannels/go-nitro/client/engine/store/safesync"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/protocols/directdefund"
	"github.com/statechannels/go-nitro/protocols/directfund"
	"github.com/statechannels/go-nitro/protocols/virtualdefund"
	"github.com/statechannels/go-nitro/protocols/virtualfund"
	"github.com/statechannels/go-nitro/rpc/serde"
	"github.com/statechannels/go-nitro/rpc/transport"

	"github.com/statechannels/go-nitro/types"

	"github.com/statechannels/go-nitro/channel/state/outcome"
)

// RpcClient is a client for making nitro rpc calls
type RpcClient struct {
	transport                    transport.Requester
	myAddress                    types.Address
	logger                       zerolog.Logger
	completedObjectivesListeners map[protocols.ObjectiveId][]chan struct{}
}

// response includes a payload or an error.
type response[T serde.ResponsePayload] struct {
	Payload T
	Error   error
}

// NewRpcClient creates a new RpcClient
func NewRpcClient(rpcServerUrl string, myAddress types.Address, logger zerolog.Logger, trans transport.Requester) (*RpcClient, error) {
	c := &RpcClient{trans, myAddress, logger, make(map[protocols.ObjectiveId][]chan struct{})}
	err := c.subscribeToNotifications()
	if err != nil {
		return nil, err
	}
	return c, nil
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

	return waitForRequest[virtualfund.ObjectiveRequest, virtualfund.ObjectiveResponse](rc, objReq)
}

// CloseVirtual closes a virtual channel
func (rc *RpcClient) CloseVirtual(id types.Destination) protocols.ObjectiveId {
	objReq := virtualdefund.NewObjectiveRequest(
		id)

	return waitForRequest[virtualdefund.ObjectiveRequest, protocols.ObjectiveId](rc, objReq)
}

// CreateLedger creates a new ledger channel
func (rc *RpcClient) CreateLedger(counterparty types.Address, ChallengeDuration uint32, outcome outcome.Exit) directfund.ObjectiveResponse {
	objReq := directfund.NewObjectiveRequest(
		counterparty,
		100,
		outcome,
		rand.Uint64(),
		common.Address{})

	return waitForRequest[directfund.ObjectiveRequest, directfund.ObjectiveResponse](rc, objReq)
}

// CloseLedger closes a ledger channel
func (rc *RpcClient) CloseLedger(id types.Destination) protocols.ObjectiveId {
	objReq := directdefund.NewObjectiveRequest(id)

	return waitForRequest[directdefund.ObjectiveRequest, protocols.ObjectiveId](rc, objReq)
}

// Pay uses the specified channel to pay the specified amount
func (rc *RpcClient) Pay(id types.Destination, amount uint64) {
	pReq := serde.PaymentRequest{Amount: amount, Channel: id}

	waitForRequest[serde.PaymentRequest, serde.PaymentRequest](rc, pReq)
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
			rpcRequest := serde.JsonRpcRequest[protocols.ObjectiveId]{}
			err := json.Unmarshal(data, &rpcRequest)
			if err != nil {
				panic(err)
			}
			for _, ch := range rc.completedObjectivesListeners[rpcRequest.Params] {
				ch <- struct{}{}
			}
			delete(rc.completedObjectivesListeners, rpcRequest.Params)
		}
	}()
	return err
}

func waitForRequest[T serde.RequestPayload, U serde.ResponsePayload](rc *RpcClient, requestData T) U {
	resChan, err := request[T, U](rc.transport, requestData, rc.logger)
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
	ch := make(chan struct{}, 1) // use a buffer of 1 so we can send on it without blocking (in case the consumer isn't ready)
	rc.completedObjectivesListeners[id] = append(rc.completedObjectivesListeners[id], ch)
	return ch
}

func (rc *RpcClient) WaitForObjectiveCompletion(expectedObjectiveId ...protocols.ObjectiveId) {
	incomplete := safesync.Map[<-chan struct{}]{}

	var wg sync.WaitGroup

	for _, id := range expectedObjectiveId {
		incomplete.Store(string(id), rc.ObjectiveCompleteChan(id))
		wg.Add(1)
	}

	incomplete.Range(
		func(id string, ch <-chan struct{}) bool {
			go func() {
				<-ch
				incomplete.Delete(string(id))
				wg.Done()
			}()
			return true
		})

	wg.Wait()

}

// request uses the supplied transport and payload to send a non-blocking JSONRPC request.
// It returns a channel that sends a response payload. If the request fails to send, an error is returned.
func request[T serde.RequestPayload, U serde.ResponsePayload](trans transport.Requester, request T, logger zerolog.Logger) (<-chan response[U], error) {
	returnChan := make(chan response[U], 1)

	var method serde.RequestMethod
	switch any(request).(type) {
	case directfund.ObjectiveRequest:
		method = serde.DirectFundRequestMethod
	case directdefund.ObjectiveRequest:
		method = serde.DirectDefundRequestMethod
	case virtualfund.ObjectiveRequest:
		method = serde.VirtualFundRequestMethod
	case virtualdefund.ObjectiveRequest:
		method = serde.VirtualDefundRequestMethod
	case serde.PaymentRequest:
		method = serde.PayRequestMethod
	default:
		return nil, fmt.Errorf("unknown request type %v", request)
	}
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

		logger.Trace().Msgf("Rpc client received response: %+v", responseData)

		jsonResponse := serde.JsonRpcResponse[U]{}
		err = json.Unmarshal(responseData, &jsonResponse)
		if err != nil {
			returnChan <- response[U]{Error: err}
		}

		returnChan <- response[U]{jsonResponse.Result, nil}
	}()

	return returnChan, nil
}
