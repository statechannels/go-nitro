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

	chainUrl := "ws://127.0.0.1:8545"
	chainPk := "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	_, _, _, err = chain.DeployContracts(context.Background(), chainUrl, "", chainPk)
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
