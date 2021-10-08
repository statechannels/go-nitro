package main

// HubRunner is a function which simulates an application running a NitroWallet in "hub" mode
func HubRunner(hubId uint, ledgerInbox *map[uint]chan LedgerChannelState) {
	NewNitroWallet(hubId, ledgerInbox, true) // We don't need to hit the API just yet. So no need to store this in a variable.
}
