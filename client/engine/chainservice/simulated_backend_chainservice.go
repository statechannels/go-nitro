package chainservice

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	NitroAdjudicator "github.com/statechannels/go-nitro/client/engine/chainservice/adjudicator"
	"github.com/statechannels/go-nitro/protocols"
)

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
