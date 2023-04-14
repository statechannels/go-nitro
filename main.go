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
	var pkString, hardhatUrl string
	var msgPort, rpcPort int
	var useNats, useDurableStore bool

	flag.BoolVar(&useNats, "usenats", true, "Specifies whether to use NATS or http/ws for the rpc server.")
	flag.BoolVar(&useDurableStore, "usedurablestore", false, "Specifies whether to use a durable store or an in-memory store.")
	flag.StringVar(&pkString, "pk", "2d999770f7b5d49b694080f987b82bbc9fc9ac2b4dcc10b0f8aba7d700f69c6d", "Specifies the private key for the client. Default is Alice's private key.")
	flag.StringVar(&hardhatUrl, "hardhaturl", "ws://127.0.0.1:8545", "Specifies the url for the hardhat node.")
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

	chainService, err := newChainService(context.Background(), *ourStore.GetAddress(), hardhatUrl)
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

	fmt.Println("Nitro as a Service listening on localhost:", rpcPort)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigs
	fmt.Printf("Received signal %s, exiting..", sig)
}

// deployAdjudicator deploys the NitroAdjudicator contract if it has not already been deployed.
// It makes use of a Create2Deployer contract for deterministic deploy addresses.
func deployAdjudicator(ctx context.Context, client *ethclient.Client, txSubmitter *bind.TransactOpts) (common.Address, error) {
	deployer, err := Create2Deployer.NewCreate2Deployer(common.HexToAddress("0x5fbdb2315678afecb367f032d93f642f64180aa3"), client)
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

// newChainService constructs a chain service using the given node url.
func newChainService(ctx context.Context, address types.Address, nodeUrl string) (chainservice.ChainService, error) {
	client, err := ethclient.Dial(nodeUrl)
	if err != nil {
		return nil, err
	}
	key, err := getHardhatFundedPrivateKey(address)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Using private key %s", common.Bytes2Hex(ethcrypto.FromECDSA(key)))
	txSubmitter, err := bind.NewKeyedTransactorWithChainID(key, big.NewInt(1337))
	if err != nil {
		return nil, err
	}
	txSubmitter.GasLimit = uint64(30_000_000) // in units

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}
	txSubmitter.GasPrice = gasPrice

	naAddress, err := deployAdjudicator(ctx, client, txSubmitter)
	if err != nil {
		return nil, err
	}

	na, err := NitroAdjudicator.NewNitroAdjudicator(naAddress, client)
	if err != nil {
		return nil, err
	}

	cs, err := chainservice.NewEthChainService(client, na, naAddress, common.Address{}, common.Address{}, txSubmitter, os.Stdout)
	if err != nil {
		return nil, err
	}
	return cs, nil
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
