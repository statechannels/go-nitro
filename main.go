package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/internal/chain"
	"github.com/statechannels/go-nitro/internal/rpc"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

func main() {
	const (
		CONFIG               = "config"
		USE_NATS             = "usenats"
		USE_DURABLE_STORE    = "usedurablestore"
		PK                   = "pk"
		CHAIN_URL            = "chainurl"
		CHAIN_AUTH_TOKEN     = "chainauthtoken"
		CHAIN_PK             = "chainpk"
		NA_ADDRESS           = "naaddress"
		VPA_ADDRESS          = "vpaaddress"
		CA_ADDRESS           = "caaddress"
		MSG_PORT             = "msgport"
		RPC_PORT             = "rpcport"
		DURABLE_STORE_FOLDER = "durablestorefolder"
	)
	var pkString, chainUrl, chainAuthToken, naAddress, vpaAddress, caAddress, chainPk, durableStoreFolder string
	var msgPort, rpcPort int
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
			Category:    "Storage:",
			Value:       false,
			Destination: &useDurableStore,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        PK,
			Usage:       "Specifies the private key used by the nitro node.",
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
			Name:        CHAIN_AUTH_TOKEN,
			Usage:       "The bearer token used for auth when making requests to the chain's RPC endpoint.",
			Category:    "Connectivity:",
			Destination: &chainAuthToken,
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
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        VPA_ADDRESS,
			Usage:       "Specifies the address of the virtual payment app.",
			Category:    "Connectivity:",
			Destination: &vpaAddress,
			Required:    true,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        CA_ADDRESS,
			Usage:       "Specifies the address of the consensus app.",
			Category:    "Connectivity:",
			Destination: &caAddress,
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
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        DURABLE_STORE_FOLDER,
			Usage:       "Specifies the folder for the durable store data storage.",
			Category:    "Storage:",
			Destination: &durableStoreFolder,
			Value:       "./data/nitro-store",
		}),
	}
	app := &cli.App{
		Name:   "go-nitro",
		Usage:  "Nitro as a service. State channel node with RPC server.",
		Flags:  flags,
		Before: altsrc.InitInputSourceWithContext(flags, altsrc.NewTomlSourceFromFlagFunc(CONFIG)),
		Action: func(cCtx *cli.Context) error {
			chainOpts := chain.ChainOpts{
				ChainUrl:       chainUrl,
				ChainAuthToken: chainAuthToken,
				ChainPk:        chainPk,
				NaAddress:      common.HexToAddress(naAddress),
				VpaAddress:     common.HexToAddress(vpaAddress),
				CaAddress:      common.HexToAddress(caAddress),
			}

			rpcServer, _, _, err := rpc.InitChainServiceAndRunRpcServer(pkString, chainOpts, useDurableStore, durableStoreFolder, useNats, msgPort, rpcPort)
			if err != nil {
				return err
			}

			HostNitroUI(uint(rpcPort))
			stopChan := make(chan os.Signal, 2)
			signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
			<-stopChan // wait for interrupt or terminate signal

			return rpcServer.Close()
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
