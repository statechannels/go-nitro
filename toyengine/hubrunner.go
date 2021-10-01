package main

import "fmt"

func HubRunner(hubId uint, ledgerInbox *map[uint]chan LedgerChannelState) {
	var store = LedgerStore{make(map[uint]LedgerChannelState)}
	var messagesHandled = 0

	// Manual setup for ledger channels
	for id := uint(1); id < num_leaves+1; id++ {
		store.ledgerChannels[id] = LedgerChannelState{hubId, ledger_hub_balance, id, ledger_leaf_balance, 1, make(map[string]uint), true, true} // turnNum = 1 so channels are funded
	}

	// Listen for ledgerUpdates. Countersign blindly! TODO
	for {
		select {
		case update := <-(*ledgerInbox)[hubId]:
			update.signedByHub = true
			store.ledgerChannels[update.leafId] = update
			fmt.Printf("ledger between hub and leaf %v updated to %v\n", update.leafId, update)
			// send back
			(*ledgerInbox)[update.leafId] <- update // this will block unless there is something receiving on this channel or if the channel is buffered.
			messagesHandled++
			fmt.Printf("%v messages handled\n", messagesHandled)

		}
	}
}
