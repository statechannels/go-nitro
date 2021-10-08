package main

// HubRunner is a function which simulates an application running a NitroWallet in "hub" mode
func HubRunner(
	hubId uint, // Network id for this node
	ledgerInbox *map[uint]chan LedgerChannelState, // Map of inbox go channels for ledger channel updates
	paymentInbox *map[uint]chan VirtualChannelState, // Map of inbox go channels for vitual channel updates
) {
	NewNitroWallet(hubId, ledgerInbox, paymentInbox, true) // We don't need to hit the API just yet. So no need to store this in a variable.
}
