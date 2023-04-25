package virtualdefund

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/statechannels/go-nitro/channel"
	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

// jsonObjective replaces the virtualfund Objective's channel pointers
// with the channel's respective IDs, making jsonObjective suitable for serialization
type jsonObjective struct {
	Status protocols.ObjectiveStatus
	V      types.Destination

	ToMyLeft             types.Destination
	ToMyRight            types.Destination
	MinimumPaymentAmount *big.Int
	MyRole               uint
}

// MarshalJSON returns a JSON representation of the VirtualDefundObjective
//
// NOTE: Marshal -> Unmarshal is a lossy process. All channel data
// (other than Id) from the fields ToMyLeft,ToMyRight are discarded
func (o Objective) MarshalJSON() ([]byte, error) {
	var left types.Destination
	var right types.Destination
	var V types.Destination

	if o.ToMyLeft != nil {
		left = o.ToMyLeft.Id
	}

	if o.ToMyRight != nil {
		right = o.ToMyRight.Id
	}

	jsonVFO := jsonObjective{
		Status:               o.Status,
		V:                    V,
		ToMyLeft:             left,
		ToMyRight:            right,
		MyRole:               o.MyRole,
		MinimumPaymentAmount: o.MinimumPaymentAmount,
	}
	return json.Marshal(jsonVFO)
}

// UnmarshalJSON populates the calling VirtualDefundObjective with the
// json-encoded data
//
// NOTE: Marshal -> Unmarshal is a lossy process. All channel data
// (other than Id) from the fields ToMyLeft,ToMyRight are discarded
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

	o.V = &channel.Channel{}
	o.V.Id = jsonVFO.V

	o.MinimumPaymentAmount = jsonVFO.MinimumPaymentAmount

	return nil
}
