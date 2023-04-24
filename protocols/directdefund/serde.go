package directdefund

import (
	"encoding/json"

	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

// jsonObjective replaces the directdefund.Objective's channel pointer with
// the channel's ID, making jsonObjective suitable for serialization
type jsonObjective struct {
	Status                protocols.ObjectiveStatus
	C                     types.Destination
	FinalTurnNum          uint64
	TransactionSumbmitted bool
}

// MarshalJSON returns a JSON representation of the DirectDefundObjective
// NOTE: Marshal -> Unmarshal is a lossy process. All channel data
// (other than Id) from the field C is discarded
func (o Objective) MarshalJSON() ([]byte, error) {
	jsonDDFO := jsonObjective{
		o.Status,
		o.C.Id,
		o.finalTurnNum,
		o.transactionSubmitted,
	}

	return json.Marshal(jsonDDFO)
}

// UnmarshalJSON populates the calling DirectDefundObjective with the
// json-encoded data
// NOTE: Marshal -> Unmarshal is a lossy process. All channel data
// (other than Id) from the field C is discarded
func (o *Objective) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	var jsonDDFO jsonObjective
	err := json.Unmarshal(data, &jsonDDFO)
	if err != nil {
		return err
	}

	o.C = &consensus_channel.ConsensusChannel{}

	o.Status = jsonDDFO.Status
	o.C.Id = jsonDDFO.C
	o.finalTurnNum = jsonDDFO.FinalTurnNum
	o.transactionSubmitted = jsonDDFO.TransactionSumbmitted

	return nil
}
