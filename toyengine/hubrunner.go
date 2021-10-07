package main

func HubRunner(hubId uint, ledgerInbox *map[uint]chan LedgerChannelState) {
	NewNitroWallet(hubId, ledgerInbox, true) // we don't need to hit the API just yet. So no need to store this in a variable.
}
