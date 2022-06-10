package chainservice

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/go-cmp/cmp"
	NitroAdjudicator "github.com/statechannels/go-nitro/client/engine/chainservice/adjudicator"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

func TestEthChainService(t *testing.T) {
	// Setup transacting EOA
	key, _ := crypto.GenerateKey()
	auth, _ := bind.NewKeyedTransactorWithChainID(key, big.NewInt(1337)) // 1337 according to docs on SimulatedBackend
	address := auth.From
	balance, _ := new(big.Int).SetString("10000000000000000000", 10) // 10 eth in wei

	// Setup "blockchain"
	gAlloc := map[common.Address]core.GenesisAccount{
		address: {Balance: balance},
	}
	blockGasLimit := uint64(4712388)
	sim := backends.NewSimulatedBackend(gAlloc, blockGasLimit)

	// Deploy Adjudicator
	naAddress, _, na, err := NitroAdjudicator.DeployNitroAdjudicator(auth, sim)
	if err != nil {
		t.Fatal(err)
	}
	sim.Commit()

	acs := NewSimulatedBackendChaneService(sim, sim, na, naAddress, auth)

	// Prepare test data to trigger EthChainService
	testDeposit := types.Funds{
		common.HexToAddress("0x00"): big.NewInt(1),
	}
	channelID := types.Destination(common.HexToHash(`4ebd366d014a173765ba1e50f284c179ade31f20441bec41664712aac6cc461d`))
	testTx := protocols.ChainTransaction{
		ChannelId: channelID,
		Deposit:   testDeposit,
		Type:      protocols.DepositTransactionType,
	}

	out := acs.SubscribeToEvents(address)
	// Submit transactiom
	acs.SendTransaction(testTx)

	// Check that the recieved event matches the expected event
	receivedEvent := <-out
	expectedEvent := DepositedEvent{CommonEvent: CommonEvent{channelID: channelID}, Holdings: testDeposit}
	if diff := cmp.Diff(expectedEvent, receivedEvent, cmp.AllowUnexported(CommonEvent{})); diff != "" {
		t.Fatalf("Clone: mismatch (-want +got):\n%s", diff)
	}

	// Not sure if this is necessary
	sim.Close()
}
