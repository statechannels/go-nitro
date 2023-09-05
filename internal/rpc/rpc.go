package rpc

import (
	"fmt"
	"log/slog"

	"github.com/statechannels/go-nitro/node"
	"github.com/statechannels/go-nitro/rpc"
	"github.com/statechannels/go-nitro/rpc/transport"
	"github.com/statechannels/go-nitro/rpc/transport/http"
	"github.com/statechannels/go-nitro/rpc/transport/nats"
)

func InitializeRpcServer(node *node.Node, rpcPort int, useNats bool) (*rpc.RpcServer, error) {
	var transport transport.Responder
	var err error

	if useNats {
		slog.Info("Initializing NATS RPC transport...")
		transport, err = nats.NewNatsTransportAsServer(rpcPort)
	} else {
		slog.Info("Initializing Http RPC transport...")
		transport, err = http.NewHttpTransportAsServer(fmt.Sprint(rpcPort))
	}
	if err != nil {
		return nil, err
	}

	rpcServer, err := rpc.NewRpcServer(node, transport)
	if err != nil {
		return nil, err
	}

	return rpcServer, nil
}
