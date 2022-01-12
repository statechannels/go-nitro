package chainservice

import (
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

type MockChainService struct {
	recieveChan chan Event                 // receiveChan is the chan used to send Events to the engine
	sendChan    chan protocols.Transaction // sendChan is the chan used to recieve Transactions from the engine

	holdings map[types.Destination]types.Funds
}

func NewMockChainService() ChainService {
	mcs := MockChainService{}
	mcs.recieveChan = make(chan Event)
	mcs.sendChan = make(chan protocols.Transaction)

	mcs.holdings = make(map[types.Destination]types.Funds)

	go mcs.ListenForTransactions()
	return mcs
}

// GetReceiveChan returns the recieveChan but narrows the type so that consumers mays only recieve on it.
func (mcs MockChainService) GetReceiveChan() <-chan Event {
	return chan Event(mcs.recieveChan)
}

// GetSendChan returns the rsendChan but narrows the type so that consumers mays only send on it.
func (mcs MockChainService) GetSendChan() chan<- protocols.Transaction {
	return mcs.sendChan
}

func (mcs MockChainService) Submit(tx protocols.Transaction) {}

func (mcs MockChainService) ListenForTransactions() {
	for tx := range mcs.sendChan {
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
		mcs.recieveChan <- syntheticEvent
	}

}
