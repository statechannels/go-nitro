package main

import (
	"math/rand"
	"time"
)

// LeafRunner is an function which simulates an application running a NitroWallet in "leaf" mode.
// It makes API calls on a configurable (e.g. random and/or periodic) schedule.
func LeafRunner(
	id uint, // Network id for this node
	ledgerInbox *map[uint]chan LedgerChannelState, // Map of inbox go channels for ledger channel updates
	paymentInbox *map[uint]chan VirtualChannelState, // Map of inbox go channels for vitual channel updates
) {
	{
		w := NewNitroWallet(id, ledgerInbox, paymentInbox, false)

		proposeTicker := time.NewTicker(leaf_propose_period)
		paymentTicker := time.NewTicker(leaf_payment_period)
		for {
			select {
			case <-proposeTicker.C: // On each tick
				randomPeer := uint(rand.Intn(num_leaves))
				w.ProposeAVirtualChannel(randomPeer, 0) // Propose a virtual channel with a random peer through hub 0

			case <-paymentTicker.C: // On each tick
				randomPeer := uint(rand.Intn(num_leaves))
				error := w.MakePayment(randomPeer) // Make a payment with a random peer
				if error != nil {
					// update stats here
				}
			}
		}
	}

}
