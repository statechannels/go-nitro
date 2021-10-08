package main

import (
	"errors"
	"fmt"
)

// NitroWallet is a state channel wallet that runs a toy model of nitro protocol.
// It is a store of channels with a set of go channels it can use to communicate with other NitroWallets.
// It exposes a toy API to consuming applications.
type NitroWallet struct {
	id              uint
	ledgerChannels  map[uint]LedgerChannelState
	virtualChannels map[uint]VirtualChannelState
	ledgerInbox     *map[uint]chan LedgerChannelState  // a channel for each peer. only listen on your own channel!
	paymentInbox    *map[uint]chan VirtualChannelState // a channel for each peer. only listen on your own channel!
	isHub           bool                               // when set, the wallet may behave differently (generally it performs more actions automatically)
}

func (w *NitroWallet) ProposeAVirtualChannel(peer uint, hub uint) error {
	cId := ChannelId(w.id, peer)

	// check: if virtual channel already exists...
	if w.virtualChannels[peer] != (VirtualChannelState{}) {
		return errors.New(`Virtual Channel already exists with that peer`) // TODO interpolate peerId
	}

	// write: virtual channel to store
	w.virtualChannels[peer] = VirtualChannelState{w.id, peer, 5, 5, 0}

	// read: ledger channel from the store
	ledger := w.ledgerChannels[hub]
	virtualChannelBal := ledger.virtualChannelBal

	// modify: reallocate 10 to a virtual channel with a single peer
	virtualChannelBal[cId] += 10
	ledger.hubBal -= 5
	ledger.leafBal -= 5
	ledger.virtualChannelBal = virtualChannelBal

	// write: ledger channel to my store
	w.ledgerChannels[hub] = ledger

	// send: to hub for signing
	(*w.ledgerInbox)[hub] <- ledger

	return nil
}

func (w *NitroWallet) ListenAndCountersignLedgerUpdates() {
	var messagesHandled = 0
	// Listen for ledgerUpdates. Countersign blindly! TODO
	for {
		select {
		case update := <-(*w.ledgerInbox)[w.id]:
			update.signedByHub = true
			w.ledgerChannels[update.leafId] = update
			fmt.Printf("ledger between hub and leaf %v updated to %v\n", update.leafId, update)
			// send back
			(*w.ledgerInbox)[update.leafId] <- update // this will block unless there is something receiving on this channel or if the channel is buffered.
			messagesHandled++
			fmt.Printf("%v messages handled\n", messagesHandled)

		}
	}
}

func (w *NitroWallet) MakePayment(peer uint) error {

	// check: if virtual channel doesn't exist, fail
	if w.virtualChannels[peer] == (VirtualChannelState{}) {
		return errors.New(`No virtual channel exists with that peer`) // TODO interpolate peerId
	}

	// read: virtual channel from the store
	virtualChannel := w.virtualChannels[peer]

	// modify: reallocate some money from me to my counterparty
	if virtualChannel.proposerId == w.id {
		virtualChannel.proposerBal -= payment_amount
		virtualChannel.joinerBal += payment_amount
	} else if virtualChannel.joinerId == w.id {
		virtualChannel.proposerBal += payment_amount
		virtualChannel.joinerBal -= payment_amount
	}

	// write: virtual channel to store
	w.virtualChannels[peer] = virtualChannel

	// send: to counterparty for signing
	var counterparty uint
	if virtualChannel.proposerId == w.id {
		counterparty = virtualChannel.joinerId
	} else if virtualChannel.joinerId == w.id {
		counterparty = virtualChannel.proposerId
	}
	(*w.paymentInbox)[counterparty] <- virtualChannel

	return nil
}

// updateAVirtualChannel := func() {}
// closeAVirtualChannel := func() {}

// NewNitroWallet creates a new NitroWallet, initializes its store and listening routines, and returns it.
func NewNitroWallet(
	id uint,
	ledgerInbox *map[uint]chan LedgerChannelState,
	paymentInbox *map[uint]chan VirtualChannelState,
	isHub bool,
) *NitroWallet {
	w := NitroWallet{
		id,
		make(map[uint]LedgerChannelState),
		make(map[uint]VirtualChannelState),
		ledgerInbox,
		paymentInbox,
		isHub,
	}

	// manual setup for ledger channels
	// TODO avoid hardcoding hub id?
	if isHub {
		// a ledger channel for each leaf
		for id := uint(1); id < num_leaves+1; id++ {
			w.ledgerChannels[id] = LedgerChannelState{0, ledger_hub_balance, id, ledger_leaf_balance, 1, make(map[string]uint), true, true} // turnNum = 1 so channels are funded
		}
		go w.ListenAndCountersignLedgerUpdates()
	} else {
		// a single ledger channel with the hub
		w.ledgerChannels[0] = LedgerChannelState{0, ledger_hub_balance, id, ledger_leaf_balance, 1, make(map[string]uint), true, true} // turnNum = 1 so channels are funded
	}

	return &w
}
