package rpc

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/statechannels/go-nitro/node"
	"github.com/statechannels/go-nitro/rpc"
	"github.com/statechannels/go-nitro/rpc/transport"
	httpTransport "github.com/statechannels/go-nitro/rpc/transport/http"
	"github.com/statechannels/go-nitro/rpc/transport/nats"
)

func InitializeRpcServer(node *node.Node, rpcPort int, useNats bool) (*rpc.RpcServer, error) {
	var transport transport.Responder
	var err error

	if useNats {
		slog.Info("Initializing NATS RPC transport...")
		transport, err = nats.NewNatsTransportAsServer(rpcPort)
	} else {
		slog.Info("Initializing Http RPC transport...")
		transport, err = httpTransport.NewHttpTransportAsServer(fmt.Sprint(rpcPort))
		if err != nil {
			return nil, err
		}
		err = blockUntilHttpServerIsReady(rpcPort)
	}
	if err != nil {
		return nil, err
	}

	rpcServer, err := rpc.NewRpcServer(node, transport)
	if err != nil {
		return nil, err
	}

	slog.Info("Completed RPC server initialization")
	return rpcServer, nil
}

// blockUntilHttpServerIsReady pings the health endpoint until the server is ready
func blockUntilHttpServerIsReady(rpcPort int) error {
	waitForServer := func() {
		time.Sleep(10 * time.Millisecond)
	}

	numAttempts := 10
	for i := 0; i < numAttempts; i++ {
		resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/health", rpcPort))
		if err != nil {
			waitForServer()
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			return nil
		}
		waitForServer()
	}
	return fmt.Errorf("http server not ready after %d attempts", numAttempts)
}
