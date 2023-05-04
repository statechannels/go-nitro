package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/statechannels/go-nitro/client"
	"github.com/statechannels/go-nitro/client/engine"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	NitroAdjudicator "github.com/statechannels/go-nitro/client/engine/chainservice/adjudicator"
	chainutils "github.com/statechannels/go-nitro/client/engine/chainservice/utils"
	p2pms "github.com/statechannels/go-nitro/client/engine/messageservice/p2p-message-service"
	"github.com/statechannels/go-nitro/client/engine/store"
	"github.com/statechannels/go-nitro/crypto"
	"github.com/statechannels/go-nitro/rpc"
	"github.com/statechannels/go-nitro/rpc/transport"
	"github.com/statechannels/go-nitro/rpc/transport/nats"
	"github.com/statechannels/go-nitro/rpc/transport/ws"
	"github.com/tidwall/buntdb"
)

func main() {
	var pkString, chainUrl, naAddress string
	var msgPort, rpcPort, chainId int
	var useNats, useDurableStore, deployContracts bool

	flag.BoolVar(&deployContracts, "deploycontracts", false, "Specifies whether to deploy the adjudicator and create2deployer contracts.")
	flag.BoolVar(&useNats, "usenats", false, "Specifies whether to use NATS or http/ws for the rpc server.")
	flag.BoolVar(&useDurableStore, "usedurablestore", false, "Specifies whether to use a durable store or an in-memory store.")
	flag.StringVar(&pkString, "pk", "2d999770f7b5d49b694080f987b82bbc9fc9ac2b4dcc10b0f8aba7d700f69c6d", "Specifies the private key for the client. Default is Alice's private key.")
	flag.StringVar(&chainUrl, "chainurl", "ws://127.0.0.1:8545", "Specifies the url of a RPC endpoint for the chain.")
	flag.StringVar(&naAddress, "naaddress", "0xC6A55E07566416274dBF020b5548eecEdB56290c", "Specifies the address of the nitro adjudicator contract. Default is the address computed by the Create2Deployer contract.")
	flag.IntVar(&msgPort, "msgport", 3005, "Specifies the tcp port for the  message service.")
	flag.IntVar(&rpcPort, "rpcport", 4005, "Specifies the tcp port for the rpc server.")
	flag.IntVar(&chainId, "chainid", 1337, "Specifies the chain id of the chain.")
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

	chainPk, err := chainutils.GetFundedTestPrivateKey(*ourStore.GetAddress())
	if err != nil {
		panic(err)
	}
	ethClient, txSubmitter, err := chainutils.ConnectToChain(context.Background(), chainUrl, chainId, chainPk)
	if err != nil {
		panic(err)
	}
	if deployContracts {
		deployedAddress, err := chainutils.DeployAdjudicator(context.Background(), ethClient, txSubmitter)
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
