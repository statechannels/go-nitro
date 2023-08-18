package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/statechannels/go-nitro/cmd/utils"
	"github.com/statechannels/go-nitro/internal/chain"
	"github.com/statechannels/go-nitro/internal/testhelpers"
	"github.com/statechannels/go-nitro/types"
	"github.com/urfave/cli/v2"
)

type participant string

const (
	alice participant = "alice"
	bob   participant = "bob"
	irene participant = "irene"
	ivan  participant = "ivan"
)

type color string

const (
	black   color = "[30m"
	red     color = "[31m"
	green   color = "[32m"
	yellow  color = "[33m"
	blue    color = "[34m"
	magenta color = "[35m"
	cyan    color = "[36m"
	white   color = "[37m"
	gray    color = "[90m"
)

var participantColor = map[participant]color{alice: blue, irene: green, ivan: cyan, bob: yellow}

const (
	FUNDED_TEST_PK  = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	ANVIL_CHAIN_URL = "ws://127.0.0.1:8545"
)

const (
	CHAIN_AUTH_TOKEN = "chainauthtoken"
	CHAIN_URL        = "chainurl"
	DEPLOYER_PK      = "chainpk"
	START_ANVIL      = "startanvil"
	HOST_UI          = "hostui"
)

func main() {
	flags := []cli.Flag{
		&cli.BoolFlag{
			Name:    START_ANVIL,
			Usage:   "Specifies whether to start a local anvil instance",
			Value:   true,
			Aliases: []string{"a"},
		},
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
		&cli.BoolFlag{
			Name:    HOST_UI,
			Usage:   "Specifies whether to host the nitro UI or not",
			Value:   false,
			Aliases: []string{"ui"},
		},
	}

	app := &cli.App{
		Name:  "start-rpc-servers",
		Usage: "Nitro as a service. State channel node with RPC server.",
		Flags: flags,

		Action: func(cCtx *cli.Context) error {
			running := []*exec.Cmd{}
			if cCtx.Bool(START_ANVIL) {
				anvilCmd, err := chain.StartAnvil()
				if err != nil {
					utils.StopCommands(running...)
					panic(err)
				}
				running = append(running, anvilCmd)
			}
			dataFolder, cleanup := testhelpers.GenerateTempStoreFolder()
			defer cleanup()
			fmt.Printf("Using data folder %s\n", dataFolder)

			chainAuthToken := cCtx.String(CHAIN_AUTH_TOKEN)
			chainUrl := cCtx.String(CHAIN_URL)
			chainPk := cCtx.String(DEPLOYER_PK)

			naAddress, vpaAddress, caAddress, err := chain.DeployContracts(context.Background(), chainUrl, chainAuthToken, chainPk)
			if err != nil {
				utils.StopCommands(running...)
				panic(err)
			}

			hostUI := cCtx.Bool(HOST_UI)

			// Setup Ivan first, he is the DHT boot peer
			client, err := setupRPCServer(ivan, participantColor[ivan], naAddress, vpaAddress, caAddress, chainUrl, chainAuthToken, dataFolder, hostUI)
			if err != nil {
				utils.StopCommands(running...)
				panic(err)
			}
			running = append(running, client)

			const IVAN_ADDRESS = "http://127.0.0.1:4008/api/v1"
			err = waitForRpcClient(IVAN_ADDRESS, 500*time.Millisecond, 5*time.Minute)
			if err != nil {
				utils.StopCommands(running...)
				panic(err)
			}
			for _, p := range []participant{alice, bob, irene} {

				client, err := setupRPCServer(p, participantColor[p], naAddress, vpaAddress, caAddress, chainUrl, chainAuthToken, dataFolder, hostUI)
				if err != nil {
					utils.StopCommands(running...)
					panic(err)
				}
				running = append(running, client)
			}

			utils.WaitForKillSignal()
			utils.StopCommands(running...)
			return nil
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

// setupRPCServer starts up an RPC server for the given participant
func setupRPCServer(p participant, c color, na, vpa, ca types.Address, chainUrl, chainAuthToken string, dataFolder string, hostUI bool) (*exec.Cmd, error) {
	args := []string{"run"}

	if hostUI {
		args = append(args, "-tags", "embed_ui")
	}

	args = append(args, ".")
	args = append(args, "-naaddress", na.String())
	args = append(args, "-vpaaddress", vpa.String())
	args = append(args, "-caaddress", ca.String())

	args = append(args, "-chainauthtoken", chainAuthToken)
	args = append(args, "-chainurl", chainUrl)

	args = append(args, "-durablestorefolder", dataFolder)

	args = append(args, "-config", fmt.Sprintf("./cmd/test-configs/%s.toml", p))

	cmd := exec.Command("go", args...)
	cmd.Stdout = newColorWriter(c, os.Stdout)
	cmd.Stderr = newColorWriter(c, os.Stderr)
	err := cmd.Start()
	if err != nil {
		return nil, err
	}
	return cmd, nil
}

// colorWriter is a writer that writes to the underlying writer with the given color
type colorWriter struct {
	writer io.Writer
	color  color
}

func (cw colorWriter) Write(p []byte) (n int, err error) {
	_, err = cw.writer.Write([]byte("\033" + string(cw.color) + string(p) + "\033[0m"))
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

// newColorWriter creates a writer that colors the output with the given color
func newColorWriter(c color, w io.Writer) colorWriter {
	return colorWriter{
		writer: w,
		color:  c,
	}
}

// waitForRpcClient waits for an RPC to be available at the given url
// It does this by performing a GET request to the url until it receives a response
func waitForRpcClient(rpcClientUrl string, interval, timeout time.Duration) error {
	timeoutTicker := time.NewTicker(timeout)
	defer timeoutTicker.Stop()
	intervalTicker := time.NewTicker(interval)
	defer intervalTicker.Stop()

	client := &http.Client{}
	for {
		select {
		case <-timeoutTicker.C:
			return errors.New("polling timed out")
		case <-intervalTicker.C:
			resp, _ := client.Get(rpcClientUrl)
			if resp != nil {
				return nil
			}
		}
	}
}
