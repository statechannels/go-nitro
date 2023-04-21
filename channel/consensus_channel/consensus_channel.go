// Package consensus_channel manages a running ledger channel.
package consensus_channel

import (
	"bytes"
	"fmt"
	"math/big"
	"sort"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/channel"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/crypto"

	"github.com/statechannels/go-nitro/types"
)

type ledgerIndex uint

var (
	ErrIncorrectChannelID = fmt.Errorf("proposal ID and channel ID do not match")
	ErrIncorrectTurnNum   = fmt.Errorf("incorrect turn number")
	ErrInvalidDeposit     = fmt.Errorf("unable to divert to guarantee: invalid deposit")
	ErrInsufficientFunds  = fmt.Errorf("insufficient funds")
	ErrDuplicateGuarantee = fmt.Errorf("duplicate guarantee detected")
	ErrGuaranteeNotFound  = fmt.Errorf("guarantee not found")
	ErrInvalidAmount      = fmt.Errorf("left amount is greater than the guarantee amount")
)

const (
	Leader   ledgerIndex = 0
	Follower ledgerIndex = 1
)

// ConsensusChannel is used to manage states in a running ledger channel.
type ConsensusChannel struct {
	channel.Channel

	// constants

	MyIndex ledgerIndex

	// variables

	// current represents the "consensus state", signed by both parties
	current SignedVars

	// a queue of proposed changes which can be applied to the current state, ordered by TurnNum.
	proposalQueue []SignedProposal
}

// newConsensusChannel constructs a new consensus channel, validating its input by
// checking that the signatures are as expected for the given fp, initialTurnNum and outcome.
func newConsensusChannel(
	s state.State,
	myIndex ledgerIndex,
	initialTurnNum uint64,
	outcome LedgerOutcome,
	signatures [2]state.Signature,
) (ConsensusChannel, error) {
	channel, err := channel.New(s, uint(myIndex))
	if err != nil {
		return ConsensusChannel{}, err
	}
	vars := Vars{TurnNum: initialTurnNum, Outcome: outcome.clone()}

	leaderAddr, err := s.RecoverSigner(signatures[Leader])
	if err != nil {
		return ConsensusChannel{}, fmt.Errorf("could not verify sig: %w", err)
	}
	if leaderAddr != s.Participants[Leader] {
		return ConsensusChannel{}, fmt.Errorf("leader did not sign initial state: %v, %v", leaderAddr, s.Participants[Leader])
	}

	followerAddr, err := s.RecoverSigner(signatures[Follower])
	if err != nil {
		return ConsensusChannel{}, fmt.Errorf("could not verify sig: %w", err)
	}
	if followerAddr != s.Participants[Follower] {
		return ConsensusChannel{}, fmt.Errorf("follower did not sign initial state: %v, %v", followerAddr, s.Participants[Leader])
	}

	current := SignedVars{
		vars,
		signatures,
	}

	return ConsensusChannel{
		Channel:       *channel,
		MyIndex:       myIndex,
		proposalQueue: make([]SignedProposal, 0),
		current:       current,
	}, nil
}

// Receive accepts a proposal signed by the ConsensusChannel counterparty,
// validates its signature, and performs updates to the proposal queue and
// consensus state.
func (c *ConsensusChannel) Receive(sp SignedProposal) error {
	if c.IsFollower() {
		return c.followerReceive(sp)
	}
	if c.IsLeader() {
		return c.leaderReceive(sp)
	}

	return fmt.Errorf("ConsensusChannel is malformed")
}

// IsProposed returns true if a proposal in the queue would lead to g being included in the receiver's outcome, and false otherwise.
//
// Specific clarification: If the current outcome already includes g, IsProposed returns false.
func (c *ConsensusChannel) IsProposed(g Guarantee) (bool, error) {
	latest, err := c.latestProposedVars()
	if err != nil {
		return false, err
	}

	return latest.Outcome.includes(g) && !c.Includes(g), nil
}

// IsProposedNext returns true if the next proposal in the queue would lead to g being included in the receiver's outcome, and false otherwise.
func (c *ConsensusChannel) IsProposedNext(g Guarantee) (bool, error) {
	vars := Vars{TurnNum: c.current.TurnNum, Outcome: c.current.Outcome.clone()}

	if len(c.proposalQueue) == 0 {
		return false, nil
	}

	p := c.proposalQueue[0]
	err := vars.HandleProposal(p.Proposal)
	if vars.TurnNum != p.TurnNum {
		return false, fmt.Errorf("proposal turn number %d does not match vars %d", p.TurnNum, vars.TurnNum)
	}

	if err != nil {
		return false, err
	}

	return vars.Outcome.includes(g) && !c.Includes(g), nil
}

