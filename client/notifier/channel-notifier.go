package notifier

import (
	"fmt"
	"sync"

	"github.com/statechannels/go-nitro/client/engine/store"
	"github.com/statechannels/go-nitro/client/query"
	"github.com/statechannels/go-nitro/internal/safesync"
	"github.com/statechannels/go-nitro/types"
)

// ChannelNotifier is used to notify multiple listeners of a channel update.
type ChannelNotifier struct {
	listeners *safesync.Map[*listeners]
	store     store.Store
}

// NewChannelNotifier constructs a channel notifier using the provided store.
func NewChannelNotifier(store store.Store) *ChannelNotifier {
	return &ChannelNotifier{listeners: &safesync.Map[*listeners]{}, store: store}
}

// RegisterForLedgerUpdates returns a buffered channel that will receive ledger channel updates when they occur.
func (cn *ChannelNotifier) RegisterForLedgerUpdates(cId types.Destination) <-chan query.LedgerChannelInfo {
	li, _ := cn.listeners.LoadOrStore(cId.String(), newListener(cId))

	lChan := li.AddListener()

	cn.listeners.Store(cId.String(), li)
	return lChan
}

// NotifyLedgerUpdated notifies all listeners of a ledger channel update.
// It will query the store for the latest ledger channel info and output an event to listeners if the ledger channel has changed.
// NOTE: NotifyLedgerUpdated is dependent on the current state of the store, so must be called before the store is updated.
func (cn *ChannelNotifier) NotifyLedgerUpdated(lId types.Destination) error {
	// Fetch the current state of the ledger channel
	latest, err := query.GetLedgerChannelInfo(lId, cn.store)
	if err != nil {
		return err
	}
	// Fetch the listeners for the ledger channel
	li, ok := cn.listeners.Load(lId.String())
	if !ok {
		return fmt.Errorf("no listeners for ledger channel %v", lId)
	}

	// We only want to notify listeners if the ledger channel has changed from the perspective of the client.
	if ledgerUpdated := li.prev == nil || li.prev.Equal(latest); ledgerUpdated {

		li.Notify(latest)

		cn.listeners.Store(lId.String(), li)
	}
	return nil
}

// listeners is a struct that holds a list of listeners for ledger channel updates.
type listeners struct {
	// Listeners is a colletion of listeners for ledger channel updates.
	Listeners []chan query.LedgerChannelInfo

	// prev is the previous ledger channel info that was sent to listeners.
	prev *query.LedgerChannelInfo
	// listenersLock is used to protect against concurrent access to the listeners slice.
	listenersLock sync.Mutex
	ledgerId      types.Destination
}

func (li *listeners) Notify(info query.LedgerChannelInfo) {
	li.listenersLock.Lock()
	defer li.listenersLock.Unlock()
	for _, list := range li.Listeners {
		list <- info
	}
	li.prev = &info
}

func (li *listeners) AddListener() chan query.LedgerChannelInfo {
	lChan := make(chan query.LedgerChannelInfo, 100)

	li.listenersLock.Lock()
	defer li.listenersLock.Unlock()

	li.Listeners = append(li.Listeners, lChan)

	return lChan
}

func newListener(ledgerId types.Destination) *listeners {
	return &listeners{Listeners: []chan query.LedgerChannelInfo{}, listenersLock: sync.Mutex{}, ledgerId: ledgerId, prev: nil}
}
