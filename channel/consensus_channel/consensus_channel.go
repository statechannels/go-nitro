// Package consensus_channel manages a running ledger channel.
package consensus_channel

import (
	"bytes"
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

var (
	ErrIncorrectChannelID = fmt.Errorf("proposal ID and channel ID do not match")
	ErrIncorrectTurnNum   = fmt.Errorf("incorrect turn number")
	ErrInvalidDeposit     = fmt.Errorf("unable to divert to guarantee: invalid deposit")
	ErrInsufficientFunds  = fmt.Errorf("insufficient funds")
	ErrDuplicateGuarantee = fmt.Errorf("duplicate guarantee detected")
	ErrGuaranteeNotFound  = fmt.Errorf("guarantee not found")
	ErrInvalidAmounts     = fmt.Errorf("left and right amounts do not add up to the guarantee amount")
)

const (
	Leader   ledgerIndex = 0
	Follower ledgerIndex = 1
)

// ConsensusChannel is used to manage states in a running ledger channel.
type ConsensusChannel struct {
	// constants

	MyIndex        ledgerIndex
	fp             state.FixedPart
	Id             types.Destination
	OnChainFunding types.Funds

	// variables

	// current represents the "consensus state", signed by both parties
	current SignedVars

	// a queue of proposed changes which can be applied to the current state
	proposalQueue []SignedProposal
}

// newConsensusChannel constructs a new consensus channel, validating its input by
// checking that the signatures are as expected for the given fp, initialTurnNum and outcome.
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
		return ConsensusChannel{}, fmt.Errorf("follower did not sign initial state: %v, %v", followerAddr, fp.Participants[Leader])
	}

	current := SignedVars{
		vars,
		signatures,
	}

	return ConsensusChannel{
		fp:            fp,
		Id:            cId,
		MyIndex:       myIndex,
		proposalQueue: make([]SignedProposal, 0),
		current:       current,
	}, nil

}

// FixedPart returns the fixed part of the channel.
func (c *ConsensusChannel) FixedPart() state.FixedPart {
	return c.fp
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
	return c.current.Outcome.includesTarget(target)
}

// HasRemovalBeenProposedFor returns whether or not a proposal exists to remove the guaranatee for the target.
func (c *ConsensusChannel) HasRemovalBeenProposedFor(target types.Destination) bool {
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
	return c.fp.Participants[Leader]
}

// Follower returns the address of the participant who recieves and contersigns
// proposals.
func (c *ConsensusChannel) Follower() common.Address {
	return c.fp.Participants[Follower]
}

func (c *ConsensusChannel) Accept(p SignedProposal) error {
	panic("UNIMPLEMENTED")
}

// sign constructs a state.State from the given vars, using the ConsensusChannel's constant
// values. It signs the resulting state using sk.
func (c *ConsensusChannel) sign(vars Vars, sk []byte) (state.Signature, error) {
	signer := crypto.GetAddressFromSecretKeyBytes(sk)
	if c.fp.Participants[c.MyIndex] != signer {
		return state.Signature{}, fmt.Errorf("attempting to sign from wrong address: %s", signer)
	}

	state := vars.AsState(c.fp)
	return state.Sign(sk)
}

