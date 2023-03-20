package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/channel"
	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/crypto"
	"github.com/statechannels/go-nitro/payments"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/protocols/directdefund"
	"github.com/statechannels/go-nitro/protocols/directfund"
	"github.com/statechannels/go-nitro/protocols/virtualdefund"
	"github.com/statechannels/go-nitro/protocols/virtualfund"
	"github.com/statechannels/go-nitro/types"
	"github.com/tidwall/buntdb"
)

type PersistStore struct {
	objectives         *buntdb.DB
	channels           *buntdb.DB
	consensusChannels  *buntdb.DB
	channelToObjective *buntdb.DB
	vouchers           *buntdb.DB

	key     string // the signing key of the store's engine
	address string // the (Ethereum) address associated to the signing key
	folder  string // the folder where the store's data is stored
}

// NewPersistStore creates a new PersistStore that uses the given folder to store its data
// It will create the folder if it does not exist
func NewPersistStore(key []byte, folder string, config buntdb.Config) Store {
	ps := PersistStore{}
	err := os.MkdirAll(folder, os.ModePerm)
	ps.checkError(err)

	ps.key = common.Bytes2Hex(key)
	ps.address = crypto.GetAddressFromSecretKeyBytes(key).String()
	ps.folder = folder

	ps.objectives = ps.openDB("objectives", config)
	ps.channels = ps.openDB("channels", config)
	ps.consensusChannels = ps.openDB("consensus_channels", config)
	ps.channelToObjective = ps.openDB("channel_to_objective", config)
	ps.vouchers = ps.openDB("vouchers", config)

	return &ps
}

func (ps *PersistStore) openDB(name string, config buntdb.Config) *buntdb.DB {
	db, err := buntdb.Open(fmt.Sprintf("%s/%s_%s.db", ps.folder, name, ps.address[2:7]))
	ps.checkError(err)
	err = db.SetConfig(config)
	ps.checkError(err)
	return db
}

func (ps *PersistStore) Close() error {
	err := ps.channels.Close()
	if err != nil {
		return err
	}
	err = ps.objectives.Close()
	if err != nil {
		return err
	}
	err = ps.consensusChannels.Close()
	if err != nil {
		return err
	}
	return ps.channelToObjective.Close()
}
func (ps *PersistStore) GetAddress() *types.Address {
	address := common.HexToAddress(ps.address)
	return &address
}

func (ps *PersistStore) GetChannelSecretKey() *[]byte {
	val := common.Hex2Bytes(ps.key)
	return &val
}

func (ps *PersistStore) GetObjectiveById(id protocols.ObjectiveId) (protocols.Objective, error) {

	var obj protocols.Objective
	err := ps.objectives.View(func(tx *buntdb.Tx) error {
		objJSON, err := tx.Get(string(id))
		if err != nil {
			return err
		}

		obj, err = decodeObjective(id, []byte(objJSON))
		if err != nil {
			return fmt.Errorf("error decoding objective %s: %w", id, err)
		}

		err = ps.populateChannelData(obj)
		if err != nil {
			// return existing objective data along with error
			return fmt.Errorf("error populating channel data for objective %s: %w", id, err)
		}
		return nil

	})
	if err != nil && errors.Is(err, buntdb.ErrNotFound) {
		return nil, ErrNoSuchObjective
	}

	return obj, nil
}

func (ps *PersistStore) SetObjective(obj protocols.Objective) error {
	// todo: locking
	objJSON, err := obj.MarshalJSON()

	if err != nil {
		return fmt.Errorf("error setting objective %s: %w", obj.Id(), err)
	}

	err = ps.objectives.Update(func(tx *buntdb.Tx) error {

		_, _, err := tx.Set(string(obj.Id()), string(objJSON), nil)
		return err
	})

	if err != nil {
		return err
	}
	for _, rel := range obj.Related() {
		switch ch := rel.(type) {
		case *channel.Channel:
			err := ps.SetChannel(ch)
			if err != nil {
				return fmt.Errorf("error setting channel %s from objective %s: %w", ch.Id, obj.Id(), err)
			}
		case *consensus_channel.ConsensusChannel:
			err := ps.SetConsensusChannel(ch)
			if err != nil {
				return fmt.Errorf("error setting consensus channel %s from objective %s: %w", ch.Id, obj.Id(), err)
			}
		default:
			return fmt.Errorf("unexpected type: %T", rel)
		}
	}

	// Objective ownership can only be transferred if the channel is not owned by another objective
	var prevOwner protocols.ObjectiveId
	var isOwned bool = false
	err = ps.channelToObjective.View(func(tx *buntdb.Tx) error {
		res, err := tx.Get(string(obj.OwnsChannel().String()))
		if err != nil {
			return nil
		}
		prevOwner = protocols.ObjectiveId(res)
		isOwned = true
		return nil
	})
	if err != nil {
		return err
	}

	if status := obj.GetStatus(); status == protocols.Approved {
		if !isOwned {
			err := ps.channelToObjective.Update(func(tx *buntdb.Tx) error {

				_, _, err := tx.Set(string(obj.OwnsChannel().String()), string(obj.Id()), nil)
				return err
			})
			ps.checkError(err)
		}
		if isOwned && prevOwner != obj.Id() {
			return fmt.Errorf("cannot transfer ownership of channel to from objective %s to %s", prevOwner, obj.Id())
		}
	}

	return nil
}

