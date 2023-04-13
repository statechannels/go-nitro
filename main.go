package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
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
	var useWs, useDurableStore bool

	flag.BoolVar(&useWs, "useWS", false, "Specifies whether to use websockets or NATS for the rpc server.")
	flag.BoolVar(&useDurableStore, "useDurableStore", false, "Specifies whether to use a durable store or an in-memory store.")
	flag.StringVar(&pkString, "pk", "2d999770f7b5d49b694080f987b82bbc9fc9ac2b4dcc10b0f8aba7d700f69c6d", "Specifies the private key for the client. Default is Alice's private key.")
	flag.StringVar(&hardhatUrl, "hardhatUrl", "ws://127.0.0.1:8545", "Specifies the url for the hardhat node.")
	flag.IntVar(&msgPort, "msgPort", 3005, "Specifies the tcp port for the  message service.")
	flag.IntVar(&rpcPort, "rpcPort", 4005, "Specifies the tcp port for the rpc server.")
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
	chainService := NewChainService(context.Background(), *ourStore.GetAddress(), hardhatUrl)
	messageservice := p2pms.NewMessageService("127.0.0.1", msgPort, *ourStore.GetAddress(), pk)

	node := client.New(
		messageservice,
		chainService,
		ourStore,
		logDestination,
		&engine.PermissivePolicy{},
		nil)

	var transport transport.Responder
	var err error
	if useWs {
		transport, err = ws.NewWebSocketTransportAsServer(fmt.Sprint(rpcPort))
	} else {
		transport, err = nats.NewNatsTransportAsServer(rpcPort)
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

func NewChainService(ctx context.Context, address types.Address, hardhatUrl string) chainservice.ChainService {
	client, err := ethclient.Dial(hardhatUrl)
	if err != nil {
		log.Fatal(err)
	}
	key := GetHardhatFundedPrivateKey(address)
	fmt.Printf("Using private key %s", common.Bytes2Hex(ethcrypto.FromECDSA(key)))
	txSubmitter, err := bind.NewKeyedTransactorWithChainID(key, big.NewInt(1337))
	if err != nil {
		log.Fatal(err)
	}
	txSubmitter.GasLimit = uint64(30_000_000) // in units

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	txSubmitter.GasPrice = gasPrice

	deployer, err := Create2Deployer.NewCreate2Deployer(common.HexToAddress("0x5fbdb2315678afecb367f032d93f642f64180aa3"), client)
	if err != nil {
		log.Fatal(err)
	}

	hexBytecode, err := hex.DecodeString(NitroAdjudicator.NitroAdjudicatorMetaData.Bin[2:])
	if err != nil {
		log.Fatal(err)
	}

	naAddress, err := deployer.ComputeAddress(&bind.CallOpts{}, [32]byte{}, ethcrypto.Keccak256Hash(hexBytecode))
	if err != nil {
		log.Fatal(err)
	}
	bytecode, err := client.CodeAt(ctx, naAddress, nil) // nil is latest block
	if err != nil {
		log.Fatal(err)
	}

	// Has NitroAdjudicator been deployed? If not, deploy it.
	if len(bytecode) == 0 {
		_, err = deployer.Deploy(txSubmitter, big.NewInt(0), [32]byte{}, hexBytecode)
		if err != nil {
			log.Fatal(err)
		}
	}

	na, err := NitroAdjudicator.NewNitroAdjudicator(naAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	cs, err := chainservice.NewEthChainService(client, na, naAddress, common.Address{}, common.Address{}, txSubmitter, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
	return cs
}

func GetHardhatFundedPrivateKey(a types.Address) *ecdsa.PrivateKey {
	// See https://hardhat.org/hardhat-network/docs/reference#accounts for defaults
	// This is the default mnemonic used by hardhat
	const HARDHAT_MNEMONIC = "test test test test test test test test test test test junk"
	// We manually set the amount of funded accounts in our hardhat config
	// If that value changes, this value must change as well
	const NUM_FUNDED = 1000
	// This is the default hd wallet path used by hardhat
	const HD_PATH = "m/44'/60'/0'/0"

	index := big.NewInt(0).Mod(a.Big(), big.NewInt(NUM_FUNDED)).Uint64()

	return derivePrivateKey(uint(index), HARDHAT_MNEMONIC, HD_PATH, NUM_FUNDED)
}

func derivePrivateKey(index uint, mnemonic string, path string, numFunded uint) *ecdsa.PrivateKey {
	if numFunded < index {
		panic(fmt.Errorf("only the first %d accounts are funded", numFunded))
	}

	wallet, err := hdwallet.NewFromMnemonic(mnemonic)

	ourPath := fmt.Sprintf("%s/%d", path, index)
	if err != nil {
		panic(err)
	}

	a, err := wallet.Derive(hdwallet.MustParseDerivationPath(ourPath), false)
	if err != nil {
		panic(err)
	}
	pk, err := wallet.PrivateKey(a)
	if err != nil {
		panic(err)
	}
	return pk
}
