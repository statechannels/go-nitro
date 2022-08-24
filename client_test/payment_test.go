package client_test

import (
	"fmt"
	"math/big"
	"math/rand"
	"testing"

	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/store"
	td "github.com/statechannels/go-nitro/internal/testdata"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/protocols/virtualdefund"
	"github.com/statechannels/go-nitro/protocols/virtualfund"
	"github.com/statechannels/go-nitro/types"
)

func TestPayments(t *testing.T) {

	// Setup logging
	logFile := "test_payments.log"
	truncateLog(logFile)
	logDestination := newLogWriter(logFile)

	chain := chainservice.NewMockChain()
	chainServiceA := chainservice.NewMockChainService(chain, alice.Address())
	chainServiceB := chainservice.NewMockChainService(chain, bob.Address())
	chainServiceI := chainservice.NewMockChainService(chain, irene.Address())

	peers := map[types.Address]string{
		alice.Address(): "localhost:3005",
		bob.Address():   "localhost:3006",
		irene.Address(): "localhost:3007",
	}

	clientA, msgA, storeA := setupClientWithSimpleTCP(alice.PrivateKey, chainServiceA, peers, logDestination, 0)
	clientB, msgB, storeB := setupClientWithSimpleTCP(bob.PrivateKey, chainServiceB, peers, logDestination, 0)
	clientI, msgI, storeI := setupClientWithSimpleTCP(irene.PrivateKey, chainServiceI, peers, logDestination, 0)
	defer msgA.Close()
	defer msgB.Close()
	defer msgI.Close()

	directlyFundALedgerChannel(t, clientA, clientI)
	directlyFundALedgerChannel(t, clientI, clientB)
	outcome := td.Outcomes.Create(alice.Address(), bob.Address(), 100, 100)
	request := virtualfund.ObjectiveRequest{

		CounterParty:      bob.Address(),
		Intermediary:      irene.Address(),
		Outcome:           outcome,
		AppDefinition:     types.Address{},
		AppData:           types.Bytes{},
		ChallengeDuration: big.NewInt(0),
		Nonce:             rand.Int63(),
	}

	firstPaid, secondPaid := int64(1), int64(4)
	totalPaid := firstPaid + secondPaid

	r := clientA.CreateVirtualChannel(request)

	waitTimeForCompletedObjectiveIds(t, &clientA, defaultTimeout, r.Id)
	waitTimeForCompletedObjectiveIds(t, &clientB, defaultTimeout, r.Id)
	waitTimeForCompletedObjectiveIds(t, &clientI, defaultTimeout, r.Id)
	clientA.Pay(r.ChannelId, big.NewInt(firstPaid))
	waitTimeForReceivedVoucher(t, &clientB, defaultTimeout, BasicVoucherInfo{big.NewInt(firstPaid), r.ChannelId})

	// The second voucher adds the first
	clientA.Pay(r.ChannelId, big.NewInt(secondPaid))
	expected := BasicVoucherInfo{big.NewInt(totalPaid), r.ChannelId}
	waitTimeForReceivedVoucher(t, &clientB, defaultTimeout, expected)

	closeId := clientA.CloseVirtualChannel(r.ChannelId)
	waitTimeForCompletedObjectiveIds(t, &clientA, defaultTimeout, closeId)
	waitTimeForCompletedObjectiveIds(t, &clientB, defaultTimeout, closeId)
	waitTimeForCompletedObjectiveIds(t, &clientI, defaultTimeout, closeId)

	for _, clientStore := range []store.Store{storeA, storeB, storeI} {

		oId := protocols.ObjectiveId(fmt.Sprintf("VirtualDefund-%s", r.ChannelId))
		o, err := clientStore.GetObjectiveById(oId)
		if err != nil {
			t.Errorf("Could not get objective: %v", err)
		}
		vdfo := o.(*virtualdefund.Objective)
		if vdfo.GetStatus() != protocols.Completed {
			t.Errorf("Expected objective %s to be completed", vdfo.Id())
		}

		// Check that the ledger outcomes get updated as expected
		switch *clientStore.GetAddress() {
		case alice.Address():
			checkAliceIreneLedgerOutcome(t, vdfo.VId(), vdfo.ToMyRight.ConsensusVars().Outcome, uint(totalPaid))
		case bob.Address():
			checkIreneBobLedgerOutcome(t, vdfo.VId(), vdfo.ToMyLeft.ConsensusVars().Outcome, uint(totalPaid))
		case irene.Address():
			checkAliceIreneLedgerOutcome(t, vdfo.VId(), vdfo.ToMyLeft.ConsensusVars().Outcome, uint(totalPaid))
			checkIreneBobLedgerOutcome(t, vdfo.VId(), vdfo.ToMyRight.ConsensusVars().Outcome, uint(totalPaid))
		}
	}

}
