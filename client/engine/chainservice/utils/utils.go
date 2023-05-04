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

// DeployAdjudicator deploys the Create2Deployer and NitroAdjudicator contracts.
// The nitro adjudicator is deployed to the address computed by the Create2Deployer contract.
func DeployAdjudicator(ctx context.Context, client *ethclient.Client, txSubmitter *bind.TransactOpts) (common.Address, error) {
	_, _, deployer, err := Create2Deployer.DeployCreate2Deployer(txSubmitter, client)
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
	key, err := ethcrypto.ToECDSA(chainPK)
	if err != nil {
		return nil, nil, err
	}
	txSubmitter, err := bind.NewKeyedTransactorWithChainID(key, big.NewInt(1337))
	if err != nil {
		return nil, nil, err
	}
	txSubmitter.GasLimit = uint64(30_000_000) // in units

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, nil, err
	}
	txSubmitter.GasPrice = gasPrice
	return client, txSubmitter, nil
}

// GetFundedTestPrivateKey selects a private key from one of 10 derived accounts using the
//
//	"test test test test test test test test test test test junk" mnemonic.
//
// Most test chains(hardhat/anvil) fund the first 10 accounts.
func GetFundedTestPrivateKey(a types.Address) ([]byte, error) {
	const TEST_MNEMONIC = "test test test test test test test test test test test junk"
	const NUM_FUNDED = 10
	const HD_PATH = "m/44'/60'/0'/0"

	index := big.NewInt(0).Mod(a.Big(), big.NewInt(NUM_FUNDED)).Uint64()

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
