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

// AutomineChainService extends EthChainService to automatically mine a block for every transaction
type AutomineChainService struct {
	*EthChainService
	sim transactionProcessor
}

// NewAutomineChainService constructs a chain service that submits transactions to a NitroAdjudicator
// and listens to events from an eventSource
func NewAutomineChainService(sim transactionProcessor, es eventSource, na *NitroAdjudicator.NitroAdjudicator, naAddress common.Address,
	txSigner *bind.TransactOpts) *AutomineChainService {
	return &AutomineChainService{sim: sim, EthChainService: NewEthChainService(na, naAddress, txSigner, es)}
}

// SendTransaction sends the transaction and blocks until it has been submitted.
func (ecs *AutomineChainService) SendTransaction(tx protocols.ChainTransaction) {
	ecs.EthChainService.SendTransaction(tx)
	ecs.sim.Commit()
}
