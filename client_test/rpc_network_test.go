package client_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/statechannels/go-nitro/client/engine/store/safesync"
	"github.com/statechannels/go-nitro/internal/testdata"
	"github.com/statechannels/go-nitro/network"
	"github.com/statechannels/go-nitro/network/serde"
	"github.com/statechannels/go-nitro/network/transport"
	natstrans "github.com/statechannels/go-nitro/network/transport/nats"
	"github.com/statechannels/go-nitro/protocols/directfund"
	"github.com/statechannels/go-nitro/types"
)

// getTopics returns a list of topics that the client/server should subscribe to.
func getTopics() []string {

	return []string{
		fmt.Sprintf("nitro.%s",
			serde.DirectFundRequestMethod),
		fmt.Sprintf("nitro.%s",
			serde.DirectDefundRequestMethod),
		fmt.Sprintf("nitro.%s",
			serde.VirtualFundRequestMethod),
		fmt.Sprintf("nitro.%s",
			serde.VirtualDefundRequestMethod),
		fmt.Sprintf("nitro.%s",
			serde.PayRequestMethod),
	}
}

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

	trp := natstrans.NewNatsTransport(nc, getTopics())

	con, err := trp.PollConnection()
	if err != nil {
		return nil, err
	}
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

	idsToMethods := safesync.Map[serde.RequestMethod]{}

	_, err = network.Request(&clientConnetion, objReq, logger, &idsToMethods)
	if err != nil {
		t.Fatal(err)
	}
}
