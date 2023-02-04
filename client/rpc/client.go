package rpc

import (
	"encoding/json"
	"fmt"
	"math/big"
	"math/rand"

	"github.com/ethereum/go-ethereum/common"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	"github.com/statechannels/go-nitro/client/engine/store/safesync"
	"github.com/statechannels/go-nitro/network"
	"github.com/statechannels/go-nitro/network/serde"
	natstrans "github.com/statechannels/go-nitro/network/transport/nats"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/protocols/directdefund"
	"github.com/statechannels/go-nitro/protocols/directfund"
	"github.com/statechannels/go-nitro/protocols/virtualdefund"
	"github.com/statechannels/go-nitro/protocols/virtualfund"
	"github.com/statechannels/go-nitro/types"

	"github.com/statechannels/go-nitro/channel/state/outcome"
)

// RpcClient is a client for making nitro rpc calls
type RpcClient struct {
	nts       *network.NetworkService
	myAddress types.Address
	chainId   *big.Int

	// responses is a collection of channels that are used to wait until a response is received from the RPC server
	responses safesync.Map[chan interface{}]

	idsToMethods safesync.Map[serde.RequestMethod]
}

// NewRpcClient creates a new RpcClient
func NewRpcClient(rpcServerUrl string, myAddress types.Address, chainId *big.Int, logger zerolog.Logger) *RpcClient {

	nc, err := nats.Connect(rpcServerUrl)
	handleError(err)
	trp := natstrans.NewNatsTransport(nc, getTopics())

	con, err := trp.PollConnection()
	handleError(err)
	nts := network.NewNetworkService(con)
	nts.Logger = logger
	c := &RpcClient{nts, myAddress, chainId, safesync.Map[chan interface{}]{}, safesync.Map[serde.RequestMethod]{}}
	c.registerHandlers()
	return c
}

// CreateLedger creates a new ledger channel
func (rc *RpcClient) CreateVirtual(intermediaries []types.Address, counterparty types.Address, ChallengeDuration uint32, outcome outcome.Exit) virtualfund.ObjectiveResponse {

	objReq := virtualfund.NewObjectiveRequest(
		intermediaries,
		counterparty,
		100,
		outcome,
		uint64(rand.Float64()), // TODO: Since numeric fields get converted to a float64 in transit we need to prevent overflow
		common.Address{})

	// Create a channel and store it in the responses map
	// We will use this channel to wait for the response
	resRec := make(chan interface{})
	rc.responses.Store(string(objReq.Id(rc.myAddress, rc.chainId)), resRec)

	requestId := rand.Uint64()
	rc.idsToMethods.Store(string(fmt.Sprintf("%d", requestId)), serde.VirtualFundRequestMethod)

	message := serde.NewJsonRpcRequest(requestId, serde.VirtualFundRequestMethod, objReq)
	data, err := json.Marshal(message)
	if err != nil {
		panic("Could not marshal direct fund request")
	}
	rc.nts.SendMessage(serde.VirtualFundRequestMethod, data)

	objRes := <-resRec
	return objRes.(virtualfund.ObjectiveResponse)
}

// CloseVirtual closes a virtual channel
func (rc *RpcClient) CloseVirtual(id types.Destination) protocols.ObjectiveId {
	objReq := virtualdefund.NewObjectiveRequest(
		id)

	// Create a channel and store it in the responses map
	// We will use this channel to wait for the response
	resRec := make(chan interface{})
	rc.responses.Store(string(objReq.Id(rc.myAddress, rc.chainId)), resRec)
	requestId := rand.Uint64()
	rc.idsToMethods.Store(string(fmt.Sprintf("%d", requestId)), serde.VirtualDefundRequestMethod)

	message := serde.NewJsonRpcRequest(requestId, serde.VirtualDefundRequestMethod, objReq)
	data, err := json.Marshal(message)
	if err != nil {
		panic("Could not marshal direct fund request")
	}
	rc.nts.SendMessage(serde.VirtualDefundRequestMethod, data)

	objRes := <-resRec
	return objRes.(protocols.ObjectiveId)
}

// CreateLedger creates a new ledger channel
func (rc *RpcClient) CreateLedger(counterparty types.Address, ChallengeDuration uint32, outcome outcome.Exit) directfund.ObjectiveResponse {

	objReq := directfund.NewObjectiveRequest(
		counterparty,
		100,
		outcome,
		uint64(rand.Float64()), // TODO: Since numeric fields get converted to a float64 in transit we need to prevent overflow
		common.Address{})

	// Create a channel and store it in the responses map
	// We will use this channel to wait for the response
	resRec := make(chan interface{})
	rc.responses.Store(string(objReq.Id(rc.myAddress, rc.chainId)), resRec)

	requestId := rand.Uint64()
	rc.idsToMethods.Store(fmt.Sprintf("%d", requestId), serde.DirectFundRequestMethod)

	message := serde.NewJsonRpcRequest(requestId, serde.DirectFundRequestMethod, objReq)
	data, err := json.Marshal(message)
	if err != nil {
		panic("Could not marshal direct fund request")
	}
	rc.nts.SendMessage(serde.DirectFundRequestMethod, data)

	objRes := <-resRec
	return objRes.(directfund.ObjectiveResponse)
}

