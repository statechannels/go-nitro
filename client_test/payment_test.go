package client_test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/store"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/protocols/virtualdefund"
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
	numOfVirtualChannels := uint(5)

	cIds := openVirtualChannels(t, clientA, clientB, clientI, numOfVirtualChannels)

	firstPaid, secondPaid := int64(1), int64(4)
	totalPaid := (firstPaid + secondPaid) * int64(numOfVirtualChannels)
	vouchers := make([]BasicVoucherInfo, 0)
	for _, cId := range cIds {
		clientA.Pay(cId, big.NewInt(firstPaid))
		vouchers = append(vouchers, BasicVoucherInfo{big.NewInt(firstPaid), cId})
	}

	for _, cId := range cIds {
		// The second voucher adds the first
		clientA.Pay(cId, big.NewInt(secondPaid))
		vouchers = append(vouchers, BasicVoucherInfo{big.NewInt(firstPaid), cId})
	}
	waitTimeForReceivedVoucher(t, &clientB, defaultTimeout, vouchers...)
	ids := make([]protocols.ObjectiveId, len(cIds))
	for i := 0; i < len(cIds); i++ {
		ids[i] = clientA.CloseVirtualChannel(cIds[i])

	}
	waitTimeForCompletedObjectiveIds(t, &clientA, defaultTimeout, ids...)
	waitTimeForCompletedObjectiveIds(t, &clientB, defaultTimeout, ids...)
	waitTimeForCompletedObjectiveIds(t, &clientI, defaultTimeout, ids...)

	for _, clientStore := range []store.Store{storeA, storeB, storeI} {
		for _, cId := range cIds {
			oId := protocols.ObjectiveId(fmt.Sprintf("VirtualDefund-%s", cId))
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

}