// ConsensusTurnNum returns the turn number of the current consensus state.
func (c *ConsensusChannel) ConsensusTurnNum() uint64 {
	return c.current.TurnNum
}

// Includes returns whether or not the consensus state includes the given guarantee.
func (c *ConsensusChannel) Includes(g Guarantee) bool {
	return c.current.Outcome.includes(g)
}

// IncludesTarget returns whether or not the consensus state includes a guarantee
// addressed to the given target.
func (c *ConsensusChannel) IncludesTarget(target types.Destination) bool {
	return c.current.Outcome.IncludesTarget(target)
}

// HasRemovalBeenProposed returns whether or not a proposal exists to remove the guaranatee for the target.
func (c *ConsensusChannel) HasRemovalBeenProposed(target types.Destination) bool {
	for _, p := range c.proposalQueue {
		if p.Proposal.Type() == RemoveProposal {
			remove := p.Proposal.ToRemove
			if remove.Target == target {
				return true
			}
		}
	}
	return false
}

// HasRemovalBeenProposedNext returns whether or not the next proposal in the queue is a remove proposal for the given target
func (c *ConsensusChannel) HasRemovalBeenProposedNext(target types.Destination) bool {
	if len(c.proposalQueue) == 0 {
		return false
	}

	p := c.proposalQueue[0]
	return p.Proposal.Type() == RemoveProposal && p.Proposal.ToRemove.Target == target
}

// IsLeader returns true if the calling client is the leader of the channel,
// and false otherwise.
func (c *ConsensusChannel) IsLeader() bool {
	return c.MyIndex == Leader
}

// IsFollower returns true if the calling client is the follower of the channel,
// and false otherwise.
func (c *ConsensusChannel) IsFollower() bool {
	return c.MyIndex == Follower
}

// Leader returns the address of the participant responsible for proposing.
func (c *ConsensusChannel) Leader() common.Address {
	return c.FixedPart.Participants[Leader]
}

// Follower returns the address of the participant who receives and contersigns
// proposals.
func (c *ConsensusChannel) Follower() common.Address {
	return c.FixedPart.Participants[Follower]
}

// FundingTargets returns a list of channels funded by the ConsensusChannel
func (c *ConsensusChannel) FundingTargets() []types.Destination {
	return c.current.Outcome.fundingTargets()
}

func (c *ConsensusChannel) Accept(p SignedProposal) error {
	panic("UNIMPLEMENTED")
}

// sign constructs a state.State from the given vars, using the ConsensusChannel's constant
// values. It signs the resulting state using sk.
func (c *ConsensusChannel) sign(vars Vars, sk []byte) (state.Signature, error) {
	signer := crypto.GetAddressFromSecretKeyBytes(sk)
	if c.FixedPart.Participants[c.MyIndex] != signer {
		return state.Signature{}, fmt.Errorf("attempting to sign from wrong address: %s", signer)
	}

	state := vars.AsState(c.FixedPart)
	return state.Sign(sk)
}

// recoverSigner returns the signer of the vars using the given signature.
func (c *ConsensusChannel) recoverSigner(vars Vars, sig state.Signature) (common.Address, error) {
	state := vars.AsState(c.FixedPart)
	return state.RecoverSigner(sig)
}

// ConsensusVars returns the vars of the consensus state
// The consensus state is the latest state that has been signed by both parties.
func (c *ConsensusChannel) ConsensusVars() Vars {
	return c.current.Vars
}

// Signatures returns the signatures on the currently supported state.
func (c *ConsensusChannel) Signatures() [2]state.Signature {
	return c.current.Signatures
}

// ProposalQueue returns the current queue of proposals, ordered by TurnNum.
func (c *ConsensusChannel) ProposalQueue() []SignedProposal {
	// Since c.proposalQueue is already ordered by TurnNum, we can simply return it.
	return c.proposalQueue
}

// latestProposedVars returns the latest proposed vars in a consensus channel
// by cloning its current vars and applying each proposal in the queue.
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
// the channel's ID.
func (c *ConsensusChannel) validateProposalID(propsal Proposal) error {
	if propsal.LedgerID != c.Channel.Id {
		return ErrIncorrectChannelID
	}

	return nil
}

