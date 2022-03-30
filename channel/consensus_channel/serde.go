package consensus_channel

import (
	"encoding/json"
	"fmt"

	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/types"
)

// jsonLedgerOutcome replaces LedgerOutcome's private fields with public ones,
// making it suitable for serialization
type jsonLedgerOutcome struct {
	AssetAddress types.Address // Address of the asset type
	Left         Balance       // Balance of participants[0]
	Right        Balance       // Balance of participants[1]
	Guarantees   map[types.Destination]Guarantee
}

// MarshalJSON returns a JSON representation of the LedgerOutcome
func (o *LedgerOutcome) MarshalJSON() ([]byte, error) {
	jsonLo := jsonLedgerOutcome{
		AssetAddress: o.assetAddress,
		Left:         o.left,
		Right:        o.right,
		Guarantees:   o.guarantees,
	}
	return json.Marshal(jsonLo)
}

// UnmarshalJSON populates the receiver with the
// json-encoded data
func (o *LedgerOutcome) UnmarshalJSON(data []byte) error {
	var jsonLo jsonLedgerOutcome
	err := json.Unmarshal(data, &jsonLo)
	if err != nil {
		return fmt.Errorf("error unmarshaling ledger outcome data")
	}

	o.assetAddress = jsonLo.AssetAddress
	o.left = jsonLo.Left
	o.right = jsonLo.Right
	o.guarantees = jsonLo.Guarantees

	return nil
}

// jsonConsensusChannel replaces ConsensusChannel's private fields with public ones,
// making it suitable for serialization
type jsonConsensusChannel struct {
	Id            types.Destination
	MyIndex       ledgerIndex
	FP            state.FixedPart
	Current       SignedVars
	ProposalQueue []SignedProposal
}

// MarshalJSON returns a JSON representation of the ConsensusChannel
func (c ConsensusChannel) MarshalJSON() ([]byte, error) {
	jsonCh := jsonConsensusChannel{
		Id:            c.Id,
		MyIndex:       c.myIndex,
		FP:            c.fp,
		Current:       c.current,
		ProposalQueue: c.proposalQueue,
	}
	return json.Marshal(jsonCh)
}

// UnmarshalJSON populates the receiver with the
// json-encoded data
func (c *ConsensusChannel) UnmarshalJSON(data []byte) error {
	var jsonCh jsonConsensusChannel
	err := json.Unmarshal(data, &jsonCh)
	if err != nil {
		return fmt.Errorf("error unmarshaling channel data")
	}

	c.Id = jsonCh.Id
	c.myIndex = jsonCh.MyIndex
	c.fp = jsonCh.FP
	c.current = jsonCh.Current
	c.proposalQueue = jsonCh.ProposalQueue

	return nil
}
