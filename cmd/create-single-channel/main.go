package main

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/statechannels/go-nitro/internal/logging"
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

const (
	HUB_STATE_CHANNEL_ADDRESS = "0x69Ae2DF44965e6f87b25bc648E097B00762583Dc"
)

func createLogger(logDestination *os.File) zerolog.Logger {
	return zerolog.New(logDestination).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Str("rpc", "client").
		Logger()
}

func createLedgerChannel(client rpc.RpcClientApi, counterPartyAddress common.Address) error {
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

func createChannels() error {
	logFile := "create-channels.log"
	logDestination := logging.NewLogWriter("./artifacts", logFile)
	defer logDestination.Close()

	var participantOpts participantOpts
	if _, err := toml.DecodeFile("./docker/local/config.toml", &participantOpts); err != nil {
		return err
	}
	url := fmt.Sprintf(":%d/api/v1", participantOpts.RpcPort)
	logger := createLogger(logDestination)
	clientConnection, err := ws.NewWebSocketTransportAsClient(url, logger)
	if err != nil {
		return err
	}
	client, err := rpc.NewRpcClient(logger, clientConnection)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	err = createLedgerChannel(client, common.HexToAddress(HUB_STATE_CHANNEL_ADDRESS))
	if err != nil {
		return err
	}

	return nil
}

func main() {
	err := createChannels()
	if err != nil {
		panic(err)
	}
}
