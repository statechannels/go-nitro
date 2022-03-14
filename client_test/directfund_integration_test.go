// Package client_test contains helpers and integration tests for go-nitro clients
package client_test // import "github.com/statechannels/go-nitro/client_test"

import (
	"math/big"
	"testing"

	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/client"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/types"
)

func directlyFundALedgerChannel(t *testing.T, alpha client.Client, beta client.Client) {
	// Set up an outcome that requires both participants to deposit
	outcome := outcome.Exit{outcome.SingleAssetExit{
		Allocations: outcome.Allocations{
			outcome.Allocation{
				Destination: types.AddressToDestination(*alpha.Address),
				Amount:      big.NewInt(100),
			},
			outcome.Allocation{
				Destination: types.AddressToDestination(*beta.Address),
				Amount:      big.NewInt(100),
			},
		},
	}}
	id := alpha.CreateDirectChannel(*beta.Address, types.Address{}, types.Bytes{}, outcome, big.NewInt(0))
	waitTimeForCompletedObjectiveIds(t, &alpha, defaultTimeout, id)
	waitTimeForCompletedObjectiveIds(t, &beta, defaultTimeout, id)
}
func TestDirectFundIntegration(t *testing.T) {
	logFile := "directfund_client_test.log"
	truncateLog(logFile)

	chain := chainservice.NewMockChain([]types.Address{alice, bob})

	clientA, messageserviceA := setupClient(aliceKey, chain, logFile)
	clientB, messageserviceB := setupClient(bobKey, chain, logFile)

	connectMessageServices(messageserviceA, messageserviceB)

	directlyFundALedgerChannel(t, clientA, clientB)

}
