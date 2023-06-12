package chainutils

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

// ConnectToChain connects to the chain at the given url and returns a client and a transactor.
func ConnectToChain(ctx context.Context, chainUrl, chainAuthToken string, chainPK []byte) (*ethclient.Client, *bind.TransactOpts, error) {
	var rpcClient *rpc.Client
	var err error

	if chainAuthToken != "" {
		fmt.Println("Adding bearer token authorization header to chain service")
		options := rpc.WithHeader("Authorization", "Bearer "+chainAuthToken)
		rpcClient, err = rpc.DialOptions(ctx, chainUrl, options)
	} else {
		rpcClient, err = rpc.DialContext(ctx, chainUrl)
	}
	if err != nil {
		return nil, nil, err
	}

	client := ethclient.NewClient(rpcClient)
	fmt.Println("Connected to ethclient at url: " + chainUrl)

	foundChainId, err := client.ChainID(context.Background())
	if err != nil {
		return nil, nil, fmt.Errorf("could not get chain id: %w", err)
	}
	fmt.Printf("Found chain id: %v\n", foundChainId)

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
