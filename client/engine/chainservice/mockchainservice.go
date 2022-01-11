package chainservice

import (
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

type MockChainService struct {
	recieveChan chan Event
	sendChan    chan protocols.Transaction

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

func (mcs MockChainService) GetReceiveChan() chan Event {
	return mcs.recieveChan
}
func (mcs MockChainService) GetSendChan() chan protocols.Transaction {
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
