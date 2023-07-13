package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/statechannels/go-nitro/internal/chain"
	"github.com/statechannels/go-nitro/internal/utils"
)

func main() {
	running := []*exec.Cmd{}

	anvilCmd, err := chain.StartAnvil()
	if err != nil {
		utils.StopCommands(running...)
		panic(err)
	}
	running = append(running, anvilCmd)

	_, _, _, err = chain.DeployContracts(context.Background())
	if err != nil {
		utils.StopCommands(running...)
		panic(err)
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
