package directfund

import (
	"encoding/json"

	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

// jsonObjective replaces the directfund.Objective's channel pointer with the
// channel's ID, making jsonObjective suitable for serialization
type jsonObjective struct {
	Status protocols.ObjectiveStatus
	C      types.Destination

	MyDepositSafetyThreshold types.Funds
	MyDepositTarget          types.Funds
	FullyFundedThreshold     types.Funds
	LatestBlockNumber        uint64
	TransactionSumbmitted    bool
}

// MarshalJSON returns a JSON representation of the DirectFundObjective
// NOTE: Marshal -> Unmarshal is a lossy process. All channel data
// (other than Id) from the field C is discarded
func (o Objective) MarshalJSON() ([]byte, error) {
	jsonDFO := jsonObjective{
		o.Status,
		o.C.Id,
		o.myDepositSafetyThreshold,
		o.myDepositTarget,
		o.fullyFundedThreshold,
		o.latestBlockNumber,
		o.transactionSubmitted,
	}
	return json.Marshal(jsonDFO)
}

// UnmarshalJSON populates the calling DirectFundObjective with the
// json-encoded data
// NOTE: Marshal -> Unmarshal is a lossy process. All channel data
// (other than Id) from the field C is discarded
func (o *Objective) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	var jsonDFO jsonObjective
	err := json.Unmarshal(data, &jsonDFO)
	if err != nil {
		return err
	}

	o.C = &consensus_channel.ConsensusChannel{}
	o.C.Id = jsonDFO.C

	o.Status = jsonDFO.Status
	o.fullyFundedThreshold = jsonDFO.FullyFundedThreshold
	o.myDepositTarget = jsonDFO.MyDepositTarget
	o.myDepositSafetyThreshold = jsonDFO.MyDepositSafetyThreshold
	o.latestBlockNumber = jsonDFO.LatestBlockNumber
	o.transactionSubmitted = jsonDFO.TransactionSumbmitted

	return nil
}
