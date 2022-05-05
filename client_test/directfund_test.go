// Package client_test contains helpers and integration tests for go-nitro clients
package client_test // import "github.com/statechannels/go-nitro/client_test"

import (
	"math/big"
	"math/rand"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/client"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	"github.com/statechannels/go-nitro/client/engine/store"
	"github.com/statechannels/go-nitro/internal/testdata"
	"github.com/statechannels/go-nitro/protocols/directfund"
	"github.com/statechannels/go-nitro/types"
)

const ledgerChannelDeposit = 5_000_000

func directlyFundALedgerChannel(t *testing.T, alpha client.Client, beta client.Client) types.Destination {
	// Set up an outcome that requires both participants to deposit
	outcome := testdata.Outcomes.Create(*alpha.Address, *beta.Address, ledgerChannelDeposit, ledgerChannelDeposit)

	request := directfund.ObjectiveRequest{
		MyAddress:         *alpha.Address,
		CounterParty:      *beta.Address,
		Outcome:           outcome,
		AppDefinition:     types.Address{},
		AppData:           types.Bytes{},
		ChallengeDuration: big.NewInt(0),
		Nonce:             rand.Int63(),
	}
	response := alpha.CreateDirectChannel(request)

	waitTimeForCompletedObjectiveIds(t, &alpha, defaultTimeout, response.Id)
	waitTimeForCompletedObjectiveIds(t, &beta, defaultTimeout, response.Id)
	return response.ChannelId
}
func TestDirectFund(t *testing.T) {

	// Setup logging
	logFile := "test_direct_fund.log"
	truncateLog(logFile)
	logDestination := newLogWriter(logFile)

	chain := chainservice.NewMockChain()
	broker := messageservice.NewBroker()

	clientA, storeA := setupClient(alice.PrivateKey, chain, broker, logDestination, 0)
	clientB, storeB := setupClient(bob.PrivateKey, chain, broker, logDestination, 0)

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
