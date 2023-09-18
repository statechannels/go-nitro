package node

import (
	"fmt"
	"log/slog"

	"github.com/statechannels/go-nitro/node"
	"github.com/statechannels/go-nitro/node/engine"
	"github.com/statechannels/go-nitro/node/engine/chainservice"
	"github.com/statechannels/go-nitro/node/engine/store"

	p2pms "github.com/statechannels/go-nitro/node/engine/messageservice/p2p-message-service"
)

func InitializeNode(chainOpts chainservice.ChainOpts, storeOpts store.StoreOpts, messageOpts p2pms.MessageOpts) (*node.Node, *store.Store, *p2pms.P2PMessageService, chainservice.ChainService, error) {
	ourStore, err := store.NewStore(storeOpts)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	slog.Info("Initializing message service on port " + fmt.Sprint(messageOpts.Port) + "...")
	messageOpts.SCAddr = *ourStore.GetAddress()
	messageService := p2pms.NewMessageService(messageOpts)

	// Compare chainOpts.ChainStartBlock to lastBlockNum seen in store. The larger of the two
	// gets passed as an argument when creating NewEthChainService
	storeBlockNum, err := ourStore.GetLastBlockNumSeen()
	if err != nil {
		return nil, nil, nil, nil, err
	}
	if storeBlockNum > chainOpts.ChainStartBlock {
		chainOpts.ChainStartBlock = storeBlockNum
	}

	slog.Info("Initializing chain service...")
	ourChain, err := chainservice.NewEthChainService(chainOpts)
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
