package main

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/cmd/utils"
	"github.com/statechannels/go-nitro/rpc"
	"github.com/statechannels/go-nitro/rpc/transport/ws"
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
	CONFIG_FILE               = "./docker/local/config.toml"
	LOG_FILE                  = "create-single-channel.log"
)

func main() {
	var participantOpts participantOpts
	if _, err := toml.DecodeFile(CONFIG_FILE, &participantOpts); err != nil {
		panic(err)
	}
	logger, logFile := utils.CreateLogger(LOG_FILE, "alice")
	defer logFile.Close()

	url := fmt.Sprintf(":%d/api/v1", participantOpts.RpcPort)
	clientConnection, err := ws.NewWebSocketTransportAsClient(url, logger)
	if err != nil {
		panic(err)
	}
	client, err := rpc.NewRpcClient(logger, clientConnection)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	err = utils.CreateLedgerChannel(client, common.HexToAddress(HUB_STATE_CHANNEL_ADDRESS))
	if err != nil {
		panic(err)
	}
}
