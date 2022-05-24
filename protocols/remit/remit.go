package remit

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/statechannels/go-nitro/channel"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

const ObjectivePrefix = "Remit-"

const (
	WaitingForNothing protocols.WaitingFor = "WaitingForNothing" // This protocol finishes immediately
)

type Objective struct {
	Status protocols.ObjectiveStatus
	C      *channel.Channel

	payer types.Destination
	payee types.Destination

	amount *big.Int
}

// GetChannelByIdFunction specifies a function that can be used to retrieve channels from a store.
type GetChannelByIdFunction func(id types.Destination) (channel *channel.Channel, ok bool)

func NewObjective(r ObjectiveRequest, preApprove bool, getChannel GetChannelByIdFunction) (Objective, error) {
	c, ok := getChannel(r.CId)
	if !ok {
		return Objective{}, errors.New("could not find channel")
	}

	o := Objective{
		Status: protocols.Unapproved,
		C:      c,
		payer:  r.Payer,
		payee:  r.Payee,
		amount: r.Amount,
	}

	if preApprove {
		o.Status = protocols.Approved
	}
	return o, nil
}

type ObjectiveRequest struct {
	CId types.Destination

	Payer types.Destination
	Payee types.Destination

	Amount *big.Int
}

// Id returns the objective id for the request.
func (r ObjectiveRequest) Id(myAddress types.Address) protocols.ObjectiveId {

	return protocols.ObjectiveId(ObjectivePrefix + r.CId.String() + ":" + r.Payer.String() + "=>" + r.Payee.String())
}

func (o *Objective) Crank(secretKey *[]byte) (protocols.Objective, protocols.SideEffects, protocols.WaitingFor, error) {
	n := o.clone()
	se := protocols.SideEffects{}
	if o.payer == o.C.MyDestination() {

		ss, err := n.C.MakePayment(n.payee, n.amount, secretKey)
		if err != nil {
			return o, protocols.SideEffects{}, WaitingForNothing, fmt.Errorf("Cannot make payment: %w", err)
		}

		payeeAddress, err := o.payee.ToAddress()

		if err != nil {
			return o, protocols.SideEffects{}, WaitingForNothing, fmt.Errorf("Cannot make payment: %w", err)
		}

		msg := protocols.CreatePaymentMessage(o.Id(), ss, payeeAddress)
		se.MessagesToSend = []protocols.Message{msg}
	}
	return &n, se, WaitingForNothing, nil

}

func (o *Objective) Update(event protocols.ObjectiveEvent) (protocols.Objective, error) {
	if o.Id() != event.ObjectiveId {
		return o, fmt.Errorf("event and objective Ids do not match: %s and %s respectively", string(event.ObjectiveId), string(o.Id()))
	}

	updated := o.clone()
	updated.C.AddSignedState(event.SignedState)
	return &updated, nil
}

// OwnsChannel returns the channel that the objective is funding.
func (o *Objective) OwnsChannel() types.Destination {
	return o.C.Id
}

// GetStatus returns the status of the objective.
func (o *Objective) GetStatus() protocols.ObjectiveStatus {
	return o.Status
}

func (o *Objective) Id() protocols.ObjectiveId {
	return protocols.ObjectiveId(ObjectivePrefix + o.C.Id.String() + ":" + o.payer.String() + "=>" + o.payee.String())
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

// IsRemitObjective inspects a objective id and returns true if the objective id is for a direct fund objective.
func IsRemitObjective(id protocols.ObjectiveId) bool {
	return strings.HasPrefix(string(id), ObjectivePrefix)
}

// ConstructObjectiveFromState takes in a state and constructs an objective from it.
func ConstructObjectiveFromState(
	s state.State,
	getChannel GetChannelByIdFunction,
) (Objective, error) {
	ch, ok := getChannel(s.ChannelId())

	if !ok {
		return Objective{}, errors.New("could not find channel")
	}

	o := Objective{
		C:      ch,
		payer:  types.AddressToDestination(ch.Participants[0]), // TODO these fields are only correct given certain assumptions
		payee:  ch.MyDestination(),
		amount: big.NewInt(0), // TODO this field is incorrect
	}
	return o, nil
}
