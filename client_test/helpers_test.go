package client_test

import (
	"io"
	"math/big"

	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/client"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	"github.com/statechannels/go-nitro/client/engine/store"
	"github.com/statechannels/go-nitro/crypto"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

// waitForCompletedObjectiveId waits for completed objectives and returns when the completed objective id matchs the id waitForCompletedObjectiveId has been given
func waitForCompletedObjectiveId(id protocols.ObjectiveId, client *client.Client) {
	got := <-client.CompletedObjectives()
	for got != id {
		got = <-client.CompletedObjectives()
	}
}

// waitForCompletedObjectiveIds waits for completed objectives and returns when the all objective ids provided have been completed.
func waitForCompletedObjectiveIds(ids []protocols.ObjectiveId, client *client.Client) { //nolint:golint,unused
	// Create a map of all objective ids to wait for and set to false
	completed := make(map[protocols.ObjectiveId]bool)
	for _, id := range ids {
		completed[id] = false
	}
	// We continue to consume completed objective ids from the chan until all have been completed
	for got := range client.CompletedObjectives() {
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

// connectMessageServices connects the message services together so any message service can communicate with another.
func connectMessageServices(services ...messageservice.TestMessageService) {
	for i, ms := range services {
		for j, ms2 := range services {
			if i != j {
				ms.Connect(ms2)
			}
		}
	}
}

// setupClient is a helper function that contructs a client and returns the new client and message service.
func setupClient(pk []byte, chain chainservice.MockChain, logDestination io.Writer) (client.Client, messageservice.TestMessageService) {
	myAddress := crypto.GetAddressFromSecretKeyBytes(pk)
	chainservice := chainservice.NewSimpleChainService(chain, myAddress)
	messageservice := messageservice.NewTestMessageService(myAddress)
	storeA := store.NewMockStore(pk)
	return client.New(messageservice, chainservice, storeA, logDestination), messageservice
}

// createVirtualOutcome is a helper function to create the outcome for two participants for a virtual channel.
func createVirtualOutcome(first types.Address, second types.Address) outcome.Exit {

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
