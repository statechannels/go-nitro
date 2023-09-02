package main

import (
	"fmt"
	"log/slog"

	"github.com/BurntSushi/toml"
	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/cmd/utils"
	"github.com/statechannels/go-nitro/internal/logging"
	"github.com/statechannels/go-nitro/rpc"
	"github.com/statechannels/go-nitro/rpc/transport/http"
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

	logging.SetupDefaultFileLogger(LOG_FILE, slog.LevelDebug)

	url := fmt.Sprintf(":%d/api/v1", participantOpts.RpcPort)
	clientConnection, err := http.NewHttpTransportAsClient(url)
	if err != nil {
		panic(err)
	}
	client, err := rpc.NewRpcClient(clientConnection)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	err = utils.CreateLedgerChannel(client, common.HexToAddress(HUB_STATE_CHANNEL_ADDRESS))
	if err != nil {
		panic(err)
	}
}
