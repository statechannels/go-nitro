package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"flag"
	"fmt"
	"math/big"
	"os"
	"os/signal"
	"syscall"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/rs/zerolog"
	"github.com/statechannels/go-nitro/client"
	"github.com/statechannels/go-nitro/client/engine"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	NitroAdjudicator "github.com/statechannels/go-nitro/client/engine/chainservice/adjudicator"
	Create2Deployer "github.com/statechannels/go-nitro/client/engine/chainservice/create2deployer"
	p2pms "github.com/statechannels/go-nitro/client/engine/messageservice/p2p-message-service"
	"github.com/statechannels/go-nitro/client/engine/store"
	"github.com/statechannels/go-nitro/crypto"
	"github.com/statechannels/go-nitro/rpc"
	"github.com/statechannels/go-nitro/rpc/transport"
	"github.com/statechannels/go-nitro/rpc/transport/nats"
	"github.com/statechannels/go-nitro/rpc/transport/ws"
	"github.com/statechannels/go-nitro/types"
	"github.com/tidwall/buntdb"
)

func main() {
	var pkString, chainUrl, naAddress string
	var msgPort, rpcPort int
	var useNats, useDurableStore, deployContracts bool

	flag.BoolVar(&deployContracts, "deploycontracts", false, "Specifies whether to deploy the adjudicator and create2deployer contracts.")
	flag.BoolVar(&useNats, "usenats", true, "Specifies whether to use NATS or http/ws for the rpc server.")
	flag.BoolVar(&useDurableStore, "usedurablestore", false, "Specifies whether to use a durable store or an in-memory store.")
	flag.StringVar(&pkString, "pk", "2d999770f7b5d49b694080f987b82bbc9fc9ac2b4dcc10b0f8aba7d700f69c6d", "Specifies the private key for the client. Default is Alice's private key.")
	flag.StringVar(&chainUrl, "chainurl", "ws://127.0.0.1:8545", "Specifies the url of a RPC endpoint for the chain.")
	flag.StringVar(&naAddress, "naaddress", "0xC6A55E07566416274dBF020b5548eecEdB56290c", "Specifies the address of the nitro adjudicator contract. Default is the address computed by the Create2Deployer contract.")
	flag.IntVar(&msgPort, "msgport", 3005, "Specifies the tcp port for the  message service.")
	flag.IntVar(&rpcPort, "rpcport", 4005, "Specifies the tcp port for the rpc server.")

	flag.Parse()

	pk := common.Hex2Bytes(pkString)
	me := crypto.GetAddressFromSecretKeyBytes(pk)

	logDestination := os.Stdout

	var ourStore store.Store
	if useDurableStore {
		dataFolder := fmt.Sprintf("./data/nitro-service/%s", me.String())
		ourStore = store.NewDurableStore(pk, dataFolder, buntdb.Config{})
	} else {
		ourStore = store.NewMemStore(pk)
	}

	ethClient, txSubmitter, err := connectToChain(context.Background(), chainUrl, *ourStore.GetAddress())
	if err != nil {
		panic(err)
	}
	if deployContracts {
		deployedAddress, err := deployAdjudicator(context.Background(), ethClient, txSubmitter)
		if err != nil {
			panic(err)
		}
		if naAddress != deployedAddress.String() {
			fmt.Printf("WARNING: The deploycontracts flag is set so the adjucator has been deployed to %s.\nThis is different from the naaddress flag which is set to %s. The naaddress flag will be ignored.\n", deployedAddress.String(), naAddress)
			naAddress = deployedAddress.String()
		}
	}

	na, err := NitroAdjudicator.NewNitroAdjudicator(common.HexToAddress(naAddress), ethClient)
	if err != nil {
		panic(err)
	}

	chainService, err := chainservice.NewEthChainService(ethClient, na, common.HexToAddress(naAddress), common.Address{}, common.Address{}, txSubmitter, os.Stdout)
	if err != nil {
		panic(err)
	}

	messageservice := p2pms.NewMessageService("127.0.0.1", msgPort, *ourStore.GetAddress(), pk, logDestination)
	node := client.New(
		messageservice,
		chainService,
		ourStore,
		logDestination,
		&engine.PermissivePolicy{},
		nil)

	var transport transport.Responder

	if useNats {
		transport, err = nats.NewNatsTransportAsServer(rpcPort)
	} else {
		transport, err = ws.NewWebSocketTransportAsServer(fmt.Sprint(rpcPort))
	}
	if err != nil {
		panic(err)
	}

	logger := zerolog.New(logDestination).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Str("client", ourStore.GetAddress().String()).
		Str("rpc", "server").
		Str("scope", "").
		Logger()
	_, err = rpc.NewRpcServer(&node, &logger, transport)
	if err != nil {
		panic(err)
	}

	fmt.Println("Nitro as a Service listening on port", rpcPort)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigs
	fmt.Printf("Received signal %s, exiting..", sig)
}

// deployAdjudicator deploys the Create2Deployer and NitroAdjudicator contracts.
// The nitro adjudicator is deployed to the address computed by the Create2Deployer contract.
func deployAdjudicator(ctx context.Context, client *ethclient.Client, txSubmitter *bind.TransactOpts) (common.Address, error) {
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

// connectToChain connects to the chain at the given url and returns a client and a transactor.
func connectToChain(ctx context.Context, chainUrl string, myAddress types.Address) (*ethclient.Client, *bind.TransactOpts, error) {
	client, err := ethclient.Dial(chainUrl)
	if err != nil {
		return nil, nil, err
	}
	key, err := getHardhatFundedPrivateKey(myAddress)
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

// getHardhatFundedPrivateKey selects a private key from one of the 1000 funded accounts in hardhat.
// It modulates the address by 1000 to select a funded account.
func getHardhatFundedPrivateKey(a types.Address) (*ecdsa.PrivateKey, error) {
	// See https://hardhat.org/hardhat-network/docs/reference#accounts for defaults
	// This is the default mnemonic used by hardhat
	const HARDHAT_MNEMONIC = "test test test test test test test test test test test junk"
	// We manually set the amount of funded accounts in our hardhat config
	// If that value changes, this value must change as well
	const NUM_FUNDED = 1000
	// This is the default hd wallet path used by hardhat
	const HD_PATH = "m/44'/60'/0'/0"

	index := big.NewInt(0).Mod(a.Big(), big.NewInt(NUM_FUNDED)).Uint64()

	wallet, err := hdwallet.NewFromMnemonic(HARDHAT_MNEMONIC)
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
	return pk, nil
}
