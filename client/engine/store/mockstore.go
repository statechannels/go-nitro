package store

import (
	"crypto/ecdsa"
	"log"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/statechannels/go-nitro/channel"
	"github.com/statechannels/go-nitro/protocols"
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

	channelSecretKey, err := crypto.ToECDSA(ms.key)
	if err != nil {
		log.Fatal("error casting public key to ECDSA")
	}
	publicKey := channelSecretKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}
	ms.address = crypto.PubkeyToAddress(*publicKeyECDSA)

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
func (ms MockStore) SetChannel(channel *channel.Channel) error {
	// TODO: Right now we're just updating the channel in memory, so setChannel can be a no-op.
	// This should be updated in https://github.com/statechannels/go-nitro/issues/191
	return nil

}

func (ms MockStore) GetChannel(channelId types.Destination) (*channel.Channel, bool) {
	// todo: locking
	for _, obj := range ms.objectives {
		for _, ch := range obj.Channels() {
			if ch.Id == channelId {
				return ch, true
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

func (ms MockStore) UpdateProgressLastMadeAt(protocols.ObjectiveId, protocols.WaitingFor) {
	// todo
}
