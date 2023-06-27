package node

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/internal/chain"
	interRpc "github.com/statechannels/go-nitro/internal/rpc"
	"github.com/statechannels/go-nitro/node"
	p2pms "github.com/statechannels/go-nitro/node/engine/messageservice/p2p-message-service"
	"github.com/statechannels/go-nitro/rpc"
	"github.com/statechannels/go-nitro/rpc/transport"
)

type NodeOpts struct {
	UseDurableStore bool
	MsgPort         int
	RpcPort         int
	Pk              string
	ChainPk         string
	ChainUrl        string
	ChainAuthToken  string
}

func RunNode(pkString string, chainOpts chain.ChainOpts,
	useDurableStore bool, useNats bool, msgPort int, rpcPort int,
) (*rpc.RpcServer, *node.Node, *p2pms.P2PMessageService, error) {
	if pkString == "" {
		panic("pk must be set")
	}
	pk := common.Hex2Bytes(pkString)

	chainService, err := chain.InitializeEthChainService(chainOpts)
	if err != nil {
		return nil, nil, nil, err
	}

	transportType := transport.Ws
	if useNats {
		transportType = transport.Nats
	}
	rpcServer, node, messageService, err := interRpc.InitializeRpcServer(pk, chainService, useDurableStore, msgPort, rpcPort, transportType)
	if err != nil {
		return nil, nil, nil, err
	}

	fmt.Println("Nitro as a Service listening on port", rpcPort)
	return rpcServer, node, messageService, nil
}