// NewBalance returns a new Balance struct with the given destination and amount.
func NewBalance(destination types.Destination, amount *big.Int) Balance {
	balanceAmount := big.NewInt(0).Set(amount)
	return Balance{
		destination: destination,
		amount:      balanceAmount,
	}
}

// Balance is a convenient, ergonomic representation of a single-asset Allocation
// of type 0, ie. a simple allocation.
type Balance struct {
	destination types.Destination
	amount      *big.Int
}

// Equal returns true if the balances are deeply equal, false otherwise.
func (b Balance) Equal(b2 Balance) bool {
	return bytes.Equal(b.destination.Bytes(), b2.destination.Bytes()) &&
		types.Equal(b.amount, b2.amount)
}

// Clone returns a deep copy of the receiver.
func (b *Balance) Clone() Balance {
	return Balance{
		destination: b.destination,
		amount:      big.NewInt(0).Set(b.amount),
	}
}

// AsAllocation converts a Balance struct into the on-chain outcome.Allocation type.
func (b Balance) AsAllocation() outcome.Allocation {
	amount := big.NewInt(0).Set(b.amount)
	return outcome.Allocation{Destination: b.destination, Amount: amount, AllocationType: outcome.NormalAllocationType}
}

// Guarantee is a convenient, ergonomic representation of a
// single-asset Allocation of type 1, ie. a guarantee.
type Guarantee struct {
	amount *big.Int
	target types.Destination
	left   types.Destination
	right  types.Destination
}

// Clone returns a deep copy of the receiver.
func (g *Guarantee) Clone() Guarantee {
	return Guarantee{
		amount: big.NewInt(0).Set(g.amount),
		target: g.target,
		left:   g.left,
		right:  g.right,
	}
}

// Target returns the target of the guarantee.
func (g Guarantee) Target() types.Destination {
	return g.target
}

// NewGuarantee constructs a new guarantee.
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

// LedgerOutcome encodes the outcome of a ledger channel involving a "leader" and "follower"
// participant.
//
// This struct does not store items in sorted order. The conventional ordering of allocation items is:
// [leader, follower, ...guaranteesSortedbyTargetDestination]
type LedgerOutcome struct {
	assetAddress types.Address // Address of the asset type
	leader       Balance       // Balance of participants[0]
	follower     Balance       // Balance of participants[1]
	guarantees   map[types.Destination]Guarantee
}

// Clone returns a deep copy of the receiver.
func (lo *LedgerOutcome) Clone() LedgerOutcome {
	clonedGuarantees := make(map[types.Destination]Guarantee)
	for key, g := range lo.guarantees {
		clonedGuarantees[key] = g.Clone()
	}
	return LedgerOutcome{
		assetAddress: lo.assetAddress,
		leader:       lo.leader.Clone(),
		follower:     lo.follower.Clone(),
		guarantees:   clonedGuarantees,
	}
}

// Leader returns the leader's balance.
func (lo *LedgerOutcome) Leader() Balance {
	return lo.leader
}

// Follower returns the follower's balance.
func (lo *LedgerOutcome) Follower() Balance {
	return lo.follower
}

// NewLedgerOutcome creates a new ledger outcome with the given asset address, balances, and guarantees.
func NewLedgerOutcome(assetAddress types.Address, leader, follower Balance, guarantees []Guarantee) *LedgerOutcome {
	guaranteeMap := make(map[types.Destination]Guarantee, len(guarantees))
	for _, g := range guarantees {
		guaranteeMap[g.target] = g
	}
	return &LedgerOutcome{
		assetAddress: assetAddress,
		leader:       leader,
		follower:     follower,
		guarantees:   guaranteeMap,
	}
}

