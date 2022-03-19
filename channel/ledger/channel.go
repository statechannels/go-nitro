package ledger

import (
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/types"
)

type LedgerChannel struct {
	Id      types.Destination
	MyIndex uint

	OnChainFunding types.Funds

	state.FixedPart

	consensus state.SignedState
}


type Proposal struct {
	turnNum uint64
	add     *struct {
		amount types.Funds
		vId    types.Destination
		left   types.Destination
		right  types.Destination
	}
	remove *struct {
		vId        types.Destination
		amountPaid types.Funds
	}
	signature state.Signature
}

