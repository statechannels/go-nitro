package store

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/channel"
	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/crypto"
	"github.com/statechannels/go-nitro/internal/safesync"
	"github.com/statechannels/go-nitro/payments"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/protocols/directdefund"
	"github.com/statechannels/go-nitro/protocols/directfund"
	"github.com/statechannels/go-nitro/protocols/virtualdefund"
	"github.com/statechannels/go-nitro/protocols/virtualfund"
	"github.com/statechannels/go-nitro/types"
)

type MemStore struct {
	objectives         safesync.Map[[]byte]
	channels           safesync.Map[[]byte]
	consensusChannels  safesync.Map[[]byte]
	channelToObjective safesync.Map[protocols.ObjectiveId]
	vouchers           safesync.Map[[]byte]

	key     string // the signing key of the store's engine
	address string // the (Ethereum) address associated to the signing key
}

func NewMemStore(key []byte) Store {
	ms := MemStore{}
	ms.key = common.Bytes2Hex(key)
	ms.address = crypto.GetAddressFromSecretKeyBytes(key).String()

	ms.objectives = safesync.Map[[]byte]{}
	ms.channels = safesync.Map[[]byte]{}
	ms.consensusChannels = safesync.Map[[]byte]{}
	ms.channelToObjective = safesync.Map[protocols.ObjectiveId]{}
	ms.vouchers = safesync.Map[[]byte]{}
	return &ms
}

func (ms *MemStore) Close() error {
	// Since this is a memory store, there is nothing to close
	return nil
}

func (ms *MemStore) GetAddress() *types.Address {
	address := common.HexToAddress(ms.address)
	return &address
}

func (ms *MemStore) GetChannelSecretKey() *[]byte {
	val := common.Hex2Bytes(ms.key)
	return &val
}

func (ms *MemStore) GetObjectiveById(id protocols.ObjectiveId) (protocols.Objective, error) {
	// todo: locking
	objJSON, ok := ms.objectives.Load(string(id))

	// return immediately if no such objective exists
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrNoSuchObjective, id)
	}

	obj, err := decodeObjective(id, objJSON)
	if err != nil {
		return nil, fmt.Errorf("error decoding objective %s: %w", id, err)
	}

	err = ms.populateChannelData(obj)
	if err != nil {
		// return existing objective data along with error
		return obj, fmt.Errorf("error populating channel data for objective %s: %w", id, err)
	}

	return obj, nil
}

func (ms *MemStore) SetObjective(obj protocols.Objective) error {
	// todo: locking
	objJSON, err := obj.MarshalJSON()
	if err != nil {
		return fmt.Errorf("error setting objective %s: %w", obj.Id(), err)
	}

	ms.objectives.Store(string(obj.Id()), objJSON)

	for _, rel := range obj.Related() {
		switch ch := rel.(type) {
		case *channel.VirtualChannel:
			err := ms.SetChannel(&ch.Channel)
			if err != nil {
				return fmt.Errorf("error setting virtual channel %s from objective %s: %w", ch.Id, obj.Id(), err)
			}
		case *channel.Channel:
			err := ms.SetChannel(ch)
			if err != nil {
				return fmt.Errorf("error setting channel %s from objective %s: %w", ch.Id, obj.Id(), err)
			}
		case *consensus_channel.ConsensusChannel:
			err := ms.SetConsensusChannel(ch)
			if err != nil {
				return fmt.Errorf("error setting consensus channel %s from objective %s: %w", ch.Id, obj.Id(), err)
			}
		default:
			return fmt.Errorf("unexpected type: %T", rel)
		}
	}

	// Objective ownership can only be transferred if the channel is not owned by another objective
	prevOwner, isOwned := ms.channelToObjective.Load(obj.OwnsChannel().String())
	if status := obj.GetStatus(); status == protocols.Approved {
		if !isOwned {
			ms.channelToObjective.Store(obj.OwnsChannel().String(), obj.Id())
		}
		if isOwned && prevOwner != obj.Id() {
			return fmt.Errorf("cannot transfer ownership of channel to from objective %s to %s", prevOwner, obj.Id())
		}
	}

	return nil
}

