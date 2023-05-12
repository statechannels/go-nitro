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
	"github.com/urfave/cli/v2/altsrc"
)

func main() {
	const (
		CONFIG            = "config"
		USE_NATS          = "usenats"
		USE_DURABLE_STORE = "usedurablestore"
		PK                = "pk"
		CHAIN_URL         = "chainurl"
		CHAIN_PK          = "chainpk"
		NA_ADDRESS        = "naaddress"
		MSG_PORT          = "msgport"
		RPC_PORT          = "rpcport"
		CHAIN_ID          = "chainid"
	)
	var pkString, chainUrl, naAddress, chainPk string
	var msgPort, rpcPort, chainId int
	var useNats, useDurableStore bool

	flags := []cli.Flag{
		&cli.StringFlag{
			Name:  CONFIG,
			Usage: "Load config options from `config.toml`",
		},
		altsrc.NewBoolFlag(&cli.BoolFlag{
			Name:        USE_NATS,
			Usage:       "Specifies whether to use NATS or http/ws for the rpc server.",
			Value:       false,
			Category:    "Connectivity:",
			Destination: &useNats,
		}),
		altsrc.NewBoolFlag(&cli.BoolFlag{
			Name:        USE_DURABLE_STORE,
			Usage:       "Specifies whether to use a durable store or an in-memory store.",
			Category:    "Storage",
			Value:       false,
			Destination: &useDurableStore,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        PK,
			Usage:       "Specifies the private key for the client.",
			Category:    "Keys:",
			Destination: &pkString,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        CHAIN_URL,
			Usage:       "Specifies the url of a RPC endpoint for the chain.",
			Value:       "ws://127.0.0.1:8545",
			DefaultText: "hardhat / anvil default",
			Category:    "Connectivity:",
			Destination: &chainUrl,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        CHAIN_PK,
			Usage:       "Specifies the private key to use when interacting with the chain.",
			Category:    "Keys:",
			Destination: &chainPk,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        NA_ADDRESS,
			Usage:       "Specifies the address of the nitro adjudicator contract.",
			Category:    "Connectivity:",
			Destination: &naAddress,
			Required:    true,
		}),
		altsrc.NewIntFlag(&cli.IntFlag{
			Name:        MSG_PORT,
			Usage:       "Specifies the tcp port for the message service.",
			Value:       3005,
			Category:    "Connectivity:",
			Destination: &msgPort,
		}),
		altsrc.NewIntFlag(&cli.IntFlag{
			Name:        RPC_PORT,
			Usage:       "Specifies the tcp port for the rpc server.",
			Value:       4005,
			Category:    "Connectivity:",
			Destination: &rpcPort,
		}),
		altsrc.NewIntFlag(&cli.IntFlag{
			Name:        CHAIN_ID,
			Usage:       "Specifies the chain id of the chain.",
			Value:       1337,
			DefaultText: "hardhat default",
			Category:    "Connectivity:",
			Destination: &chainId,
		}),
	}
	app := &cli.App{
		Name:   "go-nitro",
		Usage:  "Nitro as a service. State channel client with RPC server.",
		Flags:  flags,
		Before: altsrc.InitInputSourceWithContext(flags, altsrc.NewTomlSourceFromFlagFunc(CONFIG)),
		Action: func(cCtx *cli.Context) error {
			if pkString == "" {
				panic("pk must be set")
			}
			if chainPk == "" {
				panic("chainpk must be set")
			}
			pk := common.Hex2Bytes(pkString)
			me := crypto.GetAddressFromSecretKeyBytes(pk)

			logDestination := os.Stdout

			var ourStore store.Store

			if useDurableStore {
				fmt.Println("Initialising durable store...")
				dataFolder := fmt.Sprintf("./data/nitro-service/%s", me.String())
				ourStore = store.NewDurableStore(pk, dataFolder, buntdb.Config{})
			} else {
				fmt.Println("Initialising mem store...")
				ourStore = store.NewMemStore(pk)
			}

			fmt.Println("Connecting to chain " + chainUrl + " with chain id " + fmt.Sprint(chainId) + "...")
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

			fmt.Println("Initializing message service on port " + fmt.Sprint(msgPort) + "...")
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
				fmt.Println("Initializing NATS RPC transport...")
				transport, err = nats.NewNatsTransportAsServer(rpcPort)
			} else {
				fmt.Println("Initializing websocketNATS RPC transport...")
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

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
