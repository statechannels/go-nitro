package notifier

import (
	"sync"

	"github.com/statechannels/go-nitro/node/query"
)

// paymentChannelListeners is a struct that holds a list of listeners for payment channel info.
type paymentChannelListeners struct {
	// listeners is a list of listeners for payment channel info that we need to notify.
	listeners []chan query.PaymentChannelInfo
	// prev is the previous payment channel info that was sent to the listeners.
	prev query.PaymentChannelInfo
	// listenersLock is used to protect against concurrent access to to sibling struct members.
	listenersLock *sync.Mutex
}

// newPaymentChannelListeners constructs a new payment channel listeners struct.
func newPaymentChannelListeners() *paymentChannelListeners {
	return &paymentChannelListeners{listeners: []chan query.PaymentChannelInfo{}, listenersLock: &sync.Mutex{}}
}

// Notify notifies all listeners of a payment channel update.
// It only notifies listeners if the new info is different from the previous info.
func (li *paymentChannelListeners) Notify(info query.PaymentChannelInfo) {
	li.listenersLock.Lock()
	defer li.listenersLock.Unlock()
	if li.prev.Equal(info) {
		return
	}
	for _, list := range li.listeners {
		list <- info
	}
	li.prev = info
}

// createNewListener creates a new listener and adds it to the list of listeners.
func (li *paymentChannelListeners) createNewListener() <-chan query.PaymentChannelInfo {
	li.listenersLock.Lock()
	defer li.listenersLock.Unlock()
	// Use a buffered channel to avoid blocking the notifier.
	listener := make(chan query.PaymentChannelInfo, 1000)
	li.listeners = append(li.listeners, listener)
	return listener
}

// getOrCreateListener returns the first listener, creating one if none exist.
func (li *paymentChannelListeners) getOrCreateListener() <-chan query.PaymentChannelInfo {
	li.listenersLock.Lock()
	if len(li.listeners) != 0 {
		l := li.listeners[0]
		li.listenersLock.Unlock()
		return l
	}
	li.listenersLock.Unlock()
	return li.createNewListener()
}

// Close closes any active listeners.
func (li *paymentChannelListeners) Close() error {
	li.listenersLock.Lock()
	defer li.listenersLock.Unlock()
	for _, c := range li.listeners {
		close(c)
	}

	return nil
}

// ledgerChannelListeners is a struct that holds a list of listeners for ledger channel info.
type ledgerChannelListeners struct {
	// listeners is a list of listeners for ledger channel info that we need to notify.
	listeners []chan query.LedgerChannelInfo
	// prev is the previous ledger channel info that was sent to the listeners.
	prev query.LedgerChannelInfo
	// listenersLock is used to protect against concurrent access to sibling struct members.
	listenersLock sync.Mutex
}

// newLedgerChannelListeners constructs a new ledger channel listeners struct.
func newLedgerChannelListeners() *ledgerChannelListeners {
	return &ledgerChannelListeners{listeners: []chan query.LedgerChannelInfo{}, listenersLock: sync.Mutex{}}
}

// Notify notifies all listeners of a ledger channel update.
// It only notifies listeners if the new info is different from the previous info.
func (li *ledgerChannelListeners) Notify(info query.LedgerChannelInfo) {
	li.listenersLock.Lock()
	defer li.listenersLock.Unlock()
	if li.prev.Equal(info) {
		return
	}

	for _, list := range li.listeners {
		list <- info
	}
	li.prev = info
}

// createNewListener creates a new listener and adds it to the list of listeners.
func (li *ledgerChannelListeners) createNewListener() <-chan query.LedgerChannelInfo {
	li.listenersLock.Lock()
	defer li.listenersLock.Unlock()
	// Use a buffered channel to avoid blocking the notifier.
	listener := make(chan query.LedgerChannelInfo, 1000)
	li.listeners = append(li.listeners, listener)
	return listener
}

// getOrCreateListener returns the first listener, creating one if none exist.
func (li *ledgerChannelListeners) getOrCreateListener() <-chan query.LedgerChannelInfo {
	li.listenersLock.Lock()
	if len(li.listeners) != 0 {
		l := li.listeners[0]
		li.listenersLock.Unlock()
		return l
	}
	li.listenersLock.Unlock()
	return li.createNewListener()
}

// Close closes all listeners.
func (li *ledgerChannelListeners) Close() error {
	li.listenersLock.Lock()
	defer li.listenersLock.Unlock()
	for _, c := range li.listeners {
		close(c)
	}

	return nil
}
