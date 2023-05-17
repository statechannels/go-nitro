package chainservice

import (
	"context"
	"errors"
	"io"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	NitroAdjudicator "github.com/statechannels/go-nitro/client/engine/chainservice/adjudicator"
	ConsensusApp "github.com/statechannels/go-nitro/client/engine/chainservice/consensusapp"
	Token "github.com/statechannels/go-nitro/client/engine/chainservice/erc20"
	VirtualPaymentApp "github.com/statechannels/go-nitro/client/engine/chainservice/virtualpaymentapp"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

// This is the chain id used by the simulated backend as well as hardhat
const TEST_CHAIN_ID = 1337

var ErrUnableToAssignBigInt = errors.New("simulated_backend_chainservice: unable to assign BigInt")

type binding[T any] struct {
	Address  common.Address
	Contract *T
}

type Bindings struct {
	Adjudicator       binding[NitroAdjudicator.NitroAdjudicator]
	Token             binding[Token.Token]
	ConsensusApp      binding[ConsensusApp.ConsensusApp]
	VirtualPaymentApp binding[VirtualPaymentApp.VirtualPaymentApp]
}

type SimulatedChain interface {
	ethChain
	Commit() common.Hash
	Close() error
}

// This is used to wrap the simulated backend so that we can provide a ChainID function like a real eth client
type BackendWrapper struct {
	*backends.SimulatedBackend
}

func (b *BackendWrapper) ChainID(ctx context.Context) (*big.Int, error) {
	return big.NewInt(TEST_CHAIN_ID), nil
}

// SimulatedBackendChainService extends EthChainService to automatically mine a block for every transaction
type SimulatedBackendChainService struct {
	*EthChainService
	sim SimulatedChain
}

// NewSimulatedBackendChainService constructs a chain service that submits transactions to a NitroAdjudicator
// and listens to events from an eventSource
func NewSimulatedBackendChainService(sim SimulatedChain, bindings Bindings,
	txSigner *bind.TransactOpts, logDestination io.Writer,
) (ChainService, error) {
	ethChainService, err := newEthChainService(sim,
		bindings.Adjudicator.Contract,
		bindings.Adjudicator.Address,
		bindings.ConsensusApp.Address,
		bindings.VirtualPaymentApp.Address,
		txSigner,
		logDestination)
	if err != nil {
		return &SimulatedBackendChainService{}, err
	}

	return &SimulatedBackendChainService{sim: sim, EthChainService: ethChainService}, nil
}

// SendTransaction sends the transaction and blocks until it has been mined.
func (sbcs *SimulatedBackendChainService) SendTransaction(tx protocols.ChainTransaction) error {
	err := sbcs.EthChainService.SendTransaction(tx)
	if err != nil {
		return err
	}
	sbcs.sim.Commit()
	return nil
}

// SetupSimulatedBackend creates a new SimulatedBackend with the supplied number of transacting accounts, deploys the Nitro Adjudicator and returns both.
func SetupSimulatedBackend(numAccounts uint64) (SimulatedChain, Bindings, []*bind.TransactOpts, error) {
	accounts := make([]*bind.TransactOpts, numAccounts)
	genesisAlloc := make(map[common.Address]core.GenesisAccount)
	contractBindings := Bindings{}

	balance, success := new(big.Int).SetString("10000000000000000000", 10) // 10 eth in wei
	if !success {
		return nil, contractBindings, accounts, ErrUnableToAssignBigInt
	}

	var err error
	for i := range accounts {
		// Setup transacting EOA
		key, _ := crypto.GenerateKey()
		accounts[i], err = bind.NewKeyedTransactorWithChainID(key, big.NewInt(1337)) // 1337 according to docs on SimulatedBackend
		if err != nil {
			return nil, contractBindings, accounts, err
		}
		genesisAlloc[accounts[i].From] = core.GenesisAccount{Balance: balance}
	}

	// Setup "blockchain"
	blockGasLimit := uint64(15_000_000)
	sim := backends.NewSimulatedBackend(genesisAlloc, blockGasLimit)

	// Deploy Adjudicator
	naAddress, _, na, err := NitroAdjudicator.DeployNitroAdjudicator(accounts[0], sim)
	if err != nil {
		return nil, contractBindings, accounts, err
	}

	// Deploy ConsensusApp
	consensusAppAddress, _, ca, err := ConsensusApp.DeployConsensusApp(accounts[0], sim)
	if err != nil {
		return nil, contractBindings, accounts, err
	}

	// Deploy VirtualPaymentChannelApp
	virtualPaymentAppAddress, _, vpa, err := VirtualPaymentApp.DeployVirtualPaymentApp(accounts[0], sim)
	if err != nil {
		return nil, contractBindings, accounts, err
	}

	// Deploy a test ERC20 Token Contract
	tokenAddress, _, tokenBinding, err := Token.DeployToken(accounts[0], sim, accounts[0].From)
	if err != nil {
		return nil, contractBindings, accounts, err
	}

	// Distributed tokens to all accounts
	INITIAL_TOKEN_BALANCE := big.NewInt(10_000_000)
	for _, account := range accounts {
		accountAddress := account.From
		_, err := tokenBinding.Transfer(accounts[0], accountAddress, INITIAL_TOKEN_BALANCE)
		if err != nil {
			return nil, contractBindings, accounts, err
		}
	}

	contractBindings = Bindings{
		Adjudicator:       binding[NitroAdjudicator.NitroAdjudicator]{naAddress, na},
		Token:             binding[Token.Token]{tokenAddress, tokenBinding},
		ConsensusApp:      binding[ConsensusApp.ConsensusApp]{consensusAppAddress, ca},
		VirtualPaymentApp: binding[VirtualPaymentApp.VirtualPaymentApp]{virtualPaymentAppAddress, vpa},
	}
	sim.Commit()

	return &BackendWrapper{sim}, contractBindings, accounts, nil
}

func (sbcs *SimulatedBackendChainService) GetConsensusAppAddress() types.Address {
	return sbcs.consensusAppAddress
}

// GetVirtualPaymentAppAddress returns the address of a deployed VirtualPaymentApp
func (sbcs *SimulatedBackendChainService) GetVirtualPaymentAppAddress() types.Address {
	return sbcs.virtualPaymentAppAddress
}
