package client_test

import (
	"math/big"
	"math/rand"
	"testing"

	"github.com/statechannels/go-nitro/client"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	"github.com/statechannels/go-nitro/client/engine/store"
	td "github.com/statechannels/go-nitro/internal/testdata"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/protocols/virtualfund"
	"github.com/statechannels/go-nitro/types"
)

func openVirtualChannels(t *testing.T, clientA client.Client, clientB client.Client, clientI client.Client, numOfChannels uint) []types.Destination {
	directlyFundALedgerChannel(t, clientA, clientI)
	directlyFundALedgerChannel(t, clientI, clientB)

	objectiveIds := make([]protocols.ObjectiveId, numOfChannels)
	channelIds := make([]types.Destination, numOfChannels)
	for i := 0; i < int(numOfChannels); i++ {
		outcome := td.Outcomes.Create(alice.Address(), bob.Address(), 100, 1)
		request := virtualfund.ObjectiveRequest{
			CounterParty:      bob.Address(),
			Intermediary:      irene.Address(),
			Outcome:           outcome,
			AppDefinition:     types.Address{},
			AppData:           types.Bytes{},
			ChallengeDuration: big.NewInt(0),
			Nonce:             rand.Int63(),
		}
		response := clientA.CreateVirtualChannel(request)
		objectiveIds[i] = response.Id
		channelIds[i] = response.ChannelId
	}
	waitTimeForCompletedObjectiveIds(t, &clientA, defaultTimeout, objectiveIds...)
	waitTimeForCompletedObjectiveIds(t, &clientB, defaultTimeout, objectiveIds...)
	waitTimeForCompletedObjectiveIds(t, &clientI, defaultTimeout, objectiveIds...)

	return channelIds

}
func TestVirtualFundIntegration(t *testing.T) {

	// Setup logging
	logFile := "test_virtual_fund.log"
	truncateLog(logFile)
	logDestination := newLogWriter(logFile)

	chain := chainservice.NewMockChain()
	broker := messageservice.NewBroker()

	clientA, storeA := setupClient(alice.PrivateKey, chain, broker, logDestination, 0)
	clientB, storeB := setupClient(bob.PrivateKey, chain, broker, logDestination, 0)
	clientI, _ := setupClient(irene.PrivateKey, chain, broker, logDestination, 0)

	cIds := openVirtualChannels(t, clientA, clientB, clientI, 1)

	cId := cIds[0] // just pick a a channel

	printBalances(t, storeA, storeB, cId)
	for i := 0; i < 100; i++ {
		clientA.MakePayment(cId, bob.Destination(), big.NewInt(1))
	}
	printBalances(t, storeA, storeB, cId)
}

func printBalances(t *testing.T, storeA store.Store, storeB store.Store, cId types.Destination) {
	chA, ok := storeA.GetChannelById(cId)
	if !ok {
		t.Fatal()
	}
	ss, _ := chA.LatestSignedState()
	aBal := ss.State().Outcome.TotalAllocatedFor(alice.Destination())[types.Address{0}]
	bBal := ss.State().Outcome.TotalAllocatedFor(bob.Destination())[types.Address{0}]
	t.Log("Alice reports Alice:" + aBal.String() + " Bob:" + bBal.String())
	chB, ok := storeA.GetChannelById(cId)
	if !ok {
		t.Fatal()
	}
	ss, _ = chB.LatestSignedState()
	aBal = ss.State().Outcome.TotalAllocatedFor(alice.Destination())[types.Address{0}]
	bBal = ss.State().Outcome.TotalAllocatedFor(bob.Destination())[types.Address{0}]
	t.Log("Alice reports Alice:" + aBal.String() + " Bob:" + bBal.String())
}
