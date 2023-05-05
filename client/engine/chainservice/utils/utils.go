package chainutils

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	NitroAdjudicator "github.com/statechannels/go-nitro/client/engine/chainservice/adjudicator"
	Create2Deployer "github.com/statechannels/go-nitro/client/engine/chainservice/create2deployer"
	"github.com/statechannels/go-nitro/types"
)

const (

	// deployerIndex is the index of the HD wallet account used to deploy contracts.
	// Using the same account for deployments ensures that the same address is generated for the Create2Deployer contract.
	deployerIndex = 0
	// This is the expected address of the Create2Deployer contract.
	deployAddress = "0x5FbDB2315678afecb367f032d93F642f64180aa3"
)

// DeployAdjudicator deploys the Create2Deployer and NitroAdjudicator contracts.
// The nitro adjudicator is deployed to the address computed by the Create2Deployer contract.
// TODO: The *bind.TransactOpts param is ignored. It's only there for compatibility with the testground test.
// We should remove it when the testground test is updated.
func DeployAdjudicator(ctx context.Context, client *ethclient.Client, _ *bind.TransactOpts) (common.Address, error) {
	// We want to use the same account for deployments
	// so we use one specific account when handling deployments.
	chainPk, err := deriveFundedPk(deployerIndex)
	if err != nil {
		return types.Address{}, err
	}

	txSubmitter, err := createTxSubmitter(context.Background(), chainPk, client)
	if err != nil {
		return types.Address{}, err
	}
	deployer, err := getCreate2Deployer(ctx, client, txSubmitter)
	if err != nil {
		return types.Address{}, err
	}
	hexBytecode, err := hex.DecodeString(NitroAdjudicator.NitroAdjudicatorMetaData.Bin[2:])
	if err != nil {
		return types.Address{}, err
	}

	naAddress, err := deployer.ComputeAddress(&bind.CallOpts{}, [32]byte{}, ethcrypto.Keccak256Hash(hexBytecode))
	if err != nil {
		return types.Address{}, err
	}
	bytecode, err := client.CodeAt(ctx, naAddress, nil) // nil is latest block
	if err != nil {
		return types.Address{}, err
	}

	// Has NitroAdjudicator been deployed? If not, deploy it.
	if len(bytecode) == 0 {
		_, err = deployer.Deploy(txSubmitter, big.NewInt(0), [32]byte{}, hexBytecode)
		if err != nil {
			return types.Address{}, err
		}
	}
	return naAddress, nil
}

// getCreate2Deployer deploys or connects to the Create2Deployer contract at a known address
// If the contract is not deployed, it deploys it.
func getCreate2Deployer(ctx context.Context, client *ethclient.Client, txSubmitter *bind.TransactOpts) (*Create2Deployer.Create2Deployer, error) {
	bytecode, err := client.CodeAt(ctx, common.HexToAddress(deployAddress), nil) // nil is latest block
	if err != nil {
		return nil, err
	}
	if len(bytecode) == 0 {

		deployedAdd, _, deployer, err := Create2Deployer.DeployCreate2Deployer(txSubmitter, client)
		if err != nil {
			return nil, err
		}
		if deployedAdd != common.HexToAddress(deployAddress) {
			return nil, fmt.Errorf("deployed address %v does not match expected address %v", deployedAdd, deployAddress)
		}
		return deployer, nil
	} else {
		return Create2Deployer.NewCreate2Deployer(common.HexToAddress(deployAddress), client)
	}
}

func createTxSubmitter(ctx context.Context, chainPK []byte, client *ethclient.Client) (*bind.TransactOpts, error) {
	key, err := ethcrypto.ToECDSA(chainPK)
	if err != nil {
		return nil, err
	}

	chainId, err := client.ChainID(ctx)
	if err != nil {
		return nil, err
	}

	txSubmitter, err := bind.NewKeyedTransactorWithChainID(key, chainId)
	if err != nil {
		return nil, err
	}
	txSubmitter.GasLimit = uint64(30_000_000) // in units

	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, err
	}
	txSubmitter.GasPrice = gasPrice
	return txSubmitter, nil
}

// ConnectToChain connects to the chain at the given url and returns a client and a transactor.
func ConnectToChain(ctx context.Context, chainUrl string, chainId int, chainPK []byte) (*ethclient.Client, *bind.TransactOpts, error) {
	client, err := ethclient.Dial(chainUrl)
	if err != nil {
		return nil, nil, err
	}
	foundChainId, err := client.ChainID(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("could not get chain id: %w", err)
	}
	if foundChainId.Cmp(big.NewInt(int64(chainId))) != 0 {
		return nil, nil, fmt.Errorf("chain id mismatch: expected %d, got %d", chainId, foundChainId)
	}
	txSubmitter, err := createTxSubmitter(ctx, chainPK, client)
	if err != nil {
		return nil, nil, fmt.Errorf("could not construct tx submitter: %w", err)
	}
	return client, txSubmitter, nil
}

const (
	TEST_MNEMONIC = "test test test test test test test test test test test junk"
	NUM_FUNDED    = 10
	HD_PATH       = "m/44'/60'/0'/0"
)

// deriveFundedPk derives a private key from the test mnemonic and the given index.
func deriveFundedPk(index uint64) ([]byte, error) {
	if index >= NUM_FUNDED {
		return nil, fmt.Errorf("index %d is larger than the number of funded accounts %d", index, NUM_FUNDED)
	}
	wallet, err := hdwallet.NewFromMnemonic(TEST_MNEMONIC)
	if err != nil {
		return nil, err
	}

	ourPath := fmt.Sprintf("%s/%d", HD_PATH, index)

	derived, err := wallet.Derive(hdwallet.MustParseDerivationPath(ourPath), false)
	if err != nil {
		return nil, err
	}
	pk, err := wallet.PrivateKey(derived)
	if err != nil {
		return nil, err
	}
	return ethcrypto.FromECDSA(pk), nil
}

// GetFundedTestPrivateKey selects a private key from one of 10 derived accounts using the
//
//	"test test test test test test test test test test test junk" mnemonic.
//
// Most test chains(hardhat/anvil) fund the first 10 accounts.
func GetFundedTestPrivateKey(a types.Address) ([]byte, error) {
	index := big.NewInt(0).Mod(a.Big(), big.NewInt(NUM_FUNDED)).Uint64()
	return deriveFundedPk(index)
}
