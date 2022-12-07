// Package client_test contains helpers and integration tests for go-nitro clients
package client_test // import "github.com/statechannels/go-nitro/client_test"

import (
	"math/big"
	"testing"

	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/client"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	"github.com/statechannels/go-nitro/internal/testactors"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

func TestClientConcurrency(t *testing.T) {
	logDestination := newLogWriter("concurrency_client_test.log")
	truncateLog("concurrency_client_test.log")

	chain := chainservice.NewMockChain()
	chainServiceA := chainservice.NewMockChainService(chain, alice.Address())
	chainServiceB := chainservice.NewMockChainService(chain, alice.Address())
	broker := messageservice.NewBroker()

	clientA, _ := setupClient(testactors.Alice.PrivateKey, chainServiceA, broker, logDestination, 1)
	clientB, _ := setupClient(testactors.Bob.PrivateKey, chainServiceB, broker, logDestination, 1)
	// TODO: Only directly funded channels can be handled concurrently
	// Due to this ledger issue: https://github.com/statechannels/go-nitro/issues/366
	ids := createDirectlyFundedChannels(t, clientA, clientB, 1)
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

		res := alpha.CreateLedgerChannel(*beta.Address, 100, outcome)
		ids = append(ids, res.Id)
	}
	return ids
}
