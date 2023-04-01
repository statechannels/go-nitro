// Package integration_test contains helpers and integration tests for go-nitro clients
package integration_test // import "github.com/statechannels/go-nitro/client_test/integration_test"

import (
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/client"
	"github.com/statechannels/go-nitro/client/query"
	"github.com/statechannels/go-nitro/internal/testactors"
	td "github.com/statechannels/go-nitro/internal/testdata"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

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
			intermediaryAddresses := []types.Address{*intermediaries[0].Address}
			if tc.NumOfHops == 2 {
				intermediaries = append(intermediaries, setupIntegrationClient(tc, testactors.BrianName, infra))
				intermediaryAddresses = append(intermediaryAddresses, *intermediaries[1].Address)
			}
			defer clientA.Close()
			defer clientB.Close()

			asset := common.Address{}
			// Setup ledger channels between Alice/Bob and intermediaries
			aliceLedgers := make([]types.Destination, tc.NumOfHops)
			bobLedgers := make([]types.Destination, tc.NumOfHops)
			for i, clientI := range intermediaries {
				// Setup and check the ledger channel between Alice and the intermediary
				aliceLedgers[i] = setupLedgerChannel(t, clientA, clientI, asset)
				checkLedgerChannel(t, aliceLedgers[i], initialLedgerOutcome(*clientA.Address, *clientI.Address, asset), query.Ready, &clientA)
				// Setup and check the ledger channel between Bob and the intermediary
				bobLedgers[i] = setupLedgerChannel(t, clientI, clientB, asset)
				checkLedgerChannel(t, bobLedgers[i], initialLedgerOutcome(*clientI.Address, *clientB.Address, asset), query.Ready, &clientB)

			}

			if tc.NumOfHops == 2 {
				setupLedgerChannel(t, intermediaries[0], intermediaries[1], asset)
			}
			// Setup virtual channels
			objectiveIds := make([]protocols.ObjectiveId, tc.NumOfChannels)
			virtualIds := make([]types.Destination, tc.NumOfChannels)
			for i := 0; i < int(tc.NumOfChannels); i++ {
				outcome := td.Outcomes.Create(testactors.Alice.Address(), testactors.Bob.Address(), virtualChannelDeposit, 0, types.Address{})
				response := clientA.CreateVirtualPaymentChannel(
					intermediaryAddresses,
					testactors.Bob.Address(),
					0,
					outcome,
				)
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
					query.Ready,
					&clientA, &clientB)
			}

			// Send payments
			for i := 0; i < len(virtualIds); i++ {
				for j := 0; j < int(tc.NumOfPayments); j++ {
					clientA.Pay(virtualIds[i], big.NewInt(int64(1)))
				}
			}

			// Wait for all the vouchers to be received by bob
			for i := 0; i < len(virtualIds); i++ {
				<-clientB.ReceivedVouchers()
			}

			// Check the payment channels have the correct outcome after the payments
			for i := 0; i < len(virtualIds); i++ {
				checkPaymentChannel(t,
					virtualIds[i],
					finalPaymentOutcome(*clientA.Address, *clientB.Address, asset, tc.NumOfPayments, 1),
					query.Ready,
					&clientA, &clientB)
			}

			// Close virtual channels
			closeVirtualIds := make([]protocols.ObjectiveId, len(virtualIds))
			for i := 0; i < len(virtualIds); i++ {
				// alternative who is responsible for closing the channel
				switch i % 2 {
				case 0:
					closeVirtualIds[i] = clientA.CloseVirtualChannel(virtualIds[i])
				case 1:
					closeVirtualIds[i] = clientB.CloseVirtualChannel(virtualIds[i])
				}
			}

			waitForObjectives(t, clientA, clientB, intermediaries, closeVirtualIds)

			// Close all the ledger channels we opened

			closeLedgerChannel(t, clientA, intermediaries[0], aliceLedgers[0])
			checkLedgerChannel(t, aliceLedgers[0], finalAliceLedger(*intermediaries[0].Address, asset, tc.NumOfPayments, 1, tc.NumOfChannels), query.Complete, &clientA)

			// TODO: This is brittle, we should generalize this to
			if tc.NumOfHops == 1 {
				closeLedgerChannel(t, intermediaries[0], clientB, bobLedgers[0])
				checkLedgerChannel(t, bobLedgers[0], finalBobLedger(*intermediaries[0].Address, asset, tc.NumOfPayments, 1, tc.NumOfChannels), query.Complete, &clientB)
			}
			if tc.NumOfHops == 2 {
				closeLedgerChannel(t, intermediaries[1], clientB, bobLedgers[1])
				checkLedgerChannel(t, bobLedgers[1], finalBobLedger(*intermediaries[1].Address, asset, tc.NumOfPayments, 1, tc.NumOfChannels), query.Complete, &clientB)
			}
			for _, clientI := range intermediaries {
				clientI.Close()
			}
		})
	}
}
