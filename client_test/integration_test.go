package client_test

import (
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/client"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	p2pms "github.com/statechannels/go-nitro/client/engine/messageservice/p2p-message-service"
	"github.com/statechannels/go-nitro/client/query"
	"github.com/statechannels/go-nitro/internal/testactors"
	td "github.com/statechannels/go-nitro/internal/testdata"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

func TestSimpleIntegrationScenario(t *testing.T) {
	simpleCase := TestCase{
		Description:    "Simple test",
		Chain:          MockChain,
		MessageService: TestMessageService,
		NumOfChannels:  1,
		MessageDelay:   0,
		LogName:        "simple_integration_run.log",
		NumOfHops:      1,
		NumOfPayments:  1,
		Participants: []TestParticipant{
			{StoreType: MemStore, Actor: testactors.Alice},
			{StoreType: MemStore, Actor: testactors.Bob},
			{StoreType: MemStore, Actor: testactors.Irene},
		},
	}

	RunIntegrationTestCase(simpleCase, t)
}

func TestComplexIntegrationScenario(t *testing.T) {
	complexCase := TestCase{
		Description:    "Complex test",
		Chain:          SimulatedChain,
		MessageService: P2PMessageService,
		NumOfChannels:  5,
		MessageDelay:   0,
		LogName:        "complex_integration_run.log",
		NumOfHops:      2,
		NumOfPayments:  5,
		Participants: []TestParticipant{
			{StoreType: DurableStore, Actor: testactors.Alice},
			{StoreType: DurableStore, Actor: testactors.Bob},
			{StoreType: DurableStore, Actor: testactors.Irene},
			{StoreType: DurableStore, Actor: testactors.Brian},
		},
	}
	RunIntegrationTestCase(complexCase, t)
}

// RunIntegrationTestCase runs the integration test case.
func RunIntegrationTestCase(tc TestCase, t *testing.T) {
	// Clean up all the test data we create at the end of the test
	defer os.RemoveAll(STORE_TEST_DATA_FOLDER)

	t.Run(tc.Description, func(t *testing.T) {
		err := tc.Validate()
		if err != nil {
			t.Fatal(err)
		}
		infra := setupSharedInfra(tc)
		defer infra.Close(t)

		msgServices := make([]messageservice.MessageService, 0)

		// Setup clients
		// NOTE: We rely on the convention that Alice is the first participant, Bob the second, and the intermediaries afterwards.
		clientA, msgA := setupIntegrationClient(tc, tc.Participants[0], infra)
		defer clientA.Close()
		msgServices = append(msgServices, msgA)

		clientB, msgB := setupIntegrationClient(tc, tc.Participants[1], infra)
		defer clientB.Close()
		msgServices = append(msgServices, msgB)

		intermediaries := make([]client.Client, 0)
		for _, intermediary := range tc.Participants[2:] {
			clientI, msgI := setupIntegrationClient(tc, intermediary, infra)

			intermediaries = append(intermediaries, clientI)
			msgServices = append(msgServices, msgI)
		}

		defer func() {
			for i := range intermediaries {
				intermediaries[i].Close()
			}
		}()

		if tc.MessageService == P2PMessageService {
			p2pServices := make([]*p2pms.P2PMessageService, len(tc.Participants))
			for i, msgService := range msgServices {
				p2pServices[i] = msgService.(*p2pms.P2PMessageService)
			}

			waitForPeerInfoExchange(p2pServices...)
		}

		asset := common.Address{}
		// Setup ledger channels between Alice/Bob and intermediaries
		aliceLedgers := make([]types.Destination, tc.NumOfHops)
		bobLedgers := make([]types.Destination, tc.NumOfHops)
		for i, clientI := range intermediaries {
			// Setup and check the ledger channel between Alice and the intermediary
			aliceLedgers[i] = setupLedgerChannel(t, clientA, clientI, asset)
			checkLedgerChannel(t, aliceLedgers[i], initialLedgerOutcome(*clientA.Address, *clientI.Address, asset), query.Open, clientA)
			// Setup and check the ledger channel between Bob and the intermediary
			bobLedgers[i] = setupLedgerChannel(t, clientI, clientB, asset)
			checkLedgerChannel(t, bobLedgers[i], initialLedgerOutcome(*clientI.Address, *clientB.Address, asset), query.Open, clientB)

		}

		if tc.NumOfHops == 2 {
			setupLedgerChannel(t, intermediaries[0], intermediaries[1], asset)
		}
		// Setup virtual channels
		objectiveIds := make([]protocols.ObjectiveId, tc.NumOfChannels)
		virtualIds := make([]types.Destination, tc.NumOfChannels)
		for i := 0; i < int(tc.NumOfChannels); i++ {
			outcome := td.Outcomes.Create(testactors.Alice.Address(), testactors.Bob.Address(), virtualChannelDeposit, 0, types.Address{})
			response, err := clientA.CreatePaymentChannel(
				clientAddresses(intermediaries),
				testactors.Bob.Address(),
				0,
				outcome,
			)
			if err != nil {
				t.Fatal(err)
			}
			objectiveIds[i] = response.Id
			virtualIds[i] = response.ChannelId

		}
		// Wait for all the virtual channels to be ready
		waitForObjectives(t, clientA, clientB, intermediaries, objectiveIds)

		// Check all the virtual channels
		for i := 0; i < len(virtualIds); i++ {
			checkPaymentChannel(t,
				virtualIds[i],
				initialPaymentOutcome(*clientA.Address, *clientB.Address, asset),
				query.Open,
				clientA, clientB)
		}

		// Send payments
		for i := 0; i < len(virtualIds); i++ {
			for j := 0; j < int(tc.NumOfPayments); j++ {
				clientA.Pay(virtualIds[i], big.NewInt(int64(1)))
			}
		}

		// Wait for all the vouchers to be received by bob
		for i := 0; i < len(virtualIds)*int(tc.NumOfPayments); i++ {
			<-clientB.ReceivedVouchers()
		}

		// Check the payment channels have the correct outcome after the payments
		for i := 0; i < len(virtualIds); i++ {
			checkPaymentChannel(t,
				virtualIds[i],
				finalPaymentOutcome(*clientA.Address, *clientB.Address, asset, tc.NumOfPayments, 1),
				query.Open,
				clientA, clientB)
		}

		// Close virtual channels
		closeVirtualIds := make([]protocols.ObjectiveId, len(virtualIds))
		for i := 0; i < len(virtualIds); i++ {
			// alternative who is responsible for closing the channel
			switch i % 2 {
			case 0:
				closeVirtualIds[i], err = clientA.ClosePaymentChannel(virtualIds[i])
				if err != nil {
					t.Fatal(err)
				}
			case 1:
				closeVirtualIds[i], err = clientB.ClosePaymentChannel(virtualIds[i])
				if err != nil {
					t.Fatal(err)
				}
			}
		}

		waitForObjectives(t, clientA, clientB, intermediaries, closeVirtualIds)

		// Close all the ledger channels we opened

		closeLedgerChannel(t, clientA, intermediaries[0], aliceLedgers[0])
		checkLedgerChannel(t, aliceLedgers[0], finalAliceLedger(*intermediaries[0].Address, asset, tc.NumOfPayments, 1, tc.NumOfChannels), query.Complete, clientA)

		// TODO: This is brittle, we should generalize this to n number of intermediaries
		if tc.NumOfHops == 1 {
			closeLedgerChannel(t, intermediaries[0], clientB, bobLedgers[0])
			checkLedgerChannel(t, bobLedgers[0], finalBobLedger(*intermediaries[0].Address, asset, tc.NumOfPayments, 1, tc.NumOfChannels), query.Complete, clientB)
		}
		if tc.NumOfHops == 2 {
			closeLedgerChannel(t, intermediaries[1], clientB, bobLedgers[1])
			checkLedgerChannel(t, bobLedgers[1], finalBobLedger(*intermediaries[1].Address, asset, tc.NumOfPayments, 1, tc.NumOfChannels), query.Complete, clientB)
		}
	})
}
