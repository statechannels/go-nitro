package node_test

import (
	"fmt"
	"io"
	"log"
	"math/big"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/google/go-cmp/cmp"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/crypto"
	"github.com/statechannels/go-nitro/internal/testactors"
	"github.com/statechannels/go-nitro/internal/testdata"
	"github.com/statechannels/go-nitro/node"
	"github.com/statechannels/go-nitro/node/engine"
	"github.com/statechannels/go-nitro/node/engine/chainservice"
	"github.com/statechannels/go-nitro/node/engine/messageservice"
	p2pms "github.com/statechannels/go-nitro/node/engine/messageservice/p2p-message-service"
	"github.com/statechannels/go-nitro/node/engine/store"
	"github.com/statechannels/go-nitro/node/query"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
	"github.com/tidwall/buntdb"
)

// setupNode is a helper function that constructs a nitro node and returns the new node and its store.
func setupNode(pk []byte, chain chainservice.ChainService, msgBroker messageservice.Broker, logDestination io.Writer, meanMessageDelay time.Duration) (node.Node, store.Store) {
	myAddress := crypto.GetAddressFromSecretKeyBytes(pk)
	// TODO: Clean up test data folder?
	dataFolder := fmt.Sprintf("%s/%s/%d", DURABLE_STORE_FOLDER, myAddress.String(), rand.Uint64())
	messageservice := messageservice.NewTestMessageService(myAddress, msgBroker, meanMessageDelay)
	storeA, err := store.NewDurableStore(pk, dataFolder, buntdb.Config{})
	if err != nil {
		panic(err)
	}
	return node.New(messageservice, chain, storeA, logDestination, &engine.PermissivePolicy{}, nil), storeA
}

