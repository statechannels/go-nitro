package client_test

import (
	"fmt"
	"io"
	"log"
	"math/big"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/go-cmp/cmp"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/client"
	"github.com/statechannels/go-nitro/client/engine"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	p2pms "github.com/statechannels/go-nitro/client/engine/messageservice/p2p-message-service"
	"github.com/statechannels/go-nitro/client/engine/store"
	"github.com/statechannels/go-nitro/client/query"
	"github.com/statechannels/go-nitro/crypto"
	"github.com/statechannels/go-nitro/internal/safesync"
	"github.com/statechannels/go-nitro/internal/testactors"
	"github.com/statechannels/go-nitro/internal/testdata"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
	"github.com/tidwall/buntdb"
)

// waitWithTimeoutForCompletedObjectiveIds waits up to the given timeout for completed objectives and returns when the all objective ids provided have been completed.
// If the timeout lapses and the objectives have not all completed, the parent test will be failed.
func waitTimeForCompletedObjectiveIds(t *testing.T, client *client.Client, timeout time.Duration, ids ...protocols.ObjectiveId) {
	incomplete := safesync.Map[<-chan struct{}]{}

	var wg sync.WaitGroup

	for _, id := range ids {
		incomplete.Store(string(id), client.ObjectiveCompleteChan(id))
		wg.Add(1)
	}

	incomplete.Range(
		func(id string, ch <-chan struct{}) bool {
			go func() {
				<-ch
				incomplete.Delete(string(id))
				wg.Done()
			}()
			return true
		})

	allDone := make(chan struct{})

	go func() {
		wg.Wait()
		allDone <- struct{}{}
	}()

	select {
	case <-time.After(timeout):
		incompleteIds := make([]string, 0)
		incomplete.Range(func(key string, value <-chan struct{}) bool {
			incompleteIds = append(incompleteIds, key)
			return true
		})
		t.Fatalf("Objective ids %s failed to complete on client %s within %s", incompleteIds, client.Address, timeout)
	case <-allDone:
		return
	}
}

// setupClient is a helper function that constructs a client and returns the new client and its store.
func setupClient(pk []byte, chain chainservice.ChainService, msgBroker messageservice.Broker, logDestination io.Writer, meanMessageDelay time.Duration) (client.Client, store.Store) {
	myAddress := crypto.GetAddressFromSecretKeyBytes(pk)
	// TODO: Clean up test data folder?
	dataFolder := fmt.Sprintf("%s/%s/%d", DURABLE_STORE_FOLDER, myAddress.String(), rand.Uint64())
	messageservice := messageservice.NewTestMessageService(myAddress, msgBroker, meanMessageDelay)
	storeA := store.NewDurableStore(pk, dataFolder, buntdb.Config{})
	return client.New(messageservice, chain, storeA, logDestination, &engine.PermissivePolicy{}, nil), storeA
}

