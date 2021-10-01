package main

import (
	"math/rand"
	"time"
)

func LeafRunner(id uint, ledgerInbox *map[uint]chan LedgerChannelState) {
	w := NewNitroWallet(id, ledgerInbox)

	proposeTicker := time.NewTicker(leaf_propose_period)
	for {
		select {
		case <-proposeTicker.C:
			randomPeer := uint(rand.Intn(num_leaves))
			w.ProposeAVirtualChannel(randomPeer, 0)
		}
	}

}
