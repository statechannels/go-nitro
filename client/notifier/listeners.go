package notifier

import (
	"sync"

	"github.com/statechannels/go-nitro/client/query"
)

func newPaymentChannelListeners() *paymentChannelListeners {
	return &paymentChannelListeners{listeners: []chan query.PaymentChannelInfo{}, listenersLock: sync.Mutex{}}
}

type paymentChannelListeners struct {
	listeners []chan query.PaymentChannelInfo
	prev      query.PaymentChannelInfo
	// listenersLock is used to protect against concurrent access to the listeners slice.
	listenersLock sync.Mutex
}

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

func (li *paymentChannelListeners) getListener(index int) <-chan query.PaymentChannelInfo {
	return li.listeners[index]
}

func (li *paymentChannelListeners) createListener() <-chan query.PaymentChannelInfo {
	li.listenersLock.Lock()
	defer li.listenersLock.Unlock()
	// Use a buffered channel to avoid blocking the notifier.
	listener := make(chan query.PaymentChannelInfo, 100)
	li.listeners = append(li.listeners, listener)
	return listener
}

func newLedgerChannelListeners() *ledgerChannelListeners {
	return &ledgerChannelListeners{listeners: []chan query.LedgerChannelInfo{}, listenersLock: sync.Mutex{}}
}

type ledgerChannelListeners struct {
	listeners []chan query.LedgerChannelInfo
	prev      query.LedgerChannelInfo
	// listenersLock is used to protect against concurrent access to the listeners slice.
	listenersLock sync.Mutex
}

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

func (li *ledgerChannelListeners) createListener() <-chan query.LedgerChannelInfo {
	li.listenersLock.Lock()
	defer li.listenersLock.Unlock()
	// Use a buffered channel to avoid blocking the notifier.
	listener := make(chan query.LedgerChannelInfo, 100)
	li.listeners = append(li.listeners, listener)
	return listener
}

func (li *ledgerChannelListeners) Close() error {
	li.listenersLock.Lock()
	defer li.listenersLock.Unlock()
	for _, c := range li.listeners {
		close(c)
	}

	return nil
}

func (li *paymentChannelListeners) Close() error {
	li.listenersLock.Lock()
	defer li.listenersLock.Unlock()
	for _, c := range li.listeners {
		close(c)
	}

	return nil
}

func (li *ledgerChannelListeners) getListener(index int) <-chan query.LedgerChannelInfo {
	return li.listeners[index]
}
