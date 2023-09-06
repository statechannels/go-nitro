package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/statechannels/go-nitro/internal/chain"
	"github.com/urfave/cli/v2"
)

const (
	FUNDED_TEST_PK  = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	ANVIL_CHAIN_URL = "ws://127.0.0.1:8545"
)

const (
	CHAIN_AUTH_TOKEN = "chainauthtoken"
	CHAIN_URL        = "chainurl"
	DEPLOYER_PK      = "chainpk"
)

func main() {
	flags := []cli.Flag{
		&cli.StringFlag{
			Name:    CHAIN_AUTH_TOKEN,
			Usage:   "Specifies the auth token for the chain",
			Value:   "",
			Aliases: []string{"ct"},
		},
		&cli.StringFlag{
			Name:    CHAIN_URL,
			Usage:   "Specifies the chain url to use",
			Value:   ANVIL_CHAIN_URL,
			Aliases: []string{"cu"},
		},
		&cli.StringFlag{
			Name:     DEPLOYER_PK,
			Usage:    "Specifies the private key to use when deploying contracts",
			Category: "Keys:",
			Aliases:  []string{"dpk"},
			Value:    FUNDED_TEST_PK,
		},
	}

	app := &cli.App{
		Name:  "start-rpc-servers",
		Usage: "Nitro as a service. State channel node with RPC server.",
		Flags: flags,

		Action: func(cCtx *cli.Context) error {
			chainAuthToken := cCtx.String(CHAIN_AUTH_TOKEN)
			chainUrl := cCtx.String(CHAIN_URL)
			chainPk := cCtx.String(DEPLOYER_PK)

			_, _, _, err := chain.DeployContracts(context.Background(), chainUrl, chainAuthToken, chainPk)
			if err != nil {
				panic(err)
			}

			fmt.Println("Successfully completed!")
			return nil
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