// CloseLedger closes a ledger channel
func (rc *RpcClient) CloseLedger(id types.Destination) protocols.ObjectiveId {
	objReq := directdefund.NewObjectiveRequest(id)

	// Create a channel and store it in the responses map
	// We will use this channel to wait for the response
	resRec := make(chan interface{})
	rc.responses.Store(string(objReq.Id(rc.myAddress, rc.chainId)), resRec)
	requestId := rand.Uint64()
	rc.idsToMethods.Store(fmt.Sprintf("%d", requestId), serde.DirectDefundRequestMethod)

	message := serde.NewJsonRpcRequest(requestId, serde.DirectDefundRequestMethod, objReq)
	data, err := json.Marshal(message)
	if err != nil {
		panic("Could not marshal direct fund request")
	}
	rc.nts.SendMessage(serde.DirectDefundRequestMethod, data)

	objRes := <-resRec
	return objRes.(protocols.ObjectiveId)
}

// Pay uses the specified channel to pay the specified amount
func (rc *RpcClient) Pay(id types.Destination, amount uint64) {
	// Create a channel and store it in the responses map
	// We will use this channel to wait for the response
	resRec := make(chan interface{})
	paymentId := fmt.Sprintf("PAYMENT-%s", id)
	rc.responses.Store(paymentId, resRec)
	requestId := rand.Uint64()
	rc.idsToMethods.Store(fmt.Sprintf("%d", requestId), serde.PayRequestMethod)

	pReq := serde.PaymentRequest{Amount: amount, Channel: id}
	message := serde.NewJsonRpcRequest(requestId, serde.PayRequestMethod, pReq)
	data, err := json.Marshal(message)
	if err != nil {
		panic("Could not marshal direct fund request")
	}
	rc.nts.SendMessage(serde.PayRequestMethod, data)

	<-resRec
}

func (rc *RpcClient) Close() {
	rc.nts.Close()
}

// registerHandlers registers error and response handles for the rpc client
func (rc *RpcClient) registerHandlers() {

	rc.nts.RegisterErrorHandler(func(id uint64, data []byte) {
		panic(fmt.Sprintf("Objective failed: %v", data))
	})

	rc.nts.RegisterResponseHandler(func(id uint64, data []byte) {
		rc.nts.Logger.Trace().Msgf("Rpc client received response: %+v", data)
		method, reqFound := rc.idsToMethods.Load(fmt.Sprintf("%d", id))
		if !reqFound {
			panic(fmt.Sprintf("Could not find request for response with id %d", id))
		}

		switch method {
		case serde.DirectFundRequestMethod:
			handleResponse[directfund.ObjectiveResponse](rc, data)
		case serde.DirectDefundRequestMethod:
			handleResponse[protocols.ObjectiveId](rc, data)
		case serde.VirtualFundRequestMethod:
			handleResponse[virtualfund.ObjectiveResponse](rc, data)
		case serde.VirtualDefundRequestMethod:
			handleResponse[protocols.ObjectiveId](rc, data)
		case serde.PayRequestMethod:
			handleResponse[serde.PaymentRequest](rc, data)

		}
	})
}

// handleResponse handles a response from the rpc server for the given client
// It is not a member of the RpcClient so it can take advantage of generics
func handleResponse[T serde.ResponsePayload](rc *RpcClient, data []byte) {
	rpcResponse := serde.JsonRpcResponse[T]{}
	err := json.Unmarshal(data, &rpcResponse)
	if err != nil {
		panic("could not unmarshal objective response")
	}

	if resRec, ok := rc.responses.Load(string(getObjectiveId(rpcResponse.Result))); ok {

		resRec <- rpcResponse.Result

		rc.idsToMethods.Delete(fmt.Sprintf("%d", rpcResponse.Id))
		rc.responses.Delete(fmt.Sprintf("%v", getObjectiveId(rpcResponse.Result)))
	}
}

// getObjectiveId returns the objective id from the result of a response
func getObjectiveId(result any) protocols.ObjectiveId {
	id, isId := result.(protocols.ObjectiveId)
	if isId {
		return id
	}
	res, isRes := result.(directfund.ObjectiveResponse)
	if isRes {
		return res.Id
	}
	vRes, isVRes := result.(virtualfund.ObjectiveResponse)
	if isVRes {
		return vRes.Id
	}
	pRes, isPRes := result.(serde.PaymentRequest)
	if isPRes {
		paymentId := fmt.Sprintf("PAYMENT-%s", pRes.Channel)
		return protocols.ObjectiveId(paymentId)
	}
	panic("Could not get id from result")
}
