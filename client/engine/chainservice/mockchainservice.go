package chainservice

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

type MockChainService struct {
	recieveChan chan Event
	sendChan    chan protocols.Transaction
}

func NewMockChainService() ChainService {
	mcs := MockChainService{}
	mcs.recieveChan = make(chan Event)
	mcs.sendChan = make(chan protocols.Transaction)

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
		syntheticEvent := Event{
			ChannelId:          common.Hash{},
			Holdings:           types.Funds{},
			AdjudicationStatus: protocols.AdjudicationStatus{TurnNumRecord: 0},
		}
		if tx.Deposit.IsNonZero() {
			syntheticEvent.Holdings = tx.Deposit
		}
		mcs.recieveChan <- syntheticEvent
	}

}
