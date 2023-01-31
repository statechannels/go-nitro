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
	"github.com/statechannels/go-nitro/protocols/directfund"
	"github.com/statechannels/go-nitro/types"

	"github.com/statechannels/go-nitro/channel/state/outcome"
)

// RpcClient is a client for making nitro rpc calls
type RpcClient struct {
	nts       *network.NetworkService
	myAddress types.Address
	chainId   *big.Int

	// responses is a collection of channels that are used to wait until a response is received from the RPC server
	responses safesync.Map[chan directfund.ObjectiveResponse]
}

// NewRpcClient creates a new RpcClient
func NewRpcClient(rpcServerUrl string, myAddress types.Address, chainId *big.Int, logger zerolog.Logger) *RpcClient {

	nc, err := nats.Connect(rpcServerUrl)
	handleError(err)
	trp := natstrans.NewNatsTransport(nc, []string{fmt.Sprintf("nitro.%s", network.DirectFundRequestMethod)})

	con, err := trp.PollConnection()
	handleError(err)
	nts := network.NewNetworkService(con)
	nts.Logger = logger
	c := &RpcClient{nts, myAddress, chainId, safesync.Map[chan directfund.ObjectiveResponse]{}}
	c.registerHandlers()
	return c
}

// CreateLedger creates a new ledger channel
func (rc *RpcClient) CreateLedger(counterparty types.Address, ChallengeDuration uint32, outcome outcome.Exit) directfund.ObjectiveResponse {

	objReq := directfund.NewObjectiveRequest(
		counterparty,
		100,
		outcome,
		uint64(rand.Float64()), // TODO: Since numeric fields get converted to a float64 in transit we need to prevent overflow
		common.Address{})

	resRec := make(chan directfund.ObjectiveResponse)

	rc.responses.Store(string(objReq.Id(rc.myAddress, rc.chainId)), resRec)
	message := serde.NewDirectFundRequestMessage(rand.Uint64(), objReq)
	data, err := json.Marshal(message)
	if err != nil {
		panic("Could not marshal direct fund request")
	}
	rc.nts.SendMessage(network.DirectFundRequestMethod, data)

	objRes := <-resRec
	return objRes

}

func (rc *RpcClient) Close() {
	rc.nts.Close()
}

// registerHandlers registers error and response handles for the rpc client
func (rs *RpcClient) registerHandlers() {

	rs.nts.RegisterErrorHandler(network.DirectFundRequestMethod, func(data []byte) {
		panic(fmt.Sprintf("Objective failed: %v", data))
	})
	rs.nts.RegisterResponseHandler(func(data []byte) {
		rs.nts.Logger.Trace().Msgf("Rpc client received response: %+v", data)

		rpcResponse := serde.JsonRpcDirectFundResponse{}
		err := json.Unmarshal(data, &rpcResponse)
		if err != nil {
			panic("could not unmarshal direct fund objective response")
		}

		if resRec, ok := rs.responses.Load(fmt.Sprintf("%v", rpcResponse.ObjectiveResponse.Id)); ok {
			rs.responses.Delete(fmt.Sprintf("%v", rpcResponse.Id))
			resRec <- rpcResponse.ObjectiveResponse
		}
	})
}
