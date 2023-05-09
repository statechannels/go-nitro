package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/common"
	NitroAdjudicator "github.com/statechannels/go-nitro/client/engine/chainservice/adjudicator"
	chainutils "github.com/statechannels/go-nitro/client/engine/chainservice/utils"
	"github.com/statechannels/go-nitro/types"
)

type participant string

const (
	alice participant = "alice"
	bob   participant = "bob"
	irene participant = "irene"
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

const FUNDED_TEST_PK = "59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d"

func main() {
	running := []*exec.Cmd{}

	chainCmd := exec.Command("anvil", "--chain-id", "1337")
	chainCmd.Stdout = os.Stdout
	chainCmd.Stderr = os.Stderr
	err := chainCmd.Start()
	if err != nil {
		stopCommands(running...)
		panic(err)
	}
	running = append(running, chainCmd)
	// Give the chain a second to start up
	time.Sleep(1 * time.Second)

	adjudicatorAddress, err := deployAdjudicator(context.Background())
	if err != nil {
		stopCommands(running...)
		panic(err)
	}
	fmt.Printf("Deployed adjudicator at %v\n", adjudicatorAddress)

	aliceClient, err := setupRPCServer(alice, blue, adjudicatorAddress)
	if err != nil {
		stopCommands(running...)
		panic(err)
	}
	running = append(running, aliceClient)

	ireneClient, err := setupRPCServer(irene, green, adjudicatorAddress)
	if err != nil {
		stopCommands(running...)
		panic(err)
	}
	running = append(running, ireneClient)

	bobClient, err := setupRPCServer(bob, yellow, adjudicatorAddress)
	if err != nil {
		stopCommands(running...)
		panic(err)
	}
	running = append(running, bobClient)

	waitForKillSignal()

	stopCommands(running...)
}

// stopCommands stops the given executing commands
func stopCommands(cmds ...*exec.Cmd) {
	for _, cmd := range cmds {
		fmt.Printf("Stopping process %v\n", cmd.Args)
		err := cmd.Process.Signal(syscall.SIGINT)
		if err != nil {
			panic(err)
		}
		err = cmd.Process.Kill()
		if err != nil {
			panic(err)
		}
	}
}

// waitForKillSignal blocks until we receive a kill or interrupt signal
func waitForKillSignal() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigs
	fmt.Printf("Received signal %s, exiting..\n", sig)
}

// setupRPCServer starts up an RPC server for the given participant
func setupRPCServer(p participant, c color, na types.Address) (*exec.Cmd, error) {
	args := []string{"run", ".", "-usedurablestore"}

	switch p {
	case alice:
		args = append(args, "-msgport", "3005")
		args = append(args, "-rpcport", "4005")
		args = append(args, "-pk", "2d999770f7b5d49b694080f987b82bbc9fc9ac2b4dcc10b0f8aba7d700f69c6d")

	case irene:
		args = append(args, "-msgport", "3006")
		args = append(args, "-rpcport", "4006")
		args = append(args, "-pk", "febb3b74b0b52d0976f6571d555f4ac8b91c308dfa25c7b58d1e6a7c3f50c781")

	case bob:
		args = append(args, "-msgport", "3007")
		args = append(args, "-rpcport", "4007")
		args = append(args, "-pk", "0279651921cd800ac560c21ceea27aab0107b67daf436cdd25ce84cad30159b4")

	default:
		panic("Invalid participant")

	}
	args = append(args, "-naaddress", na.String())
	cmd := exec.Command("go", args...)
	cmd.Stdout = newColorWriter(c, os.Stdout)
	cmd.Stderr = os.Stderr
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

func newColorWriter(c color, w io.Writer) colorWriter {
	return colorWriter{
		writer: w,
		color:  c,
	}
}

// deployAdjudicator deploys th  NitroAdjudicator contract.

func deployAdjudicator(ctx context.Context) (common.Address, error) {
	client, txSubmitter, err := chainutils.ConnectToChain(context.Background(), "ws://127.0.0.1:8545", 1337, common.Hex2Bytes(FUNDED_TEST_PK))
	if err != nil {
		return types.Address{}, err
	}
	naAddress, _, _, err := NitroAdjudicator.DeployNitroAdjudicator(txSubmitter, client)
	if err != nil {
		return types.Address{}, err
	}
	return naAddress, nil
}
