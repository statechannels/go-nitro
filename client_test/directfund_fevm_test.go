package client_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/go-cmp/cmp"
	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	"github.com/statechannels/go-nitro/client/engine/store"
	"github.com/statechannels/go-nitro/internal/testdata"
)

// TestDirectFund uses the geth simulated backend
func TestFevmDirectFund(t *testing.T) {
	// Setup logging
	logFile := "test_fevm_direct_fund.log"
	truncateLog(logFile)
	logDestination := newLogWriter(logFile)
	const pkString = "716b7161580785bc96a4344eb52d23131aea0caf42a52dcf9f8aee9eef9dc3cd"

	// Setup chain service
	pk, _ := crypto.HexToECDSA(pkString)
	chainA := chainservice.NewFevmChainService(pk)
	chainB := chainservice.NewFevmChainService(pk)
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