// SetChannel sets the channel in the store.
func (ps *PersistStore) SetChannel(ch *channel.Channel) error {
	chJSON, err := ch.MarshalJSON()

	if err != nil {
		return err
	}

	err = ps.channels.Update(func(tx *buntdb.Tx) error {

		_, _, err := tx.Set(ch.Id.String(), string(chJSON), nil)
		return err
	})
	return err
}

// DestroyChannel deletes the channel with id id.
func (ps *PersistStore) DestroyChannel(id types.Destination) {
	err := ps.channels.Update(func(tx *buntdb.Tx) error {
		_, err := tx.Delete(id.String())
		return err
	})
	ps.checkError(err)
}

// SetConsensusChannel sets the channel in the store.
func (ps *PersistStore) SetConsensusChannel(ch *consensus_channel.ConsensusChannel) error {
	chJSON, err := ch.MarshalJSON()

	if err != nil {
		return err
	}

	err = ps.consensusChannels.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(ch.Id.String(), string(chJSON), nil)
		return err
	})

	return err
}

// DestroyChannel deletes the channel with id id.
func (ps *PersistStore) DestroyConsensusChannel(id types.Destination) {
	err := ps.consensusChannels.Update(func(tx *buntdb.Tx) error {
		_, err := tx.Delete(id.String())
		return err
	})
	ps.checkError(err)
}

// GetChannelById retrieves the channel with the supplied id, if it exists.
func (ps *PersistStore) GetChannelById(id types.Destination) (c *channel.Channel, ok bool) {
	ch, err := ps.getChannelById(id)

	if err != nil {
		return &channel.Channel{}, false
	}

	return &ch, true
}

// getChannelById returns the stored channel
func (ps *PersistStore) getChannelById(id types.Destination) (channel.Channel, error) {
	var chJSON string
	err := ps.channels.View(func(tx *buntdb.Tx) error {
		var err error
		chJSON, err = tx.Get(id.String())
		return err
	})

	if errors.Is(err, buntdb.ErrNotFound) {

		return channel.Channel{}, ErrNoSuchChannel
	}
	var ch channel.Channel
	err = ch.UnmarshalJSON([]byte(chJSON))

	if err != nil {
		return channel.Channel{}, fmt.Errorf("error unmarshaling channel %s", ch.Id)
	}

	return ch, nil
}

