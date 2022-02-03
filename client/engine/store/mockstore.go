package store

import (
	"crypto/ecdsa"
	"encoding/json"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/statechannels/go-nitro/channel"
	"github.com/statechannels/go-nitro/protocols"
	directfund "github.com/statechannels/go-nitro/protocols/direct-fund"
	virtualfund "github.com/statechannels/go-nitro/protocols/virtual-fund"
	"github.com/statechannels/go-nitro/types"
)

type MockStore struct {
	objectives map[protocols.ObjectiveId]string
	channels   map[types.Destination]string

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

	ms.objectives = make(map[protocols.ObjectiveId]string)
	ms.channels = make(map[types.Destination]string)

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
	objJSON, ok := ms.objectives[id]

	if strings.HasPrefix(string(id), "DirectFunding-") {
		// unmarshal directFundObjective
		var jsonDFO directfund.JSONDirectFundObjective
		json.Unmarshal([]byte(objJSON), &jsonDFO)
		channel, _ := ms.getChannelById(jsonDFO.C)
		jsonDFO.DirectFundObjective.C = &channel

		obj = jsonDFO.DirectFundObjective
	} else if strings.HasPrefix(string(id), "VirtualFund-") {
		// unmarshal virtalFundObj
		var jsonVFO virtualfund.JSONVirtualFundObjective
		json.Unmarshal([]byte(objJSON), &jsonVFO)
		channel, _ := ms.getChannelById(jsonVFO.V)
		jsonVFO.VirtualFundObjective.V = &channel

		obj = jsonVFO.VirtualFundObjective
	}

	return obj, ok
}

func (ms MockStore) SetObjective(obj protocols.Objective) error {
	// todo: locking
	bytes, err := obj.MarshalJSON()
	if err != nil {
		return err
	}
	ms.objectives[obj.Id()] = string(bytes)

	id := obj.Id()
	if strings.HasPrefix(string(id), "DirectFunding-") {
		// marshal and persist ledger chanel
		dfo, _ := obj.(directfund.DirectFundObjective)
		ms.setChannel(*dfo.C)
	} else if strings.HasPrefix(string(id), "VirtualFund-") {
		// marshal and persist virtual channel
		vfo, _ := obj.(virtualfund.VirtualFundObjective)
		ms.setChannel(*vfo.V)
	}

	return nil
}

func (ms MockStore) setChannel(c channel.Channel) error {
	bytes, err := json.Marshal(c)
	if err != nil {
		return err
	}
	ms.channels[c.Id] = string(bytes)
	return nil
}

func (ms MockStore) getChannelById(id types.Destination) (channel.Channel, error) {
	channelJSON := ms.channels[id]
	var channel channel.Channel
	err := json.Unmarshal([]byte(channelJSON), &channel)
	return channel, err
}

	return nil, false
}

func (ms MockStore) UpdateProgressLastMadeAt(protocols.ObjectiveId, protocols.WaitingFor) {
	// todo
}
