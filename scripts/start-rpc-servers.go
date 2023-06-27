package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/statechannels/go-nitro/internal/chain"
	"github.com/statechannels/go-nitro/internal/utils"
	"github.com/statechannels/go-nitro/types"
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

var (
	participants     = []participant{alice, bob, irene, ivan}
	participantColor = map[participant]color{alice: blue, irene: green, ivan: cyan, bob: yellow}
)

func main() {
	running := []*exec.Cmd{}

	anvilCmd, err := chain.StartAnvil()
	if err != nil {
		utils.StopCommands(running...)
		panic(err)
	}
	running = append(running, anvilCmd)

	naAddress, vpaAddress, caAddress, err := chain.DeployContracts(context.Background())
	if err != nil {
		utils.StopCommands(running...)
		panic(err)
	}

	for _, p := range participants {
		client, err := setupRPCServer(p, participantColor[p], naAddress, vpaAddress, caAddress)
		if err != nil {
			utils.StopCommands(running...)
			panic(err)
		}
		running = append(running, client)
	}

	waitForKillSignal()

	utils.StopCommands(running...)
}

// waitForKillSignal blocks until we receive a kill or interrupt signal
func waitForKillSignal() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigs
	fmt.Printf("Received signal %s, exiting..\n", sig)
}

// setupRPCServer starts up an RPC server for the given participant
func setupRPCServer(p participant, c color, na, vpa, ca types.Address) (*exec.Cmd, error) {
	args := []string{"run", ".", "-naaddress", na.String()}
	args = append(args, "-vpaaddress", vpa.String())
	args = append(args, "-caaddress", ca.String())
	args = append(args, "-config", fmt.Sprintf("./scripts/test-configs/%s.toml", p))

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
