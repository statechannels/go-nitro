package store

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/statechannels/go-nitro/channel"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/protocols/directfund"
	"github.com/statechannels/go-nitro/protocols/virtualfund"
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

func (ms MockStore) GetObjectiveById(id protocols.ObjectiveId) (protocols.Objective, error) {
	// todo: locking
	objJSON, ok := ms.objectives[id]

	if !ok {
		return nil, fmt.Errorf("no protocol with id %s found", id)
	}

	if strings.HasPrefix(string(id), "Direct") {
		dfo := directfund.DirectFundObjective{}

		err := dfo.UnmarshalJSON([]byte(objJSON))

		if err != nil {
			return nil, fmt.Errorf("error unmarshaling objective %s: %w", id, err)
		}

		if channel, err := ms.getChannelById(dfo.C.Id); err == nil {
			dfo.C = &channel
		} else {
			return nil, fmt.Errorf("error retrieving channel data for objective %s: %w", id, err)
		}

		return &dfo, err

	} else if strings.HasPrefix(string(id), "Virtual") {
		vfo := virtualfund.VirtualFundObjective{}

		err := vfo.UnmarshalJSON([]byte(objJSON))

		if err != nil {
			return nil, fmt.Errorf("error unmarshaling objective %s: %w", id, err)
		}

		if virtualChannel, err := ms.getChannelById(vfo.V.Id); err == nil {
			vfo.V = &virtualChannel
		} else {
			return nil, fmt.Errorf("error retrieving virtual channel data for objective %s: %w", id, err)
		}

		if leftChannel, err := ms.getChannelById(vfo.ToMyLeft.Channel.Id); err == nil {
			vfo.ToMyLeft.Channel = channel.TwoPartyLedger{Channel: leftChannel}
		} else {
			return nil, fmt.Errorf("error retrieving left-ledger channel data for objective %s: %w", id, err)
		}

		if rightChannel, err := ms.getChannelById(vfo.ToMyRight.Channel.Id); err == nil {
			vfo.ToMyRight.Channel = channel.TwoPartyLedger{Channel: rightChannel}
		} else {
			return nil, fmt.Errorf("error retrieving right-ledger channel data for objective %s: %w", id, err)
		}

		return &vfo, nil
	}

	return nil, fmt.Errorf("objective %s not a recognised objective type", id)
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
		dfo, _ := obj.(*directfund.DirectFundObjective)
		if err := ms.setChannel(*dfo.C); err != nil {
			return fmt.Errorf("failed to write channel data for %s: %w", obj.Id(), err)
		}
	} else if strings.HasPrefix(string(id), "VirtualFund-") {
		// marshal and persist virtual channel
		vfo, _ := obj.(*virtualfund.VirtualFundObjective)
		if err := ms.setChannel(*vfo.V); err != nil {
			return fmt.Errorf("failed to write virtual-channel data for %s: %w", obj.Id(), err)
		}
		// marshal and persist ledger channel(s)
		if vfo.ToMyLeft != nil {
			if err := ms.setChannel(vfo.ToMyLeft.Channel.Channel); err != nil {
				return fmt.Errorf("failed to write left ledger-channel data for %s: %w", obj.Id(), err)
			}
		}
		if vfo.ToMyRight != nil {
			if err := ms.setChannel(vfo.ToMyRight.Channel.Channel); err != nil {
				return fmt.Errorf("failed to write right ledger-channel data for %s: %w", obj.Id(), err)
			}
		}
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

func (ms MockStore) GetObjectiveByChannelId(channelId types.Destination) (protocols.Objective, error) {
	// todo: locking
	for id, objJSON := range ms.objectives {
		var obj protocols.Objective

		if strings.HasPrefix(string(id), "Direct") {
			var dfo directfund.DirectFundObjective
			if err := dfo.UnmarshalJSON([]byte(objJSON)); err != nil {
				continue
			}
			obj = &dfo
		}

		if strings.HasPrefix(string(id), "Virtual") {
			var vfo virtualfund.VirtualFundObjective
			if err := vfo.UnmarshalJSON([]byte(objJSON)); err != nil {
				continue
			}
			obj = &vfo
		}

		for _, ch := range obj.Channels() {
			if ch == channelId {
				return obj, nil
			}
		}
	}

	return nil, fmt.Errorf("no objectives found for channel %s", channelId)
}

func (ms MockStore) UpdateProgressLastMadeAt(protocols.ObjectiveId, protocols.WaitingFor) {
	// todo
}
