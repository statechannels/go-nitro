package virtualdefund

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/protocols"
)

// jsonObjective replaces the virtualfund Objective's channel pointers
// with the channel's respective IDs, making jsonObjective suitable for serialization
type jsonObjective struct {
	Status         protocols.ObjectiveStatus
	VFixed         state.FixedPart
	InitialOutcome outcome.SingleAssetExit
	FinalOutcome   outcome.SingleAssetExit
	Signatures     [3]state.Signature

	ToMyLeft             []byte
	ToMyRight            []byte
	MinimumPaymentAmount *big.Int
	MyRole               uint
}

// MarshalJSON returns a JSON representation of the VirtualDefundObjective
func (o Objective) MarshalJSON() ([]byte, error) {
	var left []byte
	var right []byte
	var err error

	if o.ToMyLeft == nil {
		left = []byte("null")
	} else {
		left, err = o.ToMyLeft.MarshalJSON()

		if err != nil {
			return nil, fmt.Errorf("error marshaling left channel of %v: %w", o, err)
		}
	}

	if o.ToMyRight == nil {
		right = []byte("null")
	} else {
		right, err = o.ToMyRight.MarshalJSON()

		if err != nil {
			return nil, fmt.Errorf("error marshaling right channel of %v: %w", o, err)
		}
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
	o.ToMyRight = &consensus_channel.ConsensusChannel{}
	if err := o.ToMyLeft.UnmarshalJSON(jsonVFO.ToMyLeft); err != nil {
		return fmt.Errorf("failed to unmarshal left ledger channel: %w", err)
	}
	if err := o.ToMyRight.UnmarshalJSON(jsonVFO.ToMyRight); err != nil {
		return fmt.Errorf("failed to unmarshal right ledger channel: %w", err)
	}

	o.Status = jsonVFO.Status

	o.MyRole = jsonVFO.MyRole
	o.Signatures = jsonVFO.Signatures
	o.InitialOutcome = jsonVFO.InitialOutcome
	o.VFixed = jsonVFO.VFixed
	o.FinalOutcome = jsonVFO.FinalOutcome
	o.MinimumPaymentAmount = jsonVFO.MinimumPaymentAmount

	return nil
}
