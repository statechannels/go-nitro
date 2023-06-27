package infra

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/statechannels/go-nitro/crypto"
	"github.com/statechannels/go-nitro/node"
	"github.com/statechannels/go-nitro/node/engine"
	"github.com/statechannels/go-nitro/node/engine/chainservice"
	NitroAdjudicator "github.com/statechannels/go-nitro/node/engine/chainservice/adjudicator"
	ConsensusApp "github.com/statechannels/go-nitro/node/engine/chainservice/consensusapp"
	chainutils "github.com/statechannels/go-nitro/node/engine/chainservice/utils"
	VirtualPaymentApp "github.com/statechannels/go-nitro/node/engine/chainservice/virtualpaymentapp"
	p2pms "github.com/statechannels/go-nitro/node/engine/messageservice/p2p-message-service"
	"github.com/statechannels/go-nitro/node/engine/store"
	"github.com/statechannels/go-nitro/rpc"
	"github.com/statechannels/go-nitro/rpc/transport"
	"github.com/statechannels/go-nitro/rpc/transport/nats"
	"github.com/statechannels/go-nitro/rpc/transport/ws"
	"github.com/statechannels/go-nitro/types"
	"github.com/tidwall/buntdb"
)

const (
	FUNDED_TEST_PK  = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	ANVIL_CHAIN_URL = "ws://127.0.0.1:8545"
)

// Start of Nitro Utilities

