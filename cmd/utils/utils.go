package utils

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/statechannels/go-nitro/internal/testdata"
	"github.com/statechannels/go-nitro/rpc"
	"github.com/statechannels/go-nitro/types"
)

// WaitForKillSignal blocks until we receive a kill or interrupt signal
func WaitForKillSignal() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigs
	fmt.Printf("Received signal %s, exiting..\n", sig)
}

// StopCommands stops the given executing commands
func StopCommands(cmds ...*exec.Cmd) {
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

func CreateLedgerChannel(client rpc.RpcClientApi, counterPartyAddress common.Address) error {
	clientAddress, err := client.Address()
	if err != nil {
		return err
	}
	ledgerChannelDeposit := uint(5_000_000)
	asset := types.Address{}
	outcome := testdata.Outcomes.Create(clientAddress, counterPartyAddress, ledgerChannelDeposit, ledgerChannelDeposit, asset)
	response, err := client.CreateLedgerChannel(counterPartyAddress, 0, outcome)
	if err != nil {
		return err
	}

	<-client.ObjectiveCompleteChan(response.Id)
	return nil
}

// waitForRpcClient waits for an RPC to be available at the given url
// It does this by performing a GET request to the url until it receives a response
func WaitForRpcClient(rpcClientUrl string, interval, timeout time.Duration) error {
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
