// Package client_test contains helpers and integration tests for go-nitro clients
package client_test // import "github.com/statechannels/go-nitro/client_test"

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/client"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	"github.com/statechannels/go-nitro/client/engine/store"
	"github.com/statechannels/go-nitro/internal/testdata"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

const ledgerChannelDeposit = 5_000_000

func directlyFundALedgerChannel(t *testing.T, alpha client.Client, beta client.Client) types.Destination {
	// Set up an outcome that requires both participants to deposit
	outcome := testdata.Outcomes.Create(*alpha.Address, *beta.Address, ledgerChannelDeposit, ledgerChannelDeposit)
	response := alpha.CreateLedgerChannel(*beta.Address, 0, outcome)

	waitTimeForCompletedObjectiveIds(t, &alpha, defaultTimeout, response.Id)
	waitTimeForCompletedObjectiveIds(t, &beta, defaultTimeout, response.Id)
	return response.ChannelId
}

type RejectingPolicyMaker struct{}

func (pm *RejectingPolicyMaker) ShouldApprove(obj protocols.Objective) bool {
	return false
}

func TestWhenObjectiveIsRejected(t *testing.T) {

	// Setup logging
	logFile := "test_direct_fund.log"
	truncateLog(logFile)
	logDestination := newLogWriter(logFile)

	chain := chainservice.NewMockChain()
	chainServiceA := chainservice.NewMockChainService(chain, alice.Address())
	chainServiceB := chainservice.NewMockChainService(chain, alice.Address())
	broker := messageservice.NewBroker()

	meanMessageDelay := time.Duration(0)
	clientA, storeA := setupClient(alice.PrivateKey, chainServiceA, broker, logDestination, meanMessageDelay)
	var storeB store.Store
	{
		messageservice := messageservice.NewTestMessageService(bob.Address(), broker, meanMessageDelay)
		storeB = store.NewMemStore(bob.PrivateKey)
		_ = client.New(messageservice, chainServiceB, storeB, logDestination, &RejectingPolicyMaker{}, nil)
	}

	outcome := testdata.Outcomes.Create(alice.Address(), bob.Address(), ledgerChannelDeposit, ledgerChannelDeposit)
	response := clientA.CreateLedgerChannel(bob.Address(), 0, outcome)

	<-clientA.WaitForObjectivesToComplete(context.Background(), response.Id)

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
}

// TestDirectFund uses the geth simulated backend
func TestDirectFund(t *testing.T) {

	// Setup logging
	logFile := "test_direct_fund.log"
	truncateLog(logFile)
	logDestination := newLogWriter(logFile)

	// Setup chain service
	sim, bindings, ethAccounts, err := chainservice.SetupSimulatedBackend(2)
	if err != nil {
		t.Fatal(err)
	}
	chainA, err := chainservice.NewSimulatedBackendChainService(sim, bindings, ethAccounts[0], logDestination)
	if err != nil {
		t.Fatal(err)
	}
	chainB, err := chainservice.NewSimulatedBackendChainService(sim, bindings, ethAccounts[1], logDestination)
	if err != nil {
		t.Fatal(err)
	}
	// End chain service setup

	broker := messageservice.NewBroker()

	clientA, storeA := setupClient(alice.PrivateKey, chainA, broker, logDestination, 0)
	clientB, storeB := setupClient(bob.PrivateKey, chainB, broker, logDestination, 0)

	directlyFundALedgerChannel(t, clientA, clientB)

	want := testdata.Outcomes.Create(*clientA.Address, *clientB.Address, ledgerChannelDeposit, ledgerChannelDeposit)
	// Ensure that we create a consensus channel in the store
	for _, store := range []store.Store{storeA, storeB} {
		var con *consensus_channel.ConsensusChannel
		var ok bool

		// each client fetches the ConsensusChannel by reference to their counterparty
		if store.GetChannelSecretKey() == &alice.PrivateKey {
			con, ok = store.GetConsensusChannel(*clientB.Address)
		} else {
			con, ok = store.GetConsensusChannel(*clientA.Address)
		}

		if !ok {
			t.Fatalf("expected a consensus channel to have been created")
		}
		vars := con.ConsensusVars()
		got := vars.Outcome.AsOutcome()

		if diff := cmp.Diff(want, got); diff != "" {
			t.Fatalf("expected outcome to be %v, got %v:\n %v", want, got, diff)
		}
		if vars.TurnNum != 1 {
			t.Fatal("expected consensus turn number to be the post fund setup 1, received #$v", vars.TurnNum)
		}
		if con.Leader() != *clientA.Address {
			t.Fatalf("Expected %v as leader, but got %v", clientA.Address, con.Leader())
		}

		if !con.OnChainFunding.IsNonZero() {
			t.Fatal("Expected nonzero on chain funding, but got zero")
		}

		if _, channelStillInStore := store.GetChannelById(con.Id); channelStillInStore {
			t.Fatalf("Expected channel to have been destroyed in %v's store, but it was not", store.GetAddress())
		}

	}

}
