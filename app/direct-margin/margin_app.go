package directmargin

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/channel/consensus_channel"
	"github.com/statechannels/go-nitro/client/engine"
	"github.com/statechannels/go-nitro/client/engine/store"
	"github.com/statechannels/go-nitro/internal"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

// TODO: Extract common errors into a common package
var ErrInvalidRequestType = internal.NewError("invalid request type")

type Balance struct {
	Remaining *big.Int
	Paid      *big.Int
}

type MarginApp struct {
	balances  map[string]*Balance
	store     store.Store
	engine    *engine.Engine
	myAddress common.Address
}

func NewMarginApp(engine *engine.Engine, myAddr common.Address, store store.Store) *MarginApp {
	return &MarginApp{
		balances:  make(map[string]*Balance),
		store:     store,
		engine:    engine,
		myAddress: myAddr,
	}
}

func (a *MarginApp) Id() string {
	return "direct-margin"
}

func (a *MarginApp) handleMarginProposal(
	ch *consensus_channel.ConsensusChannel,
	from types.Address,
	data interface{},
) {
	// 1. Check the proposal
	// 2. Accept
	// 2.1 Build the new state, and sign
	// 2.2 Call MarginAccept method with the new state

	// 3. Reject
	// 3.1 Call MarginReject method
	a.engine.SendMessages([]protocols.Message{})
}

func (a *MarginApp) HandleRequest(
	ch *consensus_channel.ConsensusChannel,
	from types.Address,
	ty string,
	data interface{},
) error {
	switch ty {
	case RequestTypeMarginProposal:
		a.handleMarginProposal(ch, from, data)

	case RequestTypeMarginAccept:

	case RequestTypeMarginReject:

	default:
		return ErrInvalidRequestType
	}

	return nil
}
