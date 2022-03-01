package client

import (
	"io"
	"log"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	"github.com/statechannels/go-nitro/client/engine/store"
	"github.com/statechannels/go-nitro/crypto"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

func connectMessageServices(services []messageservice.TestMessageService) {
	for i, ms := range services {
		for j, ms2 := range services {
			if i != j {
				ms.Connect(ms2)
			}
		}
	}
}

func setupClient(pk []byte, chain chainservice.MockChain, logDestination io.Writer) (Client, messageservice.TestMessageService) {
	myAddress := crypto.GetAddressFromSecretKeyBytes(pk)
	chainservice := chainservice.NewSimpleChainService(chain, myAddress)
	messageservice := messageservice.NewTestMessageService(myAddress)
	storeA := store.NewMockStore(pk)
	return New(messageservice, chainservice, storeA, logDestination), messageservice
}

func createOutcome(first types.Address, second types.Address) outcome.Exit {

	return outcome.Exit{outcome.SingleAssetExit{
		Allocations: outcome.Allocations{
			outcome.Allocation{
				Destination: types.AddressToDestination(first),
				Amount:      big.NewInt(5),
			},
			outcome.Allocation{
				Destination: types.AddressToDestination(second),
				Amount:      big.NewInt(5),
			},
		},
	}}
}

func TestMultiPartyVirtualFundIntegration(t *testing.T) {

	// Set up logging
	logDestination, err := os.OpenFile("virtualfund_multiparty_integration_test.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	// Reset log destination file
	err = logDestination.Truncate(0)
	if err != nil {
		log.Fatal(err)
	}

	amyKey := common.Hex2Bytes("ae1e8dc688f74ffedf8f4a67f511a060b97e1e32119e2e92fb5563405481a4df")
	amy := common.HexToAddress("0xA2A2D4f2E6FA19F62fC14d970C881AbF265E88fF")

	brianKey := common.Hex2Bytes("0aca28ba64679f63d71e671ab4dbb32aaa212d4789988e6ca47da47601c18fe2")
	brian := common.HexToAddress("0xB2B22ec3889d11f2ddb1A1Db11e80D20EF367c01")

	aliceKey := common.Hex2Bytes(`2d999770f7b5d49b694080f987b82bbc9fc9ac2b4dcc10b0f8aba7d700f69c6d`)
	alice := common.HexToAddress(`0xAAA6628Ec44A8a742987EF3A114dDFE2D4F7aDCE`)

	bobKey := common.Hex2Bytes(`0279651921cd800ac560c21ceea27aab0107b67daf436cdd25ce84cad30159b4`)
	bob := common.HexToAddress(`0xBBB676f9cFF8D242e9eaC39D063848807d3D1D94`)

	ireneKey := common.Hex2Bytes(`febb3b74b0b52d0976f6571d555f4ac8b91c308dfa25c7b58d1e6a7c3f50c781`)
	irene := common.HexToAddress(`0x111A00868581f73AB42FEEF67D235Ca09ca1E8db`)

	chain := chainservice.NewMockChain([]types.Address{alice, bob, irene, amy, brian})

	clientAlice, aliceMS := setupClient(aliceKey, chain, logDestination)
	clientBob, bobMS := setupClient(bobKey, chain, logDestination)

	clientAmy, amyMS := setupClient(amyKey, chain, logDestination)
	clientBrian, brianMS := setupClient(brianKey, chain, logDestination)

	clientIrene, ireneMS := setupClient(ireneKey, chain, logDestination)
	connectMessageServices([]messageservice.TestMessageService{aliceMS, bobMS, ireneMS, amyMS, brianMS})

	directlyFundALedgerChannel := func(alpha Client, beta Client) {
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

	directlyFundALedgerChannel(clientAlice, clientIrene)
	directlyFundALedgerChannel(clientIrene, clientBob)

	directlyFundALedgerChannel(clientAmy, clientIrene)
	directlyFundALedgerChannel(clientIrene, clientBrian)

	id := clientAlice.CreateVirtualChannel(bob, irene, types.Address{}, types.Bytes{}, createOutcome(alice, bob), big.NewInt(0))
	id2 := clientAmy.CreateVirtualChannel(brian, irene, types.Address{}, types.Bytes{}, createOutcome(amy, brian), big.NewInt(0))

	waitForCompletedObjectiveId(id, &clientAlice)
	waitForCompletedObjectiveId(id, &clientBob)
	waitForCompletedObjectiveId(id2, &clientAmy)
	waitForCompletedObjectiveId(id2, &clientBrian)

	waitForCompletedObjectiveIds([]protocols.ObjectiveId{id, id2}, &clientIrene)

}

// waitForCompletedObjectiveIds waits for completed objectives and returns when the all objective ids provided have been completed.
func waitForCompletedObjectiveIds(ids []protocols.ObjectiveId, client *Client) {
	// Create a map of all objective ids to wait for and set to false
	completed := make(map[protocols.ObjectiveId]bool)
	for _, id := range ids {
		completed[id] = false
	}
	// We continue to consume completed objective ids from the chan until all have been completed
	for got := range client.completedObjectives {
		// Mark the objective as completed
		completed[got] = true

		// If all objectives are completed we can return
		isDone := true
		for _, objectiveCompleted := range completed {
			isDone = isDone && objectiveCompleted
		}
		if isDone {
			return
		}
	}
}
