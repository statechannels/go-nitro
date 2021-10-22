// Package chainservice is a chain service responsible for submitting blockchain transactions and relaying blockchain events
package chainservice // import "github.com/statechannels/go-nitro/client/chainservice"

import (
	"math/big"

	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

// ChainEvent is an internal representation of a blockchain event
type Event struct {
	ChannelId          types.Bytes32
	Holdings           map[types.Address]big.Int // indexed by asset
	AdjudicationStatus protocols.AdjudicationStatus
}

type ChainService interface {
	GetRecieveChan() chan Event
	GetSendChan() chan protocols.Transaction
	Submit(tx protocols.Transaction)
}

type TestChainService struct{}

var recieveChan chan Event = make(chan Event)
var sendChan chan protocols.Transaction = make(chan protocols.Transaction)

func (TestChainService) GetRecieveChan() chan Event              { return recieveChan }
func (TestChainService) GetSendChan() chan protocols.Transaction { return sendChan }
func (TestChainService) Submit(tx protocols.Transaction)         {}
