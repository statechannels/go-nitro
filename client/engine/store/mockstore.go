package store

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/statechannels/go-nitro/channel"
	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/crypto"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/protocols/directdefund"
	"github.com/statechannels/go-nitro/protocols/directfund"
	"github.com/statechannels/go-nitro/protocols/virtualfund"
	"github.com/statechannels/go-nitro/types"
)

type MockStore struct {
	objectives         syncMap[[]byte]
	channels           syncMap[[]byte]
	consensusChannels  syncMap[[]byte]
	channelToObjective syncMap[protocols.ObjectiveId]

	key     []byte        // the signing key of the store's engine
	address types.Address // the (Ethereum) address associated to the signing key
}

// syncMap wraps sync.Map in order to provide type safety
type syncMap[T any] struct {
	m sync.Map
}

// Load returns the value stored in the map for a key, or nil if no
// value is present.
// The ok result indicates whether value was found in the map.
func (o *syncMap[T]) Load(id string) (value T, ok bool) {
	data, ok := o.m.Load(id)

	if !ok {
		var result T
		return result, false
	}

	value = data.(T)

	return value, ok
}

// Store sets the value for a key.
func (o *syncMap[T]) Store(key string, data T) {
	o.m.Store(key, data)
}

// Range calls f sequentially for each key and value present in the map.
// If f returns false, range stops the iteration.
//
// Range does not necessarily correspond to any consistent snapshot of the Map's
// contents: no key will be visited more than once, but if the value for any key
// is stored or deleted concurrently, Range may reflect any mapping for that key
// from any point during the Range call.
//
// Range may be O(N) with the number of elements in the map even if f returns
// false after a constant number of calls.
func (o *syncMap[T]) Range(f func(key string, value T) bool) {
	untypedF := func(key, value interface{}) bool {
		return f(key.(string), value.(T))
	}
	o.m.Range(untypedF)
}

func NewMockStore(key []byte) Store {
	ms := MockStore{}
	ms.key = key
	ms.address = crypto.GetAddressFromSecretKeyBytes(key)

	ms.objectives = syncMap[[]byte]{}
	ms.channels = syncMap[[]byte]{}
	ms.consensusChannels = syncMap[[]byte]{}
	ms.channelToObjective = syncMap[protocols.ObjectiveId]{}

	return &ms
}

func (ms *MockStore) GetAddress() *types.Address {
	return &ms.address
}

func (ms *MockStore) GetChannelSecretKey() *[]byte {
	return &ms.key
}

func (ms *MockStore) GetObjectiveById(id protocols.ObjectiveId) (protocols.Objective, error) {
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

func (ms *MockStore) SetObjective(obj protocols.Objective) error {
	// todo: locking
	objJSON, err := obj.MarshalJSON()

	if err != nil {
		return fmt.Errorf("error setting objective %s: %w", obj.Id(), err)
	}

	ms.objectives.Store(string(obj.Id()), objJSON)

	for _, rel := range obj.Related() {
		switch ch := rel.(type) {
		case *channel.Channel:
			err := ms.SetChannel(ch)
			if err != nil {
				return fmt.Errorf("error setting channel %s from objective %s: %w", ch.Id, obj.Id(), err)
			}
		case *consensus_channel.ConsensusChannel:
			err := ms.SetLedgerChannel(ch)
			if err != nil {
				return fmt.Errorf("error setting consensus channel %s from objective %s: %w", ch.Id, obj.Id(), err)
			}
		default:
			return fmt.Errorf("unexpected type: %T", rel)
		}
	}

	// Objective ownership can only be transferred if the channel is not owned by another objective
	prevOwner, isOwned := ms.channelToObjective.Load(obj.OwnsChannel().String())
	if obj.GetStatus() == protocols.Approved {
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
func (ms *MockStore) SetChannel(ch *channel.Channel) error {
	chJSON, err := ch.MarshalJSON()

	if err != nil {
		return err
	}

	ms.channels.Store(ch.Id.String(), chJSON)
	return nil
}

// SetLedgerChannel sets the channel in the store.
func (ms *MockStore) SetLedgerChannel(ch *consensus_channel.ConsensusChannel) error {
	chJSON, err := ch.MarshalJSON()

	if err != nil {
		return err
	}

	ms.consensusChannels.Store(ch.Id.String(), chJSON)
	return nil
}

// GetChannelById retrieves the channel with the supplied id, if it exists.
func (ms *MockStore) GetChannelById(id types.Destination) (c *channel.Channel, ok bool) {
	ch, err := ms.getChannelById(id)

	if err != nil {
		return &channel.Channel{}, false
	}

	return &ch, true
}

// getChannelById returns the stored channel
func (ms *MockStore) getChannelById(id types.Destination) (channel.Channel, error) {
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

// GetTwoPartyLedger returns a ledger channel between the two parties if it exists.
func (ms *MockStore) GetTwoPartyLedger(firstParty types.Address, secondParty types.Address) (*channel.TwoPartyLedger, bool) {
	var ledger *channel.TwoPartyLedger
	var ok bool

	ms.channels.Range(func(key string, chJSON []byte) bool {

		var ch channel.Channel
		err := json.Unmarshal(chJSON, &ch)

		if err != nil {
			return true // channel not found, continue looking
		}

		if len(ch.Participants) == 2 {
			// TODO: Should order matter?
			if ch.Participants[0] == firstParty && ch.Participants[1] == secondParty {
				ledger = &channel.TwoPartyLedger{Channel: ch}
				ok = true
				return false // we have found the target channel: break the Range loop
			}
		}

		return true // channel not found: continue looking
	})

	return ledger, ok
}

// GetConsensusChannelById returns a ConsensusChannel with the given channel id
func (ms *MockStore) GetConsensusChannelById(id types.Destination) (channel *consensus_channel.ConsensusChannel, err error) {

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

// GetLedgerChannel returns a ConsensusChannel between the calling client and
// the supplied counterparty, if such channel exists
func (ms *MockStore) GetLedgerChannel(counterparty types.Address) (channel *consensus_channel.ConsensusChannel, ok bool) {

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

func (ms *MockStore) GetObjectiveByChannelId(channelId types.Destination) (protocols.Objective, bool) {
	// todo: locking
	id, found := ms.channelToObjective.Load(channelId.String())
	if !found {
		return &directfund.Objective{}, false
	}

	objective, err := ms.GetObjectiveById(protocols.ObjectiveId(id))
	return objective, err == nil
}

// populateChannelData fetches stored Channel data relevent to the given
// objective and attaches it to the objective. The channel data is attached
// in-place of the objectives existing channel pointers.
func (ms *MockStore) populateChannelData(obj protocols.Objective) error {
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
		o.V = &channel.SingleHopVirtualChannel{Channel: v}

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
	default:
		return nil, fmt.Errorf("objective id %s does not correspond to a known Objective type", id)

	}
}

func (ms *MockStore) ReleaseChannelFromOwnership(channelId types.Destination) {
	ms.channelToObjective.m.Delete(channelId.String())
}
