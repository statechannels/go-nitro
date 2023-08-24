package rpc

import (
	"fmt"
	"log/slog"
	"path/filepath"

	"github.com/ethereum/go-ethereum/common"

	"github.com/statechannels/go-nitro/crypto"
	"github.com/statechannels/go-nitro/internal/chain"
	"github.com/statechannels/go-nitro/internal/logging"
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
	bootPeers []string,
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

	rpcServer, node, messageService, err := RunRpcServer(pk, chainService, useDurableStore, durableStoreFolder, msgPort, rpcPort, transportType, bootPeers)
	if err != nil {
		return nil, nil, nil, err
	}

	slog.Info("Nitro as a Service listening", "rpc-port", rpcPort, "msg-port", msgPort, "transport", transportType, "address", node.Address.String())

	return rpcServer, node, messageService, nil
}

func RunRpcServer(pk []byte, chainService chainservice.ChainService,
	useDurableStore bool, durableStoreFolder string, msgPort int, rpcPort int, transportType transport.TransportType,
	bootPeers []string,
) (*rpc.RpcServer, *node.Node, *p2pms.P2PMessageService, error) {
	me := crypto.GetAddressFromSecretKeyBytes(pk)
	logger := logging.LoggerWithAddress(slog.Default(), me)
	var ourStore store.Store
	var err error

	if useDurableStore {
		dataFolder := filepath.Join(durableStoreFolder, me.String())
		logger.Info("Initialising durable store", "dataFolder", dataFolder)

		ourStore, err = store.NewDurableStore(pk, dataFolder, buntdb.Config{})
		if err != nil {
			return nil, nil, nil, err
		}

	} else {
		logger.Info("Initialising mem store...")
		ourStore = store.NewMemStore(pk)
	}

	logger.Info("Initializing message service ", "port", msgPort)

	messageService := p2pms.NewMessageService("127.0.0.1", msgPort, *ourStore.GetAddress(), pk, bootPeers)
	node := node.New(
		messageService,
		chainService,
		ourStore,
		&engine.PermissivePolicy{})

	var transport transport.Responder

	switch transportType {
	case "nats":

		logger.Info("Initializing NATS RPC transport...")
		transport, err = nats.NewNatsTransportAsServer(rpcPort)
	case "ws":
		logger.Info("Initializing websocket RPC transport...")

		transport, err = ws.NewWebSocketTransportAsServer(fmt.Sprint(rpcPort))
	default:
		err = fmt.Errorf("unknown transport type %s", transportType)
	}

	if err != nil {
		return nil, nil, nil, err
	}

	rpcServer, err := rpc.NewRpcServer(&node, transport)
	if err != nil {
		return nil, nil, nil, err
	}
	return rpcServer, &node, messageService, nil
}
