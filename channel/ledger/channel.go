package ledger

import (
	"fmt"
	"math/big"
	"sort"

	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/types"
)

const proposerIndex = uint(0)

// ConsensusChannel is used to manage states in a running ledger channel
type ConsensusChannel struct {
	// constants
	Id             types.Destination
	MyIndex        uint
	OnChainFunding types.Funds
	state.FixedPart

	// variables
	current     SignedVars       // The "consensus state", signed by both parties
	proposalQueue []SignedProposal // A queue of proposed changes, starting from the consensus state
}

// Balance represents an Allocation of type 0, ie. a simple allocation.
type Balance struct {
	destination types.Destination
	amount      big.Int
}

// AsAllocation converts a Balance struct into the on-chain outcome.Allocation type
func (b Balance) AsAllocation() outcome.Allocation {
	return outcome.Allocation{Destination: b.destination, Amount: &b.amount, AllocationType: 0}
}

// Guarantee represents an Allocation of type 1, ie. a guarantee.
type Guarantee struct {
	amount big.Int
	target types.Destination
	left   types.Destination
	right  types.Destination
}

// AsAllocation converts a Balance struct into the on-chain outcome.Allocation type
func (g Guarantee) AsAllocation() outcome.Allocation {
	return outcome.Allocation{
		Destination:    g.target,
		Amount:         &g.amount,
		AllocationType: 1,
		Metadata:       append(g.left.Bytes(), g.right.Bytes()...),
	}
}

// LedgerOutcome encodes the outcome of a ledger channel involving a "left" and "right"
// participant.
// 
// Allocation items are not stored in sorted order. The conventional ordering of allocation items is:
// [left, right, ...guaranteesSortedbyTargetDestination]
type LedgerOutcome struct {
	assetAddress types.Address                   // Address of the asset type
	left         Balance                         // Balance of participants[0]
	right        Balance                         // Balance of participants[1]
	guarantees   map[types.Destination]Guarantee 
}

// AsOutcome converts a LedgerOutcome to an on-chain exit according to the following convention:
// - the "left" balance is first
// - the "right" balance is second
// - following [left, right] comes the guarantees in sorted order
func (o LedgerOutcome) AsOutcome() outcome.Exit {
	// The first items are [left, right] balances
	allocations := outcome.Allocations{o.left.AsAllocation(), o.right.AsAllocation()}

	// Followed by guarantees, _sorted by the target destination_
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

// DivertToGuarantee deducts g.amount from o.left's balance, and
// adds g to o.guarantees
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

// SignedVars stores 0-2 signatures for some vars in a consensus channel
type SignedVars struct {
	Vars
	Signatures [2]*state.Signature
}

// SignedProposal is a proposal with a signature on it
type SignedProposal struct {
	state.Signature
	Proposal interface{}
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
	TurnNum uint64
	VId     types.Destination
	Paid    types.Funds
}

var ErrIncorrectTurnNum = fmt.Errorf("incorrect turn number")

// Add updates Vars by including a guarantee, updating balances accordingly
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

// Add receives a Guarantee, and generates and stores a SignedProposal in
// the queue, returning the resulting SignedProposal
func (c *ConsensusChannel) Add(g Guarantee, sk []byte) (SignedProposal, error) {
	if c.MyIndex != proposerIndex {
		return SignedProposal{}, fmt.Errorf("only proposer can call Add")
	}

	vars := c.current.Vars
	var err error

	for _, p := range c.proposalQueue {
		vars, err = vars.Add(p.Proposal.(Add))
		if err != nil {
			return SignedProposal{}, err
		}
	}

	latest := c.proposalQueue[len(c.proposalQueue)-1].Proposal.(Add)

	vars, err = vars.Add(Add{Guarantee: g, turnNum: latest.turnNum + 1})
	if err != nil {
		return SignedProposal{}, err
	}

	add := Add{
		turnNum: latest.turnNum + 1,
	}

	signature, err := c.sign(vars, sk)
	if err != nil {
		return SignedProposal{}, fmt.Errorf("unable to sign state update: %f", err)
	}

	signed := SignedProposal{Proposal: add, Signature: signature}

	c.proposalQueue = append(c.proposalQueue, signed)
	return signed, nil
}

// Sign constructs a state.State from the given vars, using the ConsensusChannel's constant
// values. It signs the resulting state using pk.
func (c *ConsensusChannel) sign(vars Vars, pk []byte) (state.Signature, error) {
	fp := c.FixedPart
	state := state.State{
		// Variable
		TurnNum: vars.TurnNum,
		Outcome: vars.Outcome.AsOutcome(),

		// Constant
		ChainId:           fp.ChainId,
		Participants:      fp.Participants,
		ChannelNonce:      fp.ChannelNonce,
		ChallengeDuration: fp.ChallengeDuration,
		AppData:           types.Bytes{},
		AppDefinition:     types.Address{},
		IsFinal:           false,
	}

	return state.Sign(pk)
}

func (c *ConsensusChannel) Accept(p SignedProposal) error {
	panic("UNIMPLEMENTED")
}
