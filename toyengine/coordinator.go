package main

import (
	"fmt"
	"time"
)

const (
	num_leaves                     = 200                // The number of non-hub nodes in the state channel network
	ledger_leaf_balance            = 1e9                // The balance of the leaf in each ledger channel
	ledger_hub_balance             = 1e9                // The balance of the hub in each ledger channel
	virtual_proposer_balance       = 1e2                // The balance of the proposer in each virtual channel
	virtual_joiner_balance         = 1e2                // The balance of the joiner in each virtual channel
	leaf_propose_period            = time.Duration(1e9) // (ns) period between "ticks" that trigger a new virtual channel proposal
	leaf_payment_period            = time.Duration(1e7) // (ns) period between "ticks" that trigger a virtual channel payment
	leaf_close_period              = time.Duration(3e9) // (ns) period between "ticks" that trigger the closing of a virtual channel
	inter_node_channel_buffer_size = 1e3                // The number of messages buggered in each go channel linking nodes for messaging
)

// Example coordination:
// Spin up a single "hub" engine
// Spin up N "leaf" engines
// (skip) Each leaf creates a funded (1 gwei, 1 gwei) ledger channel with the hub
// Every 1 seconds (?) a leaf will propose a virtual channel with one of the other leaves
// Every 100ms a leaf will send a 1 wei payment in each of the virtual channels it proposed
// Every 3 seconds a leaf will close an opn virtual channel
// Metrics recorded: balances of each ledger channel (over time, or just once at the start and at the end)

// following https://gobyexample.com/stateful-goroutines

func main() {

	fmt.Println(`Setting up communciation channels...`)

	ledgerInbox := make(map[uint]chan LedgerChannelState)
	for l := uint(0); l < num_leaves+1; l++ {
		// anyone can send a message to anyone. wallets should only listen on their own channel, and never send on it
		ledgerInbox[l] = make(chan LedgerChannelState, inter_node_channel_buffer_size)
	}

	fmt.Println(`Starting hub runner...`)
	go HubRunner(0, &ledgerInbox)

	fmt.Println(`Starting leaf runners...`)
	for l := uint(1); l < num_leaves+1; l++ {
		go LeafRunner(l, &ledgerInbox)
	}

	fmt.Println(`Letting everything run for a bit...`)
	time.Sleep(time.Second * 5)

}
