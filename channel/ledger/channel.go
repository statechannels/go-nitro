package ledger

import (
	"fmt"
	"math/big"
	"sort"

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
	guarantees   map[types.Destination]Guarantee
}

func (o LedgerOutcome) Equal(other LedgerOutcome) bool {
	return o.AsOutcome().Equal(other.AsOutcome())
}

// AsOutcome converts a LedgerOutcome to an on-chain exit according to the following convention:
// - the "left" balance is first
// - the "right" balance is second
// - following [left, right] comes the guarantees in sorted order
func (o LedgerOutcome) AsOutcome() outcome.Exit {
	// The first items are [left, right] balances
	allocations := outcome.Allocations{o.left.AsAllocation(), o.right.AsAllocation()}

	// Followed by guarantees, _sorted by the
	keys := make([]types.Destination, 0, len(o.guarantees))
	for k := range o.guarantees {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i].String() < keys[j].String()
	})

	for _, target := range keys {
		allocations = append(allocations, o.guarantees[target].AsAllocation())

	}

	return outcome.Exit{
		outcome.SingleAssetExit{
			Asset:       o.assetAddress,
			Allocations: allocations,
		},
	}
}

var ErrInsufficientFunds = fmt.Errorf("unable to divert to guarantee: insufficient funds")
var ErrDuplicateGuarantee = fmt.Errorf("duplicate guarantee detected")

func (o LedgerOutcome) DivertToGuarantee(g Guarantee) (LedgerOutcome, error) {

	if g.amount.Cmp(&o.left.amount) == 1 {
		return LedgerOutcome{}, ErrInsufficientFunds
	}

	o.left.amount.Sub(&o.left.amount, &g.amount)

	_, found := o.guarantees[g.target]
	if found {
		return LedgerOutcome{}, ErrDuplicateGuarantee
	}
	o.guarantees[g.target] = g

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

var ErrIncorrectTurnNum = fmt.Errorf("incorrect turn number")

// Add updates Vars by adding a guarantee
func (vars Vars) Add(p Add) (Vars, error) {
	if p.turnNum != vars.TurnNum+1 {
		return Vars{}, ErrIncorrectTurnNum
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

