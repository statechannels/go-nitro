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
	directlyFundALedgerChannel(t, clientA, clientI, types.Address{})
	directlyFundALedgerChannel(t, clientI, clientB, types.Address{})

	objectiveIds := make([]protocols.ObjectiveId, numOfChannels)
	channelIds := make([]types.Destination, numOfChannels)
	for i := 0; i < int(numOfChannels); i++ {
		outcome := td.Outcomes.Create(alice.Address(), bob.Address(), 1, 1, types.Address{})
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
	chainServiceBr := chainservice.NewMockChainService(chain, brian.Address())
	broker := messageservice.NewBroker()

	clientA, _ := setupClient(alice.PrivateKey, chainServiceA, broker, logDestination, 0)
	defer closeClient(t, &clientA)
	irene, _ := setupClient(irene.PrivateKey, chainServiceI, broker, logDestination, 0)
	defer closeClient(t, &irene)
	ivan, _ := setupClient(brian.PrivateKey, chainServiceBr, broker, logDestination, 0)
	defer closeClient(t, &ivan)
	clientB, _ := setupClient(bob.PrivateKey, chainServiceB, broker, logDestination, 0)
	defer closeClient(t, &clientB)

	openN_HopVirtualChannels(t, []client.Client{clientA, irene, ivan, clientB}, 1)
}

// openN_HopVirtualChannels connects the n given participants in a line of ledger channels,
// then uses these ledger connections to open channels between the first participant
// and each other participant.
//
// This makes channels with 0, 1, ..., [len(participants) - 1] hops.
func openN_HopVirtualChannels(t *testing.T, participants []client.Client, channelsPerCounterparty uint) {

	// set a chain of ledger channels between incoming clients.
	// network is a line: A <-> B <-> C <-> D ... <-> X
	for i, participant := range participants {
		if i+1 < len(participants) {
			directlyFundALedgerChannel(t, participant, participants[i+1], types.Address{})
		}
	}

	alice := participants[0] // alice initiates each virtualfund operation
	counterparties := participants[1:]

	for i := 0; i < int(channelsPerCounterparty); i++ {

		// and funds channels with everyone else along the chain
		// (including her immediate neighbor)
		for j, bob := range counterparties {

			intermediaries := counterparties[0:j]
			intermediaryAddresses := clientsToAddresses(intermediaries)

			outcome := td.Outcomes.Create(*alice.Address, *bob.Address, 1, 1, types.Address{})
			response := alice.CreateVirtualPaymentChannel(
				intermediaryAddresses,
				*bob.Address,
				0,
				outcome,
			)

			waitTimeForCompletedObjectiveIds(t, &alice, defaultTimeout, response.Id)
			for _, intermediary := range intermediaries {
				waitTimeForCompletedObjectiveIds(t, &intermediary, defaultTimeout, response.Id)
			}
			waitTimeForCompletedObjectiveIds(t, &bob, defaultTimeout, response.Id)
		}
	}

}

func clientsToAddresses(clients []client.Client) []types.Address {
	ret := []types.Address{}
	for _, client := range clients {
		ret = append(ret, *client.Address)
	}
	return ret
}
