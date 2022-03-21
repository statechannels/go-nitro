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
	"github.com/statechannels/go-nitro/protocols/directfund"
	"github.com/statechannels/go-nitro/types"
)

func directlyFundALedgerChannel(t *testing.T, alpha client.Client, beta client.Client) {
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
	waitTimeForCompletedObjectiveIds(t, &alpha, defaultTimeout, id)
	waitTimeForCompletedObjectiveIds(t, &beta, defaultTimeout, id)
}
func TestDirectFundIntegration(t *testing.T) {
	logFile := "directfund_client_test.log"
	truncateLog(logFile)

	chain := chainservice.NewMockChain()
	broker := messageservice.NewBroker()

	clientA := setupClient(aliceKey, chain, broker, logFile, 1)
	clientB := setupClient(bobKey, chain, broker, logFile, 1)

	directlyFundALedgerChannel(t, clientA, clientB)

}
