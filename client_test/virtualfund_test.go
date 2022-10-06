package client_test

import (
	"testing"

	"github.com/statechannels/go-nitro/client"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	td "github.com/statechannels/go-nitro/internal/testdata"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

func openVirtualChannels(t *testing.T, clientA client.Client, clientB client.Client, clientI client.Client, numOfChannels uint) []types.Destination {
	directlyFundALedgerChannel(t, clientA, clientI)
	directlyFundALedgerChannel(t, clientI, clientB)

	objectiveIds := make([]protocols.ObjectiveId, numOfChannels)
	channelIds := make([]types.Destination, numOfChannels)
	for i := 0; i < int(numOfChannels); i++ {
		outcome := td.Outcomes.Create(alice.Address(), bob.Address(), 1, 1)
		response := clientA.CreateVirtualPaymentChannel(
			[]types.Address{irene.Address()},
			bob.Address(),
			0,
			outcome,
		)

		objectiveIds[i] = response.Id
		channelIds[i] = response.ChannelId
	}
	waitTimeForCompletedObjectiveIds(t, &clientA, defaultTimeout, objectiveIds...)
	waitTimeForCompletedObjectiveIds(t, &clientB, defaultTimeout, objectiveIds...)
	waitTimeForCompletedObjectiveIds(t, &clientI, defaultTimeout, objectiveIds...)

	return channelIds

}
func TestVirtualFundIntegration(t *testing.T) {

	// Setup logging
	logFile := "test_virtual_fund.log"
	truncateLog(logFile)
	logDestination := newLogWriter(logFile)

	chain := chainservice.NewMockChain()
	chainServiceA := chainservice.NewMockChainService(chain, alice.Address())
	chainServiceB := chainservice.NewMockChainService(chain, bob.Address())
	chainServiceI := chainservice.NewMockChainService(chain, irene.Address())
	broker := messageservice.NewBroker()

	clientA, _ := setupClient(alice.PrivateKey, chainServiceA, broker, logDestination, 0)
	clientB, _ := setupClient(bob.PrivateKey, chainServiceB, broker, logDestination, 0)
	clientI, _ := setupClient(irene.PrivateKey, chainServiceI, broker, logDestination, 0)

	openVirtualChannels(t, clientA, clientB, clientI, 1)
}
