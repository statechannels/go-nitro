package rpc

import (
	"fmt"
	"math/big"
	"math/rand"

	"github.com/ethereum/go-ethereum/common"
	"github.com/nats-io/nats.go"
	"github.com/statechannels/go-nitro/client/engine/store/safesync"
	"github.com/statechannels/go-nitro/network"
	netproto "github.com/statechannels/go-nitro/network/protocol"
	"github.com/statechannels/go-nitro/network/protocol/parser"
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
func NewRpcClient(rpcServerUrl string, myAddress types.Address, chainId *big.Int) *RpcClient {

	nc, err := nats.Connect(rpcServerUrl)
	handleError(err)
	trp := natstrans.NewNatsTransport(nc, []string{fmt.Sprintf("nitro.%s", network.DirectFundRequestMethod)})

	con, err := trp.PollConnection()
	handleError(err)
	nts := network.NewNetworkService(con, &serde.JsonRpc{})

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
	rc.nts.SendMessage(netproto.NewMessage(netproto.TypeRequest, rand.Uint64(), network.DirectFundRequestMethod, []any{&objReq}))

	objRes := <-resRec
	return objRes

}

func (rc *RpcClient) Close() {
	rc.nts.Close()
}

// registerHandlers registers error and response handles for the rpc client
func (rs *RpcClient) registerHandlers() {

	rs.nts.RegisterErrorHandler(network.DirectFundRequestMethod, func(m *netproto.Message) {
		panic(fmt.Sprintf("Objective failed: %v", *m))
	})
	rs.nts.RegisterResponseHandler(network.DirectFundRequestMethod, func(m *netproto.Message) {
		rs.nts.Logger.Trace().Msgf("Rpc client received response: %+v", m)
		if len(m.Args) < 1 {
			panic("unexpected empty args for direct funding method")

		}

		for i := 0; i < len(m.Args); i++ {

			raw := m.Args[i].(map[string]interface{})
			res := parser.ParseDirectFundResponse(raw)

			// Once we receive the response we notify the appropriate channel
			if resRec, ok := rs.responses.Load(string(res.Id)); ok {
				rs.responses.Delete(string(res.Id))
				resRec <- res
			}

		}
	})
}