func closeClient(t *testing.T, client *client.Client) {
	err := client.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func truncateLog(logFile string) {
	logDestination := newLogWriter(logFile)

	err := logDestination.Truncate(0)
	if err != nil {
		log.Fatal(err)
	}
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

func closeSimulatedChain(t *testing.T, chain chainservice.SimulatedChain) {
	if err := chain.Close(); err != nil {
		t.Fatal(err)
	}
}

func setupMessageService(tc TestCase, tp TestParticipant, si sharedTestInfrastructure) messageservice.MessageService {
	switch tc.MessageService {
	case TestMessageService:
		return messageservice.NewTestMessageService(tp.Address(), *si.broker, tc.MessageDelay)
	case P2PMessageService:
		ms := p2pms.NewMessageService(
			"127.0.0.1",
			int(tp.Port),
			tp.Address(),
			tp.PrivateKey,
		)

		return ms
	default:
		panic("Unknown message service")
	}
}

func setupChainService(tc TestCase, tp TestParticipant, si sharedTestInfrastructure) chainservice.ChainService {
	switch tc.Chain {
	case MockChain:
		return chainservice.NewMockChainService(si.mockChain, tp.Address())
	case SimulatedChain:
		logDestination := newLogWriter(tc.LogName)

		ethAcountIndex := tp.Port - testactors.START_PORT
		cs, err := chainservice.NewSimulatedBackendChainService(si.simulatedChain, *si.bindings, si.ethAccounts[ethAcountIndex], logDestination)
		if err != nil {
			panic(err)
		}
		return cs
	default:
		panic("Unknown chain service")
	}
}

func setupStore(tc TestCase, tp TestParticipant, si sharedTestInfrastructure) store.Store {
	switch tp.StoreType {
	case MemStore:
		return store.NewMemStore(tp.Actor.PrivateKey)
	case DurableStore:
		dataFolder := fmt.Sprintf("%s/%s/%d%d", STORE_TEST_DATA_FOLDER, tp.Address().String(), rand.Uint64(), time.Now().UnixNano())
		return store.NewPersistStore(tp.PrivateKey, dataFolder, buntdb.Config{})
	default:
		panic(fmt.Sprintf("Unknown store type %s", tp.StoreType))
	}
}

func setupIntegrationClient(tc TestCase, tp TestParticipant, si sharedTestInfrastructure) (client.Client, messageservice.MessageService) {
	messageService := setupMessageService(tc, tp, si)
	cs := setupChainService(tc, tp, si)
	store := setupStore(tc, tp, si)
	c := client.New(messageService, cs, store, newLogWriter(tc.LogName), &engine.PermissivePolicy{}, nil)
	return c, messageService
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

		for _, intermediary := range intermediaries {
			intermediary.ObjectiveCompleteChan(objectiveId)
		}
	}
}

func setupSharedInra(tc TestCase) sharedTestInfrastructure {
	infra := sharedTestInfrastructure{}
	switch tc.Chain {
	case MockChain:
		infra.mockChain = chainservice.NewMockChain()
	case SimulatedChain:
		sim, bindings, ethAccounts, err := chainservice.SetupSimulatedBackend(MAX_PARTICIPANTS)
		if err != nil {
			panic(err)
		}
		infra.simulatedChain = sim
		infra.bindings = &bindings
		infra.ethAccounts = ethAccounts
	default:
		panic("Unknown chain service")
	}

	if tc.MessageService == TestMessageService {

		broker := messageservice.NewBroker()
		infra.broker = &broker
	}
	return infra
}

// checkPaymentChannel checks that the ledger channel has the expected outcome and status
// It will fail if the channel does not exist
func checkPaymentChannel(t *testing.T, id types.Destination, o outcome.Exit, status query.ChannelStatus, clients ...client.Client) {
	for _, c := range clients {
		expected := expectedPaymentInfo(id, o, status)
		ledger, err := c.GetPaymentChannel(id)
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(expected, ledger, cmp.AllowUnexported(big.Int{})); diff != "" {
			panic(fmt.Errorf("payment channel diff mismatch (-want +got):\n%s", diff))
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
func checkLedgerChannel(t *testing.T, ledgerId types.Destination, o outcome.Exit, status query.ChannelStatus, clients ...client.Client) {
	for _, c := range clients {
		expected := expectedLedgerInfo(ledgerId, o, status)
		ledger, err := c.GetLedgerChannel(ledgerId)
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(expected, ledger, cmp.AllowUnexported(big.Int{})); diff != "" {
			panic(fmt.Errorf("ledger diff mismatch (-want +got):\n%s", diff))
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

func clientAddresses(clients []client.Client) []common.Address {
	addrs := make([]common.Address, len(clients))
	for i, c := range clients {
		addrs[i] = *c.Address
	}

	return addrs
}

// waitForPeerInfoExchange waits for all the P2PMessageServices to receive peer info from each other
func waitForPeerInfoExchange(numOfPeers int, services ...*p2pms.P2PMessageService) {
	for i := 0; i < numOfPeers; i++ {
		for _, s := range services {
			<-s.PeerInfoReceived()
		}
	}
}
