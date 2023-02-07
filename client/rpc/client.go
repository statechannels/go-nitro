package rpc

import (
	"math/big"
	"math/rand"

	"github.com/ethereum/go-ethereum/common"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	"github.com/statechannels/go-nitro/network"
	"github.com/statechannels/go-nitro/network/serde"
	"github.com/statechannels/go-nitro/network/transport"
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
	connection transport.Connection
	myAddress  types.Address
	chainId    *big.Int
	logger     zerolog.Logger
}

// NewRpcClient creates a new RpcClient
func NewRpcClient(rpcServerUrl string, myAddress types.Address, chainId *big.Int, logger zerolog.Logger) *RpcClient {
	nc, err := nats.Connect(rpcServerUrl)
	handleError(err)
	con := natstrans.NewNatsConnection(nc)

	c := &RpcClient{con, myAddress, chainId, logger}
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

	resChan, err := network.Request[virtualfund.ObjectiveRequest, virtualfund.ObjectiveResponse](rc.connection, objReq, rc.logger)
	if err != nil {
		panic(err)
	}

	objRes := <-resChan
	if objRes.Error != nil {
		panic(err)
	}

	return objRes.Payload
}

// CloseVirtual closes a virtual channel
func (rc *RpcClient) CloseVirtual(id types.Destination) protocols.ObjectiveId {
	objReq := virtualdefund.NewObjectiveRequest(
		id)

	resChan, err := network.Request[virtualdefund.ObjectiveRequest, protocols.ObjectiveId](rc.connection, objReq, rc.logger)
	if err != nil {
		panic(err)
	}

	objRes := <-resChan
	if objRes.Error != nil {
		panic(err)
	}

	return objRes.Payload
}

// CreateLedger creates a new ledger channel
func (rc *RpcClient) CreateLedger(counterparty types.Address, ChallengeDuration uint32, outcome outcome.Exit) directfund.ObjectiveResponse {

	objReq := directfund.NewObjectiveRequest(
		counterparty,
		100,
		outcome,
		uint64(rand.Float64()), // TODO: Since numeric fields get converted to a float64 in transit we need to prevent overflow
		common.Address{})

	resChan, err := network.Request[directfund.ObjectiveRequest, directfund.ObjectiveResponse](rc.connection, objReq, rc.logger)
	if err != nil {
		panic(err)
	}

	objRes := <-resChan
	if objRes.Error != nil {
		panic(err)
	}

	return objRes.Payload
}

// CloseLedger closes a ledger channel
func (rc *RpcClient) CloseLedger(id types.Destination) protocols.ObjectiveId {
	objReq := directdefund.NewObjectiveRequest(id)

	resChan, err := network.Request[directdefund.ObjectiveRequest, protocols.ObjectiveId](rc.connection, objReq, rc.logger)
	if err != nil {
		panic(err)
	}

	objRes := <-resChan
	if objRes.Error != nil {
		panic(err)
	}

	return objRes.Payload
}

// Pay uses the specified channel to pay the specified amount
func (rc *RpcClient) Pay(id types.Destination, amount uint64) {
	pReq := serde.PaymentRequest{Amount: amount, Channel: id}

	resChan, err := network.Request[serde.PaymentRequest, serde.PaymentRequest](rc.connection, pReq, rc.logger)
	if err != nil {
		panic(err)
	}

	res := <-resChan
	if res.Error != nil {
		panic(res.Error)
	}
}

func (rc *RpcClient) Close() {
	rc.connection.Close()
}
