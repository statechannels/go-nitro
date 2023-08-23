package rpc

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/statechannels/go-nitro/node"
	"github.com/statechannels/go-nitro/rpc"
	"github.com/statechannels/go-nitro/rpc/transport"
	"github.com/statechannels/go-nitro/rpc/transport/nats"
	"github.com/statechannels/go-nitro/rpc/transport/ws"
)

func InitializeRpcServer(node *node.Node, rpcPort int, useNats bool, logDestination *os.File) (*rpc.RpcServer, error) {
	logger := zerolog.New(logDestination).
		With().
		Timestamp().
		Logger()

	var transport transport.Responder
	var err error

	if useNats {
		logger.Info().Msg("Initializing NATS RPC transport...")
		transport, err = nats.NewNatsTransportAsServer(rpcPort)
	} else {
		logger.Info().Msg("Initializing websocket RPC transport...")
		transport, err = ws.NewWebSocketTransportAsServer(fmt.Sprint(rpcPort), logger)
	}
	if err != nil {
		return nil, err
	}

	rpcServer, err := rpc.NewRpcServer(node, &logger, transport)
	if err != nil {
		return nil, err
	}
	return rpcServer, nil
}
