package client_test

import (
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/statechannels/go-nitro/client"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	td "github.com/statechannels/go-nitro/internal/testdata"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/protocols/virtualfund"
	"github.com/statechannels/go-nitro/types"
)

const MAX_MESSAGE_DELAY = time.Millisecond * 100

// Since we are delaying messages we allow for slightly more time to complete the objective
const OBJECTIVE_TIMEOUT = time.Second * 2

func TestVirtualFundWithMessageDelays(t *testing.T) {
	// This test fails due to https://github.com/statechannels/go-nitro/issues/366
	t.Skip()
	// Set up logging
	logFile := "virtual_fund_message_delay_test.log"
	truncateLog(logFile)

	chain := chainservice.NewMockChain()
	broker := messageservice.NewBroker()

	clientA := setupClient(alice.PrivateKey, chain, broker, logFile, MAX_MESSAGE_DELAY)
	clientB := setupClient(bob.PrivateKey, chain, broker, logFile, MAX_MESSAGE_DELAY)
	clientI := setupClient(irene.PrivateKey, chain, broker, logFile, MAX_MESSAGE_DELAY)

	directlyFundALedgerChannel(t, clientA, clientI)
	directlyFundALedgerChannel(t, clientI, clientB)

	ids := createVirtualChannels(clientA, bob.Address, irene.Address, 5)
	waitTimeForCompletedObjectiveIds(t, &clientA, OBJECTIVE_TIMEOUT, ids...)
	waitTimeForCompletedObjectiveIds(t, &clientB, OBJECTIVE_TIMEOUT, ids...)
	waitTimeForCompletedObjectiveIds(t, &clientI, OBJECTIVE_TIMEOUT, ids...)

}

// createVirtualChannels creates a number of virtual channels between the given parties and returns the objective ids.
//nolint:unused // unused due to skipped test
func createVirtualChannels(client client.Client, counterParty types.Address, intermediary types.Address, amountOfChannels uint) []protocols.ObjectiveId {
	ids := make([]protocols.ObjectiveId, amountOfChannels)
	for i := uint(0); i < amountOfChannels; i++ {
		outcome := td.Outcomes.Create(*client.Address, counterParty, 1, 1)
		request := virtualfund.ObjectiveRequest{
			MyAddress:         *client.Address,
			CounterParty:      counterParty,
			Intermediary:      intermediary,
			Outcome:           outcome,
			AppDefinition:     types.Address{},
			AppData:           types.Bytes{},
			ChallengeDuration: big.NewInt(0),
			Nonce:             rand.Int63(),
		}

		ids[i] = client.CreateVirtualChannel(request)
	}
	return ids
}
