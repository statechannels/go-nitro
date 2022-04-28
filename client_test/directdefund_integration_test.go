// Package client_test contains helpers and integration tests for go-nitro clients
package client_test // import "github.com/statechannels/go-nitro/client_test"

import (
	"bytes"
	"testing"

	"github.com/statechannels/go-nitro/client"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	"github.com/statechannels/go-nitro/client/engine/store"
	"github.com/statechannels/go-nitro/types"
)

func directlyDefundALedgerChannel(t *testing.T, alpha client.Client, beta client.Client, channelId types.Destination) {

	id := alpha.CloseDirectChannel(channelId)
	waitTimeForCompletedObjectiveIds(t, &alpha, defaultTimeout, id)
	waitTimeForCompletedObjectiveIds(t, &beta, defaultTimeout, id)

}
func TestDirectDefundIntegration(t *testing.T) {

	// Setup logging
	logDestination := &bytes.Buffer{}
	t.Cleanup(flushToFileCleanupFn(logDestination, "directdefund_client_test.log"))

	chain := chainservice.NewMockChain()
	broker := messageservice.NewBroker()

	clientA, storeA := setupClient(alice.PrivateKey, chain, broker, logDestination, 0)
	clientB, storeB := setupClient(bob.PrivateKey, chain, broker, logDestination, 0)

	channelId := directlyFundALedgerChannel(t, clientA, clientB)
	directlyDefundALedgerChannel(t, clientA, clientB, channelId)

	// Ensure that we no longer have a consensus channel in the store
	// And that we have a regular Channel instead
	for _, store := range []store.Store{storeA, storeB} {

		c, channelInStore := store.GetChannelById(channelId)

		if !channelInStore {
			t.Fatalf("expected a Channel to have been created")
		}

		if c.OnChainFunding.IsNonZero() {
			t.Fatal("Expected zero on chain funding, but got nonzero")
		}

		_, err := store.GetConsensusChannelById(channelId)
		if consensusChannelStillInStore := err != nil; consensusChannelStillInStore {
			t.Fatalf("Expected ConsensusChannel to have been destroyed in %v's store, but it was not", store.GetAddress())
		}

	}

}
