package node

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/ethereum/go-ethereum/common"
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

	pk := common.Hex2Bytes(pkString)
	ourStore, err := store.NewStore(pk, useDurableStore, durableStoreFolder, buntdb.Config{})
	if err != nil {
		return nil, nil, nil, nil, err
	}

	slog.Info("Initializing message service on port " + fmt.Sprint(msgPort) + "...")
	messageService := p2pms.NewMessageService("127.0.0.1", msgPort, *ourStore.GetAddress(), pk, bootPeers)

	slog.Info("Initializing chain service and connecting to " + chainOpts.ChainUrl + "...")
	ourChain, err := chain.InitializeEthChainService(chainOpts)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	node := node.New(
		messageService,
		ourChain,
		ourStore,
		&engine.PermissivePolicy{},
	)

	return &node, &ourStore, messageService, ourChain, nil
}
