package rpc

import (
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
	nts              *network.NetworkService
	clientConnection *network.ClientConnection
	myAddress        types.Address
	chainId          *big.Int

	// responses is a collection of channels that are used to wait until a response is received from the RPC server
	responses safesync.Map[chan interface{}]

	idsToMethods safesync.Map[serde.RequestMethod]
}

// NewRpcClient creates a new RpcClient
func NewRpcClient(rpcServerUrl string, myAddress types.Address, chainId *big.Int, logger zerolog.Logger) *RpcClient {

	nc, err := nats.Connect(rpcServerUrl)
	handleError(err)
	trp := natstrans.NewNatsTransport(nc)

	con, err := trp.PollConnection()
	handleError(err)
	nts := network.NewNetworkService(con)
	clientConnection := network.ClientConnection{Connection: con}
	nts.Logger = logger
	c := &RpcClient{nts, &clientConnection, myAddress, chainId, safesync.Map[chan interface{}]{}, safesync.Map[serde.RequestMethod]{}}
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

	resChan, err := network.Request(rc.clientConnection, objReq, rc.nts.Logger)
	if err != nil {
		panic(err)
	}

	objRes := <-resChan
	if objRes.Error != nil {
		panic(err)
	}

	return objRes.Data.(serde.JsonRpcResponse[virtualfund.ObjectiveResponse]).Result
}

// CloseVirtual closes a virtual channel
func (rc *RpcClient) CloseVirtual(id types.Destination) protocols.ObjectiveId {
	objReq := virtualdefund.NewObjectiveRequest(
		id)

	resChan, err := network.Request(rc.clientConnection, objReq, rc.nts.Logger)
	if err != nil {
		panic(err)
	}

	objRes := <-resChan
	if objRes.Error != nil {
		panic(err)
	}

	return objRes.Data.(serde.JsonRpcResponse[protocols.ObjectiveId]).Result
}

// CreateLedger creates a new ledger channel
func (rc *RpcClient) CreateLedger(counterparty types.Address, ChallengeDuration uint32, outcome outcome.Exit) directfund.ObjectiveResponse {

	objReq := directfund.NewObjectiveRequest(
		counterparty,
		100,
		outcome,
		uint64(rand.Float64()), // TODO: Since numeric fields get converted to a float64 in transit we need to prevent overflow
		common.Address{})

	resChan, err := network.Request(rc.clientConnection, objReq, rc.nts.Logger)
	if err != nil {
		panic(err)
	}

	objRes := <-resChan
	if objRes.Error != nil {
		panic(err)
	}

	return objRes.Data.(serde.JsonRpcResponse[directfund.ObjectiveResponse]).Result
}

// CloseLedger closes a ledger channel
func (rc *RpcClient) CloseLedger(id types.Destination) protocols.ObjectiveId {
	objReq := directdefund.NewObjectiveRequest(id)

	resChan, err := network.Request(rc.clientConnection, objReq, rc.nts.Logger)
	if err != nil {
		panic(err)
	}

	objRes := <-resChan
	if objRes.Error != nil {
		panic(err)
	}

	return objRes.Data.(serde.JsonRpcResponse[protocols.ObjectiveId]).Result
}

// Pay uses the specified channel to pay the specified amount
func (rc *RpcClient) Pay(id types.Destination, amount uint64) {
	pReq := serde.PaymentRequest{Amount: amount, Channel: id}

	resChan, err := network.Request(rc.clientConnection, pReq, rc.nts.Logger)
	if err != nil {
		panic(err)
	}

	res := <-resChan
	if res.Error != nil {
		panic(res.Error)
	}
}

func (rc *RpcClient) Close() {
	rc.nts.Close()
}

// registerHandlers registers error and response handles for the rpc client
func (rc *RpcClient) registerHandlers() {

	rc.nts.RegisterErrorHandler(func(id uint64, data []byte) {
		panic(fmt.Sprintf("Objective failed: %v", data))
	})
}
