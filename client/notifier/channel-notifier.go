package notifier

import (
	"github.com/statechannels/go-nitro/client/engine/store"
	"github.com/statechannels/go-nitro/client/query"
	"github.com/statechannels/go-nitro/internal/safesync"
	"github.com/statechannels/go-nitro/payments"
	"github.com/statechannels/go-nitro/types"
)

const ALL_NOTIFICATIONS = "all"

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

// RegisterForAllLedgerUpdates returns a buffered channel that will receive updates for all ledger channels.
func (cn *ChannelNotifier) RegisterForAllLedgerUpdates() <-chan query.LedgerChannelInfo {
	li, _ := cn.ledgerListeners.LoadOrStore(ALL_NOTIFICATIONS, newLedgerChannelListeners())
	return li.getOrCreateListener()
}

// RegisterForLedgerUpdates returns a buffered channel that will receive updates for a specific ledger channel.
func (cn *ChannelNotifier) RegisterForLedgerUpdates(cId types.Destination) <-chan query.LedgerChannelInfo {
	li, _ := cn.ledgerListeners.LoadOrStore(cId.String(), newLedgerChannelListeners())
	return li.createNewListener()
}

// RegisterForAllPaymentUpdates returns a buffered channel that will receive updates for all payment channels.
func (cn *ChannelNotifier) RegisterForAllPaymentUpdates() <-chan query.PaymentChannelInfo {
	li, _ := cn.paymentListeners.LoadOrStore(ALL_NOTIFICATIONS, newPaymentChannelListeners())
	return li.getOrCreateListener()
}

// RegisterForLedgerUpdates returns a buffered channel that will receive updates or a specific payment channel.
func (cn *ChannelNotifier) RegisterForPaymentChannelUpdates(cId types.Destination) <-chan query.PaymentChannelInfo {
	li, _ := cn.paymentListeners.LoadOrStore(cId.String(), newPaymentChannelListeners())
	return li.createNewListener()
}

// NotifyLedgerUpdated notifies all listeners of a ledger channel update.
// It should be called whenever a ledger channel is updated.
func (cn *ChannelNotifier) NotifyLedgerUpdated(info query.LedgerChannelInfo) error {
	li, _ := cn.ledgerListeners.LoadOrStore(info.ID.String(), newLedgerChannelListeners())
	li.Notify(info)
	allLi, _ := cn.ledgerListeners.LoadOrStore(ALL_NOTIFICATIONS, newLedgerChannelListeners())
	allLi.Notify(info)

	return nil
}

// NotifyPaymentUpdated notifies all listeners of a payment channel update.
// It should be called whenever a payment channel is updated.
func (cn *ChannelNotifier) NotifyPaymentUpdated(info query.PaymentChannelInfo) error {
	li, _ := cn.paymentListeners.LoadOrStore(info.ID.String(), newPaymentChannelListeners())
	li.Notify(info)

	allLi, _ := cn.paymentListeners.LoadOrStore(ALL_NOTIFICATIONS, newPaymentChannelListeners())
	allLi.Notify(info)

	return nil
}

// Close closes the notifier and all listeners.
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