// IncludesTarget returns true when the receiver includes a guarantee that targets the given destination.
func (o *LedgerOutcome) IncludesTarget(target types.Destination) bool {
	_, found := o.guarantees[target]
	return found
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
//
// It makes the following assumptions about the exit:
//   - The first alloction entry is for the ledger leader
//   - The second alloction entry is for the ledger follower
//   - All other allocations are guarantees
func FromExit(sae outcome.SingleAssetExit) (LedgerOutcome, error) {
	var (
		leader     = Balance{destination: sae.Allocations[0].Destination, amount: sae.Allocations[0].Amount}
		follower   = Balance{destination: sae.Allocations[1].Destination, amount: sae.Allocations[1].Amount}
		guarantees = make(map[types.Destination]Guarantee)
	)

	for _, a := range sae.Allocations {
		if a.AllocationType == outcome.GuaranteeAllocationType {
			gM, err := outcome.DecodeIntoGuaranteeMetadata(a.Metadata)
			if err != nil {
				return LedgerOutcome{}, fmt.Errorf("failed to decode guarantee metadata: %w", err)
			}

			g := Guarantee{
				amount: a.Amount,
				target: a.Destination,
				left:   gM.Left,
				right:  gM.Right,
			}
			guarantees[a.Destination] = g
		}
	}

	return LedgerOutcome{leader: leader, follower: follower, guarantees: guarantees, assetAddress: sae.Asset}, nil
}

// AsOutcome converts a LedgerOutcome to an on-chain exit according to the following convention:
//   - the "leader" balance is first
//   - the "follower" balance is second
//   - guarantees follow, sorted according to their target destinations
func (o *LedgerOutcome) AsOutcome() outcome.Exit {
	// The first items are [leader, follower] balances
	allocations := outcome.Allocations{o.leader.AsAllocation(), o.follower.AsAllocation()}

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

// fundingTargets returns a list of channels funded by the LedgerOutcome
func (o *LedgerOutcome) fundingTargets() []types.Destination {
	targets := []types.Destination{}

	for dest := range o.guarantees {
		targets = append(targets, dest)
	}

	return targets
}

// Vars stores the turn number and outcome for a state in a consensus channel.
type Vars struct {
	TurnNum uint64
	Outcome LedgerOutcome
}

// Clone returns a deep copy of the receiver.
func (v *Vars) Clone() Vars {
	return Vars{
		v.TurnNum,
		v.Outcome.Clone(),
	}
}

// clone returns a deep clone of v.
func (o *LedgerOutcome) clone() LedgerOutcome {
	assetAddress := o.assetAddress

	leader := Balance{
		destination: o.leader.destination,
		amount:      big.NewInt(0).Set(o.leader.amount),
	}

	follower := Balance{
		destination: o.follower.destination,
		amount:      big.NewInt(0).Set(o.follower.amount),
	}

	guarantees := make(map[types.Destination]Guarantee)
	for d, g := range o.guarantees {
		g2 := g
		g2.amount = big.NewInt(0).Set(g.amount)
		guarantees[d] = g2
	}

	return LedgerOutcome{
		assetAddress: assetAddress,
		leader:       leader,
		follower:     follower,
		guarantees:   guarantees,
	}
}

// SignedVars stores 0-2 signatures for some vars in a consensus channel.
type SignedVars struct {
	Vars
	Signatures [2]state.Signature
}

// clone returns a deep copy of the receiver.
func (sv *SignedVars) clone() SignedVars {
	clonedSignatures := [2]state.Signature{
		sv.Signatures[0],
		sv.Signatures[1],
	}
	return SignedVars{
		sv.Vars.Clone(),
		clonedSignatures,
	}
}

// Proposal is a proposal either to add or to remove a guarantee.
//
// Exactly one of {toAdd, toRemove} should be non nil.
type Proposal struct {
	// LedgerID is the ChannelID of the ConsensusChannel which should receive the proposal.
	//
	// The target virtual channel ID is contained in the Add / Remove struct.
	LedgerID types.Destination
	ToAdd    Add
	ToRemove Remove
}

// Clone returns a deep copy of the receiver.
func (p *Proposal) Clone() Proposal {
	return Proposal{
		p.LedgerID,
		p.ToAdd.Clone(),
		p.ToRemove.Clone(),
	}
}

const (
	AddProposal    ProposalType = "AddProposal"
	RemoveProposal ProposalType = "RemoveProposal"
)

type ProposalType string

// Type returns the type of the proposal based on whether it contains an Add or a Remove proposal.
func (p *Proposal) Type() ProposalType {
	zeroAdd := Add{}
	if p.ToAdd != zeroAdd {
		return AddProposal
	} else {
		return RemoveProposal
	}
}

// Equal returns true if the supplied Proposal is deeply equal to the receiver, false otherwise.
func (p *Proposal) Equal(q *Proposal) bool {
	return p.LedgerID == q.LedgerID && p.ToAdd.equal(q.ToAdd) && p.ToRemove.equal(q.ToRemove)
}

// ChannelID returns the id of the ConsensusChannel which receive the proposal.
func (p SignedProposal) ChannelID() types.Destination {
	return p.Proposal.LedgerID
}

// SortInfo returns the channelId and turn number so the proposal can be easily sorted.
func (p SignedProposal) SortInfo() (types.Destination, uint64) {
	cId := p.Proposal.LedgerID
	turnNum := p.TurnNum
	return cId, turnNum
}

// Target returns the target channel of the proposal.
func (p *Proposal) Target() types.Destination {
	switch p.Type() {
	case "AddProposal":
		{
			return p.ToAdd.Target()
		}
	case "RemoveProposal":
		{
			return p.ToRemove.Target
		}
	default:
		{
			panic("invalid proposal type")
		}
	}
}

// SignedProposal is a Proposal with a signature on it.
type SignedProposal struct {
	state.Signature
	Proposal Proposal
	TurnNum  uint64
}

// Clone returns a deep copy of the receiver.
func (sp *SignedProposal) Clone() SignedProposal {
	sp2 := SignedProposal{sp.Signature, sp.Proposal.Clone(), sp.TurnNum}
	return sp2
}

// Add encodes a proposal to add a guarantee to a ConsensusChannel.
type Add struct {
	Guarantee
	// LeftDeposit is the portion of the Add's amount that will be deducted from left participant's ledger balance.
	//
	// The right participant's deduction is computed as the difference between the guarantee amount and LeftDeposit.
	LeftDeposit *big.Int
}

// Clone returns a deep copy of the receiver.
func (a *Add) Clone() Add {
	if a == nil || a.LeftDeposit == nil {
		return Add{}
	}
	return Add{
		a.Guarantee.Clone(),
		big.NewInt(0).Set(a.LeftDeposit),
	}
}

// NewAdd constructs a new Add proposal.
func NewAdd(g Guarantee, leftDeposit *big.Int) Add {
	return Add{
		Guarantee:   g,
		LeftDeposit: leftDeposit,
	}
}

// NewAddProposal constucts a proposal with a valid Add proposal and empty remove proposal.
func NewAddProposal(ledgerID types.Destination, g Guarantee, leftDeposit *big.Int) Proposal {
	return Proposal{ToAdd: NewAdd(g, leftDeposit), LedgerID: ledgerID}
}

// NewRemove constructs a new Remove proposal.
func NewRemove(target types.Destination, leftAmount *big.Int) Remove {
	return Remove{Target: target, LeftAmount: leftAmount}
}

// NewRemoveProposal constucts a proposal with a valid Remove proposal and empty Add proposal.
func NewRemoveProposal(ledgerID types.Destination, target types.Destination, leftAmount *big.Int) Proposal {
	return Proposal{ToRemove: NewRemove(target, leftAmount), LedgerID: ledgerID}
}

// RightDeposit computes the deposit from the right participant such that
// a.LeftDeposit + a.RightDeposit() fully funds a's guarantee.
func (a Add) RightDeposit() *big.Int {
	result := big.NewInt(0)
	result.Sub(a.amount, a.LeftDeposit)

	return result
}

func (a Add) equal(a2 Add) bool {
	return a.Guarantee.equal(a2.Guarantee) && types.Equal(a.LeftDeposit, a2.LeftDeposit)
}

func (r Remove) equal(r2 Remove) bool {
	return bytes.Equal(r.Target.Bytes(), r2.Target.Bytes()) &&
		types.Equal(r.LeftAmount, r2.LeftAmount)
}

// HandleProposal handles a proposal to add or remove a guarantee.
// It will mutate Vars by calling Add or Remove for the proposal.
func (vars *Vars) HandleProposal(p Proposal) error {
	switch p.Type() {
	case AddProposal:
		{
			return vars.Add(p.ToAdd)
		}
	case RemoveProposal:
		{
			return vars.Remove(p.ToRemove)
		}
	default:
		{
			return fmt.Errorf("invalid proposal: a proposal must be either an add or a remove proposal")
		}
	}
}

// Add mutates Vars by
//   - increasing the turn number by 1
//   - including the guarantee
//   - adjusting balances accordingly
//
// An error is returned if:
//   - the turn number is not incremented
//   - the balances are incorrectly adjusted, or the deposits are too large
//   - the guarantee is already included in vars.Outcome
//
// If an error is returned, the original vars is not mutated.
func (vars *Vars) Add(p Add) error {
	// CHECKS
	o := vars.Outcome

	_, found := o.guarantees[p.target]
	if found {
		return ErrDuplicateGuarantee
	}

	var left, right Balance

	if o.leader.destination == p.Guarantee.left {
		left = o.leader
		right = o.follower
	} else {
		left = o.follower
		right = o.leader
	}

	if types.Gt(p.LeftDeposit, p.amount) {
		return ErrInvalidDeposit
	}

	if types.Gt(p.LeftDeposit, left.amount) {
		return ErrInsufficientFunds
	}

	if types.Gt(p.RightDeposit(), right.amount) {
		return ErrInsufficientFunds
	}

	// EFFECTS

	// Increase the turn number
	vars.TurnNum += 1

	rightDeposit := p.RightDeposit()

	// Adjust balances
	if o.leader.destination == p.Guarantee.left {
		o.leader.amount.Sub(o.leader.amount, p.LeftDeposit)
		o.follower.amount.Sub(o.follower.amount, rightDeposit)
	} else {
		o.follower.amount.Sub(o.follower.amount, p.LeftDeposit)
		o.leader.amount.Sub(o.leader.amount, rightDeposit)
	}

	// Include guarantee
	o.guarantees[p.target] = p.Guarantee

	return nil
}

// Remove mutates Vars by
//   - increasing the turn number by 1
//   - removing the guarantee for the Target channel
//   - adjusting balances accordingly based on LeftAmount and RightAmount
//
// An error is returned if:
//   - the turn number is not incremented
//   - a guarantee is not found for the target
//   - the amounts are too large for the guarantee amount
//
// If an error is returned, the original vars is not mutated.
func (vars *Vars) Remove(p Remove) error {
	// CHECKS

	o := vars.Outcome

	guarantee, found := o.guarantees[p.Target]
	if !found {
		return ErrGuaranteeNotFound
	}

	if p.LeftAmount.Cmp(guarantee.amount) > 0 {
		return ErrInvalidAmount
	}

	// EFFECTS

	// Increase the turn number
	vars.TurnNum += 1

	rightAmount := big.NewInt(0).Sub(guarantee.amount, p.LeftAmount)

	// Adjust balances
	if o.leader.destination == guarantee.left {
		o.leader.amount.Add(o.leader.amount, p.LeftAmount)
		o.follower.amount.Add(o.follower.amount, rightAmount)
	} else {
		o.leader.amount.Add(o.leader.amount, rightAmount)
		o.follower.amount.Add(o.follower.amount, p.LeftAmount)
	}

	// Remove the guarantee
	delete(o.guarantees, p.Target)

	return nil
}

// Remove is a proposal to remove a guarantee for the given virtual channel.
type Remove struct {
	// Target is the address of the virtual channel being defunded
	Target types.Destination
	// LeftAmount is the amount to be credited (in the ledger channel) to the participant specified as the "left" in the guarantee.
	//
	// The amount for the "right" participant is calculated as the difference between the guarantee amount and LeftAmount.
	LeftAmount *big.Int
}

// Clone returns a deep copy of the receiver
func (r *Remove) Clone() Remove {
	if r == nil || r.LeftAmount == nil {
		return Remove{}
	}
	return Remove{
		Target:     r.Target,
		LeftAmount: big.NewInt(0).Set(r.LeftAmount),
	}
}

func (v Vars) AsState(fp state.FixedPart) state.State {
	outcome := v.Outcome.AsOutcome()
	return state.State{
		// Variable
		TurnNum: v.TurnNum,
		Outcome: outcome,

		// Constant
		Participants:      fp.Participants,
		ChannelNonce:      fp.ChannelNonce,
		ChallengeDuration: fp.ChallengeDuration,
		AppData:           types.Bytes{},
		AppDefinition:     fp.AppDefinition,
		IsFinal:           false,
	}
}

// Participants returns the channel participants.
func (c *ConsensusChannel) Participants() []types.Address {
	return c.FixedPart.Participants
}

// Clone returns a deep copy of the receiver.
func (c *ConsensusChannel) Clone() *ConsensusChannel {
	clonedProposalQueue := make([]SignedProposal, len(c.proposalQueue))
	for i, p := range c.proposalQueue {
		clonedProposalQueue[i] = p.Clone()
	}
	clonedChannel := c.Channel.Clone()
	d := ConsensusChannel{
		Channel: *clonedChannel,
		MyIndex: c.MyIndex, current: c.current.clone(), proposalQueue: clonedProposalQueue,
	}
	return &d
}

// SupportedSignedState returns the latest supported signed state.
func (cc *ConsensusChannel) SupportedSignedState() state.SignedState {
	s := cc.ConsensusVars().AsState(cc.FixedPart)
	sigs := cc.current.Signatures
	ss := state.NewSignedState(s)
	_ = ss.AddSignature(sigs[0])
	_ = ss.AddSignature(sigs[1])
	return ss
}
