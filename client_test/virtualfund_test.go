package client_test

import (
	"math/big"
	"math/rand"
	"testing"

	"github.com/statechannels/go-nitro/client"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	td "github.com/statechannels/go-nitro/internal/testdata"
	"github.com/statechannels/go-nitro/protocols/virtualfund"
	"github.com/statechannels/go-nitro/types"
)

func openAVirtualChannel(t *testing.T, clientA client.Client, clientB client.Client, clientI client.Client) types.Destination {
	directlyFundALedgerChannel(t, clientA, clientI)
	directlyFundALedgerChannel(t, clientI, clientB)

	outcome := td.Outcomes.Create(alice.Address(), bob.Address(), 1, 1)
	request := virtualfund.ObjectiveRequest{
		MyAddress:         alice.Address(),
		CounterParty:      bob.Address(),
		Intermediary:      irene.Address(),
		Outcome:           outcome,
		AppDefinition:     types.Address{},
		AppData:           types.Bytes{},
		ChallengeDuration: big.NewInt(0),
		Nonce:             rand.Int63(),
	}
	response := clientA.CreateVirtualChannel(request)

	waitTimeForCompletedObjectiveIds(t, &clientA, defaultTimeout, response.Id)
	waitTimeForCompletedObjectiveIds(t, &clientB, defaultTimeout, response.Id)
	waitTimeForCompletedObjectiveIds(t, &clientI, defaultTimeout, response.Id)

	return response.ChannelId

}
func TestVirtualFundIntegration(t *testing.T) {

	// Setup logging
	logFile := "virtualfund_client_test.log"
	truncateLog(logFile)
	logDestination := newLogWriter(logFile)

	chain := chainservice.NewMockChain()
	broker := messageservice.NewBroker()

	clientA, _ := setupClient(alice.PrivateKey, chain, broker, logDestination, 0)
	clientB, _ := setupClient(bob.PrivateKey, chain, broker, logDestination, 0)
	clientI, _ := setupClient(irene.PrivateKey, chain, broker, logDestination, 0)

	openAVirtualChannel(t, clientA, clientB, clientI)
}
