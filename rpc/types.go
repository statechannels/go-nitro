package rpc

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/statechannels/go-nitro/node/query"
	"github.com/statechannels/go-nitro/types"
)

type LedgerChannelView struct {
	ID           types.Destination
	Counterparty types.Address
	AssetAddress types.Address
	MyBalance    *hexutil.Big
	TheirBalance *hexutil.Big
	// todo: running guarantees?
}

// NewLedgerChannelView creates a "view with perspective" into a given ledger channel.
func NewLedgerChannelView(info query.LedgerChannelInfo, myAddress types.Address) (LedgerChannelView, error) {
	lcv := LedgerChannelView{}

	lcv.ID = info.ID
	lcv.AssetAddress = info.Balance.AssetAddress

	if info.Balance.Me == myAddress {
		lcv.Counterparty = info.Balance.Them
		lcv.MyBalance = info.Balance.MyBalance
		lcv.TheirBalance = info.Balance.TheirBalance
	} else if info.Balance.Them == myAddress {
		lcv.Counterparty = info.Balance.Me
		lcv.MyBalance = info.Balance.TheirBalance
		lcv.TheirBalance = info.Balance.MyBalance
	} else {
		return lcv, fmt.Errorf("%s is not a participant in ledger channel %s", myAddress, info.ID)
	}

	return lcv, nil
}
