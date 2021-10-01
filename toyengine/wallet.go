package main

type NitroWallet struct {
	id              uint
	ledgerChannels  map[uint]LedgerChannelState
	virtualChannels map[uint]VirtualChannelState
	ledgerInbox     *map[uint]chan LedgerChannelState // a channel for each peer. only listen on your own channel!
}

func (w *NitroWallet) ProposeAVirtualChannel(peer uint, hub uint) {
	cId := ChannelId(w.id, peer)
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
}

// updateAVirtualChannel := func() {}
// closeAVirtualChannel := func() {}

func NewNitroWallet(id uint, ledgerInbox *map[uint]chan LedgerChannelState) *NitroWallet {
	w := NitroWallet{
		id,
		make(map[uint]LedgerChannelState),
		make(map[uint]VirtualChannelState),
		ledgerInbox,
	}

	// manual setup for ledger channels
	// TODO avoid hardcoding hub id?
	if id != 0 {
		// for hub
		w.ledgerChannels[0] = LedgerChannelState{0, ledger_hub_balance, id, ledger_leaf_balance, 1, make(map[string]uint), true, true} // turnNum = 1 so channels are funded
	} else {
		// for leaves
		for id := uint(1); id < num_leaves+1; id++ {
			w.ledgerChannels[id] = LedgerChannelState{0, ledger_hub_balance, id, ledger_leaf_balance, 1, make(map[string]uint), true, true} // turnNum = 1 so channels are funded
		}
	}
	return &w
}
