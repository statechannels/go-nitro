package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/statechannels/go-nitro/reverseproxy"
	"github.com/urfave/cli/v2"
)

const (
	NITRO_ENDPOINT  = "nitroendpoint"
	PORT            = "port"
	DESTINATION_URL = "destinationurl"
)

func main() {
	app := &cli.App{
		Name:  "reverse-payment-proxy",
		Usage: "Runs an HTTP payment proxy that charges for HTTP requests",
		Flags: []cli.Flag{
			&cli.UintFlag{
				Name:  PORT,
				Usage: "Specifies the port to run the proxy on",
				Value: 5511,
			},
			&cli.StringFlag{
				Name:  NITRO_ENDPOINT,
				Usage: "Specifies the endpoint of the Nitro RPC server",
				Value: "localhost:4007/api/v1",
			},
			&cli.StringFlag{
				Name:  DESTINATION_URL,
				Usage: "Specifies the url to forward requests to",
				Value: "http://localhost:8081",
			},
		},
		Action: func(c *cli.Context) error {
			proxyPort := c.Uint(PORT)
			nitroEndpoint := c.String(NITRO_ENDPOINT)
			p := reverseproxy.NewReversePaymentProxy(proxyPort, nitroEndpoint, c.String(DESTINATION_URL))

			return p.Start(context.Background())
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
	waitForKillSignal()
}

// waitForKillSignal blocks until we receive a kill or interrupt signal
func waitForKillSignal() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigs
	fmt.Printf("Received signal %s, exiting..\n", sig)
}
