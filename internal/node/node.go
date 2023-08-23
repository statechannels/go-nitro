package node

import (
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/statechannels/go-nitro/internal/chain"
	"github.com/statechannels/go-nitro/node"
	"github.com/statechannels/go-nitro/node/engine"
	"github.com/statechannels/go-nitro/node/engine/chainservice"
	"github.com/statechannels/go-nitro/node/engine/store"
	"github.com/tidwall/buntdb"

	p2pms "github.com/statechannels/go-nitro/node/engine/messageservice/p2p-message-service"
)

func InitializeNode(pkString string, chainOpts chain.ChainOpts,
	useDurableStore bool, durableStoreFolder string, msgPort int, logDestination *os.File, bootPeers []string,
) (*node.Node, *store.Store, *p2pms.P2PMessageService, *chainservice.EthChainService, error) {
	if pkString == "" {
		panic("pk must be set")
	}

	logger := zerolog.New(logDestination).
		With().
		Timestamp().
		Logger()

	pk := common.Hex2Bytes(pkString)
	ourStore, err := store.NewStore(pk, logger, useDurableStore, durableStoreFolder, buntdb.Config{})
	if err != nil {
		return nil, nil, nil, nil, err
	}

	logger.Info().Msg("Initializing message service on port " + fmt.Sprint(msgPort) + "...")
	messageService := p2pms.NewMessageService("127.0.0.1", msgPort, *ourStore.GetAddress(), pk, logDestination, bootPeers)

	logger.Info().Msg("Initializing chain service and connecting to " + chainOpts.ChainUrl + "...")
	ourChain, err := chain.InitializeEthChainService(chainOpts)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	node := node.New(
		messageService,
		ourChain,
		ourStore,
		logDestination,
		&engine.PermissivePolicy{},
	)

	return &node, &ourStore, messageService, ourChain, nil
}
