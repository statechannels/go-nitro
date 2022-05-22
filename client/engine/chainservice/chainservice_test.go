package chainservice

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	NitroAdjudicator "github.com/statechannels/go-nitro/client/engine/chainservice/adjudicator"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

func TestChainService(t *testing.T) {
	// Setup transacting EOA
	key, _ := crypto.GenerateKey()
	auth := bind.NewKeyedTransactor(key)
	auth.GasPrice = big.NewInt(10000000000)
	address := auth.From
	balance := new(big.Int)
	balance.SetString("10000000000000000000", 10) // 10 eth in wei

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

	cc := NewChainConnection(na, naAddress, auth, sim)

	// Prepare test data to trigger MockChainService
	testDeposit := types.Funds{
		common.HexToAddress("0x00"): big.NewInt(1),
	}
	testTx := protocols.ChainTransaction{
		ChannelId: types.Destination(common.HexToHash(`4ebd366d014a173765ba1e50f284c179ade31f20441bec41664712aac6cc461d`)),
		Deposit:   testDeposit,
		Type:      protocols.DepositTransactionType,
	}
	cc.in <- testTx

	// Pause to allow the chain service to submit the transaction to the chain.
	time.Sleep(10 * time.Millisecond)
	sim.Commit()

	<-cc.out
}
