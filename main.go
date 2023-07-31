package main

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/internal/chain"
	"github.com/statechannels/go-nitro/internal/rpc"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

func main() {
	const (
		CONFIG = "config"

		// Connectivity
		CONNECTIVITY_CATEGORY = "Connectivity:"
		USE_NATS              = "usenats"
		CHAIN_URL             = "chainurl"
		CHAIN_AUTH_TOKEN      = "chainauthtoken"
		NA_ADDRESS            = "naaddress"
		VPA_ADDRESS           = "vpaaddress"
		CA_ADDRESS            = "caaddress"
		MSG_PORT              = "msgport"
		RPC_PORT              = "rpcport"
		GUI_PORT              = "guiport"
		BOOT_PEERS            = "bootpeers"
		USE_MDNS              = "usemdns"

		// Keys
		KEYS_CATEGORY = "Keys:"
		PK            = "pk"
		CHAIN_PK      = "chainpk"

		// Storage
		STORAGE_CATEGORY     = "Storage:"
		USE_DURABLE_STORE    = "usedurablestore"
		DURABLE_STORE_FOLDER = "durablestorefolder"
	)
	var pkString, chainUrl, chainAuthToken, naAddress, vpaAddress, caAddress, chainPk, durableStoreFolder, bootPeers string
	var msgPort, rpcPort, guiPort int
	var useNats, useDurableStore, useMdns bool

	flags := []cli.Flag{
		//nolint:exhaustruct
		&cli.StringFlag{
			Name:  CONFIG,
			Usage: "Load config options from `config.toml`",
		},
		//nolint:exhaustruct
		altsrc.NewBoolFlag(&cli.BoolFlag{
			Name:        USE_NATS,
			Usage:       "Specifies whether to use NATS or http/ws for the rpc server.",
			Value:       false,
			Category:    CONNECTIVITY_CATEGORY,
			Destination: &useNats,
		}),
		//nolint:exhaustruct
		altsrc.NewBoolFlag(&cli.BoolFlag{
			Name:        USE_DURABLE_STORE,
			Usage:       "Specifies whether to use a durable store or an in-memory store.",
			Category:    STORAGE_CATEGORY,
			Value:       false,
			Destination: &useDurableStore,
		}),
		//nolint:exhaustruct
		altsrc.NewBoolFlag(&cli.BoolFlag{
			Name:        USE_MDNS,
			Usage:       "Specifies whether to use mDNS for peer discovery (if 'false', will use kademlia-dht)",
			Category:    CONNECTIVITY_CATEGORY,
			Value:       true,
			Destination: &useMdns,
		}),
		//nolint:exhaustruct
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        PK,
			Usage:       "Specifies the private key used by the nitro node.",
			Category:    KEYS_CATEGORY,
			Destination: &pkString,
		}),
		//nolint:exhaustruct
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        CHAIN_URL,
			Usage:       "Specifies the url of a RPC endpoint for the chain.",
			Value:       "ws://127.0.0.1:8545",
			DefaultText: "hardhat / anvil default",
			Category:    CONNECTIVITY_CATEGORY,
			Destination: &chainUrl,
		}),
		//nolint:exhaustruct
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        CHAIN_AUTH_TOKEN,
			Usage:       "The bearer token used for auth when making requests to the chain's RPC endpoint.",
			Category:    CONNECTIVITY_CATEGORY,
			Destination: &chainAuthToken,
		}),
		//nolint:exhaustruct
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        CHAIN_PK,
			Usage:       "Specifies the private key to use when interacting with the chain.",
			Category:    KEYS_CATEGORY,
			Destination: &chainPk,
		}),
		//nolint:exhaustruct
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        NA_ADDRESS,
			Usage:       "Specifies the address of the nitro adjudicator contract.",
			Category:    CONNECTIVITY_CATEGORY,
			Destination: &naAddress,
			Required:    true,
		}),
		//nolint:exhaustruct
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        VPA_ADDRESS,
			Usage:       "Specifies the address of the virtual payment app.",
			Category:    CONNECTIVITY_CATEGORY,
			Destination: &vpaAddress,
			Required:    true,
		}),
		//nolint:exhaustruct
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        CA_ADDRESS,
			Usage:       "Specifies the address of the consensus app.",
			Category:    CONNECTIVITY_CATEGORY,
			Destination: &caAddress,
			Required:    true,
		}),
		//nolint:exhaustruct
		altsrc.NewIntFlag(&cli.IntFlag{
			Name:        MSG_PORT,
			Usage:       "Specifies the tcp port for the message service.",
			Value:       3005,
			Category:    CONNECTIVITY_CATEGORY,
			Destination: &msgPort,
		}),
		//nolint:exhaustruct
		altsrc.NewIntFlag(&cli.IntFlag{
			Name:        RPC_PORT,
			Usage:       "Specifies the tcp port for the rpc server.",
			Value:       4005,
			Category:    CONNECTIVITY_CATEGORY,
			Destination: &rpcPort,
		}),
		//nolint:exhaustruct
		altsrc.NewIntFlag(&cli.IntFlag{
			Name:        GUI_PORT,
			Usage:       "Specifies the tcp port for the Nitro Connect GUI.",
			Value:       5005,
			Category:    CONNECTIVITY_CATEGORY,
			Destination: &guiPort,
		}),
		//nolint:exhaustruct
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        DURABLE_STORE_FOLDER,
			Usage:       "Specifies the folder for the durable store data storage.",
			Category:    STORAGE_CATEGORY,
			Destination: &durableStoreFolder,
			Value:       "./data/nitro-store",
		}),
		//nolint:exhaustruct
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        BOOT_PEERS,
			Usage:       "Comma-delimited list of peer multiaddrs the messaging service will connect to when initialized.",
			Value:       "",
			Category:    CONNECTIVITY_CATEGORY,
			Destination: &bootPeers,
		}),
	}
	//nolint:exhaustruct
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

			var peerSlice []string
			if bootPeers != "" {
				peerSlice = strings.Split(bootPeers, ",")
			}
			rpcServer, _, _, err := rpc.InitChainServiceAndRunRpcServer(pkString, chainOpts, useDurableStore, durableStoreFolder, useNats, msgPort, rpcPort, peerSlice, useMdns)
			if err != nil {
				return err
			}

			hostNitroUI(uint(guiPort), uint(rpcPort))

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
