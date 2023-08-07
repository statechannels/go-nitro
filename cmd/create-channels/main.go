package main

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/cmd/utils"
	"github.com/statechannels/go-nitro/internal/testdata"
	"github.com/statechannels/go-nitro/rpc"
	"github.com/statechannels/go-nitro/rpc/transport/ws"
	"github.com/statechannels/go-nitro/types"
)

type participantOpts struct {
	UseDurableStore bool
	MsgPort         int
	RpcPort         int
	Pk              string
	ChainPk         string
	ChainUrl        string
	ChainAuthToken  string
}

const LOG_FILE = "create-channels.log"

func createChannels() error {
	participants := []string{"alice", "irene", "bob"}
	clients := map[string]rpc.RpcClientApi{}
	for _, participant := range participants {
		var participantOpts participantOpts

		if _, err := toml.DecodeFile(fmt.Sprintf("./cmd/test-configs/%s.toml", participant), &participantOpts); err != nil {
			return err
		}

		logger, logFile := utils.CreateLogger(LOG_FILE, participant)
		defer logFile.Close()

		url := fmt.Sprintf(":%d/api/v1", participantOpts.RpcPort)
		clientConnection, err := ws.NewWebSocketTransportAsClient(url, logger)
		if err != nil {
			return err
		}
		clients[participant], err = rpc.NewRpcClient(logger, clientConnection)
		if err != nil {
			panic(err)
		}
	}

	alice, bob, irene := clients["alice"], clients["bob"], clients["irene"]

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

	err = utils.CreateLedgerChannel(alice, ireneAddress)
	if err != nil {
		return err
	}

	err = utils.CreateLedgerChannel(irene, bobAddress)
	if err != nil {
		return err
	}

	outcome := testdata.Outcomes.Create(aliceAddress, bobAddress, 1_000, 0, types.Address{})
	response, err := alice.CreatePaymentChannel([]common.Address{ireneAddress}, bobAddress, 0, outcome)
	if err != nil {
		return err
	}
	<-alice.ObjectiveCompleteChan(response.Id)

	for _, client := range clients {
		client.Close()
	}

	return nil
}

// main creates channels between 3 participants: alice, irene, and bob
// A ledger channel is opened between alice and irene
// A ledger channel is opened between irene and bob
// A virtual channel is opened between alice and bob
func main() {
	err := createChannels()
	if err != nil {
		panic(err)
	}
}
