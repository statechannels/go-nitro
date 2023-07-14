package rpc

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/statechannels/go-nitro/crypto"
	"github.com/statechannels/go-nitro/internal/chain"
	"github.com/statechannels/go-nitro/node"
	"github.com/statechannels/go-nitro/node/engine"
	"github.com/statechannels/go-nitro/node/engine/chainservice"
	p2pms "github.com/statechannels/go-nitro/node/engine/messageservice/p2p-message-service"
	"github.com/statechannels/go-nitro/node/engine/store"
	"github.com/statechannels/go-nitro/rpc"
	"github.com/statechannels/go-nitro/rpc/transport"
	"github.com/statechannels/go-nitro/rpc/transport/nats"
	"github.com/statechannels/go-nitro/rpc/transport/ws"
	"github.com/tidwall/buntdb"
)

func InitChainServiceAndRunRpcServer(pkString string, chainOpts chain.ChainOpts,
	useDurableStore bool, durableStoreFolder string, useNats bool, msgPort int, rpcPort int,
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
	rpcServer, node, messageService, err := RunRpcServer(pk, chainService, useDurableStore, durableStoreFolder, msgPort, rpcPort, transportType, os.Stdout)
	if err != nil {
		return nil, nil, nil, err
	}

	fmt.Println("Nitro as a Service listening on port", rpcPort)
	return rpcServer, node, messageService, nil
}

func RunRpcServer(pk []byte, chainService chainservice.ChainService,
	useDurableStore bool,
	durableStoreFolder string,
	msgPort int, rpcPort int,
	transportType transport.TransportType,
	logDestination *os.File,
) (*rpc.RpcServer, *node.Node, *p2pms.P2PMessageService, error) {
	me := crypto.GetAddressFromSecretKeyBytes(pk)

	var ourStore store.Store
	var err error

	logger := zerolog.New(logDestination).
		With().
		Timestamp().
		Logger()

	if useDurableStore {
		dataFolder := filepath.Join(durableStoreFolder, me.String())
		logger.Info().Msgf("Initialising durable store in %s...", dataFolder)

		ourStore, err = store.NewDurableStore(pk, dataFolder, buntdb.Config{})
		if err != nil {
			return nil, nil, nil, err
		}

	} else {
		logger.Info().Msg("Initialising mem store...")
		ourStore = store.NewMemStore(pk)
	}

	logger.Info().Msg("Initializing message service on port " + fmt.Sprint(msgPort) + "...")
	messageService := p2pms.NewMessageService("127.0.0.1", msgPort, *ourStore.GetAddress(), pk, true, logDestination)
	node := node.New(
		messageService,
		chainService,
		ourStore,
		logDestination,
		&engine.PermissivePolicy{},
		nil)

	var transport transport.Responder

	switch transportType {
	case "nats":

		logger.Info().Msg("Initializing NATS RPC transport...")
		transport, err = nats.NewNatsTransportAsServer(rpcPort)
	case "ws":
		logger.Info().Msg("Initializing websocket RPC transport...")
		mux := http.DefaultServeMux

		transport, err = ws.NewWebSocketTransportAsServer(fmt.Sprint(rpcPort), mux)
	default:
		err = fmt.Errorf("unknown transport type %s", transportType)
	}

	if err != nil {
		return nil, nil, nil, err
	}

	serverLogger := zerolog.New(logDestination).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Str("node", ourStore.GetAddress().String()).
		Str("rpc", "server").
		Logger()

	rpcServer, err := rpc.NewRpcServer(&node, &serverLogger, transport)
	if err != nil {
		return nil, nil, nil, err
	}
	return rpcServer, &node, messageService, nil
}
