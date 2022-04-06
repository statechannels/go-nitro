package consensus_channel

import (
	"fmt"
	"math/big"
	"sort"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/crypto"

	"github.com/statechannels/go-nitro/types"
)

type ledgerIndex uint

var ErrIncorrectChannelID = fmt.Errorf("proposal ID and channel ID do not match")

const (
	Leader   ledgerIndex = 0
	Follower ledgerIndex = 1
)

// ConsensusChannel is used to manage states in a running ledger channel
type ConsensusChannel struct {
	// constants
	myIndex ledgerIndex
	fp      state.FixedPart

	Id types.Destination

	// variables
	current       SignedVars       // The "consensus state", signed by both parties
	proposalQueue []SignedProposal // A queue of proposed changes, starting from the consensus state
}

// newConsensusChannel constructs a new consensus channel, validating its input by checking that the signatures are as expected for the given fp, initialTurnNum and outcome]
func newConsensusChannel(
	fp state.FixedPart,
	myIndex ledgerIndex,
	initialTurnNum uint64,
	outcome LedgerOutcome,
	signatures [2]state.Signature,
) (ConsensusChannel, error) {

	cId, err := fp.ChannelId()
	if err != nil {
		return ConsensusChannel{}, err
	}

	vars := Vars{TurnNum: initialTurnNum, Outcome: outcome.clone()}

	leaderAddr, err := vars.AsState(fp).RecoverSigner(signatures[Leader])
	if err != nil {
		return ConsensusChannel{}, fmt.Errorf("could not verify sig: %w", err)
	}
	if leaderAddr != fp.Participants[Leader] {
		return ConsensusChannel{}, fmt.Errorf("leader did not sign initial state: %v, %v", leaderAddr, fp.Participants[Leader])
	}

	followerAddr, err := vars.AsState(fp).RecoverSigner(signatures[Follower])
	if err != nil {
		return ConsensusChannel{}, fmt.Errorf("could not verify sig: %w", err)
	}
	if followerAddr != fp.Participants[Follower] {
		return ConsensusChannel{}, fmt.Errorf("leader did not sign initial state: %v, %v", followerAddr, fp.Participants[Leader])
	}

	current := SignedVars{
		vars,
		signatures,
	}

	return ConsensusChannel{
		fp:            fp,
		Id:            cId,
		myIndex:       myIndex,
		proposalQueue: make([]SignedProposal, 0),
		current:       current,
	}, nil

}

// ConsensusTurnNum returns the turn number of the current consensus state
func (c *ConsensusChannel) ConsensusTurnNum() uint64 {
	return c.current.TurnNum
}

// Includes returns whether or not the consensus state includes the given guarantee
func (c *ConsensusChannel) Includes(g Guarantee) bool {
	return c.current.Outcome.includes(g)
}

// Leader returns the address of the participant responsible for proposing
func (c *ConsensusChannel) Leader() common.Address {
	return c.fp.Participants[Leader]
}
func (c *ConsensusChannel) Accept(p SignedProposal) error {
	panic("UNIMPLEMENTED")
}

// sign constructs a state.State from the given vars, using the ConsensusChannel's constant
// values. It signs the resulting state using sk.
func (c *ConsensusChannel) sign(vars Vars, sk []byte) (state.Signature, error) {
	signer := crypto.GetAddressFromSecretKeyBytes(sk)
	if c.fp.Participants[c.myIndex] != signer {
		return state.Signature{}, fmt.Errorf("attempting to sign from wrong address: %s", signer)
	}

	state := vars.AsState(c.fp)
	return state.Sign(sk)
}

// recoverSigner returns the signer of the vars using the given signature
func (c *ConsensusChannel) recoverSigner(vars Vars, sig state.Signature) (common.Address, error) {
	state := vars.AsState(c.fp)
	return state.RecoverSigner(sig)
}

// ConsensusVars returns the vars of the consensus state
// The consensus state is the latest state that has been signed by both parties
func (c *ConsensusChannel) ConsensusVars() Vars {
	return c.current.Vars
}

// latestProposedVars returns the latest proposed vars in a consensus channel
// by cloning its current vars and applying each proposal in the queue
func (c *ConsensusChannel) latestProposedVars() (Vars, error) {
	vars := Vars{TurnNum: c.current.TurnNum, Outcome: c.current.Outcome.clone()}

	var err error
	for _, p := range c.proposalQueue {
		err = vars.HandleProposal(p.Proposal)
		if err != nil {
			return Vars{}, err
		}
	}

	return vars, nil
}

// validateProposalID checks that the given proposal's ID matches
// the channel's ID
func (c *ConsensusChannel) validateProposalID(propsal Proposal) error {
	if propsal.ChannelID != c.Id {
		return ErrIncorrectChannelID
	}

	return nil
}

// NewBalance returns a new Balance struct with the given amount and destination
func NewBalance(destination types.Destination, amount *big.Int) Balance {
	balanceAmount := big.NewInt(0).Set(amount)
	return Balance{
		destination: destination,
		amount:      balanceAmount,
	}

}

// Balance represents an Allocation of type 0, ie. a simple allocation.
type Balance struct {
	destination types.Destination
	amount      *big.Int
}