// recoverSigner returns the signer of the vars using the given signature.
func (c *ConsensusChannel) recoverSigner(vars Vars, sig state.Signature) (common.Address, error) {
	state := vars.AsState(c.fp)
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

// ProposalQueue returns the current queue of proposals.
func (c *ConsensusChannel) ProposalQueue() []SignedProposal {
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
	if propsal.ChannelID != c.Id {
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

// Clone returns a deep copy of the receiver.
func (lo *LedgerOutcome) Clone() LedgerOutcome {
	clonedGuarantees := make(map[types.Destination]Guarantee)
	for key, g := range lo.guarantees {
		clonedGuarantees[key] = g.Clone()
	}
	return LedgerOutcome{
		assetAddress: lo.assetAddress,
		left:         lo.left.Clone(),
		right:        lo.right.Clone(),
		guarantees:   clonedGuarantees,
	}
}

// NewLedgerOutcome creates a new ledger outcome with the given asset address, balances, and guarantees.
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

// includesTarget returns true when the receiver includes a guarantee that targets the given destination.
func (o *LedgerOutcome) includesTarget(target types.Destination) bool {
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
// It makes some assumptions about the exit:
//  - The first alloction entry is for left
//  - The second alloction entry is for right
func FromExit(sae outcome.SingleAssetExit) LedgerOutcome {
	left := Balance{destination: sae.Allocations[0].Destination, amount: sae.Allocations[0].Amount}
	right := Balance{destination: sae.Allocations[1].Destination, amount: sae.Allocations[1].Amount}
	guarantees := make(map[types.Destination]Guarantee)
	for _, a := range sae.Allocations {

		if a.AllocationType == outcome.GuaranteeAllocationType {
			gM, err := outcome.DecodeIntoGuaranteeMetadata(a.Metadata)
			if err != nil {
				panic(err)
			}
			g := Guarantee{amount: a.Amount,
				target: a.Destination,
				// Instead of decoding the metadata we make an assumption that the metadata has the left/right we expect
				left:  gM.Left,
				right: gM.Right}
			guarantees[a.Destination] = g
		}

	}
	return LedgerOutcome{left: left, right: right, guarantees: guarantees, assetAddress: sae.Asset}

}

// AsOutcome converts a LedgerOutcome to an on-chain exit according to the following convention:
//  - the "left" balance is first
//  - the "right" balance is second
//  - following [left, right] comes the guarantees sorted according to their target destination
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

// SignedVars stores 0-2 signatures for some vars in a consensus channel.
type SignedVars struct {
	Vars
	Signatures [2]state.Signature
}

// clone returns a deep copy of the reciever.
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
	// ChannelID of the ConsensusChannel which should recieve the proposal.
	//
	// The target virtual channel ID is contained in the Add / Remove struct.
	ChannelID types.Destination
	ToAdd     Add
	ToRemove  Remove
}

// Clone returns a deep copy of the receiver.
func (p *Proposal) Clone() Proposal {
	return Proposal{
		p.ChannelID,
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
	return p.ChannelID == q.ChannelID && p.ToAdd.equal(q.ToAdd) && p.ToRemove.equal(q.ToRemove)
}

// ChannelID returns the id of the ConsensusChannel which receive the proposal.
func (p SignedProposal) ChannelID() types.Destination {
	return p.Proposal.ChannelID
}

// SortInfo returns the channelId and turn number so the proposal can be easily sorted.
func (p SignedProposal) SortInfo() (types.Destination, uint64) {
	cId := p.Proposal.ChannelID
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

// Clone returns a deep copy of the reciever.
func (sp *SignedProposal) Clone() SignedProposal {
	sp2 := SignedProposal{sp.Signature, sp.Proposal.Clone(), sp.TurnNum}
	return sp2
}

// Add encodes a proposal to add a guarantee to a ConsensusChannel.
//
// The balance of the guarantee is to be deducted from left.
type Add struct {
	Guarantee
	LeftDeposit *big.Int
}

// Clone returns a deep copy of the receiver.
func (a *Add) Clone() Add {
	if a == nil {
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
func NewAddProposal(channelId types.Destination, g Guarantee, leftDeposit *big.Int) Proposal {
	return Proposal{ToAdd: NewAdd(g, leftDeposit), ChannelID: channelId}
}

// NewRemove constructs a new Remove proposal.
func NewRemove(target types.Destination, leftAmount, rightAmount *big.Int) Remove {
	return Remove{Target: target, LeftAmount: leftAmount, RightAmount: rightAmount}
}

// NewRemoveProposal constucts a proposal with a valid Remove proposal and empty Add proposal.
func NewRemoveProposal(channelId types.Destination, target types.Destination, leftAmount, rightAmount *big.Int) Proposal {
	return Proposal{ToRemove: NewRemove(target, leftAmount, rightAmount), ChannelID: channelId}
}

// RightDeposit computes the deposit from the right participant such that
// a.LeftDeposit + a.RightDeposit() fully funds a's guarantee.
func (a Add) RightDeposit() *big.Int {
	result := big.NewInt(0)
	result.Sub(a.amount, a.LeftDeposit)

	return result
}

func (a Add) equal(a2 Add) bool {
	if !a.Guarantee.equal(a2.Guarantee) {
		return false
	}
	return types.Equal(a.LeftDeposit, a2.LeftDeposit)
}

func (r Remove) equal(r2 Remove) bool {
	if !bytes.Equal(r.Target.Bytes(), r2.Target.Bytes()) {

		return false
	}
	if !types.Equal(r.LeftAmount, r2.LeftAmount) {
		return false
	}
	if !types.Equal(r.RightAmount, r2.RightAmount) {
		return false
	}
	return true
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
//  - increasing the turn number by 1
//  - including the guarantee
//  - adjusting balances accordingly
//
// An error is returned if:
//  - the turn number is not incremented
//  - the balances are incorrectly adjusted, or the deposits are too large
//  - the guarantee is already included in vars.Outcome
//
// If an error is returned, the original vars is not mutated.
func (vars *Vars) Add(p Add) error {
	// CHECKS
	o := vars.Outcome

	_, found := o.guarantees[p.target]
	if found {
		return ErrDuplicateGuarantee
	}

	if types.Gt(p.LeftDeposit, p.amount) {
		return ErrInvalidDeposit
	}

	if types.Gt(p.LeftDeposit, o.left.amount) {
		return ErrInsufficientFunds
	}

	if types.Gt(p.RightDeposit(), o.right.amount) {
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
//  - increasing the turn number by 1
//  - removing the guarantee for the Target channel
//  - adjusting balances accordingly based on LeftAmount and RightAmount
//
// An error is returned if:
//  - the turn number is not incremented
//  - a guarantee is not found for the target
//  - the amounts are too large for the guarantee amount
//
// If an error is returned, the original vars is not mutated.
func (vars *Vars) Remove(p Remove) error {
	// CHECKS

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

// Remove is a proposal to remover a guarantee for the given virtual channel.
type Remove struct {
	Target types.Destination
	// LeftAmount is the amount to be credited to the left participant of the two party ledger channel
	LeftAmount *big.Int
	// RightAmount is the amount to be credited to the right participant of the two party ledger channel
	RightAmount *big.Int
}

// Clone returns a deep copy of the receiver
func (r *Remove) Clone() Remove {
	if r == nil || r.LeftAmount == nil || r.RightAmount == nil {
		return Remove{}
	}
	return Remove{
		Target:      r.Target,
		LeftAmount:  big.NewInt(0).Set(r.LeftAmount),
		RightAmount: big.NewInt(0).Set(r.RightAmount),
	}
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
		AppDefinition:     fp.AppDefinition,
		IsFinal:           false,
	}
}

// Participants returns the channel participants.
func (c *ConsensusChannel) Participants() []types.Address {
	return c.fp.Participants
}

// Clone returns a deep copy of the receiver.
func (c *ConsensusChannel) Clone() *ConsensusChannel {

	clonedProposalQueue := make([]SignedProposal, len(c.proposalQueue))
	for i, p := range c.proposalQueue {
		clonedProposalQueue[i] = p.Clone()
	}
	d := ConsensusChannel{c.MyIndex, c.fp.Clone(), c.Id, c.OnChainFunding.Clone(), c.current.clone(), clonedProposalQueue}
	return &d
}

// SupportedSignedState returns the latest supported signed state.
func (cc *ConsensusChannel) SupportedSignedState() state.SignedState {
	s := cc.ConsensusVars().AsState(cc.fp)
	sigs := cc.current.Signatures
	ss := state.NewSignedState(s)
	_ = ss.AddSignature(sigs[0])
	_ = ss.AddSignature(sigs[1])
	return ss
}
