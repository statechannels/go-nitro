package client_test

import (
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/statechannels/go-nitro/client"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	"github.com/statechannels/go-nitro/protocols/virtualfund"
	"github.com/statechannels/go-nitro/types"
)

// TestBenchmark sets up three clients, then runs a virtual funding benchmark, printing the duration
// to the screen.
func TestBenchmark(t *testing.T) {

	logFile := "virtualfund_benchmark_test.log"
	truncateLog(logFile)

	chain := chainservice.NewMockChain()
	broker := messageservice.NewBroker()

	clientAlice := setupClient(aliceKey, chain, broker, logFile)
	clientBob := setupClient(bobKey, chain, broker, logFile)
	clientIrene := setupClient(ireneKey, chain, broker, logFile)

	directlyFundALedgerChannel(t, clientAlice, clientIrene)
	directlyFundALedgerChannel(t, clientIrene, clientBob)

	done := make(chan interface{})

	n := 3
	for i := 0; i < n; i++ {
		go benchmarkVirtualChannelCreation(t, clientAlice, clientBob, irene, done)
	}

	expect(t, done, n,  time.Second*1)
}

func benchmarkVirtualChannelCreation(t *testing.T, alice, bob client.Client, irene types.Address, done chan interface{}) {
	outcome := createVirtualOutcome(*alice.Address, *bob.Address)
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
	id := alice.CreateVirtualChannel(request)

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
			t.Errorf("Ran out of time. %v out of %v completed", count, num)
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
