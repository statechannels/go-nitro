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

// ConsensusChannel is used to manage states in a running ledger channel
type ConsensusChannel struct {
	// constants
	MyIndex ledgerIndex
	state.FixedPart

	// variables
	current       SignedVars       // The "consensus state", signed by both parties
	proposalQueue []SignedProposal // A queue of proposed changes, starting from the consensus state
}

// NewConsensusChannel constructs a new consensus channel, validating its input by checking that the signatures are as expected on a prefund setup state
func NewConsensusChannel(
	fp state.FixedPart,
	myIndex ledgerIndex,
	outcome LedgerOutcome,
	signatures [2]state.Signature,
) (ConsensusChannel, error) {
	vars := Vars{TurnNum: 0, Outcome: outcome}

	leaderAddr, err := vars.asState(fp).RecoverSigner(signatures[leader])
	if err != nil {
		return ConsensusChannel{}, fmt.Errorf("could not verify sig: %w", err)
	}
	if leaderAddr != fp.Participants[leader] {
		return ConsensusChannel{}, fmt.Errorf("leader did not sign initial state: %v, %v", leaderAddr, fp.Participants[leader])
	}

	followerAddr, err := vars.asState(fp).RecoverSigner(signatures[follower])
	if err != nil {
		return ConsensusChannel{}, fmt.Errorf("could not verify sig: %w", err)
	}
	if followerAddr != fp.Participants[follower] {
		return ConsensusChannel{}, fmt.Errorf("leader did not sign initial state: %v, %v", followerAddr, fp.Participants[leader])
	}

	current := SignedVars{
		vars,
		signatures,
	}

	return ConsensusChannel{
		FixedPart:     fp,
		MyIndex:       myIndex,
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

var ErrInvalidDeposit = fmt.Errorf("unable to divert to guarantee: invalid deposit")
var ErrInsufficientFunds = fmt.Errorf("unable to divert to guarantee: insufficient funds")
var ErrDuplicateGuarantee = fmt.Errorf("duplicate guarantee detected")

// DivertToGuarantee deducts g.amount from o.left's balance, and
// adds g to o.guarantees
func (o LedgerOutcome) DivertToGuarantee(p Add) (LedgerOutcome, error) {
	g := p.Guarantee

	if types.Gt(&p.LeftDeposit, &g.amount) {
		return LedgerOutcome{}, ErrInvalidDeposit
	}

	if types.Gt(&g.amount, &o.left.amount) {
		return LedgerOutcome{}, ErrInsufficientFunds
	}

	o.left.amount.Sub(&o.left.amount, &p.LeftDeposit)
	rightDeposit := p.RightDeposit()
	o.right.amount.Sub(&o.right.amount, &rightDeposit)

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

var ErrIncorrectTurnNum = fmt.Errorf("incorrect turn number")

// Add updates Vars by including a guarantee, updating balances accordingly
func (vars Vars) Add(p Add) (Vars, error) {
	if p.turnNum != vars.TurnNum+1 {
		return Vars{}, ErrIncorrectTurnNum
	}

	vars.TurnNum += 1

	o, err := vars.Outcome.DivertToGuarantee(p)

	if err != nil {
		return Vars{}, err
	}

	vars.Outcome = o
	return vars, nil
}

// Propose receives a proposal to add a guarantee, and generates and stores a SignedProposal in
// the queue, returning the resulting SignedProposal
// Note: the TurnNum on add is ignored; the correct turn number is computed by c
func (c *ConsensusChannel) Propose(add Add, sk []byte) (SignedProposal, error) {
	if c.MyIndex != leader {
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

	latestSignedProposal := c.proposalQueue[len(c.proposalQueue)-1]
	latestProposal, ok := latestSignedProposal.Proposal.(Add)
	if !ok {
		return SignedProposal{}, fmt.Errorf("latest proposal is not an Add")
	}

	add.turnNum = latestProposal.turnNum + 1

	vars, err = vars.Add(add)
	if err != nil {
		return SignedProposal{}, fmt.Errorf("propose could not add new state vars: %w", err)
	}

	signature, err := c.sign(vars, sk)
	if err != nil {
		return SignedProposal{}, fmt.Errorf("unable to sign state update: %f", err)
	}

	signed := SignedProposal{Proposal: add, Signature: signature}

	c.proposalQueue = append(c.proposalQueue, signed)
	return signed, nil
}

// sign constructs a state.State from the given vars, using the ConsensusChannel's constant
// values. It signs the resulting state using pk.
func (c *ConsensusChannel) sign(vars Vars, pk []byte) (state.Signature, error) {
	signer := crypto.GetAddressFromSecretKeyBytes(pk)
	if c.Participants[c.MyIndex] != signer {
		return state.Signature{}, fmt.Errorf("attempting to sign from wrong address: %s", signer)
	}

	fp := c.FixedPart
	state := vars.asState(fp)
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

func (c *ConsensusChannel) Accept(p SignedProposal) error {
	panic("UNIMPLEMENTED")
}
