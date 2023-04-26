package notifier

import (
	"fmt"
	"sync"

	"github.com/statechannels/go-nitro/client/query"
	"github.com/statechannels/go-nitro/types"
)

// listeners is a struct that holds a list of listeners for ledger channel updates.
type listeners struct {
	// LedgerListeners is a colletion of listeners for ledger channel updates.
	LedgerListeners []chan query.LedgerChannelInfo
	// LedgerListeners is a colletion of listeners for ledger channel updates.
	PaymentListeners []chan query.PaymentChannelInfo
	// prevLedger is the previous ledger channel info that was sent to listeners.
	prevLedger *query.LedgerChannelInfo

	// prevLedger is the previous ledger channel info that was sent to listeners.
	prevPayment *query.PaymentChannelInfo
	// listenersLock is used to protect against concurrent access to the listeners slice.
	listenersLock sync.Mutex
	ledgerId      types.Destination
}

func (li *listeners) NotifyLedger(info query.LedgerChannelInfo) error {
	return notify(li, info)
}

func (li *listeners) NotifyPayment(info query.PaymentChannelInfo) error {
	return notify(li, info)
}

func newListener(ledgerId types.Destination) *listeners {
	return &listeners{LedgerListeners: []chan query.LedgerChannelInfo{}, listenersLock: sync.Mutex{}, ledgerId: ledgerId, prevLedger: nil}
}

// addListener adds a listener to the list of listeners for a ledger or payment channel.
func addListener[T query.ChannelInfo](li *listeners) (chan T, error) {
	lChan := make(chan T, 100)

	li.listenersLock.Lock()
	defer li.listenersLock.Unlock()
	// We have to cast to any to be able to switch on the type of a generic
	switch lType := any(lChan).(type) {
	case chan query.LedgerChannelInfo:
		li.LedgerListeners = append(li.LedgerListeners, any(lChan).(chan query.LedgerChannelInfo))
	case chan query.PaymentChannelInfo:
		li.PaymentListeners = append(li.PaymentListeners, any(lChan).(chan query.PaymentChannelInfo))
	default:
		return nil, fmt.Errorf("unknown listener type %v", lType)
	}
	return lChan, nil
}

// notify notifies all listeners of a ledger or payment channel update.
func notify[T query.ChannelInfo](li *listeners, info T) error {
	li.listenersLock.Lock()
	defer li.listenersLock.Unlock()
	// We have to cast to any to be able to switch on the type of a generic
	switch iType := any(info).(type) {
	case query.LedgerChannelInfo:
		for _, list := range li.LedgerListeners {
			list <- any(info).(query.LedgerChannelInfo)
		}
		prev := any(info).(query.LedgerChannelInfo)
		li.prevLedger = &prev
	case chan query.PaymentChannelInfo:
		for _, list := range li.PaymentListeners {
			list <- any(info).(query.PaymentChannelInfo)
		}
		prev := any(info).(query.PaymentChannelInfo)
		li.prevPayment = &prev
	default:
		return fmt.Errorf("unknown channel type %v", iType)
	}
	return nil
}
