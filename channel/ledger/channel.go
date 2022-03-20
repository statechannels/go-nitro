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

}

type Balance struct {
	destination types.Destination
	amount      big.Int
}

func (b Balance) AsAllocation() outcome.Allocation {
	return outcome.Allocation{Destination: b.destination, Amount: &b.amount, AllocationType: 0}
}

type Guarantee struct {
	amount big.Int
	target types.Destination
	left   types.Destination
	right  types.Destination
}

func (g Guarantee) AsAllocation() outcome.Allocation {
	return outcome.Allocation{
		Destination:    g.target,
		Amount:         &g.amount,
		AllocationType: 1,
		Metadata:       append(g.left.Bytes(), g.right.Bytes()...),
	}
}

type LedgerOutcome struct {
	assetAddress types.Address
	left         Balance
	right        Balance
	guarantees   []Guarantee
}

func (o LedgerOutcome) Equal(other LedgerOutcome) bool {
	return o.AsOutcome().Equal(other.AsOutcome())
}

func (o LedgerOutcome) AsOutcome() outcome.Exit {

	allocations := outcome.Allocations{o.left.AsAllocation(), o.right.AsAllocation()}
	for _, g := range o.guarantees {
		allocations = append(allocations, g.AsAllocation())

	}

	return outcome.Exit{
		outcome.SingleAssetExit{
			Asset:       o.assetAddress,
			Allocations: allocations,
		},
	}
}

func (o LedgerOutcome) DivertToGuarantee(g Guarantee) (LedgerOutcome, error) {

	if g.amount.Cmp(&o.left.amount) == 1 {
		return LedgerOutcome{}, fmt.Errorf("unable to divert to guarantee: insufficient funds")
	}

	o.left.amount.Sub(&o.left.amount, &g.amount)

	// TODO:
	// - ensure sorted
	// - check for duplication
	o.guarantees = append(o.guarantees, g)

	return o, nil
}

// Vars stores the turn number and outcome for a state in a consensus channel
type Vars struct {
	TurnNum uint64
	Outcome LedgerOutcome
}

}

// Add is a proposal to add a guarantee for the given virtual channel
// amount is to be deducted from left
type Add struct {
	turnNum uint64
	Guarantee
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

	o, err := vars.Outcome.DivertToGuarantee(p.Guarantee)

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

