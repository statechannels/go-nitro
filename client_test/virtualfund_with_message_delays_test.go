package client_test

import (
	"testing"
	"time"

	"github.com/statechannels/go-nitro/client"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	td "github.com/statechannels/go-nitro/internal/testdata"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

func TestVirtualFundWithMessageDelays(t *testing.T) {
	const MAX_MESSAGE_DELAY = time.Millisecond * 100

	// Since we are delaying messages we allow for enough time to complete the objective
	const OBJECTIVE_TIMEOUT = time.Second * 2

	// Setup logging
	logFile := "test_virtual_fund_with_message_delays.log"
	truncateLog(logFile)
	logDestination := newLogWriter(logFile)

	chain := chainservice.NewMockChain()
	chainServiceA := chainservice.NewMockChainService(chain, alice.Address())
	chainServiceB := chainservice.NewMockChainService(chain, bob.Address())
	chainServiceI := chainservice.NewMockChainService(chain, irene.Address())
	broker := messageservice.NewBroker()

	clientA, _ := setupClient(alice.PrivateKey, chainServiceA, broker, logDestination, MAX_MESSAGE_DELAY)
	defer closeClient(t, &clientA)
	clientB, _ := setupClient(bob.PrivateKey, chainServiceB, broker, logDestination, MAX_MESSAGE_DELAY)
	defer closeClient(t, &clientB)
	clientI, _ := setupClient(irene.PrivateKey, chainServiceI, broker, logDestination, MAX_MESSAGE_DELAY)
	defer closeClient(t, &clientI)

	directlyFundALedgerChannel(t, clientA, clientI, types.Address{})
	directlyFundALedgerChannel(t, clientI, clientB, types.Address{})

	ids := createVirtualChannels(clientA, bob.Address(), irene.Address(), 5)
	waitTimeForCompletedObjectiveIds(t, &clientA, OBJECTIVE_TIMEOUT, ids...)
	waitTimeForCompletedObjectiveIds(t, &clientB, OBJECTIVE_TIMEOUT, ids...)
	waitTimeForCompletedObjectiveIds(t, &clientI, OBJECTIVE_TIMEOUT, ids...)
}

// createVirtualChannels creates a number of virtual channels between the given parties and returns the objective ids.
//
//nolint:unused // unused due to skipped test
func createVirtualChannels(client client.Client, counterParty types.Address, intermediary types.Address, amountOfChannels uint) []protocols.ObjectiveId {
	ids := make([]protocols.ObjectiveId, amountOfChannels)
	for i := uint(0); i < amountOfChannels; i++ {
		outcome := td.Outcomes.Create(*client.Address, counterParty, 1, 1, types.Address{})
		ids[i] = client.CreateVirtualPaymentChannel([]types.Address{intermediary}, counterParty, 0, outcome).Id
	}
	return ids
}
