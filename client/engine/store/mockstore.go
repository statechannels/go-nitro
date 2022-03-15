package store

import (
	"fmt"

	"github.com/statechannels/go-nitro/channel"
	"github.com/statechannels/go-nitro/crypto"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/protocols/directfund"
	"github.com/statechannels/go-nitro/protocols/virtualfund"
	"github.com/statechannels/go-nitro/types"
)

type MockStore struct {
	objectives map[protocols.ObjectiveId]protocols.Objective
	channels   map[types.Destination]channel.Channel

	key     []byte        // the signing key of the store's engine
	address types.Address // the (Ethereum) address associated to the signing key
}

func NewMockStore(key []byte) Store {
	ms := MockStore{}
	ms.key = key
	ms.address = crypto.GetAddressFromSecretKeyBytes(key)

	ms.objectives = make(map[protocols.ObjectiveId]protocols.Objective)
	ms.channels = make(map[types.Destination]channel.Channel)

	return &ms
}

func (ms MockStore) GetAddress() *types.Address {
	return &ms.address
}

func (ms MockStore) GetChannelSecretKey() *[]byte {
	return &ms.key
}

func (ms MockStore) GetObjectiveById(id protocols.ObjectiveId) (obj protocols.Objective, ok bool) {
	// todo: locking
	obj, ok = ms.objectives[id]

	// return immediately if no such objective exists
	if !ok {
		return nil, ok
	}

	// populate channel data
	if dfo, isDirectFund := obj.(*directfund.Objective); isDirectFund {
		ch, err := ms.getChannelById(dfo.C.Id)

		if err != nil {
			return nil, false
		}

		dfo.C = &ch

		obj = dfo
	} else if vfo, isVirtualFund := obj.(*virtualfund.Objective); isVirtualFund {
		v, err := ms.getChannelById(vfo.V.Id)
		if err != nil {
			return nil, false
		}
		vfo.V = &channel.SingleHopVirtualChannel{Channel: v}

		if vfo.ToMyLeft != nil && vfo.ToMyLeft.Channel != nil {
			left, err := ms.getChannelById(vfo.ToMyLeft.Channel.Id)
			if err != nil {
				return nil, false
			}
			vfo.ToMyLeft.Channel = &channel.TwoPartyLedger{Channel: left}
		}

		if vfo.ToMyRight != nil && vfo.ToMyRight.Channel != nil {
			right, err := ms.getChannelById(vfo.ToMyRight.Channel.Id)
			if err != nil {
				return nil, false
			}
			vfo.ToMyRight.Channel = &channel.TwoPartyLedger{Channel: right}

		}

		obj = vfo
	}

	return obj, ok
}

func (ms MockStore) SetObjective(obj protocols.Objective) error {
	// todo: locking
	// todo: strip channel data from stored objective (avoid duplicate data-storage) (on serde PR)
	ms.objectives[obj.Id()] = obj

	for _, ch := range obj.Channels() {
		err := ms.SetChannel(ch)
		if err != nil {
			return fmt.Errorf("error setting channel %s from objective %s: %w", ch.Id, obj.Id(), err)
		}
	}

	return nil
}

// SetChannel sets the channel in the store.
func (ms *MockStore) SetChannel(ch *channel.Channel) error {
	ms.channels[ch.Id] = *ch

	return nil // temp - errors can exist / be reported when serde reintroduced
}

// getChannelById returns the stored channel
func (ms *MockStore) getChannelById(id types.Destination) (channel.Channel, error) {
	ch, ok := ms.channels[id]
	if ok {
		return ch, nil
	} else {
		return channel.Channel{}, fmt.Errorf("channel %s not found", id)
	}
}

// GetTwoPartyLedger returns a ledger channel between the two parties if it exists.
func (ms MockStore) GetTwoPartyLedger(firstParty types.Address, secondParty types.Address) (ledger *channel.TwoPartyLedger, ok bool) {
	for _, ch := range ms.channels {
		if len(ch.Participants) == 2 {
			// TODO: Should order matter?
			if ch.Participants[0] == firstParty && ch.Participants[1] == secondParty {
				return &channel.TwoPartyLedger{Channel: ch}, true
			}
		}

	}
	return nil, false
}

func (ms MockStore) GetObjectiveByChannelId(channelId types.Destination) (protocols.Objective, bool) {
	// todo: locking
	for _, obj := range ms.objectives {
		for _, ch := range obj.Channels() {
			if ch.Id == channelId {
				return obj, true
			}
		}
	}

	return nil, false
}
