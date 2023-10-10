package chain

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"math/big"
	"os"
	"os/exec"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	NitroAdjudicator "github.com/statechannels/go-nitro/node/engine/chainservice/adjudicator"
	ConsensusApp "github.com/statechannels/go-nitro/node/engine/chainservice/consensusapp"
	chainutils "github.com/statechannels/go-nitro/node/engine/chainservice/utils"
	VirtualPaymentApp "github.com/statechannels/go-nitro/node/engine/chainservice/virtualpaymentapp"
	"github.com/statechannels/go-nitro/types"
)

func StartAnvil() (*exec.Cmd, error) {
	chainCmd := exec.Command("anvil", "--chain-id", "1337", "--block-time", "1", "--silent")
	chainCmd.Stdout = os.Stdout
	chainCmd.Stderr = os.Stderr
	err := chainCmd.Start()
	if err != nil {
		return &exec.Cmd{}, nil
	}
	// If Anvil start successfully, delay by 1 second for the chain to initialize
	time.Sleep(1 * time.Second)
	return chainCmd, nil
}

// DeployContracts deploys the NitroAdjudicator, VirtualPaymentApp and ConsensusApp contracts.
func DeployContracts(ctx context.Context, chainUrl, chainAuthToken, chainPk string) (na common.Address, vpa common.Address, ca common.Address, err error) {
	ethClient, txSubmitter, err := chainutils.ConnectToChain(context.Background(), chainUrl, chainAuthToken, common.Hex2Bytes(chainPk))
	if err != nil {
		return types.Address{}, types.Address{}, types.Address{}, err
	}

	na, err = deployContract(ctx, "NitroAdjudicator", ethClient, txSubmitter, NitroAdjudicator.DeployNitroAdjudicator)
	if err != nil {
		return types.Address{}, types.Address{}, types.Address{}, err
	}

	vpa, err = deployContract(ctx, "VirtualPaymentApp", ethClient, txSubmitter, VirtualPaymentApp.DeployVirtualPaymentApp)
	if err != nil {
		return types.Address{}, types.Address{}, types.Address{}, err
	}

	ca, err = deployContract(ctx, "ConsensusApp", ethClient, txSubmitter, ConsensusApp.DeployConsensusApp)
	if err != nil {
		return types.Address{}, types.Address{}, types.Address{}, err
	}
	transferEth(ethClient, txSubmitter, na, big.NewInt(1000000000000000000))
	return
}

type contractBackend interface {
	NitroAdjudicator.NitroAdjudicator | VirtualPaymentApp.VirtualPaymentApp | ConsensusApp.ConsensusApp
}

// deployFunc is a function that deploys a contract and returns the contract address, backend, and transaction.
type deployFunc[T contractBackend] func(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *ethTypes.Transaction, *T, error)

// deployContract deploys a contract and waits for the transaction to be mined.
func deployContract[T contractBackend](ctx context.Context, name string, ethClient *ethclient.Client, txSubmitter *bind.TransactOpts, deploy deployFunc[T]) (types.Address, error) {
	a, tx, _, err := deploy(txSubmitter, ethClient)
	if err != nil {
		return types.Address{}, err
	}

	fmt.Printf("Waiting for %s deployment confirmation\n", name)
	_, err = bind.WaitMined(ctx, ethClient, tx)
	if err != nil {
		return types.Address{}, err
	}
	fmt.Printf("%s successfully deployed to %s\n", name, a.String())
	return a, nil
}

func transferEth(client *ethclient.Client, txSubmitter *bind.TransactOpts, to types.Address, amount *big.Int) {
	slog.Info("Transferring eth", "to", to, "from", txSubmitter.From, "amount", amount)
	nonce, err := client.PendingNonceAt(context.Background(), txSubmitter.From)
	if err != nil {
		log.Fatal(err)
	}

	gasLimit := uint64(21000)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	tx := ethTypes.NewTransaction(nonce, to, amount, gasLimit, gasPrice, nil)
	signedTx, err := txSubmitter.Signer(txSubmitter.From, tx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}
}
