package client_test

import (
	"math/big"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/statechannels/go-nitro/client"
	"github.com/statechannels/go-nitro/client/engine/chainservice"
	"github.com/statechannels/go-nitro/client/engine/messageservice"
	td "github.com/statechannels/go-nitro/internal/testdata"
	"github.com/statechannels/go-nitro/types"
)

func TestQueryPaymentChannels(t *testing.T) {

	// Setup logging
	logFile := "test_query.log"
	truncateLog(logFile)
	logDestination := newLogWriter(logFile)

	chain := chainservice.NewMockChain()
	chainServiceA := chainservice.NewMockChainService(chain, alice.Address())
	chainServiceB := chainservice.NewMockChainService(chain, bob.Address())
	chainServiceI := chainservice.NewMockChainService(chain, irene.Address())
	broker := messageservice.NewBroker()

	clientA, _ := setupClient(alice.PrivateKey, chainServiceA, broker, logDestination, 0)
	irene, _ := setupClient(irene.PrivateKey, chainServiceI, broker, logDestination, 0)
	clientB, _ := setupClient(bob.PrivateKey, chainServiceB, broker, logDestination, 0)
	directlyFundALedgerChannel(t, clientA, irene, types.Address{})
	directlyFundALedgerChannel(t, clientB, irene, types.Address{})

	id := clientA.CreateVirtualPaymentChannel(
		[]types.Address{*irene.Address},
		bob.Address(),
		0,
		td.Outcomes.Create(
			alice.Address(),
			bob.Address(),
			2,
			0,
			types.Address{},
		)).ChannelId

	res, err := clientA.GetPaymentChannel(id)
	if err != nil {
		t.Fatal(err)
	}

	expected := client.PaymentChannelInfo{
		ID:     id,
		Status: client.Proposed,
		Balance: client.PaymentChannelBalance{
			AssetAddress:   types.Address{},
			Payee:          *clientB.Address,
			Payer:          *clientA.Address,
			PaidSoFar:      big.NewInt(0),
			RemainingFunds: big.NewInt(2),
		}}

	if diff := cmp.Diff(expected, res, cmp.AllowUnexported(big.Int{})); diff != "" {
		t.Fatalf("Query diff mismatch (-want +got):\n%s", diff)
	}

}