// AsAllocation converts a Balance struct into the on-chain outcome.Allocation type
func (b Balance) AsAllocation() outcome.Allocation {
	amount := big.NewInt(0).Set(b.amount)
	return outcome.Allocation{Destination: b.destination, Amount: amount, AllocationType: 0}
}

// Guarantee represents an Allocation of type 1, ie. a guarantee.
type Guarantee struct {
	amount *big.Int
	target types.Destination
	left   types.Destination
	right  types.Destination
}

// NewGuarantee constructs a new guarantee
func NewGuarantee(amount *big.Int, target types.Destination, left types.Destination, right types.Destination) Guarantee {
	return Guarantee{amount, target, left, right}
}

func (g Guarantee) equal(g2 Guarantee) bool {
	if !types.Equal(g.amount, g2.amount) {
		return false
	}
	return g.target == g2.target && g.left == g2.left && g.right == g2.right
}

// AsAllocation converts a Balance struct into the on-chain outcome.Allocation type
func (g Guarantee) AsAllocation() outcome.Allocation {
	amount := big.NewInt(0).Set(g.amount)
	return outcome.Allocation{
		Destination:    g.target,
		Amount:         amount,
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

// NewLedgerOutcome creates a new ledger outcome with the given asset address and balances and guarantees
func NewLedgerOutcome(assetAddress types.Address, left, right Balance, guarantees []Guarantee) *LedgerOutcome {
	guaranteeMap := make(map[types.Destination]Guarantee, len(guarantees))
	for _, g := range guarantees {
		guaranteeMap[g.target] = g
	}
	return &LedgerOutcome{
		assetAddress: assetAddress,
		left:         left,
		right:        right,
		guarantees:   guaranteeMap,
	}
}

// includes returns true when the receiver includes g in its list of guarantees.
func (o *LedgerOutcome) includes(g Guarantee) bool {
	existing, found := o.guarantees[g.target]
	if !found {
		return false
	}

	return g.left == existing.left &&
		g.right == existing.right &&
		g.target == existing.target &&
		types.Equal(existing.amount, g.amount)
}

// FromExit creates a new LedgerOutcome from the given SingleAssetExit.
// It makes some assumptions about the exit:
// - The first alloction entry is for left
// - The second alloction entry is for right
// - We ignore guarantee metadata and just assume that it is [left,right]
func FromExit(sae outcome.SingleAssetExit) LedgerOutcome {

	left := Balance{destination: sae.Allocations[0].Destination, amount: sae.Allocations[0].Amount}
	right := Balance{destination: sae.Allocations[1].Destination, amount: sae.Allocations[1].Amount}
	guarantees := make(map[types.Destination]Guarantee)
	for _, a := range sae.Allocations {

		if a.AllocationType == outcome.GuaranteeAllocationType {
			g := Guarantee{amount: a.Amount,
				target: a.Destination,
				// Instead of decoding the metadata we make an assumption that the metadata has the left/right we expect
				left:  left.destination,
				right: right.destination}
			guarantees[a.Destination] = g
		}

	}
	return LedgerOutcome{left: left, right: right, guarantees: guarantees, assetAddress: sae.Asset}

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
var ErrInsufficientFunds = fmt.Errorf("insufficient funds")
var ErrDuplicateGuarantee = fmt.Errorf("duplicate guarantee detected")
var ErrGuaranteeNotFound = fmt.Errorf("guarantee not found")
var ErrInvalidAmounts = fmt.Errorf("left and right amounts do not add up to the guarantee amount")

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
		amount:      big.NewInt(0).Set(o.left.amount),
	}

	right := Balance{
		destination: o.right.destination,
		amount:      big.NewInt(0).Set(o.right.amount),
	}

	guarantees := make(map[types.Destination]Guarantee)
	for d, g := range o.guarantees {
		g2 := g
		g2.amount = big.NewInt(0).Set(g.amount)
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

// Proposal is a proposal either to add or to remove a guarantee.
//
// Exactly one of {toAdd, toRemove} should be non nil
type Proposal struct {
	ChannelID types.Destination
	ToAdd     Add
	ToRemove  Remove
}

// Updates the turn number on the Add or Remove proposal
func (p *Proposal) SetTurnNum(turnNum uint64) {
	if p.isAddProposal() {
		p.ToAdd.turnNum = turnNum
	} else if p.isRemoveProposal() {
		p.ToRemove.turnNum = turnNum
	}
}

// Returns the turn number on the Add or Remove proposal
func (p *Proposal) TurnNum() uint64 {
	if p.isAddProposal() {
		return p.ToAdd.turnNum
	} else {
		return p.ToRemove.turnNum
	}
}

// equal returns true if the supplied Proposal is deeply equal to the receiver, false otherwise.
func (p *Proposal) equal(q *Proposal) bool {
	return p.ToAdd.equal(q.ToAdd) && p.ToRemove == q.ToRemove
}

// isAddProposal returns true if the proposal contains a non-nil toAdd and a nil toRemove.
func (p *Proposal) isAddProposal() bool {
	return p.ToAdd != Add{} && p.ToRemove == Remove{}
}

// isRemoveProposal returns true if the proposal contains a non-nil toRemove and a nil toAdd.
func (p *Proposal) isRemoveProposal() bool {
	return p.ToAdd == Add{} && p.ToRemove != Remove{}
}

// SignedProposal is a Proposall with a signature on it
type SignedProposal struct {
	state.Signature
	Proposal Proposal
}

// Add is a proposal to add a guarantee for the given virtual channel
// amount is to be deducted from left
type Add struct {
	turnNum uint64
	Guarantee
	LeftDeposit *big.Int
}

// NewAdd constructs a new Add proposal
func NewAdd(turnNum uint64, g Guarantee, leftDeposit *big.Int) Add {
	return Add{
		turnNum:     turnNum,
		Guarantee:   g,
		LeftDeposit: leftDeposit,
	}
}

// NewAddProposal constucts a proposal with a valid Add proposal and empty remove proposal
func NewAddProposal(channelId types.Destination, turnNum uint64, g Guarantee, leftDeposit *big.Int) Proposal {
	return Proposal{ToAdd: NewAdd(turnNum, g, leftDeposit), ChannelID: channelId}
}

// NewRemove constructs a new Remove proposal
func NewRemove(turnNum uint64, target types.Destination, leftAmount, rightAmount *big.Int) Remove {
	return Remove{turnNum: turnNum, Target: target, LeftAmount: leftAmount, RightAmount: rightAmount}
}

// NewRemoveProposal constucts a proposal with a valid Remove proposal and empty Add proposal
func NewRemoveProposal(channelId types.Destination, turnNum uint64, target types.Destination, leftAmount, rightAmount *big.Int) Proposal {
	return Proposal{ToRemove: NewRemove(turnNum, target, leftAmount, rightAmount), ChannelID: channelId}
}

func (a Add) RightDeposit() *big.Int {
	result := big.NewInt(0)
	result.Sub(a.amount, a.LeftDeposit)

	return result
}

func (a Add) equal(a2 Add) bool {
	if a.turnNum != a2.turnNum {
		return false
	}
	if !a.Guarantee.equal(a2.Guarantee) {
		return false
	}
	return types.Equal(a.LeftDeposit, a2.LeftDeposit)
}

var ErrIncorrectTurnNum = fmt.Errorf("incorrect turn number")

// HandleProposal handles a proposal to add or remove a guarantee
// It will mutate Vars by calling Add or Remove for the proposal
func (vars *Vars) HandleProposal(p Proposal) error {
	if p.isAddProposal() {
		return vars.Add(p.ToAdd)

	} else if p.isRemoveProposal() {
		return vars.Remove(p.ToRemove)
	} else {
		return fmt.Errorf("invalid proposal: a proposal must be either an add or a remove proposal")
	}
}

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

	if types.Gt(p.LeftDeposit, p.amount) {
		return ErrInvalidDeposit
	}

	if types.Gt(p.amount, o.left.amount) {
		return ErrInsufficientFunds
	}

	// EFFECTS

	// Increase the turn number
	vars.TurnNum += 1

	// Adjust balances
	o.left.amount.Sub(o.left.amount, p.LeftDeposit)
	rightDeposit := p.RightDeposit()
	o.right.amount.Sub(o.right.amount, rightDeposit)

	// Include guarantee
	o.guarantees[p.target] = p.Guarantee

	return nil
}

// Remove mutates Vars by
// - increasing the turn number by 1
// - removing the guarantee for the Target channel
// - adjusting balances accordingly based on LeftAmount and RightAmount
//
// An error is returned if:
// - the turn number is not incremented
// - a guarantee is not found for the target
// - the amounts are too large for the guarantee amount
//
// If an error is returned, the original vars is not mutated
func (vars *Vars) Remove(p Remove) error {
	// CHECKS

	if p.turnNum != vars.TurnNum+1 {
		return ErrIncorrectTurnNum
	}
	o := vars.Outcome

	guarantee, found := o.guarantees[p.Target]
	if !found {
		return ErrGuaranteeNotFound
	}

	totalRemoved := big.NewInt(0).Add(p.LeftAmount, p.RightAmount)
	if totalRemoved.Cmp(guarantee.amount) != 0 {
		return ErrInvalidAmounts
	}

	// EFFECTS

	// Increase the turn number
	vars.TurnNum += 1

	// Adjust balances
	o.left.amount.Add(o.left.amount, p.LeftAmount)
	o.right.amount.Add(o.right.amount, p.RightAmount)

	// Remove the guarantee
	delete(o.guarantees, p.Target)

	return nil
}

// Remove is a proposal to remover a guarantee for the given virtual channel
type Remove struct {
	turnNum     uint64
	Target      types.Destination
	LeftAmount  *big.Int
	RightAmount *big.Int
}

func (v Vars) AsState(fp state.FixedPart) state.State {
	outcome := v.Outcome.AsOutcome()
	return state.State{
		// Variable
		TurnNum: v.TurnNum,
		Outcome: outcome,

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

// Participants returns the channel participants.
func (c *ConsensusChannel) Participants() []types.Address {
	return c.fp.Participants
}
