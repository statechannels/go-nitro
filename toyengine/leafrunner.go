package main

import (
	"math/rand"
	"time"
)

// LeafRunner is an function which simulates an application running a NitroWallet in "leaf" mode.
// It makes API calls on a configurable (e.g. random and/or periodic) schedule.
func LeafRunner(id uint, ledgerInbox *map[uint]chan LedgerChannelState) {
	w := NewNitroWallet(id, ledgerInbox, false)

	proposeTicker := time.NewTicker(leaf_propose_period)
	for {
		select {
		case <-proposeTicker.C: // On each tick
			randomPeer := uint(rand.Intn(num_leaves))
			w.ProposeAVirtualChannel(randomPeer, 0) // Propose a virtual channel with a random peer
		}
	}

}
