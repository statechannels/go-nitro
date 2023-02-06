package client_test

import (
	"math/rand"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/statechannels/go-nitro/internal/testdata"
	"github.com/statechannels/go-nitro/network"
	"github.com/statechannels/go-nitro/network/transport"
	natstrans "github.com/statechannels/go-nitro/network/transport/nats"
	"github.com/statechannels/go-nitro/protocols/directfund"
	"github.com/statechannels/go-nitro/types"
)

func newConnection() (transport.Connection, error) {
	opts := &server.Options{}
	ns, err := server.NewServer(opts)
	if err != nil {
		return nil, err
	}
	ns.Start()

	nc, err := nats.Connect(ns.ClientURL())
	if err != nil {
		return nil, err
	}

	con := natstrans.NewNatsConnection(nc)
	return con, nil
}

func TestNetworkClient(t *testing.T) {
	testOutcome := testdata.Outcomes.Create(alice.Address(), bob.Address(), 100, 100, types.Address{})

	objReq := directfund.NewObjectiveRequest(
		bob.Address(), 100,
		testOutcome,
		uint64(rand.Float64()),
		common.Address{})

	connection, err := newConnection()
	if err != nil {
		t.Fatal(err)
	}

	logger := createLogger(newLogWriter("test_network_client.log"), "alice", "client")
	clientConnetion := network.ClientConnection{Connection: connection}

	_, err = network.Request(&clientConnetion, objReq, logger)
	if err != nil {
		t.Fatal(err)
	}
}
