// Package client_test contains helpers and integration tests for go-nitro clients
package client_test // import "github.com/statechannels/go-nitro/client_test"

import (
	"math/big"
	"math/rand"
	"testing"

	"github.com/statechannels/go-nitro/client"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	"github.com/statechannels/go-nitro/internal/testdata"
	"github.com/statechannels/go-nitro/protocols/directfund"
	"github.com/statechannels/go-nitro/types"
)

func directlyFundALedgerChannel(t *testing.T, alpha client.Client, beta client.Client) {
	// Set up an outcome that requires both participants to deposit
	outcome := testdata.Outcomes.Create(*alpha.Address, *beta.Address, 5, 5)

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

	clientA := setupClient(alice.PrivateKey, chain, broker, logFile, 0)
	clientB := setupClient(bob.PrivateKey, chain, broker, logFile, 0)

	directlyFundALedgerChannel(t, clientA, clientB)

}
