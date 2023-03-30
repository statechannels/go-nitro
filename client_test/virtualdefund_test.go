// Package client_test contains helpers and integration tests for go-nitro clients
package client_test // import "github.com/statechannels/go-nitro/client_test"

import (
	"fmt"
	"io"
	"math/big"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/client"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	"github.com/statechannels/go-nitro/client/engine/store"
	"github.com/statechannels/go-nitro/internal/testactors"
	td "github.com/statechannels/go-nitro/internal/testdata"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/protocols/virtualdefund"

	"github.com/statechannels/go-nitro/types"
)

func TestVirtualDefundIntegration(t *testing.T) {
	// Setup logging
	logFile := "test_virtual_defund.log"
	truncateLog(logFile)
	logDestination := newLogWriter(logFile)
	for _, closer := range []testactors.Actor{alice, irene, bob} {
		t.Run(fmt.Sprintf("TestVirtualDefundIntegration_as_%s", closer.Name), func(t *testing.T) {
			runVirtualDefundIntegrationTestAs(t, closer.Address(), 0, defaultTimeout, logDestination)
		})
	}
}

func TestVirtualDefundIntegrationWithMessageDelay(t *testing.T) {
	// Setup logging
	logFile := "test_virtual_defund_with_message_delay.log"
	truncateLog(logFile)
	logDestination := newLogWriter(logFile)

	const MAX_MESSAGE_DELAY = time.Millisecond * 100
	// Since we are delaying messages we allow for enough time to complete the objective
	const OBJECTIVE_TIMEOUT = time.Second * 2

	runVirtualDefundIntegrationTestAs(t, alice.Address(), MAX_MESSAGE_DELAY, OBJECTIVE_TIMEOUT, logDestination)
}

// runVirtualDefundIntegrationTestAs runs a virtual defund integration test using the provided message delay, objective timeout and log destination
func runVirtualDefundIntegrationTestAs(t *testing.T, closer types.Address, messageDelay time.Duration, objectiveTimeout time.Duration, logDestination io.Writer) {
	chain := chainservice.NewMockChain()
	chainServiceA := chainservice.NewMockChainService(chain, alice.Address())
	chainServiceB := chainservice.NewMockChainService(chain, bob.Address())
	chainServiceI := chainservice.NewMockChainService(chain, irene.Address())
	broker := messageservice.NewBroker()

	clientA, storeA := setupClient(alice.PrivateKey, chainServiceA, broker, logDestination, messageDelay)
	defer closeClient(t, &clientA)
	clientB, storeB := setupClient(bob.PrivateKey, chainServiceB, broker, logDestination, messageDelay)
	defer closeClient(t, &clientB)
	clientI, storeI := setupClient(irene.PrivateKey, chainServiceI, broker, logDestination, messageDelay)
	defer closeClient(t, &clientI)

	numOfVirtualChannels := uint(5)
	paidToBob := uint(1)
	totalPaidToBob := paidToBob * numOfVirtualChannels

	cIds := openVirtualChannels(t, clientA, clientB, clientI, numOfVirtualChannels)
	for i := 0; i < len(cIds); i++ {
		clientA.Pay(cIds[i], big.NewInt(int64(paidToBob)))
	}
	ids := make([]protocols.ObjectiveId, len(cIds))
	for i := 0; i < len(cIds); i++ {
		switch closer {
		case alice.Address():
			ids[i] = clientA.CloseVirtualChannel(cIds[i])
		case bob.Address():
			ids[i] = clientB.CloseVirtualChannel(cIds[i])
		case irene.Address():
			ids[i] = clientI.CloseVirtualChannel(cIds[i])
		}
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
	chainServiceA := chainservice.NewMockChainService(chain, alice.Address())
	chainServiceB := chainservice.NewMockChainService(chain, bob.Address())
	chainServiceI := chainservice.NewMockChainService(chain, irene.Address())
	broker := messageservice.NewBroker()

	meanMessageDelay := time.Duration(0)
	clientA, storeA := setupClient(alice.PrivateKey, chainServiceA, broker, logDestination, meanMessageDelay)
	defer closeClient(t, &clientA)
	var (
		clientB client.Client
		storeB  store.Store
	)
	{
		messageservice := messageservice.NewTestMessageService(bob.Address(), broker, meanMessageDelay)
		storeB = store.NewMemStore(bob.PrivateKey)
		clientB = client.New(messageservice, chainServiceB, storeB, logDestination, &RejectingPolicyMaker{}, nil)
	}
	defer closeClient(t, &clientB)
	clientI, storeI := setupClient(irene.PrivateKey, chainServiceI, broker, logDestination, meanMessageDelay)
	defer closeClient(t, &clientI)

	directlyFundALedgerChannel(t, clientA, clientI, types.Address{})
	directlyFundALedgerChannel(t, clientB, clientI, types.Address{})

	outcome := td.Outcomes.Create(alice.Address(), bob.Address(), 1, 1, types.Address{})
	response := clientA.CreateVirtualPaymentChannel(
		[]types.Address{irene.Address()},
		bob.Address(),
		0,
		outcome,
	)

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
