package chainutils

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// ConnectToChain connects to the chain at the given url and returns a client and a transactor.
func ConnectToChain(ctx context.Context, chainUrl string, chainPK []byte) (*ethclient.Client, *bind.TransactOpts, error) {
	client, err := ethclient.Dial(chainUrl)
	if err != nil {
		return nil, nil, err
	}
	foundChainId, err := client.ChainID(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("could not get chain id: %w", err)
	}

	key, err := ethcrypto.ToECDSA(chainPK)
	if err != nil {
		return nil, nil, err
	}
	txSubmitter, err := bind.NewKeyedTransactorWithChainID(key, foundChainId)
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