// SetChannel sets the channel in the store.
func (ms *MemStore) SetChannel(ch *channel.Channel) error {
	chJSON, err := ch.MarshalJSON()
	if err != nil {
		return err
	}

	ms.channels.Store(ch.Id.String(), chJSON)
	return nil
}

// DestroyChannel deletes the channel with id id.
func (ms *MemStore) DestroyChannel(id types.Destination) error {
	ms.channels.Delete(id.String())
	return nil
}

// SetConsensusChannel sets the channel in the store.
func (ms *MemStore) SetConsensusChannel(ch *consensus_channel.ConsensusChannel) error {
	if ch.Id.IsZero() {
		return fmt.Errorf("cannot store a channel with a zero id")
	}
	chJSON, err := ch.MarshalJSON()
	if err != nil {
		return err
	}

	ms.consensusChannels.Store(ch.Id.String(), chJSON)
	return nil
}

// DestroyChannel deletes the channel with id id.
func (ms *MemStore) DestroyConsensusChannel(id types.Destination) error {
	ms.consensusChannels.Delete(id.String())
	return nil
}

// GetChannelById retrieves the channel with the supplied id, if it exists.
func (ms *MemStore) GetChannelById(id types.Destination) (c *channel.Channel, ok bool) {
	ch, err := ms.getChannelById(id)
	if err != nil {
		return &channel.Channel{}, false
	}

	return &ch, true
}

// getChannelById returns the stored channel
func (ms *MemStore) getChannelById(id types.Destination) (channel.Channel, error) {
	chJSON, ok := ms.channels.Load(id.String())

	if !ok {
		return channel.Channel{}, ErrNoSuchChannel
	}

	var ch channel.Channel
	err := ch.UnmarshalJSON(chJSON)
	if err != nil {
		return channel.Channel{}, fmt.Errorf("error unmarshaling channel %s", ch.Id)
	}

	return ch, nil
}

// GetChannelsByIds returns a collection of channels with the given ids
func (ms *MemStore) GetChannelsByIds(ids []types.Destination) ([]*channel.Channel, error) {
	toReturn := []*channel.Channel{}

	var err error

	ms.channels.Range(func(key string, chJSON []byte) bool {
		var ch channel.Channel
		err = json.Unmarshal(chJSON, &ch)
		if err != nil {
			return false
		}

		// If the channel is one of the ones we're looking for, add it to the list
		if contains(ids, ch.Id) {
			toReturn = append(toReturn, &ch)
		}

		// If we've found all the channels we need, stop looking
		if len(toReturn) == len(ids) {
			return false
		}

		return true // otherwise, continue looking
	})
	if err != nil {
		return []*channel.Channel{}, err
	}
	return toReturn, nil
}

// GetChannelsByAppDefinition returns any channels that include the given app definition
func (ms *MemStore) GetChannelsByAppDefinition(appDef types.Address) ([]*channel.Channel, error) {
	toReturn := []*channel.Channel{}
	var err error
	ms.channels.Range(func(key string, chJSON []byte) bool {
		var ch channel.Channel
		err = json.Unmarshal(chJSON, &ch)
		if err != nil {
			return false
		}
		if ch.AppDefinition == appDef {
			toReturn = append(toReturn, &ch)
		}

		return true // channel not found: continue looking
	})

	if err != nil {
		return []*channel.Channel{}, err
	}

	return toReturn, nil
}

// GetChannelsByParticipant returns any channels that include the given participant
func (ms *MemStore) GetChannelsByParticipant(participant types.Address) ([]*channel.Channel, error) {
	toReturn := []*channel.Channel{}
	ms.channels.Range(func(key string, chJSON []byte) bool {
		var ch channel.Channel
		err := json.Unmarshal(chJSON, &ch)
		if err != nil {
			return true // channel not found, continue looking
		}

		participants := ch.FixedPart.Participants
		for _, p := range participants {
			if p == participant {
				toReturn = append(toReturn, &ch)
			}
		}

		return true // channel not found: continue looking
	})

	return toReturn, nil
}

