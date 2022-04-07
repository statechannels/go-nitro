package virtualdefund

import (
	"encoding/json"
	"fmt"

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

	ToMyLeft  []byte
	ToMyRight []byte

	MyRole uint
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
		o.Status,
		o.V.Id,
		left,
		right,

		o.MyRole,
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

	o.V = &channel.SingleHopVirtualChannel{}
	o.V.Id = jsonVFO.V

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

	return nil
}