// GetChannelsByParticipant returns any channels that include the given participant
func (ps *PersistStore) GetChannelsByParticipant(participant types.Address) []*channel.Channel {
	toReturn := []*channel.Channel{}
	err := ps.channels.View(func(tx *buntdb.Tx) error {
		err := tx.Ascend("", func(key, chJSON string) bool {

			var ch channel.Channel
			err := json.Unmarshal([]byte(chJSON), &ch)

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
		return err
	})
	ps.checkError(err)
	return toReturn
}

// GetConsensusChannelById returns a ConsensusChannel with the given channel id
func (ps *PersistStore) GetConsensusChannelById(id types.Destination) (channel *consensus_channel.ConsensusChannel, err error) {

	var ch *consensus_channel.ConsensusChannel
	err = ps.consensusChannels.View(func(tx *buntdb.Tx) error {

		chJSON, err := tx.Get(id.String())

		if errors.Is(err, buntdb.ErrNotFound) {
			return ErrNoSuchChannel
		}

		ch = &consensus_channel.ConsensusChannel{}
		err = ch.UnmarshalJSON([]byte(chJSON))

		if err != nil {
			return fmt.Errorf("error unmarshaling channel %s", ch.Id)
		}
		return err
	})
	if err != nil {
		return nil, err
	}
	return ch, nil
}

// GetConsensusChannel returns a ConsensusChannel between the calling client and
// the supplied counterparty, if such channel exists
func (ps *PersistStore) GetConsensusChannel(counterparty types.Address) (channel *consensus_channel.ConsensusChannel, ok bool) {

	err := ps.consensusChannels.View(func(tx *buntdb.Tx) error {
		return tx.Ascend("", func(key, chJSON string) bool {

			var ch consensus_channel.ConsensusChannel
			err := json.Unmarshal([]byte(chJSON), &ch)

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
	})
	ps.checkError(err)
	return
}

func (ps *PersistStore) GetObjectiveByChannelId(channelId types.Destination) (protocols.Objective, bool) {
	var id protocols.ObjectiveId

	err := ps.channelToObjective.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get(channelId.String())
		id = protocols.ObjectiveId(val)

		return err
	})

	if err != nil {
		return &directfund.Objective{}, false
	}

	objective, err := ps.GetObjectiveById(protocols.ObjectiveId(id))
	return objective, err == nil
}

// populateChannelData fetches stored Channel data relevant to the given
// objective and attaches it to the objective. The channel data is attached
// in-place of the objectives existing channel pointers.
func (ps *PersistStore) populateChannelData(obj protocols.Objective) error {
	id := obj.Id()

	switch o := obj.(type) {
	case *directfund.Objective:
		ch, err := ps.getChannelById(o.C.Id)

		if err != nil {
			return fmt.Errorf("error retrieving channel data for objective %s: %w", id, err)
		}

		o.C = &ch

		return nil
	case *directdefund.Objective:

		ch, err := ps.getChannelById(o.C.Id)

		if err != nil {
			return fmt.Errorf("error retrieving channel data for objective %s: %w", id, err)
		}

		o.C = &ch

		return nil
	case *virtualfund.Objective:
		v, err := ps.getChannelById(o.V.Id)
		if err != nil {
			return fmt.Errorf("error retrieving virtual channel data for objective %s: %w", id, err)
		}
		o.V = &channel.VirtualChannel{Channel: v}

		zeroAddress := types.Destination{}

		if o.ToMyLeft != nil &&
			o.ToMyLeft.Channel != nil &&
			o.ToMyLeft.Channel.Id != zeroAddress {

			left, err := ps.GetConsensusChannelById(o.ToMyLeft.Channel.Id)
			if err != nil {
				return fmt.Errorf("error retrieving left ledger channel data for objective %s: %w", id, err)
			}
			o.ToMyLeft.Channel = left
		}

		if o.ToMyRight != nil &&
			o.ToMyRight.Channel != nil &&
			o.ToMyRight.Channel.Id != zeroAddress {
			right, err := ps.GetConsensusChannelById(o.ToMyRight.Channel.Id)
			if err != nil {
				return fmt.Errorf("error retrieving right ledger channel data for objective %s: %w", id, err)
			}
			o.ToMyRight.Channel = right
		}

		return nil
	case *virtualdefund.Objective:

		zeroAddress := types.Destination{}

		if o.ToMyLeft != nil &&
			o.ToMyLeft.Id != zeroAddress {

			left, err := ps.GetConsensusChannelById(o.ToMyLeft.Id)
			if err != nil {
				return fmt.Errorf("error retrieving left ledger channel data for objective %s: %w", id, err)
			}
			o.ToMyLeft = left
		}

		if o.ToMyRight != nil &&
			o.ToMyRight.Id != zeroAddress {
			right, err := ps.GetConsensusChannelById(o.ToMyRight.Id)
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

func (ps *PersistStore) ReleaseChannelFromOwnership(channelId types.Destination) {
	err := ps.channelToObjective.Update(func(tx *buntdb.Tx) error {
		_, err := tx.Delete(channelId.String())
		return err
	})
	ps.checkError(err)
}

// checkError is a helper function that panics if an error is not nil
// TODO: Longer term we should return errors instead of panicking
func (ps *PersistStore) checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func (ps *PersistStore) SetVoucherInfo(channelId types.Destination, v payments.VoucherInfo) error {
	return ps.vouchers.Update(func(tx *buntdb.Tx) error {
		vJSON, err := json.Marshal(v)
		if err != nil {
			return err
		}
		_, _, err = tx.Set(channelId.String(), string(vJSON), nil)

		return err
	})
}

func (ps *PersistStore) GetVoucherInfo(channelId types.Destination) (v *payments.VoucherInfo, ok bool) {
	err := ps.vouchers.View(func(tx *buntdb.Tx) error {
		vJSON, err := tx.Get(channelId.String())
		if err != nil {
			return nil
		}
		return json.Unmarshal([]byte(vJSON), &v)
	})
	if err == nil {
		ok = true
	}
	return
}

func (ps *PersistStore) RemoveVoucherInfo(channelId types.Destination) error {
	return ps.vouchers.Update(func(tx *buntdb.Tx) error {
		_, err := tx.Delete(channelId.String())
		return err
	})
}
