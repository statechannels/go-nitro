package consensus_channel

import (
	"fmt"
	"math/big"
	"sort"

	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/crypto"
	"github.com/statechannels/go-nitro/types"
)

type ledgerIndex uint

const (
	leader   ledgerIndex = 0
	follower ledgerIndex = 1
)

// consensusChannel is used to manage states in a running ledger channel
type consensusChannel struct {
	// constants
	myIndex ledgerIndex
	fp      state.FixedPart

	// variables
	current       SignedVars       // The "consensus state", signed by both parties
	proposalQueue []SignedProposal // A queue of proposed changes, starting from the consensus state
}

// newConsensusChannel constructs a new consensus channel, validating its input by checking that the signatures are as expected on a prefund setup state
func newConsensusChannel(
	fp state.FixedPart,
	myIndex ledgerIndex,
	outcome LedgerOutcome,
	signatures [2]state.Signature,
) (consensusChannel, error) {
	vars := Vars{TurnNum: 0, Outcome: outcome.clone()}

	leaderAddr, err := vars.asState(fp).RecoverSigner(signatures[leader])
	if err != nil {
		return consensusChannel{}, fmt.Errorf("could not verify sig: %w", err)
	}
	if leaderAddr != fp.Participants[leader] {
		return consensusChannel{}, fmt.Errorf("leader did not sign initial state: %v, %v", leaderAddr, fp.Participants[leader])
	}

	followerAddr, err := vars.asState(fp).RecoverSigner(signatures[follower])
	if err != nil {
		return consensusChannel{}, fmt.Errorf("could not verify sig: %w", err)
	}
	if followerAddr != fp.Participants[follower] {
		return consensusChannel{}, fmt.Errorf("leader did not sign initial state: %v, %v", followerAddr, fp.Participants[leader])
	}

	current := SignedVars{
		vars,
		signatures,
	}

	return consensusChannel{
		fp:            fp,
		myIndex:       myIndex,
		proposalQueue: make([]SignedProposal, 0),
		current:       current,
	}, nil

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

func (g Guarantee) equal(g2 Guarantee) bool {
	if !types.Equal(&g.amount, &g2.amount) {
		return false
	}
	return g.target == g2.target && g.left == g2.left && g.right == g2.right
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
	assetAddress types.Address // Address of the asset type
	left         Balance       // Balance of participants[0]
	right        Balance       // Balance of participants[1]
	guarantees   map[types.Destination]Guarantee
}

func (o *LedgerOutcome) includes(g Guarantee) bool {
	existing, found := o.guarantees[g.target]
	if !found {
		return false
	}

	return g.left == existing.left &&
		g.right == existing.right &&
		g.target == existing.target &&
		types.Equal(&existing.amount, &g.amount)
}

// AsOutcome converts a LedgerOutcome to an on-chain exit according to the following convention:
// - the "left" balance is first
// - the "right" balance is second
// - following [left, right] comes the guarantees in sorted order
func (o *LedgerOutcome) AsOutcome() outcome.Exit {
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

var ErrInvalidDeposit = fmt.Errorf("unable to divert to guarantee: invalid deposit")
var ErrInsufficientFunds = fmt.Errorf("unable to divert to guarantee: insufficient funds")
var ErrDuplicateGuarantee = fmt.Errorf("duplicate guarantee detected")

// Vars stores the turn number and outcome for a state in a consensus channel
type Vars struct {
	TurnNum uint64
	Outcome LedgerOutcome
}

// clone returns a deep clone of v
func (o *LedgerOutcome) clone() LedgerOutcome {
	assetAddress := o.assetAddress

	left := Balance{
		destination: o.left.destination,
		amount:      *big.NewInt(0).Set(&o.left.amount),
	}

	right := Balance{
		destination: o.right.destination,
		amount:      *big.NewInt(0).Set(&o.right.amount),
	}

	guarantees := make(map[types.Destination]Guarantee)
	for d, g := range o.guarantees {
		g2 := g
		g2.amount = *big.NewInt(0).Set(&g.amount)
		guarantees[d] = g2
	}

	return LedgerOutcome{
		assetAddress: assetAddress,
		left:         left,
		right:        right,
		guarantees:   guarantees,
	}
}

// SignedVars stores 0-2 signatures for some vars in a consensus channel
type SignedVars struct {
	Vars
	Signatures [2]state.Signature
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
	LeftDeposit big.Int
}

func (a Add) RightDeposit() big.Int {
	result := big.Int{}
	result.Sub(&a.amount, &a.LeftDeposit)

	return result
}

func (a Add) Equal(a2 Add) bool {
	if a.turnNum != a2.turnNum {
		return false
	}
	if !a.Guarantee.equal(a2.Guarantee) {
		return false
	}
	return types.Equal(&a.LeftDeposit, &a2.LeftDeposit)
}

var ErrIncorrectTurnNum = fmt.Errorf("incorrect turn number")

// Add mutates Vars by
// - increasing the turn number by 1
// - including the guarantee
// - adjusting balances accordingly
//
// An error is returned if:
// - the turn number is not incremented
// - the balances are incorrectly adjusted, or the deposits are too large
// - the guarantee is already included in vars.Outcome
//
// If an error is returned, the original vars is not mutated
func (vars *Vars) Add(p Add) error {
	// CHECKS
	if p.turnNum != vars.TurnNum+1 {
		return ErrIncorrectTurnNum
	}

	o := vars.Outcome

	_, found := o.guarantees[p.target]
	if found {
		return ErrDuplicateGuarantee
	}

	if types.Gt(&p.LeftDeposit, &p.amount) {
		return ErrInvalidDeposit
	}

	if types.Gt(&p.amount, &o.left.amount) {
		return ErrInsufficientFunds
	}

	// EFFECTS

	// Increase the turn number
	vars.TurnNum += 1

	// Adjust balances
	o.left.amount.Sub(&o.left.amount, &p.LeftDeposit)
	rightDeposit := p.RightDeposit()
	o.right.amount.Sub(&o.right.amount, &rightDeposit)

	// Include guarantee
	o.guarantees[p.target] = p.Guarantee

	return nil
}

// latestProposedVars returns the latest proposed vars in a consensus channel
// by cloning its current vars and applying each proposal in the queue
func (c *consensusChannel) latestProposedVars() (Vars, error) {
	vars := Vars{TurnNum: c.current.TurnNum, Outcome: c.current.Outcome.clone()}

	var err error
	for _, p := range c.proposalQueue {
		err = vars.Add(p.Proposal.(Add))
		if err != nil {
			return Vars{}, err
		}
	}

	return vars, nil
}

// sign constructs a state.State from the given vars, using the ConsensusChannel's constant
// values. It signs the resulting state using pk.
func (c *consensusChannel) sign(vars Vars, pk []byte) (state.Signature, error) {
	signer := crypto.GetAddressFromSecretKeyBytes(pk)
	if c.fp.Participants[c.myIndex] != signer {
		return state.Signature{}, fmt.Errorf("attempting to sign from wrong address: %s", signer)
	}

	state := vars.asState(c.fp)
	return state.Sign(pk)
}

func (v Vars) asState(fp state.FixedPart) state.State {
	return state.State{
		// Variable
		TurnNum: v.TurnNum,
		Outcome: v.Outcome.AsOutcome(),

		// Constant
		ChainId:           fp.ChainId,
		Participants:      fp.Participants,
		ChannelNonce:      fp.ChannelNonce,
		ChallengeDuration: fp.ChallengeDuration,
		AppData:           types.Bytes{},
		AppDefinition:     types.Address{},
		IsFinal:           false,
	}
}

func (c *consensusChannel) Accept(p SignedProposal) error {
	panic("UNIMPLEMENTED")
}
