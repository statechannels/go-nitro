package main

import (
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/cmd/utils"
	"github.com/statechannels/go-nitro/internal/logging"
	"github.com/statechannels/go-nitro/rpc"
	"github.com/statechannels/go-nitro/rpc/transport/http"
	"github.com/urfave/cli/v2"
)

const (
	COUNTERPARTY_ADDRESS = "counterpartyaddress"
	NITRO_ENDPOINT       = "nitroendpoint"
	LOG_FILE             = "create-single-channel.log"
	AMOUNT               = "amount"
)

func main() {
	logging.SetupDefaultFileLogger(LOG_FILE, slog.LevelDebug)
	flags := []cli.Flag{
		&cli.StringFlag{
			Name:    NITRO_ENDPOINT,
			Usage:   "Specifies the endpoint of the Nitro RPC server to connect to. This should be in the form 'host:port/api/v1'",
			Value:   "localhost:4005/api/v1",
			Aliases: []string{"n"},
		},
		&cli.StringFlag{
			Name:    COUNTERPARTY_ADDRESS,
			Usage:   "Specifies the address of the counterparty to create the ledger channel with.",
			Value:   "0x111A00868581f73AB42FEEF67D235Ca09ca1E8db",
			Aliases: []string{"c"},
		},
		&cli.UintFlag{
			Name:    AMOUNT,
			Value:   5_000_000,
			Usage:   "Specifies the amount of wei to deposit into the ledger channel.",
			Aliases: []string{"a"},
		},
	}

	app := &cli.App{
		Name:  "create-ledger-channel",
		Usage: "Creates a ledger channel with the specified counterparty and amount",
		Flags: flags,
		Action: func(cCtx *cli.Context) error {
			clientConnection, err := http.NewHttpTransportAsClient(cCtx.String(NITRO_ENDPOINT), 10*time.Millisecond)
			if err != nil {
				return err
			}
			client, err := rpc.NewRpcClient(clientConnection)
			if err != nil {
				return err
			}
			defer client.Close()

			err = utils.CreateLedgerChannel(client, common.HexToAddress(cCtx.String(COUNTERPARTY_ADDRESS)), cCtx.Uint(AMOUNT))
			if err != nil {
				return err
			}
			return nil
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
