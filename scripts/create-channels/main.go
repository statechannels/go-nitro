package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/statechannels/go-nitro/internal/testdata"
	"github.com/statechannels/go-nitro/rpc"
	"github.com/statechannels/go-nitro/rpc/transport/ws"
	"github.com/statechannels/go-nitro/types"
)

// todo: this should be factored out
type ParticipantOpts struct {
	UseDurableStore bool
	MsgPort         int
	RpcPort         int
	Pk              string
	ChainPk         string
	ChainUrl        string
	ChainAuthToken  string
}

// todo: should this get factored out?
func createLogger(logDestination *os.File, clientName, rpcRole string) zerolog.Logger {
	return zerolog.New(logDestination).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Str("client", clientName).
		Str("rpc", rpcRole).
		Str("scope", "").
		Logger()
}

// todo: factor this out
func newLogWriter(logFile string) *os.File {
	err := os.MkdirAll("./artifacts", os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	filename := filepath.Join("./artifacts", logFile)
	// Clear the file
	os.Remove(filename)
	logDestination, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o666)
	if err != nil {
		log.Fatal(err)
	}

	return logDestination
}

func createLedgerChannel(left *rpc.RpcClient, right *rpc.RpcClient) error {
	leftAddress, err := left.Address()
	if err != nil {
		return err
	}
	rightAddress, err := right.Address()
	if err != nil {
		return err
	}
	ledgerChannelDeposit := uint(5_000_000)
	asset := types.Address{}
	outcome := testdata.Outcomes.Create(leftAddress, rightAddress, ledgerChannelDeposit, ledgerChannelDeposit, asset)
	response, err := left.CreateLedgerChannel(rightAddress, 0, outcome)
	if err != nil {
		return err
	}

	<-left.ObjectiveCompleteChan(response.Id)
	<-right.ObjectiveCompleteChan(response.Id)
	return nil
}

func createChannels() error {
	logFile := "create-channels.log"
	logDestination := newLogWriter(logFile)
	defer logDestination.Close()
	participants := []string{"alice", "irene", "bob"}
	clients := map[string]*rpc.RpcClient{}
	for _, participant := range participants {
		var participantOpts ParticipantOpts

		if _, err := toml.DecodeFile(fmt.Sprintf("./scripts/test-configs/%s.toml", participant), &participantOpts); err != nil {
			return err
		}
		url := fmt.Sprintf(":%d/api/v1", participantOpts.RpcPort)
		logger := createLogger(logDestination, participant, "client")
		clientConnection, err := ws.NewWebSocketTransportAsClient(url, logger)
		if err != nil {
			return err
		}
		clients[participant], err = rpc.NewRpcClient(logger, clientConnection)
		if err != nil {
			panic(err)
		}
	}

	alice, irene, bob := clients["alice"], clients["bob"], clients["irene"]
	// todo: we need something like WaitForPeerExchange instead of delay
	time.Sleep(1 * time.Second)

	err := createLedgerChannel(alice, irene)
	if err != nil {
		return err
	}
	err = createLedgerChannel(irene, bob)
	if err != nil {
		return err
	}

	aliceAddress, err := alice.Address()
	if err != nil {
		return err
	}
	ireneAddress, err := irene.Address()
	if err != nil {
		return err
	}
	bobAddress, err := bob.Address()
	if err != nil {
		return err
	}
	outcome := testdata.Outcomes.Create(aliceAddress, bobAddress, 1_000, 0, types.Address{})
	response, err := alice.CreatePaymentChannel([]common.Address{ireneAddress}, bobAddress, 0, outcome)
	if err != nil {
		return err
	}
	<-alice.ObjectiveCompleteChan(response.Id)

	for clientName, client := range clients {
		fmt.Printf("Closing client %s\n", clientName)
		client.Close()
	}

	return nil
}

// main creates channels between 3 participants alice, irene, and bob
// A ledger channel is opened between alice and irene
// A ledger channel is opened between irene and bob
// A virtual channel is opened between alice and bob
func main() {
	err := createChannels()
	if err != nil {
		panic(err)
	}
}