func InitializeRpcServer(pk []byte, chainService chainservice.ChainService,
	useDurableStore bool, msgPort int, rpcPort int, transportType transport.TransportType,
) (*rpc.RpcServer, *p2pms.P2PMessageService, error) {
	me := crypto.GetAddressFromSecretKeyBytes(pk)

	logDestination := os.Stdout

	var ourStore store.Store
	var err error

	if useDurableStore {
		fmt.Println("Initialising durable store...")
		dataFolder := fmt.Sprintf("./data/nitro-service/%s", me.String())
		ourStore, err = store.NewDurableStore(pk, dataFolder, buntdb.Config{})
		if err != nil {
			return nil, nil, err
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
		return nil, nil, err
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
		return nil, nil, err
	}
	return rpcServer, messageService, nil
}

func InitializeEthChainService(chainOpts ChainOpts) (*chainservice.EthChainService, error) {
	if chainOpts.ChainPk == "" {
		return nil, fmt.Errorf("chainpk must be set")
	}

	fmt.Println("Initializing chain service and connecting to " + chainOpts.ChainUrl + "...")
	return chainservice.NewEthChainService(
		chainOpts.ChainUrl,
		chainOpts.ChainAuthToken,
		chainOpts.ChainPk,
		chainOpts.NaAddress,
		chainOpts.CaAddress,
		chainOpts.VpaAddress,
		os.Stdout)
}

type ChainOpts struct {
	ChainUrl       string
	ChainAuthToken string
	ChainPk        string
	NaAddress      common.Address
	VpaAddress     common.Address
	CaAddress      common.Address
}

func RunNode(pkString string, chainOpts ChainOpts,
	useDurableStore bool, useNats bool, msgPort int, rpcPort int,
) (*rpc.RpcServer, *p2pms.P2PMessageService, error) {
	if pkString == "" {
		panic("pk must be set")
	}
	pk := common.Hex2Bytes(pkString)

	chainService, err := InitializeEthChainService(chainOpts)
	if err != nil {
		return nil, nil, err
	}

	transportType := transport.Ws
	if useNats {
		transportType = transport.Nats
	}
	rpcServer, messageService, err := InitializeRpcServer(pk, chainService, useDurableStore, msgPort, rpcPort, transportType)
	if err != nil {
		return nil, nil, err
	}

	fmt.Println("Nitro as a Service listening on port", rpcPort)
	return rpcServer, messageService, nil
}

// waitForPeerInfoExchange waits for all the P2PMessageServices to receive peer info from each other
func WaitForPeerInfoExchange(services ...*p2pms.P2PMessageService) {
	for _, s := range services {
		for i := 0; i < len(services)-1; i++ {
			<-s.PeerInfoReceived()
		}
	}
}

// End of Nitro Utilities

// Start of Ethereum chain utilities

func StartAnvil() (*exec.Cmd, error) {
	chainCmd := exec.Command("anvil", "--chain-id", "1337")
	chainCmd.Stdout = os.Stdout
	chainCmd.Stderr = os.Stderr
	err := chainCmd.Start()
	if err == nil {
		// Give the chain a second to start up
		time.Sleep(1 * time.Second)
	}
	return chainCmd, err
}

// DeployContracts deploys the  NitroAdjudicator contract.
func DeployContracts(ctx context.Context) (na common.Address, vpa common.Address, ca common.Address, err error) {
	ethClient, txSubmitter, err := chainutils.ConnectToChain(context.Background(), "ws://127.0.0.1:8545", "", common.Hex2Bytes(FUNDED_TEST_PK))
	if err != nil {
		return types.Address{}, types.Address{}, types.Address{}, err
	}
	na, _, _, err = NitroAdjudicator.DeployNitroAdjudicator(txSubmitter, ethClient)
	if err != nil {
		return types.Address{}, types.Address{}, types.Address{}, err
	}
	fmt.Printf("Deployed NitroAdjudicator at %s\n", na.String())
	vpa, _, _, err = VirtualPaymentApp.DeployVirtualPaymentApp(txSubmitter, ethClient)
	if err != nil {
		return types.Address{}, types.Address{}, types.Address{}, err
	}
	fmt.Printf("Deployed VirtualPaymentApp at %s\n", vpa.String())
	ca, _, _, err = ConsensusApp.DeployConsensusApp(txSubmitter, ethClient)
	if err != nil {
		return types.Address{}, types.Address{}, types.Address{}, err
	}
	fmt.Printf("Deployed ConsensusApp at %s\n", ca.String())
	return
}

// End of Ethereum chain utilities

// Start of general of utilities

// StopCommands stops the given executing commands
func StopCommands(cmds ...*exec.Cmd) {
	for _, cmd := range cmds {
		fmt.Printf("Stopping process %v\n", cmd.Args)
		err := cmd.Process.Signal(syscall.SIGINT)
		if err != nil {
			panic(err)
		}
		err = cmd.Process.Kill()
		if err != nil {
			panic(err)
		}
	}
}

// End of general of utilities

type NodeOpts struct {
	UseDurableStore bool
	MsgPort         int
	RpcPort         int
	Pk              string
	ChainPk         string
	ChainUrl        string
	ChainAuthToken  string
}

func InitializeNitroNetwork() error {
	participants := []string{"alice", "bob", "irene"}
	servers := []*rpc.RpcServer{}
	msgServices := []*p2pms.P2PMessageService{}

	anvilCmd, err := StartAnvil()
	defer StopCommands(anvilCmd)
	if err != nil {
		return err
	}

	naAddress, vpaAddress, caAddress, err := DeployContracts(context.Background())
	if err != nil {
		return err
	}

	for _, participant := range participants {
		var nodeOpts NodeOpts
		if _, err := toml.DecodeFile(fmt.Sprintf("../scripts/test-configs/%s.toml", participant), &nodeOpts); err != nil {
			return err
		}
		chainOpts := ChainOpts{
			ChainUrl:       nodeOpts.ChainUrl,
			ChainAuthToken: nodeOpts.ChainAuthToken,
			ChainPk:        nodeOpts.ChainPk,
			NaAddress:      naAddress,
			CaAddress:      caAddress,
			VpaAddress:     vpaAddress,
		}
		server, msgService, err := RunNode(nodeOpts.Pk, chainOpts, nodeOpts.UseDurableStore, false, nodeOpts.MsgPort, nodeOpts.RpcPort)
		if err != nil {
			return err
		}
		servers = append(servers, server)
		msgServices = append(msgServices, msgService)
	}
	WaitForPeerInfoExchange(msgServices...)

	stopChan := make(chan os.Signal, 2)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-stopChan // wait for interrupt or terminate signal

	for _, server := range servers {
		err := server.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
