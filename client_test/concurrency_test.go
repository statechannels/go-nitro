// Package client_test contains helpers and integration tests for go-nitro clients
package client_test // import "github.com/statechannels/go-nitro/client_test"

import (
	"math/big"
	"math/rand"
	"testing"

	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/client"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/protocols/directfund"
	"github.com/statechannels/go-nitro/types"
)

func TestClientConcurrency(t *testing.T) {
	logFile := "concurrency_client_test.log"
	truncateLog(logFile)

	chain := chainservice.NewMockChain()
	broker := messageservice.NewBroker()

	clientA := setupClient(aliceKey, chain, broker, logFile, 5)
	clientB := setupClient(bobKey, chain, broker, logFile, 5)
	// TODO: Only directly funded channels can be handled concurrently
	// Due to this ledger issue: https://github.com/statechannels/go-nitro/issues/366
	ids := createDirectlyFundedChannels(t, clientA, clientB, 100)
	waitTimeForCompletedObjectiveIds(t, &clientA, defaultTimeout, ids...)
	waitTimeForCompletedObjectiveIds(t, &clientB, defaultTimeout, ids...)

}

// createDirectlyFundedChannels creates a number of channels between two clients.
// It returns a slice of objective ids for the channels that were created.
func createDirectlyFundedChannels(t *testing.T, alpha client.Client, beta client.Client, numChannels uint) []protocols.ObjectiveId {
	ids := []protocols.ObjectiveId{}
	// Set up an outcome that requires both participants to deposit
	outcome := outcome.Exit{outcome.SingleAssetExit{
		Allocations: outcome.Allocations{
			outcome.Allocation{
				Destination: types.AddressToDestination(*alpha.Address),
				Amount:      big.NewInt(5),
			},
			outcome.Allocation{
				Destination: types.AddressToDestination(*beta.Address),
				Amount:      big.NewInt(5),
			},
		},
	}}
	for i := uint(0); i < numChannels; i++ {
		request := directfund.ObjectiveRequest{
			MyAddress:         *alpha.Address,
			CounterParty:      *beta.Address,
			Outcome:           outcome,
			AppDefinition:     types.Address{},
			AppData:           types.Bytes{},
			ChallengeDuration: big.NewInt(0),
			Nonce:             rand.Int63(),
		}
		id := alpha.CreateDirectChannel(request)
		ids = append(ids, id)
	}
	return ids
}
