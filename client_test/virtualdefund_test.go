// Package client_test contains helpers and integration tests for go-nitro clients
package client_test // import "github.com/statechannels/go-nitro/client_test"

import (
	"fmt"
	"io"
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/client"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	"github.com/statechannels/go-nitro/client/engine/store"
	td "github.com/statechannels/go-nitro/internal/testdata"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/protocols/virtualdefund"
	"github.com/statechannels/go-nitro/protocols/virtualfund"
	"github.com/statechannels/go-nitro/types"
)

func TestVirtualDefundIntegration(t *testing.T) {

	// Setup logging
	logFile := "test_virtual_defund.log"
	truncateLog(logFile)
	logDestination := newLogWriter(logFile)
	runVirtualDefundIntegrationTest(t, 0, defaultTimeout, logDestination)

}

func TestVirtualDefundIntegrationWithMessageDelay(t *testing.T) {

	// Setup logging
	logFile := "test_virtual_defund_with_message_delay.log"
	truncateLog(logFile)
	logDestination := newLogWriter(logFile)

	const MAX_MESSAGE_DELAY = time.Millisecond * 100
	// Since we are delaying messages we allow for enough time to complete the objective
	const OBJECTIVE_TIMEOUT = time.Second * 2

	runVirtualDefundIntegrationTest(t, MAX_MESSAGE_DELAY, OBJECTIVE_TIMEOUT, logDestination)

}

// runVirtualDefundIntegrationTest runs a virtual defund integration test using the provided message delay, objective timeout and log destination
func runVirtualDefundIntegrationTest(t *testing.T, messageDelay time.Duration, objectiveTimeout time.Duration, logDestination io.Writer) {
	chain := chainservice.NewMockChain()
	broker := messageservice.NewBroker()

	clientA, storeA := setupClient(alice.PrivateKey, chain, broker, logDestination, messageDelay)
	clientB, storeB := setupClient(bob.PrivateKey, chain, broker, logDestination, messageDelay)
	clientI, storeI := setupClient(irene.PrivateKey, chain, broker, logDestination, messageDelay)

	numOfVirtualChannels := uint(5)
	paidToBob := uint(1)
	totalPaidToBob := paidToBob * numOfVirtualChannels

	cIds := openVirtualChannels(t, clientA, clientB, clientI, numOfVirtualChannels)

	ids := make([]protocols.ObjectiveId, len(cIds))
	for i := 0; i < len(cIds); i++ {
		ids[i] = clientA.CloseVirtualChannel(cIds[i], big.NewInt(int64(paidToBob)))

	}
	waitTimeForCompletedObjectiveIds(t, &clientA, objectiveTimeout, ids...)
	waitTimeForCompletedObjectiveIds(t, &clientB, objectiveTimeout, ids...)
	waitTimeForCompletedObjectiveIds(t, &clientI, objectiveTimeout, ids...)

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
				checkIreneBobLedgerOutcome(t, vdfo.VId(), vdfo.ToMyLeft.ConsensusVars().Outcome, totalPaidToBob)
			case irene.Address():
				checkAliceIreneLedgerOutcome(t, vdfo.VId(), vdfo.ToMyLeft.ConsensusVars().Outcome, totalPaidToBob)
				checkIreneBobLedgerOutcome(t, vdfo.VId(), vdfo.ToMyRight.ConsensusVars().Outcome, totalPaidToBob)
			}
		}

	}
}

// checkAliceIreneLedgerOutcome checks the ledger outcome between alice and irene is as expected
func checkAliceIreneLedgerOutcome(t *testing.T, vId types.Destination, outcome consensus_channel.LedgerOutcome, totalPaidToBob uint) {
	if outcome.IncludesTarget(vId) {
		t.Errorf("The outcome %+v should not contain a guarantee for the virtual channel %s", outcome, vId)
	}
	expectedLeaderBalance := consensus_channel.NewBalance(alice.Destination(), big.NewInt(int64(ledgerChannelDeposit-totalPaidToBob)))
	if diff := cmp.Diff(expectedLeaderBalance, outcome.Leader()); diff != "" {
		t.Errorf("Unexpected leader balance: %s", diff)
	}

	expectedFollowerBalance := consensus_channel.NewBalance(irene.Destination(), big.NewInt(int64(ledgerChannelDeposit+totalPaidToBob)))
	if diff := cmp.Diff(expectedFollowerBalance, outcome.Follower()); diff != "" {
		t.Errorf("Unexpected follower balance: %s", diff)
	}
}