// GetConsensusChannelById returns a ConsensusChannel with the given channel id
func (ms *MemStore) GetConsensusChannelById(id types.Destination) (channel *consensus_channel.ConsensusChannel, err error) {
	chJSON, ok := ms.consensusChannels.Load(id.String())

	if !ok {
		return &consensus_channel.ConsensusChannel{}, ErrNoSuchChannel
	}

	ch := &consensus_channel.ConsensusChannel{}
	err = ch.UnmarshalJSON(chJSON)

	if err != nil {
		return &consensus_channel.ConsensusChannel{}, fmt.Errorf("error unmarshaling channel %s", ch.Id)
	}

	return ch, nil
}

// GetConsensusChannel returns a ConsensusChannel between the calling node and
// the supplied counterparty, if such channel exists
func (ms *MemStore) GetConsensusChannel(counterparty types.Address) (channel *consensus_channel.ConsensusChannel, ok bool) {
	ms.consensusChannels.Range(func(key string, chJSON []byte) bool {
		var ch consensus_channel.ConsensusChannel
		err := json.Unmarshal(chJSON, &ch)
		if err != nil {
			return true // channel not found, continue looking
		}

		participants := ch.Participants()
		if len(participants) == 2 {
			if participants[0] == counterparty || participants[1] == counterparty {
				channel = &ch
				ok = true
				return false // we have found the target channel: break the Range loop
			}
		}

		return true // channel not found: continue looking
	})

	return
}

func (ms *MemStore) GetAllConsensusChannels() ([]*consensus_channel.ConsensusChannel, error) {
	toReturn := []*consensus_channel.ConsensusChannel{}
	var err error
	ms.consensusChannels.Range(func(key string, chJSON []byte) bool {
		var ch consensus_channel.ConsensusChannel

		err = json.Unmarshal(chJSON, &ch)
		if err != nil {
			return false
		}

		toReturn = append(toReturn, &ch)
		return true // channel not found: continue looking
	})
	if err != nil {
		return nil, err
	}
	return toReturn, nil
}

func (ms *MemStore) GetObjectiveByChannelId(channelId types.Destination) (protocols.Objective, bool) {
	// todo: locking
	id, found := ms.channelToObjective.Load(channelId.String())
	if !found {
		return &directfund.Objective{}, false
	}

	objective, err := ms.GetObjectiveById(protocols.ObjectiveId(id))
	return objective, err == nil
}

