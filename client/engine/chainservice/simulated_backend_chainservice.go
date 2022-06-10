package chainservice

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	NitroAdjudicator "github.com/statechannels/go-nitro/client/engine/chainservice/adjudicator"
	"github.com/statechannels/go-nitro/protocols"
)

var ErrUnableToAssignBigInt = errors.New("simulated_backend_chainservice: unable to assign BigInt")

type transactionProcessor interface {
	Commit()
}

// SimulatedBackendChaneService extends EthChainService to automatically mine a block for every transaction
type SimulatedBackendChaneService struct {
	*EthChainService
	sim transactionProcessor
}

// NewSimulatedBackendChaneService constructs a chain service that submits transactions to a NitroAdjudicator
// and listens to events from an eventSource
func NewSimulatedBackendChaneService(sim transactionProcessor, es eventSource, na *NitroAdjudicator.NitroAdjudicator, naAddress common.Address,
	txSigner *bind.TransactOpts) *SimulatedBackendChaneService {
	return &SimulatedBackendChaneService{sim: sim, EthChainService: NewEthChainService(na, naAddress, txSigner, es)}
}

// SendTransaction sends the transaction and blocks until it has been submitted.
func (ecs *SimulatedBackendChaneService) SendTransaction(tx protocols.ChainTransaction) {
	ecs.EthChainService.SendTransaction(tx)
	ecs.sim.Commit()
}

func SetupSimulatedBackend(numAccounts uint64) (*backends.SimulatedBackend, *NitroAdjudicator.NitroAdjudicator,
	common.Address, []*bind.TransactOpts, error) {
	accounts := make([]*bind.TransactOpts, numAccounts)
	genesisAlloc := make(map[common.Address]core.GenesisAccount)

	// TODO return error
	balance, success := new(big.Int).SetString("10000000000000000000", 10) // 10 eth in wei
	if !success {
		return nil, nil, common.Address{}, accounts, ErrUnableToAssignBigInt
	}

	var err error
	for i := range accounts {
		// Setup transacting EOA
		key, _ := crypto.GenerateKey()
		accounts[i], err = bind.NewKeyedTransactorWithChainID(key, big.NewInt(1337)) // 1337 according to docs on SimulatedBackend
		if err != nil {
			return nil, nil, common.Address{}, accounts, err
		}
		genesisAlloc[accounts[i].From] = core.GenesisAccount{Balance: balance}
	}

	// Setup "blockchain"
	blockGasLimit := uint64(4712388)
	sim := backends.NewSimulatedBackend(genesisAlloc, blockGasLimit)

	// Deploy Adjudicator
	naAddress, _, na, err := NitroAdjudicator.DeployNitroAdjudicator(accounts[0], sim)
	if err != nil {
		return nil, nil, common.Address{}, accounts, err
	}
	sim.Commit()
	return sim, na, naAddress, accounts, nil
}
