package store

import (
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

type MockStore struct {
	objectives map[protocols.ObjectiveId]protocols.Objective

	key []byte // the signing key of the store's engine
}

func NewMockStore(key []byte) Store {
	ms := MockStore{}
	ms.key = key
	ms.objectives = make(map[protocols.ObjectiveId]protocols.Objective)

	return ms
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

func (ms MockStore) GetObjectiveByChannelId(channelId types.Destination) (protocols.Objective, bool) {
	// todo: locking
	for _, obj := range ms.objectives {
		for _, ch := range obj.Channels() {
			if ch == channelId {
				return obj, true
			}
		}
	}

	return nil, false
}

func (ms MockStore) UpdateProgressLastMadeAt(protocols.ObjectiveId, protocols.WaitingFor) {
	// todo
}