// populateChannelData fetches stored Channel data relevant to the given
// objective and attaches it to the objective. The channel data is attached
// in-place of the objectives existing channel pointers.
func (ms *MemStore) populateChannelData(obj protocols.Objective) error {
	id := obj.Id()

	switch o := obj.(type) {
	case *directfund.Objective:
		ch, err := ms.getChannelById(o.C.Id)
		if err != nil {
			return fmt.Errorf("error retrieving channel data for objective %s: %w", id, err)
		}

		o.C = &ch

		return nil
	case *directdefund.Objective:

		ch, err := ms.getChannelById(o.C.Id)
		if err != nil {
			return fmt.Errorf("error retrieving channel data for objective %s: %w", id, err)
		}

		o.C = &ch

		return nil
	case *virtualfund.Objective:
		v, err := ms.getChannelById(o.V.Id)
		if err != nil {
			return fmt.Errorf("error retrieving virtual channel data for objective %s: %w", id, err)
		}
		o.V = &channel.VirtualChannel{Channel: v}

		zeroAddress := types.Destination{}

		if o.ToMyLeft != nil &&
			o.ToMyLeft.Channel != nil &&
			o.ToMyLeft.Channel.Id != zeroAddress {

			left, err := ms.GetConsensusChannelById(o.ToMyLeft.Channel.Id)
			if err != nil {
				return fmt.Errorf("error retrieving left ledger channel data for objective %s: %w", id, err)
			}
			o.ToMyLeft.Channel = left
		}

		if o.ToMyRight != nil &&
			o.ToMyRight.Channel != nil &&
			o.ToMyRight.Channel.Id != zeroAddress {
			right, err := ms.GetConsensusChannelById(o.ToMyRight.Channel.Id)
			if err != nil {
				return fmt.Errorf("error retrieving right ledger channel data for objective %s: %w", id, err)
			}
			o.ToMyRight.Channel = right
		}

		return nil
	case *virtualdefund.Objective:
		v, err := ms.getChannelById(o.V.Id)
		if err != nil {
			return fmt.Errorf("error retrieving virtual channel data for objective %s: %w", id, err)
		}
		o.V = &channel.VirtualChannel{Channel: v}

		zeroAddress := types.Destination{}

		if o.ToMyLeft != nil &&
			o.ToMyLeft.Id != zeroAddress {

			left, err := ms.GetConsensusChannelById(o.ToMyLeft.Id)
			if err != nil {
				return fmt.Errorf("error retrieving left ledger channel data for objective %s: %w", id, err)
			}
			o.ToMyLeft = left
		}

		if o.ToMyRight != nil &&
			o.ToMyRight.Id != zeroAddress {
			right, err := ms.GetConsensusChannelById(o.ToMyRight.Id)
			if err != nil {
				return fmt.Errorf("error retrieving right ledger channel data for objective %s: %w", id, err)
			}
			o.ToMyRight = right
		}
		return nil
	default:
		return fmt.Errorf("objective %s did not correctly represent a known Objective type", id)
	}
}

// decodeObjective is a helper which encapsulates the deserialization
// of Objective JSON data. The decoded objectives will not have any
// channel data other than the channel Id.
func decodeObjective(id protocols.ObjectiveId, data []byte) (protocols.Objective, error) {
	switch {
	case directfund.IsDirectFundObjective(id):
		dfo := directfund.Objective{}
		err := dfo.UnmarshalJSON(data)
		return &dfo, err
	case directdefund.IsDirectDefundObjective(id):
		ddfo := directdefund.Objective{}
		err := ddfo.UnmarshalJSON(data)
		return &ddfo, err
	case virtualfund.IsVirtualFundObjective(id):
		vfo := virtualfund.Objective{}
		err := vfo.UnmarshalJSON(data)
		return &vfo, err
	case virtualdefund.IsVirtualDefundObjective(id):
		dvfo := virtualdefund.Objective{}
		err := dvfo.UnmarshalJSON(data)
		return &dvfo, err
	default:
		return nil, fmt.Errorf("objective id %s does not correspond to a known Objective type", id)

	}
}

func (ms *MemStore) ReleaseChannelFromOwnership(channelId types.Destination) error {
	ms.channelToObjective.Delete(channelId.String())
	return nil
}

func (ms *MemStore) SetVoucherInfo(channelId types.Destination, v payments.VoucherInfo) error {
	jsonData, err := json.Marshal(v)
	if err != nil {
		return err
	}
	ms.vouchers.Store(channelId.String(), jsonData)
	return nil
}

func (ms *MemStore) GetVoucherInfo(channelId types.Destination) (v *payments.VoucherInfo, ok bool) {
	data, ok := ms.vouchers.Load(channelId.String())
	if !ok {
		return nil, false
	}

	v = &payments.VoucherInfo{}
	err := json.Unmarshal(data, v)
	if err != nil {
		return nil, false
	}
	return v, true
}

func (ms *MemStore) RemoveVoucherInfo(channelId types.Destination) error {
	ms.vouchers.Delete(channelId.String())
	return nil
}

// contains is a helper function which returns true if the given item is included in col
func contains[T types.Destination | protocols.ObjectiveId](col []T, item T) bool {
	for _, i := range col {
		if i == item {
			return true
		}
	}
	return false
}
