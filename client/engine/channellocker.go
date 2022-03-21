package engine

import (
	"bytes"
	"sort"
	"sync"

	"github.com/statechannels/go-nitro/channel"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

// ChannelLocker is a utility class that allows for locking of channels to prevent concurrent updates to the same channels.
// It avoids deadlocks by always acquiring channel locks in the same order.
type ChannelLocker struct {
	channelLocks sync.Map
}

// NewChannelLocker returns a new ChannelLocker
func NewChannelLocker() *ChannelLocker {
	return &ChannelLocker{
		channelLocks: sync.Map{},
	}
}

// Lock acquires a locks on all relevant channels for an objective. This will block until a lock on all channels is acquired.
func (l *ChannelLocker) Lock(objective protocols.Objective) {

	channelIds := getChannelIds(objective.Channels())
	// We sort the channel ids to ensure that we always acquire locks in the same order to prevent deadlocks
	sorted := sortChannelIds(channelIds)

	for _, channelId := range sorted {
		result, _ := l.channelLocks.LoadOrStore(channelId, &sync.Mutex{})
		lock := result.(*sync.Mutex)
		lock.Lock()
	}
}

// Unlock releases the lock on the given channels for an objective.
func (l *ChannelLocker) Unlock(objective protocols.Objective) {
	channelIds := getChannelIds(objective.Channels())
	sorted := sortChannelIds(channelIds)

	for _, channelId := range sorted {
		result, _ := l.channelLocks.Load(channelId)
		lock := result.(*sync.Mutex)
		lock.Unlock()
	}
}

// SortChannelIds is a helper function to sort the channel ids.
// This is used to ensure that locks are acquired in the same order.
func sortChannelIds(channelIds []types.Destination) []types.Destination {
	sorted := make([]types.
		Destination, len(channelIds))
	copy(sorted, channelIds)
	sort.Slice(sorted, func(i, j int) bool { return bytes.Compare(channelIds[i].Bytes(), channelIds[j].Bytes()) < 0 })
	return sorted
}

// getChannelIds is a helper function to get the channel ids from a collection of channels.
func getChannelIds(channels []*channel.Channel) []types.Destination {
	channelIds := make([]types.Destination, len(channels))
	for i, channel := range channels {
		channelIds[i] = channel.Id
	}
	return channelIds
}
