package notifier

import (
	"github.com/statechannels/go-nitro/client/engine/store"
	"github.com/statechannels/go-nitro/client/query"
	"github.com/statechannels/go-nitro/internal/safesync"
	"github.com/statechannels/go-nitro/payments"
	"github.com/statechannels/go-nitro/types"
)

// ChannelNotifier is used to notify multiple listeners of a channel update.
type ChannelNotifier struct {
	ledgerListeners  *safesync.Map[*ledgerChannelListeners]
	paymentListeners *safesync.Map[*paymentChannelListeners]
	store            store.Store
	vm               *payments.VoucherManager
}

// NewChannelNotifier constructs a channel notifier using the provided store.
func NewChannelNotifier(store store.Store, vm *payments.VoucherManager) *ChannelNotifier {
	return &ChannelNotifier{
		ledgerListeners:  &safesync.Map[*ledgerChannelListeners]{},
		paymentListeners: &safesync.Map[*paymentChannelListeners]{},
		store:            store,
		vm:               vm,
	}
}

func (cn *ChannelNotifier) RegisterForAllLedgerUpdates() <-chan query.LedgerChannelInfo {
	li, _ := cn.ledgerListeners.LoadOrStore("all", newLedgerChannelListeners())

	newList := li.createListener()
	cn.ledgerListeners.Store("all", li)
	return newList
}

// RegisterForLedgerUpdates returns a buffered channel that will receive ledger channel updates when they occur.
func (cn *ChannelNotifier) RegisterForLedgerUpdates(cId types.Destination) <-chan query.LedgerChannelInfo {
	li, _ := cn.ledgerListeners.LoadOrStore(cId.String(), newLedgerChannelListeners())

	newList := li.createListener()
	cn.ledgerListeners.Store(cId.String(), li)
	return newList
}

func (cn *ChannelNotifier) RegisterForAllPaymentUpdates() <-chan query.PaymentChannelInfo {
	li, loaded := cn.paymentListeners.LoadOrStore("all", newPaymentChannelListeners())
	if loaded {
		// Use the existing listener instead of always spawning a new one
		return li.getListener(0)
	}

	newList := li.createListener()
	cn.paymentListeners.Store("all", li)
	return newList
}

// RegisterForLedgerUpdates returns a buffered channel that will receive ledger channel updates when they occur.
func (cn *ChannelNotifier) RegisterForPaymentChannelUpdates(cId types.Destination) <-chan query.PaymentChannelInfo {
	li, _ := cn.paymentListeners.LoadOrStore(cId.String(), newPaymentChannelListeners())

	newList := li.createListener()
	cn.paymentListeners.Store(cId.String(), li)
	return newList
}

// NotifyLedgerUpdated notifies all listeners of a ledger channel update.

func (cn *ChannelNotifier) NotifyLedgerUpdated(info query.LedgerChannelInfo) error {
	// Fetch the listeners for the ledger channel
	li, _ := cn.ledgerListeners.LoadOrStore(info.ID.String(), newLedgerChannelListeners())
	li.Notify(info)
	cn.ledgerListeners.Store(info.ID.String(), li)

	allLi, ok := cn.ledgerListeners.Load("all")
	if !ok {
		return nil
	}
	allLi.Notify(info)
	cn.ledgerListeners.Store("all", allLi)

	return nil
}

// NotifyPaymentUpdated notifies all listeners of a payment channel update.
// It will query the store for the latest payment channel info and output an event to listeners if the payment channel has changed.
// NOTE: NotifyPaymentUpdated is dependent on the current state of the store, so must be called before the store is updated.
func (cn *ChannelNotifier) NotifyPaymentUpdated(info query.PaymentChannelInfo) error {
	// Fetch the listeners for the payment channel
	li, _ := cn.paymentListeners.LoadOrStore(info.ID.String(), newPaymentChannelListeners())
	li.Notify(info)
	cn.paymentListeners.Store(info.ID.String(), li)

	allLi, _ := cn.paymentListeners.LoadOrStore("all", newPaymentChannelListeners())
	allLi.Notify(info)
	cn.paymentListeners.Store("all", allLi)

	return nil
}

func (cn *ChannelNotifier) Close() error {
	var err error
	cn.ledgerListeners.Range(func(k string, v *ledgerChannelListeners) bool {
		err = v.Close()
		return err == nil
	})
	cn.paymentListeners.Range(func(k string, v *paymentChannelListeners) bool {
		err = v.Close()
		return err == nil
	})
	return err
}
