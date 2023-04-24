package virtualdefund

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

// jsonObjective replaces the virtualfund Objective's channel pointers
// with the channel's respective IDs, making jsonObjective suitable for serialization
type jsonObjective struct {
	Status         protocols.ObjectiveStatus
	VFixed         state.FixedPart
	InitialOutcome outcome.SingleAssetExit
	FinalOutcome   outcome.SingleAssetExit
	Signatures     []state.Signature

	ToMyLeft             types.Destination
	ToMyRight            types.Destination
	MinimumPaymentAmount *big.Int
	MyRole               uint
}

// MarshalJSON returns a JSON representation of the VirtualDefundObjective
func (o Objective) MarshalJSON() ([]byte, error) {
	var left types.Destination
	var right types.Destination

	if o.ToMyLeft != nil {
		left = o.ToMyLeft.Id
	}

	if o.ToMyRight != nil {
		right = o.ToMyRight.Id
	}

	jsonVFO := jsonObjective{
		Status:               o.Status,
		VFixed:               o.VFixed,
		Signatures:           o.Signatures,
		FinalOutcome:         o.FinalOutcome,
		InitialOutcome:       o.InitialOutcome,
		ToMyLeft:             left,
		ToMyRight:            right,
		MyRole:               o.MyRole,
		MinimumPaymentAmount: o.MinimumPaymentAmount,
	}
	return json.Marshal(jsonVFO)
}

// UnmarshalJSON populates the calling VirtualDefundObjective with the
// json-encoded data
func (o *Objective) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	var jsonVFO jsonObjective
	if err := json.Unmarshal(data, &jsonVFO); err != nil {
		return fmt.Errorf("failed to unmarshal the VirtualDefundObjective: %w", err)
	}

	o.ToMyLeft = &consensus_channel.ConsensusChannel{}
	o.ToMyLeft.Id = jsonVFO.ToMyLeft
	o.ToMyRight = &consensus_channel.ConsensusChannel{}
	o.ToMyRight.Id = jsonVFO.ToMyRight

	o.Status = jsonVFO.Status

	o.MyRole = jsonVFO.MyRole
	o.Signatures = jsonVFO.Signatures
	o.InitialOutcome = jsonVFO.InitialOutcome
	o.VFixed = jsonVFO.VFixed
	o.FinalOutcome = jsonVFO.FinalOutcome
	o.MinimumPaymentAmount = jsonVFO.MinimumPaymentAmount

	return nil
}
