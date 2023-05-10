package main

import (
	"context"
	"fmt"
	"log"
	"os"

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
	"github.com/urfave/cli/v2"
)

func main() {
	const (
		USE_NATS          = "usenats"
		USE_DURABLE_STORE = "usedureablestore"
		PK                = "pk"
		CHAIN_URL         = "chainurl"
		CHAIN_PK          = "chainpk"
		NA_ADDRESS        = "naaddress"
		MSG_PORT          = "msgport"
		RPC_PORT          = "rpcport"
		CHAIN_ID          = "chainid"
	)
	app := &cli.App{
		Name:  "go-nitro",
		Usage: "Nitro as a service. State channel client with RPC server.",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:     USE_NATS,
				Usage:    "Specifies whether to use NATS or http/ws for the rpc server.",
				Category: "Connectivity:",
			},
			&cli.BoolFlag{
				Name:     USE_DURABLE_STORE,
				Usage:    "Specifies whether to use a durable store or an in-memory store.",
				Category: "Storage",
			},
			&cli.StringFlag{
				Name:        PK,
				Usage:       "Specifies the private key for the client. Default is Alice's private key.",
				DefaultText: "2d999770f7b5d49b694080f987b82bbc9fc9ac2b4dcc10b0f8aba7d700f69c6d",
				Category:    "Keys:",
			},
			&cli.StringFlag{
				Name:        CHAIN_URL,
				Usage:       "Specifies the url of a RPC endpoint for the chain.",
				DefaultText: "ws://127.0.0.1:8545",
				Category:    "Connectivity:",
			},
			&cli.StringFlag{
				Name:        CHAIN_PK,
				Usage:       "Specifies the private key to use when interacting with the chain. Default is a hardhat/anvil funded account.",
				DefaultText: "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
				Category:    "Keys:",
			},
			&cli.StringFlag{
				Name:        NA_ADDRESS,
				Usage:       "Specifies the address of the nitro adjudicator contract.",
				DefaultText: "0xC6A55E07566416274dBF020b5548eecEdB56290c",
				Category:    "Connectivity:",
			},
			&cli.IntFlag{
				Name:        MSG_PORT,
				Usage:       "Specifies the tcp port for the message service.",
				DefaultText: "3005",
				Category:    "Connectivity:",
			},
			&cli.IntFlag{
				Name:        RPC_PORT,
				Usage:       "Specifies the tcp port for the rpc server.",
				DefaultText: "4005",
				Category:    "Connectivity:",
			},
			&cli.IntFlag{
				Name:        CHAIN_ID,
				Usage:       "Specifies the chain id of the chain.",
				DefaultText: "1337",
				Category:    "Connectivity:",
			},
		},
		Action: func(cCtx *cli.Context) error {
			pkString := cCtx.String(PK)
			chainUrl := cCtx.String(CHAIN_URL)
			naAddress := cCtx.String(NA_ADDRESS)
			chainPk := cCtx.String(CHAIN_PK)

			msgPort := cCtx.Int(MSG_PORT)
			rpcPort := cCtx.Int(RPC_PORT)
			chainId := cCtx.Int(CHAIN_ID)

			useNats := cCtx.Bool(USE_NATS)
			useDurableStore := cCtx.Bool(USE_DURABLE_STORE)

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

			ethClient, txSubmitter, err := chainutils.ConnectToChain(context.Background(), chainUrl, chainId, common.Hex2Bytes(chainPk))
			if err != nil {
				panic(err)
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
				return err
			}

			fmt.Println("Nitro as a Service listening on port", rpcPort)
			// NOT SURE IF WE NEED THIS?
			// sigs := make(chan os.Signal, 1)
			// signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
			// sig := <-sigs
			// fmt.Printf("Received signal %s, exiting..", sig)

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
