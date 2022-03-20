package ledger

import (
	"fmt"

	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/types"
)

type LedgerChannel struct {
	Id      types.Destination
	MyIndex uint

	OnChainFunding types.Funds

	state.FixedPart

	consensus state.SignedState
}

// Vars stores the turn number and outcome for a state in a consensus channel
type Vars struct {
	TurnNum uint64
	Outcome outcome.Exit
}

// Add is a proposal to add a guarantee for the given virtual channel
// amount is to be deducted from left
type Add struct {
	turnNum uint64
	amount  types.Funds
	vId     types.Destination
	left    types.Destination
	right   types.Destination
}

// Remove is a proposal to remove a guarantee from the ledger channel
// - `paid` funds are added to the right participant's balance
// - the remainder are returned to the left participant's balance
type Remove struct {
	turnNum uint64
	vId     types.Destination
	paid    types.Funds
}

// Add updates Vars by adding a guarantee
func (vars Vars) Add(p Add) (Vars, error) {
	if p.turnNum != vars.TurnNum+1 {
		return Vars{}, fmt.Errorf("incorrect turn number")
	}

	vars.TurnNum += 1

	// TODO: Check for duplicate entries
	o, err := vars.Outcome.DivertToGuarantee(p.left, p.right, p.amount, types.Funds{}, p.vId)

	if err != nil {
		return Vars{}, err
	}

	vars.Outcome = o
	return vars, nil
}

func (vars Vars) Remove(p Remove) (Vars, error) {
	panic("UNIMPLEMENTED")
}

}

