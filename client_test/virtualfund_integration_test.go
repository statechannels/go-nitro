package client_test

import (
	"log"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/client"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	"github.com/statechannels/go-nitro/client/engine/store"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

// TestVirtualFundIntegration is a work in progress:
// It should:
// [x] spin up three clients with connected test services
// [x] directly fund a pair of ledger channels
// [x] call an API method such as clientA.CreateVirtualChannel
// [x] assert on an appropriate objective completing in all clients
func TestVirtualFundIntegration(t *testing.T) {

	// Set up logging
	logDestination, err := os.OpenFile("virtualfund_client_test.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	// Reset log destination file
	err = logDestination.Truncate(0)
	if err != nil {
		log.Fatal(err)
	}

	aKey := common.Hex2Bytes(`2d999770f7b5d49b694080f987b82bbc9fc9ac2b4dcc10b0f8aba7d700f69c6d`)
	a := common.HexToAddress(`0xAAA6628Ec44A8a742987EF3A114dDFE2D4F7aDCE`)

	bKey := common.Hex2Bytes(`0279651921cd800ac560c21ceea27aab0107b67daf436cdd25ce84cad30159b4`)
	b := common.HexToAddress(`0xBBB676f9cFF8D242e9eaC39D063848807d3D1D94`)

	iKey := common.Hex2Bytes(`febb3b74b0b52d0976f6571d555f4ac8b91c308dfa25c7b58d1e6a7c3f50c781`)
	i := common.HexToAddress(`0x111A00868581f73AB42FEEF67D235Ca09ca1E8db`)

	chain := chainservice.NewMockChain([]types.Address{a, b, i})

	chainservA := chainservice.NewSimpleChainService(chain, a)
	messageserviceA := messageservice.NewTestMessageService(a)
	storeA := store.NewMockStore(aKey)
	clientA := client.New(messageserviceA, chainservA, storeA, logDestination)

	chainservB := chainservice.NewSimpleChainService(chain, b)
	messageserviceB := messageservice.NewTestMessageService(b)
	storeB := store.NewMockStore(bKey)
	clientB := client.New(messageserviceB, chainservB, storeB, logDestination)

	chainservI := chainservice.NewSimpleChainService(chain, i)
	messageserviceI := messageservice.NewTestMessageService(i)
	storeI := store.NewMockStore(iKey)
	clientI := client.New(messageserviceI, chainservI, storeI, logDestination)

	messageserviceA.Connect(messageserviceB)
	messageserviceA.Connect(messageserviceI)

	messageserviceB.Connect(messageserviceA)
	messageserviceB.Connect(messageserviceI)

	messageserviceI.Connect(messageserviceA)
	messageserviceI.Connect(messageserviceB)

	directlyFundALedgerChannel := func(alpha client.Client, beta client.Client) {
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

	directlyFundALedgerChannel(clientA, clientI)
	directlyFundALedgerChannel(clientI, clientB)

	outcome := outcome.Exit{outcome.SingleAssetExit{
		Allocations: outcome.Allocations{
			outcome.Allocation{
				Destination: types.AddressToDestination(a),
				Amount:      big.NewInt(5),
			},
			outcome.Allocation{
				Destination: types.AddressToDestination(b),
				Amount:      big.NewInt(5),
			},
		},
	}}
	id := clientA.CreateVirtualChannel(b, i, types.Address{}, types.Bytes{}, outcome, big.NewInt(0))
	waitForCompletedObjectiveId(id, &clientA)
	waitForCompletedObjectiveId(id, &clientB)
	waitForCompletedObjectiveId(id, &clientI)

}

// waitForCompletedObjectiveId waits for completed objectives and returns when the completed objective id matchs the id waitForCompletedObjectiveId has been given
func waitForCompletedObjectiveId(id protocols.ObjectiveId, client *client.Client) {
	got := <-client.CompletedObjectives()
	for got != id {
		got = <-client.CompletedObjectives()
	}
}
