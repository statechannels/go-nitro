package integration_test

import (
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/go-cmp/cmp"
	p2pcrypto "github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/client"
	"github.com/statechannels/go-nitro/client/engine"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	p2pms "github.com/statechannels/go-nitro/client/engine/messageservice/p2p-message-service"
	"github.com/statechannels/go-nitro/client/engine/store"
	"github.com/statechannels/go-nitro/client/query"
	"github.com/statechannels/go-nitro/internal/testactors"
	"github.com/statechannels/go-nitro/internal/testdata"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
	"github.com/tidwall/buntdb"
)

func getActorInfo(name testactors.ActorName, tc TestCase) (actor testactors.Actor, participant TestParticipant) {
	switch name {
	case testactors.AliceName:
		actor = testactors.Alice
	case testactors.BobName:
		actor = testactors.Bob
	case testactors.IreneName:
		actor = testactors.Irene
	case testactors.BrianName:
		actor = testactors.Brian
	default:
		panic("Unknown actor")
	}

	found := false
	for _, p := range tc.Participants {
		if p.Name == name {
			participant = p
			found = true
			break
		}
	}
	if !found {
		panic("Unknown participant")
	}

	return
}

func setupMessageService(tc TestCase, actorName testactors.ActorName, si sharedInra) messageservice.MessageService {
	actor, _ := getActorInfo(actorName, tc)
	switch tc.MessageService {
	case TestMessageService:
		return messageservice.NewTestMessageService(actor.Address(), *si.broker, tc.MessageDelay)
	case P2PMessageService:
		ms := p2pms.NewMessageService(
			"127.0.0.1",
			int(actor.Port),
			actor.Address(),
			actor.PrivateKey,
		)

		ms.AddPeers(si.peers)
		return ms
	default:
		panic("Unknown message service")
	}
}

func setupChainService(tc TestCase, actorName testactors.ActorName, si sharedInra) chainservice.ChainService {
	a, _ := getActorInfo(actorName, tc)
	switch tc.Chain {
	case MockChain:
		return chainservice.NewMockChainService(si.mockChain, a.Address())
	case SimulatedChain:
		logDestination := newLogWriter(tc.LogName)

		ethAcountIndex := a.Port - testactors.START_PORT
		cs, err := chainservice.NewSimulatedBackendChainService(*si.simulatedChain, *si.bindings, si.ethAccounts[ethAcountIndex], logDestination)
		if err != nil {
			panic(err)
		}
		return cs
	default:
		panic("Unknown chain service")
	}
}

func setupStore(tc TestCase, actorName testactors.ActorName, si sharedInra) store.Store {
	a, p := getActorInfo(actorName, tc)

	switch p.StoreType {
	case MemStore:
		return store.NewMemStore(a.PrivateKey)
	case DurableStore:
		dataFolder := fmt.Sprintf("%s/%s/%d%d", STORE_TEST_DATA_FOLDER, a.Address().String(), rand.Uint64(), time.Now().UnixNano())
		return store.NewPersistStore(a.PrivateKey, dataFolder, buntdb.Config{})
	default:
		panic(fmt.Sprintf("Unknown store type %s", p.StoreType))
	}
}

func setupIntegrationClient(tc TestCase, actorName testactors.ActorName, si sharedInra) client.Client {
	messageService := setupMessageService(tc, actorName, si)
	cs := setupChainService(tc, actorName, si)
	store := setupStore(tc, actorName, si)
	return client.New(messageService, cs, store, newLogWriter(tc.LogName), &engine.PermissivePolicy{}, nil)
}

func newLogWriter(logFile string) *os.File {
	err := os.MkdirAll("../artifacts", os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	filename := filepath.Join("../artifacts", logFile)
	logDestination, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o666)
	if err != nil {
		log.Fatal(err)
	}

	return logDestination
}

func initialLedgerOutcome(alpha, beta, asset types.Address) outcome.Exit {
	return testdata.Outcomes.Create(alpha, beta, ledgerChannelDeposit, ledgerChannelDeposit, asset)
}

func finalAliceLedger(intermediary, asset types.Address, numPayments, paymentAmount, numChannels uint) outcome.Exit {
	return testdata.Outcomes.Create(
		testactors.Alice.Address(),
		intermediary,
		ledgerChannelDeposit-(numPayments*paymentAmount*numChannels),
		ledgerChannelDeposit+(numPayments*paymentAmount*numChannels),
		asset)
}

func finalBobLedger(intermediary, asset types.Address, numPayments, paymentAmount, numChannels uint) outcome.Exit {
	return testdata.Outcomes.Create(
		intermediary,
		testactors.Bob.Address(),

		ledgerChannelDeposit-(numPayments*paymentAmount*numChannels),
		ledgerChannelDeposit+(numPayments*paymentAmount*numChannels),

		asset)
}

func initialPaymentOutcome(alpha, beta, asset types.Address) outcome.Exit {
	return testdata.Outcomes.Create(alpha, beta, virtualChannelDeposit, 0, asset)
}

func finalPaymentOutcome(alpha, beta, asset types.Address, numPayments, paymentAmount uint) outcome.Exit {
	return testdata.Outcomes.Create(
		alpha,
		beta,
		virtualChannelDeposit-numPayments*paymentAmount,
		numPayments*paymentAmount,
		asset)
}

