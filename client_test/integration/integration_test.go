// Package integration_test contains helpers and integration tests for go-nitro clients
package integration_test // import "github.com/statechannels/go-nitro/client_test/integration_test"

import (
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/client"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	p2pms "github.com/statechannels/go-nitro/client/engine/messageservice/p2p-message-service"
	"github.com/statechannels/go-nitro/internal/testactors"
	td "github.com/statechannels/go-nitro/internal/testdata"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

func setupSharedInra(tc TestRun) sharedInra {
	infra := sharedInra{}
	switch tc.Chain {
	case MockChain:
		infra.mockChain = chainservice.NewMockChain()
	case SimulatedChain:
		sim, bindings, ethAccounts, err := chainservice.SetupSimulatedBackend(3)
		if err != nil {
			panic(err)
		}
		infra.simulatedChain = &sim
		infra.bindings = &bindings
		infra.ethAccounts = ethAccounts
	default:
		panic("Unknown chain service")
	}

	switch tc.MessageService {
	case TestMessageService:
		broker := messageservice.NewBroker()
		infra.broker = &broker
	case P2PMessageService:

		infra.peers = make([]p2pms.PeerInfo, len(tc.Participants))
		for i, p := range tc.Participants {

			actor := getActor(p.Name)

			infra.peers[i] = p2pms.PeerInfo{
				Port:      int(actor.Port),
				IpAddress: "127.0.0.1",
				Address:   actor.Address(),
				Id:        p2pms.Id(actor.Address()),
			}
		}
	}

	return infra
}

func TestClientIntegration(t *testing.T) {

	// Clean up all the test data we create at the end of the test
	defer os.RemoveAll(STORE_TEST_DATA_FOLDER)

	for _, tc := range cases {
		t.Run(tc.Description, func(t *testing.T) {

			infra := setupSharedInra(tc)

			// Setup clients
			clientA := setupIntegrationClient(tc, testactors.AliceName, infra)
			clientB := setupIntegrationClient(tc, testactors.BobName, infra)

			intermediaries := []client.Client{setupIntegrationClient(tc, testactors.IreneName, infra)}
			if tc.NumOfHops == 2 {
				intermediaries = append(intermediaries, setupIntegrationClient(tc, testactors.BrianName, infra))
			}
			defer clientA.Close()
			defer clientB.Close()
			for _, clientI := range intermediaries {
				defer clientI.Close()
			}

			// Setup ledger channels between Alice/Bob and intermediaries
			aliceLedgers := make([]types.Destination, tc.NumOfHops)
			bobLedgers := make([]types.Destination, tc.NumOfHops)
			for i, clientI := range intermediaries {
				aliceLedgers[i] = setupLedgerChannel(t, clientA, clientI, common.Address{})
				bobLedgers[i] = setupLedgerChannel(t, clientB, clientI, common.Address{})

			}

			// Setup virtual channels
			objectiveIds := make([]protocols.ObjectiveId, tc.NumOfChannels)
			virtualIds := make([]types.Destination, tc.NumOfChannels)
			for i := 0; i < int(tc.NumOfChannels); i++ {
				outcome := td.Outcomes.Create(testactors.Alice.Address(), testactors.Bob.Address(), 1, 1, types.Address{})
				response := clientA.CreateVirtualPaymentChannel(
					[]types.Address{testactors.Irene.Address()},
					testactors.Bob.Address(),
					0,
					outcome,
				)
				objectiveIds[i] = response.Id
				virtualIds[i] = response.ChannelId
			}

			// TODO: This could probably be updated to use the new ObjectiveCompleted chan
			waitTimeForCompletedObjectiveIds(t, &clientA, defaultTimeout, objectiveIds...)
			waitTimeForCompletedObjectiveIds(t, &clientB, defaultTimeout, objectiveIds...)
			for _, clientI := range intermediaries {
				waitTimeForCompletedObjectiveIds(t, &clientI, defaultTimeout, objectiveIds...)
			}

			// Send payments
			for i := 0; i < len(virtualIds); i++ {
				for j := 0; j < int(tc.NumOfPayments); j++ {
					clientA.Pay(virtualIds[i], big.NewInt(int64(1)))
				}
			}

			// Close virtuall channels
			closeVirtualIds := make([]protocols.ObjectiveId, len(virtualIds))
			for i := 0; i < len(virtualIds); i++ {
				// alternative who is responsible for closing the channel
				switch i % (2 + int(tc.NumOfHops)) {
				case 0:
					closeVirtualIds[i] = clientA.CloseVirtualChannel(virtualIds[i])
				case 1:
					closeVirtualIds[i] = clientB.CloseVirtualChannel(virtualIds[i])
				case 2:
					closeVirtualIds[i] = intermediaries[0].CloseVirtualChannel(virtualIds[i])
				case 3:
					closeVirtualIds[i] = intermediaries[1].CloseVirtualChannel(virtualIds[i])
				}

			}
			// TODO: This could probably be updated to use the new ObjectiveCompleted chan
			waitTimeForCompletedObjectiveIds(t, &clientA, defaultTimeout, closeVirtualIds...)
			waitTimeForCompletedObjectiveIds(t, &clientB, defaultTimeout, closeVirtualIds...)
			for _, clientI := range intermediaries {
				waitTimeForCompletedObjectiveIds(t, &clientI, defaultTimeout, closeVirtualIds...)
			}

			// Close all the ledger channels we opened
			for i, l := range aliceLedgers {
				closeLedgerChannel(t, clientA, intermediaries[i], l)
			}
			for i, l := range bobLedgers {
				closeLedgerChannel(t, clientB, intermediaries[i], l)
			}

		})
	}
}
