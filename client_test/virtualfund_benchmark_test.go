package client_test

import (
	"bytes"
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/statechannels/go-nitro/client"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	"github.com/statechannels/go-nitro/internal/testdata"
	"github.com/statechannels/go-nitro/protocols/virtualfund"
	"github.com/statechannels/go-nitro/types"
)

// TestBenchmark sets up three clients, then runs a virtual funding benchmark, printing the duration
// to the screen.
func TestBenchmark(t *testing.T) {

	// Setup logging
	logDestination := &bytes.Buffer{}
	t.Cleanup(flushToFileCleanupFn(logDestination, "virtualfund_benchmark_test.log"))

	chain := chainservice.NewMockChain()
	broker := messageservice.NewBroker()

	clientAlice, _ := setupClient(alice.PrivateKey, chain, broker, logDestination, 0)
	clientBob, _ := setupClient(bob.PrivateKey, chain, broker, logDestination, 0)
	clientIrene, _ := setupClient(irene.PrivateKey, chain, broker, logDestination, 0)

	directlyFundALedgerChannel(t, clientAlice, clientIrene)
	directlyFundALedgerChannel(t, clientIrene, clientBob)

	done := make(chan interface{})

	n := 1
	for i := 0; i < n; i++ {
		go benchmarkVirtualChannelCreation(t, clientAlice, clientBob, irene.Address, done)
	}

	expect(t, done, n, time.Second*1)
}

// benchmarkVirtualChannelCreation creates a new virtual channel with the given actors, and
// times how long it takes for the objective to complete (from Bob's point of view)
// The resulting time is printed to the test runner's output
func benchmarkVirtualChannelCreation(t *testing.T, alice, bob client.Client, irene types.Address, done chan interface{}) {
	outcome := testdata.Outcomes.Create(*alice.Address, *bob.Address, 1, 1)
	request := virtualfund.ObjectiveRequest{
		MyAddress:         *alice.Address,
		CounterParty:      *bob.Address,
		Intermediary:      irene,
		Outcome:           outcome,
		AppDefinition:     types.Address{},
		AppData:           types.Bytes{},
		ChallengeDuration: big.NewInt(0),
		Nonce:             rand.Int63(),
	}
	id := alice.CreateVirtualChannel(request).Id

	defer elapsed(t, string(id))()

	for got := range bob.CompletedObjectives() {
		if got == id {
			done <- nil
			return
		}
	}
}

// Returns after `done` has received `num` messages.
//
// To ensure it eventually returns, it will error after a timeout, which resets
// whenever `done` receives a message. So, it will return after at most
// `num * defaultTimeout` time has elapsed.
func expect(t *testing.T, done chan interface{}, num int, timeout time.Duration) {
	count := 0
	for {
		select {
		case <-done:
			count += 1
			if count == num {
				return
			}
		case <-time.After(timeout):
			t.Fatalf("Ran out of time. %v out of %v completed", count, num)
			t.FailNow()
		}
	}
}

func elapsed(t *testing.T, what string) func() {
	start := time.Now()
	return func() {
		t.Logf("%s took %v\n", what, time.Since(start))
	}
}
