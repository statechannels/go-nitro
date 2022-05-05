// Package client_test contains helpers and integration tests for go-nitro clients
package client_test // import "github.com/statechannels/go-nitro/client_test"

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	"github.com/statechannels/go-nitro/client/engine/store"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/protocols/virtualdefund"
	"github.com/statechannels/go-nitro/types"
)

func TestVirtualDefundIntegration(t *testing.T) {

	// Setup logging
	logFile := "test_virtual_defund.log"
	truncateLog(logFile)
	logDestination := newLogWriter(logFile)

	chain := chainservice.NewMockChain()
	broker := messageservice.NewBroker()

	clientA, storeA := setupClient(alice.PrivateKey, chain, broker, logDestination, 0)
	clientB, storeB := setupClient(bob.PrivateKey, chain, broker, logDestination, 0)
	clientI, storeI := setupClient(irene.PrivateKey, chain, broker, logDestination, 0)

	numOfVirtualChannels := uint(5)
	paidToBob := uint(1)
	totalPaidToBob := paidToBob * numOfVirtualChannels

	// TODO: This test only supports defunding 1 virtual channel due to https://github.com/statechannels/go-nitro/issues/637
	cIds := openVirtualChannels(t, clientA, clientB, clientI, numOfVirtualChannels)

	ids := make([]protocols.ObjectiveId, len(cIds))
	for i := 0; i < len(cIds); i++ {
		ids[i] = clientA.CloseVirtualChannel(cIds[i], big.NewInt(int64(paidToBob)))

	}
	waitTimeForCompletedObjectiveIds(t, &clientA, defaultTimeout, ids...)
	waitTimeForCompletedObjectiveIds(t, &clientB, defaultTimeout, ids...)
	waitTimeForCompletedObjectiveIds(t, &clientI, defaultTimeout, ids...)

	for _, clientStore := range []store.Store{storeA, storeB, storeI} {
		for _, cId := range cIds {
			oId := protocols.ObjectiveId(fmt.Sprintf("VirtualDefund-%s", cId))
			o, err := clientStore.GetObjectiveById(oId)
			if err != nil {
				t.Errorf("Could not get objective: %v", err)
			}
			vdfo := o.(*virtualdefund.Objective)
			if vdfo.GetStatus() != protocols.Completed {
				t.Errorf("Expected objective %s to be completed", vdfo.Id())
			}

			// Check that the ledger outcomes get updated as expected
			switch *clientStore.GetAddress() {
			case alice.Address():
				checkAliceIreneLedgerOutcome(t, vdfo.VId(), vdfo.ToMyRight.ConsensusVars().Outcome, totalPaidToBob)
			case bob.Address():
				checkIreneBobLedger(t, vdfo.VId(), vdfo.ToMyLeft.ConsensusVars().Outcome, totalPaidToBob)
			case irene.Address():
				checkAliceIreneLedgerOutcome(t, vdfo.VId(), vdfo.ToMyLeft.ConsensusVars().Outcome, totalPaidToBob)
				checkIreneBobLedger(t, vdfo.VId(), vdfo.ToMyRight.ConsensusVars().Outcome, totalPaidToBob)
			}
		}

	}

}

// checkAliceIreneLedgerOutcome checks the ledger outcome between alice and irene is as expected
func checkAliceIreneLedgerOutcome(t *testing.T, vId types.Destination, outcome consensus_channel.LedgerOutcome, totalPaidToBob uint) {
	if outcome.IncludesTarget(vId) {
		t.Errorf("The outcome %+v should not contain a guarantee for the virtual channel %s", outcome, vId)
	}
	expectedLeaderBalance := consensus_channel.NewBalance(alice.Destination(), big.NewInt(int64(5-totalPaidToBob)))
	if diff := cmp.Diff(expectedLeaderBalance, outcome.Leader()); diff != "" {
		t.Errorf("Unexpected leader balance: %s", diff)
	}

	expectedFollowerBalance := consensus_channel.NewBalance(irene.Destination(), big.NewInt(int64(5+totalPaidToBob)))
	if diff := cmp.Diff(expectedFollowerBalance, outcome.Follower()); diff != "" {
		t.Errorf("Unexpected follower balance: %s", diff)
	}
}

// checkIreneBobLedger checks the ledger outcome between irene and bob is as expected
func checkIreneBobLedger(t *testing.T, vId types.Destination, outcome consensus_channel.LedgerOutcome, totalPaidToBob uint) {
	if outcome.IncludesTarget(vId) {
		t.Errorf("The outcome %+v should not contain a guarantee for the virtual channel %s", outcome, vId)
	}
	expectedLeaderBalance := consensus_channel.NewBalance(irene.Destination(), big.NewInt(int64(5-totalPaidToBob)))
	if diff := cmp.Diff(expectedLeaderBalance, outcome.Leader()); diff != "" {
		t.Errorf("Unexpected leader balance: %s", diff)
	}

	expectedFollowerBalance := consensus_channel.NewBalance(bob.Destination(), big.NewInt(int64(5+totalPaidToBob)))
	if diff := cmp.Diff(expectedFollowerBalance, outcome.Follower()); diff != "" {
		t.Errorf("Unexpected follower balance: %s", diff)
	}
}
