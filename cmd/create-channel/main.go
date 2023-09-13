package main

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/cmd/utils"
	"github.com/statechannels/go-nitro/internal/logging"
	"github.com/statechannels/go-nitro/rpc"
	"github.com/statechannels/go-nitro/rpc/transport/http"
)

const (
	HUB_STATE_CHANNEL_ADDRESS = "0x89825D58A7E2C198a06125be9CD0631317f9A07B"
	NODE_IP_ADDRESS           = "192.81.214.172"
	NODE_RPC_PORT             = 4005
	CONFIG_FILE               = "./docker/local/config.toml"
	LOG_FILE                  = "create-single-channel.log"
)

func main() {
	logging.SetupDefaultFileLogger(LOG_FILE, slog.LevelDebug)

	url := fmt.Sprintf("%s:%d/api/v1", NODE_IP_ADDRESS, NODE_RPC_PORT)
	clientConnection, err := http.NewHttpTransportAsClient(url, 10*time.Millisecond)
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
