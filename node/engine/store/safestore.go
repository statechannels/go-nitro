package store

import (
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/channel"
	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/internal/safesync"
	"github.com/statechannels/go-nitro/payments"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

type SafeStore struct {
	store   Store
	mutexes safesync.Map[*sync.Mutex]
}

func NewSafeStore(key []byte) Store {
	return &SafeStore{
		store:   NewMemStore(key),
		mutexes: safesync.Map[*sync.Mutex]{},
	}
}

func (ss *SafeStore) GetChannelSecretKey() *[]byte { return ss.store.GetChannelSecretKey() }
func (ss *SafeStore) GetAddress() *types.Address   { return (*common.Address)(ss.store.GetAddress()) }

// Objective store
func (ss *SafeStore) GetObjectiveById(id protocols.ObjectiveId) (protocols.Objective, error) {
	mu, _ := ss.mutexes.LoadOrStore(string(id), &(sync.Mutex{}))
	mu.Lock()
	defer mu.Unlock()
	return ss.store.GetObjectiveById(id)
}

func (ss *SafeStore) GetObjectiveByChannelId(id types.Destination) (obj protocols.Objective, ok bool) {
	return ss.store.GetObjectiveByChannelId(id)
}

func (ss *SafeStore) SetObjective(o protocols.Objective) error {
	mu, _ := ss.mutexes.LoadOrStore(string(o.Id()), &(sync.Mutex{}))
	mu.Lock()
	defer mu.Unlock()
	return ss.store.SetObjective(o)
}

// Channel store
func (ss *SafeStore) GetChannelsByIds(ids []types.Destination) ([]*channel.Channel, error) {
	for _, id := range ids {
		mu, _ := ss.mutexes.LoadOrStore(id.String(), &(sync.Mutex{}))
		mu.Lock()
		defer mu.Unlock()
	}
	return ss.store.GetChannelsByIds(ids)
}

func (ss *SafeStore) GetChannelById(id types.Destination) (c *channel.Channel, ok bool) {
	mu, _ := ss.mutexes.LoadOrStore(id.String(), &(sync.Mutex{}))
	mu.Lock()
	defer mu.Unlock()
	return ss.store.GetChannelById(id)
}

func (ss *SafeStore) GetChannelsByParticipant(participant types.Address) ([]*channel.Channel, error) {
	// TODO lock all mutexes??
	return ss.store.GetChannelsByParticipant(participant)
}

func (ss *SafeStore) SetChannel(c *channel.Channel) error {
	mu, _ := ss.mutexes.LoadOrStore(c.Id.String(), &(sync.Mutex{}))
	mu.Lock()
	defer mu.Unlock()
	return ss.store.SetChannel(c)
}
func (ss *SafeStore) DestroyChannel(id types.Destination) error { return ss.store.DestroyChannel(id) }

func (ss *SafeStore) GetChannelsByAppDefinition(appDef types.Address) ([]*channel.Channel, error) {
	return ss.store.GetChannelsByAppDefinition(appDef)
}

func (ss *SafeStore) ReleaseChannelFromOwnership(id types.Destination) error {
	mu, _ := ss.mutexes.LoadOrStore(id.String(), &(sync.Mutex{}))
	mu.Lock()
	defer mu.Unlock()
	return ss.store.ReleaseChannelFromOwnership(id)
}

// Consensus Channel Store

func (ss *SafeStore) GetAllConsensusChannels() ([]*consensus_channel.ConsensusChannel, error) {
	// TODO lock all mutexes??
	return ss.store.GetAllConsensusChannels()
}

func (ss *SafeStore) GetConsensusChannel(counterparty types.Address) (channel *consensus_channel.ConsensusChannel, ok bool) {
	// TODO lock all mutexes??
	return ss.store.GetConsensusChannel(counterparty)
}

func (ss *SafeStore) GetConsensusChannelById(id types.Destination) (channel *consensus_channel.ConsensusChannel, err error) {
	mu, _ := ss.mutexes.LoadOrStore(id.String(), &(sync.Mutex{}))
	mu.Lock()
	defer mu.Unlock()
	return ss.store.GetConsensusChannelById(id)
}

func (ss *SafeStore) SetConsensusChannel(c *consensus_channel.ConsensusChannel) error {
	mu, _ := ss.mutexes.LoadOrStore(c.Id.String(), &(sync.Mutex{}))
	mu.Lock()
	defer mu.Unlock()
	return ss.store.SetConsensusChannel(c)
}

func (ss *SafeStore) DestroyConsensusChannel(id types.Destination) error {
	mu, _ := ss.mutexes.LoadOrStore(id.String(), &(sync.Mutex{}))
	mu.Lock()
	defer mu.Unlock()
	return ss.store.DestroyConsensusChannel(id)
}

func (ss *SafeStore) Close() error { return ss.store.Close() }

// Voucher store
func (ss *SafeStore) SetVoucherInfo(channelId types.Destination, v payments.VoucherInfo) error {
	return ss.store.SetVoucherInfo(channelId, v)
}

func (ss *SafeStore) GetVoucherInfo(channelId types.Destination) (v *payments.VoucherInfo, err error) {
	return ss.store.GetVoucherInfo(channelId)
}

func (ss *SafeStore) RemoveVoucherInfo(channelId types.Destination) error {
	return ss.store.RemoveVoucherInfo(channelId)
}

// Updating an Objective

func (ss *SafeStore) GetAndLockObjectiveById(id protocols.ObjectiveId) (protocols.Objective, error) {
	mu, _ := ss.mutexes.LoadOrStore(string(id), &(sync.Mutex{}))
	mu.Lock()
	// DO NOT UNLOCK MUTEX
	return ss.store.GetObjectiveById(id)
}

func (ss *SafeStore) SetAndUnlockObjective(o protocols.Objective) error {
	mu, _ := ss.mutexes.LoadOrStore(string(o.Id()), &(sync.Mutex{}))
	defer mu.Unlock() // WILL PANIC IF THE MUTEX IS ALREADY UNLOCKED
	return ss.store.SetObjective(o)
}