// checkIreneBobLedgerOutcome checks the ledger outcome between irene and bob is as expected
func checkIreneBobLedgerOutcome(t *testing.T, vId types.Destination, outcome consensus_channel.LedgerOutcome, totalPaidToBob uint) {
	if outcome.IncludesTarget(vId) {
		t.Errorf("The outcome %+v should not contain a guarantee for the virtual channel %s", outcome, vId)
	}
	expectedLeaderBalance := consensus_channel.NewBalance(irene.Destination(), big.NewInt(int64(ledgerChannelDeposit-totalPaidToBob)))
	if diff := cmp.Diff(expectedLeaderBalance, outcome.Leader()); diff != "" {
		t.Errorf("Unexpected leader balance: %s", diff)
	}

	expectedFollowerBalance := consensus_channel.NewBalance(bob.Destination(), big.NewInt(int64(ledgerChannelDeposit+totalPaidToBob)))
	if diff := cmp.Diff(expectedFollowerBalance, outcome.Follower()); diff != "" {
		t.Errorf("Unexpected follower balance: %s", diff)
	}
}

func TestWhenVirtualDefundObjectiveIsRejected(t *testing.T) {

	// Setup logging
	logFile := "test_rejected_virtualdefund_fund.log"
	truncateLog(logFile)
	logDestination := newLogWriter(logFile)

	chain := chainservice.NewMockChain()
	broker := messageservice.NewBroker()

	meanMessageDelay := time.Duration(0)
	clientA, storeA := setupClient(alice.PrivateKey, chain, broker, logDestination, meanMessageDelay)
	var (
		clientB client.Client
		storeB  store.Store
	)
	{
		messageservice := messageservice.NewTestMessageService(bob.Address(), broker, meanMessageDelay)
		storeB = store.NewMemStore(bob.PrivateKey)
		clientB = client.New(messageservice, chain, storeB, logDestination, &RejectingPolicyMaker{}, nil)
	}
	clientI, storeI := setupClient(irene.PrivateKey, chain, broker, logDestination, meanMessageDelay)

	directlyFundALedgerChannel(t, clientA, clientI)
	directlyFundALedgerChannel(t, clientB, clientI)

	outcome := td.Outcomes.Create(alice.Address(), bob.Address(), 1, 1)
	request := virtualfund.ObjectiveRequest{
		CounterParty:      bob.Address(),
		Intermediary:      irene.Address(),
		Outcome:           outcome,
		AppDefinition:     types.Address{},
		AppData:           types.Bytes{},
		ChallengeDuration: big.NewInt(0),
		Nonce:             rand.Int63(),
	}
	response := clientA.CreateVirtualChannel(request)

	waitTimeForCompletedObjectiveIds(t, &clientA, time.Second, response.Id)
	waitTimeForCompletedObjectiveIds(t, &clientB, time.Second, response.Id)
	waitTimeForCompletedObjectiveIds(t, &clientI, time.Second, response.Id)

	obj, _ := storeA.GetObjectiveById(response.Id)

	if obj.GetStatus() != protocols.Rejected {
		t.Error("expected objective to be rejected")
		t.FailNow()
	}

	obj, _ = storeB.GetObjectiveById(response.Id)

	if obj.GetStatus() != protocols.Rejected {
		t.Error("expected objective to be rejected")
		t.FailNow()
	}

	obj, _ = storeI.GetObjectiveById(response.Id)

	if obj.GetStatus() != protocols.Rejected {
		t.Error("expected objective to be rejected")
		t.FailNow()
	}
}
