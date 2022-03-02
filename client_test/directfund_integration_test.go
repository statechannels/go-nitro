// Package client_test contains helpers and integration tests for go-nitro clients
package client_test // import "github.com/statechannels/go-nitro/client_test"

import (
	"log"
	"math/big"
	"os"
	"testing"

	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/client"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	"github.com/statechannels/go-nitro/client/engine/store"
	"github.com/statechannels/go-nitro/crypto"
	"github.com/statechannels/go-nitro/types"
)

func directlyFundALedgerChannel(alpha client.Client, beta client.Client) {
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
	id := alpha.CreateDirectChannel(*beta.Address, types.Address{}, types.Bytes{}, outcome, big.NewInt(0))
	waitForCompletedObjectiveId(id, &alpha)
	waitForCompletedObjectiveId(id, &beta)
}
func TestDirectFundIntegration(t *testing.T) {

	// Set up logging
	logDestination, err := os.OpenFile("directfund_client_test.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	// Reset log destination file
	err = logDestination.Truncate(0)
	if err != nil {
		log.Fatal(err)
	}

	aKey, a := crypto.GeneratePrivateKeyAndAddress()
	bKey, b := crypto.GeneratePrivateKeyAndAddress()
	chain := chainservice.NewMockChain([]types.Address{a, b})

	chainservA := chainservice.NewSimpleChainService(chain, a)
	messageserviceA := messageservice.NewTestMessageService(a)
	storeA := store.NewMockStore(aKey)
	clientA := client.New(messageserviceA, chainservA, storeA, logDestination)

	chainservB := chainservice.NewSimpleChainService(chain, b)
	messageserviceB := messageservice.NewTestMessageService(b)
	storeB := store.NewMockStore(bKey)
	clientB := client.New(messageserviceB, chainservB, storeB, logDestination)

	connectMessageServices(messageserviceA, messageserviceB)

	directlyFundALedgerChannel(clientA, clientB)

}
