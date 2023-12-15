package consensus_channel

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/types"
)

// jsonAdd replaces Add's private fields with public ones,
// making it suitable for serialization
// embedded structs are moved to name fields for easier serialization
type jsonAdd struct {
	Guarantee   Guarantee
	LeftDeposit *big.Int
}

// MarshalJSON returns a JSON representation of the Add
func (a Add) MarshalJSON() ([]byte, error) {
	jsonA := jsonAdd{
		a.Guarantee, a.LeftDeposit,
	}
	return json.Marshal(jsonA)
}

// UnmarshalJSON populates the receiver with the
// json-encoded data
func (a *Add) UnmarshalJSON(data []byte) error {
	var jsonA jsonAdd
	err := json.Unmarshal(data, &jsonA)
	if err != nil {
		return fmt.Errorf("error unmarshaling guarantee data: %w", err)
	}

	a.Guarantee = jsonA.Guarantee
	a.LeftDeposit = jsonA.LeftDeposit

	return nil
}

// jsonRemove replaces Remove's private fields with public ones,
// making it suitable for serialization
// embedded structs are moved to name fields for easier serialization
type jsonRemove struct {
	Target     types.Destination
	LeftAmount *big.Int
}

// MarshalJSON returns a JSON representation of the Remove
func (r Remove) MarshalJSON() ([]byte, error) {
	jsonR := jsonRemove(r)

	return json.Marshal(jsonR)
}

// UnmarshalJSON populates the receiver with the
// json-encoded data
func (r *Remove) UnmarshalJSON(data []byte) error {
	var jsonR jsonRemove
	err := json.Unmarshal(data, &jsonR)
	if err != nil {
		return fmt.Errorf("error unmarshaling remove data: %w", err)
	}

	r.Target = jsonR.Target
	r.LeftAmount = jsonR.LeftAmount

	return nil
}

// jsonProposal replaces Proposal's private fields with public ones,
// making it suitable for serialization
type jsonProposal struct {
	LedgerID   types.Destination
	ToAdd      Add
	ToRemove   Remove
	AddHTLC    HTLC
	RemoveHTLC []byte
}

// MarshalJSON returns a JSON representation of the Proposal
func (p Proposal) MarshalJSON() ([]byte, error) {
	jsonP := jsonProposal(p)

	return json.Marshal(jsonP)
}

// UnmarshalJSON populates the receiver with the
// json-encoded data
func (p *Proposal) UnmarshalJSON(data []byte) error {
	var jsonP jsonProposal
	err := json.Unmarshal(data, &jsonP)
	if err != nil {
		return fmt.Errorf("error unmarshaling guarantee data: %w", err)
	}

	p.LedgerID = jsonP.LedgerID
	p.ToAdd = jsonP.ToAdd
	p.ToRemove = jsonP.ToRemove

	return nil
}

// jsonBalance replaces Balance's private fields with public ones,
// making it suitable for serialization
type jsonBalance struct {
	Destination types.Destination
	Amount      *big.Int
}

// MarshalJSON returns a JSON representation of the Balance
func (b Balance) MarshalJSON() ([]byte, error) {
	jsonB := jsonBalance{
		b.destination, b.amount,
	}
	return json.Marshal(jsonB)
}

// UnmarshalJSON populates the receiver with the
// json-encoded data
func (b *Balance) UnmarshalJSON(data []byte) error {
	var jsonB jsonBalance
	err := json.Unmarshal(data, &jsonB)
	if err != nil {
		return fmt.Errorf("error unmarshaling guarantee data: %w", err)
	}

	b.destination = jsonB.Destination
	b.amount = jsonB.Amount

	return nil
}

// jsonGuarantee replaces Guarantee's private fields with public ones,
// making it suitable for serialization
type jsonGuarantee struct {
	Amount *big.Int
	Target types.Destination
	Left   types.Destination
	Right  types.Destination
}

// MarshalJSON returns a JSON representation of the Guarantee
func (g Guarantee) MarshalJSON() ([]byte, error) {
	jsonG := jsonGuarantee{
		g.amount, g.target, g.left, g.right,
	}
	return json.Marshal(jsonG)
}

// UnmarshalJSON populates the receiver with the
// json-encoded data
func (g *Guarantee) UnmarshalJSON(data []byte) error {
	var jsonG jsonGuarantee
	err := json.Unmarshal(data, &jsonG)
	if err != nil {
		return fmt.Errorf("error unmarshaling guarantee data: %w", err)
	}

	g.amount = jsonG.Amount
	g.target = jsonG.Target
	g.left = jsonG.Left
	g.right = jsonG.Right

	return nil
}

// jsonLedgerOutcome replaces LedgerOutcome's private fields with public ones,
// making it suitable for serialization
type jsonLedgerOutcome struct {
	AssetAddress types.Address // Address of the asset type
	Leader       Balance       // Balance of participants[0]
	Follower     Balance       // Balance of participants[1]
	Guarantees   map[types.Destination]Guarantee
}

// MarshalJSON returns a JSON representation of the LedgerOutcome
func (l LedgerOutcome) MarshalJSON() ([]byte, error) {
	jsonLo := jsonLedgerOutcome{
		AssetAddress: l.assetAddress,
		Leader:       l.leader,
		Follower:     l.follower,
		Guarantees:   l.guarantees,
	}
	return json.Marshal(jsonLo)
}

// UnmarshalJSON populates the receiver with the
// json-encoded data
func (l *LedgerOutcome) UnmarshalJSON(data []byte) error {
	var jsonLo jsonLedgerOutcome
	err := json.Unmarshal(data, &jsonLo)
	if err != nil {
		return fmt.Errorf("error unmarshaling ledger outcome data: %w", err)
	}

	l.assetAddress = jsonLo.AssetAddress
	l.leader = jsonLo.Leader
	l.follower = jsonLo.Follower
	l.guarantees = jsonLo.Guarantees

	return nil
}

// jsonConsensusChannel replaces ConsensusChannel's private fields with public ones,
// making it suitable for serialization
type jsonConsensusChannel struct {
	Id             types.Destination
	OnChainFunding types.Funds
	MyIndex        ledgerIndex
	FP             state.FixedPart
	Current        SignedVars
	ProposalQueue  []SignedProposal
}

// MarshalJSON returns a JSON representation of the ConsensusChannel
func (c ConsensusChannel) MarshalJSON() ([]byte, error) {
	jsonCh := jsonConsensusChannel{
		MyIndex:        c.MyIndex,
		FP:             c.fp,
		Id:             c.Id,
		OnChainFunding: c.OnChainFunding,
		Current:        c.current,
		ProposalQueue:  c.proposalQueue,
	}
	return json.Marshal(jsonCh)
}

// UnmarshalJSON populates the receiver with the
// json-encoded data
func (c *ConsensusChannel) UnmarshalJSON(data []byte) error {
	var jsonCh jsonConsensusChannel
	err := json.Unmarshal(data, &jsonCh)
	if err != nil {
		return fmt.Errorf("error unmarshaling channel data: %w", err)
	}

	c.Id = jsonCh.Id
	c.OnChainFunding = jsonCh.OnChainFunding
	c.MyIndex = jsonCh.MyIndex
	c.fp = jsonCh.FP
	c.current = jsonCh.Current
	c.proposalQueue = jsonCh.ProposalQueue

	return nil
}
