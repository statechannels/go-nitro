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
	listeners *safesync.Map[*listeners]
	store     store.Store
	vm        *payments.VoucherManager
}

// NewChannelNotifier constructs a channel notifier using the provided store.
func NewChannelNotifier(store store.Store, vm *payments.VoucherManager) *ChannelNotifier {
	return &ChannelNotifier{listeners: &safesync.Map[*listeners]{}, store: store, vm: vm}
}

// RegisterForLedgerUpdates returns a buffered channel that will receive ledger channel updates when they occur.
func (cn *ChannelNotifier) RegisterForLedgerUpdates(cId types.Destination) (<-chan query.LedgerChannelInfo, error) {
	return register[query.LedgerChannelInfo](cn, cId)
}

// RegisterForLedgerUpdates returns a buffered channel that will receive ledger channel updates when they occur.
func (cn *ChannelNotifier) RegisterForPaymentChannelUpdates(cId types.Destination) (<-chan query.PaymentChannelInfo, error) {
	return register[query.PaymentChannelInfo](cn, cId)
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
	// If no one has registered for this channel, we don't need to notify anyone.
	if !ok {
		return nil
	}

	// We only want to notify listeners if the ledger channel has changed from the perspective of the client.
	if ledgerUpdated := li.prevLedger == nil || li.prevLedger.Equal(latest); ledgerUpdated {

		err = li.NotifyLedger(latest)
		if err != nil {
			return err
		}
		cn.listeners.Store(lId.String(), li)
	}
	return nil
}

// NotifyPaymentUpdated notifies all listeners of a payment channel update.
// It will query the store for the latest payment channel info and output an event to listeners if the payment channel has changed.
// NOTE: NotifyPaymentUpdated is dependent on the current state of the store, so must be called before the store is updated.
func (cn *ChannelNotifier) NotifyPaymentUpdated(pId types.Destination) error {
	// Fetch the current state of the payment channel
	latest, err := query.GetPaymentChannelInfo(pId, cn.store, cn.vm)
	if err != nil {
		return err
	}
	// Fetch the listeners for the ledger channel
	li, ok := cn.listeners.Load(pId.String())
	// If no one has registered for this channel, we don't need to notify anyone.
	if !ok {
		return nil
	}

	// We only want to notify listeners if the payment channel has changed from the perspective of the client.
	if channelUpdated := li.prevPayment == nil || li.prevPayment.Equal(latest); channelUpdated {

		err = li.NotifyPayment(latest)
		if err != nil {
			return err
		}
		cn.listeners.Store(pId.String(), li)
	}
	return nil
}

func register[T query.ChannelInfo](cn *ChannelNotifier, cId types.Destination) (<-chan T, error) {
	li, _ := cn.listeners.LoadOrStore(cId.String(), newListener(cId))

	lChan, err := addListener[T](li)
	if err != nil {
		return nil, err
	}

	cn.listeners.Store(cId.String(), li)
	return lChan, nil
}
