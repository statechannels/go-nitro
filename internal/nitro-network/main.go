package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/BurntSushi/toml"
	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/internal/chain"
	interRpc "github.com/statechannels/go-nitro/internal/rpc"
	"github.com/statechannels/go-nitro/internal/testdata"
	"github.com/statechannels/go-nitro/internal/utils"
	"github.com/statechannels/go-nitro/node"
	p2pms "github.com/statechannels/go-nitro/node/engine/messageservice/p2p-message-service"
	"github.com/statechannels/go-nitro/rpc"
	"github.com/statechannels/go-nitro/types"
)

type participantOpts struct {
	UseDurableStore bool
	MsgPort         int
	RpcPort         int
	Pk              string
	ChainPk         string
	ChainUrl        string
	ChainAuthToken  string
}

const (
	FUNDED_TEST_PK  = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	ANVIL_CHAIN_URL = "ws://127.0.0.1:8545"
)

func main() {
	err := InitializeNitroNetwork()
	if err != nil {
		panic(err)
	}
}

// InitializeNitroNetwork initializes a nitro network with 3 participants alice, irene, and bob
// A ledger channel is opened between alice and irene
// A ledger channel is opened between irene and bob
// A virtual channel is opened between alice and bob
func InitializeNitroNetwork() error {
	participants := []string{"alice", "irene", "bob"}
	servers := []*rpc.RpcServer{}
	nodes := []*node.Node{}
	msgServices := []*p2pms.P2PMessageService{}

	anvilCmd, err := chain.StartAnvil()
	if err != nil {
		return err
	}
	defer utils.StopCommands(anvilCmd)

	naAddress, vpaAddress, caAddress, err := chain.DeployContracts(context.Background(), ANVIL_CHAIN_URL, "", FUNDED_TEST_PK)
	if err != nil {
		return err
	}

	for _, participant := range participants {
		var nodeOpts participantOpts
		if _, err := toml.DecodeFile(fmt.Sprintf("./scripts/test-configs/%s.toml", participant), &nodeOpts); err != nil {
			return err
		}
		chainOpts := chain.ChainOpts{
			ChainUrl:       nodeOpts.ChainUrl,
			ChainAuthToken: nodeOpts.ChainAuthToken,
			ChainPk:        nodeOpts.ChainPk,
			NaAddress:      naAddress,
			CaAddress:      caAddress,
			VpaAddress:     vpaAddress,
		}
		server, node, msgService, err := interRpc.InitChainServiceAndRunRpcServer(nodeOpts.Pk, chainOpts, nodeOpts.UseDurableStore, false, nodeOpts.MsgPort, nodeOpts.RpcPort)
		if err != nil {
			return err
		}
		servers = append(servers, server)
		nodes = append(nodes, node)
		msgServices = append(msgServices, msgService)
	}
	utils.WaitForPeerInfoExchange(msgServices...)

	alice, irene, bob := nodes[0], nodes[1], nodes[2]
	err = createLedgerChannel(alice, irene)
	if err != nil {
		return err
	}

	err = createLedgerChannel(irene, bob)
	if err != nil {
		return err
	}

	outcome := testdata.Outcomes.Create(*alice.Address, *bob.Address, 1_000, 0, types.Address{})
	response, err := alice.CreatePaymentChannel([]common.Address{*irene.Address}, *bob.Address, 0, outcome)
	if err != nil {
		return err
	}
	<-alice.ObjectiveCompleteChan(response.Id)
	<-bob.ObjectiveCompleteChan(response.Id)

	fmt.Printf("Created payment channel between Alice and Bob")

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

func createLedgerChannel(left *node.Node, right *node.Node) error {
	ledgerChannelDeposit := uint(5_000_000)
	asset := types.Address{}
	outcome := testdata.Outcomes.Create(*left.Address, *right.Address, ledgerChannelDeposit, ledgerChannelDeposit, asset)
	response, err := left.CreateLedgerChannel(*right.Address, 0, outcome)
	if err != nil {
		return err
	}

	<-left.ObjectiveCompleteChan(response.Id)
	<-right.ObjectiveCompleteChan(response.Id)

	fmt.Printf("Created ledged channel between %s and %s", left.Address, right.Address)
	return nil
}
