package main

import (
	"fmt"
	"time"
)

const (
	num_leaves          = 200
	ledger_leaf_balance = 1e9
	ledger_hub_balance  = 1e9
	leaf_propose_period = `10s`
	leaf_payment_period = `100ms`
	leaf_close_period   = `30s`
)

// Spin up a single "hub" engine
// Spin up N "leaf" engines
// (skip) Each leaf creates a funded (1 gwei, 1 gwei) ledger channel with the hub
// Every 10 seconds (?) a leaf will propose a virtual channel with one of the other leaves
// Every 100ms a leaf will send a 1 wei payment in each of the virtual channels it proposed
// Every 30 seconds a leaf will close an opn virtual channel
// Metrics recorded: balances of each ledger channel (over time, or just once at the start and at the end)

// following https://gobyexample.com/stateful-goroutines

type LedgerChannelProposal struct {
	leaf int // id
}

func hub(id int, ledgerChannelProposals chan LedgerChannelProposal) {
	var ledgerChannels = make(map[int]bool) // the state of the hub is a mapping from leaf id to a bool (connected or not?)
	for {
		select {
		case proposal := <-ledgerChannelProposals:
			ledgerChannels[proposal.leaf] = true
			fmt.Println(`hub connected to`, proposal.leaf)
		}
	}
}

func leaf(id int, ledgerChannelProposals chan LedgerChannelProposal) {
	// var virtualChannels = make(map[int]bool)
	ledgerChannelProposal := LedgerChannelProposal{id}
	ledgerChannelProposals <- ledgerChannelProposal
}

func main() {
	fmt.Println(`Starting hub...`)
	proposals := make(chan LedgerChannelProposal)
	go hub(0, proposals)
	for l := 1; l < num_leaves+1; l++ {
		go leaf(l, proposals)
	}

	time.Sleep(time.Second)

}
