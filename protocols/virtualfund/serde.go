package virtualfund

import (
	"encoding/json"
	"fmt"

	"github.com/statechannels/go-nitro/channel"
	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

// jsonConnection is a serialization-friendly struct representation
// of a Connection
type jsonConnection struct {
	Channel       types.Destination
	GuaranteeInfo GuaranteeInfo
}

// MarshalJSON returns a JSON representation of the Connection
// NOTE: Marshal -> Unmarshal is a lossy process. All channel data
// other than the ID is dropped
func (c Connection) MarshalJSON() ([]byte, error) {
	jsonC := jsonConnection{c.Channel.Id, c.GuaranteeInfo}
	bytes, err := json.Marshal(jsonC)
	if err != nil {
		return []byte{}, err
	}

	return bytes, err
}

// UnmarshalJSON populates the calling Connection with the
// json-encoded data
//
// NOTE: Marshal -> Unmarshal is a lossy process. All channel data from
//
//	(other than Id) is discarded
func (c *Connection) UnmarshalJSON(data []byte) error {
	c.Channel = &consensus_channel.ConsensusChannel{}

	if string(data) == "null" {
		// populate a well-formed but blank-addressed Connection
		c.Channel.Id = types.Destination{}
		return nil
	}

	var jsonC jsonConnection
	err := json.Unmarshal(data, &jsonC)
	if err != nil {
		return err
	}

	c.Channel.Id = jsonC.Channel
	c.GuaranteeInfo = jsonC.GuaranteeInfo

	return nil
}

// jsonObjective replaces the virtualfund Objective's channel pointers
// with the channel's respective IDs, making jsonObjective suitable for serialization
type jsonObjective struct {
	Status protocols.ObjectiveStatus
	V      types.Destination

	ToMyLeft  []byte
	ToMyRight []byte

	N      uint
	MyRole uint

	A0 types.Funds
	B0 types.Funds
}

// MarshalJSON returns a JSON representation of the VirtualFundObjective
// NOTE: Marshal -> Unmarshal is a lossy process. All channel data from
// the virtual and ledger channels (other than Ids) is discarded
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
		o.n,
		o.MyRole,
		o.a0,
		o.b0,
	}
	return json.Marshal(jsonVFO)
}

// UnmarshalJSON populates the calling VirtualFundObjective with the
// json-encoded data
//
// NOTE: Marshal -> Unmarshal is a lossy process. All channel data from
//
//	the virtual and ledger channels (other than Ids) is discarded
func (o *Objective) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	var jsonVFO jsonObjective
	if err := json.Unmarshal(data, &jsonVFO); err != nil {
		return fmt.Errorf("failed to unmarshal the VirtualFundObjective: %w", err)
	}

	o.V = &channel.VirtualChannel{}
	o.V.Id = jsonVFO.V

	o.ToMyLeft = &Connection{}
	o.ToMyRight = &Connection{}
	if err := o.ToMyLeft.UnmarshalJSON(jsonVFO.ToMyLeft); err != nil {
		return fmt.Errorf("failed to unmarshal left ledger channel: %w", err)
	}
	if err := o.ToMyRight.UnmarshalJSON(jsonVFO.ToMyRight); err != nil {
		return fmt.Errorf("failed to unmarshal right ledger channel: %w", err)
	}

	// connection.Unmarshal cannot populate a nil connection, so instead returns
	// a non-connection to a channel with the zero-Id. virtualfund.objective expects
	// these connection objects to be nil.
	zeroAddress := types.Destination{}
	if o.ToMyLeft.Channel.Id == zeroAddress {
		o.ToMyLeft = nil
	}
	if o.ToMyRight.Channel.Id == zeroAddress {
		o.ToMyRight = nil
	}

	o.Status = jsonVFO.Status
	o.n = jsonVFO.N
	o.MyRole = jsonVFO.MyRole
	o.a0 = jsonVFO.A0
	o.b0 = jsonVFO.B0

	return nil
}
