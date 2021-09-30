package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	num_leaves                     = 200
	ledger_leaf_balance            = 1e9
	ledger_hub_balance             = 1e9
	virtual_proposer_balance       = 1e2
	virtual_joiner_balance         = 1e2
	leaf_propose_period            = time.Duration(1e9) // ns
	leaf_payment_period            = time.Duration(1e7) // ns
	leaf_close_period              = time.Duration(3e9) // ns
	inter_node_channel_buffer_size = 1e3
)

// Spin up a single "hub" engine
// Spin up N "leaf" engines
// (skip) Each leaf creates a funded (1 gwei, 1 gwei) ledger channel with the hub
// Every 1 seconds (?) a leaf will propose a virtual channel with one of the other leaves
// Every 100ms a leaf will send a 1 wei payment in each of the virtual channels it proposed
// Every 3 seconds a leaf will close an opn virtual channel
// Metrics recorded: balances of each ledger channel (over time, or just once at the start and at the end)

// following https://gobyexample.com/stateful-goroutines

func channelId(joiner uint, proposer uint) string {
	return fmt.Sprintf("%v-%v", joiner, proposer)
}

type LedgerStore struct {
	ledgerChannels map[uint]LedgerChannelState
}
type LedgerChannelState struct {
	hubId, hubBal, leafId, leafBal, turnNum uint            // hubs and leaves have integer ids
	virtualChannelBal                       map[string]uint // channels have 2d uint ids [joiner][proposer]
	signedByHub                             bool
	signedByLeaf                            bool
}

type VirtualChannelState struct {
	proposerId, joinerId, proposerBal, joinerBal, turnNum uint
}

type LedgerRequest struct {
	virtualChannelProposer, virtualChannelJoiner uint
	amount                                       int // note this is a *signed* quantity
	sucess                                       chan bool
}

func hub(ledgerUpdatesToHub <-chan LedgerChannelState, ledgerUpdatesToLeaves map[uint]chan LedgerChannelState) {
	const hubId = 0
	var store = LedgerStore{make(map[uint]LedgerChannelState)}

	// Manual setup for ledger channels
	for id := uint(1); id < num_leaves+1; id++ {
		store.ledgerChannels[id] = LedgerChannelState{hubId, ledger_hub_balance, id, ledger_leaf_balance, 1, make(map[string]uint), true, true} // turnNum = 1 so channels are funded
	}

	// Listen for ledgerUpdates. Countersign blindly! TODO
	for {
		select {
		case update := <-ledgerUpdatesToHub:
			update.signedByHub = true
			store.ledgerChannels[update.leafId] = update
			fmt.Printf(`ledger between hub and leaf %v updated to %v`, update.leafId, update)
			fmt.Println()
			// send back
			ledgerUpdatesToLeaves[update.leafId] <- update // this will block unless there is something receiving on this channel or if the channel is buffered.
		}
	}
}

func leaf(id uint, ledgerUpdatesToHub chan<- LedgerChannelState, ledgerUpdatesToMe <-chan LedgerChannelState) {
	const hubId = 0
	var ledgerChannel = LedgerChannelState{hubId, ledger_hub_balance, id, ledger_leaf_balance, 1, make(map[string]uint), false, true}

	proposeAVirtualChannel := func() {
		// TODO communicate with peer
		randomPeer := uint(rand.Intn(num_leaves)) // TODO check we don't already have a channel open
		cId := channelId(id, randomPeer)
		// read: from the store
		ledger := ledgerChannel
		virtualChannelBal := ledger.virtualChannelBal
		// modify: reallocate 10 to a virtual channel with a single peer
		virtualChannelBal[cId] += 10
		ledger.hubBal -= 5
		ledger.leafBal -= 5
		ledger.virtualChannelBal = virtualChannelBal
		// store: in my store
		ledgerChannel = ledger
		// send: to hub for signing
		ledgerUpdatesToHub <- ledger
	}
	// updateAVirtualChannel := func() {}
	// closeAVirtualChannel := func() {}

	proposeAVirtualChannel()

	// var virtualChannels = make(map[int]bool)

}

func main() {
	fmt.Println(`Starting hub...`)
	//  The coordinator makes buffered channels. This means engines can send or receive without blocking until the buffer is saturated. This is because we are simulating processes running on different machines? Does that make sense?
	ledgerUpdatesToHub := make(chan LedgerChannelState, inter_node_channel_buffer_size) // the coordinator creates a channel for the leaves to communicate with the hub
	ledgerUpdatesToLeaves := make(map[uint]chan LedgerChannelState)                     // the coordinator makes enough channels for the hub to communicate with each leaf
	for l := uint(1); l < num_leaves+1; l++ {
		ledgerUpdatesToLeaves[l] = make(chan LedgerChannelState, inter_node_channel_buffer_size)
	}
	// the coordinator makes enough channels for the leaves to be fully connected, and store these in a mapping keyed by proposer and joiner
	go hub(ledgerUpdatesToHub, ledgerUpdatesToLeaves)
	fmt.Println(`Starting leaves...`)
	for l := uint(1); l < num_leaves+1; l++ {
		go leaf(l, ledgerUpdatesToHub, ledgerUpdatesToLeaves[l])
	}

	time.Sleep(time.Second * 1)

}