func setupLedgerChannel(t *testing.T, alpha client.Client, beta client.Client, asset common.Address) types.Destination {
	// Set up an outcome that requires both participants to deposit
	outcome := initialLedgerOutcome(*alpha.Address, *beta.Address, asset)

	response := alpha.CreateLedgerChannel(*beta.Address, 0, outcome)

	<-alpha.ObjectiveCompleteChan(response.Id)
	<-beta.ObjectiveCompleteChan(response.Id)

	return response.ChannelId
}

func closeLedgerChannel(t *testing.T, alpha client.Client, beta client.Client, channelId types.Destination) {
	response := alpha.CloseLedgerChannel(channelId)

	<-alpha.ObjectiveCompleteChan(response)
	<-beta.ObjectiveCompleteChan(response)
}

func waitForObjectives(t *testing.T, a, b client.Client, intermediaries []client.Client, objectiveIds []protocols.ObjectiveId) {
	for _, objectiveId := range objectiveIds {
		<-a.ObjectiveCompleteChan(objectiveId)

		<-b.ObjectiveCompleteChan(objectiveId)

		// TODO: For some reason we're not beeing an ojective complete event from the second intermediary
		<-intermediaries[0].ObjectiveCompleteChan(objectiveId)
		// <-intermediaries[1].ObjectiveCompleteChan(objectiveId)
	}
}

func setupSharedInra(tc TestCase) sharedInra {
	infra := sharedInra{}
	switch tc.Chain {
	case MockChain:
		infra.mockChain = chainservice.NewMockChain()
	case SimulatedChain:
		sim, bindings, ethAccounts, err := chainservice.SetupSimulatedBackend(MAX_PARTICIPANTS)
		if err != nil {
			panic(err)
		}
		infra.simulatedChain = &sim
		infra.bindings = &bindings
		infra.ethAccounts = ethAccounts
	default:
		panic("Unknown chain service")
	}

	switch tc.MessageService {
	case TestMessageService:
		broker := messageservice.NewBroker()
		infra.broker = &broker
	case P2PMessageService:

		infra.peers = make([]p2pms.PeerInfo, len(tc.Participants))
		for i, p := range tc.Participants {

			actor, _ := getActorInfo(p.Name, tc)
			messageKey, err := p2pcrypto.UnmarshalSecp256k1PrivateKey(actor.PrivateKey)
			if err != nil {
				panic(err)
			}
			id, err := peer.IDFromPrivateKey(messageKey)
			if err != nil {
				panic(err)
			}
			infra.peers[i] = p2pms.PeerInfo{
				Port:      int(actor.Port),
				IpAddress: "127.0.0.1",
				Address:   actor.Address(),
				Id:        id,
			}
		}
	}

	return infra
}

// checkPaymentChannel checks that the ledger channel has the expected outcome and status
// It will fail if the channel does not exist
func checkPaymentChannel(t *testing.T, id types.Destination, o outcome.Exit, status query.ChannelStatus, clients ...*client.Client) {
	for _, c := range clients {
		expected := expectedPaymentInfo(id, o, status)
		ledger, err := c.GetPaymentChannel(id)
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(expected, ledger, cmp.AllowUnexported(big.Int{})); diff != "" {
			t.Fatalf("Payment channel diff mismatch (-want +got):\n%s", diff)
		}
	}
}

// expectedLedgerInfo constructs a LedgerChannelInfo so we can easily compare it to the result of GetLedgerChannel
func expectedLedgerInfo(id types.Destination, outcome outcome.Exit, status query.ChannelStatus) query.LedgerChannelInfo {
	clientAdd, _ := outcome[0].Allocations[0].Destination.ToAddress()
	hubAdd, _ := outcome[0].Allocations[1].Destination.ToAddress()

	return query.LedgerChannelInfo{
		ID:     id,
		Status: status,
		Balance: query.LedgerChannelBalance{
			AssetAddress:  types.Address{},
			Hub:           hubAdd,
			Client:        clientAdd,
			ClientBalance: outcome[0].Allocations[0].Amount,
			HubBalance:    outcome[0].Allocations[1].Amount,
		},
	}
}

// checkLedgerChannel checks that the ledger channel has the expected outcome and status
// It will fail if the channel does not exist
func checkLedgerChannel(t *testing.T, ledgerId types.Destination, o outcome.Exit, status query.ChannelStatus, clients ...*client.Client) {
	for _, c := range clients {
		expected := expectedLedgerInfo(ledgerId, o, status)
		ledger, err := c.GetLedgerChannel(ledgerId)
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(expected, ledger, cmp.AllowUnexported(big.Int{})); diff != "" {
			t.Fatalf("Ledger diff mismatch (-want +got):\n%s", diff)
		}
	}
}

// expectedPaymentInfo constructs a LedgerChannelInfo so we can easily compare it to the result of GetPaymentChannel
func expectedPaymentInfo(id types.Destination, outcome outcome.Exit, status query.ChannelStatus) query.PaymentChannelInfo {
	payer, _ := outcome[0].Allocations[0].Destination.ToAddress()
	payee, _ := outcome[0].Allocations[1].Destination.ToAddress()

	return query.PaymentChannelInfo{
		ID:     id,
		Status: status,
		Balance: query.PaymentChannelBalance{
			AssetAddress:   types.Address{},
			Payee:          payee,
			Payer:          payer,
			RemainingFunds: outcome[0].Allocations[0].Amount,
			PaidSoFar:      outcome[0].Allocations[1].Amount,
		},
	}
}
