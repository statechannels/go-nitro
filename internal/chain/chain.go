package chain

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/node/engine/chainservice"
	NitroAdjudicator "github.com/statechannels/go-nitro/node/engine/chainservice/adjudicator"
	ConsensusApp "github.com/statechannels/go-nitro/node/engine/chainservice/consensusapp"
	chainutils "github.com/statechannels/go-nitro/node/engine/chainservice/utils"
	VirtualPaymentApp "github.com/statechannels/go-nitro/node/engine/chainservice/virtualpaymentapp"
	"github.com/statechannels/go-nitro/types"
)

const (
	FUNDED_TEST_PK  = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	ANVIL_CHAIN_URL = "ws://127.0.0.1:8545"
)

type ChainOpts struct {
	ChainUrl       string
	ChainAuthToken string
	ChainPk        string
	NaAddress      common.Address
	VpaAddress     common.Address
	CaAddress      common.Address
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

func StartAnvil() (*exec.Cmd, error) {
	chainCmd := exec.Command("anvil", "--chain-id", "1337")
	chainCmd.Stdout = os.Stdout
	chainCmd.Stderr = os.Stderr
	err := chainCmd.Start()
	if err == nil {
		// If Anvil start successfully, delay by 1 second for the chain to initialize
		time.Sleep(1 * time.Second)
		return chainCmd, nil
	}
	return chainCmd, err
}

// DeployContracts deploys the NitroAdjudicator, VirtualPaymentApp and ConsensusApp contracts.
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
