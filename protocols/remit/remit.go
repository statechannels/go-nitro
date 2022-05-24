package remit

import (
	"encoding/json"
	"math/big"

	"github.com/statechannels/go-nitro/channel"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

const ObjectivePrefix = "Remit-"

type Objective struct {
	Status protocols.ObjectiveStatus
	C      *channel.Channel

	payer types.Destination
	payee types.Destination

	amount *big.Int
}

func (o *Objective) Crank(secretKey *[]byte) (protocols.Objective, protocols.SideEffects, protocols.WaitingFor, error) {
}

func (o *Objective) Update(event protocols.ObjectiveEvent) (protocols.Objective, error) {}

// OwnsChannel returns the channel that the objective is funding.
func (o *Objective) OwnsChannel() types.Destination {
	return o.C.Id
}

// GetStatus returns the status of the objective.
func (o *Objective) GetStatus() protocols.ObjectiveStatus {
	return o.Status
}

func (o *Objective) Id() protocols.ObjectiveId {
	return protocols.ObjectiveId(ObjectivePrefix + o.C.Id.String() + ":" + o.payer.String() + "=>" + o.payee.String() + o.amount.String())
}

// clone returns a deep copy of the receiver.
func (o *Objective) clone() Objective {
	clone := Objective{}
	clone.Status = o.Status

	cClone := o.C.Clone()
	clone.C = cClone

	clone.payer = o.payer
	clone.payee = o.payee

	clone.amount = big.NewInt(0).Set(o.amount)

	return clone
}

func (o *Objective) Approve() protocols.Objective {
	updated := o.clone()
	// todo: consider case of s.Status == Rejected
	updated.Status = protocols.Approved

	return &updated
}

func (o *Objective) Reject() protocols.Objective {
	updated := o.clone()
	updated.Status = protocols.Rejected
	return &updated
}

type jsonObjective struct {
	Status protocols.ObjectiveStatus
	C      types.Destination
	Payer  types.Destination
	Payee  types.Destination
	Amount *big.Int
}

// MarshalJSON returns a JSON representation of the DirectFundObjective
//
// NOTE: Marshal -> Unmarshal is a lossy process. All channel data
//       (other than Id) from the field C is discarded
func (o Objective) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		jsonObjective{
			o.Status,
			o.C.Id,
			o.payer,
			o.payee,
			o.amount,
		})
}

// UnmarshalJSON populates the calling DirectFundObjective with the
// json-encoded data
//
// NOTE: Marshal -> Unmarshal is a lossy process. All channel data
//       (other than Id) from the field C is discarded
func (o *Objective) UnmarshalJSON(data []byte) error {

	if string(data) == "null" {
		return nil
	}

	var jsonO jsonObjective
	err := json.Unmarshal(data, &jsonO)

	if err != nil {
		return err
	}

	o.C = &channel.Channel{}
	o.C.Id = jsonO.C

	o.Status = jsonO.Status
	o.payer = jsonO.Payer
	o.payee = jsonO.Payee
	o.amount = jsonO.Amount
	return nil
}

func (o *Objective) Related() []protocols.Storable {
	return []protocols.Storable{o.C}
}
