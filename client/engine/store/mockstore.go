package store

import (
	"strings"

	"github.com/statechannels/go-nitro/channel"
	"github.com/statechannels/go-nitro/crypto"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/protocols/directfund"
	"github.com/statechannels/go-nitro/protocols/virtualfund"
	"github.com/statechannels/go-nitro/types"
)

type MockStore struct {
	objectives map[protocols.ObjectiveId]protocols.Objective

	key     []byte        // the signing key of the store's engine
	address types.Address // the (Ethereum) address associated to the signing key
}

func NewMockStore(key []byte) Store {
	ms := MockStore{}
	ms.key = key
	ms.address = crypto.GetAddressFromSecretKeyBytes(key)

	ms.objectives = make(map[protocols.ObjectiveId]protocols.Objective)

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
	return obj, ok
}

func (ms MockStore) SetObjective(obj protocols.Objective) error {
	// todo: locking
	ms.objectives[obj.Id()] = obj
	return nil
}

// SetChannel sets the channel in the store.
func (ms *MockStore) SetChannel(ch *channel.Channel) error {
	// TODO: This is a temporary implementation that is pretty clunky.
	// This should be replaced in https://github.com/statechannels/go-nitro/pull/227
	for _, obj := range ms.objectives {
		if strings.HasPrefix(string(obj.Id()), "DirectFunding-") {
			dfO := obj.(*directfund.Objective)
			if dfO.C.Id == ch.Id {
				dfO.C = ch
				err := ms.SetObjective(dfO)
				if err != nil {
					return err
				}

			}
		} else if strings.HasPrefix(string(obj.Id()), "VirtualFund-") {
			vfO := obj.(*virtualfund.Objective)
			if vfO.V.Id == ch.Id {
				vfO.V = &channel.SingleHopVirtualChannel{Channel: *ch}
				err := ms.SetObjective(vfO)
				if err != nil {
					return err
				}

			}
			if vfO.ToMyLeft != nil && vfO.ToMyLeft.Channel.Id == ch.Id {
				vfO.ToMyLeft.Channel = &channel.TwoPartyLedger{Channel: *ch}
				err := ms.SetObjective(vfO)
				if err != nil {
					return err
				}

			}
			if vfO.ToMyRight != nil && vfO.ToMyRight.Channel.Id == ch.Id {
				vfO.ToMyRight.Channel = &channel.TwoPartyLedger{Channel: *ch}
				err := ms.SetObjective(vfO)
				if err != nil {
					return err
				}

			}
		}
	}
	return nil
}

// GetTwoPartyLedger returns a ledger channel between the two parties if it exists.
func (ms MockStore) GetTwoPartyLedger(firstParty types.Address, secondParty types.Address) (ledger *channel.TwoPartyLedger, ok bool) {
	for _, obj := range ms.objectives {
		for _, ch := range obj.Channels() {
			if len(ch.Participants) == 2 {
				// TODO: Should order matter?
				if ch.Participants[0] == firstParty && ch.Participants[1] == secondParty {
					return &channel.TwoPartyLedger{Channel: *ch}, true
				}
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
