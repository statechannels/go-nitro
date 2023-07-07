package rpc

import (
	"fmt"
	"os"

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
	rpcServer, node, messageService, err := RunRpcServer(pk, chainService, useDurableStore, msgPort, rpcPort, transportType, os.Stdout)
	if err != nil {
		return nil, nil, nil, err
	}

	fmt.Println("Nitro as a Service listening on port", rpcPort)
	return rpcServer, node, messageService, nil
}

func RunRpcServer(pk []byte, chainService chainservice.ChainService,
	useDurableStore bool, msgPort int, rpcPort int, transportType transport.TransportType, logDestination *os.File,
) (*rpc.RpcServer, *node.Node, *p2pms.P2PMessageService, error) {
	me := crypto.GetAddressFromSecretKeyBytes(pk)

	var ourStore store.Store
	var err error

	if useDurableStore {
		fmt.Println("Initialising durable store...")
		dataFolder := fmt.Sprintf("./data/nitro-service/%s", me.String())
		ourStore, err = store.NewDurableStore(pk, dataFolder, buntdb.Config{})
		if err != nil {
			return nil, nil, nil, err
		}

	} else {
		fmt.Println("Initialising mem store...")
		ourStore = store.NewMemStore(pk)
	}

	fmt.Println("Initializing message service on port " + fmt.Sprint(msgPort) + "...")
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

		fmt.Println("Initializing NATS RPC transport...")
		transport, err = nats.NewNatsTransportAsServer(rpcPort)
	case "ws":
		fmt.Println("Initializing websocket RPC transport...")
		transport, err = ws.NewWebSocketTransportAsServer(fmt.Sprint(rpcPort))
	default:
		err = fmt.Errorf("unknown transport type %s", transportType)
	}

	if err != nil {
		return nil, nil, nil, err
	}

	logger := zerolog.New(logDestination).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Str("node", ourStore.GetAddress().String()).
		Str("rpc", "server").
		Str("scope", "").
		Logger()

	rpcServer, err := rpc.NewRpcServer(&node, &logger, transport)
	if err != nil {
		return nil, nil, nil, err
	}
	return rpcServer, &node, messageService, nil
}