func closeNode(t *testing.T, node *node.Node) {
	err := node.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func newLogWriter(logFile string) *os.File {
	err := os.MkdirAll("../artifacts", os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	filename := filepath.Join("../artifacts", logFile)
	// Clear the file
	os.Remove(filename)
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

func setupMessageService(tc TestCase, tp TestParticipant, si sharedTestInfrastructure, logWriter io.Writer) messageservice.MessageService {
	switch tc.MessageService {
	case TestMessageService:
		return messageservice.NewTestMessageService(tp.Address(), *si.broker, tc.MessageDelay)
	case P2PMessageService:
		ms := p2pms.NewMessageService(
			"127.0.0.1",
			int(tp.Port),
			tp.Address(),
			tp.PrivateKey,
			true,
			logWriter,
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

		ethAccountIndex := tp.Port - testactors.START_PORT
		cs, err := chainservice.NewSimulatedBackendChainService(si.simulatedChain, *si.bindings, si.ethAccounts[ethAccountIndex], logDestination)
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
		s, err := store.NewDurableStore(tp.PrivateKey, dataFolder, buntdb.Config{})
		if err != nil {
			panic(err)
		}
		return s
	default:
		panic(fmt.Sprintf("Unknown store type %s", tp.StoreType))
	}
}

func setupIntegrationNode(tc TestCase, tp TestParticipant, si sharedTestInfrastructure) (node.Node, messageservice.MessageService) {
	messageService := setupMessageService(tc, tp, si, newLogWriter(tc.LogName))
	cs := setupChainService(tc, tp, si)
	store := setupStore(tc, tp, si)
	n := node.New(messageService, cs, store, newLogWriter(tc.LogName), &engine.PermissivePolicy{}, nil)
	return n, messageService
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

func setupLedgerChannel(t *testing.T, alpha node.Node, beta node.Node, asset common.Address) types.Destination {
	// Set up an outcome that requires both participants to deposit
	outcome := initialLedgerOutcome(*alpha.Address, *beta.Address, asset)

	response, err := alpha.CreateLedgerChannel(*beta.Address, 0, outcome)
	if err != nil {
		t.Fatal(err)
	}

	<-alpha.ObjectiveCompleteChan(response.Id)
	<-beta.ObjectiveCompleteChan(response.Id)

	return response.ChannelId
}

func closeLedgerChannel(t *testing.T, alpha node.Node, beta node.Node, channelId types.Destination) {
	response, err := alpha.CloseLedgerChannel(channelId)
	if err != nil {
		t.Fatal(err)
	}

	<-alpha.ObjectiveCompleteChan(response)
	<-beta.ObjectiveCompleteChan(response)
}

func waitForObjectives(t *testing.T, a, b node.Node, intermediaries []node.Node, objectiveIds []protocols.ObjectiveId) {
	for _, objectiveId := range objectiveIds {
		<-a.ObjectiveCompleteChan(objectiveId)

		<-b.ObjectiveCompleteChan(objectiveId)

		for _, intermediary := range intermediaries {
			<-intermediary.ObjectiveCompleteChan(objectiveId)
		}
	}
}

func setupSharedInfra(tc TestCase) sharedTestInfrastructure {
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
func checkPaymentChannel(t *testing.T, id types.Destination, o outcome.Exit, status query.ChannelStatus, clients ...node.Node) {
	for _, c := range clients {
		expected := createPaychInfo(id, o, status)
		ledger, err := c.GetPaymentChannel(id)
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(expected, ledger, cmp.AllowUnexported(big.Int{})); diff != "" {
			panic(fmt.Errorf("payment channel diff mismatch (-want +got):\n%s", diff))
		}
	}
}

// createLedgerInfo constructs a LedgerChannelInfo so we can easily compare it to the result of GetLedgerChannel
func createLedgerInfo(id types.Destination, outcome outcome.Exit, status query.ChannelStatus) query.LedgerChannelInfo {
	clientAdd, _ := outcome[0].Allocations[0].Destination.ToAddress()
	hubAdd, _ := outcome[0].Allocations[1].Destination.ToAddress()

	return query.LedgerChannelInfo{
		ID:     id,
		Status: status,
		Balance: query.LedgerChannelBalance{
			AssetAddress:  types.Address{},
			Hub:           hubAdd,
			Client:        clientAdd,
			ClientBalance: (*hexutil.Big)(outcome[0].Allocations[0].Amount),
			HubBalance:    (*hexutil.Big)(outcome[0].Allocations[1].Amount),
		},
	}
}

type channelStatusShorthand struct {
	clientA uint
	clientB uint
	status  query.ChannelStatus
}

// createLedgerStory returns a sequence of LedgerChannelInfo structs according
// to the supplied states.
func createLedgerStory(
	id types.Destination,
	clientAddr, hubAddr common.Address,
	states []channelStatusShorthand,
) []query.LedgerChannelInfo {
	story := make([]query.LedgerChannelInfo, len(states))
	for i, state := range states {
		story[i] = createLedgerInfo(
			id,
			simpleOutcome(clientAddr, hubAddr, state.clientA, state.clientB),
			state.status,
		)
	}

	return story
}

// checkLedgerChannel checks that the ledger channel has the expected outcome and status
// It will fail if the channel does not exist
func checkLedgerChannel(t *testing.T, ledgerId types.Destination, o outcome.Exit, status query.ChannelStatus, clients ...node.Node) {
	for _, c := range clients {
		expected := createLedgerInfo(ledgerId, o, status)
		ledger, err := c.GetLedgerChannel(ledgerId)
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(expected, ledger, cmp.AllowUnexported(big.Int{})); diff != "" {
			panic(fmt.Errorf("ledger diff mismatch (-want +got):\n%s", diff))
		}
	}
}

// createPaychInfo constructs a PaymentChannelInfo so we can easily compare it to the result of GetPaymentChannel
func createPaychInfo(id types.Destination, outcome outcome.Exit, status query.ChannelStatus) query.PaymentChannelInfo {
	payer, _ := outcome[0].Allocations[0].Destination.ToAddress()
	payee, _ := outcome[0].Allocations[1].Destination.ToAddress()

	return query.PaymentChannelInfo{
		ID:     id,
		Status: status,
		Balance: query.PaymentChannelBalance{
			AssetAddress:   types.Address{},
			Payee:          payee,
			Payer:          payer,
			RemainingFunds: (*hexutil.Big)(outcome[0].Allocations[0].Amount),
			PaidSoFar:      (*hexutil.Big)(outcome[0].Allocations[1].Amount),
		},
	}
}

// createPaychStory returns a sequence of PaymentChannelInfo structs according
// to the supplied states.
func createPaychStory(
	id types.Destination,
	payerAddr, payeeAddr common.Address,
	states []channelStatusShorthand,
) []query.PaymentChannelInfo {
	story := make([]query.PaymentChannelInfo, len(states))
	for i, state := range states {
		story[i] = createPaychInfo(
			id,
			simpleOutcome(payerAddr, payeeAddr, state.clientA, state.clientB),
			state.status,
		)
	}
	return story
}

func clientAddresses(clients []node.Node) []common.Address {
	addrs := make([]common.Address, len(clients))
	for i, c := range clients {
		addrs[i] = *c.Address
	}

	return addrs
}

// waitForPeerInfoExchange waits for all the P2PMessageServices to receive peer info from each other
func waitForPeerInfoExchange(services ...*p2pms.P2PMessageService) {
	for _, s := range services {
		for i := 0; i < len(services)-1; i++ {
			<-s.PeerInfoReceived()
		}
	}
}
