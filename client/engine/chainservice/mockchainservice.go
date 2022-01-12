package chainservice

import (
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

type MockChainService struct {
	out chan Event                 // out is the chan used to send Events to the engine
	in  chan protocols.Transaction // in is the chan used to recieve Transactions from the engine

	holdings map[types.Destination]types.Funds
}

func NewMockChainService() ChainService {
	mcs := MockChainService{}
	mcs.out = make(chan Event)
	mcs.in = make(chan protocols.Transaction)

	mcs.holdings = make(map[types.Destination]types.Funds)

	go mcs.ListenForTransactions()
	return mcs
}

// Out() returns the but chan but narrows the type so that external consumers mays only recieve on it.
func (mcs MockChainService) Out() <-chan Event {
	return chan Event(mcs.out)
}

// In returns the in chan but narrows the type so that external consumers mays only send on it.
func (mcs MockChainService) In() chan<- protocols.Transaction {
	return mcs.in
}

func (mcs MockChainService) ListenForTransactions() {
	for tx := range mcs.in {
		channelId := tx.ChannelId
		syntheticEvent := Event{
			ChannelId:          channelId,
			Holdings:           types.Funds{},
			AdjudicationStatus: protocols.AdjudicationStatus{TurnNumRecord: 0},
		}
		if tx.Deposit.IsNonZero() {
			mcs.holdings[channelId] = mcs.holdings[channelId].Add(tx.Deposit)
			syntheticEvent.Holdings = mcs.holdings[channelId]
		}
		mcs.out <- syntheticEvent
	}

}
