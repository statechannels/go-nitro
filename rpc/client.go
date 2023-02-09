package rpc

import (
	"encoding/json"
	"math/big"
	"math/rand"

	"github.com/ethereum/go-ethereum/common"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/protocols/directdefund"
	"github.com/statechannels/go-nitro/protocols/directfund"
	"github.com/statechannels/go-nitro/protocols/virtualdefund"
	"github.com/statechannels/go-nitro/protocols/virtualfund"
	"github.com/statechannels/go-nitro/rpc/network"
	"github.com/statechannels/go-nitro/rpc/serde"
	"github.com/statechannels/go-nitro/rpc/transport"

	natstrans "github.com/statechannels/go-nitro/rpc/transport/nats"
	"github.com/statechannels/go-nitro/types"

	"github.com/statechannels/go-nitro/channel/state/outcome"
)

// RpcClient is a client for making nitro rpc calls
type RpcClient struct {
	connection          transport.Connection
	myAddress           types.Address
	chainId             *big.Int
	logger              zerolog.Logger
	completedObjectives chan protocols.ObjectiveId
}

// NewRpcClient creates a new RpcClient
func NewRpcClient(rpcServerUrl string, myAddress types.Address, chainId *big.Int, logger zerolog.Logger) (*RpcClient, error) {
	nc, err := nats.Connect(rpcServerUrl)
	if err != nil {
		return nil, err
	}
	con := natstrans.NewNatsConnection(nc)
	c := &RpcClient{con, myAddress, chainId, logger, make(chan protocols.ObjectiveId, 100)}
	err = c.subscribeToNotifications()
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

func (rc *RpcClient) CompletedObjectives() <-chan protocols.ObjectiveId {
	return rc.completedObjectives
}

func (rc *RpcClient) Close() {
	rc.connection.Close()
}

func (rc *RpcClient) subscribeToNotifications() error {
	err := rc.connection.Subscribe(serde.ObjectiveCompleted, func(data []byte) {
		rpcRequest := serde.JsonRpcRequest[protocols.ObjectiveId]{}
		err := json.Unmarshal(data, &rpcRequest)
		if err != nil {
			panic(err)
		}
		rc.completedObjectives <- rpcRequest.Params
	})
	return err
}

func waitForRequest[T serde.RequestPayload, U serde.ResponsePayload](rc *RpcClient, requestData T) U {
	resChan, err := network.Request[T, U](rc.connection, requestData, rc.logger)
	if err != nil {
		panic(err)
	}

	res := <-resChan
	if res.Error != nil {
		panic(res.Error)
	}

	return res.Payload
}
