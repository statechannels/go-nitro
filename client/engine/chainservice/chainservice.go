// Package chainservice is a chain service responsible for submitting blockchain transactions and relaying blockchain events.
package chainservice // import "github.com/statechannels/go-nitro/client/chainservice"

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"

	NitroAdjudicator "github.com/statechannels/go-nitro/client/engine/chainservice/adjudicator"
)

// Event dictates which methods all chain events must implement
type Event interface {
	ChannelID() types.Destination
}

// CommonEvent declares fields shared by all chain events
type CommonEvent struct {
	channelID          types.Destination
	AdjudicationStatus protocols.AdjudicationStatus
	BlockNum           uint64
}

// DepositedEvent is an internal representation of the deposited blockchain event
type DepositedEvent struct {
	CommonEvent
	Holdings types.Funds // indexed by asset
}

func (de DepositedEvent) ChannelID() types.Destination {
	return de.channelID
}

// AllocationUpdated is an internal representation of the AllocatonUpdated blockchain event
type AllocationUpdatedEvent struct {
	CommonEvent
	Holdings types.Funds // indexed by asset
}

func (de AllocationUpdatedEvent) ChannelID() types.Destination {
	return de.channelID
}

// todo implement other event types
// Concluded
// ChallengeRegistered
// ChallengeCleared

// ChainEventHandler describes an objective that can handle chain events
type ChainEventHandler interface {
	UpdateWithChainEvent(event Event) (protocols.Objective, error)
}

type ChainService interface {
	Out() <-chan Event
	In() chan<- protocols.ChainTransaction
}

type EventProducer interface {
	SubscribeFilterLogs(ctx context.Context, query ethereum.FilterQuery, ch chan<- ethTypes.Log) (ethereum.Subscription, error)
}

type ChainConnection struct {
	in  chan protocols.ChainTransaction
	out chan Event
	na  *NitroAdjudicator.NitroAdjudicator
	to  *bind.TransactOpts
}

func NewChainConnection(na *NitroAdjudicator.NitroAdjudicator, naAddress common.Address, to *bind.TransactOpts, ep EventProducer) ChainConnection {
	cc := ChainConnection{}
	cc.in = make(chan protocols.ChainTransaction)
	cc.out = make(chan Event)
	cc.na = na
	cc.to = to

	go cc.listenForTx()
	go cc.listenForEvents(na, naAddress, to, ep)

	return cc
}

func (cc ChainConnection) listenForTx() {
	for tx := range cc.in {
		cc.handleTx(tx)
	}
}

// handleTx responds to the given tx.
func (cc ChainConnection) handleTx(tx protocols.ChainTransaction) {
	switch tx.Type {
	case protocols.DepositTransactionType:
		for address, amount := range tx.Deposit {
			cc.to.Value = amount
			_, err := cc.na.Deposit(cc.to, address, tx.ChannelId, big.NewInt(0), amount)
			if err != nil {
				panic(err)
			}
		}
	default:
		panic("unexpected chain transaction")
	}
}

func (cc ChainConnection) listenForEvents(na *NitroAdjudicator.NitroAdjudicator, naAddress common.Address, to *bind.TransactOpts, ep EventProducer) {
	query := ethereum.FilterQuery{
		Addresses: []common.Address{naAddress},
	}
	logs := make(chan ethTypes.Log)
	sub, err := ep.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatal(err)
	}
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLog := <-logs:
			fmt.Println(vLog) // pointer to event log

			// send dummy event
			event := DepositedEvent{}
			cc.out <- event
		}
	}
}

func (cc ChainConnection) In() chan<- protocols.ChainTransaction {
	return cc.in
}

func (cc ChainConnection) Out(a types.Address) <-chan Event {
	return cc.out
}
